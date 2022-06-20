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
			if got := ParseGroup(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIdsGroupInt(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []int
	}{
		{
			"empty",
			"",
			[]int{},
		},
		{
			"1",
			"1",
			[]int{1},
		},
		{
			"> 1",
			"1,10,11,12",
			[]int{1, 10, 11, 12},
		},
		{
			"> 1 contain space",
			"1, 10, 11 ,  12",
			[]int{1, 10, 11, 12},
		},
		{
			"> 1 contain dump value",
			"1, 10, 11 ,  12, 11,1",
			[]int{1, 10, 11, 12},
		},
		{
			"> 1 contain zero value",
			"0, 1, 10, 11 ,  12, 0",
			[]int{1, 10, 11, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseGroupInt(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGroupInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
