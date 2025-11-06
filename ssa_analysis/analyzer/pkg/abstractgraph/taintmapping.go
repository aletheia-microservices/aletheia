package abstractgraph

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"analyzer/pkg/common"
	"analyzer/pkg/utils"
)

// TaintMapping is a mapping between primary taint (key) that already existed
// and the list of secondary taints (value) that were recently propagated
//
// in practice, this means that all objects with primary taint will also inherit
// the secondary taint and possibly originate a foreign key when written to the database
type TaintMapping struct {
	mapping     map[AbstractTaint][]AbstractTaint // do not use pointers so that we can compare easily
	mappingKeys []AbstractTaint                   // to track order of keys in "mapping" above
}

func NewTaintMapping() *TaintMapping {
	return &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
}

func (tm *TaintMapping) GetMappingKeys() []AbstractTaint {
	return tm.mappingKeys
}

func (tm *TaintMapping) GetMappingForKey(key AbstractTaint) []AbstractTaint {
	return tm.mapping[key]
}

func (tm *TaintMapping) AddIfNotExists(key AbstractTaint, valElem AbstractTaint, after bool) {
	if mappingVal, ok := tm.mapping[key]; ok {
		if !slices.Contains(mappingVal, valElem) {
			if after {
				tm.mapping[key] = append(mappingVal, valElem)
			} else {
				tm.mapping[key] = append([]AbstractTaint{valElem}, mappingVal...)
			}
		}
	} else {
		if after {
			tm.mappingKeys = append(tm.mappingKeys, key)
		} else {
			tm.mappingKeys = append([]AbstractTaint{key}, tm.mappingKeys...)
		}
		tm.mapping[key] = []AbstractTaint{valElem}
	}
}

func (tm *TaintMapping) Merge(other *TaintMapping, after bool) {
	for _, otherKey := range other.GetMappingKeys() {
		otherValLst := other.GetMappingForKey(otherKey)
		for _, otherValElem := range otherValLst {
			tm.AddIfNotExists(otherKey, otherValElem, after)
		}
	}
}

func (tm *TaintMapping) String() string {
	if len(tm.mapping) == 0 {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{\n")

	keys := make([]AbstractTaint, 0, len(tm.mapping))
	for k := range tm.mapping {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].GetDatabasePath() < keys[j].GetDatabasePath()
	})

	for _, key := range keys {
		var valsStr []string
		for _, val := range tm.mapping[key] {
			valsStr = append(valsStr, val.GetDatabasePath())
		}

		sort.Strings(valsStr)

		builder.WriteString(fmt.Sprintf("  %s: [%s]\n", key.GetDatabasePath(), strings.Join(valsStr, ", ")))
	}

	builder.WriteString("}")
	return builder.String()
}

/* func ComputeTaintMapping(obj *AbstractObject) *TaintMapping {
	taintMapping := &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}

	for objpath1, primaryTaintsLst := range obj.GetPrimaryTaints() {
		for objpath2, secondaryTaintsLst := range obj.GetSecondaryTaints() {
			if objpath1 == objpath2 {
				for _, primaryTaint := range primaryTaintsLst {
					taintMapping
				}
			}
		}
	}
	return taintMapping
} */

func MergeTaints(obj *AbstractObject, otherTaintsMap map[string][]*AbstractTaint, primary bool, traced bool) *TaintMapping {
	fmt.Printf("[TAINTMAPPING] merging taints (primary=%t, traced=%t): %v\n", primary, traced, otherTaintsMap)
	var taintMapping *TaintMapping

	if !primary {
		taintMapping = &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
	}

	for objpath, otherTaints := range otherTaintsMap {
		fmt.Printf("[TAINTMAPPING] checking existing taints for objpath (%s)\n", objpath)
		existingTaints := obj.taints[objpath]

		exists := func(otherTaint *AbstractTaint) (string, bool) {
			for _, existingTaint := range existingTaints {
				if existingTaint.Equals(otherTaint) {
					return objpath, true
				}
				fmt.Printf("[TAINTMAPPING] checking if upper path (%s) vs (%s)\n", existingTaint.fieldpath, otherTaint.fieldpath)
				if ok, subpath := existingTaint.IsUpperPath(otherTaint); ok {
					return objpath + subpath, false
				}
			}
			return objpath, false
		}

		for _, otherTaint := range otherTaints {
			if objpath, equal := exists(otherTaint); !equal {
				// need to create new AbstractTaint to avoid just
				// storing the pointer and modifying its fields
				newTaint := NewAbstractTaint(
					otherTaint.fieldpath,
					otherTaint.dbcallID,
					otherTaint.dbOpType,
					primary, traced,
				)

				fmt.Printf("\t[TAINTMAPPING] [OBJ={%s}] adding new taint (%s, traced=%t) on obj path (%s): %v\n", obj.String(), common.OperationTypeToString(newTaint.dbOpType), newTaint.traced, objpath, newTaint)
				obj.taints[objpath] = append(obj.taints[objpath], newTaint)

				// NOTE: explores all upper paths
				//
				// trace info for arguments and (especially) returns
				// can still lead to secondary taints that we still want to track
				if !primary {
					fmt.Printf("\t[TAINTMAPPING] [OBJ={%s}] adding mapping for objpath={%s} // taint={%s} // traced={%t}\n", obj.String(), objpath, newTaint.LongString(), traced)
					var ok = true
					var subpath = ""
					for ok {
						for _, existingTaint := range obj.taints[objpath] {
							// filter by writes to reduce number of foreign keys for now
							if existingTaint.IsPrimary() && !traced {
								if subpath == "" {
									taintMapping.AddIfNotExists(*existingTaint, *newTaint, true)
									fmt.Printf("\t\t[TAINTMAPPING] [OBJ={%s}] [1] upperpath={%s} // subpath={%s} // existingTaint={%s} // traced={%t}\n", obj.String(), objpath, subpath, existingTaint.LongString(), traced)
								} else {
									lowerTaint := existingTaint.Copy()
									// TODO: verify if this is needed (can't recall now)
									lowerTaint.AddSuffixToDatabasePath(subpath)
									taintMapping.AddIfNotExists(*lowerTaint, *newTaint, true)
								}
							} else if traced {
								fmt.Printf("\t\t[TAINTMAPPING] [OBJ={%s}] [3] upperpath={%s} // subpath={%s} // existingTaint={%s} // traced={%t}\n", obj.String(), objpath, subpath, existingTaint.LongString(), traced)
								if subpath == "" {
									taintMapping.AddIfNotExists(*newTaint, *existingTaint, true)
								} else {
									lowerTaint := *existingTaint
									lowerTaint.fieldpath = lowerTaint.fieldpath + subpath
									// [TO BE IMPROVED]
									// for some reason it works better when we change the
									// position between newTaint and lowerTaint in call args
									// e.g., SockShop3: order_db.orders.ID REFERENCES ship_db.shipments.ID
									//
									// i think this is because of the order
									// when tainting primary vs. traced
									taintMapping.AddIfNotExists(*newTaint, lowerTaint, true)
								}

							}
						}
						objpath, subpath, ok = utils.ExtractUpperPath(objpath)
					}
				}
			}
		}
	}
	return taintMapping
}
