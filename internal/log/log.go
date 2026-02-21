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

package log

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

func InitLogger() {
	handler := log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.InfoLevel,
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05",
	})
	logger := slog.New(handler)

	slog.SetDefault(logger)
}
