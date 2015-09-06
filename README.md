Go Pager
========

This package implements a very basic pager, useful for creating paginations on lists in various
applications. The pager automatically inserts ``…`` in the sequence of links to prevent excessive
amounts of links being generated.

There is no output logic implemented here. You put in a few basic numbers and get back a list of
abstract ``Link`` structures. Rendering those is your task.

Examples
--------

```go
import "github.com/xrstf/pager"

currentPage   := 0 // pages start at zero
totalElements := 1234

// a basic pager sets a few sane default values;
// use pager.NewPager() to control all of them
myPager := pager.NewBasicPager(currentPage, totalElements)

// links is a []pager.Link slice
links := myPager.Links()

// render as you like
for _, link := links {
	if link.Type == pager.LinkFirst {
		writeSomewhere("<a href='/list'>&lt;&lt;</a>")
	} else if link.Type == pager.LinkPrev {
		writeSomewhere("<a href='/list?page=" + strconv.Itoa(link.Page) + "'>&lt;</a>")
	}
	// etc.
}
```

License
-------

The code is licensed under the MIT license.
