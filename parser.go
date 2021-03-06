package gosnow

/*
#cgo CFLAGS: -I./drafter/src/ -I./drafter/ext/snowcrash/src/
#cgo LDFLAGS: -ldrafter -L"./drafter/build/out/Release/"
#include <stdlib.h>
#include <stdio.h>
#include "cdrafter.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// Warning Codes
// see https://github.com/apiaryio/snowcrash/blob/master/src/SourceAnnotation.h#L128
const (
	NoWarning                 = 0
	APINameWarning            = 1
	DuplicateWarning          = 2
	FormattingWarning         = 3
	RedefinitionWarning       = 4
	IgnoringWarning           = 5
	EmptyDefinitionWarning    = 6
	NotEmptyDefinitionWarning = 7
	LogicalErrorWarning       = 8
	DeprecatedWarning         = 9
	IndentationWarning        = 10
	AmbiguityWarning          = 11
	URIWarning                = 12
)

// Error Codes
// see https://github.com/apiaryio/snowcrash/blob/master/src/SourceAnnotation.h#L113
const (
	NoError          = 0
	ApplicationError = 1
	BusinessError    = 2
	SymbolError      = 3
)

// LocationPoint makes up an array of possibly non-continuous blocks of the source API Blueprint.
type LocationPoint struct {
	// Zero-based index of the character where warning has occurred.
	Index int
	//  Number of the characters from index where warning has occurred.
	Length int
}

// SourceAnnotation is a building block of the json for API Blueprints
type SourceAnnotation struct {
	Message  string
	Code     int
	Location []LocationPoint
}

// Ok returns whether or not the return code indicated success
func (s *SourceAnnotation) Ok() bool {
	return s.Code == 0
}

// ParseResult contains all the information resulting from the parse
// in a fully typed structured format
// see: https://github.com/apiaryio/api-blueprint-ast/blob/master/Parse%20Result.md
type ParseResult struct {
	AST       Blueprint
	SourceMap BlueprintSourcemap
	// Description of a parsing error as occurred during parsing.
	// If this field is present && code different from 0 then the
	// content of ast field should be ignored.
	Error SourceAnnotation
	// Ordered array of parser warnings as occurred during the parsing.
	Warnings []SourceAnnotation
}

func newPR(data []byte) (*ParseResult, error) {
	pr := new(ParseResult)
	if err := json.Unmarshal(data, pr); err != nil {
		return nil, err
	}
	return pr, nil
}

const (
	// ScRenderDescriptionsOptionKey causes the parser to render markdown in description
	ScRenderDescriptionsOptionKey = 1 << 0
	// RequireBlueprintNameOptionKey causes the parser to treat missing blueprint name as error
	RequireBlueprintNameOptionKey = 1 << 1
	// ExportSourcemapOptionKey causes the parser to export source maps AST
	ExportSourcemapOptionKey = 1 << 2
)

// RawOptionParse parses the inputted string and passes it to the drafter library for parsing.
// The raw JSON result is then returned from the function
func RawOptionParse(source string, flags int) ([]byte, error) {
	cs := C.CString(source)
	var result string
	cr := C.CString(result)
	ret := int(C.drafter_c_parse(cs, C.sc_blueprint_parser_options(flags), &cr))
	if ret != 0 {
		return nil, fmt.Errorf("drafter_c_parse failed with code: %v", ret)
	}

	res := C.GoString(cr)
	C.free(unsafe.Pointer(cr)) /* we MUST release allocted memory for result */

	return []byte(res), nil
}

// OptionParse parses the inputted string using the drafter library with the
// specified options to affect the parse and unmarshals the returned json
// into a struct containing all the necessary information
func OptionParse(source string, options int) (*ParseResult, error) {
	data, err := RawOptionParse(source, options)
	if err != nil {
		return nil, err
	}
	pr, err := newPR(data)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

// Parse is a wrapper around OptionParse passing in 0 for
// the options value, indicating normal parsing behavior
func Parse(source string) (*ParseResult, error) {
	return OptionParse(source, 0)
}
