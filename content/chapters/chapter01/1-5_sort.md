---
title: ソート
lines: 3
steps:
status: done
techniques: [sort-by-swap]
---

## 問題

文字列を a → b → c の順に並べ替える。

| 入力 | 出力 |
| --- | --- |
| `caba` | `aabc` |
| `abccba` | `aabbcc` |

## 解答

行数チャレンジ: 目標 3 行（達成）

### 解答1（3 行）

```ab
ca=ac
cb=bc
ba=ab
```

## 考え方

入れ替えの基本。

後ろになっている文字が前に来ているパターンを入れ替えることで対応できる。


cは2パターンあり、bは1パターン、aは0パターンのため、解答の3行となる。