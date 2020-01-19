package modfile_test

import (
	"testing"

	"github.com/gostaticanalysis/modfile"
	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
	xmodfile "golang.org/x/mod/modfile"
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
			ignoreErrT := testutil.Filter(t, func(format string, args ...interface{}) bool {
				// ignore errors
				return false
			})
			rs := analysistest.Run(ignoreErrT, testdata, modfile.Analyzer, tt.pkg)
			if len(rs) != 1 {
				t.Fatal("unexpected result:", rs)
			}

			switch {
			case tt.hasMod && rs[0].Result == nil:
				t.Error("modfile cannot parse")
			case !tt.hasMod && rs[0].Result != nil:
				t.Errorf("an unexpected modfile has parsed: %#v", rs[0].Result.(*xmodfile.File).Module)
			}
		})
	}
}
