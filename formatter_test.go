package gosearch

import (
	"fmt"
	"testing"
)

func TestDefaultFormatter_Format(t *testing.T) {
	f := NewDefaultFormatter(RenderParts(RenderAll), NotRenderEmpty())

	for _, p := range []Package{
		{
			Name:       "demo",
			ImportPath: "github.com/demo",
			HomeSite:   "https://github.com/demo",
			Synopsis:   "A demo Package for Test",
			License:    "Apache 2.0",
		}, {
			Name:       "demo",
			ImportPath: "github.com/demo",
			HomeSite:   "https://github.com/demo",
			Synopsis:   "A demo Package for Test",
			License:    "Apache 2.0",
		},
	} {
		fmt.Println(f.Format(p))
	}
}
