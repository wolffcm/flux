// DO NOT EDIT: This file is autogenerated via the builtin command.

package secrets

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
					Column: 12,
					Line:   3,
				},
				File:   "secrets.flux",
				Source: "package secrets\n\nbuiltin get",
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
						Column: 12,
						Line:   3,
					},
					File:   "secrets.flux",
					Source: "builtin get",
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
							Column: 12,
							Line:   3,
						},
						File:   "secrets.flux",
						Source: "get",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "get",
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=go",
		Name:     "secrets.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 16,
						Line:   1,
					},
					File:   "secrets.flux",
					Source: "package secrets",
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
							Column: 16,
							Line:   1,
						},
						File:   "secrets.flux",
						Source: "secrets",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "secrets",
			},
		},
	}},
	Package: "secrets",
	Path:    "influxdata/influxdb/secrets",
}
