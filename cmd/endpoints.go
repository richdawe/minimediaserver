package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/richdawe/minimediaserver/services/catalog"
)

//go:embed root.tmpl.html
var rootTemplate string

type TemplateRenderer struct {
	templates *template.Template
}

func (tr *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func getRoot(c echo.Context, catalogService *catalog.CatalogService) error {
	/*
		htmlDoc := "<html><head><title>Title</title></head><body><em>Oh hai</em></body></html>"
		return c.HTML(http.StatusOK, htmlDoc)
	*/
	//fmt.Printf("tracks = %+v\n", catalogService.GetTracks())
	// TODO: how to catch errors in template rendering?
	return c.Render(http.StatusOK, "root", catalogService.GetTracks())
}

func getTrack(c echo.Context) error {
	id := c.Param("id")
	// TODO: template music player for a track
	return c.String(http.StatusOK, fmt.Sprintf("dat track %s is here: %s", id, id))
}

func setupEndpoints(catalogService *catalog.CatalogService) (*echo.Echo, error) {
	t, err := template.New("root").Parse(rootTemplate)
	if err != nil {
		return nil, err
	}
	tr := &TemplateRenderer{
		templates: t,
	}

	e := echo.New()
	e.Renderer = tr
	e.GET("/", func(c echo.Context) error {
		return getRoot(c, catalogService)
	})
	e.GET("/tracks/:id", getTrack)
	return e, nil
}
