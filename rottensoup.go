//This file is part of rottensoup ©2021 Jörg Walter

package rottensoup

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"regexp"
)

func walkTree(n *html.Node, pre, post func(*html.Node) bool) bool {
	if pre != nil && !pre(n) {
		return false
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !walkTree(c, pre, post) {
			return false
		}
	}
	if post != nil && !post(n) {
		return false
	}
	return true
}

//BEGIN: Functions to be used as pre or post with walktree

//matchAttrVal add Node n and child nodes to the nodecollection nc if the regex val matches
//the value of the attribute specified with key
func matchAttrVal(nodes *[]*html.Node, key string, val *regexp.Regexp) func(n *html.Node) bool {
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

//END: Functions to be used as pre or post with walktree*

func AttrVal(node *html.Node, attribute string) string {
	for _, attr := range node.Attr {
		if attr.Key == attribute {
			return attr.Val
		}
	}
	return ""
}

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
	walkTree(n, byID, nil)
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
	walkTree(n, pre, nil)
	return match
}

func ElementsByAttrMatch(n *html.Node, key string, val *regexp.Regexp) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	walkTree(n, matchAttrVal(&nodes, key, val), nil)
	return nodes
}

//Returns the first child node of n that matches tag tag.
func FirstElementByTag(n *html.Node, tag ...atom.Atom) *html.Node {
	elements := ElementsByTag(n, tag...)
	if len(elements) < 1 {
		return nil
	}
	return elements[0]
}

//Returns all child nodes of n that match tag tag.
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
	walkTree(n, pre, nil)
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
	walkTree(n, pre, nil)
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

//Returns the node's next sibling which tag matches tag. Returns nil if no such node was found.
func NextSiblingByTag(n *html.Node, tag atom.Atom) *html.Node {
	sibl := n.NextSibling
	if sibl == nil {
		return nil
	}
	if sibl.DataAtom == tag {
		return sibl
	}
	return NextSiblingByTag(sibl, tag)
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
