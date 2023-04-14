package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gosom/scrapemate"
	"github.com/gosom/scrapemate-example-scrapemelive/scrapemelive"
	"github.com/gosom/scrapemate/adapters/cache/leveldbcache"
	fetcher "github.com/gosom/scrapemate/adapters/fetchers/nethttp"
	parser "github.com/gosom/scrapemate/adapters/parsers/goqueryparser"
	provider "github.com/gosom/scrapemate/adapters/providers/memory"
)

func main() {
	err := run()
	if err == nil || errors.Is(err, scrapemate.ErrorExitSignal) {
		os.Exit(0)
		return
	}
	os.Exit(1)
}

func run() error {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("deferred cancel"))

	provider := provider.New()

	go func() {
		job := &scrapemelive.ProductCollectJob{
			Job: scrapemate.Job{
				ID:     uuid.New().String(),
				Method: http.MethodGet,
				URL:    "https://scrapeme.live/shop/",
				Headers: map[string]string{
					"User-Agent": scrapemate.DefaultUserAgent,
				},
				Timeout:    10 * time.Second,
				MaxRetries: 3,
			},
		}
		provider.Push(ctx, job)
	}()

	httpFetcher := fetcher.New(&http.Client{
		Timeout: 10 * time.Second,
	})

	cacher, err := leveldbcache.NewLevelDBCache("__leveldb_cache")
	if err != nil {
		return err
	}

	mate, err := scrapemate.New(
		scrapemate.WithContext(ctx, cancel),
		scrapemate.WithJobProvider(provider),
		scrapemate.WithHttpFetcher(httpFetcher),
		scrapemate.WithConcurrency(10),
		scrapemate.WithHtmlParser(parser.New()),
		scrapemate.WithCache(cacher),
	)

	if err != nil {
		return err
	}

	resultsDone := make(chan struct{})
	go func() {
		defer close(resultsDone)
		if err := writeCsv(mate.Results()); err != nil {
			cancel(err)
			return
		}
	}()

	err = mate.Start()
	<-resultsDone
	return err
}

func writeCsv(results <-chan scrapemate.Result) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	headersWritten := false
	for result := range results {
		if result.Data == nil {
			continue
		}
		product, ok := result.Data.(scrapemelive.Product)
		if !ok {
			return fmt.Errorf("unexpected data type: %T", result.Data)
		}
		if !headersWritten {
			if err := w.Write(product.CsvHeaders()); err != nil {
				return err
			}
			headersWritten = true
		}
		if err := w.Write(product.CsvRow()); err != nil {
			return err
		}
		w.Flush()
	}
	return w.Error()
}
