package filter

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Constants for property value types
const (
	StringValueType = "string_value"
	DoubleValueType = "double_value"
	IntValueType    = "int_value"
	BoolValueType   = "bool_value"
	ArrayValueType  = "array_value"
)

// Define the lexer for SQL WHERE clauses
var sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "whitespace", Pattern: `\s+`},
	{Name: "Comment", Pattern: `--[^\r\n]*`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Float", Pattern: `[-+]?\d*\.\d+([eE][-+]?\d+)?|[-+]?\d+[eE][-+]?\d+`},
	{Name: "Int", Pattern: `[-+]?\d+`},
	{Name: "String", Pattern: `'([^'\\]|\\.)*'|"([^"\\]|\\.)*"`},
	{Name: "EscapedIdent", Pattern: "`([^`\\\\]|\\\\.)*`"},
	{Name: "Operators", Pattern: `>=|<=|!=|<>|=|>|<`},
	{Name: "Punct", Pattern: `[().,]`},
})

// Global parser instance - built once, reused everywhere (thread-safe)
var (
	globalParser *participle.Parser[FilterRoot]
	parserOnce   sync.Once
)

// initParser builds the parser, called in getParser using sync.Once for thread safety
func initParser() {
	globalParser = participle.MustBuild[FilterRoot](
		participle.Lexer(sqlLexer),
		participle.Elide("whitespace", "Comment"),
		participle.CaseInsensitive("OR", "AND", "LIKE", "ILIKE", "IN", "true", "false"),
		participle.CaseInsensitive(StringValueType, DoubleValueType, IntValueType, BoolValueType, ArrayValueType),
	)
}

// getParser returns the singleton parser instance (thread-safe)
func getParser() *participle.Parser[FilterRoot] {
	parserOnce.Do(initParser)
	return globalParser
}

// Grammar structures for filter expressions.
//
// The grammar implements operator precedence through a hierarchy of types:
//
//	FilterRoot
//	  └── Expression
//	        └── OrExpression    (lowest precedence - evaluated last)
//	              └── AndExpression  (higher precedence - evaluated first)
//	                    └── Term (comparison or parenthesized group)
//
// This ensures AND binds tighter than OR, so:
//
//	"a = 1 OR b = 2 AND c = 3" is parsed as "a = 1 OR (b = 2 AND c = 3)"
//
// Even simple expressions like "name = 'test'" flow through the full hierarchy,
// just with empty Right slices at each level.

//nolint:govet
type FilterRoot struct {
	Expression *Expression `@@`
}

// Expression is the top-level grammar rule, delegating to OrExpression for precedence handling.
//
//nolint:govet
type Expression struct {
	Or *OrExpression `@@`
}

// OrExpression handles OR operators (lowest precedence).
// Left is always present; Right contains additional AND-expressions joined by OR.
//
//nolint:govet
type OrExpression struct {
	Left  *AndExpression   `@@`
	Right []*AndExpression `("OR" @@)*`
}

// AndExpression handles AND operators (higher precedence than OR).
// Left is always present; Right contains additional terms joined by AND.
//
//nolint:govet
type AndExpression struct {
	Left  *Term   `@@`
	Right []*Term `("AND" @@)*`
}

// Term is either a parenthesized sub-expression (Group) or a leaf Comparison.
//
//nolint:govet
type Term struct {
	Group      *Expression `"(" @@ ")"`
	Comparison *Comparison `| @@`
}

//nolint:govet
type Comparison struct {
	Left     *PropertyRef `@@`
	Operator string       `@("=" | "!=" | "<>" | ">=" | "<=" | ">" | "<" | "LIKE" | "ILIKE" | "IN")`
	Right    *Value       `@@`
}

//nolint:govet
type PropertyRef struct {
	EscapedName string   `@EscapedIdent`
	Name        string   `| @Ident`
	Path        []string `("." @Ident)*`
	Type        string   `("." @("string_value" | "double_value" | "int_value" | "bool_value"))?`
}

//nolint:govet
type Value struct {
	String    *string    `@String`
	Integer   *int64     `| @Int`
	Float     *float64   `| @Float`
	Boolean   *string    `| @("true" | "false" | "TRUE" | "FALSE")`
	ValueList *ValueList `| @@`
}

//nolint:govet
type ValueList struct {
	Values []*SingleValue `"(" (@@  ("," @@)*)? ")"`
}

//nolint:govet
type SingleValue struct {
	String  *string  `@String`
	Integer *int64   `| @Int`
	Float   *float64 `| @Float`
	Boolean *string  `| @("true" | "false" | "TRUE" | "FALSE")`
}

// FilterExpression represents a parsed filter expression (keeping for compatibility)
type FilterExpression struct {
	Left     *FilterExpression
	Right    *FilterExpression
	Operator string
	Property string
	Value    any
	IsLeaf   bool
}

// PropertyReference represents a property reference with type information
type PropertyReference struct {
	Name         string
	IsCustom     bool
	ValueType    string             // StringValueType, DoubleValueType, IntValueType, BoolValueType, ArrayValueType
	ExplicitType string             // Non-empty if the type was explicitly specified (e.g., "property.double_value")
	IsEscaped    bool               // whether the property name was escaped with backticks
	PropertyDef  PropertyDefinition // Full property definition for advanced handling
}

// Parse parses a filter query string and returns the root expression
// This function is thread-safe and reuses a singleton parser instance
func Parse(input string) (*FilterExpression, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}

	parser := getParser()
	filterRoot, err := parser.ParseString("", input)
	if err != nil {
		return nil, fmt.Errorf("error parsing filter query: %w", err)
	}

	return convertOrExpression(filterRoot.Expression.Or), nil
}

// convertOrExpression converts the participle AST OrExpression to our FilterExpression
func convertOrExpression(expr *OrExpression) *FilterExpression {
	left := convertAndExpression(expr.Left)

	for _, right := range expr.Right {
		rightExpr := convertAndExpression(right)
		left = &FilterExpression{
			Left:     left,
			Right:    rightExpr,
			Operator: "OR",
			IsLeaf:   false,
		}
	}

	return left
}

func convertAndExpression(expr *AndExpression) *FilterExpression {
	left := convertTerm(expr.Left)

	for _, right := range expr.Right {
		rightExpr := convertTerm(right)
		left = &FilterExpression{
			Left:     left,
			Right:    rightExpr,
			Operator: "AND",
			IsLeaf:   false,
		}
	}

	return left
}

func convertTerm(term *Term) *FilterExpression {
	if term.Group != nil {
		return convertOrExpression(term.Group.Or)
	}

	return convertComparison(term.Comparison)
}

func convertComparison(comp *Comparison) *FilterExpression {
	propRef := convertPropertyRef(comp.Left, comp.Right)
	value := convertValue(comp.Right)

	// Build the full property path including any nested paths
	propertyName := propRef.Name

	// Add type suffix if specified (for explicit type queries)
	if comp.Left.Type != "" {
		propertyName = propRef.Name + "." + comp.Left.Type
	}

	return &FilterExpression{
		Property: propertyName,
		Operator: comp.Operator,
		Value:    value,
		IsLeaf:   true,
	}
}

func convertPropertyRef(prop *PropertyRef, value *Value) *PropertyReference {
	var name string
	var isEscaped bool
	if prop.EscapedName != "" {
		// Remove backticks from escaped name
		name = strings.Trim(prop.EscapedName, "`")
		// Handle escape sequences in the name
		name = strings.ReplaceAll(name, `\.`, `.`)
		name = strings.ReplaceAll(name, `\\`, `\`)
		isEscaped = true
	} else {
		name = prop.Name
		// Build the full path if there are nested properties
		if len(prop.Path) > 0 {
			name = name + "." + strings.Join(prop.Path, ".")
		}
		isEscaped = false
	}

	var valueType string
	var isCustom bool

	if prop.Type != "" {
		// Explicit type specified
		valueType = prop.Type
		// We still need to determine if it's custom based on the property mapping
		// This is a bit tricky since we don't have entity type context here
		// For now, assume if explicit type is given, it could be either
		isCustom = true // Will be properly determined later in query builder
	} else {
		// Use the new property mapping system - but we need entity type context
		// For now, use a fallback approach and let the query builder handle it properly
		isCustom = true // Will be properly determined later in query builder
		valueType = inferValueType(value)
	}

	return &PropertyReference{
		Name:      name,
		IsCustom:  isCustom,
		ValueType: valueType,
		IsEscaped: isEscaped,
	}
}

func convertValue(val *Value) interface{} {
	if val.String != nil {
		return unquoteStringValue(*val.String)
	}

	if val.Integer != nil {
		return *val.Integer
	}

	if val.Float != nil {
		return *val.Float
	}

	if val.Boolean != nil {
		return strings.ToLower(*val.Boolean) == "true"
	}

	if val.ValueList != nil {
		// Convert list of values to slice
		var values []any
		for _, singleVal := range val.ValueList.Values {
			values = append(values, convertSingleValue(singleVal))
		}
		return values
	}

	return nil
}

func convertSingleValue(val *SingleValue) interface{} {
	if val.String != nil {
		return unquoteStringValue(*val.String)
	}

	if val.Integer != nil {
		return *val.Integer
	}

	if val.Float != nil {
		return *val.Float
	}

	if val.Boolean != nil {
		return strings.ToLower(*val.Boolean) == "true"
	}

	return nil
}

// unquoteStringValue removes quotes and handles escape sequences
func unquoteStringValue(str string) string {
	// Remove quotes from string
	result := strings.Trim(str, `"'`)
	// Handle escape sequences
	result = strings.ReplaceAll(result, `\"`, `"`)
	result = strings.ReplaceAll(result, `\'`, `'`)
	result = strings.ReplaceAll(result, `\\`, `\`)
	return result
}

// inferValueType determines the appropriate value type based on the actual value
func inferValueType(val *Value) string {
	if val.ValueList != nil && len(val.ValueList.Values) > 0 {
		// For lists, infer type from the first value
		return inferSingleValueType(val.ValueList.Values[0])
	}
	// Convert Value to SingleValue for type inference
	singleVal := &SingleValue{
		String:  val.String,
		Integer: val.Integer,
		Float:   val.Float,
		Boolean: val.Boolean,
	}
	return inferSingleValueType(singleVal)
}

// inferSingleValueType determines the appropriate value type for a SingleValue
func inferSingleValueType(val *SingleValue) string {
	if val.String != nil {
		return StringValueType
	}
	if val.Integer != nil {
		return IntValueType
	}
	if val.Float != nil {
		return DoubleValueType
	}
	if val.Boolean != nil {
		return BoolValueType
	}
	return StringValueType // default to string
}
