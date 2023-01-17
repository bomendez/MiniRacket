package main

import (
	"bufio"
	"fmt"
	"os"

	"my.com/cs5400/minrkt"
)

func main() {
	fmt.Println("Welcome to minimalistic racket!")
	env := &minrkt.Environment{}
	env.Variables = make(map[string]interface{})
	env.Functions = make(map[string]minrkt.FuncParamExpr)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if len(line) != 0 {
			tokens, err := minrkt.Tokenizer(line)
			if err != nil {
				fmt.Printf("Input Error: %v\n", err)
				continue
			}
			remaingToks, exp, err := minrkt.Parser(tokens)
			// if []Tokens is not empty throw error
			if len(remaingToks) != 0 {
				fmt.Println("Parse Error: missing )")
			}
			if err != nil {
				fmt.Printf("Parse Error: %v\n", err)
				continue
			}
			result, err := minrkt.Evaluator(exp, env)
			if err != nil {
				fmt.Println(err)
			} else if result != nil {
				fmt.Println(result)
			}
		}
	}
}
