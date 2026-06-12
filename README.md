# Monkey Interpreter in Go

Coding along with "Writing An Interpreter In Go"(Thorsten Ball).

Thorsten Ball 著『Writing An Interpreter In Go』の実装リポジトリ。

## Getting Started to REPL

- **Start**:

```bash
make start
```

- **REPL**:

```
Hello <username>! This is the Monkey programming language!
Feel free to type in commands
Ctrl+D to exit!
> let fib = fn(n, a0, a1) {
.   if (n > 0) {
.     fib(n - 1, a1, a0 + a1)
.   } else {
.     a0
.   }
. }
fn(n, a0, a1) {
  if ((n > 0)) { fib((n - 1), a1, (a0 + a1)) } else { a0 }
}
> let a = 150;
150
> fib(a, 0, 1)
6792540214324356296
> fib(a)
ERROR: invalid function call: parameters=3. args=1
```

## Tests

```bash
go test -timeout 5s -run TestXX github.com/k20ku/monkey/<package> -v
```

Example:

```bash
go test -timeout 5s -run TestParsingPrefixExpressions github.com/k20ku/monkey/parser -v
```

## Monkey Language AST Structure

```ocaml
Node :=
      *Program(Statements: []Statement)
    | Statement
    | Expression

Statement :=
      *LetStatement(Name: *identifier, Value: Expression)
    | *ReturnStatement(ReturnValue: Expression)
    | *ExpressionStatement(Expression: Expression)
    | *BlockStatement(Statements: []Statement)

Expression :=
      *Identifier
    | *IntegerLiteral
    | *PrefixExpression(Right: Expression)
    | *InfixExpressoion(Left: Expression, Right: Expression)
    | *Boolean
    | *IfExpression(
        Condition: Expression
        Consequence: *BlockStatement
        Alternative: *BlockStatement
      )
    | *FunctionLiteral(
        Parameters: []*identifier
        Body: *BlockStatement  
      )
    | *CallExpression(
        Function: Expression
        Arguments: []Expression
      )
```

## Motivation

当初は書籍の内容を忠実に実装することを目的としていたが、実装を進める中で REPL や Parser の改善、開発体験の向上にも取り組んでいる。

## Current Status

- 現在は Evaluator を実装完了。
- REPLの複数行対応完了（Unix系・Windows系）

## Improvements

### 1. REPL

書籍版 REPL は `scanner` ベースで実装されている。

シンプルで分かりやすい一方、

- 矢印キーがエスケープシーケンスとして表示される
- 入力履歴が扱いづらい
- 実用的な REPL としては不便

という問題があった。

そのため readline 系ライブラリを調査し、Go 製 REPL である [Gore](https://github.com/x-motemen/gore) が利用している [liner](https://github.com/peterh/liner) を採用した。

**採用理由**:

- 現在も保守されているわけではないが
- 利用実績がある
- ほぼnative Goで書かれており依存が非常に少ない
- 必要以上に高機能ではない

### 2. ContLine

継続行対応．

```monkey
monkey> fn add(x0, x1) {
......    return x0 + x1;
......  }
```

#### 経緯

liner 導入後、継続行機能の実装を試みた。

- 当初は Gore の [contLiner](https://github.com/x-motemen/gore/blob/main/liner.go) を移植すれば動くと考えていたが、実際には無限ループなどが発生しうまく動作しなかった。
- 調査を進めた結果、contLiner は単独で成立しているわけではなく、`text/scanner` を利用した `lexer` と密接に結合していることが判明した。

そこで [実験用リポジトリ](https://github.com/k20ku/gotorepl) を作成し、

```
Input
↓
text/scanner
↓
Lexer
↓
contLiner
↓
Gore
↓
EntryPoint
```

という構成を再現した。

#### 展望

- 現在は複数行入力と自動インデントが動作することを確認している。
- 今後は `session.go` を用いて，`:mode lex` などによるモード切替を実装する予定．

### 3. Parser

#### 経緯

- 継続行の実装中、不完全な入力を大量に扱うようになった。
- その結果、[`parseBlockStatement`](https://github.com/k20ku/monkey/blob/fd971177dc7b0f85ce5d850a1d187d2d6fb1bd8f/parser/parser.go#L385) 周辺の境界条件に気づいた。

例えば、

```monkey
fn(x) {
return x
```

のような入力でも `EOF` 到達によってブロック解析が終了し、場合によっては AST が生成されてしまう。これは無限ループ防止のための `EOF` 判定による副作用である。

#### 現在

- 呼び出し側で閉じ波括弧の存在を確認し、適切なエラーを出すようにしている。
- 既存テストへの影響がないことも確認済み。
- 本来は `parseBlockStatement` 自体が責任を持つべきかもしれないが、言語拡張途中であるため変更範囲を最小限に留めている。

### 4. AST

- Evaluator 実装中、AST 全体の構造を把握しづらくなった。
- そのため OCaml の代数的データ型風の表記で AST を整理した。
- これは見栄えのためではなく、Evaluator 実装時の認知負荷を下げることが目的である。

## Development Process

- REPL 関連の実験は本体リポジトリではなく別リポジトリ [gotorepl](https://github.com/k20ku/gotorepl) で行っている。
- これは大学時代のプロジェクトで、大きな試験機能を直接本体へ投入した結果保守が困難になった経験による。

そのため現在は、

1. 実験環境で検証
2. 問題点を把握
3. 本体へ統合

という流れを取っている。
