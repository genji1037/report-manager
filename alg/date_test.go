package alg

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	fmt.Println(time.Now().In(time.UTC).Format("2006-01-02 15:04:05"))
}

func TestConvertLayout(t *testing.T) {
	type args struct {
		str       string
		oldLayout string
		newLayout string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				str:       "2021-09-01",
				oldLayout: "2006-01-02",
				newLayout: "2006-1-2",
			},
			want:    "2021-9-1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertLayout(tt.args.str, tt.args.oldLayout, tt.args.newLayout)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertLayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertLayout() got = %v, want %v", got, tt.want)
			}
		})
	}
}
