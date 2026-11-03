package main

import (
	"os"
	"path/filepath"
	"testing"
)

// fixture はテスト用のコンテンツ。ファイル名の並びと番号の並びをわざとずらしている。
var fixture = map[string]string{
	"chapters/chapter01/_chapter.md": "---\ntitle: はじめて\n---\n概要\n",
	"chapters/chapter01/1-2_second.md": `---
title: 二問目
lines: 2
steps: 10
status: wip
techniques: [swap, swap-2]
---
## 解答
` + "```ab\nba=ab\n```\n",
	"chapters/chapter01/1-1_first.md": `---
lines: 1
status: done
techniques: [swap]
---
本文
`,
	"chapters/chapter02/2-1.md": `---
title: 次の章
techniques: [mark]
---
`,
	"techniques/chapter01/swap.md":   "---\ntitle: 入れ替え\n---\n本文\n",
	"techniques/chapter01/swap-2.md": "本文だけ\n",
	"techniques/chapter02/mark.md":   "---\ntitle: 目印\nsyntax: [1, 2]\n---\n",
	"syntax/chapter01.md":            "---\ntitle: 基本\n---\n## a=b\n",
	"syntax/chapter02.md":            "## 追加\n",
}

// writeContent は files を一時ディレクトリに書き出し、そのパスを返す。
func writeContent(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func loadFixture(t *testing.T) *Site {
	t.Helper()
	site, err := loadSite(writeContent(t, fixture))
	if err != nil {
		t.Fatalf("loadSite: %v", err)
	}
	return site
}

func withFile(files map[string]string, name, body string) map[string]string {
	m := make(map[string]string, len(files)+1)
	for k, v := range files {
		m[k] = v
	}
	m[name] = body
	return m
}
