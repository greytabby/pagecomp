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
	pathRules  []*regexp.Regexp
	paramRules []string
}

// NewComparator return comprator
func NewComparator() *Comparator {
	return &Comparator{
		pathRules:  make([]*regexp.Regexp, 0),
		paramRules: make([]string, 0),
	}
}

// Equal return whether page a and b is diffrence under the compare rule
func (c *Comparator) Equal(a, b Page) bool {
	// ルールがない場合そのままPageを比べる
	if len(c.pathRules) == 0 && len(c.paramRules) == 0 {
		return a.Equal(b)
	}

	if equalPath(a, b) {
		return equalParams(c.paramRules, a.params, b.params)
	}

	// pathが異なっているかどうか
	matchPath := false
	for _, regex := range c.pathRules {
		if regex.MatchString(a.path) && regex.MatchString(b.path) {
			matchPath = true
			break
		}
	}

	// paramが異なっているかどうか
	matchParam := equalParams(c.paramRules, a.params, b.params)

	return matchPath && matchParam
}

// AddPathRule add path paramater rule to comparator
func (c *Comparator) AddPathRule(pattern string) error {
	regex, err := pathPatternToRegex(pattern)
	if err != nil {
		return err
	}
	c.pathRules = append(c.pathRules, regex)
	return nil
}

// AddPathRule add query, post, put, delete paramater rule to comparator
func (c *Comparator) AddParamRule(paramKey string) {
	c.paramRules = append(c.paramRules, paramKey)
}

func pathPatternToRegex(pattern string) (*regexp.Regexp, error) {
	parts := strings.Split(pattern, "/")

	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "[^/]+"
			parts[i] = expr
		}
	}

	pattern = "^" + strings.Join(parts, "/") + "$"
	return regexp.Compile(pattern)
}

// ruleにあるparamが2つのpageで一致しているかどうかを返す
func equalParams(ruleParams []string, a, b map[string]string) bool {
	if !hasRuleKey(ruleParams, a, b) {
		return true
	}
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

func hasRuleKey(ruleParams []string, a, b map[string]string) bool {
	for _, param := range ruleParams {
		if _, ok := a[param]; ok {
			return true
		}
		if _, ok := b[param]; ok {
			return true
		}
	}
	return false
}

func equalPath(a, b Page) bool {
	return a.path == b.path
}
