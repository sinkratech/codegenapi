package actions

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	_ "embed"

	"github.com/golang-cz/textcase"
	"github.com/urfave/cli/v2"
)

var (
	//go:embed genfile/api/opt.tmpl
	optTmpl string
)

type GenInterfaceMethod struct {
	InterfaceName string
	VarName       string
}

type GenInterfaceData struct {
	PackageName string
	Method      []GenInterfaceMethod
}

func GenInterfaceImpl(ctx *cli.Context) error {
	packageName := ctx.Args().First()
	if packageName == "" {
		return fmt.Errorf("missing argument for package name")
	}

	rootDir, err := findRootProjectDir()
	if err != nil {
		return fmt.Errorf("cannot find root project dir: %w", err)
	}

	targetDir := filepath.Join(rootDir, "api", packageName)
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", targetDir, err)
	}

	interfaces, err := getAllInterfacesName(targetDir)
	if err != nil {
		return fmt.Errorf("cannot get interfaces: %w", err)
	}

	var data GenInterfaceData
	data.PackageName = packageName

	for it := range slices.Values(interfaces) {
		method := GenInterfaceMethod{
			InterfaceName: textcase.PascalCase(it),
			VarName:       textcase.CamelCase(it),
		}

		data.Method = append(data.Method, method)
	}

	err = generateDepsWithInterfaceInject(filepath.Join(targetDir, "deps.go"), data)
	if err != nil {
		return fmt.Errorf("cannot generate deps with interface inject: %w", err)
	}

	return nil
}

func generateDepsWithInterfaceInject(targetFile string, data GenInterfaceData) error {
	t, err := template.New("opt").Parse(optTmpl)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	err = t.Execute(file, data)
	return err
}

func getAllInterfacesName(targetDir string) ([]string, error) {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", targetDir, err)
	}

	var interfaces []string

	for entry := range slices.Values(entries) {
		filename := entry.Name()
		if filename != "deps.go" && filename != "entrypoint.go" {
			it, err := astParseOnlyTypeInterface(filepath.Join(targetDir, filename))
			if err != nil {
				return nil, err
			}

			interfaces = append(interfaces, it...)
		}
	}

	return interfaces, nil
}

func astParseOnlyTypeInterface(targetFile string) ([]string, error) {
	fileset := token.NewFileSet()
	node, err := parser.ParseFile(fileset, targetFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var interfaces []string

	ast.Inspect(node, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		if decl.Tok != token.TYPE {
			return true
		}

		// Use interface only for global scope
		if isInsideFunction(decl, node) {
			return true
		}

		for spec := range slices.Values(decl.Specs) {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			_, ok = typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			interfaces = append(interfaces, typeSpec.Name.Name)
		}

		return true
	})

	return interfaces, nil
}

func isInsideFunction(decl *ast.GenDecl, file *ast.File) bool {
	for f := range slices.Values(file.Decls) {
		if funcDecl, ok := f.(*ast.FuncDecl); ok {
			if decl.Pos() > funcDecl.Pos() && decl.Pos() < funcDecl.End() {
				return true
			}
		}
	}
	return false
}
