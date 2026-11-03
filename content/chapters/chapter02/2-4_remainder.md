---
title: 余り
lines: 6
steps:
status: done
techniques: [normalize-chars, remove-in-groups, specific-first, empty-left]
---

## 問題

長さを 3 で割った余り（0〜2）を答える。

| 入力 | 出力 |
| --- | --- |
| `aaa` | `0` |
| `bbbbb` | `2` |

## 解答

行数チャレンジ: 目標 6 行（達成）

### 解答1（6 行）

```ab
a=b
c=b

bbb=
bb=(return)2
b=(return)1
=(return)0
```

## 考え方

なんでもいいが、文字を揃えてしまって3で割る。(今回は `b` にそろえた)

3で割るには3文字ずつ消してしまう。

その後は残った文字数で判定すればよい。