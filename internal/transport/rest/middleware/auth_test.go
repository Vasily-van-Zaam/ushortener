package middleware_test

import (
	"reflect"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/middleware"
)

func Test_generateRandom(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Generate random",
			args: args{
				size: 256,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := middleware.GenerateRandom(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateRandom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != 256 {
				t.Errorf("generateRandom() length = %v, want %v", len(got), tt.args.size)
			}
		})
	}
}

func Test_cripto_DencriptEncript(t *testing.T) {
	type args struct {
		src    []byte
		secret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "cripto",
			args: args{
				src:    []byte("hello world"),
				secret: "secret",
			},
			want: "hello world",
		},
		{
			name: "cripto",
			args: args{
				src:    []byte("123456789"),
				secret: "secret1",
			},
			want: "123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := middleware.NewCripto([]byte(tt.args.secret))
			if err != nil {
				t.Errorf("cripto.Dencript() error = %v", err)
			}
			gotE := e.Encript(tt.args.src)
			got, err := e.Dencript(gotE)
			if (err != nil) != tt.wantErr {
				t.Errorf("cripto.Dencript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("cripto.Dencript() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
