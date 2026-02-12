package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vlah-software-house/paddle-api-mock/internal/middleware"
	"github.com/vlah-software-house/paddle-api-mock/internal/models"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	writeJSON(w, status, models.PaddleResponse{
		Data: data,
		Meta: models.Meta{
			RequestID: middleware.GetRequestID(r.Context()),
		},
	})
}

func respondList(w http.ResponseWriter, r *http.Request, data interface{}, total int) {
	writeJSON(w, http.StatusOK, models.PaddleListResponse{
		Data: data,
		Meta: models.Meta{
			RequestID: middleware.GetRequestID(r.Context()),
		},
		Pagination: &models.PaginationInfo{
			PerPage:        50,
			HasMore:        false,
			EstimatedTotal: total,
		},
	})
}

func respondError(w http.ResponseWriter, r *http.Request, status int, errType, code, detail string) {
	writeJSON(w, status, models.ErrorResponse{
		Error: models.ErrorDetail{
			Type:   errType,
			Code:   code,
			Detail: detail,
		},
		Meta: models.Meta{
			RequestID: middleware.GetRequestID(r.Context()),
		},
	})
}

func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
