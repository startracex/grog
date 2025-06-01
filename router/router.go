package router

import (
	"sort"
	"strings"
)

const (
	MatchStrict = iota
	MatchSingle
	MatchMulti
)

func NewRouter[T any]() *Router[T] {
	return &Router[T]{
		children: make([]*Router[T], 0),
	}
}

type Router[T any] struct {
	Pattern  string
	Part     string
	children []*Router[T]
	Match    int
	Value    T
}

func (rt *Router[T]) Insert(pattern string, parts []string, height int, value T) {
	defer rt.Sort()

	if len(parts) == height {
		rt.Pattern = pattern
		rt.Value = value
		return
	}
	part := parts[height]
	spec := rt.findStrict(part)
	if spec == nil {
		spec = &Router[T]{
			Part:  part,
			Match: Dynamic(part).matchType,
		}
		rt.children = append(rt.children, spec)
	}
	spec.Insert(pattern, parts, height+1, value)
}

func (rt *Router[T]) Search(parts []string, height int) *Router[T] {
	if len(parts) == height || Dynamic(rt.Part).matchType == MatchMulti {
		if rt.Pattern == "" {
			return nil
		}
		return rt
	}
	part := parts[height]
	children := rt.filterWild(part)
	for _, child := range children {
		result := child.Search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (rt *Router[T]) findStrict(part string) *Router[T] {
	for _, child := range rt.children {
		if child.Part == part {
			return child
		}
	}
	return nil
}

func (rt *Router[T]) filterWild(part string) []*Router[T] {
	nodes := make([]*Router[T], 0)
	for _, child := range rt.children {
		if child.Part == part || child.Match > 0 {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (rt *Router[T]) Sort() {
	sort.Slice(rt.children, func(a, b int) bool {
		return rt.children[a].Match < rt.children[b].Match
	})
	for _, child := range rt.children {
		child.Sort()
	}
}

type DynamicType struct {
	key       string
	matchType int
}

func Dynamic(key string) DynamicType {
	if len(key) > 0 {
		if affix(key, "{", "}") || affix(key, "[", "]") {
			key = key[1 : len(key)-1]
			result := Dynamic(key)
			if result.matchType == MatchStrict {
				result.matchType = MatchSingle
			}
			return result
		}
		if len(key) > 1 {
			a := key[0]
			if a == ':' {
				return DynamicType{key[1:], MatchSingle}
			}
			if a == '*' {
				return DynamicType{key[1:], MatchMulti}
			}
			if strings.HasPrefix(key, "...") || strings.HasSuffix(key, "...") {
				return DynamicType{strings.TrimPrefix(strings.TrimSuffix(key, "..."), "..."), MatchMulti}
			}
		}
	}
	return DynamicType{
		key:       key,
		matchType: MatchStrict,
	}
}

func affix(key, prefix, suffix string) bool {
	return strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix)
}

func ParseParams(path string, pattern string) map[string]string {
	pathSplit := SplitSlash(path)
	patternSplit := SplitSlash(pattern)
	params := make(map[string]string)
	for i := range patternSplit {
		info := Dynamic(patternSplit[i])
		if info.matchType == MatchSingle {
			params[info.key] = pathSplit[i]
		} else if info.matchType == MatchMulti {
			params[info.key] = strings.Join(pathSplit[i:], "/")
			break
		}
	}
	return params
}

// SplitSlash split slash rune
func SplitSlash(s string) []string {
	vs := strings.Split(s, "/")
	var parts []string
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
		}
	}
	return parts
}
