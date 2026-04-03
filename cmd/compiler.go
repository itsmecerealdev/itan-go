package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/itsmecerealdev/itan-go/lib/node"
	"github.com/itsmecerealdev/itan-go/lib/parser"
	"github.com/itsmecerealdev/itan-go/lib/token"
)

func main() {
	var tokens []token.Token
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		tokens = tokenizeByLine(filePath)
		tokens = append(tokens, token.Token{
			Name: "End",
			Type: token.End,
		})
	} else {
		log.Fatal("compiler requires file name (go run cmd/compiler.go filename). include the extension (filename.it) to compile")
	}
	for _, t := range tokens {
		fmt.Printf("Token: %s, Name: %s, Line: %d, Col: %d\n", t.Type, t.Name, t.Line, t.Col)
	}
	root := parser.ParseProgram(tokens)
	printer := node.Printer{}
	declarator := node.Declarator{}
	semantic := node.Semantic{}
	root.Accept(&printer)
	root.Accept(&declarator)
	root.Accept(&semantic)
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
		if strings.TrimSpace(line) == "" {
			lineNumber++
			continue
		}
		if checkIfComment(line) {
			lineNumber++
			continue
		}
		tempTokens, err := token.Tokenize(line, lineNumber)
		if err != nil {
			log.Fatal(err)
		}
		retTokens = append(retTokens, tempTokens...)
		lineNumber++;
	}
	return retTokens
}

func checkIfComment(line string) bool {
	count := 0
	res := ""
	for _, c := range line {
		if !unicode.IsSpace(rune(c)) {
			count++
			res += string(c)
		}
		if count >= 2 {
			break
		}
	}
	return res == "//"
}
