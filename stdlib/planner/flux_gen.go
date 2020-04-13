// DO NOT EDIT: This file is autogenerated via the builtin command.

package planner

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
					Column: 35,
					Line:   4,
				},
				File:   "planner.flux",
				Source: "package planner\n\noption disableLogicalRules = [\"\"]\noption disablePhysicalRules = [\"\"]",
				Start: ast.Position{
					Column: 1,
					Line:   1,
				},
			},
		},
		Body: []ast.Statement{&ast.OptionStatement{
			Assignment: &ast.VariableAssignment{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 34,
							Line:   3,
						},
						File:   "planner.flux",
						Source: "disableLogicalRules = [\"\"]",
						Start: ast.Position{
							Column: 8,
							Line:   3,
						},
					},
				},
				ID: &ast.Identifier{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 27,
								Line:   3,
							},
							File:   "planner.flux",
							Source: "disableLogicalRules",
							Start: ast.Position{
								Column: 8,
								Line:   3,
							},
						},
					},
					Name: "disableLogicalRules",
				},
				Init: &ast.ArrayExpression{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 34,
								Line:   3,
							},
							File:   "planner.flux",
							Source: "[\"\"]",
							Start: ast.Position{
								Column: 30,
								Line:   3,
							},
						},
					},
					Elements: []ast.Expression{&ast.StringLiteral{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 33,
									Line:   3,
								},
								File:   "planner.flux",
								Source: "\"\"",
								Start: ast.Position{
									Column: 31,
									Line:   3,
								},
							},
						},
						Value: "",
					}},
				},
			},
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 34,
						Line:   3,
					},
					File:   "planner.flux",
					Source: "option disableLogicalRules = [\"\"]",
					Start: ast.Position{
						Column: 1,
						Line:   3,
					},
				},
			},
		}, &ast.OptionStatement{
			Assignment: &ast.VariableAssignment{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 35,
							Line:   4,
						},
						File:   "planner.flux",
						Source: "disablePhysicalRules = [\"\"]",
						Start: ast.Position{
							Column: 8,
							Line:   4,
						},
					},
				},
				ID: &ast.Identifier{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 28,
								Line:   4,
							},
							File:   "planner.flux",
							Source: "disablePhysicalRules",
							Start: ast.Position{
								Column: 8,
								Line:   4,
							},
						},
					},
					Name: "disablePhysicalRules",
				},
				Init: &ast.ArrayExpression{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 35,
								Line:   4,
							},
							File:   "planner.flux",
							Source: "[\"\"]",
							Start: ast.Position{
								Column: 31,
								Line:   4,
							},
						},
					},
					Elements: []ast.Expression{&ast.StringLiteral{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 34,
									Line:   4,
								},
								File:   "planner.flux",
								Source: "\"\"",
								Start: ast.Position{
									Column: 32,
									Line:   4,
								},
							},
						},
						Value: "",
					}},
				},
			},
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 35,
						Line:   4,
					},
					File:   "planner.flux",
					Source: "option disablePhysicalRules = [\"\"]",
					Start: ast.Position{
						Column: 1,
						Line:   4,
					},
				},
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=rust",
		Name:     "planner.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 16,
						Line:   1,
					},
					File:   "planner.flux",
					Source: "package planner",
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
						File:   "planner.flux",
						Source: "planner",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "planner",
			},
		},
	}},
	Package: "planner",
	Path:    "planner",
}
