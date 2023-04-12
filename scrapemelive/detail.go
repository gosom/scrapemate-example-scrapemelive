package scrapemelive

import (
	"context"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gosom/kit/logging"
	"github.com/gosom/scrapemate"
)

type ProductJob struct {
	scrapemate.Job
}

func (o *ProductJob) Process(ctx context.Context, resp scrapemate.Response) (any, []scrapemate.IJob, error) {
	log := ctx.Value("log").(logging.Logger)
	log.Info("processing product job")
	doc, ok := resp.Document.(*goquery.Document)
	if !ok {
		return nil, nil, errors.New("failed to convert response to goquery document")
	}
	product := parseProduct(doc)
	return product, nil, nil
}
