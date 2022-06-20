package infra

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// ParseGroup 解析以','分隔的整数, 去除0值, 去除重复值
func ParseGroup[T constraints.Integer](s string) []T {
	if s == "" {
		return []T{}
	}

	ss := strings.Split(s, ",")
	res := make([]T, 0, len(ss))
	mp := make(map[T]struct{})
	for i := 0; i < len(ss); i++ {
		v, err := strconv.ParseInt(strings.TrimSpace(ss[i]), 0, 64)
		if err != nil || v == 0 {
			continue
		}
		vv := T(v)
		if _, ok := mp[vv]; ok {
			continue
		}
		mp[vv] = struct{}{}
		res = append(res, vv)
	}
	return res
}

// JoinGroup 以','分隔的字符串为字符串, 去除0值, 去除重复值
func JoinGroup[T constraints.Integer](vs []T) string {
	sep := ","
	switch len(vs) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(int64(vs[0]), 10)
	}
	strElems := make([]string, 0, len(vs))
	n := 0
	mp := make(map[T]struct{})
	for i := 0; i < len(vs); i++ {
		if vs[i] == 0 {
			continue
		}
		if _, ok := mp[vs[i]]; ok {
			continue
		}
		mp[vs[i]] = struct{}{}
		v := strconv.FormatInt(int64(vs[i]), 10)
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
