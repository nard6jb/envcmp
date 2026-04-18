// Package diff compares two parsed env maps and returns a list of differences.
package diff

// Status describes the kind of difference found for a key.
type Status string

const (
	// Missing means the key exists in the left env but not the right.
	Missing Status = "missing"
	// Extra means the key exists in the right env but not the left.
	Extra Status = "extra"
	// Changed means the key exists in both but with different values.
	Changed Status = "changed"
)

// Result holds the diff outcome for a single key.
type Result struct {
	Key        string
	Status     Status
	LeftValue  string
	RightValue string
}

// Compare returns the differences between left and right env maps.
func Compare(left, right map[string]string) []Result {
	var results []Result

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			results = append(results, Result{Key: k, Status: Missing, LeftValue: lv})
		} else if lv != rv {
			results = append(results, Result{Key: k, Status: Changed, LeftValue: lv, RightValue: rv})
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			results = append(results, Result{Key: k, Status: Extra, RightValue: rv})
		}
	}

	return results
}

// HasDiff returns true if any differences exist.
func HasDiff(results []Result) bool {
	return len(results) > 0
}
