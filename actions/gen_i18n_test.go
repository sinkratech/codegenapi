package actions

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

var parsedLangContent = []ParsedLanguageContent{
	{
		Lang:          "id",
		RawMethodName: "foo.Bar",
		MethodName:    "Foo_Bar",
		Message:       "bar",
		ContainArgs:   false,
		Default:       false,
	},
	{
		Lang:          "en",
		RawMethodName: "foo.Bar",
		MethodName:    "Foo_Bar",
		Message:       "baz",
		ContainArgs:   false,
		Default:       true,
	},
}

func Test_transformIntoFinalLangContent(t *testing.T) {
	langContent := []FileLanguageContent{
		{
			Lang:    "id",
			Default: false,
			Contents: map[string]string{
				"foo.Bar": "bar",
			},
		},
		{
			Lang:    "en",
			Default: true,
			Contents: map[string]string{
				"foo.Bar": "baz",
			},
		},
	}

	got, err := transformIntoFinalLangContent(langContent)
	require.NoError(t, err)
	require.Equal(t, parsedLangContent, got)
}

func Test_getDefaultLang(t *testing.T) {
	got := getDefaultLang(parsedLangContent)
	require.Equal(t, "en", got)
}

func Test_getUniqueMethod(t *testing.T) {
	got := getUniqueMethods(parsedLangContent)
	expected := []Method{
		{
			Name:        "Foo_Bar",
			RawName:     "foo.Bar",
			ContainArgs: false,
		},
	}

	require.Equal(t, expected, got)
}

func Test_rearrangeByLanguage(t *testing.T) {
	got := rearrangeByLanguage(parsedLangContent)
	expected := LanguageMap{
		"id": {"foo.Bar": "bar"},
		"en": {"foo.Bar": "baz"},
	}

	require.Equal(t, expected, got)
}
