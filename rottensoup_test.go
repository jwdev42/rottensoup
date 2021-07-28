//This file is part of rottensoup ©2021 Jörg Walter

package rottensoup

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

const htmlDir = "test" //directory where the testfiles are at.

func parseTestFile(name string) (*html.Node, error) {
	path := filepath.Join(htmlDir, name)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot open test file: %s", err)
	}
	defer f.Close()
	return html.Parse(f)
}

func TestMatchAttrs(t *testing.T) {

	nodeattr := make([]html.Attribute, 0, 2)
	nodeattr = append(nodeattr, html.Attribute{"test", "id", "1337"})
	nodeattr = append(nodeattr, html.Attribute{"", "src", "https://example.net/image.jpg"})
	n := &html.Node{Data: "test", Attr: nodeattr}

	musthave := make([]html.Attribute, len(nodeattr))
	copy(musthave, nodeattr)

	if !MatchAttrs(n, musthave...) {
		t.Errorf("%v and %v were expected to be equal", n.Attr, musthave)
	}

	musthave = append(musthave, html.Attribute{"", "alt", "test"})

	if MatchAttrs(n, musthave...) {
		t.Errorf("%v and %v were not expected to be equal", n.Attr, musthave)
	}

	empty := make([]html.Attribute, 0)
	if !MatchAttrs(n, empty...) {
		t.Error("A match against an empty slice should always return true")
	}
}

func TestNextElementSibling(t *testing.T) {
	const testDoc = "test.html"
	const testID = "siblings"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}
	ul := ElementByID(root, testID)
	if ul == nil {
		t.Errorf("No Element found with id \"%s\"!", testID)
	}
	var li [3]*html.Node

	li[0] = FirstElementByTag(ul, atom.Li)
	li[1] = NextElementSibling(li[0])
	li[2] = NextElementSibling(li[1])

	for i, e := range li {
		text := FirstNodeByType(e, html.TextNode)
		expect := fmt.Sprintf("Sibling %d", i+1)
		if text == nil {
			t.Error("Expected text node")
			continue
		}
		if text.Data != expect {
			t.Errorf("Expected \"%s\", got \"%s\"", expect, text.Data)
		}
	}
}

func TestPrevElementSibling(t *testing.T) {
	const testDoc = "test.html"
	const testID = "pre2"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}

	e := ElementByID(root, testID)
	if e == nil {
		t.Fatal("No start element found for testing")
	}

	expectTag := func(n *html.Node, tag atom.Atom) {
		if n.DataAtom != tag {
			t.Fatalf("Previous sibling: Expected node of tag \"%s\", got tag \"%s\"", tag.String(), n.DataAtom.String())
		}
	}

	prev := PrevElementSibling(e)
	if prev == nil {
		t.Fatal("No previous sibling found")
	}

	expectTag(prev, atom.Br)
	prev = PrevElementSibling(prev)
	expectTag(prev, atom.A)

	const expectedHref = "https://google.com"
	if AttrVal(prev, "href") != expectedHref {
		t.Fatalf("Expected \"%s\", got \"%s\"", expectedHref, AttrVal(prev, "href"))
	}
}

func TestElementsByAttrMatch(t *testing.T) {
	const matches = 4
	const testDoc = "attr_match.html"

	search := regexp.MustCompile("caption-[a-z]+")

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}

	res := ElementsByAttrMatch(root, "class", search)
	if len(res) != matches {
		t.Errorf("Expected %d matches, got %d", matches, len(res))
	}
	for i, e := range res {
		text := FirstNodeByType(e, html.TextNode)
		expect := fmt.Sprintf("Match %d", i+1)
		if text == nil {
			t.Error("Expected text node")
			continue
		}
		if text.Data != expect {
			t.Errorf("Expected \"%s\", got \"%s\"", expect, text.Data)
		}
	}
}

func TestNextSiblingByTag(t *testing.T) {
	const testDoc = "test.html"
	const testID = "TestNextSiblingByTag"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}
	parent := ElementByID(root, testID)

	expect := func(start *html.Node, key, val string, tag ...atom.Atom) {
		res := NextSiblingByTag(start, tag...)
		if res == nil {
			t.Error("Did not find a matching sibling element")
			return
		}
		found := false
		for _, attr := range res.Attr {
			if attr.Namespace == "" && attr.Key == key {
				found = true
				if attr.Val != val {
					t.Errorf("Expected \"%s\" for attribute \"%s\", got \"%s\"", val, key, attr.Val)
					return
				}
			}
		}
		if !found {
			t.Errorf("Did not find a matching sibling element with an attribute \"%s\" of value \"%s\"", key, val)
		}
	}

	expect(parent.FirstChild, "href", "https://example.net", atom.A)
	expect(parent.FirstChild, "id", "pre1", atom.Pre)
	expect(parent.FirstChild, "id", "pre1", atom.A, atom.Pre)
	expect(parent.FirstChild, "id", "pre1", atom.Pre, atom.A)
	expect(parent.FirstChild, "id", "StartTestNextSiblingByTag", atom.Pre, atom.A, atom.P)
}
