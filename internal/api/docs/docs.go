// Package docs serves the Redoc API reference generated from
// docs/openapi.yaml (repo root) via `make docs-build`. The generated page and
// the pinned redoc.standalone.js bundle are embedded so the docs work offline
// and from within the compiled binary.
package docs

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static
var assets embed.FS

const (
	contentTypeHTML = "text/html; charset=utf-8"
	contentTypeJS   = "application/javascript; charset=utf-8"
)

// Register mounts GET /docs (the rendered API reference) and its vendored
// Redoc bundle. Run `make docs-build` after editing docs/openapi.yaml to keep
// the embedded page in sync with the spec.
func Register(r *gin.Engine) {
	index := mustRead("static/index.html")
	script := mustRead("static/redoc.standalone.js")

	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, contentTypeHTML, index)
	})
	r.GET("/docs/redoc.standalone.js", func(c *gin.Context) {
		c.Data(http.StatusOK, contentTypeJS, script)
	})
}

func mustRead(path string) []byte {
	content, err := fs.ReadFile(assets, path)
	if err != nil {
		panic("docs: missing embedded asset " + path + ": " + err.Error())
	}
	return content
}
