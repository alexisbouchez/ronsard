package qbe_test

import (
	"os"
	"ronsard/internal/lexer"
	"ronsard/internal/parser"
	"ronsard/internal/qbe"
	"testing"
)

func TestQBEGenerator(t *testing.T) {
	exampleFilePath := "./../../examples/01_bonjour.rd"
	expectedQBE := `
data $str = { b "Bonjour!", b 0 }
export function w $main() {
@start
%r1 =l add 0, $str
%u1 =w call $puts(l %r1)
ret 0
}
`

	exampleFile, err := os.ReadFile(exampleFilePath)
	if err != nil {
		t.Fatal(err)
	}
	exampleFileContents := string(exampleFile)

	myLexer := lexer.NewLexer(exampleFileContents)

	var tokens []lexer.Token
	for {
		token := myLexer.NextToken()
		if token.Type == lexer.EOF {
			break
		}
		tokens = append(tokens, token)
	}

	myParser := parser.NewParser(tokens)
	ast, err := myParser.Parse()
	if err != nil {
		t.Fatalf("Error parsing input: %v", err)
	}

	qbeGen := qbe.NewQBEGenerator(ast)
	qbeIL, err := qbeGen.Generate()
	if err != nil {
		t.Fatalf("Error generating QBE IL: %v", err)
	}

	if qbeIL != expectedQBE {
		t.Errorf("Generated QBE IL does not match expected output.\nExpected:\n%s\nGot:\n%s", expectedQBE, qbeIL)
	}
}
