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

package app

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cqroot/tinyserver/internal/middleware"
	"github.com/cqroot/tinyserver/pkg/netutil"
	"github.com/gin-gonic/gin"
)

//go:embed style.css
var cssFiles embed.FS
var cssContent string

const dirListTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.5, user-scalable=yes">
  <title>{{ .Path }}</title>
  <style>{{ .Css }}</style>
</head>
<body>
  <h1>Index of {{ .Path }}</h1>
  <ul>
    {{ range .Items }}
      <li>
        <span class="icon">{{ if .IsDir }}📁 {{ else }}📄 {{ end }}</span>
        <a href="{{ .Name }}" class="{{ if .IsDir }}dir{{ else }}file{{ end }}">{{ .Name }}</a>
      </li>
    {{ end }}
  </ul>
</body>
</html>
`

func init() {
	cssBytes, err := cssFiles.ReadFile("style.css")
	if err != nil {
		panic(fmt.Sprintf("failed to read embedded style.css: %v", err))
	}
	cssContent = string(cssBytes)
}

type App struct {
	workDir string
}

func New(workDir string) (*App, error) {
	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	a := App{
		workDir: absWorkDir,
	}
	return &a, nil
}

func IsWildcardHosts(host string) bool {
	wildcardHosts := []string{"", "0.0.0.0", "::"}
	return slices.Contains(wildcardHosts, host)
}

func (a App) IsAvailablePath(path string) bool {
	return !strings.HasPrefix(filepath.Base(path), ".")
}

func (a App) HandleDir(c *gin.Context, reqPath string, localPath string) {
	if !strings.HasSuffix(reqPath, "/") {
		c.Redirect(http.StatusMovedPermanently, reqPath+"/")
		return
	}

	files, err := os.ReadDir(localPath)
	if err != nil {
		slog.Error("Internal server error.", slog.String("err", err.Error()))
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	slices.SortFunc(files, func(a, b os.DirEntry) int {
		if a.IsDir() && !b.IsDir() {
			return -1
		}
		if !a.IsDir() && b.IsDir() {
			return 1
		}
		return strings.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
	})

	type item struct {
		Name  string
		IsDir bool
	}
	items := make([]item, 0, len(files))
	if reqPath != "/" {
		items = append(items, item{
			Name:  "..",
			IsDir: true,
		})
	}
	for _, f := range files {
		name := f.Name()
		if !a.IsAvailablePath(name) {
			continue
		}
		if f.IsDir() {
			name += "/"
		}
		items = append(items, item{
			Name:  name,
			IsDir: f.IsDir(),
		})
	}

	c.HTML(http.StatusOK, "dirlist", gin.H{
		"Path":  reqPath,
		"Items": items,
		"Css":   template.CSS(cssContent),
	})
}

func (a App) HandleFunc(c *gin.Context) {
	reqPath := c.Param("path")
	if reqPath == "" {
		reqPath = "/"
	}

	localPath, err := filepath.Abs(filepath.Join(a.workDir, reqPath))
	if err != nil {
		slog.Error("Internal server error.", slog.String("err", err.Error()))
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	if !strings.HasPrefix(localPath, a.workDir) {
		slog.Error("Access denied.", slog.String("path", localPath))
		c.String(http.StatusForbidden, "access denied")
		return
	}

	if !a.IsAvailablePath(localPath) {
		c.String(http.StatusNotFound, "404 not found")
		return
	}

	info, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		c.String(http.StatusNotFound, "404 not found")
		return
	}
	if err != nil {
		slog.Error("Internal server error.", slog.String("err", err.Error()))
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	if !info.IsDir() {
		c.File(localPath)
		return
	}

	a.HandleDir(c, reqPath, localPath)
}

func LogAppInfo(bindIp string, bindPort int, whitelist []string) {
	slog.Info("Starting TinyServer.",
		slog.String("bind_addr", net.JoinHostPort(bindIp, strconv.Itoa(bindPort))),
		slog.String("whitelist", strings.Join(whitelist, ", ")),
	)

	if IsWildcardHosts(bindIp) {
		ips, err := netutil.GetLocalIPs()
		if err != nil {
			return
		}

		slog.Info("Available bind addresses:")
		for _, ip := range ips {
			slog.Info(fmt.Sprintf("  %s http://%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render("•"),
				net.JoinHostPort(ip, strconv.Itoa(bindPort)),
			))
		}
	}
}

func (a App) Run(bindIp string, bindPort int, whitelist []string) error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.WhitelistMiddleware(whitelist))
	engine.Use(middleware.LoggerMiddleware())

	tmpl := template.Must(template.New("dirlist").Parse(dirListTemplate))
	engine.SetHTMLTemplate(tmpl)
	engine.GET("/*path", a.HandleFunc)

	LogAppInfo(bindIp, bindPort, whitelist)
	addr := net.JoinHostPort(bindIp, strconv.Itoa(bindPort))
	return engine.Run(addr)
}
