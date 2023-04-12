package scrapemelive

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseProductLinks(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-category.html")
	links := parseProductLinks(doc)
	require.Len(t, links, 16)
}

func Test_parseNextPage(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-category.html")
	nextPage := parseNextPage(doc)
	require.Equal(t, "https://scrapeme.live/shop/page/2/", nextPage)
}
