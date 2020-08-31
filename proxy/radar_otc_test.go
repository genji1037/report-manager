package proxy

import (
	"reflect"
	"testing"
)

func TestGetRetryOrFailedTransfer(t *testing.T) {
	tests := []struct {
		name    string
		want    []BlockChainTransfer
		wantErr bool
	}{
		{
			name:    "simple",
			want:    []BlockChainTransfer{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRetryOrFailedTransfer()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRetryOrFailedTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRetryOrFailedTransfer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
