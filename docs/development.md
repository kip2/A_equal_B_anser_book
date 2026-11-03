# 開発者向けドキュメント

サイト生成ツールの構成、テスト、CI についてまとめています。記事の書き方は [writing.md](writing.md) を参照してください。

## リポジトリ構成

```text
content/                     # 記事
├── chapters/chapterNN/      #   チャプター概要（_chapter.md）とステージ（N-M_<name>.md）
├── techniques/chapterNN/    #   テクニック（置いたフォルダ ＝ 使えるようになるチャプター）
└── syntax/chapterNN.md      #   チャプターで追加される文法
builder/                     # Go 製のサイト生成ツール
├── main.go                  #   コマンドの入口（-content / -out / -serve / -addr）
├── content.go               #   Markdown・frontmatter の読み込みと、ページ同士の関連付け
├── render.go                #   テンプレートで HTML を出力
├── templates/               #   各ページのテンプレート
├── static/                  #   CSS と JavaScript（開閉・コピーボタン）
└── *_test.go                #   テスト
docs/                        # 記事の追加手順（writing.md）、雛形（templates/）、このドキュメント
.claude/skills/              # Claude Code 用の記事追加スキル
.github/                     # CI・公開・Dependabot の設定
dist/                        # 生成されたサイト（git 管理外）
```

## サイトの生成

```sh
go run ./builder             # dist/ に出力
go run ./builder -serve      # 出力後に http://localhost:8080/ で確認
```

## テスト

```sh
go test ./...                      # すべて
go test -skip '^TestBuild' ./...   # ユニットテストのみ
go test -run '^TestBuild' ./...    # 結合テストのみ
```

### ユニットテスト（[builder/content_test.go](../builder/content_test.go)）

| テスト | 確認していること |
| --- | --- |
| `TestSplitFrontMatter` | frontmatter と本文の切り分け。frontmatter なし・CRLF・BOM 付き・空・閉じていない場合も正しく扱えるか |
| `TestIntText` / `TestStatusLabel` / `TestDisplayTitle` | 行数の `-` 表示、状態の表示名、チャプター・文法のタイトル表記 |
| `TestLoadSite` | コンテンツの読み込みと関連付け。チャプター・ステージの並び順、チャプターをまたぐ前後リンク、`title` / `status` の省略時の値、ステージ ⇔ テクニック ⇔ 文法の相互リンク、使ったテクニックの重複除去 |
| `TestLoadSiteEmpty` | `content/` が空でもエラーにならないか |
| `TestLoadSiteErrors` | 書き間違いをエラーにできるか。存在しないテクニック名、不正な `status`、ファイル名とフォルダのチャプター番号の不一致、テクニック名の重複、frontmatter の項目名の打ち間違い、YAML の書式エラー、存在しない文法チャプターの指定 |
| `TestDict` | テンプレート用の補助関数の引数チェック |

### 結合テスト（[builder/build_test.go](../builder/build_test.go)）

| テスト | 確認していること |
| --- | --- |
| `TestBuildFixture` | テスト用コンテンツからサイトを生成し、全ページが出力されるか、ステージ・テクニックのページに解答コード・関連リンク・前後リンク・記録が出ているか |
| `TestBuildRemovesStaleFiles` | 再生成時に前回の出力が残らないか |
| `TestBuildContent` | **実際の `content/` からサイトを生成できるか**。記事の書き間違いはここで検出される |

結合テストでは、生成した全ページについて次も検査します。

- 内部リンク（`href` / `src`）の参照先ファイルがすべて存在するか（リンク切れ）
- `#アンカー` 付きリンクの参照先に、その `id` があるか
- サイトの外を指す相対リンクがないか
- テンプレートの不具合で出る文字列（`ZgotmplZ`、`<no value>`）が含まれていないか

テスト用のコンテンツは [builder/helpers_test.go](../builder/helpers_test.go) にあります。

## CI（[.github/workflows/ci.yml](../.github/workflows/ci.yml)）

`main` への push とプルリクエストで実行されます。

| ジョブ | 内容 |
| --- | --- |
| 静的チェック | `gofmt`（整形）、`go mod tidy`（依存の整理漏れ）、`go vet`、`staticcheck` |
| テスト | ユニットテストと結合テスト（`-race` 付き） |
| 脆弱性チェック | `govulncheck` で依存ライブラリの既知の脆弱性を検査 |

## 公開（[.github/workflows/pages.yml](../.github/workflows/pages.yml)）

`main` への push で実行されます。テストを先に実行し、通った場合だけサイトを生成して GitHub Pages に公開します。

初回のみ、リポジトリの Settings → Pages → Source を「GitHub Actions」にしてください。

## 依存関係の更新（[.github/dependabot.yml](../.github/dependabot.yml)）

Dependabot が Go モジュールと GitHub Actions の更新を月に1回プルリクエストで提案します。CI が通っていればマージして構いません。
