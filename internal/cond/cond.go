package cond

import (
	"golang.org/x/net/html"
	"regexp"
)

//Returns a function that collects all nodes who have an attribute that matches key and where its value matches the regex val.
func AttrValByRegex(nodes *[]*html.Node, key string, val *regexp.Regexp) func(n *html.Node) bool {
	return func(node *html.Node) bool {
		for _, a := range node.Attr {
			if a.Key == key && val.MatchString(a.Val) {
				*nodes = append(*nodes, node)
				return true
			}
		}
		return true
	}
}
