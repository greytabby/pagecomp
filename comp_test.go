package pagecomp

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestComperatorEqual(t *testing.T) {
	cases := []struct {
		name string
		p1   Page
		p2   Page
		want bool
	}{
		{name: "params一致、path完全一致", p1: Page{path: "/a", params: pageParams("x", "y")}, p2: Page{path: "/a", params: pageParams("x", "y")}, want: true},
		{name: "paramsなし、path完全一致", p1: Page{path: "/a", params: nil}, p2: Page{path: "/a", params: nil}, want: true},
		{name: "params不一致、path完全一致", p1: Page{path: "/a", params: pageRandomParams("x", "y")}, p2: Page{path: "/a", params: pageRandomParams("x", "y")}, want: false},
		{name: "params不一致、path不一致", p1: Page{path: "/a", params: pageRandomParams("x")}, p2: Page{path: "/b", params: pageRandomParams("x")}, want: false},
	}

	c := NewComparator()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := c.Equal(tc.p1, tc.p2)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEqualParams(t *testing.T) {
	cases := []struct {
		name  string
		rules []string
		p1    map[string]string
		p2    map[string]string
		want  bool
	}{
		{name: "ルールなし;Param不一致", rules: nil, p1: pageParams("a"), p2: pageParams("b"), want: true},
		{name: "ルールあり;Param不一致", rules: []string{"a", "b"}, p1: pageRandomParams("a", "b"), p2: pageRandomParams("a", "b"), want: false},
		{name: "ルールあり;Param一致", rules: []string{"a", "b"}, p1: pageParams("a", "b"), p2: pageParams("a", "b"), want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := equalParams(tc.rules, tc.p1, tc.p2)
			assert.Equal(t, tc.want, got)
		})
	}
}

// 指定されたkeyに固定の値が入ったmapを返す
func pageParams(keys ...string) map[string]string {
	params := make(map[string]string)
	for _, key := range keys {
		params[key] = "test"
	}
	return params
}

// 指定されたkeyにランダムな値が入ったmapを返す
func pageRandomParams(keys ...string) map[string]string {
	params := make(map[string]string)
	for _, key := range keys {
		params[key] = uuid.New().String()
	}
	return params
}

func TestMutch(t *testing.T) {
	rule, err := newRule("/a/b/c/:id/d/:bbb/a", "x", "y")
	assert.NoError(t, err)

	cases := []struct {
		name string
		page Page
		want bool
	}{
		{name: "Should match", page: Page{path: "/a/b/c/id/d/bbb/a", params: pageParams("x", "y")}, want: true},
		{name: "not path match", page: Page{path: "/a/b/c/id/d/a", params: pageParams("x", "y")}, want: false},
		{name: "not param match", page: Page{path: "/a/b/c/id/d/bbb/a", params: pageParams("x")}, want: false},
		{name: "path match", page: Page{path: "/a/b/c/xxxxxxxxx/d/bbbbbbbbbbbb/a", params: pageParams("x", "y")}, want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.Match(tc.page)
			assert.Equal(t, tc.want, got)
		})
	}
}

func ExampleEqual() {
	c := NewComparator()
	err := c.AddRule("/a/:id", "x")
	if err != nil {
		panic(err)
	}

	p1 := Page{path: "/a/abc", params: map[string]string{"x": "example"}}
	p2 := Page{path: "/a/1234567", params: map[string]string{"x": "example"}}
	result := c.Equal(p1, p2)
	fmt.Println(result)

	// Output:
	// true
}
