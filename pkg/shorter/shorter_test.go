package shorter_test

import (
	"log"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/pkg/shorter"
)

func Test_shorter_Convet(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				id: "1",
			},
		},
		{
			args: args{
				id: "1000",
			},
		},
		{
			args: args{
				id: "1000000",
			},
		},
		{
			args: args{
				id: "10000000000000000",
			},
		},
		{
			args: args{
				id: "95546546565465465",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shorter.NewShorter59()
			idStr := s.Convert(tt.args.id)
			id := s.UnConvert(idStr)
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
		want string
	}{
		{
			args: args{
				id: "$ZZZZZZZZZZ",
			},
			want: "511116753300641400",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shorter.NewShorter59()
			// got := s.ToInt(tt.args.id)
			// log.Println(tt.args.id, "=", got)
			if got := s.UnConvert(tt.args.id); got != tt.want {
				t.Errorf("shorter.ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shorter_Some_Convert(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				id: "123434343",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := shorter.New([]string{
				"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F",
			})
			got := s.Convert(tt.args.id)
			got1 := s.UnConvert(got)
			log.Println(got, got1)
			// got := s.ToInt(tt.args.id)
			// log.Println(tt.args.id, "=", got)
			// if got := s.Convert(tt.args.id); got != tt.want {
			// 	t.Errorf("shorter.ToInt() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func BenchmarkConvet(b *testing.B) {
	s := shorter.NewShorter59()

	b.ResetTimer()

	s.Convert("1234567890190000000")
}
