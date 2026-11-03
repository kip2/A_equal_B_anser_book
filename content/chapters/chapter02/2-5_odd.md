---
title: 奇数
lines: 10
steps:
status: done
techniques: [sort-by-swap, remove-in-groups, specific-first, empty-left]
---

## 問題

a・b・c それぞれの個数が奇数か 0 なら true、そうでなければ false。

| 入力 | 出力 |
| --- | --- |
| `abc` | `true` |
| `abac` | `false` |

## 解答

行数チャレンジ: 目標 10 行（達成）

### 解答1（10 行）

```ab
# ソート
ca=ac
cb=bc
ba=ab

# 3文字以上を変換
aaa=a
bbb=b
ccc=c

# 残った個数で判定する
aa=(return)false
bb=(return)false
cc=(return)false
=(return)true
```

## 考え方

まずソートを行う。

次に、3文字以上の文字を変換する。

3文字以上のものがあれば何度も適用されるため、実質、3の余りが出現する。

あとは個数で判定すればよい。