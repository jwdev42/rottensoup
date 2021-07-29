//This file is part of rottensoup ©2021 Jörg Walter

package cond

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"regexp"
)

//Returns a function that collects all nodes who have an attribute that matches key and where its value matches the regex val.
func AttrValByRegex(nodes *[]*html.Node, namespace, key string, val *regexp.Regexp) func(n *html.Node) bool {
	return func(node *html.Node) bool {
		for _, a := range node.Attr {
			if a.Namespace == namespace && a.Key == key && val.MatchString(a.Val) {
				*nodes = append(*nodes, node)
				return true
			}
		}
		return true
	}
}

//Returns a function that adds every given node to nodes if it contains all attributes specified in attr.
//If first is true, search will stop after first match.
func MatchAttrs(nodes *[]*html.Node, first bool, attr ...html.Attribute) func(*html.Node) bool {
	return func(n *html.Node) bool {
		for _, a := range attr {
			found := false
			for _, na := range n.Attr {
				if a == na {
					found = true
					break
				}
			}
			if !found {
				return true
			}
		}
		*nodes = append(*nodes, n)
		if first {
			return false
		}
		return true
	}
}

func MatchTag(nodes *[]*html.Node, first bool, tag ...atom.Atom) func(*html.Node) bool {
	return func(node *html.Node) bool {
		for _, t := range tag {
			if node.DataAtom == t {
				*nodes = append(*nodes, node)
				if first {
					return false
				}
				return true
			}
		}
		return true
	}
}

/* --- Filters --- */

func TagFilter(tag atom.Atom, f func(*html.Node) bool) func(*html.Node) bool {
	return func(n *html.Node) bool {
		if n.DataAtom != tag {
			return true
		}
		return f(n)
	}
}

func TypeFilter(t html.NodeType, f func(*html.Node) bool) func(*html.Node) bool {
	return func(n *html.Node) bool {
		if n.Type != t {
			return true
		}
		return f(n)
	}
}
