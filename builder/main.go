// builder は content/ の Markdown から静的サイトを生成する。
//
//	go run ./builder           # dist/ に出力
//	go run ./builder -serve    # 出力したあと http://localhost:8080 で確認
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	content := flag.String("content", "content", "コンテンツのディレクトリ")
	out := flag.String("out", "dist", "出力先のディレクトリ")
	serve := flag.Bool("serve", false, "生成後にローカルサーバーを起動する")
	addr := flag.String("addr", "localhost:8080", "-serve のアドレス")
	baseURL := flag.String("base-url", "", "公開先の URL（OGP タグの絶対 URL に使う。空なら OGP タグを出さない）")
	repoURL := flag.String("repo-url", "", "リポジトリの URL（フッターのリンクに使う。空ならリンクを出さない）")
	flag.Parse()

	site, err := loadSite(*content)
	if err != nil {
		log.Fatal(err)
	}
	if err := build(site, *out, buildOptions{BaseURL: *baseURL, RepoURL: *repoURL}); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s に出力しました（チャプター %d / ステージ %d / テクニック %d / 文法 %d）",
		*out, len(site.Chapters), len(site.Stages), len(site.Techniques), len(site.Syntaxes))

	if *serve {
		log.Printf("http://%s/ で確認できます（Ctrl+C で終了）", *addr)
		log.Fatal(http.ListenAndServe(*addr, http.FileServer(http.Dir(*out))))
	}
}
