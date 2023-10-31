package controller

import (
	"embed"
	"my-project-name/app/shared/archetype/container"
	einar "my-project-name/app/shared/archetype/echo_server"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

//go:embed *.html
var html embed.FS

//go:embed *.css
var css embed.FS

func init() {
	einar.EmbeddedPatterns = append(einar.EmbeddedPatterns, einar.EmbeddedPattern{
		Content: html,
		Pattern: "content.html",
	})
	einar.EmbeddedPatterns = append(einar.EmbeddedPatterns, einar.EmbeddedPattern{
		Content: css,
		Pattern: "content.css",
	})
	container.InjectInboundAdapter(func() error {
		einar.Echo.GET("/content", render)
		einar.Echo.GET("/content.css", echo.WrapHandler(http.FileServer(http.FS(css))))
		return nil
	}, container.InjectionProps{
		DependencyID: uuid.NewString(),
	})
}

func render(c echo.Context) error {
	data := map[string]interface{}{
		"layoutComponentDefault":  "content",
		"contentComponentDefault": "content/content-children-b", //Investigar la manera de pasar el contentComponentDefault por contexto
	}
	standalone := c.Request().Header.Get("standalone")
	if standalone == "true" {
		return c.Render(http.StatusOK, "content.html", data)
	}
	return c.Render(http.StatusOK, "layout.html", data)
}
