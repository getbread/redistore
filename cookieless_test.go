package redistore

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookielessSessionIDStore(t *testing.T) {
	store := NewCookielessSessionIDStore("TEST", "ABCD")

	t.Run("load empty", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		sessionID, err := store.Load(req)
		assert.Empty(t, sessionID)
		assert.Equal(t, ErrHeaderNotFound, err)
	})
	t.Run("load invalid", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-SESSION-ID-TEST", "XYZ:value")
		sessionID, err := store.Load(req)
		assert.Empty(t, sessionID)
		assert.Equal(t, ErrInvalidSessionID, err)
	})
	t.Run("load valid", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-SESSION-ID-TEST", "94d5574a0ef464c629296fc9d263517944b94d1df9f3472fb7fb2d90af42ca36:value")
		sessionID, err := store.Load(req)
		assert.NotEmpty(t, sessionID)
		assert.NoError(t, err)
	})

	t.Run("save", func(t *testing.T) {
		rec := httptest.NewRecorder()
		store.Save("value", rec)
		assert.Equal(t, "94d5574a0ef464c629296fc9d263517944b94d1df9f3472fb7fb2d90af42ca36:value", rec.Header().Get("X-SESSION-ID-TEST"))
	})
}
