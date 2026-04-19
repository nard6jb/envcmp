package audit

import (
	"fmt"
	"io"
	"strings"
)

// Print writes a human-readable summary of the log to w.
func Print(w io.Writer, l *Log, color bool) {
	if len(l.Entries) == 0 {
		fmt.Fprintln(w, "No audit entries found.")
		return
	}
	for i, e := range l.Entries {
		ts := e.Timestamp.Format("2006-01-02 15:04:05 UTC")
		files := strings.Join(e.Files, ", ")
		if color {
			fmt.Fprintf(w, "\033[1m[%d]\033[0m %s — \033[36m%s\033[0m (%s)\n", i+1, ts, e.Command, files)
		} else {
			fmt.Fprintf(w, "[%d] %s — %s (%s)\n", i+1, ts, e.Command, files)
		}
		for k, v := range e.Changes {
			fmt.Fprintf(w, "    ~ %s: %s\n", k, v)
		}
		for _, iss := range e.Issues {
			if color {
				fmt.Fprintf(w, "    \033[33m! %s\033[0m\n", iss)
			} else {
				fmt.Fprintf(w, "    ! %s\n", iss)
			}
		}
	}
}
