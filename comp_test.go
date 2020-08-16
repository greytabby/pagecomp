package pagecomp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestComperatorEqualWithNoRule(t *testing.T) {
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

func TestComparatorEqualWithPathRule(t *testing.T) {
	c := NewComparator()
	err := c.AddPathRule("/a/b/:param1")
	assert.NoError(t, err)
	err = c.AddPathRule("/a/b/:param1/c/:param2")
	assert.NoError(t, err)

	cases := []struct {
		name string
		p1   Page
		p2   Page
		want bool
	}{
		{name: "path同じ,param同じ", p1: Page{path: "/a/b/param", params: pageParams("x", "y")}, p2: Page{path: "/a/b/param", params: pageParams("x", "y")}, want: true},
		{name: "path同じ,param不一致", p1: Page{path: "/a/b/param", params: pageRandomParams("x")}, p2: Page{path: "/a/b/param", params: pageRandomParams("x")}, want: true},
		{name: "pathパラメータが不一致", p1: Page{path: "/a/b/param1", params: nil}, p2: Page{path: "/a/b/param2", params: nil}, want: true},
		{name: "pathが不一致1", p1: Page{path: "/a/b/param1", params: nil}, p2: Page{path: "/a/b/param1/c/param2", params: nil}, want: false},
		{name: "pathが不一致2", p1: Page{path: "/a", params: nil}, p2: Page{path: "/b", params: nil}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := c.Equal(tc.p1, tc.p2)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestComparatorEqualWithParamRule(t *testing.T) {
	c := NewComparator()
	c.AddParamRule("page_id")

	cases := []struct {
		name string
		p1   Page
		p2   Page
		want bool
	}{
		{name: "path同じ,param同じ", p1: Page{path: "/a", params: pageParams("page_id")}, p2: Page{path: "/a", params: pageParams("page_id")}, want: true},
		{name: "path同じ,param不一致", p1: Page{path: "/a", params: pageRandomParams("page_id")}, p2: Page{path: "/a", params: pageRandomParams("page_id")}, want: false},
		{name: "paramなし", p1: Page{path: "/a", params: nil}, p2: Page{path: "/a", params: nil}, want: true},
		{name: "pathが不一致,paramが同じ", p1: Page{path: "/a", params: pageParams("page_id")}, p2: Page{path: "/b", params: pageParams("page_id")}, want: false},
		{name: "pathが不一致,paramが不一致", p1: Page{path: "/a", params: pageRandomParams("page_id")}, p2: Page{path: "/b", params: pageRandomParams("page_id")}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := c.Equal(tc.p1, tc.p2)
			assert.Equal(t, tc.want, got)
		})
	}
}
func TestComparatorEqualWithNoRule(t *testing.T) {

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

func ExampleEqual() {
	c := NewComparator()
	err := c.AddPathRule("/a/:id")
	if err != nil {
		panic(err)
	}
	c.AddParamRule("x")

	p1 := Page{path: "/a/abc", params: map[string]string{"x": "example"}}
	p2 := Page{path: "/a/1234567", params: map[string]string{"x": "example"}}
	result := c.Equal(p1, p2)
	fmt.Println(result)

	// Output:
	// true
}

func TestRegex(t *testing.T) {
	regex := regexp.MustCompile(`/a/b/[^/]+`)
	t.Log(regex.MatchString("/a/b/test"))
	t.Log(regex.MatchString("/a/b/test/c/d"))
}
