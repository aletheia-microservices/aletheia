package abstractgraph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type TaintPair struct {
	Taint1 AbstractTaint
	Taint2 AbstractTaint
}

// memoization of previously computed taint pairs to avoid propagating over again
// it is resetted on NewTaintMapping which occurs when visiting a node or matching queue push/pops
var seenTaintPairsOnJoin = make(map[TaintPair]bool)

// TaintMapping is a mapping between primary taint (key) that already existed
// and the list of secondary taints (value) that were recently propagated
//
// in practice, this means that all objects with primary taint will also inherit
// the secondary taint and possibly originate a foreign key when written to the databases
type TaintMapping struct {
	mapping     map[AbstractTaint][]AbstractTaint // do not use pointers so that we can compare easily
	mappingKeys []AbstractTaint                   // to track order of keys in "mapping" above
}

func NewTaintMapping() *TaintMapping {
	seenTaintPairsOnJoin = make(map[TaintPair]bool)
	return &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
}

func (tm *TaintMapping) Clear() {
	tm = &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
}

func (tm *TaintMapping) GetMappingKeys() []AbstractTaint {
	return tm.mappingKeys
}

func (tm *TaintMapping) GetMappingForKey(key AbstractTaint) []AbstractTaint {
	return tm.mapping[key]
}

func (tm *TaintMapping) AddIfNotExists(key AbstractTaint, valElem AbstractTaint, after bool, join bool) {
	if key == valElem {
		return
	}

	if _, exists := seenTaintPairsOnJoin[TaintPair{Taint1: key, Taint2: valElem}]; exists {
		return
	}

	logrus.Tracef("[TM] adding taint mapping (%s) -> (%s)\n", key.String(), valElem.String())
	if mappingVal, ok := tm.mapping[key]; ok {
		var exists bool
		for _, t := range mappingVal {
			if t.EqualExceptPrimaryAndTrace(&valElem) {
				exists = true
				break
			}
		}
		if !exists {
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

	if join {
		seenTaintPairsOnJoin[TaintPair{Taint1: key, Taint2: valElem}] = true
	}
}

func (tm *TaintMapping) Join(other *TaintMapping, after bool) {
	for _, otherKey := range other.GetMappingKeys() {
		otherValLst := other.GetMappingForKey(otherKey)
		for _, otherValElem := range otherValLst {
			tm.AddIfNotExists(otherKey, otherValElem, after, true)
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
		var num int
		var valsStr []string
		for _, val := range tm.mapping[key] {
			valsStr = append(valsStr, val.LongLongString())
			num++
		}

		sort.Strings(valsStr)

		builder.WriteString(fmt.Sprintf("  %s: [%d]\n", key.GetDatabasePath(), num))
	}

	builder.WriteString("}")
	return builder.String()
}
