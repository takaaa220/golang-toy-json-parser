package parser

import (
	"reflect"
	"testing"

	"github.com/takaaa220/golang-toy-json-parser/lexer"
)

func TestParser_Parse(t *testing.T) {

	tests := []struct {
		input   string
		lexer   *lexer.Lexer
		want    interface{}
		wantErr bool
	}{
		{
			input:   "null",
			want:    nil,
			wantErr: false,
		},
		{
			input:   "true",
			want:    true,
			wantErr: false,
		},
		{
			input:   `{"key1": "value", "key2": 123, "key3": [1, true, false, null, "string"]}`,
			want:    map[string]interface{}{"key1": "value", "key2": float64(123), "key3": []interface{}{float64(1), true, false, nil, "string"}},
			wantErr: false,
		},
		{
			input:   `[true, null, {"key1": "value"}, 1.01]`,
			want:    []interface{}{true, nil, map[string]interface{}{"key1": "value"}, float64(1.01)},
			wantErr: false,
		},
		{
			input:   `{key: "value"}`,
			want:    nil,
			wantErr: true,
		},
		{
			input:   `[1, 2, 3,]`,
			want:    nil,
			wantErr: true,
		},
		{
			input:   `{"key": "value",}`,
			want:    nil,
			wantErr: true,
		},
		{
			input:   `{"key": "value", "key2": "value2"} [1, 2, 3]`,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := NewParser(lexer.NewLexer(tt.input))

			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
