package alg

import "testing"

func TestStrArr2Str(t *testing.T) {
	type args struct {
		strs []string
		sep  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				strs: []string{},
				sep:  ",",
			},
			want: "",
		},
		{
			args: args{
				strs: []string{"cy"},
				sep:  ",",
			},
			want: "cy",
		},
		{
			args: args{
				strs: []string{"cy", "xy"},
				sep:  ",",
			},
			want: "cy,xy",
		},
		{
			args: args{
				strs: []string{"cy", "xy", "sf"},
				sep:  ",",
			},
			want: "cy,xy,sf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrArr2Str(tt.args.strs, "'", tt.args.sep); got != tt.want {
				t.Errorf("StrArr2Str() = %v, want %v", got, tt.want)
			}
		})
	}
}
