//This file is part of rottensoup ©2021 Jörg Walter

package rottensoup

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func TestAttrVal(t *testing.T) {
	e := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "id",
				Val: "test",
			},
		},
	}
	if AttrVal(e, "", "id") != "test" {
		t.Errorf("Expected value \"test\" for attribute \"id\", got \"%s\" instead", AttrVal(e, "", "id"))
	}
	if AttrVal(e, "", "class") != "" {
		t.Error("AttrVal returned nonempty string for unavailable attribute")
	}
}

func TestHasAttr(t *testing.T) {
	e := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "id",
				Val: "test",
			},
		},
	}

	if !HasAttr(e, "", "id") {
		t.Error("Element has an attr \"id\", but HasAttr returns false.")
	}
	if HasAttr(e, "", "class") {
		t.Error("Element hasn't an attr \"class\", but HasAttr returns true.")
	}
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
		t.Fatalf("No Element found with id \"%s\"!", testID)
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

	solitude := new(html.Node)
	if NextElementSibling(solitude) != nil {
		t.Error("Unexpected non-nil return value")
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
	if AttrVal(prev, "", "href") != expectedHref {
		t.Fatalf("Expected \"%s\", got \"%s\"", expectedHref, AttrVal(prev, "", "href"))
	}

	solitude := new(html.Node)
	if PrevElementSibling(solitude) != nil {
		t.Error("Unexpected non-nil return value")
	}
}

func TestElementByID(t *testing.T) {
	const testDoc = "test.html"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}

	exists := ElementByID(root, "siblings")
	if exists == nil {
		t.Error("ElementByID didn't find existing element")
	}
	doesnotexist := ElementByID(root, "imnothere")
	if doesnotexist != nil {
		t.Error("ElementByID returned an element that shouldn't exist")
	}
}

func TestElementsByAttr(t *testing.T) {
	const testDoc = "by_tag_and_attr.html"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}

	attrs := []html.Attribute{
		html.Attribute{Key: "class", Val: "cell"},
		html.Attribute{Key: "lang", Val: "de"},
		html.Attribute{Key: "title", Val: "test"},
	}
	invalidAttr := html.Attribute{Key: "foo", Val: "bar"}

	c1 := ElementsByAttr(root, attrs[0])
	c2 := ElementsByAttr(root, attrs[0], attrs[1])
	c3 := ElementsByAttr(root, attrs...)

	if i := len(c1); i != 19 {
		t.Errorf("c1: Expected length of 19, got %d", i)
	}
	if i := len(c2); i != 10 {
		t.Errorf("c1: Expected length of 10, got %d", i)
	}
	if i := len(c3); i != 5 {
		t.Errorf("c1: Expected length of 5, got %d", i)
	}

	if c1[0] != FirstElementByAttr(root, attrs[0]) {
		t.Errorf("FirstElementByAttr returns wrong element")
	}

	if doesnotexist := ElementsByAttr(root, invalidAttr); doesnotexist != nil {
		t.Error("ElementsByAttr returned an element that should not exist")
	}

	if doesnotexist := FirstElementByAttr(root, invalidAttr); doesnotexist != nil {
		t.Error("FirstElementByAttr returned an element that should not exist")
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

	res := ElementsByAttrMatch(root, "", "class", search)
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

	if invalid := ElementsByAttrMatch(root, "", "class", regexp.MustCompile("imnothere")); invalid != nil {
		t.Error("Expected no match")
	}

}

func TestElementsByClassName(t *testing.T) {
	const testDoc = "classnames.html"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}

	numInChildText := func(n *html.Node, num int) bool {
		textNode := FirstNodeByType(n, html.TextNode)
		if textNode == nil {
			panic("No text node found")
		}
		return strings.Contains(textNode.Data, fmt.Sprintf("%d", num))
	}

	testClassValues := func(prefix string, expectedElements int, nums ...int) {
		classes := make([]string, len(nums))
		for i, num := range nums {
			classes[i] = fmt.Sprintf("%s%d", prefix, num)
		}
		elems := ElementsByClassName(root, classes...)
		first := FirstElementByClassName(root, classes...)
		if elems != nil && first != elems[0] || elems == nil && first != nil {
			t.Error("FirstElementByClassName returned the wrong element")
		}
		if l := len(elems); l != expectedElements {
			t.Errorf("ElementsByClassName should have returned %d elements, it returned %d elements instead", expectedElements, l)
		}
		for _, e := range elems {
			for _, i := range nums {
				if !numInChildText(e, i) {
					t.Errorf("ElementsByClassName returned an element that should be a member of class \"%s\", but isn't", classes[i])
				}
			}
		}
	}
	testClassValues("class", 0, 0)
	testClassValues("class", 4, 1)
	testClassValues("class", 3, 2)
	testClassValues("class", 4, 3)
	testClassValues("class", 1, 4)
	testClassValues("class", 2, 1, 2, 3)
	testClassValues("class", 1, 1, 2, 3, 4)
}

func TestElementsByTag(t *testing.T) {
	const testDoc = "by_tag.html"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}

	test := func(expect int, tag ...atom.Atom) {
		tags := ElementsByTag(root, tag...)
		if tags != nil {
			if len(tags) != expect {
				t.Errorf("Expected %d elements, got %d", expect, len(tags))
			}
			if FirstElementByTag(root, tag...) != tags[0] {
				t.Error("FirstElementByTag has wrong element")
			}
		} else {
			if expect != 0 {
				t.Error("Expected nil but got a non-nil result")
			}
			if FirstElementByTag(root, tag...) != nil {
				t.Error("FirstElementByTag was expected to be nil")
			}
		}
	}

	test(0, atom.Autocomplete)
	test(4, atom.P)
	test(8, atom.P, atom.Div)

}

func TestElementsByTagAndAttr(t *testing.T) {
	const testDoc = "by_tag_and_attr.html"

	root, err := parseTestFile(testDoc)
	if err != nil {
		t.Fatal(err)
	}

	attrs := []html.Attribute{
		html.Attribute{Key: "class", Val: "cell"},
		html.Attribute{Key: "lang", Val: "de"},
		html.Attribute{Key: "title", Val: "test"},
	}

	c1 := ElementsByTagAndAttr(root, atom.Td, attrs[0])
	c2 := ElementsByTagAndAttr(root, atom.Td, attrs[0], attrs[1])
	c3 := ElementsByTagAndAttr(root, atom.Td, attrs...)

	verifyCells := func(nodes []*html.Node, offset, length int) {
		if len(nodes) != length {
			t.Errorf("Invalid result length: Expected %d, got %d", length, len(nodes))
			return
		}
		for i, v := range nodes {
			cellID := i + 1 + offset
			if v.DataAtom != atom.Td {
				t.Errorf("Tag mismatch in table cell %d: Expected \"%s\", got \"%s\"", cellID, atom.Td.String(), v.DataAtom.String())
				continue
			}
			text := FirstNodeByType(v, html.TextNode)
			if text == nil {
				t.Errorf("No text node found in table cell %d", cellID)
				continue
			}
			if text.Data != fmt.Sprintf("%d", cellID) {
				t.Errorf("Text mismatch in table cell %d: Expected \"%d\", got \"%s\"", cellID, cellID, text.Data)
			}
		}
	}

	verifyCells(c1, 0, 16) //c1 must match cells 1 to 16
	verifyCells(c2, 4, 8)  //c2 must match cells 5 to 12
	verifyCells(c3, 8, 4)  //c3 must match cells 9 to 12

	invalid := ElementsByTagAndAttr(root, atom.Td, html.Attribute{Key: "foo", Val: "bar"})
	if invalid != nil {
		t.Error("Expected return value nil")
	}
	if FirstElementByTagAndAttr(root, atom.Td, html.Attribute{Key: "foo", Val: "bar"}) != nil {
		t.Error("Expected return value nil")
	}
	if c1[0] != FirstElementByTagAndAttr(root, atom.Td, attrs[0]) {
		t.Errorf("FirstElementByTagAndAttr returns wrong element")
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

	if nothing := NextSiblingByTag(parent.FirstChild, atom.Table); nothing != nil {
		t.Error("Expected return value nil")
	}

	expect(parent.FirstChild, "href", "https://example.net", atom.A)
	expect(parent.FirstChild, "id", "pre1", atom.Pre)
	expect(parent.FirstChild, "id", "pre1", atom.A, atom.Pre)
	expect(parent.FirstChild, "id", "pre1", atom.Pre, atom.A)
	expect(parent.FirstChild, "id", "StartTestNextSiblingByTag", atom.Pre, atom.A, atom.P)
}
