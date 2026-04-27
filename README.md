# Monkey lang

## test

```bash
go test -timeout 5s -run TestXxx github.com/k20ku/monkey/<package> -v
```

e.g.

```bash
go test -timeout 5s \
   -run TestParsingPrefixExpressions github.com/k20ku/monkey/parser -v
```

## syntax

```pesudo
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
```
