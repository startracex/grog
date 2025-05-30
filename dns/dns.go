package dns

import (
	"net"
	"strings"
)

type dnsNode[T any] struct {
	label    string
	children map[string]*dnsNode[T]
	value    T
	valid    bool
}

type DNS[T any] struct {
	root *dnsNode[T]
}

func NewDNS[T any]() *DNS[T] {
	return &DNS[T]{
		root: &dnsNode[T]{
			children: make(map[string]*dnsNode[T]),
		},
	}
}

func (t *DNS[T]) Insert(domain string, value T) {

	if domain == "*" {
		t.root.children["*"] = &dnsNode[T]{
			label:    "*",
			children: make(map[string]*dnsNode[T]),
			value:    value,
			valid:    true,
		}
		return
	}

	labels := splitDomain(domain)

	if len(labels) > 1 && labels[0] == "@" {
		t.insertLabels(labels[1:], value)
		return
	}

	t.insertLabels(labels, value)
}

func (t *DNS[T]) insertLabels(labels []string, value T) {
	current := t.root

	for i := len(labels) - 1; i >= 0; i-- {
		label := labels[i]

		if _, exists := current.children[label]; !exists {
			current.children[label] = &dnsNode[T]{
				label:    label,
				children: make(map[string]*dnsNode[T]),
			}
		}

		current = current.children[label]
	}

	current.value = value
	current.valid = true
}

func (t *DNS[T]) Match(domain string) (T, bool) {
	labels := splitDomain(domain)
	return t.matchRecursive(t.root, labels, len(labels)-1)
}

func (t *DNS[T]) matchRecursive(node *dnsNode[T], labels []string, level int) (T, bool) {
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

func splitDomain(domain string) []string {
	domain = strings.Trim(domain, ".")
	if domain == "" {
		return []string{}
	}
	return strings.Split(domain, ".")
}

func GetDomain(hostport string) string {
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
