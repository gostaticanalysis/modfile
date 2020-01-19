package modfile_test

import (
	"testing"

	"github.com/gostaticanalysis/modfile"
	xmodfile "golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	cases := []struct {
		pkg    string
		hasMod bool
	}{
		{"a", true},
		{"b", false},
	}

	testdata := analysistest.TestData()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.pkg, func(t *testing.T) {
			rs := analysistest.Run(t, testdata, modfile.Analyzer, tt.pkg)
			if len(rs) != 1 {
				t.Fatal("unexpected result:", rs)
			}

			switch {
			case tt.hasMod && rs[0].Result == nil:
				t.Error("modfile cannot parse")
			case !tt.hasMod && rs[0].Result.(*xmodfile.File) != nil: // memo: typed nil
				t.Errorf("an unexpected modfile has parsed: %#v", rs[0].Result.(*xmodfile.File).Module)
			}
		})
	}
}
