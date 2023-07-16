package save_test

import (
	"bytes"
	"cmd/service/internal/http-server/handlers/url/save"
	"cmd/service/internal/http-server/handlers/url/save/mocks"
	"cmd/service/internal/lib/logger/handlers/slogdiscard"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
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
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty Url",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is required",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlSaverMock := mocks.NewUrlSaver(t)
			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On(
					"SaveUrl", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()

			}
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)
			input := fmt.Sprintf(`{"url":"%s","alias":"%s"}`, tc.url, tc.alias)
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Responce

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}