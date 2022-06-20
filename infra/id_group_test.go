package infra

import (
	"reflect"
	"testing"
)

func TestParseIdsGroup(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []int64
	}{
		{
			"empty",
			"",
			[]int64{},
		},
		{
			"1",
			"1",
			[]int64{1},
		},
		{
			"> 1",
			"1,10,11,12",
			[]int64{1, 10, 11, 12},
		},
		{
			"> 1 contain space",
			"1, 10, 11 ,  12",
			[]int64{1, 10, 11, 12},
		},
		{
			"> 1 contain dump value",
			"1, 10, 11 ,  12, 11,1",
			[]int64{1, 10, 11, 12},
		},
		{
			"> 1 contain zero value",
			"0, 1, 10, 11 ,  12, 0",
			[]int64{1, 10, 11, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseGroup[int64](tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdsGroup(t *testing.T) {
	tests := []struct {
		name string
		ids  []int64
		want string
	}{
		{
			"empty nil",
			nil,
			"",
		},
		{
			"empty",
			[]int64{},
			"",
		},
		{
			"1",
			[]int64{1},
			"1",
		},
		{
			"> 1",
			[]int64{1, 10, 11, 12},
			"1,10,11,12",
		},
		{
			"> 1 contain dump value",
			[]int64{1, 10, 11, 11, 12, 1},
			"1,10,11,12",
		},
		{
			"> 1 contain zero/dump value",
			[]int64{0, 1, 10, 0, 11, 11, 0, 12},
			"1,10,11,12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinGroup(tt.ids); got != tt.want {
				t.Errorf("JoinGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
