package sql

import (
	"fmt"
	"github.com/hdget/hdutils/json"
	"testing"
)

func TestNewBatchUpdater(t *testing.T) {
	type args struct {
		table      string
		updateSet  string
		whenColumn string
	}
	tests := []struct {
		name string
		args args
		want BatchUpdater
	}{
		{
			name: "TestNewBatchUpdater",
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cases := []*mysqlBatchUpdateCase{
				{1, 2},
				{2, json.JsonObject()},
			}

			u := NewMysqlBatchUpdater("table", "name", "id")
			for _, i2 := range cases {
				u.Add(i2.WhenValue, i2.ThenValue)
			}
			sql, err := u.Generate()
			if err != nil {
				panic(err)
			}
			fmt.Println(sql)
		})
	}
}
