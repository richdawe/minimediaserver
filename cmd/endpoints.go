package main

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/richdawe/minimediaserver/services/catalog"
)

//go:embed templates/*
var templatesContent embed.FS

//go:embed static/*
var staticContent embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

func (tr *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func getRoot(c echo.Context, catalogService *catalog.CatalogService) error {
	// TODO: how to catch errors in template rendering?
	return c.Render(http.StatusOK, "root.tmpl.html", make(map[string]string, 0))
}

func getTracks(c echo.Context, catalogService *catalog.CatalogService) error {
	// TODO: how to catch errors in template rendering?
	tracks, _ := catalogService.GetTracks()
	return c.Render(http.StatusOK, "tracks.tmpl.html", tracks)
}

func getTracksByID(c echo.Context, catalogService *catalog.CatalogService) error {
	id := c.Param("id")
	track, err := catalogService.GetTrack(id)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that doesn't exist
		return err
	}
	// TODO: available data types => different query parameters in template
	return c.Render(http.StatusOK, "tracksbyid.tmpl.html", track)
}

func getTracksByIDData(c echo.Context, catalogService *catalog.CatalogService) error {
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
	return c.Stream(http.StatusOK, track.MIMEType, r)
}

func getPlaylists(c echo.Context, catalogService *catalog.CatalogService) error {
	// TODO: how to catch errors in template rendering?
	_, playlists := catalogService.GetTracks()
	return c.Render(http.StatusOK, "playlists.tmpl.html", playlists)
}

func getPlaylistsByID(c echo.Context, catalogService *catalog.CatalogService) error {
	id := c.Param("id")
	playlist, err := catalogService.GetPlaylist(id)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that doesn't exist
		return err
	}
	// TODO: available data types => different query parameters in template
	return c.Render(http.StatusOK, "playlistsbyid.tmpl.html", playlist)
}

func setupEndpoints(catalogService *catalog.CatalogService) (*echo.Echo, error) {
	t, err := template.ParseFS(templatesContent, "templates/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	tr := &TemplateRenderer{
		templates: t,
	}

	e := echo.New()
	e.Renderer = tr

	e.Pre(middleware.RemoveTrailingSlash())

	// Don't wait process requests indefinitely.
	e.Server.ReadTimeout = time.Duration(60 * time.Second)
	e.Server.WriteTimeout = time.Duration(60 * time.Second)

	e.GET("/", func(c echo.Context) error {
		return getRoot(c, catalogService)
	})
	e.GET("/tracks", func(c echo.Context) error {
		return getTracks(c, catalogService)
	})
	e.GET("/tracks/:id", func(c echo.Context) error {
		return getTracksByID(c, catalogService)
	})
	e.GET("/tracks/:id/data", func(c echo.Context) error {
		return getTracksByIDData(c, catalogService)
	})
	e.GET("/playlists", func(c echo.Context) error {
		return getPlaylists(c, catalogService)
	})
	e.GET("/playlists/:id", func(c echo.Context) error {
		return getPlaylistsByID(c, catalogService)
	})
	e.GET("/static/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		path := filepath.Join("static", filename)
		data, err := staticContent.ReadFile(path)
		if err != nil {
			// TODO: return appropriate error for e.g.: file that can't be found
			return err
		}

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
		mimeType := "application/binary"
		switch {
		case strings.HasSuffix(filename, ".css"):
			mimeType = "text/css"
		case strings.HasSuffix(filename, ".html"):
			mimeType = "text/html"
		case strings.HasSuffix(filename, ".js"):
			mimeType = "text/javascript"
		}

		return c.Blob(http.StatusOK, mimeType, data)
	})
	return e, nil
}
