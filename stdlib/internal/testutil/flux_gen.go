// DO NOT EDIT: This file is autogenerated via the builtin command.

package testutil

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
					Column: 19,
					Line:   5,
				},
				File:   "testutil.flux",
				Source: "package testutil\n\nbuiltin fail\nbuiltin yield\nbuiltin makeRecord",
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
					File:   "testutil.flux",
					Source: "builtin fail",
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
						File:   "testutil.flux",
						Source: "fail",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "fail",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   4,
					},
					File:   "testutil.flux",
					Source: "builtin yield",
					Start: ast.Position{
						Column: 1,
						Line:   4,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   4,
						},
						File:   "testutil.flux",
						Source: "yield",
						Start: ast.Position{
							Column: 9,
							Line:   4,
						},
					},
				},
				Name: "yield",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 19,
						Line:   5,
					},
					File:   "testutil.flux",
					Source: "builtin makeRecord",
					Start: ast.Position{
						Column: 1,
						Line:   5,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 19,
							Line:   5,
						},
						File:   "testutil.flux",
						Source: "makeRecord",
						Start: ast.Position{
							Column: 9,
							Line:   5,
						},
					},
				},
				Name: "makeRecord",
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=rust",
		Name:     "testutil.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 17,
						Line:   1,
					},
					File:   "testutil.flux",
					Source: "package testutil",
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
						File:   "testutil.flux",
						Source: "testutil",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "testutil",
			},
		},
	}},
	Package: "testutil",
	Path:    "internal/testutil",
}