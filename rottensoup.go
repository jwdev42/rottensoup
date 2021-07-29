//This file is part of rottensoup ©2021 Jörg Walter

//Collection of functions that help navigating through an html5 DOM tree.
package rottensoup

import (
	"github.com/jwdev42/rottensoup/internal/cond"
	"github.com/jwdev42/rottensoup/internal/nav"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"regexp"
)

//Returns the corresponding attribute value if node n has an attribute of the given namespace and key.
//If no such attribute exists, an empty string will be returned.
func AttrVal(n *html.Node, namespace, key string) string {
	for _, attr := range n.Attr {
		if attr.Namespace == namespace && attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

//Executes depth-first search on all child nodes of n, returns the first element found with the given id.
//If no suitable element was found, nil will be returned.
func ElementByID(n *html.Node, id string) *html.Node {
	var elem *html.Node
	byID := func(n *html.Node) bool {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == id {
				elem = n
				return false
			}
		}
		return true
	}
	nav.DFS(n, byID, nil)
	return elem
}

//Executes depth-first search on all child nodes of n, returns the first node that matches NodeType t or nil if no such node was found.
func FirstNodeByType(n *html.Node, t html.NodeType) *html.Node {
	var match *html.Node
	pre := func(n *html.Node) bool {
		if n.Type == t {
			match = n
			return false
		}
		return true
	}
	nav.DFS(n, pre, nil)
	return match
}

//Executes depth-first search on all child nodes of n, returns a slice that contains all elements that have an attribute
//where namespace and key are equal to the function args and where the attribute's val matches the given regular expression.
//If no proper element was found, nil will be returned.
func ElementsByAttrMatch(n *html.Node, namespace, key string, val *regexp.Regexp) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	nav.DFS(n, cond.AttrValByRegex(&nodes, namespace, key, val), nil)
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

//Executes depth-first search on all child nodes of n and returns the first element that matches at least one of the given tags.
//Returns nil if no such element was found.
func FirstElementByTag(n *html.Node, tag ...atom.Atom) *html.Node {
	var node *html.Node
	pre := func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}
		for _, v := range tag {
			if n.DataAtom == v {
				node = n
				return false
			}
		}
		return true
	}
	nav.DFS(n, pre, nil)
	return node
}

//Executes depth-first search on all child nodes of n and returns all elements that match at least one of the given tags.
//Returns nil if no such element was found.
func ElementsByTag(n *html.Node, tag ...atom.Atom) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	pre := func(n *html.Node) bool {
		for _, t := range tag {
			if n.DataAtom == t {
				nodes = append(nodes, n)
				break
			}
		}
		return true
	}
	nav.DFS(n, pre, nil)
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

//Executes depth-first search on all child nodes of n and returns all elements that match
//the given tag and contain all given attributes. Returns nil if no matches were found.
func ElementsByTagAndAttr(n *html.Node, tag atom.Atom, attr ...html.Attribute) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	pre := func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.DataAtom == tag {
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
			nodes = append(nodes, n)
		}
		return true
	}
	nav.DFS(n, pre, nil)
	return nodes
}

//Returns true if node n has an attribute that matches namespace and key, returns false otherwise.
func HasAttr(n *html.Node, namespace, key string) bool {
	for _, attr := range n.Attr {
		if attr.Namespace == namespace && attr.Key == key {
			return true
		}
	}
	return false
}

//Returns true if node n contains all given attributes, returns false otherwise.
func MatchAttrs(n *html.Node, attr ...html.Attribute) bool {
	attrs_to_match := make(map[html.Attribute]bool)
	for _, a := range attr {
		attrs_to_match[a] = false
	}

	for _, a := range n.Attr {
		if _, ok := attrs_to_match[a]; ok {
			attrs_to_match[a] = true
		}
	}

	for _, v := range attrs_to_match {
		if !v {
			return false
		}
	}
	return true
}

//Returns the node's next sibling where at least one of the given tags match. Returns nil if no such node was found.
func NextSiblingByTag(n *html.Node, tag ...atom.Atom) *html.Node {
	var node *html.Node
	nav.Siblings(n, false, cond.MatchTag(&node, tag...), nil)
	return node
}

//Returns the node's next sibling that is an element. Returns nil if no such element was found.
func NextElementSibling(n *html.Node) *html.Node {
	sibl := n.NextSibling
	if sibl == nil {
		return nil
	}
	if sibl.Type == html.ElementNode {
		return sibl
	}
	return NextElementSibling(sibl)
}

//Returns the node's next previous sibling that is an element. Returns nil if no such element was found.
func PrevElementSibling(n *html.Node) *html.Node {
	sibl := n.PrevSibling
	if sibl == nil {
		return nil
	}
	if sibl.Type == html.ElementNode {
		return sibl
	}
	return PrevElementSibling(sibl)
}
