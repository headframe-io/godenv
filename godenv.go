package godenv

import (
	"fmt"
	"io"

	"github.com/headframe-io/godenv/internal/ast"
	"github.com/headframe-io/godenv/internal/parser"
	"github.com/headframe-io/godenv/internal/scanner"
)

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader) (map[string]string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	s := scanner.New(string(input))
	p := parser.New(s)

	statement, err := p.Parse()
	if err != nil {
		return nil, err
	}

	fileStmt, ok := statement.(*ast.FileStatement)
	if !ok {
		return nil, fmt.Errorf("unexpected statement: %T", statement)
	}

	values := make(map[string]string, len(fileStmt.Statements))

	for _, stmt := range fileStmt.Statements {
		if assign, ok := stmt.(*ast.AssignStatement); ok {
			values[assign.Name] = assign.Value
		}
	}

	return values, nil
}
