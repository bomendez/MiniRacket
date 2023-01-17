package minrkt

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParserNoErrors(t *testing.T) {
	emptyEnv := &Environment{}

	// Two Operands
	inputTokList1 := []Token{
		Token{TokenType(0), "("},
		Token{TokenType(7), "="},
		Token{TokenType(2), "5"},
		Token{TokenType(2), "5"},
		Token{TokenType(1), ")"},
	}
	var tokRem1 = []Token{}
	var expTrue Exp
	var tokList1 []Exp
	var tokListOperand1A Exp
	tokListOperand1A = &expNumConst{5}
	tokList1 = append(tokList1, tokListOperand1A)
	var tokListOperand1B Exp
	tokListOperand1B = &expNumConst{5}
	tokList1 = append(tokList1, tokListOperand1B)
	expTrue = &expOperator{TOK_EQ, tokList1}

	// Nested Operators
	inputTokList2 := []Token{
		Token{TokenType(0), "("},
		Token{TokenType(4), "-"},
		Token{TokenType(2), "10"},
		Token{TokenType(0), "("},
		Token{TokenType(5), "*"},
		Token{TokenType(2), "0.5"},
		Token{TokenType(2), "50"},
		Token{TokenType(1), ")"},
		Token{TokenType(1), ")"},
	}
	var tokRem2 = []Token{}
	var expFalse Exp
	var tokList2 []Exp
	var tokListOperand2A Exp
	tokListOperand2A = &expNumConst{10}
	tokList2 = append(tokList2, tokListOperand2A)
	var tokListOperand2B Exp
	// operand node
	var expSub2 Exp
	var tokListSub2 []Exp
	var tokListOperandSub2A Exp
	tokListOperandSub2A = &expNumConst{0.5}
	tokListSub2 = append(tokListSub2, tokListOperandSub2A)
	var tokListOperandSub2B Exp
	tokListOperandSub2B = &expNumConst{50}
	tokListSub2 = append(tokListSub2, tokListOperandSub2B)
	expSub2 = &expOperator{TOK_MUL, tokListSub2}
	// add to parent
	tokListOperand2B = expSub2
	tokList2 = append(tokList2, tokListOperand2B)
	expFalse = &expOperator{TOK_SUB, tokList2}

	// Three Operands
	inputTokList3 := []Token{
		Token{TokenType(0), "("},
		Token{TokenType(5), "*"},
		Token{TokenType(2), "2"},
		Token{TokenType(0), "("},
		Token{TokenType(6), "/"},
		Token{TokenType(2), "10"},
		Token{TokenType(2), "5"},
		Token{TokenType(1), ")"},
		Token{TokenType(2), "2"},
		Token{TokenType(1), ")"},
	}
	var tokRem3 = []Token{}
	var exp3 Exp
	var tokList3 []Exp
	var tokListOperand3A Exp
	tokListOperand3A = &expNumConst{2}
	tokList3 = append(tokList3, tokListOperand3A)
	var tokListOperand3B Exp
	// operand node 2
	var expSub3 Exp
	var tokListSub3 []Exp
	var tokListOperandSub3A Exp
	tokListOperandSub3A = &expNumConst{10}
	tokListSub3 = append(tokListSub3, tokListOperandSub3A)
	var tokListOperandSub3B Exp
	tokListOperandSub3B = &expNumConst{5}
	tokListSub3 = append(tokListSub3, tokListOperandSub3B)
	expSub3 = &expOperator{TOK_DIV, tokListSub3}
	// add to parent
	tokListOperand3B = expSub3
	tokList3 = append(tokList3, tokListOperand3B)
	tokList3 = append(tokList3, &expNumConst{2})
	exp3 = &expOperator{TOK_MUL, tokList3}

	// One operand
	inputTokList4 := []Token{
		Token{TokenType(0), "("},
		Token{TokenType(3), "+"},
		Token{TokenType(2), "5"},
		Token{TokenType(1), ")"},
	}
	tokRem4 := []Token{}
	var exp4 Exp
	exp4 = &expNumConst{5}

	var tests = []struct {
		a       []Token
		wantTok []Token
		wantE   Exp
		wantErr error
	}{
		{inputTokList1, tokRem1, expTrue, nil},
		{inputTokList2, tokRem2, expFalse, nil},
		{inputTokList3, tokRem3, exp3, nil},
		{inputTokList4, tokRem4, exp4, nil},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			tok, e, err := Parser(tt.a)
			eVal, _ := e.Eval(emptyEnv)
			ttVal, _ := tt.wantE.Eval(emptyEnv)
			if !reflect.DeepEqual(tok, tt.wantTok) || eVal != ttVal || err != tt.wantErr {
				t.Errorf("got %v %v %v, want %v %v %v", tok, e, err, tt.wantTok, tt.wantE, tt.wantErr)
			}
		})
	}
}

func TestParserMissingLParen(t *testing.T) {
	inputTokList := []Token{
		Token{TokenType(3), "+"},
		Token{TokenType(2), "1"},
		Token{TokenType(2), "5"},
	}
	tokRem := []Token{}
	var exp Exp
	wantErr := &ParseError{"missing ("}
	gotT, gotE, gotErr := Parser(inputTokList)
	if !reflect.DeepEqual(gotT, tokRem) {
		t.Errorf("Parser(): got remainder: %v want: %v", gotT, tokRem)
	}
	if gotE != exp {
		t.Errorf("Parser(): got Exp: %v want: %v", gotE, exp)
	}
	if gotErr.Error() != wantErr.Error() {
		t.Errorf("Parser(): got error: %v want: %v", gotErr, wantErr)
	}
}

func TestParserMissingOperator(t *testing.T) {
	inputTokList := []Token{
		Token{TokenType(0), "("},
		Token{TokenType(2), "1"},
		Token{TokenType(2), "5"},
		Token{TokenType(1), ")"},
	}
	tokRem := []Token{}
	var exp Exp
	wantErr := &ParseError{"missing operator"}
	gotT, gotE, gotErr := Parser(inputTokList)
	if !reflect.DeepEqual(gotT, tokRem) {
		t.Errorf("Parser(): got remainder: %v want: %v", gotT, tokRem)
	}
	if gotE != exp {
		t.Errorf("Parser(): got Exp: %v want: %v", gotE, exp)
	}
	if gotErr.Error() != wantErr.Error() {
		t.Errorf("Parser(): got error: %v want: %v", gotErr, wantErr)
	}
}

// Parser should not throw an error for extra characters.
// But should return the extra characters as tokens to main()
func TestParserExtraToken(t *testing.T) {
	inputTokList := []Token{
		Token{TokenType(0), "("},
		Token{TokenType(3), "+"},
		Token{TokenType(2), "1"},
		Token{TokenType(2), "5"},
		Token{TokenType(1), ")"},
		Token{TokenType(2), "5"},
	}
	emptyEnv := &Environment{}
	tokRem := []Token{Token{TokenType(2), "5"}}
	var exp Exp
	var tokList1 []Exp
	var tokListOperand1A Exp
	tokListOperand1A = &expNumConst{1}
	tokList1 = append(tokList1, tokListOperand1A)
	var tokListOperand1B Exp
	tokListOperand1B = &expNumConst{5}
	tokList1 = append(tokList1, tokListOperand1B)
	exp = &expOperator{TOK_ADD, tokList1}
	// wantErr := &ParseError{"incomplete statement"}
	gotT, gotE, gotErr := Parser(inputTokList)
	gotVal, _ := gotE.Eval(emptyEnv)
	ttVal, _ := exp.Eval(emptyEnv)
	if !reflect.DeepEqual(gotT, tokRem) {
		t.Errorf("Parser(): got remainder: %v want: %v", gotT, tokRem)
	}
	if gotVal != ttVal {
		t.Errorf("Parser(): got Exp: %v want: %v", gotE, exp)
	}
	if gotErr != nil {
		t.Errorf("Parser(): got error: %v", gotErr)
	}
}

func TestBuildOperandNode(t *testing.T) {
	var tests = []struct {
		a        Token
		wantBool Exp
	}{
		{Token{TokenType(2), "5.1"}, &expNumConst{5.1}},
		{Token{TokenType(16), "true"}, &expBoolConst{true}},
		{Token{TokenType(16), "#t"}, &expBoolConst{true}},
		{Token{TokenType(17), "false"}, &expBoolConst{false}},
		{Token{TokenType(17), "#f"}, &expBoolConst{false}},
		// {Token{TokenType(19), "x"}, &expVarConst{"x"}},
		// {Token{TokenType(19), "num_Const1"}, &expVarConst{"num_Const1"}},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			ttWantBool, _ := tt.wantBool.Eval(&Environment{})
			e := buildOperandNode(tt.a)
			eVal, _ := e.Eval(&Environment{})
			if eVal != ttWantBool {
				t.Errorf("Build operand got %v, want %v", e, ttWantBool)
			}
		})
	}
}

func TestEvalNum(t *testing.T) {
	var tests = []struct {
		a    *expNumConst
		want float64
	}{
		{&expNumConst{5}, 5},
		{&expNumConst{0}, 0},
		{&expNumConst{1}, 1},
		{&expNumConst{87.46}, 87.46},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestEvalVar(t *testing.T) {
	env := &Environment{}
	env.CallStack = make([]map[string]interface{}, 0)
	env.CallStack = append(env.CallStack, map[string]interface{}{"a": 10})
	env.Variables = make(map[string]interface{})
	env.Variables["x"] = 1
	env.Variables["five_5"] = 5.5
	env.Variables["true"] = true
	env.Variables["false"] = false
	var tests = []struct {
		a    *expVar
		want interface{}
	}{
		{&expVar{"a"}, 10},
		{&expVar{"x"}, 1},
		{&expVar{"five_5"}, 5.5},
		{&expVar{"true"}, true},
		{&expVar{"false"}, false},
		{&expVar{"undefined"}, "undefined"},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(env)
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestEvalFunc(t *testing.T) {
	env := &Environment{}
	env.Functions = make(map[string]FuncParamExpr)
	env.CallStack = make([]map[string]interface{}, 0)

	var add5Exp Exp
	operandList := []Exp{&expVar{"a"}, &expNumConst{5.5}}
	add5Exp = &expOperator{TOK_ADD, operandList}
	funcStr := FuncParamExpr{[]string{"a"}, add5Exp}
	env.Functions["add5"] = funcStr

	var times_5_p_5 Exp
	operandList2 := []Exp{&expVar{"b"}, &expVar{"five_5"}}
	times_5_p_5 = &expOperator{TOK_MUL, operandList2}
	funcStr2 := FuncParamExpr{[]string{"b"}, times_5_p_5}
	env.Functions["times5_5"] = funcStr2

	env.Variables = make(map[string]interface{})
	env.Variables["x"] = 1
	env.Variables["five_5"] = 5.5
	env.Variables["true"] = true
	env.Variables["false"] = false

	var tests = []struct {
		a    *expFunc
		want interface{}
	}{
		{&expFunc{"add5", []Exp{&expNumConst{1}}}, 6.5},
		{&expFunc{"times5_5", []Exp{&expNumConst{5}}}, 27.5},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(env)
			want := tt.want
			if !reflect.DeepEqual(e, want) {
				t.Errorf("got %v, want %v", e, want)
			}
		})
	}
}

func TestAddOperator(t *testing.T) {
	var tests = []struct {
		a    *expOperator
		want float64
	}{
		{&expOperator{TOK_ADD, []Exp{&expNumConst{5}}}, 5},
		{&expOperator{TOK_ADD, []Exp{}}, 0},
		{&expOperator{TOK_ADD, []Exp{&expNumConst{5}, &expNumConst{5}}}, 10},
		{&expOperator{TOK_ADD, []Exp{&expNumConst{2.5}, &expNumConst{5.5}, &expNumConst{2}}}, 10},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestSubtractOperator(t *testing.T) {
	var tests = []struct {
		a    *expOperator
		want float64
	}{
		{&expOperator{TOK_SUB, []Exp{&expNumConst{5}}}, -5},
		{&expOperator{TOK_SUB, []Exp{&expNumConst{0}}}, 0},
		{&expOperator{TOK_SUB, []Exp{&expNumConst{5}, &expNumConst{5}}}, 0},
		{&expOperator{TOK_SUB, []Exp{&expNumConst{10.5}, &expNumConst{5.5}, &expNumConst{2}}}, 3},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestMultipyOperator(t *testing.T) {
	var tests = []struct {
		a    *expOperator
		want float64
	}{
		{&expOperator{TOK_MUL, []Exp{&expNumConst{5}}}, 5},
		{&expOperator{TOK_MUL, []Exp{}}, 1},
		{&expOperator{TOK_MUL, []Exp{&expNumConst{1}}}, 1},
		{&expOperator{TOK_MUL, []Exp{&expNumConst{5}, &expNumConst{5}}}, 25},
		{&expOperator{TOK_MUL, []Exp{&expNumConst{0.5}, &expNumConst{0.5}, &expNumConst{4}}}, 1},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestDivideOperator(t *testing.T) {
	var tests = []struct {
		a    *expOperator
		want float64
	}{
		{&expOperator{TOK_DIV, []Exp{&expNumConst{5}}}, 0},
		//{&expOperator{TOK_DIV, []Exp{&expNumConst{0}, &expNumConst{8.64}}}, 0},
		{&expOperator{TOK_DIV, []Exp{&expNumConst{5}, &expNumConst{5}}}, 1},
		{&expOperator{TOK_DIV, []Exp{&expNumConst{100}, &expNumConst{25}, &expNumConst{2}}}, 2},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestArithmeticComparisons(t *testing.T) {
	var tests = []struct {
		a    *expOperator
		want bool
	}{
		{&expOperator{TOK_EQ, []Exp{&expNumConst{5}, &expNumConst{5}}}, true},
		{&expOperator{TOK_EQ, []Exp{&expNumConst{5}, &expNumConst{5.5}}}, false},
		{&expOperator{TOK_GTEQ, []Exp{&expNumConst{5}, &expNumConst{5}}}, true},
		{&expOperator{TOK_GTEQ, []Exp{&expNumConst{5.5}, &expNumConst{5}}}, true},
		{&expOperator{TOK_GTEQ, []Exp{&expNumConst{1}, &expNumConst{5}}}, false},
		{&expOperator{TOK_LTEQ, []Exp{&expNumConst{5}, &expNumConst{5}}}, true},
		{&expOperator{TOK_LTEQ, []Exp{&expNumConst{4}, &expNumConst{5}}}, true},
		{&expOperator{TOK_LTEQ, []Exp{&expNumConst{225}, &expNumConst{22.4}}}, false},
		{&expOperator{TOK_GT, []Exp{&expNumConst{5.5}, &expNumConst{5}}}, true},
		{&expOperator{TOK_GT, []Exp{&expNumConst{5}, &expNumConst{5}}}, false},
		{&expOperator{TOK_GT, []Exp{&expNumConst{2}, &expNumConst{5}}}, false},
		{&expOperator{TOK_LT, []Exp{&expNumConst{2}, &expNumConst{5}}}, true},
		{&expOperator{TOK_LT, []Exp{&expNumConst{5}, &expNumConst{5}}}, false},
		{&expOperator{TOK_LT, []Exp{&expNumConst{10.5}, &expNumConst{5}}}, false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestBooleanOperators(t *testing.T) {
	var expTrue Exp
	var tokList1 []Exp
	var tokListOperand1A Exp
	tokListOperand1A = &expNumConst{1}
	tokList1 = append(tokList1, tokListOperand1A)
	var tokListOperand1B Exp
	tokListOperand1B = &expNumConst{1}
	tokList1 = append(tokList1, tokListOperand1B)
	expTrue = &expOperator{TOK_EQ, tokList1}

	var expFalse Exp
	var tokList2 []Exp
	var tokListOperand2A Exp
	tokListOperand2A = &expNumConst{2}
	tokList2 = append(tokList2, tokListOperand2A)
	var tokListOperand2B Exp
	tokListOperand2B = &expNumConst{2}
	tokList2 = append(tokList2, tokListOperand2B)
	expFalse = &expOperator{TOK_GT, tokList2}

	var tests = []struct {
		a    *expOperator
		want bool
	}{
		{&expOperator{TOK_AND, []Exp{expTrue, expTrue}}, true},
		{&expOperator{TOK_AND, []Exp{expTrue, expFalse}}, false},
		{&expOperator{TOK_AND, []Exp{expFalse, expTrue}}, false},
		{&expOperator{TOK_AND, []Exp{expFalse, expFalse}}, false},
		{&expOperator{TOK_OR, []Exp{expTrue, expTrue}}, true},
		{&expOperator{TOK_OR, []Exp{expTrue, expFalse}}, true},
		{&expOperator{TOK_OR, []Exp{expFalse, expTrue}}, true},
		{&expOperator{TOK_OR, []Exp{expFalse, expFalse}}, false},
		{&expOperator{TOK_NOT, []Exp{expFalse}}, true},
		{&expOperator{TOK_NOT, []Exp{expTrue}}, false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestIfOperator(t *testing.T) {
	var expTrue Exp
	var tokList1 []Exp
	var tokListOperand1A Exp
	tokListOperand1A = &expNumConst{1}
	tokList1 = append(tokList1, tokListOperand1A)
	var tokListOperand1B Exp
	tokListOperand1B = &expNumConst{1}
	tokList1 = append(tokList1, tokListOperand1B)
	expTrue = &expOperator{TOK_EQ, tokList1}

	var expFalse Exp
	var tokList2 []Exp
	var tokListOperand2A Exp
	tokListOperand2A = &expNumConst{2}
	tokList2 = append(tokList2, tokListOperand2A)
	var tokListOperand2B Exp
	tokListOperand2B = &expNumConst{2}
	tokList2 = append(tokList2, tokListOperand2B)
	expFalse = &expOperator{TOK_GT, tokList2}

	var expTen Exp
	var tokList3 []Exp
	var tokListOperand3A Exp
	tokListOperand3A = &expNumConst{5}
	tokList3 = append(tokList3, tokListOperand3A)
	var tokListOperand3B Exp
	tokListOperand3B = &expNumConst{5}
	tokList3 = append(tokList3, tokListOperand3B)
	expTen = &expOperator{TOK_ADD, tokList3}

	var tests = []struct {
		a    *expOperator
		want interface{}
	}{
		{&expOperator{TOK_IF, []Exp{expTrue, expTrue, expTen}}, true},
		{&expOperator{TOK_IF, []Exp{expFalse, expTrue, expFalse}}, false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			e, _ := tt.a.Eval(&Environment{})
			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}

func TestIsOperator(t *testing.T) {
	var tests = []struct {
		a    Token
		want bool
	}{
		{Token{TokenType(0), "("}, false},
		{Token{TokenType(1), ")"}, false},
		{Token{TokenType(2), "6.2"}, false},
		{Token{TokenType(3), "+"}, true},
		{Token{TokenType(4), "-"}, true},
		{Token{TokenType(5), "*"}, true},
		{Token{TokenType(6), "/"}, true},
		{Token{TokenType(7), "="}, true},
		{Token{TokenType(8), ">="}, true},
		{Token{TokenType(9), "<="}, true},
		{Token{TokenType(10), ">"}, true},
		{Token{TokenType(11), "<"}, true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			tok := isOperator(tt.a)
			if tok != tt.want {
				t.Errorf("got %v, want %v", tok, tt.want)
			}
		})
	}
}

func TestIsOperand(t *testing.T) {
	if got, want := isOperand(Token{TokenType(2), "2.5"}), true; got != want {
		t.Errorf("isOperand(): got: %v want: %v", got, want)
	}
}

func TestIsNotOperand(t *testing.T) {
	if got, want := isOperand(Token{TokenType(0), "("}), false; got != want {
		t.Errorf("isOperand(): got: %v want: %v", got, want)
	}
}
