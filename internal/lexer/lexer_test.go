package lexer_test

import (
	"os"
	"ronsard/internal/lexer"
	"testing"
)

func TestLexer(t *testing.T) {
	testExampleAgainstLexer(t, "./../../examples/01_bonjour.rd", []lexer.Token{
		{Type: lexer.KEYWORD, Value: "afficher"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.STRING_LITERAL, Value: `"Bonjour!"`},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.SEMICOLON, Value: ";"},
	})
}

func testExampleAgainstLexer(
	t *testing.T,
	exampleFilePath string,
	expectedTokens []lexer.Token,
) {
	exampleFile, err := os.ReadFile(exampleFilePath)
	if err != nil {
		t.Fatal(err)
	}
	exampleFileContents := string(exampleFile)

	myLexer := lexer.NewLexer(exampleFileContents)

	for i, expected := range expectedTokens {
		token := myLexer.NextToken()
		if token.Type != expected.Type || token.Value != expected.Value {
			t.Errorf("Test %d: expected token %+v, got %+v", i+1, expected, token)
		}
	}

	token := myLexer.NextToken()
	if token.Type != lexer.EOF {
		t.Errorf("Expected EOF token, got %+v", token)
	}
}
