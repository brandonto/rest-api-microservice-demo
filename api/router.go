package api

import (
    "github.com/brandonto/rest-api-microservice-demo/db"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/render"
)

func CreateRouter(svcDb *db.Db) chi.Router {
    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(render.SetContentType(render.ContentTypeJSON))

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
