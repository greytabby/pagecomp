package pagecomp

import (
	"regexp"
	"strings"
)

// PageComparator is interface of page compare
type PageComparator interface {
	Equal(a, b Page) bool
}

// Comparator reqpresents page comperator
type Comparator struct {
	rules []*compareRule
}

type compareRule struct {
	pathRule  *regexp.Regexp
	paramRule []string
}

// NewComparator return comprator
func NewComparator() *Comparator {
	return &Comparator{
		rules: make([]*compareRule, 0),
	}
}

// Equal return whether page a and b is diffrence under the compare rule
func (c *Comparator) Equal(a, b Page) bool {
	// ルールがない場合そのままPageを比べる
	if len(c.rules) == 0 {
		return a.Equal(b)
	}

	for _, rule := range c.rules {
		if rule.Match(a) && rule.Match(b) {
			if equalParams(rule.Params(), a.params, b.params) {
				return true
			}
		}
	}
	return false
}

// AddRule add rule to comparetor
func (c *Comparator) AddRule(pattern string, params ...string) error {
	rule, err := newRule(pattern, params...)
	if err != nil {
		return err
	}
	c.rules = append(c.rules, rule)
	return nil
}

func newRule(pattern string, params ...string) (*compareRule, error) {
	parts := strings.Split(pattern, "/")

	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			parts[i] = expr
		}
	}

	pattern = strings.Join(parts, "/")
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	rule := &compareRule{
		pathRule:  regex,
		paramRule: params,
	}
	return rule, nil
}

// ruleにあるparamが2つのpageで一致しているかどうかを返す
func equalParams(ruleParams []string, a, b map[string]string) bool {
	for _, param := range ruleParams {
		if _, ok := a[param]; !ok {
			return false
		}
		if _, ok := b[param]; !ok {
			return false
		}

		if a[param] != b[param] {
			return false
		}
	}
	return true
}

// pageがruleに合致しているかどうかを返す
func (r *compareRule) Match(p Page) bool {
	for _, param := range r.paramRule {
		if _, ok := p.params[param]; !ok {
			return false
		}
	}

	if !r.pathRule.Match(([]byte(p.path))) {
		return false
	}

	return true
}

func (r *compareRule) Params() []string {
	return r.paramRule
}
