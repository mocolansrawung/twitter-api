package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/evermos/boilerplate-go/internal/domain/tweet"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type TweetHandler struct {
	TweetService tweet.TweetService
	// AuthMiddleware *middleware.Authentication
}

func ProvideTweetHandler(tweetService tweet.TweetService) TweetHandler {
	return TweetHandler{
		TweetService: tweetService,
	}
}

// Router for Tweet domain
func (h *TweetHandler) Router(r chi.Router) {
	r.Route("/tweets", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/", h.CreateTweet)
			r.Get("/", h.ResolveAllTweets)
			r.Get("/{id}", h.ResolveTweetByID)
		})
	})
}

func (h *TweetHandler) CreateTweet(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat tweet.TweetRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, _ := uuid.NewV4()

	tweet, err := h.TweetService.Create(requestFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, tweet)
}

func (h *TweetHandler) ResolveAllTweets(w http.ResponseWriter, r *http.Request) {
	pageString := r.URL.Query().Get("page")
	page, err := convertIdParamsToInt(pageString)
	if err != nil || page < 0 {
		page = 0
	}

	limitString := r.URL.Query().Get("limit")
	limit, err := convertIdParamsToInt(limitString)
	if err != nil || limit <= 0 {
		limit = 10
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	tweets, err := h.TweetService.ResolveTweets(page, limit, sort, order)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, tweets)
}

func (h *TweetHandler) ResolveTweetByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	tweet, err := h.TweetService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, tweet)
}

// Internal function
func convertIdParamsToInt(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error during conversion")
		return 0, err
	}

	return id, nil
}
