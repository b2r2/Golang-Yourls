package yourls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	baseURL = "" //Your baseurl "http://your-own-domain-here.com/yourls-api.php"
	token   = "" // Your secret signature token from yourls
	url     = "https://google.com"
)

func NewMock(a, u string) *UserData {
	ud := New(token, baseURL)
	ud.Action = a
	ud.URL = u
	return ud
}

func TestDo(t *testing.T) {
	actions := []string{"shorturl", "expand", "url-stats", "stats", "db-stats"}
	for _, a := range actions {
		testCase := struct {
			name string
			in   *UserData
			out  int
		}{
			name: a,
			in:   NewMock(a, url),
			out:  http.StatusOK,
		}
		c, err := testCase.in.Get()
		assert.NoError(t, err)
		assert.Equal(t, c.StatusCode, testCase.out)
		url = c.Shorturl
	}
}
