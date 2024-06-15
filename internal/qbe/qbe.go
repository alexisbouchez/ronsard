package qbe

import (
	"fmt"
	"ronsard/internal/parser"
)

type QBEGenerator struct {
	ast       *parser.Node
	code      string
	count     int
	dataCount int
}

func NewQBEGenerator(ast *parser.Node) *QBEGenerator {
	return &QBEGenerator{
		ast: ast,
	}
}

func (g *QBEGenerator) Generate() (string, error) {
	g.emitPreamble()
	err := g.emitNode(g.ast)
	if err != nil {
		return "", err
	}
	g.emitPostamble() // Make sure to call this to add the closing part
	return g.code, nil
}

func (g *QBEGenerator) emitPreamble() {
	g.code += `
export function w $main() {
@start
`
}

func (g *QBEGenerator) emitNode(node *parser.Node) error {
	switch node.Type {
	case parser.PRINT_STMT:
		for _, child := range node.Children {
			if err := g.emitNode(child); err != nil {
				return err
			}
		}
	case parser.STRING_EXPR:
		g.count++
		g.code = fmt.Sprintf("data $data_%d = { b %s, b 0 }\n", g.dataCount, node.Value) + g.code
		g.code += fmt.Sprintf("%%r%d =l add 0, $data_%d\n", g.count, g.dataCount)
		g.code += fmt.Sprintf("%%u%d =w call $puts(l %%r%d)\n", g.count, g.count)
	default:
		return fmt.Errorf("unhandled node type: %v", node.Type)
	}
	return nil
}

func (g *QBEGenerator) emitPostamble() {
	g.code += "ret 0\n}\n"
}
