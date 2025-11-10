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

type Router[T any] struct {
	Part     string
	Match    uint8
	Pattern  string
	Value    T
	Children []*Router[T]
}

func NewRouter[T any]() *Router[T] {
	return &Router[T]{}
}

func (r *Router[T]) Insert(pattern string, value T) {
	r.insert(pattern, pattern, value)
	r.Sort()
}

func (r *Router[T]) Search(path string) *Router[T] {
	return r.search(path)
}

func (r *Router[T]) insert(path, pattern string, value T) {
	if path == "" {
		r.Pattern = pattern
		r.Value = value
		return
	}

	part, remaining := nextPart(path)
	child := r.findChild(part)
	if child == nil {
		child = &Router[T]{
			Part:  part,
			Match: Dynamic(part).MatchType,
		}
		r.Children = append(r.Children, child)
	}
	child.insert(remaining, pattern, value)
}

func (r *Router[T]) search(path string) *Router[T] {
	if path == "" {
		if r.Pattern != "" {
			return r
		}
		return nil
	}

	part, remaining := nextPart(path)
	for _, child := range r.Children {
		switch child.Match {
		case MatchStrict:
			if child.Part == part {
				if result := child.search(remaining); result != nil {
					return result
				}
			}
		case MatchSingle:
			if result := child.search(remaining); result != nil {
				return result
			}
		case MatchMulti:
			return child
		}
	}
	return nil
}

func (r *Router[T]) findChild(part string) *Router[T] {
	for _, child := range r.Children {
		if child.Part == part {
			return child
		}
	}
	return nil
}

func (r *Router[T]) Sort() {
	sort.SliceStable(r.Children, func(i, j int) bool {
		return r.Children[i].Match < r.Children[j].Match
	})

	for _, child := range r.Children {
		child.Sort()
	}
}

func nextPart(path string) (part, remaining string) {
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		return "", ""
	}

	idx := strings.IndexByte(path, '/')
	if idx == -1 {
		return path, ""
	}
	return path[:idx], path[idx+1:]
}

type DynamicType struct {
	Key       string
	MatchType uint8
}

func Dynamic(key string) DynamicType {
	if len(key) > 1 {
		first, last := key[0], key[len(key)-1]
		if (first == '{' && last == '}') || (first == '[' && last == ']') {
			key = key[1 : len(key)-1]
			result := Dynamic(key)
			if result.MatchType == MatchStrict {
				result.MatchType = MatchSingle
			}
			return result
		}
		if first == ':' {
			return DynamicType{key[1:], MatchSingle}
		}
		if first == '*' {
			return DynamicType{key[1:], MatchMulti}
		}
		if strings.HasPrefix(key, "...") {
			return DynamicType{key[3:], MatchMulti}
		}
		if strings.HasSuffix(key, "...") {
			return DynamicType{key[:len(key)-3], MatchMulti}
		}
	}
	return DynamicType{
		Key:       key,
		MatchType: MatchStrict,
	}
}

func ParseParams(path, pattern string) map[string]string {
	params := make(map[string]string)
	var pathPart, patternPart string

	for {
		pathPart, path = nextPart(path)
		patternPart, pattern = nextPart(pattern)

		if patternPart == "" {
			break
		}
		if pathPart == "" {
			break
		}

		info := Dynamic(patternPart)
		switch info.MatchType {
		case MatchSingle:
			params[info.Key] = pathPart
		case MatchMulti:
			params[info.Key] = pathPart
			if path != "" {
				params[info.Key] += "/" + path
			}
			return params
		}
	}
	return params
}
