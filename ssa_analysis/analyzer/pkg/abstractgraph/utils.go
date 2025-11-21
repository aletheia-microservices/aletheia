package abstractgraph

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Examples:
// - t3
// - t4.t7
// - t4.t8
// - t4.t8.t9
// Explanation:
// in t4.t7, t4 is the variable in the exposed service method, 
// and t7 is the variable in the internal called method
func parseT(t string) []int {
	splits := strings.Split(t, ".")
	res := make([]int, 0, len(splits))

	for _, sub_t := range splits {
		if strings.HasPrefix(sub_t, "t") {
			sub_t = sub_t[1:]
		} else {
			logrus.Fatalf("unexpected absence of 't' prefix in string (t=%s) (sub_t=%s)\n", t, sub_t)
		}
		n, err := strconv.Atoi(sub_t)
		if err != nil {
			logrus.Fatalf("could not extract number from string (t=%s) (sub_t=%s)\n", t, sub_t)
		}
		res = append(res, n)
	}
	return res
}

func lessT(t1 string, t2 string) bool {
	t1_lst := parseT(t1)
	t2_lst := parseT(t2)

	for i := 0; i < len(t1_lst) && i < len(t2_lst); i++ {
		if t1_lst[i] < t2_lst[i] {
			return true
		}
		if t1_lst[i] > t2_lst[i] {
			return false
		}
	}
	// if all equal up to min length => shorter one comes first
	// e.g., t4.t8 < t4.t8.t9
	return len(t1_lst) < len(t2_lst)
}

func greaterT(t1 string, t2 string) bool {
	return lessT(t2, t1)
}

func equalT(t1 string, t2 string) bool {
	return t1 == t2
}
