package main

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"strconv"
)

var digits []string = []string{"1","2","3","4","5","6","7","8","9","0"}

var delimiter []string = []string{" ", "\n", "\t", "\r"}

var builtin map[Symbol]func([]interface{}) Number = make(map[Symbol]func([]interface{}) Number)

func init() {
	builtin[Symbol{"+"}] = plus
	builtin[Symbol{"-"}] = minus
	builtin[Symbol{"*"}] = multiply
	builtin[Symbol{"/"}] = divide

	return
}

var Env map[Symbol]interface{} = make(map[Symbol]interface{})
//type Node struct {
//	tokens []interface{}
//}

//type Root Node


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


var index int = 0

func parse(content string) []interface{} {
	var node []interface{}
	
	// for ;index<len(content); index++ {
	// 	c := string(content[index])
	// 	if c == " " || c == "\n" || c == "\t" || c=="\r" {
	// 		continue
	// 	} else {
	// 		if c == "(" {
	// 		break
	// 		}
	// 	}
	// }
	
	index++

	for ;index<len(content); index++ {
		c := string(content[index])
		switch c {
		case " ","\n","\t","\r" : {}
		case "(" : {
			node =append(node,parse(content))
		}
		case ")": {
//			index++
			return node
		}
		case "\"" : {
			var str string
			str += "\""
			index++
			for ;string(content[index]) != "\""; index++ {
				str += string(content[index])
			}
			node = append(node, String{str})
		}
		default: {
			var token string
			for ;index<len(content); index++ {
				c := string(content[index])
				if c == " " || c == "\n" || c == "\t" || c == "\r" {
					break
				} else if c == ")" {
					index--
					break
				} else {
					token += c
				}
			}
			if in(string(token[0]), digits) {
				n,err := strconv.Atoi(token)
				checkErr(err)
				node = append(node, Number{n})
			} else {
				node = append(node, Symbol{token})
			}
		}
		}
	}
	return node
}

func eval(content string) {
	root := parse(content)
	fmt.Println("Eval!")
	fmt.Println(root)
	fmt.Println(doEval(root))
}

func doEval(node []interface{}) interface{} {
	fmt.Println("The length of node is:", len(node))
	var result interface{}
	for i:=1; i<len(node); i++ {
		switch t := node[i].(type) {
		case []interface{} : {
			result := doEval(t)
			node[i] = result
		}
		case Symbol : {
			node[i] = Env[t]
		}
		}
	}
	
	switch operator := node[0].(type) {
	case []interface{} :  {
		if token,ok := operator[0].(Symbol); ok && token.value != "lambda" {
			log.Fatal("The first element should be a function")
		}
		subnode,ok := operator[1].([]interface{})
		checkFalse(ok)
		fmt.Println("The length of subnode is:", len(subnode))
		for k,v := range(subnode) {
			varible,ok := v.(Symbol)
			checkFalse(ok)
			Env[varible] = node[k+1]
		}
		subnode2, ok := operator[2].([]interface{})
		checkFalse(ok)
		result = doEval(subnode2)
	}
	case Symbol :{
		if oper := builtin[operator]; oper != nil {
			result = oper(node[1:])
			return result
		} else if oper := Env[operator]; oper != nil {
			symbol,ok := oper.(Func)
			checkFalse(ok)
			node[0] = symbol.node
			result = doEval(node)
		} else {
			log.Fatal("Undefined function")
		}
	}
		
		default : {
			log.Fatal("The first argument must be a function")
		}
	}
	return result
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

//type Parser struct {
//	Tokens []string
//	index  int
//}

// func GetTokens(expr string) []string {
// 	var tokens []string
// 	for i:=0; i<len(expr); i++ {
// 		c := string(expr[i])
// 		if c != " " && c != "\n" && c != "\t" && c != "\r" {
// 			var token string
// 			for ;i<len(expr); i++ {
// 				c := string(expr[i])
// 				if c == " " || c == "\n" || c == "\t" || c == "\r" {
// 					break
// 				} else {
// 					token += c
// 				}
// 			}
// 			tokens = append(tokens, token)
// 		}
// 	}
// 	return tokens
// }


//func (p *Parser) append(token string) {
//	fmt.Println("The token is: ", token)
//	p.Tokens = append(p.Tokens, token)
//}

func main () {
	if len(os.Args) != 2 {
		log.Fatal("Usage: gscheme <scheme file>")
	}

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
