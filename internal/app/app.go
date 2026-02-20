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
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/cqroot/tinyserver/internal/middleware"
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
	cssBytes, _ := cssFiles.ReadFile("style.css")
	cssContent = string(cssBytes)
}

func IsWildcardHosts(host string) bool {
	wildcardHosts := []string{"", "0.0.0.0", "::"}
	return slices.Contains(wildcardHosts, host)
}

func HandleFunc(c *gin.Context) {
	reqPath := c.Param("path")
	if reqPath == "" {
		reqPath = "/"
	}

	localPath := filepath.Join(".", reqPath)

	info, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		c.String(http.StatusNotFound, "404 not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if !info.IsDir() {
		c.File(localPath)
		return
	}

	files, err := os.ReadDir(localPath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

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

func Run(bindIp string, bindPort int, whitelist []string) error {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(middleware.WhitelistMiddleware(whitelist))

	tmpl := template.Must(template.New("dirlist").Parse(dirListTemplate))
	app.SetHTMLTemplate(tmpl)
	app.GET("/*path", HandleFunc)

	PrintAppInfo(bindIp, bindPort, whitelist)
	addr := net.JoinHostPort(bindIp, strconv.Itoa(bindPort))
	return app.Run(addr)
}
