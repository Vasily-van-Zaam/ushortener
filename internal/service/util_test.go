package service_test

import (
	"log"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/service"
)

func Test_shorter_Convet(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				id: 1,
			},
		},
		{
			args: args{
				id: 1000,
			},
		},
		{
			args: args{
				id: 1000000,
			},
		},
		{
			args: args{
				id: 10000000000000000,
			},
		},
		{
			args: args{
				id: 95546546565465465,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewShorter()
			idStr := s.ToString(tt.args.id)
			id := s.ToInt(idStr)
			log.Println(id, "==", tt.args.id, "==", idStr)

			if id != tt.args.id {
				t.Errorf("shorter.ToInt(%v) == %v, want tt.args.id == %v", idStr, tt.args.id, id)
			}
		})
	}
}

func Test_shorter_ToInt(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			args: args{
				id: "$ZZZZZZZZZZ",
			},
			want: 511116753300641400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewShorter()
			// got := s.ToInt(tt.args.id)
			// log.Println(tt.args.id, "=", got)
			if got := s.ToInt(tt.args.id); got != tt.want {
				t.Errorf("shorter.ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConvet(b *testing.B) {
	s := service.NewShorter()

	b.ResetTimer()

	s.ToString(1234567890190000000)
}
