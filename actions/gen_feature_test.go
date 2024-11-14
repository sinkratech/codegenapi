package actions

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_findRootProjectDir(t *testing.T) {
	packageName, err := findRootProjectDir()
	require.NoError(t, err)
	require.False(t, packageName == "github.com/sinkratech/codegen-api")
}

func Test_generateDepsTmpl(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(dir)
	})

	data := GenFeatureData{PackageName: "foosha"}

	err = generateDepsTmpl(dir, data)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dir, "deps.go"))
	require.NoError(t, err)

	require.Contains(t, string(content), "// DO NOT EDIT. Edit at your own risk.")
	require.Contains(t, string(content), "package foosha")
}

func Test_generateEntrypointTmpl(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(dir)
	})

	data := GenFeatureData{PackageName: "foosha"}

	err = generateEntrypointTmpl(dir, data)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dir, "entrypoint.go"))
	require.NoError(t, err)

	require.Contains(t, string(content), "package foosha")
	require.Contains(t, string(content), `import "github.com/danielgtaylor/huma/v2"`)
	require.Contains(t, string(content), `type Deps struct`)
	require.Contains(t, string(content), `func (d *Deps) Routes(router huma.API)`)
}
