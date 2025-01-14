package parsing

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Schematic struct {
	Elements []*Element `@@*`
}

type Element struct {
	Type       string     `@Ident`
	Properties []Property `('(' (@@ (',' @@)*)? ')')?`
	Actions    []Action   `( @@+ )?`
}

type Property struct {
	Key   string `@Ident "="`
	Value string `@String`
}

type Action struct {
	Type  string  `'.' @Ident`
	Units float64 `('(' @Number? ')')?`
}

var (
	schemGoLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Ident", `[a-zA-Z_][a-zA-Z_0-9]*`},
		{"String", `"[^"]*"`},
		{"Number", `[-+]?[.0-9]+\b`},
		{"Punct", `\[|]|[-!()+/*=,]`},
		{"comment", `#[^\n]+`},
		{"whitespace", `\s+`},
	})
	schemGoParser = participle.MustBuild[Schematic](
		participle.Lexer(schemGoLexer),
		participle.Unquote("String"),
	)
)

func ReadSchematic(schematicFilePath string) (schematic *Schematic, err error) {

	schemFile, err := os.Open(schematicFilePath)
	if err != nil {
		panic(err)
	}

	return schemGoParser.Parse(schematicFilePath, schemFile)
}
