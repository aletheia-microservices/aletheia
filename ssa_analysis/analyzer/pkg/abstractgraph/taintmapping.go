package abstractgraph

import (
	"fmt"
	"slices"
	"sort"
	"strings"
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
	if key == valElem {
		return
	}
	fmt.Printf("[TM] adding taint mapping (%s) -> (%s)\n", key.String(), valElem.String())
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
