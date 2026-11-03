# A=B Answer Book

ゲーム「A=B」の自分の解答集とテクニック集です。

> [!WARNING]
> ネタバレ注意：全ステージの解答と解説を含みます。

🔗 サイト：<https://kip2.github.io/A_equal_B_anser_book/>

## 構成

```text
content/
├── chapters/chapterNN/
│   ├── _chapter.md          # チャプター名・概要（省略可）
│   └── N-M_<name>.md        # ステージの解答と解説
├── techniques/chapterNN/
│   └── <slug>.md            # テクニック（置いたフォルダ ＝ 使えるようになるチャプター）
└── syntax/
    └── chapterNN.md         # チャプターで追加される文法
builder/                     # Go製のサイト生成ツール
docs/                        # 記事の追加手順と雛形
```

## 記事の追加

手順と書式は [docs/writing.md](docs/writing.md) を参照してください。雛形は [docs/templates/](docs/templates/) にあります。
Claude Code では `/add-article` で追加を頼めます。

ページ同士のリンク（ステージ ⇔ テクニック ⇔ 文法）や進捗表は、frontmatterから自動で生成されます。

## ローカルで確認

```sh
go run ./builder -serve      # http://localhost:8080/
```

`go run ./builder` だけなら `dist/` に出力して終了します。

## テスト

```sh
go test ./...                      # すべて
go test -skip '^TestBuild' ./...   # ユニットテストのみ
go test -run '^TestBuild' ./...    # 結合テスト（サイト生成とリンク切れチェック）のみ
```

結合テストは実際の `content/` も生成して検査するので、存在しないテクニック名や frontmatter の書き間違い、リンク切れもここで検出できます。

## 公開

`main` に push すると GitHub Actions がビルドして GitHub Pages に公開します。
初回のみ、リポジトリの Settings → Pages → Source を「GitHub Actions」にしてください。
