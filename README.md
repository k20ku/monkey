# Monkey lang

## test

```bash
go test -timeout 5s -run TestXX monkey/xx -v
```

e.g.

```bash
go test -timeout 5s -run TestParsingPrefixExpressions monkey/parser -v
```

## syntax

```pesudo
Statement :=
    | *LetStatement(Name: *identifier, Value: Expression)
    | *ReturnStatement(ReturnValue: Expression)
    | *ExpressionStatement(Expression: Expression)
    | *BlockStatement(Statements: []Statement)

Expression :=
    | *Identifier
    | *IntegerLiteral
    | *PrefixExpression(Right: Expression)
    | *InfixExpressoion(Left: Expression, Right: Expression)
    | *Boolean
    | *IfExpression(
    |   Condition: Expression
    |   Consequence: *BlockStatement
    |   Alternative: *BlockStatement
    | )
```
