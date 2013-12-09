package main

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"runtime/debug"
)

var digits []string = []string{"1","2","3","4","5","6","7","8","9","0"}

var delimiter []string = []string{" ", "\n", "\t", "\r"}

// var builtin map[Symbol]interface{} = make(map[Symbol]interface{})

// func init() {
// 	builtin[Symbol{"+"}] = plus
// 	builtin[Symbol{"-"}] = minus
// 	builtin[Symbol{"*"}] = multiply
// 	builtin[Symbol{"/"}] = divide
// 	builtin[Symbol{"car"}] = car
// 	builtin[Symbol{"cdr"}] = cdr
// 	builtin[Symbol{"define"}] = define

// 	return
// }

func newEnv(parent *Env) Env {
	env := new(Env)
	env.maps = make(map[Symbol]interface{})
	env.parent = parent
	return *env
}

type Env struct {
	maps map[Symbol]interface{}
	parent *Env
}

func (e *Env) Lookup (k Symbol) interface{} {
	if e.maps[k] != nil {
		return e.maps[k]
	} else if e.parent != nil {
		return e.parent.Lookup(k)
	} else {
		return nil
	}
}

func (e *Env) Bind(k Symbol, v interface{}) {
	e.maps[k] =v
	return
}

type Number struct {
	value int
}

type Symbol struct {
	value string
}

type String struct {
	value string
}

type Func struct {
	node []interface{}
}

type Bool struct {
	value bool
}

type Nil struct {
	value string
}

var True Bool = Bool{true}

var False Bool = Bool{false}

var index int = 0

func newParser(content string) Parser {
	var tokens []string
	var token string
	i := 0
	
	for i=0; i<len(content); i++ {
		switch c := string(content[i]); c {
		case "(",")" : token = c
		case " ","\n","\t","\r" : { token = "" }
		case "\"" : {
			token = "\""
			i++
			for ;string(content[i])!="\""; i++ {
				token += string(content[i])
			}
		}
		default: {
			token = ""
			for i<len(content) {
				c := string(content[i])
				if c == " " || c == "\n" || c == "\t" || c == "\r" {
					break
				} else if c == ")" || c =="(" {
					i--
					break
				} else {
					token += c
					i++
				}
			}
		}
		}
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	p := Parser{tokens, 0}
	return p
}

type Parser struct {
	Tokens []string
	index  int
}

func (p *Parser) GetToken() string {
	var token string
	if p.index == len(p.Tokens) {
		token = ""
	} else {
		token = p.Tokens[p.index]
		p.index++
	}
	return token
}

func (p *Parser) UngetToken() {
	if p.index > 0 {
		p.index--
	}
	return
}

func (p *Parser) Parse() []interface{} {
	var node []interface{}
	var token string
	digits,_ := regexp.Compile("[0-9]+")
	strings,_ := regexp.Compile("\".*\"")
Loop:
	for {
		token = p.GetToken()
		switch {
		case token == "" : break Loop 
		case token == "(" : {
			result := p.Parse()
			node = append(node, result)
		}
		case token == ")" : return node
		case token == "#t": node = append(node, True)
		case token == "#f": node = append(node, False)
		case digits.MatchString(token): {
			n,err := strconv.Atoi(token)
			checkErr(err)
			node = append(node, Number{n})
		}
		case strings.MatchString(token): {
			node = append(node,String{token})
		}
		default: node = append(node, Symbol{token})
		}
	}
	return node
}

func eval(content string) {
	parser := newParser(content)
	fmt.Println("Eval!")
	root := parser.Parse()
	env := newEnv(nil)
	
	for _,v := range root {
//		fmt.Println(v)
		fmt.Println(doEval(v,env))
	}
}

func doEval(node interface{}, env Env) interface{} {
 	fmt.Println("Start Eval:", node)
//	fmt.Println(env)
 	var result interface{}
	
	switch t := node.(type) {
	case Symbol : {
		value := env.Lookup(t)
		fmt.Printf("The value of symbol: %v is %v\n", t, value)
		if value == nil {
			debug.PrintStack()
			log.Fatal("Undefined Symbol: ", t)
		}
		result = value
	}
		
	case Number, String, Bool : { result = t }
	case []interface{} : {
		if v,ok := t[0].(Symbol); ok {
			switch v.value {
			case "define" : {
				env.Bind(t[1].(Symbol), t[2])
				fmt.Println("Binding function: ", t[2])
				return Nil{"nil"}
			}
			case "cond" : {
				for i:=1; i<len(t)-1; i++ {
//					fmt.Println("Environ:", env)
					if doEval(t[i].([]interface{})[0], env) == True {
						fmt.Println("Enter branch:", t[i])
						return doEval(t[i].([]interface{})[1], env)
					}
				}
				if (t[len(t)-1].([]interface{})[0] == Symbol{"else"}) {
					fmt.Println("Enter else:", t[len(t)-1])
					return doEval(t[len(t)-1].([]interface{})[1], env)
				}
				result = Nil{"nil"}
			}
			case "quote": return t[1]
			default: {
				if segs := env.Lookup(v); segs != nil {
					var expr []interface{}
					fmt.Println("The function is: ",segs)
					expr = append(expr,segs)
					expr = append(expr,t[1:]...)
					fmt.Println("To be excuted:", expr)
					return doEval(expr, newEnv(&env))
				} 
			}
			}
		}

		switch v := t[0].(type) {
		case []interface{} : {
			fmt.Println("Whole expr:", t)
			fmt.Println("Lambda? :", v)
			for k,v := range v[1].([]interface{}) {
				env.Bind(v.(Symbol), doEval(t[k+1], env))
			}
			fmt.Println("lambda biding:", env)
			result = doEval(v[2], newEnv(&env))
		}
			
		case Symbol: {
			var expr []interface{}
			for i:=1; i<len(t); i++ {
				expr = append(expr, doEval(t[i], env))
			}
//			fmt.Println(t)
			switch v.value {
			case "+" : result = plus(expr)
			case "-" : result = minus(expr)
			case "*" : result = multiply(expr)
			case "/" : result = divide(expr)
			case "car" : result = car(expr)
			case "cdr" : result = cdr(expr)
			case "cons" : result = cons(expr)
			case "eq?" : result = eqf(expr)
			case "null?" : result = nullf(expr)
			case "atom?" : result = atomf(expr)
			case "zero?" : result = zerof(expr)
			default: {
				log.Fatal("Undefined function:", v.value)
			}
			}
		}
		default: {
			log.Fatal("The first argument of a list should be an operator")
		}
		}
	}
	}
	return result
}

func zerof(args []interface{}) interface{} {
	zero,ok := args[0].(Number)
	if !ok {
		log.Fatal("Argument of zero? should be number.")
	}
	return Bool{zero.value == 0}
}

func car(args []interface{}) interface{} {
	if len(args) == 0 {
		log.Fatal("car on a null list")
	}
	
	if list,ok := args[0].([]interface{}); ok {
		return list[0]
	} else {
		log.Fatal("Can only car on a list")
		return -1
	}
}

func cdr(args []interface{}) []interface{} {
	if list,ok := args[0].([]interface{}); ok && len(list)>0 {
		return list[1:]
	}
	return nil
}

func cons(args []interface{}) []interface{} {
	list,ok := args[1].([]interface{})
	if !ok { panic("Error")	}
	newlist := make([]interface{}, len(list)+1)
	newlist[0] = args[0]
	copy(newlist[1:], list[:])

	return newlist
}

func cond(args []interface{}) interface{} {
	for _,value := range args {
		seg,_ := value.([]interface{})
		if test,ok := seg[0].(Bool); ok && test == True {
			return seg[1]
		} else {
			if Else,ok := seg[0].(Symbol); ok && (Symbol{"else"} == Else) {
			return seg[1]
			}
		}
	}
	return Nil{"nil"}
}

func eqf(args []interface{}) Bool {
	if reflect.DeepEqual(args[0], args[1]) {
		return True
	} else {
		return False
	}
}

func nullf(args []interface{}) Bool {
	list, ok := args[0].([]interface{})
	if !ok {
		log.Fatal("The arguments of null? should be a list")
	}
	return Bool{len(list)==0}
}

func atomf(args []interface{}) Bool {
	switch args[0].(type) {
	case Number,Symbol,String : {return True}
	default: return False
	}
}

func plus(args []interface{}) Number {
	number,ok := args[0].(Number)
	checkFalse(ok)
	result := number.value
	for i:=1; i<len(args); i++ {
		number,ok := args[i].(Number)
		if  !ok {
			log.Fatal("arguments for add must be numbers")
		}
		result = result + number.value
	}
	return Number{result}
}

func minus(args []interface{}) Number{
	number,ok := args[0].(Number)
	checkFalse(ok)
	result := number.value
	for i:=1; i<len(args); i++ {
		number,ok := args[i].(Number)
		if !ok {
			log.Fatal("arguments for add must be numbers")
		}
		result = result - number.value
	}
	return Number{result}

}


func multiply(args []interface{}) Number{
	number,ok := args[0].(Number)
	checkFalse(ok)
	result := number.value
	for i:=1; i<len(args); i++ {
		number,ok := args[i].(Number);
		if !ok {
			log.Fatal("arguments for add must be numbers")
		}

		result = result * number.value
	}
	return Number{result}
}


func divide(args []interface{}) Number{
	number,ok := args[0].(Number)
	checkFalse(ok)
	result := number.value
	for i:=1; i<len(args); i++ {
		number,ok := args[i].(Number)
		if !ok {
			log.Fatal("arguments for add must be numbers")
		}
		result = result / number.value
	}
	return Number{result}
}


func main () {
	if len(os.Args) != 2 {
		log.Fatal("Usage: gscheme <scheme file>")
	}
	log.SetFlags(log.Lshortfile)

	content,err := ioutil.ReadFile(os.Args[1])
	checkErr(err)

	eval(string(content))
	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkFalse(ok bool) {
	if !ok {
		panic("Not OK!")
	}
}

func in(a string, list []string) bool {
	for _,value :=range list {
		if a == value {
			return true
		}
	}
	return false
}
