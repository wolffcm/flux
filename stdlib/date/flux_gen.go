// DO NOT EDIT: This file is autogenerated via the builtin command.

package date

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
					Column: 15,
					Line:   37,
				},
				File:   "date.flux",
				Source: "package date\n\nbuiltin second\nbuiltin minute\nbuiltin hour\nbuiltin weekDay\nbuiltin monthDay\nbuiltin yearDay\nbuiltin month\nbuiltin year\nbuiltin week\nbuiltin quarter\nbuiltin millisecond\nbuiltin microsecond\nbuiltin nanosecond\nbuiltin truncate\n\nSunday    = 0\nMonday    = 1\nTuesday   = 2\nWednesday = 3\nThursday  = 4\nFriday    = 5\nSaturday  = 6\n\nJanuary   = 1\nFebruary  = 2\nMarch     = 3\nApril     = 4\nMay       = 5\nJune      = 6\nJuly      = 7\nAugust    = 8\nSeptember = 9\nOctober   = 10\nNovember  = 11\nDecember  = 12",
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
					File:   "date.flux",
					Source: "builtin second",
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
						File:   "date.flux",
						Source: "second",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "second",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 15,
						Line:   4,
					},
					File:   "date.flux",
					Source: "builtin minute",
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
							Column: 15,
							Line:   4,
						},
						File:   "date.flux",
						Source: "minute",
						Start: ast.Position{
							Column: 9,
							Line:   4,
						},
					},
				},
				Name: "minute",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   5,
					},
					File:   "date.flux",
					Source: "builtin hour",
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
							Column: 13,
							Line:   5,
						},
						File:   "date.flux",
						Source: "hour",
						Start: ast.Position{
							Column: 9,
							Line:   5,
						},
					},
				},
				Name: "hour",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 16,
						Line:   6,
					},
					File:   "date.flux",
					Source: "builtin weekDay",
					Start: ast.Position{
						Column: 1,
						Line:   6,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 16,
							Line:   6,
						},
						File:   "date.flux",
						Source: "weekDay",
						Start: ast.Position{
							Column: 9,
							Line:   6,
						},
					},
				},
				Name: "weekDay",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 17,
						Line:   7,
					},
					File:   "date.flux",
					Source: "builtin monthDay",
					Start: ast.Position{
						Column: 1,
						Line:   7,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 17,
							Line:   7,
						},
						File:   "date.flux",
						Source: "monthDay",
						Start: ast.Position{
							Column: 9,
							Line:   7,
						},
					},
				},
				Name: "monthDay",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 16,
						Line:   8,
					},
					File:   "date.flux",
					Source: "builtin yearDay",
					Start: ast.Position{
						Column: 1,
						Line:   8,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 16,
							Line:   8,
						},
						File:   "date.flux",
						Source: "yearDay",
						Start: ast.Position{
							Column: 9,
							Line:   8,
						},
					},
				},
				Name: "yearDay",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   9,
					},
					File:   "date.flux",
					Source: "builtin month",
					Start: ast.Position{
						Column: 1,
						Line:   9,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   9,
						},
						File:   "date.flux",
						Source: "month",
						Start: ast.Position{
							Column: 9,
							Line:   9,
						},
					},
				},
				Name: "month",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   10,
					},
					File:   "date.flux",
					Source: "builtin year",
					Start: ast.Position{
						Column: 1,
						Line:   10,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 13,
							Line:   10,
						},
						File:   "date.flux",
						Source: "year",
						Start: ast.Position{
							Column: 9,
							Line:   10,
						},
					},
				},
				Name: "year",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   11,
					},
					File:   "date.flux",
					Source: "builtin week",
					Start: ast.Position{
						Column: 1,
						Line:   11,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 13,
							Line:   11,
						},
						File:   "date.flux",
						Source: "week",
						Start: ast.Position{
							Column: 9,
							Line:   11,
						},
					},
				},
				Name: "week",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 16,
						Line:   12,
					},
					File:   "date.flux",
					Source: "builtin quarter",
					Start: ast.Position{
						Column: 1,
						Line:   12,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 16,
							Line:   12,
						},
						File:   "date.flux",
						Source: "quarter",
						Start: ast.Position{
							Column: 9,
							Line:   12,
						},
					},
				},
				Name: "quarter",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 20,
						Line:   13,
					},
					File:   "date.flux",
					Source: "builtin millisecond",
					Start: ast.Position{
						Column: 1,
						Line:   13,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 20,
							Line:   13,
						},
						File:   "date.flux",
						Source: "millisecond",
						Start: ast.Position{
							Column: 9,
							Line:   13,
						},
					},
				},
				Name: "millisecond",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 20,
						Line:   14,
					},
					File:   "date.flux",
					Source: "builtin microsecond",
					Start: ast.Position{
						Column: 1,
						Line:   14,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 20,
							Line:   14,
						},
						File:   "date.flux",
						Source: "microsecond",
						Start: ast.Position{
							Column: 9,
							Line:   14,
						},
					},
				},
				Name: "microsecond",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 19,
						Line:   15,
					},
					File:   "date.flux",
					Source: "builtin nanosecond",
					Start: ast.Position{
						Column: 1,
						Line:   15,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 19,
							Line:   15,
						},
						File:   "date.flux",
						Source: "nanosecond",
						Start: ast.Position{
							Column: 9,
							Line:   15,
						},
					},
				},
				Name: "nanosecond",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 17,
						Line:   16,
					},
					File:   "date.flux",
					Source: "builtin truncate",
					Start: ast.Position{
						Column: 1,
						Line:   16,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 17,
							Line:   16,
						},
						File:   "date.flux",
						Source: "truncate",
						Start: ast.Position{
							Column: 9,
							Line:   16,
						},
					},
				},
				Name: "truncate",
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   18,
					},
					File:   "date.flux",
					Source: "Sunday    = 0",
					Start: ast.Position{
						Column: 1,
						Line:   18,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 7,
							Line:   18,
						},
						File:   "date.flux",
						Source: "Sunday",
						Start: ast.Position{
							Column: 1,
							Line:   18,
						},
					},
				},
				Name: "Sunday",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   18,
						},
						File:   "date.flux",
						Source: "0",
						Start: ast.Position{
							Column: 13,
							Line:   18,
						},
					},
				},
				Value: int64(0),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   19,
					},
					File:   "date.flux",
					Source: "Monday    = 1",
					Start: ast.Position{
						Column: 1,
						Line:   19,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 7,
							Line:   19,
						},
						File:   "date.flux",
						Source: "Monday",
						Start: ast.Position{
							Column: 1,
							Line:   19,
						},
					},
				},
				Name: "Monday",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   19,
						},
						File:   "date.flux",
						Source: "1",
						Start: ast.Position{
							Column: 13,
							Line:   19,
						},
					},
				},
				Value: int64(1),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   20,
					},
					File:   "date.flux",
					Source: "Tuesday   = 2",
					Start: ast.Position{
						Column: 1,
						Line:   20,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 8,
							Line:   20,
						},
						File:   "date.flux",
						Source: "Tuesday",
						Start: ast.Position{
							Column: 1,
							Line:   20,
						},
					},
				},
				Name: "Tuesday",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   20,
						},
						File:   "date.flux",
						Source: "2",
						Start: ast.Position{
							Column: 13,
							Line:   20,
						},
					},
				},
				Value: int64(2),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   21,
					},
					File:   "date.flux",
					Source: "Wednesday = 3",
					Start: ast.Position{
						Column: 1,
						Line:   21,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 10,
							Line:   21,
						},
						File:   "date.flux",
						Source: "Wednesday",
						Start: ast.Position{
							Column: 1,
							Line:   21,
						},
					},
				},
				Name: "Wednesday",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   21,
						},
						File:   "date.flux",
						Source: "3",
						Start: ast.Position{
							Column: 13,
							Line:   21,
						},
					},
				},
				Value: int64(3),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   22,
					},
					File:   "date.flux",
					Source: "Thursday  = 4",
					Start: ast.Position{
						Column: 1,
						Line:   22,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 9,
							Line:   22,
						},
						File:   "date.flux",
						Source: "Thursday",
						Start: ast.Position{
							Column: 1,
							Line:   22,
						},
					},
				},
				Name: "Thursday",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   22,
						},
						File:   "date.flux",
						Source: "4",
						Start: ast.Position{
							Column: 13,
							Line:   22,
						},
					},
				},
				Value: int64(4),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   23,
					},
					File:   "date.flux",
					Source: "Friday    = 5",
					Start: ast.Position{
						Column: 1,
						Line:   23,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 7,
							Line:   23,
						},
						File:   "date.flux",
						Source: "Friday",
						Start: ast.Position{
							Column: 1,
							Line:   23,
						},
					},
				},
				Name: "Friday",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   23,
						},
						File:   "date.flux",
						Source: "5",
						Start: ast.Position{
							Column: 13,
							Line:   23,
						},
					},
				},
				Value: int64(5),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   24,
					},
					File:   "date.flux",
					Source: "Saturday  = 6",
					Start: ast.Position{
						Column: 1,
						Line:   24,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 9,
							Line:   24,
						},
						File:   "date.flux",
						Source: "Saturday",
						Start: ast.Position{
							Column: 1,
							Line:   24,
						},
					},
				},
				Name: "Saturday",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   24,
						},
						File:   "date.flux",
						Source: "6",
						Start: ast.Position{
							Column: 13,
							Line:   24,
						},
					},
				},
				Value: int64(6),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   26,
					},
					File:   "date.flux",
					Source: "January   = 1",
					Start: ast.Position{
						Column: 1,
						Line:   26,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 8,
							Line:   26,
						},
						File:   "date.flux",
						Source: "January",
						Start: ast.Position{
							Column: 1,
							Line:   26,
						},
					},
				},
				Name: "January",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   26,
						},
						File:   "date.flux",
						Source: "1",
						Start: ast.Position{
							Column: 13,
							Line:   26,
						},
					},
				},
				Value: int64(1),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   27,
					},
					File:   "date.flux",
					Source: "February  = 2",
					Start: ast.Position{
						Column: 1,
						Line:   27,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 9,
							Line:   27,
						},
						File:   "date.flux",
						Source: "February",
						Start: ast.Position{
							Column: 1,
							Line:   27,
						},
					},
				},
				Name: "February",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   27,
						},
						File:   "date.flux",
						Source: "2",
						Start: ast.Position{
							Column: 13,
							Line:   27,
						},
					},
				},
				Value: int64(2),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   28,
					},
					File:   "date.flux",
					Source: "March     = 3",
					Start: ast.Position{
						Column: 1,
						Line:   28,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 6,
							Line:   28,
						},
						File:   "date.flux",
						Source: "March",
						Start: ast.Position{
							Column: 1,
							Line:   28,
						},
					},
				},
				Name: "March",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   28,
						},
						File:   "date.flux",
						Source: "3",
						Start: ast.Position{
							Column: 13,
							Line:   28,
						},
					},
				},
				Value: int64(3),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   29,
					},
					File:   "date.flux",
					Source: "April     = 4",
					Start: ast.Position{
						Column: 1,
						Line:   29,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 6,
							Line:   29,
						},
						File:   "date.flux",
						Source: "April",
						Start: ast.Position{
							Column: 1,
							Line:   29,
						},
					},
				},
				Name: "April",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   29,
						},
						File:   "date.flux",
						Source: "4",
						Start: ast.Position{
							Column: 13,
							Line:   29,
						},
					},
				},
				Value: int64(4),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   30,
					},
					File:   "date.flux",
					Source: "May       = 5",
					Start: ast.Position{
						Column: 1,
						Line:   30,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 4,
							Line:   30,
						},
						File:   "date.flux",
						Source: "May",
						Start: ast.Position{
							Column: 1,
							Line:   30,
						},
					},
				},
				Name: "May",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   30,
						},
						File:   "date.flux",
						Source: "5",
						Start: ast.Position{
							Column: 13,
							Line:   30,
						},
					},
				},
				Value: int64(5),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   31,
					},
					File:   "date.flux",
					Source: "June      = 6",
					Start: ast.Position{
						Column: 1,
						Line:   31,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 5,
							Line:   31,
						},
						File:   "date.flux",
						Source: "June",
						Start: ast.Position{
							Column: 1,
							Line:   31,
						},
					},
				},
				Name: "June",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   31,
						},
						File:   "date.flux",
						Source: "6",
						Start: ast.Position{
							Column: 13,
							Line:   31,
						},
					},
				},
				Value: int64(6),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   32,
					},
					File:   "date.flux",
					Source: "July      = 7",
					Start: ast.Position{
						Column: 1,
						Line:   32,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 5,
							Line:   32,
						},
						File:   "date.flux",
						Source: "July",
						Start: ast.Position{
							Column: 1,
							Line:   32,
						},
					},
				},
				Name: "July",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   32,
						},
						File:   "date.flux",
						Source: "7",
						Start: ast.Position{
							Column: 13,
							Line:   32,
						},
					},
				},
				Value: int64(7),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   33,
					},
					File:   "date.flux",
					Source: "August    = 8",
					Start: ast.Position{
						Column: 1,
						Line:   33,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 7,
							Line:   33,
						},
						File:   "date.flux",
						Source: "August",
						Start: ast.Position{
							Column: 1,
							Line:   33,
						},
					},
				},
				Name: "August",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   33,
						},
						File:   "date.flux",
						Source: "8",
						Start: ast.Position{
							Column: 13,
							Line:   33,
						},
					},
				},
				Value: int64(8),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   34,
					},
					File:   "date.flux",
					Source: "September = 9",
					Start: ast.Position{
						Column: 1,
						Line:   34,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 10,
							Line:   34,
						},
						File:   "date.flux",
						Source: "September",
						Start: ast.Position{
							Column: 1,
							Line:   34,
						},
					},
				},
				Name: "September",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   34,
						},
						File:   "date.flux",
						Source: "9",
						Start: ast.Position{
							Column: 13,
							Line:   34,
						},
					},
				},
				Value: int64(9),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 15,
						Line:   35,
					},
					File:   "date.flux",
					Source: "October   = 10",
					Start: ast.Position{
						Column: 1,
						Line:   35,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 8,
							Line:   35,
						},
						File:   "date.flux",
						Source: "October",
						Start: ast.Position{
							Column: 1,
							Line:   35,
						},
					},
				},
				Name: "October",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 15,
							Line:   35,
						},
						File:   "date.flux",
						Source: "10",
						Start: ast.Position{
							Column: 13,
							Line:   35,
						},
					},
				},
				Value: int64(10),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 15,
						Line:   36,
					},
					File:   "date.flux",
					Source: "November  = 11",
					Start: ast.Position{
						Column: 1,
						Line:   36,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 9,
							Line:   36,
						},
						File:   "date.flux",
						Source: "November",
						Start: ast.Position{
							Column: 1,
							Line:   36,
						},
					},
				},
				Name: "November",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 15,
							Line:   36,
						},
						File:   "date.flux",
						Source: "11",
						Start: ast.Position{
							Column: 13,
							Line:   36,
						},
					},
				},
				Value: int64(11),
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 15,
						Line:   37,
					},
					File:   "date.flux",
					Source: "December  = 12",
					Start: ast.Position{
						Column: 1,
						Line:   37,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 9,
							Line:   37,
						},
						File:   "date.flux",
						Source: "December",
						Start: ast.Position{
							Column: 1,
							Line:   37,
						},
					},
				},
				Name: "December",
			},
			Init: &ast.IntegerLiteral{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 15,
							Line:   37,
						},
						File:   "date.flux",
						Source: "12",
						Start: ast.Position{
							Column: 13,
							Line:   37,
						},
					},
				},
				Value: int64(12),
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=go",
		Name:     "date.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   1,
					},
					File:   "date.flux",
					Source: "package date",
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
							Column: 13,
							Line:   1,
						},
						File:   "date.flux",
						Source: "date",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "date",
			},
		},
	}},
	Package: "date",
	Path:    "date",
}
