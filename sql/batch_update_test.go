package sql

import (
	"github.com/hdget/hdutils/json"
	"testing"
)

func Test_mysqlBatchUpdater_Generate(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "Test_mysqlBatchUpdater_Generate",
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewMysqlBatchUpdater("order").Set("kind", 1).Case("request", "id")
			u.When("1", json.JsonObject())
			got, err := u.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Generate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
