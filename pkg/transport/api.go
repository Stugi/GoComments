package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"stugi/go-comment/pkg/model"
	"stugi/go-comment/pkg/service"
	"stugi/go-comment/pkg/transport/middleware"

	"github.com/gorilla/mux"
)

type API struct {
	service *service.Service
	r       *mux.Router
}

func New(service *service.Service) *API {
	api := API{
		service: service,
		r:       mux.NewRouter(),
	}
	api.endpoints()

	// middleware
	api.r.Use(middleware.RequestIDMiddleware)
	api.r.Use(middleware.LoggingMiddleware)

	return &api
}

func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// Принять ID новости и текст комментария или ID родительского комментария (на которые кто-то отвечает), и сохранить это в БД.
	api.r.HandleFunc("/comments", api.CreateComment).Methods(http.MethodPost, http.MethodOptions)
	//Метод для выгрузки всех комментарием по ID новости.
	//Пока представим, что комментариев не много и нам не нужно делать пагинацию.
	api.r.HandleFunc("/comments/news/{id}", api.GetComments).Methods(http.MethodGet, http.MethodOptions)
}

func (api *API) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := r.Body
	defer body.Close()
	// empty body
	if body == http.NoBody {
		sendError(ctx, w, fmt.Errorf("empty body"), http.StatusBadRequest)
		return
	}

	var comment model.Comment
	// body -> comment
	err := json.NewDecoder(body).Decode(&comment)
	if err != nil {
		sendError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	idComment, err := api.service.AddComment(ctx, comment)
	if err != nil {
		sendError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"id": idComment})
}

func (api *API) GetComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	rawID := vars["id"]
	id, err := strconv.Atoi(rawID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := api.service.GetCommentsByNewsID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func sendError(ctx context.Context, w http.ResponseWriter, err error, status int) {
	var response model.ErrorResponse
	response.Code = status
	response.Message = err.Error()
	// Попробуем извлечь RequestID из контекста
	if requestID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		response.RequestID = requestID
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}
