package dbf

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"testing"
)

/**
 * a column name can be shorter than 10 runes and earlier terminated by byte(0) character
 */
func Test_newColumn(t *testing.T) {
	type args struct {
		rawData []byte
		enc     encoding.Encoding
	}
	type result struct {
		name   string
		typeOf ColumnType
	}

	tests := []struct {
		name    string
		args    args
		want    result
	}{
		{ name: "KUHNUMMER", args: args{[]byte{75, 85, 72, 78, 85, 77, 77, 69, 82, 0, 0, 78, 5, 0, 193, 57, 4, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KUHNUMMER", TypeNumber}},
		{ name: "RD_NR", args: args{[]byte{82, 68, 95, 78, 82, 0, 0, 0, 0, 0, 0, 78, 9, 0, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"RD_NR", TypeNumber}},
		{ name: "BULLENNR4", args: args{[]byte{66, 85, 76, 76, 69, 78, 78, 82, 52, 0, 0, 78, 12, 0, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLENNR4", TypeNumber}},
		{ name: "BULLENNR3", args: args{[]byte{66, 85, 76, 76, 69, 78, 78, 82, 51, 0, 0, 78, 15, 0, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLENNR3", TypeNumber}},
		{ name: "BULLENNR2", args: args{[]byte{66, 85, 76, 76, 69, 78, 78, 82, 50, 0, 0, 78, 18, 0, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLENNR2", TypeNumber}},
		{ name: "BULLENNR1", args: args{[]byte{66, 85, 76, 76, 69, 78, 78, 82, 49, 0, 0, 78, 21, 0, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLENNR1", TypeNumber}},
		{ name: "SPEICHER", args: args{[]byte{83, 80, 69, 73, 67, 72, 69, 82, 0, 82, 0, 76, 24, 0, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"SPEICHER", TypeBool}},
		{ name: "KOMMENT4", args: args{[]byte{75, 79, 77, 77, 69, 78, 84, 52, 0, 82, 0, 67, 25, 0, 193, 57, 22, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KOMMENT4", TypeText}},
		{ name: "KOMMENT3", args: args{[]byte{75, 79, 77, 77, 69, 78, 84, 51, 0, 82, 0, 67, 47, 0, 193, 57, 22, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KOMMENT3", TypeText}},
		{ name: "KOMMENT2", args: args{[]byte{75, 79, 77, 77, 69, 78, 84, 50, 0, 82, 0, 67, 69, 0, 193, 57, 22, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KOMMENT2", TypeText}},
		{ name: "KOMMENT1", args: args{[]byte{75, 79, 77, 77, 69, 78, 84, 49, 0, 82, 0, 67, 91, 0, 193, 57, 22, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KOMMENT1", TypeText}},
		{ name: "TRAG", args: args{[]byte{84, 82, 65, 71, 0, 78, 68, 49, 0, 83, 67, 76, 113, 0, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"TRAG", TypeBool}},
		{ name: "BULLE4", args: args{[]byte{66, 85, 76, 76, 69, 52, 0, 0, 46, 68, 66, 67, 114, 0, 193, 57, 12, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLE4", TypeText}},
		{ name: "BULLE3", args: args{[]byte{66, 85, 76, 76, 69, 51, 0, 0, 46, 68, 66, 67, 126, 0, 193, 57, 12, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLE3", TypeText}},
		{ name: "BULLE2", args: args{[]byte{66, 85, 76, 76, 69, 50, 0, 0, 46, 68, 66, 67, 138, 0, 193, 57, 12, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLE2", TypeText}},
		{ name: "DATUM_1BES", args: args{[]byte{68, 65, 84, 85, 77, 95, 49, 66, 69, 83, 0, 68, 150, 0, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"DATUM_1BES", TypeDate}},
		{ name: "BULLE1", args: args{[]byte{66, 85, 76, 76, 69, 49, 0, 0, 0, 0, 0, 67, 158, 0, 193, 57, 12, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLE1", TypeText}},
		{ name: "DATUM_2BES", args: args{[]byte{68, 65, 84, 85, 77, 95, 50, 66, 69, 83, 0, 68, 170, 0, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"DATUM_2BES", TypeDate}},
		{ name: "DATUM_3BES", args: args{[]byte{68, 65, 84, 85, 77, 95, 51, 66, 69, 83, 0, 68, 178, 0, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"DATUM_3BES", TypeDate}},
		{ name: "DATUM_4BES", args: args{[]byte{68, 65, 84, 85, 77, 95, 52, 66, 69, 83, 0, 68, 186, 0, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"DATUM_4BES", TypeDate}},
		{ name: "ABKALBUNG", args: args{[]byte{65, 66, 75, 65, 76, 66, 85, 78, 71, 0, 0, 68, 194, 0, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"ABKALBUNG", TypeDate}},
		{ name: "TRA", args: args{[]byte{84, 82, 65, 0, 0, 0, 0, 0, 0, 0, 0, 67, 202, 0, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"TRA", TypeText}},
		{ name: "NR", args: args{[]byte{78, 82, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 203, 0, 193, 57, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"NR", TypeNumber}},
		{ name: "AKTBESDAT", args: args{[]byte{65, 75, 84, 66, 69, 83, 68, 65, 84, 0, 0, 68, 205, 0, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"AKTBESDAT", TypeDate}},
		{ name: "BULLE", args: args{[]byte{66, 85, 76, 76, 69, 0, 0, 0, 0, 0, 0, 67, 213, 0, 193, 57, 12, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"BULLE", TypeText}},
		{ name: "STATUS", args: args{[]byte{83, 84, 65, 84, 85, 83, 0, 0, 0, 0, 0, 78, 225, 0, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"STATUS", TypeNumber}},
		{ name: "LEBENSNR", args: args{[]byte{76, 69, 66, 69, 78, 83, 78, 82, 0, 0, 0, 78, 226, 0, 193, 57, 9, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"LEBENSNR", TypeNumber}},
		{ name: "TRAGETAGE", args: args{[]byte{84, 82, 65, 71, 69, 84, 65, 71, 69, 0, 0, 78, 235, 0, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"TRAGETAGE", TypeNumber}},
		{ name: "KALBGE", args: args{[]byte{75, 65, 76, 66, 71, 69, 0, 0, 0, 0, 0, 78, 238, 0, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KALBGE", TypeNumber}},
		{ name: "GEBURTSV", args: args{[]byte{71, 69, 66, 85, 82, 84, 83, 86, 0, 0, 0, 78, 239, 0, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"GEBURTSV", TypeNumber}},
		{ name: "LAKTAGE", args: args{[]byte{76, 65, 75, 84, 65, 71, 69, 0, 0, 0, 0, 78, 240, 0, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"LAKTAGE", TypeNumber}},
		{ name: "LK_AK", args: args{[]byte{76, 75, 95, 65, 75, 0, 0, 0, 0, 0, 0, 68, 243, 0, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"LK_AK", TypeDate}},
		{ name: "VERBLEIB", args: args{[]byte{86, 69, 82, 66, 76, 69, 73, 66, 0, 0, 0, 78, 251, 0, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"VERBLEIB", TypeNumber}},
		{ name: "KALBNR", args: args{[]byte{75, 65, 76, 66, 78, 82, 0, 0, 0, 0, 0, 78, 252, 0, 193, 57, 9, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KALBNR", TypeNumber}},
		{ name: "ZKZ", args: args{[]byte{90, 75, 90, 0, 0, 0, 0, 0, 0, 0, 0, 78, 5, 1, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"ZKZ", TypeNumber}},
		{ name: "KALBKOM", args: args{[]byte{75, 65, 76, 66, 75, 79, 77, 0, 0, 0, 0, 67, 8, 1, 193, 57, 30, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KALBKOM", TypeText}},
		{ name: "EKALTER", args: args{[]byte{69, 75, 65, 76, 84, 69, 82, 0, 0, 0, 0, 78, 38, 1, 193, 57, 4, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"EKALTER", TypeNumber}},
		{ name: "EKA", args: args{[]byte{69, 75, 65, 0, 0, 0, 0, 0, 0, 0, 0, 78, 42, 1, 193, 57, 4, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"EKA", TypeNumber}},
		{ name: "KALBTAGE1", args: args{[]byte{75, 65, 76, 66, 84, 65, 71, 69, 49, 0, 0, 78, 46, 1, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KALBTAGE1", TypeNumber}},
		{ name: "KALBTAGE2", args: args{[]byte{75, 65, 76, 66, 84, 65, 71, 69, 50, 0, 0, 78, 49, 1, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KALBTAGE2", TypeNumber}},
		{ name: "KALBTAGE3", args: args{[]byte{75, 65, 76, 66, 84, 65, 71, 69, 51, 0, 0, 78, 52, 1, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KALBTAGE3", TypeNumber}},
		{ name: "KALBTAGE4", args: args{[]byte{75, 65, 76, 66, 84, 65, 71, 69, 52, 0, 0, 78, 55, 1, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"KALBTAGE4", TypeNumber}},
		{ name: "DTRAGETAGE", args: args{[]byte{68, 84, 82, 65, 71, 69, 84, 65, 71, 69, 0, 78, 58, 1, 193, 57, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"DTRAGETAGE", TypeNumber}},
		{ name: "RASSE", args: args{[]byte{82, 65, 83, 83, 69, 0, 0, 0, 0, 0, 0, 78, 61, 1, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"RASSE", TypeNumber}},
		{ name: "TROCKENGES", args: args{[]byte{84, 82, 79, 67, 75, 69, 78, 71, 69, 83, 0, 68, 62, 1, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"TROCKENGES", TypeDate}},
		{ name: "PAG", args: args{[]byte{80, 65, 71, 0, 0, 0, 0, 0, 0, 0, 0, 78, 70, 1, 193, 57, 4, 2, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"PAG", TypeNumber}},
		{ name: "IM", args: args{[]byte{73, 77, 0, 0, 0, 0, 0, 0, 0, 0, 0, 68, 74, 1, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"IM", TypeDate}},
		{ name: "TS", args: args{[]byte{84, 83, 0, 0, 0, 0, 0, 0, 0, 0, 0, 76, 82, 1, 193, 57, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"TS", TypeBool}},
		{ name: "IM2", args: args{[]byte{73, 77, 50, 0, 0, 0, 0, 0, 0, 0, 0, 68, 83, 1, 193, 57, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, charmap.CodePage850}, want: result{"IM2", TypeDate}},


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newColumn(tt.args.rawData, tt.args.enc)
			if err != nil {
				t.Errorf("newColumn() error = %v", err)
				return
			}
			assert.Equal(t, tt.want.name, got.Name)
			assert.Equal(t, tt.want.typeOf, got.Type)
		})
	}
}
