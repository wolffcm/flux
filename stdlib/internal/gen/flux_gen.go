// DO NOT EDIT: This file is autogenerated via the builtin command.

package gen

import (
	ast "github.com/wolffcm/flux/ast"
	runtime "github.com/wolffcm/flux/runtime"
)

func init() {
	runtime.RegisterPackage(pkgAST)
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
					Column: 15,
					Line:   3,
				},
				File:   "gen.flux",
				Source: "package gen\n\nbuiltin tables",
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
						Column: 15,
						Line:   3,
					},
					File:   "gen.flux",
					Source: "builtin tables",
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
							Column: 15,
							Line:   3,
						},
						File:   "gen.flux",
						Source: "tables",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "tables",
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=rust",
		Name:     "gen.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 12,
						Line:   1,
					},
					File:   "gen.flux",
					Source: "package gen",
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
							Column: 12,
							Line:   1,
						},
						File:   "gen.flux",
						Source: "gen",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "gen",
			},
		},
	}},
	Package: "gen",
	Path:    "internal/gen",
}
