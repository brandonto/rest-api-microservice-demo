package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/brandonto/rest-api-microservice-demo/db"
	"github.com/brandonto/rest-api-microservice-demo/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Middleware to extract URL query params prior to processing
//
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		// Parses all relevant query params
		//
		detailed := ListMessagesDetailedQueryParamDefault
		if detailedQueryParam := r.URL.Query().Get("detailed"); detailedQueryParam != "" {
			// Converts "detailed" query param to a boolean value
			//
			detailed, err = stringToBool(detailedQueryParam)
			if err != nil {
				// Respond with status Bad Request - no response payload
				//
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		limit := ListMessagesLimitQueryParamDefault
		if limitQueryParam := r.URL.Query().Get("limit"); limitQueryParam != "" {
			// Converts "limit" query param to a uint64 value
			//
			limit, err = strconv.ParseUint(limitQueryParam, 10, 64)
			if err != nil {
				// Respond with status Bad Request - no response payload
				//
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Just check against upper limit... lower limit should be caught by
			// the ParseUint function above
			//
			if limit > ListMessagesLimitQueryParamMax {
				// Respond with status Bad Request - no response payload
				//
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// If limit is 0, we can just return an empty array now with status
			// OK
			//
			if limit == 0 {
				render.Status(r, http.StatusOK)
				render.JSON(w, r, []model.Message{})
				return
			}
		}

		afterId := ListMessagesAfterIdQueryParamDefault
		if afterIdQueryParam := r.URL.Query().Get("afterId"); afterIdQueryParam != "" {
			// Converts "afterId" query param to a uint64 value
			//
			afterId, err = strconv.ParseUint(afterIdQueryParam, 10, 64)
			if err != nil {
				// Respond with status Bad Request - no response payload
				//
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// Adds query params to the request context and forward to next
		// http.Handler
		//
		ctx := context.WithValue(r.Context(), "detailed", detailed)
		ctx = context.WithValue(ctx, "limit", limit)
		ctx = context.WithValue(ctx, "afterId", afterId)
		next.ServeHTTP(w, r.WithContext(ctx))
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
