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

// JoinGroup 以','分隔的字符串为字符串, 去除0值, 去除重复值
func JoinGroup(vs []int64) string {
	sep := ","
	switch len(vs) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(vs[0], 10)
	}
	strElems := make([]string, 0, len(vs))
	n := 0
	mp := make(map[int64]struct{})
	for i := 0; i < len(vs); i++ {
		if vs[i] == 0 {
			continue
		}
		if _, ok := mp[vs[i]]; ok {
			continue
		}
		mp[vs[i]] = struct{}{}
		v := strconv.FormatInt(vs[i], 10)
		strElems = append(strElems, v)
		n += len(v)
	}
	n += len(sep) * (len(strElems) - 1)

	var b strings.Builder
	b.Grow(n)
	b.WriteString(strElems[0])
	for _, s := range strElems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	return b.String()
}

// JoinGroupInt 以','分隔的字符串为字符串, 去除0值, 去除重复值
func JoinGroupInt(vs []int) string {
	sep := ","
	switch len(vs) {
	case 0:
		return ""
	case 1:
		return strconv.Itoa(vs[0])
	}
	strElems := make([]string, 0, len(vs))
	n := 0
	mp := make(map[int]struct{})
	for i := 0; i < len(vs); i++ {
		if vs[i] == 0 {
			continue
		}
		if _, ok := mp[vs[i]]; ok {
			continue
		}
		mp[vs[i]] = struct{}{}
		v := strconv.Itoa(vs[i])
		strElems = append(strElems, v)
		n += len(v)
	}
	n += len(sep) * (len(strElems) - 1)

	var b strings.Builder
	b.Grow(n)
	b.WriteString(strElems[0])
	for _, s := range strElems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	return b.String()
}
