package goup

import (
	"sort"
	"strings"
)

const (
	matchStrict = iota
	matchSingle
	matchMulti
)

type RouteTree struct {
	Pattern  string
	part     string
	children []*RouteTree
	match    int
	sorted   bool
}

func (rt *RouteTree) Set(s string) {
	rt.Insert(s, SplitSlash(s), 0)
}

func (rt *RouteTree) Get(s string) *RouteTree {
	return rt.Search(SplitSlash(s), 0)
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
		rt.sorted = false
	}
	spec.Insert(pattern, parts, height+1)
}

func (rt *RouteTree) Search(parts []string, height int) *RouteTree {
	if !rt.sorted {
		rt.Sort()
	}

	if len(parts) == height || Dynamic(rt.part).matchType == matchMulti {
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
	if rt == nil {
		return
	}
	sort.Slice(rt.children, func(a, b int) bool {
		return rt.children[a].match < rt.children[b].match
	})
	for _, child := range rt.children {
		child.Sort()
	}
	rt.sorted = true
}

// func warnConflict(s, h string) {
// 	indexOf := strings.Index(s, h)
// 	lengthOf := len(h)
// 	fmt.Println("  " + s)
// 	fmt.Print("  ")
// 	for indexOf > 0 {
// 		fmt.Print(" ")
// 		indexOf--
// 	}
// 	for lengthOf > 0 {
// 		fmt.Print("^")
// 		lengthOf--
// 	}
// 	fmt.Println()
// }

type dynamicInfo struct {
	key       string
	carry     int
	multi     bool
	matchType int
}

func Dynamic(key string) dynamicInfo {
	if len(key) > 0 {
		if (strings.HasPrefix(key, "{") && strings.HasSuffix(key, "}")) || strings.HasPrefix(key, "[") && strings.HasSuffix(key, "]") {
			key = key[1 : len(key)-1]
			return Dynamic(key)
		}
		if len(key) > 1 {
			a := key[0]
			if a == ':' {
				return dynamicInfo{key, 1, false, matchSingle}
			}
			if a == '*' {
				return dynamicInfo{key, 1, true, matchMulti}
			}
			if strings.HasPrefix(key, "...") {
				return dynamicInfo{key, 3, true, matchMulti}
			}
		}
	}
	return dynamicInfo{
		key:       "",
		carry:     0,
		multi:     false,
		matchType: matchStrict,
	}
}

func ParseParams(path string, pattern string) map[string]string {
	pathSpit := SplitSlash(path)
	patternSpit := SplitSlash(pattern)
	params := make(map[string]string)
	for i := range patternSpit {
		info := Dynamic(patternSpit[i])
		if info.key == "" {
			continue
		}
		params[info.key[info.carry:]] = pathSpit[i]
		if info.multi {
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
