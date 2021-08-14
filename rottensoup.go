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
	nodes := make([]*html.Node, 0, 1)
	attr := html.Attribute{Key: "id", Val: id}
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.MatchAttrs(&nodes, true, attr)), nil)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

//Executes depth-first search on all child nodes of n, returns the first node that matches NodeType t or nil if no such node was found.
func FirstNodeByType(n *html.Node, t html.NodeType) *html.Node {
	var match *html.Node
	pre := func(n *html.Node) bool {
		match = n
		return false
	}
	nav.DFS(n, cond.TypeFilter(t, pre), nil)
	return match
}

//Executes depth-first search on all child nodes of n and returns all elements that
//contain all given attributes. Returns nil if no matches were found.
func ElementsByAttr(n *html.Node, attr ...html.Attribute) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.MatchAttrs(&nodes, false, attr...)), nil)
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

//Executes depth-first search on all child nodes of n, returns a slice that contains all elements that have an attribute
//where namespace and key are equal to the function args and where the attribute's val matches the given regular expression.
//If no proper element was found, nil will be returned.
func ElementsByAttrMatch(n *html.Node, namespace, key string, val *regexp.Regexp) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.AttrValByRegex(&nodes, namespace, key, val)), nil)
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

//Executes depth-first search on all child nodes of n and returns the first element that
//contains all given attributes. Returns nil if no match was found.
func FirstElementByAttr(n *html.Node, attr ...html.Attribute) *html.Node {
	nodes := make([]*html.Node, 0, 1)
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.MatchAttrs(&nodes, true, attr...)), nil)
	if len(nodes) < 1 {
		return nil
	}
	return nodes[0]
}

func FirstElementByClassName(n *html.Node, name ...string) *html.Node {
	nodes := make([]*html.Node, 0, 1)
	nav.DFS(n, cond.MatchClassNames(&nodes, true, name...), nil)
	if len(nodes) < 1 {
		return nil
	}
	return nodes[0]
}

//Executes depth-first search on all child nodes of n and returns the first element that matches at least one of the given tags.
//Returns nil if no such element was found.
func FirstElementByTag(n *html.Node, tag ...atom.Atom) *html.Node {
	nodes := make([]*html.Node, 0, 1)
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.MatchTag(&nodes, true, tag...)), nil)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

//Executes depth-first search on all child nodes of n and returns the first element that matches
//the given tag and contains all given attributes. Returns nil if no match was found.
func FirstElementByTagAndAttr(n *html.Node, tag atom.Atom, attr ...html.Attribute) *html.Node {
	nodes := make([]*html.Node, 0, 1)
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.TagFilter(tag, cond.MatchAttrs(&nodes, true, attr...))), nil)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

func ElementsByClassName(n *html.Node, name ...string) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	nav.DFS(n, cond.MatchClassNames(&nodes, false, name...), nil)
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

//Executes depth-first search on all child nodes of n and returns all elements that match at least one of the given tags.
//Returns nil if no such element was found.
func ElementsByTag(n *html.Node, tag ...atom.Atom) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.MatchTag(&nodes, false, tag...)), nil)
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

//Executes depth-first search on all child nodes of n and returns all elements that match
//the given tag and contain all given attributes. Returns nil if no matches were found.
func ElementsByTagAndAttr(n *html.Node, tag atom.Atom, attr ...html.Attribute) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	nav.DFS(n, cond.TypeFilter(html.ElementNode, cond.TagFilter(tag, cond.MatchAttrs(&nodes, false, attr...))), nil)
	if len(nodes) == 0 {
		return nil
	}
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
	nodes := make([]*html.Node, 0, 1)
	cond.MatchAttrs(&nodes, true, attr...)(n)
	if len(nodes) > 0 {
		return true
	}
	return false
}

//Returns the node's next sibling where at least one of the given tags match. Returns nil if no such node was found.
func NextSiblingByTag(n *html.Node, tag ...atom.Atom) *html.Node {
	nodes := make([]*html.Node, 0, 1)
	nav.Siblings(n, false, cond.TypeFilter(html.ElementNode, cond.MatchTag(&nodes, true, tag...)), nil)
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
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
