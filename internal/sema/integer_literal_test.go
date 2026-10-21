package sema

import (
	"testing"

	"github.com/tomiya7688/neo-Cobol/internal/ast"
	"github.com/tomiya7688/neo-Cobol/internal/types"
)

// {
//   責務: [
//     TestLiteralFitsTargetUsesSharedIntegerRadix: 意味解析が共通の整数基数規則を使うことを検証する
//   ]
//   処理: [
//     1: 先頭ゼロ付き10進と明示prefix付き整数を各型へ代入可能か判定する
//     2: BYTE範囲外とLONG範囲外を拒否する
//   ]
// }
func TestLiteralFitsTargetUsesSharedIntegerRadix(t *testing.T) {
	cases := []struct {
		literal string
		target  types.Kind
		want    bool
	}{
		{literal: "08", target: types.Byte, want: true},
		{literal: "010", target: types.Byte, want: true},
		{literal: "-010", target: types.Integer, want: true},
		{literal: "0xFF", target: types.Byte, want: true},
		{literal: "0b100000000", target: types.Byte, want: false},
		{literal: "9223372036854775808", target: types.Long, want: false},
	}
	for _, testCase := range cases {
		t.Run(testCase.literal, func(t *testing.T) {
			got := literalFitsTarget(ast.NumberLiteral{Value: testCase.literal}, testCase.target)
			if got != testCase.want {
				t.Fatalf("literalFitsTarget(%q, %s) = %t, want %t", testCase.literal, testCase.target, got, testCase.want)
			}
		})
	}
}
