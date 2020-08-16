package pagecomp

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPage(t *testing.T) {
	t.Run("Query Parameter", func(t *testing.T) {
		r, err := http.NewRequest("GET", "https://www.example.com/a/b/c/d?q1=q1value&q2=q2value", nil)
		assert.NoError(t, err)
		r.ParseForm()

		page := NewPage(r)
		assert.Equal(t, "/a/b/c/d", page.path)
		assert.Equal(t, "q1value", page.params["q1"])
		assert.Equal(t, "q2value", page.params["q2"])
	})

	t.Run("Post Parameter", func(t *testing.T) {
		postValues := url.Values{}
		postValues.Add("p1", "p1value")
		postValues.Add("p2", "p2value")
		r, err := http.NewRequest("POST", "https://www.example.com/a/b/c/d", bytes.NewBufferString(postValues.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		assert.NoError(t, err)

		page := NewPage(r)
		assert.Equal(t, "/a/b/c/d", page.path)
		assert.Equal(t, "p1value", page.params["p1"])
		assert.Equal(t, "p2value", page.params["p2"])
	})
}
