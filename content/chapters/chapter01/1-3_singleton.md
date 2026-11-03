---
title: 重複削除
lines: 3
steps:
status: done
techniques: [shrink-runs]
---

## 問題

同じ文字が続いている部分を 1 文字にまとめる。

| 入力 | 出力 |
| --- | --- |
| `aabbbcc` | `abc` |
| `abccba` | `abcba` |

## 解答

行数チャレンジ: 目標 3 行（達成）

### 解答1（3 行）

```ab
aa=a
bb=b
cc=c
```

## 考え方

悩むところなし。