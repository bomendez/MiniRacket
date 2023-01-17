package minrkt

func Evaluator(root Exp, env *Environment) (interface{}, error) {
	var result interface{}
	var err error
	result, err = root.Eval(env)
	if result == true {
		result = "#t"
	} else if result == false {
		result = "#f"
	}
	return result, err
}
