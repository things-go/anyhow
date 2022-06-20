package infra

import (
	"strconv"
	"strings"
)

// ParseGroup 解析以','分隔的整数, 去除0值, 去除重复值
func ParseGroup(s string) []int64 {
	if s == "" {
		return []int64{}
	}

	ss := strings.Split(s, ",")
	res := make([]int64, 0, len(ss))
	mp := make(map[int64]struct{})
	for i := 0; i < len(ss); i++ {
		v, err := strconv.ParseInt(strings.TrimSpace(ss[i]), 0, 64)
		if err != nil || v == 0 {
			continue
		}
		if _, ok := mp[v]; ok {
			continue
		}
		mp[v] = struct{}{}
		res = append(res, v)
	}
	return res
}

// ParseGroupInt 解析以','分隔的整数, 去除0值, 去除重复值
func ParseGroupInt(s string) []int {
	if s == "" {
		return []int{}
	}

	ss := strings.Split(s, ",")
	res := make([]int, 0, len(ss))
	mp := make(map[int]struct{})
	for i := 0; i < len(ss); i++ {
		v, err := strconv.Atoi(strings.TrimSpace(ss[i]))
		if err != nil || v == 0 {
			continue
		}
		if _, ok := mp[v]; ok {
			continue
		}
		mp[v] = struct{}{}
		res = append(res, v)
	}
	return res
}
