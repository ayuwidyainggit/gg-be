package sql_helper

import "testing"

func TestParseSort(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Test ParseSort",
			args: args{
				input: "field:order",
			},
			want:  "field",
			want1: "desc",
		},
		{
			name: "Test ParseSort",
			args: args{
				input: "field",
			},
			want:  "field",
			want1: "desc",
		},
		{
			name: "Test ParseSort",
			args: args{
				input: "field:asc",
			},
			want:  "field",
			want1: "asc",
		},
		{
			name: "Test ParseSort",
			args: args{
				input: "field:desc",
			},
			want:  "field",
			want1: "desc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseSort(tt.args.input)
			if got != tt.want {
				t.Errorf("ParseSort() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseSort() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
