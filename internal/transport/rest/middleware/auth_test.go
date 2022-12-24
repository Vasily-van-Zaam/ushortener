package middleware

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

func TestAuth_Handle(t *testing.T) {
	type fields struct {
		Config *core.Config
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				Config: tt.fields.Config,
			}
			if got := a.Handle(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}
