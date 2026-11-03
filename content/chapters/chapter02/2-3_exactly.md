---
title: ちょうど3文字
lines: 5
steps:
status: done
techniques: [normalize-chars, other-chars]
---

## 問題

長さがちょうど 3 なら true、それ以外は false。

| 入力 | 出力 |
| --- | --- |
| `abc` | `true` |
| `ab` | `false` |

## 解答

行数チャレンジ: 目標 5 行（達成）

### 解答1（5 行）

```ab
b=a
c=a

aaa=true
a=(return)false
truetrue=(return)false
```

## 考え方

すべて `a` に変換する。

`aaa` のパターンを `true` とする。

しかし、この状態で `return` すると `aaaa` などのパターンも `true` となってしまう。

そのため、一度 `true` に変換しておいて、 `a` が残るパターンを `false` とする。

しかし、6文字のパターンでは `truetrue` となるため、その点だけケアする。

もっときれいに解ける方法ありそう。
