package main

import(
	"fmt"
	"github.com/itsmecerealdev/itan-go/lib/token"
	"github.com/itsmecerealdev/itan-go/lib/node"
	"github.com/itsmecerealdev/itan-go/lib/parser"
)

func main() {
	buffer := "1234567 + 123456789; int32 x = 5; x = 5;"
	tokens, err := token.Tokenize(buffer)
	if (err != nil) {
		fmt.Printf("Error: %s\n", err);
		return
	}
	for _, t := range tokens {
		fmt.Printf("Token: %s, Name: %s\n", t.Type, t.Name)
	}
	var root node.ProgramNode = parser.ParseProgram(tokens)
	fmt.Println(root)
	root.Visit()
}
