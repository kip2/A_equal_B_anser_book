---
title: サンプルステージ
lines: 1
steps: 3
status: done            # done / wip / todo
techniques: [sort-by-swap]
---

## 問題

（サンプル）`a` と `b` からなる文字列を、`a` が前に来るように並べ替える。

| 入力 | 出力 |
|---|---|
| `ba` | `ab` |
| `bab` | `abb` |

## 解答

```ab
ba=ab
```

## 考え方

隣り合う `ba` を `ab` に入れ替え続ければ、最後は `a` がすべて前に集まる。
