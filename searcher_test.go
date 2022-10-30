package gosearch

import (
	"context"
	"net/http"
	"testing"
)

func TestOfficialSearcher_Search(t *testing.T) {
	searcher := NewOfficialSearcher(http.DefaultClient)
	packages, err := searcher.Search(context.Background(), "cobra")
	if err != nil {
		t.Fatal(err)
	}
	f := NewDefaultFormatter(RenderParts(RenderAll), NotRenderEmpty())
	for _, p := range packages {
		t.Log(f.Format(p))
		t.Log(f.ShortFormatter.Format(p))
	}
}
