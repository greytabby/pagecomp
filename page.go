package pagecomp

import "reflect"

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
