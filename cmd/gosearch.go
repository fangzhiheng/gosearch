package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/fangzhiheng/gosearch"
	"github.com/spf13/cobra"
	"math"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	var option SearchOption
	root := &cobra.Command{
		Use:     os.Args[0],
		Example: os.Args[0] + " gin cobra",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			option.Keywords = args
			return Search(cmd.Context(), option)
		},
	}

	flags := root.Flags()
	flags.BoolVar(&option.Insecure, "insecure", false, "Skip https cert verify")
	flags.BoolVarP(&option.RenderShortly, "short", "s", false, "Render results shortly")
	flags.BoolVarP(&option.NotRenderTitle, "omittitle", "r", false, "Render raw results")

	err := root.ExecuteContext(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "err: %+v", err)
		stop()
	}
}

type SearchOption struct {
	Keywords       []string
	NotRenderTitle bool
	Insecure       bool
	RenderShortly  bool
}

func Search(ctx context.Context, option SearchOption) error {
	keywords := parseKeywords(option)
	if len(keywords) == 0 {
		return nil
	}
	client := createHttpClient(option)
	searcher := gosearch.NewOfficialSearcher(client)
	parallelism := int(math.Min(float64(runtime.NumCPU()), float64(len(keywords))))
	keywordsCh := make(chan string, parallelism)
	packagesCh := make(chan []gosearch.Package, parallelism)

	searcherWaitGroup := sync.WaitGroup{}
	for i := 0; i < parallelism; i++ {
		searcherWaitGroup.Add(1)
		go func() {
			defer searcherWaitGroup.Done()
			for {
				select {
				case _, _ = <-ctx.Done():
					return
				case keyword, ok := <-keywordsCh:
					if !ok {
						return
					}
					packages, err := searcher.Search(ctx, keyword)
					if err != nil {
						continue
					}
					packagesCh <- packages
				}
			}
		}()
	}

	go func() {
		searcherWaitGroup.Wait()
		close(packagesCh)
	}()

	for _, keyword := range keywords {
		keywordsCh <- keyword
	}
	close(keywordsCh)

	var formatter = createFormatter(option)
	var handledPackage = map[string]struct{}{}
	for packages := range packagesCh {
		for _, pkg := range packages {
			if _, ok := handledPackage[pkg.ImportPath]; ok {
				continue
			}
			handledPackage[pkg.ImportPath] = struct{}{}
			_, _ = fmt.Fprintln(os.Stdout, formatter.Format(pkg))
		}
	}
	return nil
}

func parseKeywords(option SearchOption) []string {
	var keywords []string
	for _, fk := range option.Keywords {
		ks := strings.Split(fk, ",")
		keywords = append(keywords, ks...)
	}
	return keywords
}

func createFormatter(option SearchOption) gosearch.Formatter {
	var formatter gosearch.Formatter
	var renderOptions = []gosearch.FormatterOption{
		gosearch.RenderParts(gosearch.RenderAll),
		gosearch.NotRenderEmpty(),
	}
	if option.NotRenderTitle {
		renderOptions = append(renderOptions, gosearch.NotRenderTitle())
	}
	if option.RenderShortly {
		formatter = gosearch.NewShortFormatter(renderOptions...)
	} else {
		formatter = gosearch.NewDefaultFormatter(renderOptions...)
	}
	return formatter
}

func createHttpClient(option SearchOption) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: option.Insecure,
			},
		},
	}
}
