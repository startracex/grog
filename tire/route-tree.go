package tire

import (
	"sort"
	"strings"
)

const (
	MatchStrict = iota
	MatchSingle
	MatchMulti
)

func NewRouteTree() *RouteTree {
	return &RouteTree{
		children: make([]*RouteTree, 0),
	}
}

type RouteTree struct {
	Pattern  string
	part     string
	children []*RouteTree
	match    int
}

func (rt *RouteTree) InsertPattern(pattern string) {
	parts := SplitSlash(pattern)
	rt.Insert(pattern, parts, 0)
}

func (rt *RouteTree) Insert(pattern string, parts []string, height int) {
	defer rt.Sort()

	if len(parts) == height {
		rt.Pattern = pattern
		return
	}
	part := parts[height]
	spec := rt.findStrict(part)
	if spec == nil {
		spec = &RouteTree{
			part:  part,
			match: Dynamic(part).matchType,
		}
		rt.children = append(rt.children, spec)
	}
	spec.Insert(pattern, parts, height+1)
}

func (rt *RouteTree) SearchPattern(pattern string) *RouteTree {
	parts := SplitSlash(pattern)
	return rt.Search(parts, 0)
}

func (rt *RouteTree) Search(parts []string, height int) *RouteTree {
	if len(parts) == height || Dynamic(rt.part).matchType == MatchMulti {
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

func (rt *RouteTree) findStrict(part string) *RouteTree {
	for _, child := range rt.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

func (rt *RouteTree) filterWild(part string) []*RouteTree {
	nodes := make([]*RouteTree, 0)
	for _, child := range rt.children {
		if child.part == part || child.match > 0 {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (rt *RouteTree) Sort() {
	sort.Slice(rt.children, func(a, b int) bool {
		return rt.children[a].match < rt.children[b].match
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
			if strings.HasPrefix(key, "...") {
				return DynamicType{key[3:], MatchMulti}
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
