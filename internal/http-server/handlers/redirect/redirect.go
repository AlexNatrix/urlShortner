package redirect

import (
	resp "cmd/service/internal/lib/api/responce"
	"cmd/service/internal/lib/logger/sl"
	"cmd/service/internal/storage"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlGetter
type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias empty")
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		resUrl, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got URL", slog.String("url", resUrl))

		http.Redirect(w, r, resUrl, http.StatusFound)

	}
}
