package scrapemelive

import (
	"context"
	"errors"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/gosom/kit/logging"
	"github.com/gosom/scrapemate"
)

type ProductCollectJob struct {
	scrapemate.Job
}

func (o *ProductCollectJob) Process(ctx context.Context, resp scrapemate.Response) (any, []scrapemate.IJob, error) {
	log := ctx.Value("log").(logging.Logger)
	log.Info("processing collect job")
	doc, ok := resp.Document.(*goquery.Document)
	if !ok {
		return nil, nil, errors.New("failed to convert response to goquery document")
	}
	var nextJobs []scrapemate.IJob
	links := parseProductLinks(doc)
	for _, link := range links {
		nextJobs = append(nextJobs, &ProductJob{
			Job: scrapemate.Job{
				ID:     uuid.New().String(),
				Method: "GET",
				URL:    link,
				Headers: map[string]string{
					"User-Agent": scrapemate.DefaultUserAgent,
				},
				Timeout:    10 * time.Second,
				MaxRetries: 3,
				Priority:   0,
			},
		})
	}
	nextPage := parseNextPage(doc)
	if nextPage != "" {
		nextJobs = append(nextJobs, &ProductCollectJob{
			Job: scrapemate.Job{
				ID:     uuid.New().String(),
				Method: "GET",
				URL:    nextPage,
				Headers: map[string]string{
					"User-Agent": scrapemate.DefaultUserAgent,
				},
				Timeout:    10 * time.Second,
				MaxRetries: 3,
				Priority:   1,
			},
		})
	}

	return nil, nextJobs, nil
}

func parseProductLinks(doc *goquery.Document) []string {
	var links []string
	doc.Find("a.woocommerce-LoopProduct-link").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		links = append(links, link)
	})
	return links
}

func parseNextPage(doc *goquery.Document) string {
	return doc.Find("a.next.page-numbers").AttrOr("href", "")
}
