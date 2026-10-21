package types

import (
	"errors"
	"strconv"
	"testing"
)

// {
//   責務: [
//     TestParseIntegerLiteral: 明示prefixと先頭ゼロを含む整数の値を検証する
//   ]
//   処理: [
//     1: 10進数と2進・8進・16進の入力を解析する
//     2: 境界値を含む解析結果を期待値と比較する
//   ]
// }
func TestParseIntegerLiteral(t *testing.T) {
	cases := []struct {
		literal string
		want    int64
	}{
		{literal: "0", want: 0},
		{literal: "00", want: 0},
		{literal: "010", want: 10},
		{literal: "08", want: 8},
		{literal: "09", want: 9},
		{literal: "+010", want: 10},
		{literal: "-010", want: -10},
		{literal: "0b1010", want: 10},
		{literal: "0B1010", want: 10},
		{literal: "0o10", want: 8},
		{literal: "0O10", want: 8},
		{literal: "0xFF", want: 255},
		{literal: "0Xff", want: 255},
		{literal: "9223372036854775807", want: 1<<63 - 1},
		{literal: "-9223372036854775808", want: -1 << 63},
	}
	for _, testCase := range cases {
		t.Run(testCase.literal, func(t *testing.T) {
			got, err := ParseIntegerLiteral(testCase.literal)
			if err != nil {
				t.Fatalf("ParseIntegerLiteral(%q) returned error: %v", testCase.literal, err)
			}
			if got != testCase.want {
				t.Fatalf("ParseIntegerLiteral(%q) = %d, want %d", testCase.literal, got, testCase.want)
			}
		})
	}
}

// {
//   責務: [
//     TestParseIntegerLiteralDistinguishesInvalidDigitsAndOverflow: 不正な桁と範囲外を区別する
//   ]
//   処理: [
//     1: 不正な桁とLONG範囲外の入力を解析する
//     2: strconvのsyntax errorとrange errorを比較する
//   ]
// }
func TestParseIntegerLiteralDistinguishesInvalidDigitsAndOverflow(t *testing.T) {
	cases := []struct {
		literal string
		want    error
	}{
		{literal: "0b2", want: strconv.ErrSyntax},
		{literal: "0o8", want: strconv.ErrSyntax},
		{literal: "0xG", want: strconv.ErrSyntax},
		{literal: "9223372036854775808", want: strconv.ErrRange},
	}
	for _, testCase := range cases {
		t.Run(testCase.literal, func(t *testing.T) {
			_, err := ParseIntegerLiteral(testCase.literal)
			var numericError *strconv.NumError
			if !errors.As(err, &numericError) || numericError.Err != testCase.want {
				t.Fatalf("ParseIntegerLiteral(%q) error = %v, want NumError %v", testCase.literal, err, testCase.want)
			}
		})
	}
}

// {
//   責務: [
//     TestInferNumberLiteralUsesDecimalUnlessBasePrefixIsExplicit: 整数型推論の基数規則を検証する
//   ]
//   処理: [
//     1: 先頭ゼロ付き10進数と明示8進数を推論する
//     2: 推論結果を期待する整数型と比較する
//   ]
// }
func TestInferNumberLiteralUsesDecimalUnlessBasePrefixIsExplicit(t *testing.T) {
	cases := []struct {
		literal string
		want    Kind
	}{
		{literal: "010", want: Integer},
		{literal: "08", want: Integer},
		{literal: "-010", want: Integer},
		{literal: "0o10", want: Integer},
	}
	for _, testCase := range cases {
		got, err := InferNumberLiteral(testCase.literal)
		if err != nil {
			t.Fatalf("InferNumberLiteral(%q) returned error: %v", testCase.literal, err)
		}
		if got != testCase.want {
			t.Fatalf("InferNumberLiteral(%q) = %v, want %v", testCase.literal, got, testCase.want)
		}
	}
}
