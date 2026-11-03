package main

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// 結合テスト: コンテンツを読み込んでサイトを生成し、出力を検査する。

func TestBuildFixture(t *testing.T) {
	out := t.TempDir()
	if err := build(loadFixture(t), out, buildOptions{}); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{
		"index.html",
		"chapters/index.html",
		"chapters/1/index.html",
		"chapters/1/1-1/index.html",
		"chapters/1/1-2/index.html",
		"chapters/2/2-1/index.html",
		"techniques/index.html",
		"techniques/swap/index.html",
		"techniques/mark/index.html",
		"syntax/index.html",
		"syntax/1/index.html",
		"syntax/2/index.html",
		"static/style.css",
		"static/site.js",
	} {
		if _, err := os.Stat(filepath.Join(out, filepath.FromSlash(name))); err != nil {
			t.Errorf("%s が出力されていない", name)
		}
	}

	checkSite(t, out)

	stage := readFile(t, out, "chapters/1/1-2/index.html")
	for _, want := range []string{
		`<code class="language-ab">`,       // 解答のコードブロック
		`href="../../../techniques/swap/"`, // 使ったテクニック
		`href="../../../syntax/1/"`,        // 文法
		`href="../../../chapters/1/1-1/"`,  // 前のステージ
		`href="../../../chapters/2/2-1/"`,  // 次のステージ（チャプターをまたぐ）
		`2行 / 10ステップ / 挑戦中`,
	} {
		if !strings.Contains(stage, want) {
			t.Errorf("1-2 のページに %q がない", want)
		}
	}

	technique := readFile(t, out, "techniques/mark/index.html")
	for _, want := range []string{`href="../../syntax/1/"`, `href="../../syntax/2/"`, `href="../../chapters/2/2-1/"`} {
		if !strings.Contains(technique, want) {
			t.Errorf("mark のページに %q がない", want)
		}
	}
}

func TestBuildRemovesStaleFiles(t *testing.T) {
	out := t.TempDir()
	stale := filepath.Join(out, "old.html")
	if err := os.WriteFile(stale, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := build(loadFixture(t), out, buildOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(stale); err == nil {
		t.Error("前回の出力が残っている")
	}
}

// TestBuildContent は実際の content/ がエラーなく生成でき、リンク切れがないことを確かめる。
func TestBuildContent(t *testing.T) {
	site, err := loadSite(filepath.Join("..", "content"))
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := build(site, out, buildOptions{}); err != nil {
		t.Fatal(err)
	}
	checkSite(t, out)
}

var (
	linkRe = regexp.MustCompile(`(?:href|src)="([^"]*)"`)
	idRe   = regexp.MustCompile(`id="([^"]+)"`)
)

// checkSite は生成された全ページについて、テンプレートの不具合とリンク切れを検査する。
func checkSite(t *testing.T, out string) {
	t.Helper()
	pages := map[string]string{}
	err := filepath.WalkDir(out, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(p) != ".html" {
			return err
		}
		rel, _ := filepath.Rel(out, p)
		pages[filepath.ToSlash(rel)] = readFile(t, out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for page, html := range pages {
		// html/template が危険と判断した値や、存在しないフィールドの参照
		for _, bad := range []string{"ZgotmplZ", "<no value>"} {
			if strings.Contains(html, bad) {
				t.Errorf("%s: %q が含まれている", page, bad)
			}
		}

		for _, m := range linkRe.FindAllStringSubmatch(html, -1) {
			link := m[1]
			if strings.Contains(link, "://") || strings.HasPrefix(link, "mailto:") {
				continue
			}
			target, frag, _ := strings.Cut(link, "#")
			resolved := page
			if target != "" {
				resolved = path.Join(path.Dir(page), target)
				if strings.HasSuffix(target, "/") {
					resolved = path.Join(resolved, "index.html")
				}
			}
			if strings.HasPrefix(resolved, "..") {
				t.Errorf("%s: %q がサイトの外を指している", page, link)
				continue
			}
			if _, err := os.Stat(filepath.Join(out, filepath.FromSlash(resolved))); err != nil {
				t.Errorf("%s: リンク切れ %q", page, link)
				continue
			}
			if frag != "" && !hasID(pages[resolved], frag) {
				t.Errorf("%s: %q のアンカーがない", page, link)
			}
		}
	}
}

func hasID(html, id string) bool {
	for _, m := range idRe.FindAllStringSubmatch(html, -1) {
		if m[1] == id {
			return true
		}
	}
	return false
}

func readFile(t *testing.T, out, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestBuildOGP(t *testing.T) {
	t.Run("base-url あり", func(t *testing.T) {
		out := t.TempDir()
		if err := build(loadFixture(t), out, buildOptions{BaseURL: "https://example.com/book"}); err != nil {
			t.Fatal(err)
		}
		stage := readFile(t, out, "chapters/1/1-2/index.html")
		for _, want := range []string{
			`<meta property="og:title" content="1-2 二問目 | A=B Unofficial Answer Book">`,
			`<meta property="og:url" content="https://example.com/book/chapters/1/1-2/">`,
			`<meta property="og:image" content="https://example.com/book/static/ogp-unofficial.png">`,
			`<meta name="twitter:card" content="summary_large_image">`,
		} {
			if !strings.Contains(stage, want) {
				t.Errorf("1-2 のページに %q がない", want)
			}
		}
		top := readFile(t, out, "index.html")
		for _, want := range []string{
			`<meta property="og:title" content="A=B Unofficial Answer Book">`,
			`<meta property="og:url" content="https://example.com/book/">`,
		} {
			if !strings.Contains(top, want) {
				t.Errorf("トップページに %q がない", want)
			}
		}
		if _, err := os.Stat(filepath.Join(out, "static", "ogp-unofficial.png")); err != nil {
			t.Error("OGP 画像が出力されていない")
		}
	})

	t.Run("base-url なし", func(t *testing.T) {
		out := t.TempDir()
		if err := build(loadFixture(t), out, buildOptions{}); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(readFile(t, out, "index.html"), "og:") {
			t.Error("base-url なしなのに OGP タグが出ている")
		}
	})
}

func TestBuildFooterLinks(t *testing.T) {
	const game = `href="https://store.steampowered.com/app/1720850/AB/"`

	t.Run("repo-url あり", func(t *testing.T) {
		out := t.TempDir()
		if err := build(loadFixture(t), out, buildOptions{RepoURL: "https://github.com/example/book"}); err != nil {
			t.Fatal(err)
		}
		for _, page := range []string{"index.html", "chapters/1/1-2/index.html"} {
			html := readFile(t, out, page)
			for _, want := range []string{game, `href="https://github.com/example/book"`} {
				if !strings.Contains(html, want) {
					t.Errorf("%s に %q がない", page, want)
				}
			}
		}
	})

	t.Run("repo-url なし", func(t *testing.T) {
		out := t.TempDir()
		if err := build(loadFixture(t), out, buildOptions{}); err != nil {
			t.Fatal(err)
		}
		html := readFile(t, out, "index.html")
		if !strings.Contains(html, game) {
			t.Errorf("元のゲームへのリンク %q がない", game)
		}
		if strings.Contains(html, "ソースコード") {
			t.Error("repo-url なしなのに GitHub へのリンクが出ている")
		}
	})
}
