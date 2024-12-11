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
			match: dynamic(part).matchType,
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

	if len(parts) == height || dynamic(rt.part).matchType == matchMulti {
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

type dynamicInfo struct {
	key       string
	carry     int
	multi     bool
	matchType int
}

func dynamic(key string) dynamicInfo {
	if len(key) > 0 {
		if affix(key, "{", "}") {
			key = key[1 : len(key)-1]
			if affix(key, "[", "]") {
				return dynamic(key)
			} else {
				return dynamicInfo{key, 0, true, matchSingle}
			}
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

func affix(key, prefix, suffix string) bool {
	return strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix)
}

func ParseParams(path string, pattern string) map[string]string {
	pathSplit := SplitSlash(path)
	patternSplit := SplitSlash(pattern)
	params := make(map[string]string)
	for i := range patternSplit {
		info := dynamic(patternSplit[i])
		if info.key == "" {
			continue
		}
		if !info.multi {
			params[info.key[info.carry:]] = pathSplit[i]
		} else {
			params[info.key[info.carry:]] = strings.Join(pathSplit[i:], "/")
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
