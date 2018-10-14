package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/parser"
)

func TestParser_ParseField(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantField *parser.Field
		wantErr   bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without fieldNumber",
			input:   "foo.bar nested_message = ;",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; string fieldNumber",
			input:   "foo.bar nested_message = a;",
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: "foo.bar nested_message = 2;",
			wantField: &parser.Field{
				Type:        "foo.bar",
				FieldName:   "nested_message",
				FieldNumber: "2",
			},
		},
		{
			name:  "parsing another excerpt from the official reference",
			input: "repeated int32 samples = 4 [packed=true];",
			wantField: &parser.Field{
				IsRepeated:  true,
				Type:        "int32",
				FieldName:   "samples",
				FieldNumber: "4",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "packed",
						Constant:   "true",
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
			got, err := parser.ParseField()
			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantField) {
				t.Errorf("got %v, but want %v", got, test.wantField)
			}

			if !parser.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
