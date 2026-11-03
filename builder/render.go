package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

var pageNames = []string{
	"index",
	"chapters", "chapter", "stage",
	"techniques", "technique",
	"syntaxes", "syntax",
}

type Crumb struct {
	Label string
	URL   string
}

// Page はテンプレートに渡すデータ。Root はそのページからサイトのルートへの相対パス。
type Page struct {
	Root    string
	Title   string
	Section string
	Crumbs  []Crumb
	Site    *Site

	Chapter   *Chapter
	Stage     *Stage
	Technique *Technique
	Syntax    *Syntax
}

var funcs = template.FuncMap{
	// dict は部分テンプレートに複数の値を渡すために使う。
	"dict": func(kv ...any) (map[string]any, error) {
		if len(kv)%2 != 0 {
			return nil, fmt.Errorf("dict: 引数はキーと値のペアで指定してください")
		}
		m := make(map[string]any, len(kv)/2)
		for i := 0; i < len(kv); i += 2 {
			k, ok := kv[i].(string)
			if !ok {
				return nil, fmt.Errorf("dict: キーは文字列にしてください: %v", kv[i])
			}
			m[k] = kv[i+1]
		}
		return m, nil
	},
}

func loadTemplates() (map[string]*template.Template, error) {
	base, err := template.New("").Funcs(funcs).ParseFS(templateFS, "templates/layout.html", "templates/partials.html")
	if err != nil {
		return nil, err
	}
	tmpls := map[string]*template.Template{}
	for _, name := range pageNames {
		t, err := template.Must(base.Clone()).ParseFS(templateFS, "templates/"+name+".html")
		if err != nil {
			return nil, err
		}
		tmpls[name] = t
	}
	return tmpls, nil
}

func build(site *Site, out string) error {
	tmpls, err := loadTemplates()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(out); err != nil {
		return err
	}

	write := func(name, path string, p Page) error {
		p.Root = strings.Repeat("../", strings.Count(path, "/"))
		if p.Root == "" {
			p.Root = "./"
		}
		p.Site = site
		var buf bytes.Buffer
		if err := tmpls[name].ExecuteTemplate(&buf, "layout", p); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		dir := filepath.Join(out, filepath.FromSlash(path))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dir, "index.html"), buf.Bytes(), 0o644)
	}

	chaptersCrumb := Crumb{"チャプター", "chapters/"}
	techniquesCrumb := Crumb{"テクニック", "techniques/"}
	syntaxCrumb := Crumb{"文法", "syntax/"}

	if err := write("index", "", Page{}); err != nil {
		return err
	}

	if err := write("chapters", "chapters/", Page{Title: "チャプター", Section: "chapters",
		Crumbs: []Crumb{{Label: "チャプター"}}}); err != nil {
		return err
	}
	for _, c := range site.Chapters {
		if err := write("chapter", c.URL(), Page{Title: c.DisplayTitle(), Section: "chapters", Chapter: c,
			Crumbs: []Crumb{chaptersCrumb, {Label: c.Name()}}}); err != nil {
			return err
		}
		for _, s := range c.Stages {
			if err := write("stage", s.URL(), Page{Title: s.Label(), Section: "chapters", Chapter: c, Stage: s,
				Crumbs: []Crumb{chaptersCrumb, {c.Name(), c.URL()}, {Label: s.ID()}}}); err != nil {
				return err
			}
		}
	}

	if err := write("techniques", "techniques/", Page{Title: "テクニック総覧", Section: "techniques",
		Crumbs: []Crumb{{Label: "テクニック"}}}); err != nil {
		return err
	}
	for _, t := range site.Techniques {
		if err := write("technique", t.URL(), Page{Title: t.Title, Section: "techniques", Technique: t,
			Crumbs: []Crumb{techniquesCrumb, {Label: t.Title}}}); err != nil {
			return err
		}
	}

	if err := write("syntaxes", "syntax/", Page{Title: "文法総覧", Section: "syntax",
		Crumbs: []Crumb{{Label: "文法"}}}); err != nil {
		return err
	}
	for _, sx := range site.Syntaxes {
		if err := write("syntax", sx.URL(), Page{Title: sx.DisplayTitle(), Section: "syntax", Syntax: sx,
			Crumbs: []Crumb{syntaxCrumb, {Label: sx.Chapter.Name()}}}); err != nil {
			return err
		}
	}

	return copyStatic(out)
}

func copyStatic(out string) error {
	return fs.WalkDir(staticFS, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dst := filepath.Join(out, filepath.FromSlash(path))
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		data, err := staticFS.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
}
