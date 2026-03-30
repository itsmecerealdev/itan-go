package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/itsmecerealdev/itan-go/lib/node"
	"github.com/itsmecerealdev/itan-go/lib/parser"
	"github.com/itsmecerealdev/itan-go/lib/token"
)

func main() {
	// buffer := "1234567 + 123456789; int32 x = 5; x = 5;"
	var tokens []token.Token
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		tokens = tokenizeByLine(filePath)
	} else {
		log.Fatal("compiler requires file name (go run cmd/compiler.go filename). include the extension (filename.it) to compile")
	}
	// tokens, err := token.Tokenize(buffer)
	// if (err != nil) {
		// fmt.Printf("Error: %s\n", err);
		// return
	// }
	for _, t := range tokens {
		fmt.Printf("Token: %s, Name: %s, Line: %d, Col: %d\n", t.Type, t.Name, t.Line, t.Col)
	}
	var root node.ProgramNode = parser.ParseProgram(tokens)
	fmt.Println(root)
	root.Visit()
}

func tokenizeByLine(filePath string) []token.Token {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var retTokens []token.Token
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		tempTokens, err := token.Tokenize(line, lineNumber)
		if err != nil {
			log.Fatal(err)
		}
		retTokens = append(retTokens, tempTokens...)
		lineNumber++;
	}
	return retTokens
}
