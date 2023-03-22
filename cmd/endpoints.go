package main

import (
	_ "embed"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/richdawe/minimediaserver/services/catalog"
)

//go:embed root.tmpl.html
var rootTemplate string

//go:embed track.tmpl.html
var trackTemplate string

type TemplateRenderer struct {
	templates *template.Template
}

func (tr *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func getRoot(c echo.Context, catalogService *catalog.CatalogService) error {
	// TODO: how to catch errors in template rendering?
	return c.Render(http.StatusOK, "root", catalogService.GetTracks())
}

func getTrack(c echo.Context, catalogService *catalog.CatalogService) error {
	id := c.Param("id")
	track, err := catalogService.GetTrack(id)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that doesn't exist
		return err
	}
	// TODO: available data types => different query parameters in template
	return c.Render(http.StatusOK, "track", track)
}

func getTrackData(c echo.Context, catalogService *catalog.CatalogService) error {
	id := c.Param("id")
	track, err := catalogService.GetTrack(id)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that doesn't exist
		return err
	}
	r, err := catalogService.ReadTrack(track)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that can't be read
		return err
	}
	return c.Stream(http.StatusOK, "audio/flac", r) // TODO: get data type from track
}

func setupEndpoints(catalogService *catalog.CatalogService) (*echo.Echo, error) {
	t, err := template.New("root").Parse(rootTemplate)
	if err != nil {
		return nil, err
	}
	t.New("track").Parse(trackTemplate)
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
	e.GET("/tracks/:id", func(c echo.Context) error {
		return getTrack(c, catalogService)
	})
	e.GET("/tracks/:id/data", func(c echo.Context) error {
		return getTrackData(c, catalogService)
	})
	return e, nil
}
