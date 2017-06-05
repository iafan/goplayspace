package help

var langspecRoot = "https://golang.org/ref/spec#"

var langspec = map[string]string{

	// Keywords

	"break":       "Break_statements",
	"case":        "Type_switches",
	"chan":        "Channel_types",
	"const":       "Constant_declarations",
	"continue":    "Continue_statements",
	"default":     "Expression_switches",
	"defer":       "Defer_statements",
	"else":        "If_statements",
	"fallthrough": "Fallthrough_statements",
	"for":         "For_statements",
	"func":        "Function_declarations",
	"go":          "Go_statements",
	"goto":        "Goto_statements",
	"if":          "If_statements",
	"import":      "Import_declarations",
	"interface":   "Interface_types",
	"map":         "Map_types",
	"package":     "Package_clause",
	"range":       "For_range",
	"return":      "Return_statements",
	"select":      "Select_statements",
	"struct":      "Struct_types",
	"switch":      "Switch_statements",
	"type":        "Types",
	"var":         "Variable_declarations",

	// Predeclared identifiers

	// - Types
	"bool":   "Boolean_types",
	"error":  "Errors",
	"string": "String_types",

	"byte":       "Numeric_types",
	"complex64":  "Numeric_types",
	"complex128": "Numeric_types",
	"float32":    "Numeric_types",
	"float64":    "Numeric_types",
	"int":        "Numeric_types",
	"int8":       "Numeric_types",
	"int16":      "Numeric_types",
	"int32":      "Numeric_types",
	"int64":      "Numeric_types",
	"rune":       "Numeric_types",
	"uint":       "Numeric_types",
	"uint8":      "Numeric_types",
	"uint16":     "Numeric_types",
	"uint32":     "Numeric_types",
	"uint64":     "Numeric_types",
	"uintptr":    "Numeric_types",

	// - Constants
	"false": "Boolean_types",
	"iota":  "Iota",
	"true":  "Boolean_types",

	// - Zero value
	"nil": "Pointer_types",

	// - Functions
	"append":  "Appending_and_copying_slices",
	"cap":     "Length_and_capacity",
	"close":   "Close",
	"complex": "Complex_numbers",
	"copy":    "Appending_and_copying_slices",
	"delete":  "Deletion_of_map_elements",
	"imag":    "Complex_numbers",
	"len":     "Length_and_capacity",
	"make":    "Making_slices_maps_and_channels",
	"new":     "Allocation",
	"panic":   "Handling_panics",
	"print":   "Bootstrapping",
	"println": "Bootstrapping",
	"real":    "Complex_numbers",
	"recover": "Handling_panics",

	// Other statements/operators

	"...": "Passing_arguments_to_..._parameters",
	"[]":  "Slice_types",
	"{}":  "Interface_types",
	".":   "Qualified_identifiers",
	":=":  "Short_variable_declarations",
	";":   "Semicolons",
	"'":   "Rune_literals",
	`"`:   "String_literals",
	"`":   "String_literals",
	"{":   "Blocks",
	"}":   "Blocks",
	"(":   "Notation",
	")":   "Notation",
	"//":  "Comments",
	"/*":  "Comments",
	"*/":  "Comments",
	"++":  "IncDec_statements",
	"--":  "IncDec_statements",
	"=":   "Assignments",
	"<-":  "Send_statements", // also "Receive_operator"
	"_":   "Blank_identifier",

	`\a`: "Rune_literals",
	`\b`: "Rune_literals",
	`\f`: "Rune_literals",
	`\n`: "Rune_literals",
	`\r`: "Rune_literals",
	`\t`: "Rune_literals",
	`\v`: "Rune_literals",
	`\\`: "Rune_literals",
	`\'`: "Rune_literals",
	`\"`: "Rune_literals",

	"||": "Logical_operators",
	"&&": "Logical_operators",
	"!":  "Logical_operators",

	"*": "Address_operators",
	"&": "Address_operators",

	"+": "Arithmetic_operators",
	"-": "Arithmetic_operators",
	//"*": //"Arithmetic_operators",
	"/": "Arithmetic_operators",
	"%": "Arithmetic_operators",
	//"&": //"Arithmetic_operators",
	"|":  "Arithmetic_operators",
	"^":  "Arithmetic_operators",
	"&^": "Arithmetic_operators",

	"<<": "Arithmetic_operators",
	">>": "Arithmetic_operators",

	"==": "Comparison_operators",
	"!=": "Comparison_operators",
	"<":  "Comparison_operators",
	"<=": "Comparison_operators",
	">":  "Comparison_operators",
	">=": "Comparison_operators",
}
