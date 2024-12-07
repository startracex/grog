package core

import (
	"sort"
	"strings"
)

const (
	matchStrict = iota
	matchSingle
	matchMulti
)

type Node struct {
	Pattern  string
	part     string
	children []*Node
	match    int
	sorted   bool
}

func (n *Node) Set(s string) {
	n.Insert(s, SplitSlash(s), 0)
}

func (n *Node) Get(s string) *Node {
	return n.Search(SplitSlash(s), 0)
}

func (n *Node) Insert(pattern string, parts []string, height int) {
	defer n.Sort()

	if len(parts) == height {
		n.Pattern = pattern
		return
	}
	part := parts[height]
	spec := n.findStrict(part)
	if spec == nil {
		spec = &Node{
			part:     part,
			match:    Dynamic(part).matchType,
		}
		n.children = append(n.children, spec)
		n.sorted = false
	}
	spec.Insert(pattern, parts, height+1)
}

func (n *Node) Search(parts []string, height int) *Node {
	if !n.sorted {
		n.Sort()
	}

	if len(parts) == height || Dynamic(n.part).matchType == matchMulti {
		if n.Pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.filterWild(part)
	for _, child := range children {
		result := child.Search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *Node) findStrict(part string) *Node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

func (n *Node) filterWild(part string) []*Node {
	nodes := make([]*Node, 0)
	for _, child := range n.children {
		if child.part == part || child.match > 0 {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *Node) Sort() {
	if n == nil {
		return
	}
	sort.Slice(n.children, func(a, b int) bool {
		return n.children[a].match < n.children[b].match
	})
	for _, child := range n.children {
		child.Sort()
	}
	n.sorted = true
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
