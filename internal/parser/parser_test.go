package parser_test

import (
	"os"
	"ronsard/internal/lexer"
	"ronsard/internal/parser"
	"testing"
)

func TestParser(t *testing.T) {
	testExampleAgainstParser(t, "./../../examples/01_bonjour.rd", &parser.Node{
		Type: parser.PRINT_STMT,
		Children: []*parser.Node{
			{
				Type:  parser.STRING_EXPR,
				Value: `"Bonjour!"`,
			},
		},
	})
}

func testExampleAgainstParser(
	t *testing.T,
	exampleFilePath string,
	expectedAST *parser.Node,
) {
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

	if !compareNodes(ast, expectedAST) {
		t.Errorf("AST does not match expected structure. Got %+v, expected %+v", ast, expectedAST)
	}
}

// compareNodes compares two AST nodes recursively.
func compareNodes(a, b *parser.Node) bool {
	if a.Type != b.Type || a.Value != b.Value {
		return false
	}
	if len(a.Children) != len(b.Children) {
		return false
	}
	for i := range a.Children {
		if !compareNodes(a.Children[i], b.Children[i]) {
			return false
		}
	}
	return true
}
