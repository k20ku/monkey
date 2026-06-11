# Monkey Interpreter in Go

Coding along with "Writing An Interpreter In Go"(Thorsten Ball).

Thorsten Ball 著『Writing An Interpreter In Go』の実装リポジトリ。

## Run REPL

- **Start**:

```bash
go run .
```

- **REPL**:

```
Hello <username>! This is the Monkey programming language!
Feel free to type in commands
Ctrl+D to exit!
'/lex', '/parse', '/eval', '/'
monkey> 
```

- **There are some REPL modes**:
        - `/`: default mode (current is `EVAL` mode) (input REPL to `/`)
        - `EVAL`: evaluates inputs (input REPL to `/eval`)
        - `PARSE`: parse input and prints JSON AST (input REPL to `/parse`)
        - `LEX`: lex input and prints Tokens (input REPL to `/lex`)

## Tests

```bash
go test -timeout 5s -run TestXX github.com/k20ku/monkey/<package> -v
```

Example:

```bash
go test -timeout 5s -run TestParsingPrefixExpressions github.com/k20ku/monkey/parser -v
```

## Motivation

当初は書籍の内容を忠実に実装することを目的としていたが、実装を進める中で REPL や Parser の改善、開発体験の向上にも取り組んでいる。

## Current Status

- 現在は Evaluator を実装中。
- Function Evaluation が残っており、その完了後に REPL の改善を本体へ統合する予定。

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

### 2. ContLine (Future)

将来的な継続行対応．

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
- ただし、これはまだ実験用リポジトリでの検証段階であり、本体への統合は Evaluator 完了後に行う予定。

```commonlisp
(repl) > (defun fib (n)
........   (if (<= n 1)
........     n
........     (+ fib (- n 1)
........       fib (- n 2)
........     )
........   )
........ )
```

```monkey
(repl) > fn fib(n) {
........   return if (n <= 1) {
........     n
........   } else {
........     fib(n - 1) + fib(n - 2)
........   }
........ }
```

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

## Development Process

- REPL 関連の実験は本体リポジトリではなく別リポジトリ [gotorepl](https://github.com/k20ku/gotorepl) で行っている。
- これは大学時代のプロジェクトで、大きな試験機能を直接本体へ投入した結果保守が困難になった経験による。

そのため現在は、

1. 実験環境で検証
2. 問題点を把握
3. 本体へ統合

という流れを取っている。
