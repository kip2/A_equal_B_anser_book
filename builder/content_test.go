package main

import (
	"slices"
	"strings"
	"testing"
)

func TestSplitFrontMatter(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantFM   string
		wantBody string
	}{
		{"frontmatter なし", "本文\n", "", "本文\n"},
		{"通常", "---\ntitle: a\n---\n本文\n", "title: a", "本文\n"},
		{"CRLF", "---\r\ntitle: a\r\n---\r\n本文\r\n", "title: a", "本文\n"},
		{"BOM 付き", "\uFEFF---\ntitle: a\n---\n本文", "title: a", "本文"},
		{"空の frontmatter", "---\n---\n本文", "", "本文"},
		{"閉じていない", "---\ntitle: a\n本文", "", "---\ntitle: a\n本文"},
		{"途中の --- は無視", "本文\n---\n続き", "", "本文\n---\n続き"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body := splitFrontMatter([]byte(tt.src))
			if fm != tt.wantFM || body != tt.wantBody {
				t.Errorf("got (%q, %q), want (%q, %q)", fm, body, tt.wantFM, tt.wantBody)
			}
		})
	}
}

func TestIntText(t *testing.T) {
	n := 12
	if got := intText(nil); got != "-" {
		t.Errorf("intText(nil) = %q", got)
	}
	if got := intText(&n); got != "12" {
		t.Errorf("intText(12) = %q", got)
	}
}

func TestStatusLabel(t *testing.T) {
	for status, want := range map[string]string{"done": "クリア", "wip": "挑戦中", "todo": "未着手"} {
		if got := (&Stage{Status: status}).StatusLabel(); got != want {
			t.Errorf("%s: got %q, want %q", status, got, want)
		}
	}
}

func TestDisplayTitle(t *testing.T) {
	c := &Chapter{Num: 3}
	if got := c.DisplayTitle(); got != "チャプター3" {
		t.Errorf("タイトルなし: %q", got)
	}
	c.Title = "応用"
	if got := c.DisplayTitle(); got != "チャプター3: 応用" {
		t.Errorf("タイトルあり: %q", got)
	}
	if got := (&Syntax{Chapter: c}).DisplayTitle(); got != "チャプター3の文法" {
		t.Errorf("文法: %q", got)
	}
}

func TestLoadSite(t *testing.T) {
	site := loadFixture(t)

	t.Run("チャプター", func(t *testing.T) {
		if len(site.Chapters) != 2 {
			t.Fatalf("チャプター数 = %d", len(site.Chapters))
		}
		c1, c2 := site.Chapters[0], site.Chapters[1]
		if c1.Num != 1 || c2.Num != 2 {
			t.Errorf("並び順: %d, %d", c1.Num, c2.Num)
		}
		if c1.Title != "はじめて" || !strings.Contains(string(c1.Body), "概要") {
			t.Errorf("_chapter.md が読めていない: %q %q", c1.Title, c1.Body)
		}
		if c1.Next != c2 || c2.Prev != c1 || c1.Prev != nil || c2.Next != nil {
			t.Error("チャプターの前後リンクが不正")
		}
		if got := c1.DoneCount(); got != 1 {
			t.Errorf("DoneCount = %d", got)
		}
	})

	t.Run("ステージ", func(t *testing.T) {
		var ids []string
		for _, s := range site.Stages {
			ids = append(ids, s.ID())
		}
		if want := []string{"1-1", "1-2", "2-1"}; !slices.Equal(ids, want) {
			t.Fatalf("ステージ = %v, want %v", ids, want)
		}
		s11, s12, s21 := site.Stages[0], site.Stages[1], site.Stages[2]
		if s11.Title != "first" {
			t.Errorf("title 省略時はファイル名から: %q", s11.Title)
		}
		if s21.Status != "todo" {
			t.Errorf("status 省略時は todo: %q", s21.Status)
		}
		if s12.LinesText() != "2" || s12.StepsText() != "10" || s21.LinesText() != "-" {
			t.Errorf("行数/ステップ数: %s %s %s", s12.LinesText(), s12.StepsText(), s21.LinesText())
		}
		if s12.Next != s21 || s21.Prev != s12 {
			t.Error("チャプターをまたいだ前後リンクが不正")
		}
		if s11.URL() != "chapters/1/1-1/" {
			t.Errorf("URL = %q", s11.URL())
		}
		if got := s21.EarlierSyntaxes(site); len(got) != 1 || got[0].Chapter.Num != 1 {
			t.Errorf("EarlierSyntaxes = %v", got)
		}
	})

	t.Run("テクニック", func(t *testing.T) {
		var slugs []string
		for _, tc := range site.Techniques {
			slugs = append(slugs, tc.Slug)
		}
		if want := []string{"swap", "swap-2", "mark"}; !slices.Equal(slugs, want) {
			t.Fatalf("テクニック = %v, want %v", slugs, want)
		}
		swap, swap2, mark := site.Techniques[0], site.Techniques[1], site.Techniques[2]
		if swap2.Title != "swap-2" {
			t.Errorf("title 省略時はファイル名から: %q", swap2.Title)
		}
		if len(swap.UsedIn) != 2 || len(mark.UsedIn) != 1 {
			t.Errorf("UsedIn: swap=%d mark=%d", len(swap.UsedIn), len(mark.UsedIn))
		}
		if len(swap.Syntaxes) != 1 || swap.Syntaxes[0].Chapter.Num != 1 {
			t.Error("syntax 省略時は置いたフォルダのチャプター")
		}
		if len(mark.Syntaxes) != 2 {
			t.Errorf("syntax: [1, 2] が反映されていない: %d", len(mark.Syntaxes))
		}
		c1 := site.Chapters[0]
		if !slices.Equal(c1.Techniques, []*Technique{swap, swap2}) {
			t.Error("Chapter.Techniques が不正")
		}
		if got := c1.UsedTechniques(); !slices.Equal(got, []*Technique{swap, swap2}) {
			t.Errorf("UsedTechniques は重複なし: %d 件", len(got))
		}
	})

	t.Run("文法", func(t *testing.T) {
		if len(site.Syntaxes) != 2 {
			t.Fatalf("文法数 = %d", len(site.Syntaxes))
		}
		sx1 := site.Syntaxes[0]
		if sx1.Title != "基本" || site.Chapters[0].Syntax != sx1 {
			t.Error("チャプター1の文法が紐付いていない")
		}
		// swap, swap-2 (省略時) と mark (明示) が関連テクニックになる
		if len(sx1.Techniques) != 3 {
			t.Errorf("関連テクニック = %d 件", len(sx1.Techniques))
		}
		if sx1.Next != site.Syntaxes[1] {
			t.Error("文法の前後リンクが不正")
		}
	})
}

func TestLoadSiteEmpty(t *testing.T) {
	site, err := loadSite(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(site.Chapters)+len(site.Stages)+len(site.Techniques)+len(site.Syntaxes) != 0 {
		t.Error("空のディレクトリなのにコンテンツがある")
	}
}

func TestLoadSiteErrors(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		body    string
		wantErr string
	}{
		{"存在しないテクニック", "chapters/chapter01/1-3.md", "---\ntechniques: [nope]\n---\n", `"nope"`},
		{"不正な status", "chapters/chapter01/1-3.md", "---\nstatus: clear\n---\n", "status"},
		{"チャプター番号の不一致", "chapters/chapter01/2-5.md", "本文", "一致しません"},
		{"テクニックの重複", "techniques/chapter02/swap.md", "本文", "重複"},
		{"frontmatter の打ち間違い", "chapters/chapter01/1-3.md", "---\nline: 3\n---\n", "line"},
		{"YAML の書式エラー", "chapters/chapter01/1-3.md", "---\ntitle: [\n---\n", "frontmatter"},
		{"存在しない文法", "techniques/chapter01/x.md", "---\nsyntax: [9]\n---\n", "チャプター9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadSite(writeContent(t, withFile(fixture, tt.file, tt.body)))
			if err == nil {
				t.Fatal("エラーにならなかった")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("エラー %q に %q が含まれない", err, tt.wantErr)
			}
		})
	}
}

func TestDict(t *testing.T) {
	dict := funcs["dict"].(func(...any) (map[string]any, error))
	m, err := dict("a", 1, "b", "x")
	if err != nil || m["a"] != 1 || m["b"] != "x" {
		t.Errorf("dict = %v, %v", m, err)
	}
	if _, err := dict("a"); err == nil {
		t.Error("奇数個の引数がエラーにならない")
	}
	if _, err := dict(1, 2); err == nil {
		t.Error("文字列以外のキーがエラーにならない")
	}
}
