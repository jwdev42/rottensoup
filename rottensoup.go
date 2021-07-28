//This file is part of rottensoup ©2021 Jörg Walter

package rottensoup

import (
	"github.com/jwdev42/rottensoup/internal/cond"
	"github.com/jwdev42/rottensoup/internal/nav"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"regexp"
)

//Returns the attribute value for node n's attribute referenced by key key. Returns an empty string if no such attribute exists.
func AttrVal(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

//Returns the first child node of n that has an id attribute that matches the given id string. Returns nil if no such node exists.
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

//Returns the first child node of n that matches NodeType t or nil if no such node was found.
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

func ElementsByAttrMatch(n *html.Node, key string, val *regexp.Regexp) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	nav.DFS(n, cond.AttrValByRegex(&nodes, key, val), nil)
	return nodes
}

//Returns the first child node of n that matches at least one of the given tags.
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

//Returns all child nodes of n that match at least one of the given tags.
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
	return nodes
}

//Returns all child nodes of n that match tag tag and contain all attributes in attr.
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

//Returns true if Node n has an attribute with key key, returns false otherwise.
func HasAttr(n *html.Node, key string) bool {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return true
		}
	}
	return false
}

//Returns true if Node n contains all attributes that are passed to attr, returns false otherwise.
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
