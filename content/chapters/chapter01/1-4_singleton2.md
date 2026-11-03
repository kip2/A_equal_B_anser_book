---
title: 重複削除 2
lines: 2
steps:
status: done
techniques: [shrink-runs]
---

## 問題

2 文字以上続いている a をまとめて消す。

| 入力 | 出力 |
| --- | --- |
| `aaabcaa` | `bc` |
| `baaba` | `bba` |

## 解答

行数チャレンジ: 目標 2 行（達成）

### 解答1（2 行）

```ab
aaa=aa
aa=
```

## 考え方

「連続した」なので、a単体を消すことでは対応できない。

偶数を`aa=`で対応し、奇数を偶数に変換することで消去。