package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/richdawe/minimediaserver/internal/httprange"
	"github.com/richdawe/minimediaserver/internal/offsetlimitreader"
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

func getRoot(c echo.Context, catalogService catalog.CatalogService) error {
	// TODO: how to catch errors in template rendering?
	return c.Render(http.StatusOK, "root.tmpl.html", make(map[string]string, 0))
}

func getTracks(c echo.Context, catalogService catalog.CatalogService) error {
	// TODO: how to catch errors in template rendering?
	tracks, _ := catalogService.GetTracks()
	return c.Render(http.StatusOK, "tracks.tmpl.html", tracks)
}

func getTracksByID(c echo.Context, catalogService catalog.CatalogService) error {
	id := c.Param("id")
	track, err := catalogService.GetTrack(id)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that doesn't exist
		return err
	}
	// TODO: available data types => different query parameters in template
	return c.Render(http.StatusOK, "tracksbyid.tmpl.html", track)
}

func getTracksByIDData(c echo.Context, catalogService catalog.CatalogService, cacheMaxAge int) error {
	id := c.Param("id")
	track, err := catalogService.GetTrack(id)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that doesn't exist
		return err
	}

	// Parse any requested byte ranges.
	// This article was really helpful in adding this functionality;
	// <https://www.zeng.dev/post/2023-http-range-and-play-mp4-in-browser/>
	var httpRanges []httprange.HttpRange

	rangeVal := c.Request().Header.Get("Range")
	if rangeVal != "" {
		httpRanges, err = httprange.ParseRange(rangeVal, track.DataLen)
		if err != nil {
			// TODO: return appropriate error
			return err
		}
	}

	r, err := catalogService.ReadTrack(track)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that can't be read
		return err
	}

	// Ensure that the first requested range is returned.
	responseCode := http.StatusOK
	if len(httpRanges) > 0 {
		start := httpRanges[0].Start
		length := httpRanges[0].Length
		maxChunkSize := int64(1024 * 1024)
		if length > maxChunkSize {
			length = maxChunkSize
		}

		r = offsetlimitreader.New(r, start, length)

		// TODO: test coverage for range header; should use httprange code from go?
		// TODO: include range in HTTP logs
		rangeResponse := fmt.Sprintf("bytes %d-%d/%d", start, start+length-1, track.DataLen)
		c.Response().Header().Add("Content-Range", rangeResponse)
		c.Response().Header().Add("Content-Length", strconv.FormatInt(length, 10))
		responseCode = http.StatusPartialContent
	}

	// Allow ranges to be requested.
	c.Response().Header().Add("Accept-Ranges", "bytes")
	// Allow the track data to be cached by the client.
	c.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", cacheMaxAge))
	return c.Stream(responseCode, track.MIMEType, r)
}

func getPlaylists(c echo.Context, catalogService catalog.CatalogService) error {
	// TODO: how to catch errors in template rendering?
	_, playlists := catalogService.GetTracks()
	return c.Render(http.StatusOK, "playlists.tmpl.html", playlists)
}

func getPlaylistsByID(c echo.Context, catalogService catalog.CatalogService) error {
	id := c.Param("id")
	playlist, err := catalogService.GetPlaylist(id)
	if err != nil {
		// TODO: return appropriate error for e.g.: track that doesn't exist
		return err
	}
	// TODO: available data types => different query parameters in template
	return c.Render(http.StatusOK, "playlistsbyid.tmpl.html", playlist)
}

func templateAddInt(a, b int) int {
	return a + b
}

func setupEndpoints(config Config, catalogService catalog.CatalogService) (*echo.Echo, error) {
	t := template.New("endpoints").Funcs(template.FuncMap{
		"addInt": templateAddInt,
	})
	t, err := t.ParseFS(templatesContent, "templates/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	tr := &TemplateRenderer{
		templates: t,
	}

	e := echo.New()
	e.Renderer = tr

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())

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
		return getTracksByIDData(c, catalogService, config.CacheMaxAge)
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
