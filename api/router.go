package api

import (
    "github.com/brandonto/rest-api-microservice-demo/db"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/render"
)

func NewRouter(svcDb *db.Db) chi.Router {
    r := chi.NewRouter()

    // Use go-chi's built in Logger middleware to enable lightweight logging of
    // HTTP requests and responses
    //
    r.Use(middleware.Logger)

    // All HTTP responses in this API with a payload is JSON formatted. This is
    // still safe for empty HTTP responses because empty responses in this
    // application doesn't go through go-chi's render package.
    //
    r.Use(render.SetContentType(render.ContentTypeJSON))

    // Configure API routes
    //
    r.Route("/messages", func(r chi.Router) {
        r.With(Paginate).Get("/", ListMessages(svcDb)) // GET /messages
        r.Post("/", CreateMessage(svcDb))              // POST /messages

        r.Route("/{messageId}", func(r chi.Router) {
            r.Use(GetMessageCtxFunc(svcDb))
            r.Get("/", GetMessage(svcDb))              // GET /messages/{messageId}
            r.Put("/", UpdateMessage(svcDb))           // PUT /messages/{messageId}
            r.Delete("/", DeleteMessage(svcDb))        // DELETE /messages/{messageId}
        })
    })

    return r
}
