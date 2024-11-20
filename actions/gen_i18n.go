package actions

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"

	_ "embed"

	"github.com/golang-cz/textcase"
	"github.com/urfave/cli/v2"
)

var (
	argsRegex = regexp.MustCompile(`%[0-9]*[vTsqbdoxXeEgGfnc%]`)

	//go:embed genfile/i18n/base.tmpl
	i18nTmpl string
)

type FileLanguageContent struct {
	Lang     string            `json:"lang"`
	Default  bool              `json:"default"`
	Contents map[string]string `json:"contents"`
}

type ParsedLanguageContent struct {
	Lang          string
	RawMethodName string
	MethodName    string
	Message       string
	ContainArgs   bool
	Default       bool
}

type LanguageMap = map[string]map[string]string

type Method struct {
	Name        string
	RawName     string
	ContainArgs bool
}

type GenI18nData struct {
	DefaultLang string
	Methods     []Method
	LanguageMap LanguageMap
}

func GenI18n(ctx *cli.Context) error {
	dir, err := findRootProjectDir()
	if err != nil {
		return fmt.Errorf("cannot find root project dir: %w", err)
	}

	targetDir := filepath.Join(dir, "i18n")
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", targetDir, err)
	}

	byteContent, err := os.ReadFile(filepath.Join(dir, "i18n.json"))
	if err != nil {
		return fmt.Errorf("cannot read i18n.json: %w", err)
	}

	var langContent []FileLanguageContent
	err = json.Unmarshal(byteContent, &langContent)
	if err != nil {
		return fmt.Errorf("cannot parse i18n.json: %w", err)
	}

	finalContent, err := transformIntoFinalLangContent(langContent)
	if err != nil {
		return fmt.Errorf("cannot transform content: %w", err)
	}

	data := GenI18nData{
		DefaultLang: getDefaultLang(finalContent),
		Methods:     getUniqueMethods(finalContent),
		LanguageMap: rearrangeByLanguage(finalContent),
	}

	err = generateI18nTmpl(dir, data)
	if err != nil {
		return fmt.Errorf("cannot generate i18n to i18n.go: %w", err)
	}

	return nil
}

func generateI18nTmpl(targetDir string, data GenI18nData) error {
	t, err := template.New("i18n").Parse(i18nTmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(targetDir, "i18n", "i18n.go"))
	if err != nil {
		return err
	}

	defer file.Close()

	err = t.Execute(file, data)
	return err
}

func transformIntoFinalLangContent(langContent []FileLanguageContent) ([]ParsedLanguageContent, error) {
	var finals []ParsedLanguageContent

	for it := range slices.Values(langContent) {
		for k, v := range it.Contents {
			splittedKey := strings.Split(k, ".")
			if len(splittedKey) != 2 {
				return nil, fmt.Errorf("bad format for %s", k)
			}

			feature, method := splittedKey[0], splittedKey[1]
			var final ParsedLanguageContent

			final.Lang = it.Lang
			final.Default = it.Default
			final.RawMethodName = k
			final.MethodName = fmt.Sprintf("%s_%s", textcase.PascalCase(feature), textcase.PascalCase(method))
			final.Message = v
			final.ContainArgs = argsRegex.MatchString(v)

			finals = append(finals, final)
		}
	}

	return finals, nil
}

func getDefaultLang(contents []ParsedLanguageContent) string {
	lang := contents[0].Lang

	for it := range slices.Values(contents) {
		if it.Default {
			return it.Lang
		}
	}

	return lang
}

func getUniqueMethods(contents []ParsedLanguageContent) []Method {
	uniqueMethods := make(map[string]Method)

	for it := range slices.Values(contents) {
		uniqueMethods[it.MethodName] = Method{
			Name:        it.MethodName,
			RawName:     it.RawMethodName,
			ContainArgs: it.ContainArgs,
		}
	}

	return slices.Collect(maps.Values(uniqueMethods))
}

func rearrangeByLanguage(finalLangContents []ParsedLanguageContent) LanguageMap {
	langMap := make(LanguageMap)

	for it := range slices.Values(finalLangContents) {
		lang := it.Lang
		key := it.RawMethodName
		val := it.Message

		mapContent, ok := langMap[lang]
		if !ok {
			langMap[lang] = map[string]string{
				key: val,
			}
		} else {
			mapContent[key] = val
		}
	}

	return langMap
}
