package pagecomp

import (
	"net/http"
	"reflect"
)

// Page represent a page to compare
type Page struct {
	path   string
	params map[string]string
}

// Equal return page compare result
func (p Page) Equal(another Page) bool {
	if p.path != another.path {
		return false
	}
	if !reflect.DeepEqual(p.params, another.params) {
		return false
	}
	return true
}

// NewPage create page from `*http.request`
func NewPage(r *http.Request) Page {
	path := r.URL.Path
	params := make(map[string]string)
	r.ParseForm()

	for k, v := range r.Form {
		params[k] = v[0]
	}

	return Page{
		path:   path,
		params: params,
	}
}
