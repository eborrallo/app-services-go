package view

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"

	"github.com/gin-gonic/gin"
)

// TemplRender implements the render.Render interface.
type TemplRender struct {
	Code int
	Data templ.Component
}

// Render implements the render.Render interface.
func (t TemplRender) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	w.WriteHeader(t.Code)
	if t.Data != nil {
		return t.Data.Render(context.Background(), w)
	}
	return nil
}

// WriteContentType implements the render.Render interface.
func (t TemplRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// Instance implements the render.Render interface.
func (t *TemplRender) Instance(name string, data interface{}) render.Render {
	if templData, ok := data.(templ.Component); ok {
		return &TemplRender{
			Code: http.StatusOK,
			Data: templData,
		}
	}
	return nil
}

// RegisterRoutes registers routes for the server.
func RegisterRoutes(s *gin.Engine) {

	// Define HTML renderer for template engine.
	s.HTMLRender = &TemplRender{}

	// Handle static files.
	s.Static("/static", "./static")

	// Handle index page view.
	s.GET("/", indexViewHandler)

	// Handle API endpoints.
	s.GET("/api/hello-world", showContentAPIHandler)
}

