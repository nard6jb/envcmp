package env

import "fmt"

// PatchOp represents a single patch operation on an env map.
type PatchOp struct {
	Key    string
	Value  string
	Delete bool
}

// PatchResult holds the outcome of applying a patch operation.
type PatchResult struct {
	Key      string
	OldValue string
	NewValue string
	Deleted  bool
	Added    bool
	Changed  bool
	Skipped  bool
}

// PatchOptions controls patch behaviour.
type PatchOptions struct {
	AllowNew       bool // allow adding keys not present in base
	AllowDelete    bool // allow deleting keys
	ErrorOnMissing bool // error if a non-delete op targets a missing key
}

// Patch applies a list of PatchOps to a copy of base and returns the
// patched map along with a per-operation result slice.
func Patch(base map[string]string, ops []PatchOp, opts PatchOptions) (map[string]string, []PatchResult, error) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	results := make([]PatchResult, 0, len(ops))

	for _, op := range ops {
		old, exists := out[op.Key]

		if op.Delete {
			if !opts.AllowDelete {
				return nil, nil, fmt.Errorf("patch: delete not allowed for key %q", op.Key)
			}
			if !exists {
				results = append(results, PatchResult{Key: op.Key, Skipped: true})
				continue
			}
			delete(out, op.Key)
			results = append(results, PatchResult{Key: op.Key, OldValue: old, Deleted: true})
			continue
		}

		if !exists {
			if !opts.AllowNew {
				if opts.ErrorOnMissing {
					return nil, nil, fmt.Errorf("patch: key %q not found in base", op.Key)
				}
				results = append(results, PatchResult{Key: op.Key, Skipped: true})
				continue
			}
			out[op.Key] = op.Value
			results = append(results, PatchResult{Key: op.Key, NewValue: op.Value, Added: true})
			continue
		}

		out[op.Key] = op.Value
		results = append(results, PatchResult{Key: op.Key, OldValue: old, NewValue: op.Value, Changed: old != op.Value})
	}

	return out, results, nil
}
