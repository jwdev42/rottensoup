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

func MatchTag(n **html.Node, tag ...atom.Atom) func(*html.Node) bool {
	return func(node *html.Node) bool {
		if node.Type != html.ElementNode {
			return true
		}
		for _, t := range tag {
			if node.DataAtom == t {
				*n = node
				return false
			}
		}
		return true
	}
}
