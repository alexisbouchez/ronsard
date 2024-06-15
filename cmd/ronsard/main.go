package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"ronsard/internal/lexer"
	"ronsard/internal/parser"
	"ronsard/internal/qbe"
)

func main() {
	// Input file path
	if len(os.Args) < 2 {
		fmt.Println("Usage: ronsard <source-file.rd>")
		return
	}
	sourceFilePath := os.Args[1]
	outputBinaryPath := "a.out"

	// Read the Ronsard source file
	sourceFile, err := os.ReadFile(sourceFilePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	sourceContents := string(sourceFile)

	// Lexing
	myLexer := lexer.NewLexer(sourceContents)
	var tokens []lexer.Token
	for {
		token := myLexer.NextToken()
		if token.Type == lexer.EOF {
			break
		}
		tokens = append(tokens, token)
	}

	// Parsing
	myParser := parser.NewParser(tokens)
	ast, err := myParser.Parse()
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}

	// QBE IL generation
	qbeGen := qbe.NewQBEGenerator(ast)
	qbeIL, err := qbeGen.Generate()
	if err != nil {
		fmt.Printf("Error generating QBE IL: %v\n", err)
		return
	}

	// Use a buffer to hold the assembly code produced by QBE
	var asmBuffer bytes.Buffer

	// Create a pipe to avoid intermediate files
	qbeCmd := exec.Command("qbe", "-")
	qbeCmd.Stdin = bytes.NewReader([]byte(qbeIL))
	qbeCmd.Stdout = &asmBuffer
	qbeCmd.Stderr = os.Stderr

	// Run QBE to generate assembly
	if err = qbeCmd.Run(); err != nil {
		fmt.Printf("Error running QBE: %v\n", err)
		return
	}

	// Create the cc command to compile the assembly to a binary
	ccCmd := exec.Command("cc", "-o", outputBinaryPath, "-x", "assembler", "-")
	ccCmd.Stdin = &asmBuffer
	ccCmd.Stdout = os.Stdout
	ccCmd.Stderr = os.Stderr

	// Run cc to produce the final binary
	if err = ccCmd.Run(); err != nil {
		fmt.Printf("Error running GCC: %v\n", err)
		return
	}

	fmt.Printf("Compilation successful! Executable: %s\n", outputBinaryPath)
}
