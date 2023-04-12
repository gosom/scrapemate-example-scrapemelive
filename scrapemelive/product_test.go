package scrapemelive

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func openTestFile(t *testing.T, filename string) *goquery.Document {
	t.Helper()
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	require.NoError(t, err)
	return doc
}

func Test_parseProduct(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-product.html")
	product := parseProduct(doc)
	require.Equal(t, "Charmeleon", product.Title)
	require.Equal(t, "£165.00", product.Price)
	require.Equal(t, "Charmeleon mercilessly destroys its foes using its sharp claws. If it encounters a strong foe, it turns aggressive. In this excited state, the flame at the tip of its tail flares with a bluish white color.", product.ShortDescription)
	require.Equal(t, "6565", product.Sku)
	require.ElementsMatch(t, []string{"Pokemon", "Flame"}, product.Categories)
	require.ElementsMatch(t, []string{"Blaze", "charmeleon", "Flame"}, product.Tags)
}

func Test_parseTitle(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-product.html")
	require.Equal(t, "Charmeleon", parseTitle(doc))
}

func Test_parsePrice(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-product.html")
	require.Equal(t, "£165.00", parsePrice(doc))
}

func Test_parseShortDescription(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-product.html")
	require.Equal(t, "Charmeleon mercilessly destroys its foes using its sharp claws. If it encounters a strong foe, it turns aggressive. In this excited state, the flame at the tip of its tail flares with a bluish white color.", parseShortDescription(doc))
}

func Test_parseSku(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-product.html")
	require.Equal(t, "6565", parseSku(doc))
}

func Test_parseCategories(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-product.html")
	require.ElementsMatch(t, []string{"Pokemon", "Flame"}, parseCategories(doc))
}

func Test_parseTags(t *testing.T) {
	t.Parallel()
	doc := openTestFile(t, "../testdata/sample-product.html")
	require.ElementsMatch(t, []string{"Blaze", "charmeleon", "Flame"}, parseTags(doc))
}
