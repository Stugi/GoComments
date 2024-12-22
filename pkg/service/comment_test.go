package service

import "testing"

func Test_checkComments(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "example1",
			args: args{
				text: "example1",
			},
			want: false,
		},
		{
			name: "example2",
			args: args{
				text: "example2 qwerty, йцукен, zxvbnm",
			},
			want: true,
		},
		{
			name: "example3",
			args: args{
				text: "example3 qwertYйцук",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkComments(tt.args.text); got != tt.want {
				t.Errorf("checkComments() = %v, want %v", got, tt.want)
			}
		})
	}
}
