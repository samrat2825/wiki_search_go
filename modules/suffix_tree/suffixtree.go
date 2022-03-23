package suffixtree

import (
	"strings"
	"unicode/utf8"
)

type generalizedSuffixTree struct {
	root       *node 
	activeLeaf *node 
}

func (t *generalizedSuffixTree) Search(word string, numElements int) []int {
	node := t.searchNode(word)
	if node == nil {
		return nil
	}
	return node.getData(numElements)
}

func (t *generalizedSuffixTree) searchNode(word string) *node {
	var currentNode *node = t.root
	var currentEdge *edge
	var i int

	for i < len(word) {
		rune, _ := utf8.DecodeRuneInString(word[i:])
		currentEdge = currentNode.getEdge(rune)
		if currentEdge == nil {
			return nil
		} else {
			label := string(currentEdge.label)
			lenToMatch := len(word) - i
			if lenToMatch > len(label) {
				lenToMatch = len(label)
			}
			if word[i:i+lenToMatch] != label[:lenToMatch] {
				return nil
			}

			if len(label) >= len(word)-i {
				return currentEdge.node
			} else {
				currentNode = currentEdge.node
				i += lenToMatch
			}
		}
	}

	return nil
}

func (t *generalizedSuffixTree) Put(key string, index int) {
	t.activeLeaf = t.root
	s := t.root
	runes := []rune(key)

	var text []rune
	for k, r := range runes {
		text = append(text, r)
		s, text = t.update(s, text, runes[k:], index)
		s, text = t.canonize(s, text)
	}

	if t.activeLeaf.suffix == nil && t.activeLeaf != t.root && t.activeLeaf != s {
		t.activeLeaf.suffix = s
	}
}

func (t *generalizedSuffixTree) update(inputNode *node, stringPart []rune, rest []rune, value int) (s *node, runes []rune) {
	s = inputNode
	runes = stringPart
	newRune := stringPart[len(stringPart)-1]

	oldroot := t.root
	endpoint, r := t.testAndSplit(s, stringPart[:len(stringPart)-1], newRune, rest, value)

	var leaf *node
	for !endpoint {
		tempEdge := r.getEdge(newRune)
		if tempEdge != nil {
			leaf = tempEdge.node
		} else {
			leaf = newNode()
			leaf.addRef(value)
			newedge := newEdge(rest, leaf)
			r.addEdge(newRune, newedge)
		}

		if t.activeLeaf != t.root {
			t.activeLeaf.suffix = leaf
		}
		t.activeLeaf = leaf

		if oldroot != t.root {
			oldroot.suffix = r
		}

		oldroot = r
		if s.suffix == nil { 
			runes = runes[1:]
		} else {
			n, b := t.canonize(s.suffix, safeCutLastChar(runes))
			s = n
			runes = append(b, runes[len(runes)-1])
		}

		endpoint, r = t.testAndSplit(s, safeCutLastChar(runes), newRune, rest, value)
	}

	if oldroot != t.root {
		oldroot.suffix = r
	}

	return
}

func (t *generalizedSuffixTree) canonize(s *node, runes []rune) (*node, []rune) {

	currentNode := s
	if len(runes) > 0 {
		g := s.getEdge(runes[0])
		for g != nil && strings.Index(string(runes), string(g.label)) == 0 {
			runes = runes[len(g.label):]
			currentNode = g.node
			if len(runes) > 0 {
				g = currentNode.getEdge(runes[0])
			}
		}
	}
	return currentNode, runes
}

func (t *generalizedSuffixTree) testAndSplit(inputs *node, stringPart []rune, r rune, remainder []rune, value int) (bool, *node) {
	s, str := t.canonize(inputs, stringPart)

	if len(str) > 0 {
		g := s.getEdge(str[0])

		if len(g.label) > len(str) && g.label[len(str)] == r {
			return true, s
		} else {
			newlabel := g.label[len(str):]

			w := newNode()
			newedge := newEdge(str, w)
			s.addEdge(str[0], newedge)
			g.label = newlabel
			w.addEdge(newlabel[0], g)

			return false, w
		}
	} else {
		e := s.getEdge(r)
		if e == nil {
			return false, s
		} else {
			if string(remainder) == string(e.label) {
				e.node.addRef(value)
				return true, s
			} else if strings.Index(string(remainder), string(e.label)) == 0 {
				return true, s
			} else if strings.Index(string(e.label), string(remainder)) == 0 {
				newNode := newNode()
				newNode.addRef(value)
				newEdge := newEdge(remainder, newNode)
				s.addEdge(r, newEdge)

				e.label = e.label[len(remainder):]
				newNode.addEdge(e.label[0], e)
				return false, s
			} else {
				return true, s
			}
		}
	}

}

func safeCutLastChar(runes []rune) []rune {
	if len(runes) == 0 {
		return nil
	}
	return runes[:len(runes)-1]
}

func NewGeneralizedSuffixTree() *generalizedSuffixTree {
	t := &generalizedSuffixTree{}
	t.root = newNode()
	t.activeLeaf = t.root
	return t
}