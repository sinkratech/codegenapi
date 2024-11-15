package actions

import (
	_ "embed"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_astParseOnlyTypeInterface(t *testing.T) {
	target, err := findRootProjectDir()
	require.NoError(t, err)

	path := filepath.Join(target, "actions", "testfile", "service.go")

	interfaces, err := astParseOnlyTypeInterface(path)
	require.NoError(t, err)

	require.NotEmpty(t, interfaces)
	require.Contains(t, interfaces, "Foo")
	require.Contains(t, interfaces, "Bar")
	require.NotContains(t, interfaces, "Inline")
}

func Test_getAllInterfacesName(t *testing.T) {
	target, err := findRootProjectDir()
	require.NoError(t, err)

	path := filepath.Join(target, "actions", "testfile")

	interfaces, err := getAllInterfacesName(path)
	require.NoError(t, err)

	require.NotEmpty(t, interfaces)
	require.Contains(t, interfaces, "Foo")
	require.Contains(t, interfaces, "Bar")
	require.NotContains(t, interfaces, "Inline")
}
