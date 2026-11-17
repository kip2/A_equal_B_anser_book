---
title: 昇順
lines: 9
steps:
status: done
techniques: [sort-by-swap, cancel-pairs, specific-first, empty-left]
---

## 問題

個数が a < b < c になっていれば true、そうでなければ false。

| 入力 | 出力 |
| --- | --- |
| `abbccc` | `true` |
| `ccccbbb` | `true` |

## 解答


### 解答1（10 行）

```ab
# バブルソート
ca=ac
cb=bc
ba=ab

# abcを一組ずつ消す
xc=
ab=x
xb=bx

# 残ったbcの形を一組ずつ消す
byc=y
bc=y

# bccの形であればtrue
yc=(return)true
# それ以外をfalse
=(return)false

```

### 解答2（9 行）

```ab
# バブルソート
ca=ac
cb=bc
ba=ab

# abcを一組ずつ消す
xc=
ab=x
xb=bx

# 残ったbcの形を一組ずつ消す
bbc=b

# bccの形であればtrue
bcc=(return)true
# それ以外をfalse
=(return)false

```

## 考え方

まずソートして、a・b・c の順に並べる。

次に、a・b・c を1つずつ組にして消していく。a と b を `x` にまとめ、`x` を b の右へ移動させて、c と一緒に消す。

a がなくなったら、残った b と c を1つずつ組にして消す。

最後に c が残っていれば、b より c が多いので true。それ以外はすべて false。

### 解答2 の改良点

解答1 では、b と c の組を `y` に置き換えてから消していた。解答2 では `bbc=b` で、b を1つ残したまま b と c の組を消す。

組を消し終わると `b` のあとに c が並んだ形になるので、`bcc` があれば c のほうが多い。b を目印として使い回すので、`y` が要らなくなり、1行減る。

