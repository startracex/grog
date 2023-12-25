package core

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	Pattern  string
	part     string
	children []*Node
	isWild   bool
}

func (n *Node) Set(s string) {
	n.Insert(s, SplitPattern(s), 0)
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
	isWild, _, _ := WildOf(part)

	child := n.findWildChild()
	if child != nil && child.isWild && isWild {
		fmt.Println("WARNING: The following routes may conflict.")
		warnConflict(child.Pattern, child.part)
		warnConflict(pattern, part)
	}

	child = n.findSpecificChild(part)
	if child == nil {
		child = &Node{part: part, isWild: isWild}
		n.children = append(n.children, child)
	}

	child.Insert(pattern, parts, height+1)
}

func (n *Node) Search(parts []string, height int) *Node {
	_, _, isMulti := WildOf(n.part)
	if len(parts) == height || isMulti {
		if n.Pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.Search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *Node) findWildChild() *Node {
	for _, child := range n.children {
		if child.isWild {
			return child
		}
	}
	return nil
}

func (n *Node) findSpecificChild(part string) *Node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// SplitPattern split slash rune
func SplitPattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	var parts []string
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			_, _, isMulti := WildOf(item)
			if isMulti {
				break
			}
		}
	}
	return parts
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

// WildOf returns is wild (bool), wild key (string), is multi (bool)
func WildOf(s string) (bool, string, bool) {
	if len(s) == 0 {
		return false, "", false
	}
	if s[0] == '*' {
		return true, s[1:], true
	}
	if s[0] == ':' {
		return true, s[1:], false
	}
	if s[0] == '{' && s[len(s)-1] == '}' {
		return true, s[1 : len(s)-1], false
	}
	if strings.HasPrefix(s, "...") {
		return true, s[3:], true
	}
	return false, "", false
}

func (n *Node) Sort() {
	if n == nil {
		return
	}
	list := n.children
	sort.Slice(n.children, func(i, j int) bool {
		if !n.children[i].isWild && n.children[j].isWild {
			return true
		} else if n.children[i].isWild && !n.children[j].isWild {
			return false
		} else {
			return len(n.children[i].Pattern) < len(n.children[j].Pattern)
		}
	})
	if len(list) > 0 {
		for i := range list {
			list[i].Sort()
		}
	}
}

func warnConflict(s, h string) {
	indexOf := strings.Index(s, h)
	lengthOf := len(h)
	fmt.Println("  " + s)
	fmt.Print("  ")
	for indexOf > 0 {
		fmt.Print(" ")
		indexOf--
	}
	for lengthOf > 0 {
		fmt.Print("^")
		lengthOf--
	}
	fmt.Println()
}
