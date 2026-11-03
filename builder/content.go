package main

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

var (
	chapterDirRe = regexp.MustCompile(`^chapter(\d+)$`)
	stageFileRe  = regexp.MustCompile(`^(\d+)-(\d+)(?:_(.+))?\.md$`)
	syntaxFileRe = regexp.MustCompile(`^chapter(\d+)\.md$`)
)

var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

type Site struct {
	Chapters   []*Chapter
	Stages     []*Stage
	Techniques []*Technique
	Syntaxes   []*Syntax

	chapters   map[int]*Chapter
	techniques map[string]*Technique
}

type Chapter struct {
	Num        int
	Title      string
	Body       template.HTML
	Stages     []*Stage
	Syntax     *Syntax
	Techniques []*Technique // このチャプターから使えるようになったテクニック
	Prev, Next *Chapter
}

type Stage struct {
	Chapter    *Chapter
	Num        int
	Title      string
	Lines      *int
	Steps      *int
	Status     string
	Techniques []*Technique
	Body       template.HTML
	Prev, Next *Stage

	techniqueSlugs []string
}

type Technique struct {
	Slug     string
	Title    string
	Since    *Chapter
	Syntaxes []*Syntax
	Body     template.HTML
	UsedIn   []*Stage

	syntaxNums []int
}

type Syntax struct {
	Chapter    *Chapter
	Title      string
	Body       template.HTML
	Techniques []*Technique // この文法に関連するテクニック
	Prev, Next *Syntax
}

func (c *Chapter) URL() string  { return fmt.Sprintf("chapters/%d/", c.Num) }
func (c *Chapter) Name() string { return fmt.Sprintf("チャプター%d", c.Num) }

func (c *Chapter) DisplayTitle() string {
	if c.Title == "" {
		return c.Name()
	}
	return c.Name() + ": " + c.Title
}

func (c *Chapter) Label() string { return c.DisplayTitle() }

func (c *Chapter) DoneCount() int {
	n := 0
	for _, s := range c.Stages {
		if s.Status == "done" {
			n++
		}
	}
	return n
}

// UsedTechniques はこのチャプターのステージで使ったテクニックを重複なしで返す。
func (c *Chapter) UsedTechniques() []*Technique {
	var ts []*Technique
	for _, s := range c.Stages {
		for _, t := range s.Techniques {
			if !slices.Contains(ts, t) {
				ts = append(ts, t)
			}
		}
	}
	return ts
}

func (s *Stage) ID() string    { return fmt.Sprintf("%d-%d", s.Chapter.Num, s.Num) }
func (s *Stage) URL() string   { return s.Chapter.URL() + s.ID() + "/" }
func (s *Stage) Label() string { return s.ID() + " " + s.Title }

func (s *Stage) LinesText() string { return intText(s.Lines) }
func (s *Stage) StepsText() string { return intText(s.Steps) }

func (s *Stage) StatusLabel() string {
	switch s.Status {
	case "done":
		return "クリア"
	case "wip":
		return "挑戦中"
	default:
		return "未着手"
	}
}

// EarlierSyntaxes はこのステージより前のチャプターで追加された文法を返す。
func (s *Stage) EarlierSyntaxes(site *Site) []*Syntax {
	var ss []*Syntax
	for _, sx := range site.Syntaxes {
		if sx.Chapter.Num < s.Chapter.Num {
			ss = append(ss, sx)
		}
	}
	return ss
}

func (t *Technique) URL() string { return "techniques/" + t.Slug + "/" }

func (s *Syntax) URL() string  { return fmt.Sprintf("syntax/%d/", s.Chapter.Num) }
func (s *Syntax) Name() string { return s.Chapter.Name() + "の文法" }

func (s *Syntax) Label() string { return s.Name() }

func (s *Syntax) DisplayTitle() string {
	if s.Title == "" {
		return s.Name()
	}
	return s.Name() + ": " + s.Title
}

func (s *Site) StageCount() int { return len(s.Stages) }

func (s *Site) DoneCount() int {
	n := 0
	for _, c := range s.Chapters {
		n += c.DoneCount()
	}
	return n
}

func intText(p *int) string {
	if p == nil {
		return "-"
	}
	return strconv.Itoa(*p)
}

func loadSite(dir string) (*Site, error) {
	site := &Site{chapters: map[int]*Chapter{}, techniques: map[string]*Technique{}}
	if err := site.loadChapters(filepath.Join(dir, "chapters")); err != nil {
		return nil, err
	}
	if err := site.loadSyntaxes(filepath.Join(dir, "syntax")); err != nil {
		return nil, err
	}
	if err := site.loadTechniques(filepath.Join(dir, "techniques")); err != nil {
		return nil, err
	}
	if err := site.link(); err != nil {
		return nil, err
	}
	return site, nil
}

func (site *Site) chapter(num int) *Chapter {
	c, ok := site.chapters[num]
	if !ok {
		c = &Chapter{Num: num}
		site.chapters[num] = c
		site.Chapters = append(site.Chapters, c)
	}
	return c
}

func (site *Site) loadChapters(dir string) error {
	entries, err := readDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		m := chapterDirRe.FindStringSubmatch(e.Name())
		if !e.IsDir() || m == nil {
			continue
		}
		ch := site.chapter(atoi(m[1]))
		chDir := filepath.Join(dir, e.Name())
		files, err := readDir(chDir)
		if err != nil {
			return err
		}
		for _, f := range files {
			path := filepath.Join(chDir, f.Name())
			if f.Name() == "_chapter.md" {
				var meta struct {
					Title string `yaml:"title"`
				}
				if ch.Body, err = readDoc(path, &meta); err != nil {
					return err
				}
				ch.Title = meta.Title
				continue
			}
			m := stageFileRe.FindStringSubmatch(f.Name())
			if f.IsDir() || m == nil {
				continue
			}
			if atoi(m[1]) != ch.Num {
				return fmt.Errorf("%s: ファイル名のチャプター番号がフォルダ(%s)と一致しません", path, e.Name())
			}
			var meta struct {
				Title      string   `yaml:"title"`
				Lines      *int     `yaml:"lines"`
				Steps      *int     `yaml:"steps"`
				Status     string   `yaml:"status"`
				Techniques []string `yaml:"techniques"`
			}
			body, err := readDoc(path, &meta)
			if err != nil {
				return err
			}
			switch meta.Status {
			case "":
				meta.Status = "todo"
			case "done", "wip", "todo":
			default:
				return fmt.Errorf("%s: status は done / wip / todo のいずれかにしてください: %q", path, meta.Status)
			}
			if meta.Title == "" {
				meta.Title = m[3]
			}
			st := &Stage{
				Chapter:        ch,
				Num:            atoi(m[2]),
				Title:          meta.Title,
				Lines:          meta.Lines,
				Steps:          meta.Steps,
				Status:         meta.Status,
				Body:           body,
				techniqueSlugs: meta.Techniques,
			}
			ch.Stages = append(ch.Stages, st)
		}
	}
	return nil
}

func (site *Site) loadSyntaxes(dir string) error {
	entries, err := readDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		m := syntaxFileRe.FindStringSubmatch(e.Name())
		if e.IsDir() || m == nil {
			continue
		}
		var meta struct {
			Title string `yaml:"title"`
		}
		body, err := readDoc(filepath.Join(dir, e.Name()), &meta)
		if err != nil {
			return err
		}
		ch := site.chapter(atoi(m[1]))
		ch.Syntax = &Syntax{Chapter: ch, Title: meta.Title, Body: body}
		site.Syntaxes = append(site.Syntaxes, ch.Syntax)
	}
	return nil
}

func (site *Site) loadTechniques(dir string) error {
	entries, err := readDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		m := chapterDirRe.FindStringSubmatch(e.Name())
		if !e.IsDir() || m == nil {
			continue
		}
		ch := site.chapter(atoi(m[1]))
		chDir := filepath.Join(dir, e.Name())
		files, err := readDir(chDir)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() || filepath.Ext(f.Name()) != ".md" {
				continue
			}
			path := filepath.Join(chDir, f.Name())
			slug := strings.TrimSuffix(f.Name(), ".md")
			if _, dup := site.techniques[slug]; dup {
				return fmt.Errorf("%s: テクニック %q が重複しています", path, slug)
			}
			var meta struct {
				Title  string `yaml:"title"`
				Syntax []int  `yaml:"syntax"`
			}
			body, err := readDoc(path, &meta)
			if err != nil {
				return err
			}
			if meta.Title == "" {
				meta.Title = slug
			}
			t := &Technique{Slug: slug, Title: meta.Title, Since: ch, Body: body, syntaxNums: meta.Syntax}
			site.techniques[slug] = t
			site.Techniques = append(site.Techniques, t)
		}
	}
	return nil
}

// link は読み込んだコンテンツ同士の参照を解決し、並び順と前後リンクを決める。
func (site *Site) link() error {
	slices.SortFunc(site.Chapters, func(a, b *Chapter) int { return cmp.Compare(a.Num, b.Num) })
	slices.SortFunc(site.Syntaxes, func(a, b *Syntax) int { return cmp.Compare(a.Chapter.Num, b.Chapter.Num) })
	slices.SortFunc(site.Techniques, func(a, b *Technique) int {
		return cmp.Or(cmp.Compare(a.Since.Num, b.Since.Num), cmp.Compare(a.Slug, b.Slug))
	})

	for i, c := range site.Chapters {
		if i > 0 {
			c.Prev = site.Chapters[i-1]
			c.Prev.Next = c
		}
		slices.SortFunc(c.Stages, func(a, b *Stage) int { return cmp.Compare(a.Num, b.Num) })
		site.Stages = append(site.Stages, c.Stages...)
	}
	for i, s := range site.Stages {
		if i > 0 {
			s.Prev = site.Stages[i-1]
			s.Prev.Next = s
		}
	}
	for i, sx := range site.Syntaxes {
		if i > 0 {
			sx.Prev = site.Syntaxes[i-1]
			sx.Prev.Next = sx
		}
	}

	for _, t := range site.Techniques {
		t.Since.Techniques = append(t.Since.Techniques, t)
		if len(t.syntaxNums) == 0 {
			if sx := t.Since.Syntax; sx != nil {
				t.Syntaxes = []*Syntax{sx}
			}
		}
		for _, n := range t.syntaxNums {
			c, ok := site.chapters[n]
			if !ok || c.Syntax == nil {
				return fmt.Errorf("テクニック %q: チャプター%dの文法ファイル (content/syntax/chapter%02d.md) がありません", t.Slug, n, n)
			}
			t.Syntaxes = append(t.Syntaxes, c.Syntax)
		}
		for _, sx := range t.Syntaxes {
			sx.Techniques = append(sx.Techniques, t)
		}
	}

	for _, s := range site.Stages {
		for _, slug := range s.techniqueSlugs {
			t, ok := site.techniques[slug]
			if !ok {
				return fmt.Errorf("ステージ %s: テクニック %q が content/techniques にありません", s.ID(), slug)
			}
			s.Techniques = append(s.Techniques, t)
			t.UsedIn = append(t.UsedIn, s)
		}
	}
	return nil
}

func readDir(dir string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	return entries, err
}

// readDoc は frontmatter 付きの Markdown を読み、frontmatter を meta に、本文を HTML にする。
func readDoc(path string, meta any) (template.HTML, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	fm, body := splitFrontMatter(src)
	if fm != "" {
		dec := yaml.NewDecoder(strings.NewReader(fm))
		dec.KnownFields(true)
		if err := dec.Decode(meta); err != nil {
			return "", fmt.Errorf("%s: frontmatter: %w", path, err)
		}
	}
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(body), &buf); err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return template.HTML(buf.String()), nil
}

func splitFrontMatter(src []byte) (fm, body string) {
	s := strings.TrimPrefix(string(src), "\uFEFF")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return "", s
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			return strings.Join(lines[1:i], "\n"), strings.Join(lines[i+1:], "\n")
		}
	}
	return "", s
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
