package api

import (
    "context"
    "net/http"
    "strconv"

    "github.com/brandonto/rest-api-microservice-demo/db"
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
func GetMessageCtxFunc(svcDb *db.Db) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            var detailedMessage *model.DetailedMessage

            if messageIdParam := chi.URLParam(r, "messageId"); messageIdParam != "" {

                // Converts messageId URL param to a uint64 id
                //
                messageId, err := strconv.ParseUint(messageIdParam, 10, 64)
                if err != nil {
                    // Bails out early with a 400 if messageId in the URL params
                    // is not an integer - no response payload.
                    //
                    w.WriteHeader(http.StatusBadRequest)
                    return
                }

                // Retrieves message from the database
                //
                detailedMessage, err = svcDb.GetMessage(messageId)
                if err != nil {
                    // Bails out early with a 404 if message does not exist in
                    // the database - no response payload.
                    //
                    w.WriteHeader(http.StatusNotFound)
                    return
                }
            } else {
                // Bails out early with a 404 if messageId does not exist in the
                // URL params - no response payload.
                //
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            // Adds message to the request context and forward to next
            // http.Handler
            //
            ctx := context.WithValue(r.Context(), "detailedMessage", detailedMessage)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
