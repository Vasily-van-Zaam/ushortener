package psql

import (
	"context"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestStore_GetStats(t *testing.T) {
	type fields struct {
		config *core.Config
		db     *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.Stats
		wantErr bool
	}{
		{},
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// s := &Store{
			// 	config: tt.fields.config,
			// 	db:     tt.fields.db,
			// }
			// got, err := s.GetStats(tt.args.ctx)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("Store.GetStats() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Store.GetStats() = %v, want %v", got, tt.want)
			// }
		})
	}
}
