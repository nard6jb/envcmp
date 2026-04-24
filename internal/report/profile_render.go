package report

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envcmp/internal/profile"
)

// RenderProfiles prints all profiles in the store to stdout.
func RenderProfiles(s *profile.Store, color bool) {
	RenderProfilesTo(os.Stdout, s, color)
}

// RenderProfilesTo writes profile listings to the provided writer.
func RenderProfilesTo(w io.Writer, s *profile.Store, color bool) {
	if len(s.Profiles) == 0 {
		if color {
			fmt.Fprintln(w, "\033[33mNo profiles defined.\033[0m")
		} else {
			fmt.Fprintln(w, "No profiles defined.")
		}
		return
	}

	profiles := make([]profile.Profile, len(s.Profiles))
	copy(profiles, s.Profiles)
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	for _, p := range profiles {
		if color {
			fmt.Fprintf(w, "\033[36m%-20s\033[0m %s\n", p.Name, strings.Join(p.Files, ", "))
		} else {
			fmt.Fprintf(w, "%-20s %s\n", p.Name, strings.Join(p.Files, ", "))
		}
	}
}

// RenderProfileAdded prints a confirmation message after adding a profile.
func RenderProfileAdded(w io.Writer, name string, color bool) {
	if color {
		fmt.Fprintf(w, "\033[32mProfile %q saved.\033[0m\n", name)
	} else {
		fmt.Fprintf(w, "Profile %q saved.\n", name)
	}
}

// RenderProfileRemoved prints a confirmation message after removing a profile.
func RenderProfileRemoved(w io.Writer, name string, found bool, color bool) {
	if !found {
		if color {
			fmt.Fprintf(w, "\033[31mProfile %q not found.\033[0m\n", name)
		} else {
			fmt.Fprintf(w, "Profile %q not found.\n", name)
		}
		return
	}
	if color {
		fmt.Fprintf(w, "\033[33mProfile %q removed.\033[0m\n", name)
	} else {
		fmt.Fprintf(w, "Profile %q removed.\n", name)
	}
}
