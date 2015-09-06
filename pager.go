package pager

import "math"
import "strconv"

type LinkType int

const (
	Regular LinkType = iota
	LinkFirst
	LinkPrev
	LinkNext
	LinkLast
	LinkEllipsis
)

type Link struct {
	Page    int
	Type    LinkType
	Enabled bool
	Active  bool
}

func (l *Link) String() string {
	chr := ""

	switch l.Type {
	case Regular:
		chr = strconv.Itoa(l.Page)
	case LinkFirst:
		chr = "^"
	case LinkPrev:
		chr = "<"
	case LinkNext:
		chr = ">"
	case LinkLast:
		chr = "$"
	case LinkEllipsis:
		chr = "…"
	}

	if l.Enabled {
		if l.Active {
			return "[" + chr + "]"
		} else {
			return "(" + chr + ")"
		}
	} else {
		return "_" + chr + "_"
	}
}

type Pager struct {
	currentPage    int
	totalElements  int
	perPage        int
	maxLinks       int
	linksLeftRight int
	linksOnEnd     int
	pages          int
}

func NewPager(currentPage int, totalElements int, perPage int, maxLinks int, linksLeftRight int, linksOnEnd int) Pager {
	pages := int(math.Ceil(float64(totalElements) / float64(perPage)))

	if currentPage > (pages - 1) {
		currentPage = pages - 1
	}

	if currentPage < 0 {
		currentPage = 0
	}

	if maxLinks < 5 {
		maxLinks = 5
	}

	return Pager{currentPage, totalElements, perPage, maxLinks, linksLeftRight, linksOnEnd, pages}
}

func NewBasicPager(currentPage int, totalElements int) Pager {
	return NewPager(currentPage, totalElements, 10, 10, 2, 2)
}

func (p *Pager) Links() []Link {
	result := make([]Link, 0, 4) // we need at least the first/prev/next/last links

	// create the static first/prev links
	goBack := p.currentPage > 0
	result = append(result, Link{0, LinkFirst, goBack, false})
	result = append(result, Link{p.currentPage - 1, LinkPrev, goBack, false})

	// To make things simple, we now collect page numbers to link to instead
	// of full structs; this makes removing duplicates easier; we wrap the pages
	// into structs later on.
	pages := make([]int, 0)

	// easy, we have fewer links that the maximum
	if p.pages <= p.maxLinks {
		// one link per page
		for i := 0; i < p.pages; i++ {
			pages = append(pages, i)
		}

		// create at least one link to the current page
		if p.pages == 0 {
			pages = append(pages, 0)
		}
	} else {
		// always place a link to the first page
		pages = append(pages, 0)

		// create a number of links starting from 1 (now we have [0, 1, 2, 3, 4, ])
		for i := 0; i < p.linksOnEnd; i++ {
			pages = append(pages, i+1)
		}

		// create links surrounding the current page
		begin := p.currentPage - p.linksLeftRight
		end := p.currentPage + p.linksLeftRight

		// if there is a gap between the already built start of the list and the
		// $begin of our list leading up to the current page, place an ellipsis
		// (-1 will be translated later on to the LinkEllipsis constant)
		if (begin - 1) > p.linksOnEnd {
			pages = append(pages, -1)
		}

		// create the [(cur-3),(cur-2),(cur-1),(cur),(cur+1),... links]
		for i := begin; i <= end; i++ {
			if i > 0 && i < p.pages {
				pages = append(pages, i)
			}
		}

		// place an ellipsis if there's a gap
		if end < (p.pages - p.linksOnEnd - 2) {
			pages = append(pages, -1)
		}

		// create links leading to the end ([..., (max-3),(max-2),(max-1)])
		for i := p.linksOnEnd; i > 0; i-- {
			pages = append(pages, p.pages-i-1)
		}

		// create link to the last page
		pages = append(pages, p.pages-1)

		// remove duplicates
		pages = removeDuplicates(pages)
	}

	// wrap the page numbers into structs
	for _, page := range pages {
		ltype := Regular
		enabled := true

		if page == -1 {
			ltype = LinkEllipsis
			enabled = false
		}

		result = append(result, Link{page, ltype, enabled, page == p.currentPage})
	}

	// create the static next/last links
	goNext := p.currentPage < (p.pages - 1)
	result = append(result, Link{p.currentPage + 1, LinkNext, goNext, false})
	result = append(result, Link{p.pages - 1, LinkLast, goNext, false})

	return result
}

func removeDuplicates(input []int) []int {
	notebook := map[int]bool{}
	result := []int{}

	for _, value := range input {
		if _, seen := notebook[value]; !seen {
			result = append(result, value)
			notebook[value] = true
		}
	}

	return result
}
