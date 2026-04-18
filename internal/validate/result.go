package validate

// Result holds the outcome of a validation run.
type Result struct {
	// Valid is true when the target env satisfies all reference keys
	// and contains no unexpected extra keys.
	Valid bool

	// MissingKeys are keys present in the reference but absent in the target.
	MissingKeys []string

	// ExtraKeys are keys present in the target but absent in the reference.
	ExtraKeys []string
}
