// DO NOT EDIT: This file is autogenerated via the builtin command.

package generate

import (
	flux "github.com/wolffcm/flux"
	ast "github.com/wolffcm/flux/ast"
)

func init() {
	flux.RegisterPackage(pkgAST)
}

var pkgAST = &ast.Package{
	BaseNode: ast.BaseNode{
		Errors: nil,
		Loc:    nil,
	},
	Files: []*ast.File{&ast.File{
		BaseNode: ast.BaseNode{
			Errors: nil,
			Loc: &ast.SourceLocation{
				End: ast.Position{
					Column: 13,
					Line:   3,
				},
				File:   "generate.flux",
				Source: "package generate\n\nbuiltin from",
				Start: ast.Position{
					Column: 1,
					Line:   1,
				},
			},
		},
		Body: []ast.Statement{&ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   3,
					},
					File:   "generate.flux",
					Source: "builtin from",
					Start: ast.Position{
						Column: 1,
						Line:   3,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 13,
							Line:   3,
						},
						File:   "generate.flux",
						Source: "from",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "from",
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=go",
		Name:     "generate.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 17,
						Line:   1,
					},
					File:   "generate.flux",
					Source: "package generate",
					Start: ast.Position{
						Column: 1,
						Line:   1,
					},
				},
			},
			Name: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 17,
							Line:   1,
						},
						File:   "generate.flux",
						Source: "generate",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "generate",
			},
		},
	}},
	Package: "generate",
	Path:    "generate",
}
