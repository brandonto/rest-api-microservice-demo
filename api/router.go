package api

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/brandonto/rest-api-microservice-demo/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewRouter(svcDb *db.Db, enableLogger bool) chi.Router {
	r := chi.NewRouter()

	// Use go-chi's built in Logger middleware to enable lightweight logging of
	// HTTP requests and responses
	//
	if enableLogger {
		r.Use(middleware.Logger)
	}

	// All HTTP responses in this API with a payload is JSON formatted. This is
	// still safe for empty HTTP responses because empty responses in this
	// application doesn't go through go-chi's render package.
	//
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Reads openapi.json
	//
	openApiFileBytes, err := ioutil.ReadFile("docs/openapi.json")
	if err != nil {
		log.Fatal(errors.New("Unable to open openapi.json file"))
	}

	// GET /openapi.json
	//
	r.Get("/openapi.json", func (w http.ResponseWriter, r *http.Request) {
		w.Write(openApiFileBytes)
	})

	// GET /swagger
	//
	workDir, _ := os.Getwd()
	htmlDir := http.Dir(filepath.Join(workDir, "html"))
	FileServer(r, "/swagger", htmlDir)

	// Configure API routes
	//
	r.Route("/messages", func(r chi.Router) {
		r.With(Paginate).Get("/", ListMessages(svcDb)) // GET /messages
		r.Post("/", CreateMessage(svcDb))              // POST /messages

		r.Route("/{messageId}", func(r chi.Router) {
			r.Use(GetMessageCtxFunc(svcDb))
			r.Get("/", GetMessage(svcDb))       // GET /messages/{messageId}
			r.Put("/", UpdateMessage(svcDb))    // PUT /messages/{messageId}
			r.Delete("/", DeleteMessage(svcDb)) // DELETE /messages/{messageId}
		})
	})

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
//
// This code snippet is taken from:
// https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
//
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
