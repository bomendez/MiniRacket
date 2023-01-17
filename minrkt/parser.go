package minrkt

import (
	"fmt"
	"strconv"
)

type Environment struct {
	Variables map[string]interface{}
	Functions map[string]FuncParamExpr
	CallStack []map[string]interface{}
}

type ArgumentError struct {
	c string
}

func (e *ArgumentError) Error() string {
	return fmt.Sprintf("Too few arguments after operator: %s", e.c)
}

type ParseError struct {
	c string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("with %s", e.c)
}

type EvalError struct {
	c string
}

func (e *EvalError) Error() string {
	return fmt.Sprintf("%s", e.c)
}

type Exp interface {
	Eval(*Environment) (interface{}, error)
}

type expVar struct {
	name string
}

type expFunc struct {
	name      string
	arguments []Exp
}

type expBoolConst struct {
	val bool
}

type expNumConst struct {
	val float64
}

type expOperator struct {
	opType   TokenType
	operands []Exp
}

type expDefineVar struct {
	name string
	val  Exp
}

type expDefineFunc struct {
	name       string
	expression Exp
	paramNames []string
}

type FuncParamExpr struct {
	params     []string
	expression Exp
}

func (e *expVar) Eval(env *Environment) (interface{}, error) {
	var funcVarMap map[string]interface{}
	if len(env.CallStack) > 0 {
		funcVarMap = env.CallStack[len(env.CallStack)-1]
		// } else {
		// 	funcVarMap = make(map[string]interface{})
	}
	funcVar, ok1 := funcVarMap[e.name]
	if !ok1 { // check global variables
		varVal, ok := env.Variables[e.name]
		if !ok {
			return e.name, &EvalError{e.name + " undefined"}
		} else {
			return varVal, nil
		}
	} else {
		return funcVar, nil
	}
}

func (e *expFunc) Eval(env *Environment) (interface{}, error) {
	funcStruct, ok := env.Functions[e.name]
	if !ok {
		return e.name, &EvalError{e.name + " undefined"}
	}
	funcExpression := funcStruct.expression
	funcParams := funcStruct.params
	localParams := make(map[string]interface{})
	// populate map of argument to parameter assignment
	for i, param := range funcParams {
		arg, err := e.arguments[i].Eval(env)
		if err == nil {
			localParams[param] = arg
		}
	}

	// push localParams to env
	env.CallStack = append(env.CallStack, localParams)
	// evaluate Exp
	result, err := funcExpression.Eval(env)
	// pop localParams from env upon return
	env.CallStack = env.CallStack[:len(env.CallStack)-1]

	return result, err
}

func (e *expBoolConst) Eval(_ *Environment) (interface{}, error) {
	var val interface{} = e.val
	return val, nil
}

func (e *expNumConst) Eval(_ *Environment) (interface{}, error) {
	var val interface{} = e.val
	return val, nil
}

func (e *expDefineVar) Eval(env *Environment) (interface{}, error) {
	iName := e.name
	iValue, err := e.val.Eval(env)
	env.Variables[iName] = iValue
	return nil, err
}

func (e *expDefineFunc) Eval(env *Environment) (interface{}, error) {
	varName := e.name
	// create struct for map value {Exp, []string}
	mapVal := FuncParamExpr{e.paramNames, e.expression}
	env.Functions[varName] = mapVal // struct of param list and Exp
	return nil, nil
}

func (e *expOperator) Eval(env *Environment) (interface{}, error) {
	var result interface{}
	var err error
	switch opType := e.opType; opType {
	case TOK_ADD:
		var sum float64
		for i := 0; i < len(e.operands); i++ {
			var iSum interface{}
			iSum, err = e.operands[i].Eval(env)
			subSum, ok := iSum.(float64)
			if ok {
				sum = sum + subSum
			} else {
				err = &EvalError{"Operands not Converted"}
			}
		}
		result = sum

	case TOK_SUB:
		var i int
		var diff float64
		if len(e.operands) == 0 {
			diff = 0
			err = &EvalError{"subtraction requires at least one operand"}
		} else if len(e.operands) == 1 {
			i = 0
		} else {
			var iDiff interface{}
			iDiff, err = e.operands[0].Eval(env)
			subDiff, ok := iDiff.(float64)
			if ok {
				diff = subDiff
				i = 1
			} else {
				err = &EvalError{"Operands not Converted"}
			}
		}
		for ; i < len(e.operands); i++ {
			var iDiff interface{}
			iDiff, err = e.operands[i].Eval(env)
			subDiff, ok := iDiff.(float64)
			if ok {
				diff = diff - subDiff
			} else {
				err = &EvalError{"Operands not Converted"}
			}
		}
		result = diff
	case TOK_MUL:
		product := 1.0
		for i := 0; i < len(e.operands); i++ {
			var iProduct interface{}
			iProduct, err = e.operands[i].Eval(env)
			subProduct, ok := iProduct.(float64)
			if ok {
				product = product * subProduct
			} else {
				err = &EvalError{"Operands not Converted"}
			}
		}
		result = product
	case TOK_DIV: // assumes non-empty operandList
		var i int
		var quotient float64
		if len(e.operands) == 0 {
			quotient = 0
			err = &EvalError{"division requires at least one operand"}
		} else if len(e.operands) == 1 {
			i = 0
		} else {
			var iQuotient interface{}
			iQuotient, err = e.operands[0].Eval(env)
			subQuotient, ok := iQuotient.(float64)
			if ok {
				quotient = subQuotient
				i = 1
			} else {
				err = &EvalError{"Operands not Converted"}
			}
		}
		for ; i < len(e.operands); i++ {
			var iQuotient interface{}
			iQuotient, err = e.operands[i].Eval(env)
			subQuotient, ok := iQuotient.(float64)
			if ok {
				quotient = quotient / subQuotient
			} else {
				err = &EvalError{"Operands not Converted"}
			}
		}
		result = quotient
	case TOK_EQ:
		var boolResult bool
		if len(e.operands) != 2 {
			err = &EvalError{"= requires 2 operands"}
		} else {
			var iBool1 interface{}
			var iBool2 interface{}
			iBool1, err = e.operands[0].Eval(env)
			iBool2, err = e.operands[1].Eval(env)
			type1 := fmt.Sprintf("%T", iBool1)
			type2 := fmt.Sprintf("%T", iBool2)
			if type1 == type2 {
				boolResult = iBool1 == iBool2
			} else {
				err = &ParseError{"mismatched types"}
			}
		}
		result = boolResult
	case TOK_GTEQ:
		var boolResult bool
		if len(e.operands) != 2 {
			err = &EvalError{"= requires 2 operands"}
		} else {
			var iBool1 interface{}
			var iBool2 interface{}
			iBool1, err = e.operands[0].Eval(env)
			iBool2, err = e.operands[1].Eval(env)
			type1 := fmt.Sprintf("%T", iBool1)
			type2 := fmt.Sprintf("%T", iBool2)
			if type1 == "float64" && type2 == "float64" {
				subBool1, ok1 := iBool1.(float64)
				subBool2, ok2 := iBool2.(float64)
				if ok1 && ok2 {
					boolResult = subBool1 >= subBool2
				}
			} else {
				err = &ParseError{"mismatched types"}
			}
		}
		result = boolResult
	case TOK_LTEQ:
		var boolResult bool
		if len(e.operands) != 2 {
			err = &EvalError{"= requires 2 operands"}
		} else {
			var iBool1 interface{}
			var iBool2 interface{}
			iBool1, err = e.operands[0].Eval(env)
			iBool2, err = e.operands[1].Eval(env)
			type1 := fmt.Sprintf("%T", iBool1)
			type2 := fmt.Sprintf("%T", iBool2)
			if type1 == "float64" && type2 == "float64" {
				subBool1, ok1 := iBool1.(float64)
				subBool2, ok2 := iBool2.(float64)
				if ok1 && ok2 {
					boolResult = subBool1 <= subBool2
				}
			} else {
				err = &ParseError{"mismatched types"}
			}
		}
		result = boolResult
	case TOK_GT:
		var boolResult bool
		if len(e.operands) != 2 {
			err = &EvalError{"= requires 2 operands"}
		} else {
			var iBool1 interface{}
			var iBool2 interface{}
			iBool1, err = e.operands[0].Eval(env)
			iBool2, err = e.operands[1].Eval(env)
			type1 := fmt.Sprintf("%T", iBool1)
			type2 := fmt.Sprintf("%T", iBool2)
			if type1 == "float64" && type2 == "float64" {
				subBool1, ok1 := iBool1.(float64)
				subBool2, ok2 := iBool2.(float64)
				if ok1 && ok2 {
					boolResult = subBool1 > subBool2
				}
			} else {
				err = &ParseError{"mismatched types"}
			}
		}
		result = boolResult
	case TOK_LT:
		var boolResult bool
		if len(e.operands) != 2 {
			err = &EvalError{"= requires 2 operands"}
		} else {
			var iBool1 interface{}
			var iBool2 interface{}
			iBool1, err = e.operands[0].Eval(env)
			iBool2, err = e.operands[1].Eval(env)
			type1 := fmt.Sprintf("%T", iBool1)
			type2 := fmt.Sprintf("%T", iBool2)
			if type1 == "float64" && type2 == "float64" {
				subBool1, ok1 := iBool1.(float64)
				subBool2, ok2 := iBool2.(float64)
				if ok1 && ok2 {
					boolResult = subBool1 < subBool2
				}
			} else {
				err = &ParseError{"mismatched types"}
			}
		}
		result = boolResult
	case TOK_AND:
		var boolResult bool
		if len(e.operands) != 2 {
			err = &EvalError{"'and' requires 2 operands"}
		} else {
			var iBool1 interface{}
			var iBool2 interface{}
			iBool1, err = e.operands[0].Eval(env)
			subBool1, ok1 := iBool1.(bool)
			if !ok1 {
				err = &EvalError{"first 'and' operand must be boolean"}
			} else {
				if !subBool1 {
					boolResult = false
				} else {
					iBool2, err = e.operands[1].Eval(env)
					subBool2, ok2 := iBool2.(bool)
					if !ok2 {
						err = &EvalError{"second 'and' operand must be boolean"}
					} else {
						if !subBool2 {
							boolResult = false
						} else {
							boolResult = true
						}
					}
				}
			}
		}
		result = boolResult
	case TOK_OR:
		var boolResult bool
		if len(e.operands) != 2 {
			err = &EvalError{"'or' requires 2 operands"}
		} else {
			var iBool1 interface{}
			var iBool2 interface{}
			iBool1, err = e.operands[0].Eval(env)
			subBool1, ok1 := iBool1.(bool)
			if !ok1 {
				err = &EvalError{"first 'or' operand must be boolean"}
			} else {
				if !subBool1 {
					iBool2, err = e.operands[1].Eval(env)
					subBool2, ok2 := iBool2.(bool)
					if !ok2 {
						err = &EvalError{"second 'or' operand must be boolean"}
					} else {
						if subBool2 {
							boolResult = true
						} else {
							boolResult = false
						}
					}
				} else {
					boolResult = true
				}
			}
		}
		result = boolResult
	case TOK_NOT:
		var boolResult bool
		if len(e.operands) != 1 {
			err = &EvalError{"'not' requires 1 operand"}
		} else {
			var iBool1 interface{}
			iBool1, err = e.operands[0].Eval(env)
			subBool1, ok1 := iBool1.(bool)
			if !ok1 {
				err = &EvalError{"'not' operand must be boolean"}
			} else {
				boolResult = !subBool1
			}
		}
		result = boolResult
	case TOK_IF:
		var ifResult interface{}
		if len(e.operands) != 3 {
			err = &EvalError{"'if' requires 3 operands"}
		} else {
			var iBool1 interface{}
			iBool1, err = e.operands[0].Eval(env)
			subBool1, ok1 := iBool1.(bool)
			if !ok1 {
				err = &EvalError{"'not' operand must be boolean"}
			} else {
				if subBool1 {
					ifResult, err = e.operands[1].Eval(env)
				} else {
					ifResult, err = e.operands[2].Eval(env)
				}
			}
		}
		result = ifResult
	}
	return result, err
}

func isLeftParenthesis(tok Token) bool {
	if tok.tokType == TOK_LPAREN {
		return true
	} else {
		return false
	}
}

func isIdentifier(tok Token) bool {
	return tok.tokType == TOK_VAR
}

func isOperand(tok Token) bool {
	if tok.tokType == TOK_NUM ||
		tok.tokType == TOK_TRUE ||
		tok.tokType == TOK_FALSE {
		return true
	} else {
		return false
	}
}

func isOperator(tok Token) bool {
	if tok.tokType == TOK_ADD ||
		tok.tokType == TOK_SUB ||
		tok.tokType == TOK_MUL ||
		tok.tokType == TOK_DIV ||
		tok.tokType == TOK_EQ ||
		tok.tokType == TOK_GTEQ ||
		tok.tokType == TOK_LTEQ ||
		tok.tokType == TOK_GT ||
		tok.tokType == TOK_LT ||
		tok.tokType == TOK_AND ||
		tok.tokType == TOK_OR ||
		tok.tokType == TOK_IF {
		return true
	} else {
		return false
	}
}

func isDefine(tok Token) bool {
	return tok.tokType == TOK_DEFINE
}

func buildOperandNode(token Token) Exp {
	currOp := token
	var opNode Exp
	switch tokType := token.tokType; tokType {
	case TOK_NUM:
		value, err := strconv.ParseFloat(currOp.val, 64)

		if err == nil {
			opNode = &expNumConst{value}
		}
	case TOK_TRUE:
		opNode = &expBoolConst{true}
	case TOK_FALSE:
		opNode = &expBoolConst{false}
	}
	return opNode
}

func buildVar(token Token) Exp {
	return &expVar{token.val}
}

// Assumes non-empty input
func Parser(tokens []Token) ([]Token, Exp, error) {
	// if tokens is empty, return empty slice and set error
	if len(tokens) == 0 {
		var exp Exp
		return []Token{}, exp, &ParseError{"incomplete statement"}
	}
	currToken := tokens[0]
	if isOperand(currToken) {
		return tokens[1:], buildOperandNode(currToken), nil
	}
	if isIdentifier(currToken) {
		return tokens[1:], buildVar(currToken), nil
	}
	if len(tokens) == 1 {
		if currToken.tokType != TOK_RPAREN {
			var exp Exp
			return []Token{}, exp, &ParseError{"missing value after ("}
		}
	}

	if isLeftParenthesis(currToken) {
		operatorToken := tokens[1]
		var err error
		if isIdentifier(operatorToken) { // parse function call
			funcName := operatorToken.val
			leftOver := tokens[2:]

			// parse arguments excluding right parenthesis
			var funcArguments []Exp
			if len(leftOver) > 1 {
				for leftOver[0].tokType != TOK_RPAREN {
					if !isOperand(leftOver[0]) {
						var exp Exp
						return []Token{}, exp, &ParseError{"invalid function arguments"}
					}
					var funcExpression Exp
					leftOver, funcExpression, err = Parser(leftOver)
					funcArguments = append(funcArguments, funcExpression)
				}
			}

			if len(leftOver) == 0 || leftOver[0].tokType != TOK_RPAREN {
				var exp Exp
				return []Token{}, exp, &ParseError{"missing closing )"}
			} else {
				return leftOver[1:], &expFunc{funcName, funcArguments}, err
			}
		} else if isDefine(operatorToken) {
			var varName string
			var varExpression Exp
			var varOperands []string
			leftOver := tokens[2:]
			if len(leftOver) < 2 {
				var exp Exp
				return []Token{}, exp, &ParseError{"define requires two inputs"}
			}
			if isLeftParenthesis(leftOver[0]) { // parse function signature
				identToken := leftOver[1]
				if !isIdentifier(identToken) {
					var exp Exp
					return []Token{}, exp, &ParseError{"missing procedure identifier"}
				}

				// add Name
				varName = identToken.val
				leftOver = leftOver[2:]

				// add parameters
				for leftOver[0].tokType != TOK_RPAREN {
					if !isIdentifier(leftOver[0]) {
						var exp Exp
						return []Token{}, exp, &ParseError{"invalid function parameters"}
					}
					varOperands = append(varOperands, leftOver[0].val)
					leftOver = leftOver[1:]
				}

				if len(leftOver) == 0 || leftOver[0].tokType != TOK_RPAREN {
					var exp Exp
					return []Token{}, exp, &ParseError{"missing closing )"}
				}

				// skip right parenthesis
				leftOver = leftOver[1:]

				if len(leftOver) == 0 || !isLeftParenthesis(leftOver[0]) {
					var exp Exp
					return []Token{}, exp, &ParseError{"missing function expression"}
				}

				// parse function expression
				leftOver, varExpression, err = Parser(leftOver)

				if len(leftOver) == 0 || leftOver[0].tokType != TOK_RPAREN {
					var exp Exp
					return []Token{}, exp, &ParseError{"missing closing )"}
				} else {
					return []Token{}, &expDefineFunc{varName, varExpression, varOperands}, nil
				}

			} else { // parse variable definition
				if len(leftOver) == 0 || !isIdentifier(leftOver[0]) {
					var exp Exp
					return []Token{}, exp, &ParseError{"missing variable name"}
				}

				varName = leftOver[0].val
				leftOver = leftOver[1:]
				leftOver, varExpression, err = Parser(leftOver)

				if len(leftOver) == 0 || leftOver[0].tokType != TOK_RPAREN {
					var exp Exp
					return []Token{}, exp, &ParseError{"missing closing )"}
				} else {
					return []Token{}, &expDefineVar{varName, varExpression}, nil
				}
			}
		}
		if !isOperator(operatorToken) {
			var exp Exp
			return []Token{}, exp, &ParseError{"missing operator"}
		} else {
			var root Exp
			var operandList []Exp
			leftOver := tokens[2:]

			for len(leftOver) == 0 || leftOver[0].tokType != TOK_RPAREN {
				var subTree2 Exp
				leftOver, subTree2, err = Parser(leftOver)
				operandList = append(operandList, subTree2)
			}
			root = &expOperator{operatorToken.tokType, operandList}
			// ensure []Tokens remaining from second rec call is ")"
			if len(leftOver) == 0 || leftOver[0].tokType != TOK_RPAREN {
				var exp Exp
				return []Token{}, exp, &ParseError{"missing closing )"}
			}
			return leftOver[1:], root, err
		}
	} else {
		var exp Exp
		return []Token{}, exp, &ParseError{"missing ("}
	}
}
