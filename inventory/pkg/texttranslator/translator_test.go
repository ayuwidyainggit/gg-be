package texttranslator

import (
	"testing"
)

func BenchmarkTranslator2(b *testing.B) {
	const text string = `c_pro_code: 103 is already exist`

	trans, err := New(2)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < b.N; i++ {
		trans.Translate(ID, text)

	}
}

func BenchmarkTranslator(b *testing.B) {
	const text string = `c_pro_code: 103 is already exist`

	trans, err := New(1)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < b.N; i++ {
		trans.Translate(ID, text)

	}
}

func TestTranslator_Translate(t *testing.T) {
	type fields struct {
		Version string
	}
	type args struct {
		lang    string
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "test googletranslatefree version",
			fields: fields{Version: V1},
			args:   args{lang: ID, message: "c_pro_code: 103 is already exist"},
			want:   "c_pro_code: 103 sudah ada",
		},
		{
			name:   "test Conight/go-googletrans version",
			fields: fields{Version: V2},
			args:   args{lang: ID, message: "c_pro_code: 103 is already exist"},
			want:   "c_pro_code: 103 sudah ada",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Translator{
				Version: tt.fields.Version,
			}
			if got := tr.Translate(tt.args.lang, tt.args.message); got != tt.want {
				t.Errorf("Translator.Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}
