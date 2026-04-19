// Package audit records command invocations and their outcomes to a
// persistent JSON log file. Each run appends an Entry containing the
// timestamp, command name, files involved, and any detected changes or
// issues. The log can be inspected later for compliance or debugging.
//
// Usage:
//
//	err := audit.Append(".envcmp-audit.json", audit.Entry{
//		Command: "diff",
//		Files:   []string{"staging.env", "prod.env"},
//	})
package audit
