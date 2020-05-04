package loader

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/tinygo-org/tinygo/cgo"
	"golang.org/x/tools/go/packages"
)

// Program holds all packages and some metadata about the program as a whole.
type Program struct {
	Build        *build.Context
	Tests        bool
	Packages     map[string]*Package
	MainPkg      *Package
	sorted       []*Package
	fset         *token.FileSet
	TypeChecker  types.Config
	Dir          string // current working directory (for error reporting)
	TINYGOROOT   string // root of the TinyGo installation or root of the source code
	CFlags       []string
	ClangHeaders string
}

// Package holds a loaded package, its imports, and its parsed files.
type Package struct {
	*Program
	*packages.Package
	Files []*ast.File
	Pkg   *types.Package
	types.Info
}

// Load loads the given package with all dependencies (including the runtime
// package). Call .Parse() afterwards to parse all Go files (including CGo
// processing, if necessary).
func (p *Program) Load(importPath string) error {
	if p.Packages == nil {
		p.Packages = make(map[string]*Package)
	}

	err := p.loadPackage(importPath)
	if err != nil {
		return err
	}
	p.MainPkg = p.sorted[len(p.sorted)-1]
	if _, ok := p.Packages["runtime"]; !ok {
		// The runtime package wasn't loaded. Although `go list -deps` seems to
		// return the full dependency list, there is no way to get those
		// packages from the go/packages package. Therefore load the runtime
		// manually and add it to the list of to-be-compiled packages
		// (duplicates are already filtered).
		return p.loadPackage("runtime")
	}
	return nil
}

func (p *Program) loadPackage(importPath string) error {
	cgoEnabled := "0"
	if p.Build.CgoEnabled {
		cgoEnabled = "1"
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode:       packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps,
		Env:        append(os.Environ(), "GOROOT="+p.Build.GOROOT, "GOOS="+p.Build.GOOS, "GOARCH="+p.Build.GOARCH, "CGO_ENABLED="+cgoEnabled),
		BuildFlags: []string{"-tags", strings.Join(p.Build.BuildTags, " ")},
		Tests:      p.Tests,
	}, importPath)
	if err != nil {
		return err
	}
	var pkg *packages.Package
	if p.Tests {
		// We need the second package. Quoting from the docs:
		// > For example, when using the go command, loading "fmt" with Tests=true
		// > returns four packages, with IDs "fmt" (the standard package),
		// > "fmt [fmt.test]" (the package as compiled for the test),
		// > "fmt_test" (the test functions from source files in package fmt_test),
		// > and "fmt.test" (the test binary).
		pkg = pkgs[1]
	} else {
		if len(pkgs) != 1 {
			return fmt.Errorf("expected exactly one package while importing %s, got %d", importPath, len(pkgs))
		}
		pkg = pkgs[0]
	}
	var importError *Errors
	var addPackages func(pkg *packages.Package)
	addPackages = func(pkg *packages.Package) {
		if _, ok := p.Packages[pkg.PkgPath]; ok {
			return
		}
		pkg2 := p.newPackage(pkg)
		p.Packages[pkg.PkgPath] = pkg2
		if len(pkg.Errors) != 0 {
			if importError != nil {
				// There was another error reported already. Do not report
				// errors from multiple packages at once.
				return
			}
			importError = &Errors{
				Pkg: pkg2,
			}
			for _, err := range pkg.Errors {
				importError.Errs = append(importError.Errs, err)
			}
			return
		}

		// Get the list of imports (sorted alphabetically).
		names := make([]string, 0, len(pkg.Imports))
		for name := range pkg.Imports {
			names = append(names, name)
		}
		sort.Strings(names)

		// Add all the imports.
		for _, name := range names {
			addPackages(pkg.Imports[name])
		}

		p.sorted = append(p.sorted, pkg2)
	}
	addPackages(pkg)
	if importError != nil {
		return importError
	}
	return nil
}

// newPackage instantiates a new *Package object with initialized members.
func (p *Program) newPackage(pkg *packages.Package) *Package {
	return &Package{
		Program: p,
		Package: pkg,
		Info: types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object),
			Scopes:     make(map[ast.Node]*types.Scope),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
		},
	}
}

// Sorted returns a list of all packages, sorted in a way that no packages come
// before the packages they depend upon.
func (p *Program) Sorted() []*Package {
	return p.sorted
}

// Parse parses all packages and typechecks them.
//
// The returned error may be an Errors error, which contains a list of errors.
//
// Idempotent.
func (p *Program) Parse() error {
	// Parse all packages.
	for _, pkg := range p.Sorted() {
		err := pkg.Parse()
		if err != nil {
			return err
		}
	}

	if p.Tests {
		err := p.swapTestMain()
		if err != nil {
			return err
		}
	}

	// Typecheck all packages.
	for _, pkg := range p.Sorted() {
		err := pkg.Check()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Program) swapTestMain() error {
	var tests []string

	isTestFunc := func(f *ast.FuncDecl) bool {
		// TODO: improve signature check
		if strings.HasPrefix(f.Name.Name, "Test") && f.Name.Name != "TestMain" {
			return true
		}
		return false
	}
	for _, f := range p.MainPkg.Files {
		for i, d := range f.Decls {
			switch v := d.(type) {
			case *ast.FuncDecl:
				if isTestFunc(v) {
					tests = append(tests, v.Name.Name)
				}
				if v.Name.Name == "main" {
					// Remove main
					if len(f.Decls) == 1 {
						f.Decls = make([]ast.Decl, 0)
					} else {
						f.Decls[i] = f.Decls[len(f.Decls)-1]
						f.Decls = f.Decls[:len(f.Decls)-1]
					}
				}
			}
		}
	}

	// TODO: Check if they defined a TestMain and call it instead of testing.TestMain
	const mainBody = `package main

import (
	"testing"
)

func main () {
	m := &testing.M{
		Tests: []testing.TestToCall{
{{range .TestFunctions}}
			{Name: "{{.}}", Func: {{.}}},
{{end}}
		},
	}

	testing.TestMain(m)
}
`
	tmpl := template.Must(template.New("testmain").Parse(mainBody))
	b := bytes.Buffer{}
	tmplData := struct {
		TestFunctions []string
	}{
		TestFunctions: tests,
	}

	err := tmpl.Execute(&b, tmplData)
	if err != nil {
		return err
	}
	path := filepath.Join(p.MainPkg.Dir, "$testmain.go")

	if p.fset == nil {
		p.fset = token.NewFileSet()
	}

	newMain, err := parser.ParseFile(p.fset, path, b.Bytes(), parser.AllErrors)
	if err != nil {
		return err
	}
	p.MainPkg.Files = append(p.MainPkg.Files, newMain)

	return nil
}

// parseFile is a wrapper around parser.ParseFile.
func (p *Program) parseFile(path string, mode parser.Mode) (*ast.File, error) {
	if p.fset == nil {
		p.fset = token.NewFileSet()
	}

	rd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	relpath := path
	if filepath.IsAbs(path) {
		rp, err := filepath.Rel(p.Dir, path)
		if err == nil {
			relpath = rp
		}
	}
	return parser.ParseFile(p.fset, relpath, rd, mode)
}

// Parse parses and typechecks this package.
//
// Idempotent.
func (p *Package) Parse() error {
	if len(p.Files) != 0 {
		return nil
	}

	// Load the AST.
	// TODO: do this in parallel.
	if p.PkgPath == "unsafe" {
		// Special case for the unsafe package. Don't even bother loading
		// the files.
		p.Pkg = types.Unsafe
		return nil
	}

	files, err := p.parseFiles()
	if err != nil {
		return err
	}
	p.Files = files

	return nil
}

// Check runs the package through the typechecker. The package must already be
// loaded and all dependencies must have been checked already.
//
// Idempotent.
func (p *Package) Check() error {
	if p.Pkg != nil {
		return nil
	}

	var typeErrors []error
	checker := p.TypeChecker
	checker.Error = func(err error) {
		typeErrors = append(typeErrors, err)
	}

	// Do typechecking of the package.
	checker.Importer = p

	typesPkg, err := checker.Check(p.PkgPath, p.fset, p.Files, &p.Info)
	if err != nil {
		if err, ok := err.(Errors); ok {
			return err
		}
		return Errors{p, typeErrors}
	}
	p.Pkg = typesPkg
	return nil
}

// parseFiles parses the loaded list of files and returns this list.
func (p *Package) parseFiles() ([]*ast.File, error) {
	// TODO: do this concurrently.
	var files []*ast.File
	var fileErrs []error

	var cgoFiles []*ast.File
	for _, file := range p.GoFiles {
		f, err := p.parseFile(file, parser.ParseComments)
		if err != nil {
			fileErrs = append(fileErrs, err)
			continue
		}
		if err != nil {
			fileErrs = append(fileErrs, err)
			continue
		}
		for _, importSpec := range f.Imports {
			if importSpec.Path.Value == `"C"` {
				cgoFiles = append(cgoFiles, f)
			}
		}
		files = append(files, f)
	}
	if len(cgoFiles) != 0 {
		cflags := append(p.CFlags, "-I"+filepath.Dir(p.GoFiles[0]))
		if p.ClangHeaders != "" {
			cflags = append(cflags, "-Xclang", "-internal-isystem", "-Xclang", p.ClangHeaders)
		}
		generated, errs := cgo.Process(files, p.Program.Dir, p.fset, cflags)
		if errs != nil {
			fileErrs = append(fileErrs, errs...)
		}
		files = append(files, generated)
	}
	if len(fileErrs) != 0 {
		return nil, Errors{p, fileErrs}
	}

	return files, nil
}

// Import implements types.Importer. It loads and parses packages it encounters
// along the way, if needed.
func (p *Package) Import(to string) (*types.Package, error) {
	if to == "unsafe" {
		return types.Unsafe, nil
	}
	if _, ok := p.Imports[to]; ok {
		return p.Packages[p.Imports[to].PkgPath].Pkg, nil
	} else {
		return nil, errors.New("package not imported: " + to)
	}
}
