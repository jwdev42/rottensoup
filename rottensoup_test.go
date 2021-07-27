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
		t.Error(err)
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
		if e.FirstChild.Type != html.TextNode {
			t.Error("Expected text node")
		}
		expect := fmt.Sprintf("Sibling %d", i+1)
		if e.FirstChild.Data != expect {
			t.Errorf("Expected \"%s\", got \"%s\"", expect, e.FirstChild.Data)
		}
	}
}

func TestElementsByAttrMatch(t *testing.T) {
	const matches = 4
	const testDoc = "attr_match.html"

	search := regexp.MustCompile("caption-[a-z]+")

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Error(err)
	}

	res := ElementsByAttrMatch(root, "class", search)
	if len(res) != matches {
		t.Errorf("Expected %d matches, got %d", matches, len(res))
	}
	for i, e := range res {
		if e.FirstChild.Type != html.TextNode {
			t.Error("Expected text node")
		}
		expect := fmt.Sprintf("Match %d", i+1)
		if e.FirstChild.Data != expect {
			t.Errorf("Expected \"%s\", got \"%s\"", expect, e.FirstChild.Data)
		}
	}
}
