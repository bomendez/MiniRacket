package minrkt

import (
	"fmt"
	"testing"
)

func TestEvaluator(t *testing.T) {
	emptyEnv := &Environment{}

	var want1 interface{}
	var want1Sub float64 = 6
	want1 = want1Sub
	var exp1 Exp
	var tokList1 []Exp
	var tokListOperand1A Exp
	tokListOperand1A = &expNumConst{1}
	tokList1 = append(tokList1, tokListOperand1A)
	var tokListOperand1B Exp
	tokListOperand1B = &expNumConst{5}
	tokList1 = append(tokList1, tokListOperand1B)
	exp1 = &expOperator{TOK_ADD, tokList1}

	var want2 interface{}
	var want2Sub float64 = 52
	want2 = want2Sub
	var exp2 Exp
	var tokList2 []Exp
	var tokListOperand2A Exp
	tokListOperand2A = &expNumConst{2}
	tokList2 = append(tokList2, tokListOperand2A)
	var tokListOperand2B Exp
	// operand 2 is operator subTree
	var expSub2 Exp
	var tokListSub2 []Exp
	var tokListOperandSub2A Exp
	tokListOperandSub2A = &expNumConst{10}
	tokListSub2 = append(tokListSub2, tokListOperandSub2A)
	var tokListOperandSub2B Exp
	tokListOperandSub2B = &expNumConst{5}
	tokListSub2 = append(tokListSub2, tokListOperandSub2B)
	expSub2 = &expOperator{TOK_MUL, tokListSub2}
	// add to parent
	tokListOperand2B = expSub2
	tokList2 = append(tokList2, tokListOperand2B)
	exp2 = &expOperator{TOK_ADD, tokList2}

	var want3 interface{}
	var want3Sub string = "#t"
	want3 = want3Sub
	var exp3 Exp
	var tokList3 []Exp
	var tokListOperand3A Exp
	tokListOperand3A = &expNumConst{5}
	tokList3 = append(tokList3, tokListOperand3A)
	var tokListOperand3B Exp
	tokListOperand3B = &expNumConst{5}
	tokList3 = append(tokList3, tokListOperand3B)
	exp3 = &expOperator{TOK_EQ, tokList3}

	var want4 interface{}
	var want4Sub string = "#t"
	want4 = want4Sub
	var exp4 Exp
	var tokList4 []Exp
	var tokListOperand4A Exp
	tokListOperand4A = &expNumConst{5.5}
	tokList4 = append(tokList4, tokListOperand4A)
	var tokListOperand4B Exp
	tokListOperand4B = &expNumConst{5}
	tokList4 = append(tokList4, tokListOperand4B)
	exp4 = &expOperator{TOK_GT, tokList4}

	var want5 interface{}
	var want5Sub string = "#t"
	want5 = want5Sub
	var exp5 Exp
	var tokList5 []Exp
	var tokListOperand5A Exp
	tokListOperand5A = &expNumConst{2}
	tokList5 = append(tokList5, tokListOperand5A)
	var tokListOperand5B Exp
	tokListOperand5B = &expNumConst{5}
	tokList5 = append(tokList5, tokListOperand5B)
	exp5 = &expOperator{TOK_LT, tokList5}

	var wantTrue interface{}
	var wantTrueSub string = "#t"
	wantTrue = wantTrueSub
	var expTrue Exp = &expBoolConst{true}

	var wantFalse interface{}
	var wantFalseSub string = "#f"
	wantFalse = wantFalseSub
	var expFalse Exp = &expBoolConst{false}

	var tests = []struct {
		a       Exp
		b       *Environment
		wantF   interface{}
		wantErr error
	}{
		{exp1, emptyEnv, want1, nil},
		{exp2, emptyEnv, want2, nil},
		{exp3, emptyEnv, want3, nil},
		{exp4, emptyEnv, want4, nil},
		{exp5, emptyEnv, want5, nil},
		{expTrue, emptyEnv, wantTrue, nil},
		{expFalse, emptyEnv, wantFalse, nil},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			f, err := Evaluator(tt.a, tt.b)
			switch f.(type) {
			case float64:
				fmt.Printf("float64: %v\n", f)
			case bool:
				fmt.Printf("bool: %v\n", f)
			case string:
				fmt.Printf("string: %v\n", f)
			}
			if f != tt.wantF || err != tt.wantErr {
				t.Errorf("got %v %v, want %v %v", f, err, tt.wantF, tt.wantErr)
			}
		})
	}
}
