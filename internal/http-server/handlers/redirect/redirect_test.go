package redirect_test

import (
	"cmd/service/internal/http-server/handlers/redirect"
	"cmd/service/internal/http-server/handlers/redirect/mocks"
	"cmd/service/internal/lib/api"
	"cmd/service/internal/lib/logger/handlers/slogdiscard"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewUrlGetter(t)
			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetUrl", tc.alias).Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))
			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToUrl, err := api.GetRedirect(ts.URL + "/" + tc.alias)

			require.NoError(t, err)

			assert.Equal(t, tc.url, redirectedToUrl)

		})
	}
}
