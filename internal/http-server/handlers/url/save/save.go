package save

import (
	resp "cmd/service/internal/lib/api/responce"
	"cmd/service/internal/lib/logger/sl"
	"cmd/service/internal/lib/random"
	"cmd/service/internal/storage"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

// TODO: md config?
const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlSaver
type UrlSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Responce struct {
	resp.Responce
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.save.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationErrors(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomStr(aliasLength)
		}
		id, err := urlSaver.SaveUrl(req.URL, alias)
		if errors.Is(err, storage.ErrUrlExist) {
			log.Info("url already exist", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}
		log.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, Responce{
			Responce: resp.OK(),
			Alias:    alias,
		})

	}
}
