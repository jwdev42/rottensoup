package nav

import (
	"golang.org/x/net/html"
)

//Perform depth-first search on child nodes of n.
func DFS(n *html.Node, pre, post func(*html.Node) bool) bool {
	if pre != nil && !pre(n) {
		return false
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !DFS(c, pre, post) {
			return false
		}
	}
	if post != nil && !post(n) {
		return false
	}
	return true
}

func Siblings(n *html.Node, reverse bool, pre, post func(*html.Node) bool) {
	if reverse {
		siblings(n.PrevSibling, reverse, pre, post)
	} else {
		siblings(n.NextSibling, reverse, pre, post)
	}
}

func siblings(n *html.Node, reverse bool, pre, post func(*html.Node) bool) {
	if n == nil {
		return
	}

	if pre != nil && !pre(n) {
		return
	}

	var next *html.Node
	if reverse {
		next = n.PrevSibling
	} else {
		next = n.NextSibling
	}

	siblings(next, reverse, pre, post)

	if post != nil {
		post(n)
	}
}
