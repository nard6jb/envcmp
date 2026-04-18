// Package diff provides functionality to compare two parsed .env files
// and report missing, extra, and changed keys between them.
package diff

// Result holds the diff outcome between two env maps.
type Result struct {
	MissingInRight []string            // keys present in left but not in right
	MissingInLeft  []string            // keys present in right but not in left
	Changed        map[string][2]string // key -> [leftVal, rightVal]
}

// Compare takes two env maps (key->value) and returns a Result describing
// their differences.
func Compare(left, right map[string]string) Result {
	res := Result{
		Changed: make(map[string][2]string),
	}

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			res.MissingInRight = append(res.MissingInRight, k)
			continue
		}
		if lv != rv {
			res.Changed[k] = [2]string{lv, rv}
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			res.MissingInLeft = append(res.MissingInLeft, k)
		}
	}

	return res
}

// HasDiff returns true if the Result contains any differences.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 || len(r.MissingInLeft) > 0 || len(r.Changed) > 0
}
