package vdom

import (
	"fmt"
	"testing"
)

// TestParse tests the tree returned from the Parse function for various different
// inputs.
func TestParse(t *testing.T) {
	// We'll use table-driven testing here.
	testCases := []struct {
		// A human-readable name describing this test case
		name string
		// The src html to be parsed
		src []byte
		// The expected tree to be returned from the Parse function
		expectedTree *Tree
	}{
		{
			name: "Element root",
			src:  []byte("<div></div>"),
			expectedTree: &Tree{
				Roots: []Node{
					&Element{
						Name: "div",
					},
				},
			},
		},
		{
			name: "Text root",
			src:  []byte("Hello"),
			expectedTree: &Tree{
				Roots: []Node{
					&Text{
						Value: []byte("Hello"),
					},
				},
			},
		},
		{
			name: "Comment root",
			src:  []byte("<!--comment-->"),
			expectedTree: &Tree{
				Roots: []Node{
					&Comment{
						Value: []byte("comment"),
					},
				},
			},
		},
		{
			name: "ProcInst root",
			src:  []byte("<?target inst?>"),
			expectedTree: &Tree{
				Roots: []Node{
					&ProcInst{
						Target: "target",
						Inst:   []byte("inst"),
					},
				},
			},
		},
		{
			name: "Directive root",
			src:  []byte("<!doctype html>"),
			expectedTree: &Tree{
				Roots: []Node{
					&Directive{
						Value: []byte("doctype html"),
					},
				},
			},
		},
		{
			name: "ul with nested li's",
			src:  []byte("<ul><li>one</li><li>two</li><li>three</li></ul>"),
			expectedTree: &Tree{
				Roots: []Node{
					&Element{
						Name: "ul",
						children: []Node{
							&Element{
								Name: "li",
								children: []Node{
									&Text{
										Value: []byte("one"),
									},
								},
							},
							&Element{
								Name: "li",
								children: []Node{
									&Text{
										Value: []byte("two"),
									},
								},
							},
							&Element{
								Name: "li",
								children: []Node{
									&Text{
										Value: []byte("three"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Element with attrs",
			src:  []byte(`<div class="container" id="main" data-custom-attr="foo"></div>`),
			expectedTree: &Tree{
				Roots: []Node{
					&Element{
						Name: "div",
						Attrs: []Attr{
							{Name: "class", Value: "container"},
							{Name: "id", Value: "main"},
							{Name: "data-custom-attr", Value: "foo"},
						},
					},
				},
			},
		},
		{
			name: "Script tag with escaped characters",
			src:  []byte(`<script type="text/javascript">function((){console.log("&lt;Hello brackets&gt;")})()</script>`),
			expectedTree: &Tree{
				Roots: []Node{
					&Element{
						Name: "script",
						Attrs: []Attr{
							{Name: "type", Value: "text/javascript"},
						},
						children: []Node{
							&Text{
								Value: []byte(`function((){console.log("<Hello brackets>")})()`),
							},
						},
					},
				},
			},
		},
		{
			name: "Form with autoclosed tags",
			src:  []byte(`<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`),
			expectedTree: &Tree{
				Roots: []Node{
					&Element{
						Name: "form",
						Attrs: []Attr{
							{Name: "method", Value: "post"},
						},
						children: []Node{
							&Element{
								Name: "input",
								Attrs: []Attr{
									{Name: "type", Value: "text"},
									{Name: "name", Value: "firstName"},
								},
							},
							&Element{
								Name: "input",
								Attrs: []Attr{
									{Name: "type", Value: "text"},
									{Name: "name", Value: "lastName"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple roots",
			src:  []byte("<!doctype html><div></div>Hello<!--comment--><?target inst?>"),
			expectedTree: &Tree{
				Roots: []Node{
					&Directive{
						Value: []byte("doctype html"),
					},
					&Element{
						Name: "div",
					},
					&Text{
						Value: []byte("Hello"),
					},
					&Comment{
						Value: []byte("comment"),
					},
					&ProcInst{
						Target: "target",
						Inst:   []byte("inst"),
					},
				},
			},
		},
	}
	// Iterate through each test case
	for i, tc := range testCases {
		// Parse the input from tc.src
		gotTree, err := Parse(tc.src)
		if err != nil {
			t.Errorf("Unexpected error in Parse: %s", err.Error())
		}
		// Check that the resulting tree matches what we expect
		if match, msg := tc.expectedTree.Compare(gotTree); !match {
			t.Errorf("Error in test case %d (%s): HTML was not parsed correctly.\n%s", i, tc.name, msg)
		}
	}
}

// TestHTML tests the HTML method for each node in a parsed tree for various different
// inputs.
func TestHTML(t *testing.T) {
	// We'll use table-driven testing here.
	testCases := []struct {
		// A human-readable name describing this test case
		name string
		// The src html to be parsed
		src []byte
		// A function which should check the results of the HTML method of each
		// node in the parsed tree, and return an error if any results are incorrect.
		testFunc func(*Tree) error
	}{
		{
			name: "Element root",
			src:  []byte("<div></div>"),
			testFunc: func(tree *Tree) error {
				expectedHTML := []byte("<div></div>")
				return expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root element")
			},
		},
		{
			name: "Text root",
			src:  []byte("Hello"),
			testFunc: func(tree *Tree) error {
				expectedHTML := []byte("Hello")
				return expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root text node")
			},
		},
		{
			name: "Comment root",
			src:  []byte("<!--comment-->"),
			testFunc: func(tree *Tree) error {
				expectedHTML := []byte("<!--comment-->")
				return expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root comment node")
			},
		},
		{
			name: "ProcInst root",
			src:  []byte("<?target inst?>"),
			testFunc: func(tree *Tree) error {
				expectedHTML := []byte("<?target inst?>")
				return expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root proc inst")
			},
		},
		{
			name: "Directive root",
			src:  []byte("<!doctype html>"),
			testFunc: func(tree *Tree) error {
				expectedHTML := []byte("<!doctype html>")
				return expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root directive")
			},
		},
		{
			name: "ul with nested li's",
			src:  []byte("<ul><li>one</li><li>two</li><li>three</li></ul>"),
			testFunc: func(tree *Tree) error {
				{
					// Test the root of the tree, the ul element
					expectedHTML := []byte("<ul><li>one</li><li>two</li><li>three</li></ul>")
					if err := expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "the root ul element"); err != nil {
						return err
					}
				}
				lis := tree.Roots[0].Children()
				{
					// Test each li element
					expectedHTML := [][]byte{
						[]byte("<li>one</li>"),
						[]byte("<li>two</li>"),
						[]byte("<li>three</li>"),
					}
					for i, li := range lis {
						desc := fmt.Sprintf("li element %d", i)
						if err := expectHTMLEquals(expectedHTML[i], li.HTML(), desc); err != nil {
							return err
						}
					}
				}
				{
					// Test each text node inside the li elements
					expectedHTML := [][]byte{
						[]byte("one"),
						[]byte("two"),
						[]byte("three"),
					}
					for i, li := range lis {
						gotHTML := li.Children()[0].HTML()
						desc := fmt.Sprintf("the text inside li element %d", i)
						if err := expectHTMLEquals(expectedHTML[i], gotHTML, desc); err != nil {
							return err
						}
					}
				}
				return nil
			},
		},
		{
			name: "Element with attrs",
			src:  []byte(`<div class="container" id="main" data-custom-attr="foo"></div>`),
			testFunc: func(tree *Tree) error {
				expectedHTML := []byte(`<div class="container" id="main" data-custom-attr="foo"></div>`)
				return expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root element")
			},
		},
		{
			name: "Script tag with escaped characters",
			src:  []byte(`<script type="text/javascript">function((){console.log("&lt;Hello brackets&gt;")})()</script>`),
			testFunc: func(tree *Tree) error {
				{
					// Test the root element
					expectedHTML := []byte(`<script type="text/javascript">function((){console.log("<Hello brackets>")})()</script>`)
					if err := expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root script element"); err != nil {
						return err
					}
				}
				{
					// Test the text node inside the root element
					expectedHTML := []byte(`function((){console.log("<Hello brackets>")})()`)
					if err := expectHTMLEquals(expectedHTML, tree.Roots[0].Children()[0].HTML(), "text node inside script element"); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			name: "Form with autoclosed tags",
			src:  []byte(`<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`),
			testFunc: func(tree *Tree) error {
				{
					// Test the root element
					expectedHTML := []byte(`<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`)
					if err := expectHTMLEquals(expectedHTML, tree.Roots[0].HTML(), "root script element"); err != nil {
						return err
					}
				}
				{
					inputs := tree.Roots[0].Children()
					// Test each child input element
					expectedHTML := [][]byte{
						[]byte(`<input type="text" name="firstName">`),
						[]byte(`<input type="text" name="lastName">`),
					}
					for i, input := range inputs {
						desc := fmt.Sprintf("input element %d", i)
						if err := expectHTMLEquals(expectedHTML[i], input.HTML(), desc); err != nil {
							return err
						}
					}
				}
				return nil
			},
		},
		{
			name: "Multiple roots",
			src:  []byte("<!doctype html><div></div>Hello<!--comment--><?target inst?>"),
			testFunc: func(tree *Tree) error {
				expectedHTML := [][]byte{
					[]byte(`<!doctype html>`),
					[]byte(`<div></div>`),
					[]byte(`Hello`),
					[]byte(`<!--comment-->`),
					[]byte(`<?target inst?>`),
				}
				for i, root := range tree.Roots {
					desc := fmt.Sprintf("root node %d of type %T", i, root)
					if err := expectHTMLEquals(expectedHTML[i], root.HTML(), desc); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}
	// Iterate through each test case
	for i, tc := range testCases {
		// Parse the input from tc.src
		gotTree, err := Parse(tc.src)
		if err != nil {
			t.Errorf("Unexpected error in Parse: %s", err.Error())
		}
		// Use the testFunc to test for certain conditions
		if err := tc.testFunc(gotTree); err != nil {
			t.Errorf("Error in test case %d (%s):\n%s", i, tc.name, err.Error())
		}
	}
}

// expectHTMLEquals returns an error if expected does not equal got. description should be
// a human-readable description of the element that was tested.
func expectHTMLEquals(expected []byte, got []byte, description string) error {
	if string(expected) != string(got) {
		return fmt.Errorf("HTML for %s was not correct.\n\tExpected: %s\n\tBut got:  %s", description, string(expected), string(got))
	}
	return nil
}

// TestInnerHTML tests the InnerHTML method for each element in a parsed tree for various different
// inputs.
func TestInnerHTML(t *testing.T) {
	// We'll use table-driven testing here.
	testCases := []struct {
		// A human-readable name describing this test case
		name string
		// The src html to be parsed
		src []byte
		// A function which should check the results of the InnerHTML method of each
		// node in the parsed tree, and return an error if any results are incorrect.
		testFunc func(*Tree) error
	}{
		{
			name: "Element root",
			src:  []byte("<div></div>"),
			testFunc: func(tree *Tree) error {
				expectedInner := []byte("")
				el := tree.Roots[0].(*Element)
				return expectInnerHTMLEquals(expectedInner, el.InnerHTML(), "root element")
			},
		},
		{
			name: "ul with nested li's",
			src:  []byte("<ul><li>one</li><li>two</li><li>three</li></ul>"),
			testFunc: func(tree *Tree) error {
				{
					// Test the root of the tree, the ul element
					expectedInner := []byte("<li>one</li><li>two</li><li>three</li>")
					el := tree.Roots[0].(*Element)
					if err := expectInnerHTMLEquals(expectedInner, el.InnerHTML(), "the root ul element"); err != nil {
						return err
					}
				}
				lis := tree.Roots[0].Children()
				{
					// Test each li element
					expectedInner := [][]byte{
						[]byte("one"),
						[]byte("two"),
						[]byte("three"),
					}
					for i, li := range lis {
						el := li.(*Element)
						desc := fmt.Sprintf("li element %d", i)
						if err := expectInnerHTMLEquals(expectedInner[i], el.InnerHTML(), desc); err != nil {
							return err
						}
					}
				}
				return nil
			},
		},
		{
			name: "Inner element with attrs",
			src:  []byte(`<div><div class="container" id="main" data-custom-attr="foo"></div></div>`),
			testFunc: func(tree *Tree) error {
				expectedInner := []byte(`<div class="container" id="main" data-custom-attr="foo"></div>`)
				el := tree.Roots[0].(*Element)
				return expectInnerHTMLEquals(expectedInner, el.InnerHTML(), "root element")
			},
		},
		{
			name: "Script tag with escaped characters",
			src:  []byte(`<script type="text/javascript">function((){console.log("&lt;Hello brackets&gt;")})()</script>`),
			testFunc: func(tree *Tree) error {
				expectedInner := []byte(`function((){console.log("<Hello brackets>")})()`)
				el := tree.Roots[0].(*Element)
				if err := expectInnerHTMLEquals(expectedInner, el.InnerHTML(), "root script element"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Form with autoclosed tags",
			src:  []byte(`<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`),
			testFunc: func(tree *Tree) error {
				{
					// Test the root element
					expectedInner := []byte(`<input type="text" name="firstName"><input type="text" name="lastName">`)
					el := tree.Roots[0].(*Element)
					if err := expectInnerHTMLEquals(expectedInner, el.InnerHTML(), "root script element"); err != nil {
						return err
					}
				}
				{
					inputs := tree.Roots[0].Children()
					// Test each child input element
					for i, input := range inputs {
						el := input.(*Element)
						desc := fmt.Sprintf("input element %d", i)
						if err := expectInnerHTMLEquals([]byte{}, el.InnerHTML(), desc); err != nil {
							return err
						}
					}
				}
				return nil
			},
		},
	}
	// Iterate through each test case
	for i, tc := range testCases {
		// Parse the input from tc.src
		gotTree, err := Parse(tc.src)
		if err != nil {
			t.Errorf("Unexpected error in Parse: %s", err.Error())
		}
		// Use the testFunc to test for certain conditions
		if err := tc.testFunc(gotTree); err != nil {
			t.Errorf("Error in test case %d (%s):\n%s", i, tc.name, err.Error())
		}
	}
}

// expectInnerHTMLEquals returns an error if expected does not equal got. description should be
// a human-readable description of the element that was tested.
func expectInnerHTMLEquals(expected []byte, got []byte, description string) error {
	if string(expected) != string(got) {
		return fmt.Errorf("InnerHTML for %s was not correct.\n\tExpected: %s\n\tBut got:  %s", description, string(expected), string(got))
	}
	return nil
}

func TestSelector(t *testing.T) {
	// We'll use table-driven testing here.
	testCases := []struct {
		// A human-readable name describing this test case
		name string
		// The src html to be parsed
		src []byte
		// A function which should check the results of the HTML method of each
		// node in the parsed tree, and return an error if any results are incorrect.
		testFunc func(*Tree) error
	}{}
}
