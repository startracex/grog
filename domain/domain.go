package domain

import (
	"net"
	"strings"
)

type Node[T any] struct {
	label    string
	children map[string]*Node[T]
	value    T
	valid    bool
}

type Domain[T any] struct {
	root *Node[T]
}

func New[T any]() *Domain[T] {
	return &Domain[T]{
		root: &Node[T]{
			children: make(map[string]*Node[T]),
		},
	}
}

func (t *Domain[T]) Insert(domain string, value T) {
	if domain == "*" {
		t.root.children["*"] = &Node[T]{
			label:    "*",
			children: make(map[string]*Node[T]),
			value:    value,
			valid:    true,
		}
		return
	}

	labels := split(domain)

	if len(labels) > 1 && labels[0] == "@" {
		t.insertLabels(labels[1:], value)
		return
	}

	t.insertLabels(labels, value)
}

func (t *Domain[T]) insertLabels(labels []string, value T) {
	current := t.root

	for i := len(labels) - 1; i >= 0; i-- {
		label := labels[i]

		if _, exists := current.children[label]; !exists {
			current.children[label] = &Node[T]{
				label:    label,
				children: make(map[string]*Node[T]),
			}
		}

		current = current.children[label]
	}

	current.value = value
	current.valid = true
}

func (t *Domain[T]) Match(domain string) (T, bool) {
	labels := split(domain)
	return t.matchRecursive(t.root, labels, len(labels)-1)
}

func (t *Domain[T]) matchRecursive(node *Node[T], labels []string, level int) (T, bool) {
	if level < 0 {
		return node.value, node.valid
	}

	currentLabel := labels[level]

	if child, exists := node.children[currentLabel]; exists {
		if value, ok := t.matchRecursive(child, labels, level-1); ok {
			return value, true
		}
	}

	if child, exists := node.children["*"]; exists {
		if level == 0 {
			return child.value, true
		}
		if value, ok := t.matchRecursive(child, labels, level-1); ok {
			return value, true
		}
	}

	if child, exists := node.children["+"]; exists {
		return child.value, true
	}

	var _0 T
	return _0, false
}

func split(domain string) []string {
	domain = strings.Trim(domain, ".")
	if domain == "" {
		return []string{}
	}
	return strings.Split(domain, ".")
}

func Clean(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}
	for i := len(host) - 1; i >= 0; i-- {
		c := host[i]
		if c < '0' || c > '9' {
			return host
		}
	}
	return ""
}
