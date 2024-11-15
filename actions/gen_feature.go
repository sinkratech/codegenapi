package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	_ "embed"

	"github.com/urfave/cli/v2"
)

var (
	//go:embed genfile/api/deps.tmpl
	depsTmpl string

	//go:embed genfile/api/entrypoint.tmpl
	entrypointTmpl string
)

type GenFeatureData struct {
	PackageName string
}

func GenFeature(ctx *cli.Context) error {
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

	data := GenFeatureData{PackageName: packageName}

	err = generateDepsTmpl(targetDir, data)
	if err != nil {
		return fmt.Errorf("cannot generate deps.go for %s: %w", targetDir, err)
	}

	err = generateEntrypointTmpl(targetDir, data)
	if err != nil {
		return fmt.Errorf("cannot generate entrypoint.go %s: %w", targetDir, err)
	}

	return nil
}

func generateDepsTmpl(targetDir string, data GenFeatureData) error {
	t, err := template.New("deps").Parse(depsTmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(targetDir, "deps.go"))
	if err != nil {
		return err
	}

	defer file.Close()

	err = t.Execute(file, data)
	return err
}

func generateEntrypointTmpl(targetDir string, data GenFeatureData) error {
	t, err := template.New("entrypoint").Parse(entrypointTmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(targetDir, "entrypoint.go"))
	if err != nil {
		return err
	}

	defer file.Close()

	err = t.Execute(file, data)
	return err
}

func findRootProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}

		dir = parentDir
	}

	return "", fmt.Errorf("cannot find parent project (which contains go.mod file)")
}
