package c

import (
	"testing"

	"github.com/tomiya7688/neo-Cobol/internal/nir"
	"github.com/tomiya7688/neo-Cobol/internal/types"
)

// {
//   責務: [
//     TestEmitIntegerLiteralUsesSharedRadix: C出力が共通の整数基数規則を使うことを検証する
//   ]
//   処理: [
//     1: 先頭ゼロ付き10進と明示prefix付き整数をC式へ出力する
//     2: 出力値の文字列を期待値と比較する
//   ]
// }
func TestEmitIntegerLiteralUsesSharedRadix(t *testing.T) {
	cases := []struct {
		literal string
		want    string
	}{
		{literal: "010", want: "10"},
		{literal: "08", want: "8"},
		{literal: "-010", want: "-10"},
		{literal: "0o10", want: "8"},
		{literal: "0xFF", want: "255"},
	}
	for _, testCase := range cases {
		t.Run(testCase.literal, func(t *testing.T) {
			got, err := emitValue(nir.Value{Type: types.Integer, Text: testCase.literal}, nil)
			if err != nil {
				t.Fatalf("emitValue(%q) returned error: %v", testCase.literal, err)
			}
			if got != testCase.want {
				t.Fatalf("emitValue(%q) = %q, want %q", testCase.literal, got, testCase.want)
			}
		})
	}
}
