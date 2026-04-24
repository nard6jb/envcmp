package env

// RenameOptions controls how keys are renamed in an env map.
type RenameOptions struct {
	// Mapping is a map of old key name -> new key name.
	Mapping map[string]string
	// DropUnmapped removes keys not present in Mapping when true.
	DropUnmapped bool
}

// RenameResult holds the outcome of a Rename operation.
type RenameResult struct {
	// Out is the resulting env map after renaming.
	Out map[string]string
	// Renamed lists keys that were successfully renamed (old -> new).
	Renamed map[string]string
	// Dropped lists keys that were removed because DropUnmapped was set.
	Dropped []string
	// Skipped lists keys whose target name already existed in the map.
	Skipped []string
}

// Rename applies key renaming to src according to opts.
// It never mutates src.
func Rename(src map[string]string, opts RenameOptions) RenameResult {
	out := make(map[string]string, len(src))
	result := RenameResult{
		Renamed: make(map[string]string),
	}

	// First pass: copy unmapped keys (or drop them).
	for k, v := range src {
		if _, mapped := opts.Mapping[k]; mapped {
			continue
		}
		if opts.DropUnmapped {
			result.Dropped = append(result.Dropped, k)
		} else {
			out[k] = v
		}
	}

	// Second pass: apply renames.
	for oldKey, newKey := range opts.Mapping {
		v, exists := src[oldKey]
		if !exists {
			continue
		}
		if _, conflict := out[newKey]; conflict {
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}
		out[newKey] = v
		result.Renamed[oldKey] = newKey
	}

	result.Out = out
	return result
}
