package time

import (
	"reflect"
	"testing"
	"time"
)

func TestGetBetweenDays(t *testing.T) {
	type xargs struct {
		beginDate string
		args      []string
	}
	tests := []struct {
		name    string
		xargs   xargs
		want    []string
		wantErr bool
	}{
		{
			name: "test between days",
			xargs: xargs{
				beginDate: "2022-02-27",
				args:      []string{"2022-03-01"},
			},
			want:    []string{"2022-02-27", "2022-02-28", "2022-03-01"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBetweenDays(LayoutIsoDate, tt.xargs.beginDate, tt.xargs.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBetweenDays() panic = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBetweenDays() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBeginEndDayTime(t *testing.T) {
	type args struct {
		strBeginDate string
		strEndDate   string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		want1   time.Time
		wantErr bool
	}{
		{
			name: "TestToBeginEndDayTime",
			args: args{
				strBeginDate: "2024-x01-01",
				strEndDate:   "2024-01-02",
			},
			want:    time.Time{},
			want1:   time.Time{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ToDayBeginEndTime(tt.args.strBeginDate, tt.args.strEndDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBeginEndDayTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBeginEndDayTime() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToBeginEndDayTime() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
