package modfile

import (
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:       "modfile",
	Doc:        Doc,
	Run:        run,
	ResultType: reflect.TypeOf((*modfile.File)(nil)),
}

const Doc = "modfile is ..."

func run(pass *analysis.Pass) (result interface{}, rerr error) {
	result = (*modfile.File)(nil)

	cmd := exec.Command("go", "list", "-m", "-f", "{{.GoMod}}", pass.Pkg.Path())
	cmd.Env = append([]string{"GO111MODULE", "auto"}, os.Environ()...)
	pass.Fset.Iterate(func(f *token.File) bool {
		fname := f.Name()
		if filepath.Ext(fname) == ".go" &&
			!strings.HasSuffix(fname, "_test.go") {
			cmd.Dir = filepath.Dir(fname)
			return false
		}
		return true
	})

	output, err := cmd.Output()
	if err != nil {
		// ignore err
		return
	}
	modfilename := strings.TrimSpace(string(output))

	data, err := ioutil.ReadFile(modfilename)
	if err != nil {
		// ignore err
		return
	}

	f, err := modfile.Parse(modfilename, data, nil)
	if err != nil {
		// ignore err
		return
	}

	return f, nil
}
