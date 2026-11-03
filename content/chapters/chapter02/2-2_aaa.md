---
title: AAA
lines: 4
steps:
status: done
techniques: [normalize-chars, empty-left]
---

## 問題

a が 3 個以上あれば true、なければ false。

| 入力 | 出力 |
| --- | --- |
| `aaa` | `true` |
| `aaaa` | `true` |

## 解答

行数チャレンジ: 目標 4 行（達成）

### 解答1（4 行）

```ab
b=
c=
aaa=(return)true
=(return)false
```

## 考え方

単純に `b` と `c` を消せばよい。

