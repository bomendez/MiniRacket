package minrkt

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	var tests = []struct {
		a       string
		wantTok []Token
		wantE   error
	}{
		{"()", []Token{
			{TokenType(0), "("},
			{TokenType(1), ")"},
		}, nil},
		{"(+ 1 2)", []Token{
			{TokenType(0), "("},
			{TokenType(3), "+"},
			{TokenType(2), "1"},
			{TokenType(2), "2"},
			{TokenType(1), ")"},
		}, nil},
		{"(= 5.6 9)", []Token{
			{TokenType(0), "("},
			{TokenType(7), "="},
			{TokenType(2), "5.6"},
			{TokenType(2), "9"},
			{TokenType(1), ")"},
		}, nil},
		{"-", []Token{{TokenType(4), "-"}}, nil},
		{"*", []Token{{TokenType(5), "*"}}, nil},
		{"/", []Token{{TokenType(6), "/"}}, nil},
		{"=", []Token{{TokenType(7), "="}}, nil},
		{">=", []Token{{TokenType(8), ">="}}, nil},
		{"<=", []Token{{TokenType(9), "<="}}, nil},
		{">", []Token{{TokenType(10), ">"}}, nil},
		{"<", []Token{{TokenType(11), "<"}}, nil},
		{"and", []Token{{TokenType(12), "and"}}, nil},
		{"or", []Token{{TokenType(13), "or"}}, nil},
		{"not", []Token{{TokenType(14), "not"}}, nil},
		{"if", []Token{{TokenType(15), "if"}}, nil},
		{"true", []Token{{TokenType(16), "true"}}, nil},
		{"#t", []Token{{TokenType(16), "#t"}}, nil},
		{"false", []Token{{TokenType(17), "false"}}, nil},
		{"#f", []Token{{TokenType(17), "#f"}}, nil},
		{"define", []Token{{TokenType(18), "define"}}, nil},
		{"x", []Token{{TokenType(19), "x"}}, nil},
		{"numCount", []Token{{TokenType(19), "numCount"}}, nil},
		{"Num_count2", []Token{{TokenType(19), "Num_count2"}}, nil},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.a)
		t.Run(testname, func(t *testing.T) {
			tok, e := Tokenizer(tt.a)
			if !(reflect.DeepEqual(tok, tt.wantTok)) || e != tt.wantE {
				t.Errorf("got %v %v, want %v %v", tok, e, tt.wantTok, tt.wantE)
			}
		})
	}
}

func TestNextToken(t *testing.T) {
	var tests = []struct {
		a       string
		wantTok Token
		wantS   string
		wantE   error
	}{
		{"(+ 1 2)", Token{TokenType(0), "("}, "+ 1 2)", nil},
		{"= 1 2)", Token{TokenType(7), "="}, " 1 2)", nil},
		{"+12)", Token{TokenType(3), "+"}, "12)", nil},
		{"22.9 3.2", Token{TokenType(2), "22.9"}, " 3.2", nil},
		{"", Token{}, "", nil},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.a)
		t.Run(testname, func(t *testing.T) {
			tok, s, e := NextToken(tt.a)
			if tok != tt.wantTok || s != tt.wantS || e != tt.wantE {
				t.Errorf("got %v %s %v, want %v %s %v", tok, s, e, tt.wantTok, tt.wantS, tt.wantE)
			}
		})
	}
}

func TestGetTokenIndx(t *testing.T) {
	var tests = []struct {
		a    []string
		want int
	}{
		{[]string{"(", "("}, 0},
		{[]string{"(", "", "(", "", ""}, 1},
		{[]string{"77", "", "", "77"}, 2},
		{[]string{"7.2", "", "", "7.2"}, 2},
		{[]string{"/", "", "", "", "", "", "", "/"}, 6},
		{[]string{"=", "", "", "", "", "", "", "", "="}, 7},
		{[]string{">=", "", "", "", "", "", "", "", "", ">="}, 8},
		{[]string{"<=", "", "", "", "", "", "", "", "", "", "<="}, 9},
		{[]string{">", "", "", "", "", "", "", "", "", "", "", ">"}, 10},
		{[]string{"<", "", "", "", "", "", "", "", "", "", "", "", "<"}, 11},
		{[]string{"and", "", "", "", "", "", "", "", "", "", "", "", "", "and"}, 12},
		{[]string{"or", "", "", "", "", "", "", "", "", "", "", "", "", "", "or"}, 13},
		{[]string{"not", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "not"}, 14},
		{[]string{"if", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "if"}, 15},
		{[]string{"true", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "true"}, 16},
		{[]string{"#t", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "#t"}, 16},
		{[]string{"false", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "false"}, 17},
		{[]string{"#f", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "#f"}, 17},
		{[]string{"define", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "define"}, 18},
		{[]string{"x", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "x"}, 19},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			ans := getTokenIndx(tt.a)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
