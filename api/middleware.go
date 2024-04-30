package api

import (
    "context"
    "net/http"

    "github.com/brandonto/rest-api-microservice-demo/model"

    "github.com/go-chi/chi/v5"
    //"github.com/go-chi/render"
)

func Paginate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    })
}

// Middleware to load message specified by ID {messageId} prior to processing.
// Also bails out early with 404 Not Found if message does not exist.
// 
func MessageCtx(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var message *model.Message
        var err error

        if messageId := chi.URLParam(r, "messageId"); messageId != "" {
            message = &model.Message{Id: 1, Payload: "test"}
            err = nil
        } else {
            // Bails out early with a 404 if messageId does not exist in the URL
            // params.
            //
		    w.WriteHeader(http.StatusNotFound)
            return
        }

        if err != nil {
            // Bails out early with a 404 if message does not exist.
            //
		    w.WriteHeader(http.StatusNotFound)
            return
        }

        // Adds message to the request context and forward to next http.Handler
        //
        ctx := context.WithValue(r.Context(), "message", message)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
