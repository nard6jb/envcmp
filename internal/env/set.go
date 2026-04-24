package env

import "fmt"

// SetOption configures the behaviour of the Set operation.
type SetOption func(*setOptions)

type setOptions struct {
	overwrite bool
	validateKey bool
}

// WithOverwrite allows Set to replace an existing value.
// By default, Set returns an error if the key already exists.
func WithOverwrite() SetOption {
	return func(o *setOptions) {
		o.overwrite = true
	}
}

// WithKeyValidation enables basic key validation (non-empty, no whitespace).
func WithKeyValidation() SetOption {
	return func(o *setOptions) {
		o.validateKey = true
	}
}

// Set adds or updates a single key/value pair in the provided map.
// It returns an error if the key already exists and WithOverwrite is not set,
// or if key validation is enabled and the key is invalid.
//
// The source map is modified in place; a new map is also returned for
// convenient chaining.
func Set(env map[string]string, key, value string, opts ...SetOption) (map[string]string, error) {
	cfg := &setOptions{}
	for _, o := range opts {
		o(cfg)
	}

	if cfg.validateKey {
		if err := validateEnvKey(key); err != nil {
			return env, err
		}
	}

	if env == nil {
		env = make(map[string]string)
	}

	if _, exists := env[key]; exists && !cfg.overwrite {
		return env, fmt.Errorf("key %q already exists; use WithOverwrite to replace it", key)
	}

	env[key] = value
	return env, nil
}

// SetMany applies a batch of key/value pairs to env using the same options.
// Processing stops at the first error.
func SetMany(env map[string]string, pairs map[string]string, opts ...SetOption) (map[string]string, error) {
	if env == nil {
		env = make(map[string]string)
	}
	for k, v := range pairs {
		var err error
		env, err = Set(env, k, v, opts...)
		if err != nil {
			return env, err
		}
	}
	return env, nil
}

// validateEnvKey returns an error for keys that would be invalid in a .env file.
func validateEnvKey(key string) error {
	if key == "" {
		return fmt.Errorf("env key must not be empty")
	}
	for _, r := range key {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return fmt.Errorf("env key %q must not contain whitespace", key)
		}
	}
	return nil
}
