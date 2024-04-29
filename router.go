package main

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
    r := chi.NewRouter()

    r.Use(middleware.Logger)

    r.Route("/messages", func(r chi.Router) {
        r.With(Paginate).Get("/", ListMessages) // GET /messages
        r.Post("/", CreateMessage)              // POST /messages

        r.Route("/{messageId}", func(r chi.Router) {
            r.Use(MessageCtx)
            r.Get("/", GetMessage)              // GET /messages/{messageId}
            r.Put("/", UpdateMessage)           // PUT /messages/{messageId}
            r.Delete("/", DeleteMessage)        // DELETE /messages/{messageId}
        })
    })

    return r
}
