package gosearch

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

const (
	officialSiteUrl = "https://pkg.go.dev"
)

// Searcher 查询器，给定关键字查询相关go mod
type Searcher interface {
	Search(ctx context.Context, keyword string) ([]Package, error)
}

// 官方包管理器查询器，通过模拟查询https://pkg.go.dev实现
type officialSearcher struct {
	cli  *http.Client
	conv *converter
}

// Search 查询官方页面，通过解析HTML来提取包信息
func (searcher *officialSearcher) Search(ctx context.Context, keyword string) ([]Package, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, officialSiteUrl+"/search?q="+keyword, nil)
	if err != nil {
		return nil, err
	}

	resp, err := searcher.cli.Do(request)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var pkgs = make([]Package, 0, 24)
	doc.Find(".SearchResults .SearchSnippet").Each(func(i int, s *goquery.Selection) {
		pkgs = append(pkgs, searcher.conv.Convert(s))
	})
	return pkgs, nil
}

// 官方包查询需要解析html页面，每个包分别在 .SearchResults .SearchSnippet 下定义
// 可以通过浏览器查看页面源代码获得
type converter struct {
	headerReplacer *strings.Replacer
	spaceRegexp    *regexp.Regexp
}

func (conv *converter) Convert(snap *goquery.Selection) Package {
	var pkg Package

	headerContainer := snap.Find(".SearchSnippet-headerContainer")
	header := conv.headerReplacer.Replace(headerContainer.Text())
	headerSpaceIdx := strings.IndexRune(header, ' ')
	pkg.Name = header[0:headerSpaceIdx]
	pkg.ImportPath = header[headerSpaceIdx+1:]
	href, exists := headerContainer.Find("a").Attr("href")
	if exists {
		pkg.HomeSite = officialSiteUrl + href
	}

	synopsis := snap.Find(".SearchSnippet-synopsis").Text()
	pkg.Synopsis = strings.TrimSpace(synopsis)

	infoDom := snap.Find(".SearchSnippet-infoLabel")
	pkg.License = strings.TrimSpace(infoDom.Find(".snippet-license").Text())

	type applier func(v string)

	appliers := []applier{
		func(v string) {
			pkg.ImportedBy = strings.TrimSpace(v)
		},
		func(v string) {
			pkg.Version = strings.TrimSpace(v)
		},
		func(v string) {
			pkg.Published = strings.TrimSpace(v)
		},
	}
	applyNodes := infoDom.Find("strong").Nodes
	for i := 0; i < len(applyNodes) && i < len(appliers); i++  {
		appliers[i](applyNodes[i].FirstChild.Data)
	}
	return pkg
}

func newConverter() *converter {
	return &converter{
		headerReplacer: strings.NewReplacer("\n", "", " ", "", "(", " ", ")", ""),
		spaceRegexp:    regexp.MustCompile("[\n\\s]+"),
	}
}

func NewOfficialSearcher(cli *http.Client) Searcher {
	return &officialSearcher{
		cli:  cli,
		conv: newConverter(),
	}
}
