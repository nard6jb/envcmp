// Package filter provides utilities to narrow down env key-value maps
// before diffing or validation. Filtering can be done by key prefix,
// an explicit allow-list of keys, or an exclusion list.
//
// Typical usage:
//
//	filtered := filter.Apply(envMap, filter.Options{
//		Prefix:  "APP_",
//		Exclude: []string{"APP_DEBUG"},
//	})
package filter
