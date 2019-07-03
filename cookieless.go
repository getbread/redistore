package redistore

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// errors
var (
	ErrHeaderNotFound   = errors.New("header not found")
	ErrInvalidSessionID = errors.New("invalid session id")
)

// A CookielessSessionIDStore stores session ids in http headers
type CookielessSessionIDStore struct {
	name     string
	password string
}

// NewCookielessSessionIDStore creates a new CookielessSessionIDStore
func NewCookielessSessionIDStore(name, password string) *CookielessSessionIDStore {
	return &CookielessSessionIDStore{
		name:     name,
		password: password,
	}
}

// Load attempts to load a session id from an http request
func (store CookielessSessionIDStore) Load(r *http.Request) (string, error) {
	value := r.Header.Get(store.headerName())
	if value == "" {
		return "", ErrHeaderNotFound
	}

	idx := strings.IndexByte(value, ':')
	if idx < 0 {
		return "", ErrInvalidSessionID
	}

	mac, sessionID := value[:idx], value[idx+1:]
	if !store.verify(mac, sessionID) {
		return "", ErrInvalidSessionID
	}

	return sessionID, nil
}

// Save attemps to save a session id to an http response writer
func (store CookielessSessionIDStore) Save(sessionID string, w http.ResponseWriter) {
	mac := store.sign(sessionID)
	w.Header().Set(store.headerName(), fmt.Sprintf("%s:%s", mac, sessionID))
}

func (store CookielessSessionIDStore) sign(message string) (mac string) {
	h := hmac.New(sha256.New, []byte(store.password))
	io.WriteString(h, message)
	signature := h.Sum(nil)
	return hex.EncodeToString(signature)
}

func (store CookielessSessionIDStore) verify(mac, message string) bool {
	mac1, err := hex.DecodeString(mac)
	if err != nil {
		return false
	}

	mac2, err := hex.DecodeString(store.sign(message))
	if err != nil {
		return false
	}

	return hmac.Equal(mac1, mac2)
}

func (store CookielessSessionIDStore) headerName() string {
	return fmt.Sprintf("X-SESSION-ID-%s", store.name)
}
