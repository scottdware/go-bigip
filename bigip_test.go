package bigip

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenSession(t *testing.T) {
	now := time.Now()
	timeout := 1200
	startTime, err := now.MarshalJSON()
	require.NoError(t, err)
	h := handler{
		StartTime:  string(startTime),
		Timeout:    timeout,
		LastUpdate: now.Unix() * microToSeconds,
		Expiry:     now.Add(time.Duration(timeout)*time.Second).Unix() * microToSeconds,
	}
	user := "user"
	password := "password"
	testServer := httptest.NewServer(&h)
	var b *BigIP
	t.Run("should login", func(t *testing.T) {
		b, err = NewTokenSession(testServer.URL, user, password, "tmos", nil)
		require.NoError(t, err)
		assert.Equal(t, 1, h.loginCounter)
		assert.Equal(t, 0, h.refreshCounter)
		assert.Equal(t, user, b.User)
		assert.Equal(t, password, b.Password)
		assert.True(t, b.TokenExpiry.Sub(now) > 1000*time.Second)
	})

	t.Run("should refresh with unexpired token", func(t *testing.T) {
		now := time.Now()
		startTime, err := now.MarshalJSON()
		require.NoError(t, err)
		h.StartTime = string(startTime)
		h.Timeout = 3600
		h.LastUpdate = now.Unix() * microToSeconds
		h.Expiry = now.Add(time.Duration(h.Timeout)*time.Second).Unix() * microToSeconds
		expiryBefore := b.TokenExpiry
		err = b.RefreshTokenSession(time.Duration(h.Timeout) * time.Second)
		require.NoError(t, err)
		assert.Equal(t, 1, h.loginCounter)
		assert.Equal(t, 1, h.refreshCounter)
		assert.True(t, b.TokenExpiry.Sub(expiryBefore) > 0)
	})

	t.Run("refresh should login with expired token", func(t *testing.T) {
		b.TokenExpiry = now.Add(-10000 * time.Second) // token is expred
		err = b.RefreshTokenSession(3600 * time.Second)
		require.NoError(t, err)
		assert.Equal(t, 2, h.loginCounter)
		assert.Equal(t, 1, h.refreshCounter)
	})
	t.Run("should login if refresh failes", func(t *testing.T) {
		h.Timeout = 0 // force refresh error (see handler below)
		err = b.RefreshTokenSession(3600 * time.Second)
		require.NoError(t, err)
		assert.Equal(t, 3, h.loginCounter)
		assert.Equal(t, 2, h.refreshCounter)
	})
}

type handler struct {
	StartTime      string
	Timeout        int
	LastUpdate     int64
	Expiry         int64
	loginCounter   int
	refreshCounter int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.URL.Path == "/mgmt/shared/authn/login" {
		h.loginCounter++
		tmpl, err := template.New("response").Parse(loginResponseTmpl)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, h)
		return
	}
	if r.Method == "PATCH" && strings.HasPrefix(r.URL.Path, "/mgmt/shared/authz/tokens/") {
		h.refreshCounter++
		if h.Timeout <= 0 {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		tmpl, err := template.New("response").Parse(refreshResponse)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, h)
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
}

const loginResponseTmpl = `{
  "username": "a-user",
  "loginReference": {
    "link": "https://localhost/mgmt/cm/system/authn/providers/tmos/1f44a60e-11a7-3c51-a49f-82983026b41b/login"
  },
  "loginProviderName": "tmos",
  "token": {
    "token": "KZ44TOKEN7SNOTZNR7D7UP24SC",
    "name": "KZ44TOKEN7SNOTZNR7D7UP24SC",
    "userName": "a-user",
    "authProviderName": "tmos",
    "user": {
      "link": "https://localhost/mgmt/cm/system/authn/providers/tmos/1f44a60e-11a7-3c51-a49f-82983026b41b/users/40994460-44e3-3a09-8837-2189119665e5"
    },
    "timeout": {{ .Timeout }},
    "startTime": {{ .StartTime }},
    "address": "192.168.1.2",
    "partition": "[All]",
    "generation": 1,
    "lastUpdateMicros": {{ .LastUpdate }},
    "expirationMicros": {{ .Expiry }},
    "kind": "shared:authz:tokens:authtokenitemstate",
    "selfLink": "https://localhost/mgmt/shared/authz/tokens/KZ44TOKEN7SNOTZNR7D7UP24SC"
  },
  "generation": 0,
  "lastUpdateMicros": 0
}
`

const refreshResponse = `{
  "token": "KZ44TOKEN7SNOTZNR7D7UP24SC",
  "name": "KZ44TOKEN7SNOTZNR7D7UP24SC",
  "userName": "a-user",
  "authProviderName": "tmos",
  "user": {
    "link": "https://localhost/mgmt/cm/system/authn/providers/tmos/1f44a60e-11a7-3c51-a49f-82983026b41b/users/40994460-44e3-3a09-8837-2189119665e5"
  },
  "timeout": {{ .Timeout }},
  "startTime": {{ .StartTime }},
  "address": "192.168.1.2",
  "partition": "[All]",
  "generation": 2,
  "lastUpdateMicros": {{ .LastUpdate }},
  "expirationMicros": {{ .Expiry }},
  "kind": "shared:authz:tokens:authtokenitemstate",
  "selfLink": "https://localhost/mgmt/shared/authz/tokens/KZ44TOKEN7SNOTZNR7D7UP24SC"
}
`
