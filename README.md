# A=B Answer Book

プログラミングパズルゲーム [A=B](https://store.steampowered.com/app/1720850/AB/)（Steam）の、自分の解答集とテクニック集です。

> [!WARNING]
> ネタバレ注意：全ステージの解答と解説を含みます。

🔗 サイト：<https://kip2.github.io/A_equal_B_anser_book/>

## サイトの見方

画面上部のタイトルとメニューから、4つのページに移動できます。

| 画面の表示 | ページ | 内容 |
| --- | --- | --- |
| `A=B Answer Book` | トップ | 全体の目次と進み具合 |
| `チャプター` | チャプター | チャプターごとのステージ一覧 |
| `テクニック` | テクニック総覧 | 解き方のテクニック集 |
| `文法` | 文法総覧 | チャプターごとに追加される文法 |

### トップ

チャプター・テクニック・文法の目次です。クリアしたステージ数もここで分かります。

### チャプター

チャプターの一覧です。チャプター名を選ぶと、そのチャプターのページに移動します。

- **チャプターのページ**：ステージ一覧、そのチャプターで使えるようになったテクニック、追加された文法
- **ステージのページ**：問題の要約、解答コード、考え方。下に、使ったテクニックと文法へのリンクがあります

### テクニック

テクニックを「使えるようになるチャプター」ごとに並べています。

見出しを開くと解説が読めます。個別ページには、そのテクニックを使ったステージと関連する文法へのリンクがあります。

### 文法

チャプターごとに追加される文法をまとめています。

個別ページには、そのチャプターのステージと関連するテクニックへのリンクがあります。

### 共通の操作

- `[+]` / `[-]` の付いた見出しは、クリックすると開閉します。
- `[すべて開く]` / `[すべて閉じる]` で、ページ内の見出しをまとめて開閉できます。
- 解答コードの右上のコピーアイコン（重なった2枚の紙）で、コードをコピーできます。
- 上部のパンくず（`トップ / チャプター / チャプター1 / 1-1`）から上の階層に戻れます。
- 下部の `<<` / `>>` で、前後のページに移動できます。

### 表記

| 表記 | 意味 |
| --- | --- |
| `3行 / 12ステップ` | 解答の行数とステップ数。`-` は未記録 |
| `クリア` / `挑戦中` / `未着手` | ステージの状態 |
| `チャプターN〜` | そのテクニックがチャプター N の文法から使えること |

## 記事の追加

手順・書式・雛形は [docs/writing.md](docs/writing.md) にまとめています。

Claude Code では `/add-article` で追加を頼めます。

## 自分の解答集を作る場合

このリポジトリの `template` ブランチが、サンプル記事だけが入った仕組み一式です。

1. `template` ブランチを取得し、自分のリポジトリの `main` として push します。

   ```sh
   git clone --branch template --single-branch https://github.com/kip2/A_equal_B_anser_book.git my-answer-book
   cd my-answer-book
   git switch -c main
   git remote set-url origin <自分のリポジトリの URL>
   git push -u origin main
   ```

2. 自分のリポジトリの Settings → Pages → Source を「GitHub Actions」にします。

3. `content/` のサンプル記事を消して、自分の記事を書きます（[docs/writing.md](docs/writing.md)）。

4. README のサイト URL などを自分のものに書き換えます。

`main` に push すると、テストが通ったあと GitHub Pages に公開されます。

### 参考

- [docs/writing.md](docs/writing.md)：記事の追加手順と書式
- [docs/development.md](docs/development.md)：リポジトリ構成、テスト、CI、公開の仕組み
- [.github/workflows/ci.yml](.github/workflows/ci.yml)：CI の設定
- [.github/workflows/pages.yml](.github/workflows/pages.yml)：GitHub Pages への公開の設定
