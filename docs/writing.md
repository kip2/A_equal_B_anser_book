# 記事の追加手順

## 早見表

| 追加するもの | 置き場所 | 雛形 |
|---|---|---|
| ステージ | `content/chapters/chapterNN/N-M_<name>.md` | [templates/stage.md](templates/stage.md) |
| テクニック | `content/techniques/chapterNN/<slug>.md` | [templates/technique.md](templates/technique.md) |
| 文法 | `content/syntax/chapterNN.md` | [templates/syntax.md](templates/syntax.md) |
| チャプター名・概要 | `content/chapters/chapterNN/_chapter.md` | [templates/chapter.md](templates/chapter.md) |

- `NN` はチャプター番号を2桁にしたもの（`chapter01`, `chapter12`）。ファイル名の `N-M` は桁をそろえなくてよい（`1-1`, `12-3`）。
- ページ間のリンク（ステージ ⇔ テクニック ⇔ 文法）、一覧、進捗表はすべて自動で作られる。

## 手順

### 1. ファイルを作る

雛形をコピーして、置き場所に配置する。

```sh
cp docs/templates/stage.md content/chapters/chapter01/1-2_<name>.md
```

### 2. 書く

#### ステージ

| 項目 | 内容 |
|---|---|
| `title` | ステージ名。省略するとファイル名の `<name>` 部分 |
| `lines` / `steps` | 行数・ステップ数。空欄なら「-」と表示 |
| `status` | `done`（クリア）/ `wip`（挑戦中）/ `todo`（未着手） |
| `techniques` | 使ったテクニックのファイル名（拡張子なし）のリスト |

- 解答のコードブロックは ` ```ab ` で書く。サイトでは行番号とコピーボタンが付く。
- 行数・ステップ数違いなど解答が複数あるときは、`## 解答` の下に `### 行数優先` のように見出しを分けて並べる。

#### テクニック

- 置いたフォルダのチャプター＝「そのテクニックが使えるようになるチャプター」。
- ファイル名（`<slug>`）が URL とステージからの参照名になる。英小文字とハイフンで付ける（例: `sort-by-swap`）。あとから変えると参照しているステージも直す必要がある。
- 複数チャプターの文法を使う場合は `syntax: [1, 3]` のように書く。

#### 文法

- 1チャプターにつき1ファイル。そのチャプターで**新しく追加された**文法だけを書く。

### 3. 確認する

```sh
go run ./builder -serve     # http://localhost:8080/ で表示を確認
go test ./...               # 書き間違い・リンク切れのチェック
```

存在しないテクニック名、`status` の打ち間違い、frontmatter の項目名の打ち間違いなどは、ビルドかテストでファイル名付きのエラーになる。

### 4. 公開する

```sh
git add content
git commit -m "Add 1-2"
git push
```

`main` に push すると、テストが通ったあとサイトが自動で公開される。

## Claude に頼む場合

Claude Code で `/add-article` を実行するか、「1-2 の解答を追加して」のように頼む。
解答コード・行数・ステップ数・考え方など、分かっている情報を渡すと、ファイルの作成から `go test` での確認まで行う。
（スキルの定義: [.claude/skills/add-article/SKILL.md](../.claude/skills/add-article/SKILL.md)）
