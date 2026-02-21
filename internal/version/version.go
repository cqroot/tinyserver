/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package version

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtWith = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

type Info struct {
	Version   string
	Commit    string
	Date      string
	BuiltWith string
}

func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		BuiltWith: builtWith,
	}
}

func (i Info) String() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	sb := strings.Builder{}
	sb.WriteString("\n  ")
	sb.WriteString(style.Render("• Version:      "))
	sb.WriteString(i.Version)

	sb.WriteString("\n  ")
	sb.WriteString(style.Render("• Commit:       "))
	sb.WriteString(i.Commit)

	sb.WriteString("\n  ")
	sb.WriteString(style.Render("• Built at:     "))
	sb.WriteString(i.Date)

	sb.WriteString("\n  ")
	sb.WriteString(style.Render("• Built with:   "))
	sb.WriteString(i.BuiltWith)
	return sb.String()
}
