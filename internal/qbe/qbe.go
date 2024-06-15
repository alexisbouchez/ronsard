package qbe

import (
	"fmt"
	"ronsard/internal/parser"
	"strings"
)

type QBEGenerator struct {
	ast   *parser.Node
	code  strings.Builder
	count int
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
	return g.code.String(), nil
}

func (g *QBEGenerator) emitPreamble() {
	g.code.WriteString(`
data $str = { b "Bonjour!", b 0 }
export function w $main() {
@start
`)
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
		g.code.WriteString(fmt.Sprintf("%%r%d =l add 0, $str\n", g.count))
		g.code.WriteString(fmt.Sprintf("%%u%d =w call $puts(l %%r%d)\n", g.count, g.count))
	default:
		return fmt.Errorf("unhandled node type: %v", node.Type)
	}
	return nil
}

func (g *QBEGenerator) emitPostamble() {
	g.code.WriteString("ret 0\n}\n")
}
