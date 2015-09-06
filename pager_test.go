package pager

import "testing"
import "strings"

type pagerInput struct {
	currentPage    int
	totalElements  int
	perPage        int
	maxLinks       int
	linksLeftRight int
	linksOnEnd     int
	Expected       string
}

func TestLinks(t *testing.T) {
	testcases := []pagerInput{
		pagerInput{0, 0, 10, 10, 2, 2, "_^_,_<_,[0],_>_,_$_"},
		pagerInput{0, 3, 10, 10, 2, 2, "_^_,_<_,[0],_>_,_$_"},
		pagerInput{0, 3, 2, 10, 2, 2, "_^_,_<_,[0],(1),(>),($)"},
		pagerInput{0, 3, 1, 10, 2, 2, "_^_,_<_,[0],(1),(2),(>),($)"},
	}

	for _, test := range testcases {
		pager := NewPager(test.currentPage, test.totalElements, test.perPage, test.maxLinks, test.linksLeftRight, test.linksOnEnd)
		assertLinks(t, pager.Links(), test.Expected)
	}
}

func assertLinks(t *testing.T, links []Link, expected string) {
	serialized := serializeLinks(links)

	if serialized != expected {
		t.Errorf("Link sequence does not meet expectations:\nExpected: %s\nActual  : %s", expected, serialized)
	}
}

func serializeLinks(links []Link) string {
	results := make([]string, 0)

	for _, link := range links {
		results = append(results, link.String())
	}

	return strings.Join(results, ",")
}
