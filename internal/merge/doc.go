// Package merge provides utilities for combining multiple .env files into a
// single unified map.
//
// It supports two conflict resolution strategies:
//
//   - StrategyFirst: the first file to define a key wins.
//   - StrategyLast:  the last file to define a key wins (override semantics).
//
// Conflicts — keys defined in more than one source — are collected and
// returned in Result.Conflicts so callers can surface them to the user.
package merge
