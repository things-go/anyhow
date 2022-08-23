package pickup

import (
	"strings"
)

const sixStar = "******"

// HideCard 隐藏证件号码.证件号码都为数字+字母
// len(card) == 0: ******
// len(card) <= 4: *(x len(card))
// len(card) <= 10: xxx******
// len(card) >= 10 xxx******xxx
func HideCard(card string) string {
	length := len(card)
	switch {
	case length == 0:
		return sixStar
	case length <= 4:
		return strings.Repeat("*", length)
	case length <= 10:
		return BuildHideString(card[:4], "", length-4)
	default: // length > 10
		return BuildHideString(card[:4], card[length-3:], length-7)
	}
}

// HideMobile 隐藏手机号.
// len(mobile) <= 7: ******
// len(mobile) >= 7 xxx*****xxxx
func HideMobile(mobile string) string {
	length := len(mobile)
	switch {
	case length == 0:
		return sixStar
	case length < 7:
		return strings.Repeat("*", length)
	default:
		return BuildHideString(mobile[:3], mobile[length-4:], length-7)
	}
}

// HideName 隐藏真实名称(如姓名、账号、公司等).
// ""                ==> ******
// 李                ==> 李
// 李四               ==> 李*
// 张三丰              ==> 张*丰
// 公孙先生             ==> 公孙**
// helloWorld           ==> hel****rld
// 北京搜狗科技公司         ==> 北京****公司
// 北京搜狗科技发展有限公司  ==> 北京搜******限公司
// 工商发展银行深圳南山科苑梅龙路支行  ==> 工商发展*********龙路支行
func HideName(s string) string {
	if s == "" {
		return sixStar
	}
	runs := []rune(s)
	length := len(runs)
	switch {
	case length <= 2:
		return BuildHideString(string(runs[:1]), "", length-1)
	case length == 3:
		return BuildHideString(string(runs[:1]), string(runs[length-1:]), length-2)
	case length < 5:
		return BuildHideString(string(runs[:2]), "", length-2)
	case length < 10:
		return BuildHideString(string(runs[:2]), string(runs[length-2:]), length-4)
	case length < 16:
		return BuildHideString(string(runs[:3]), string(runs[length-3:]), length-6)
	default:
		return BuildHideString(string(runs[:4]), string(runs[length-4:]), length-8)
	}
}

// HideLastString 隐藏字符串最后 length 位
// s == "" || length <= 0: ******
// len(runs) <= length: * (x length)
// other: xxxx[*(x length)]
func HideLastString(s string, length int) string {
	if s == "" || length <= 0 {
		return sixStar
	}
	runs := []rune(s)
	if len(runs) <= length {
		return strings.Repeat("*", length)
	}
	return BuildHideString(string(runs[:len(runs)-length]), "", length)
}

// ""  ==>  ******
// 李  ==>  ******
// 李四  ==>  李*
// 张三丰  ==>  张*丰
// 公孙先生  ==>  公**生
// helloWorld  ==>  he*****ld
// 北京搜狗科技公司  ==>  北京****公司
// 北京搜狗科技发展有限公司  ==>  北京搜******限公司
// 工商发展银行深圳南山科苑梅龙路支行  ==>  工商发展********龙路支行
func HideMiddleString(s string) string {
	runs := []rune(s)
	hideLen := len(runs) / 2
	showLen := len(runs) - hideLen
	if hideLen == 0 || showLen == 0 {
		return sixStar
	}
	subLen := showLen / 2
	if subLen == 0 {
		return BuildHideString(string(runs[:showLen]), "", hideLen)
	}
	return BuildHideString(string(runs[:subLen]), string(runs[len(runs)-subLen:]), hideLen)
}

// BuildHideString 生成隐藏字符串,中间使用 '*' 表示
func BuildHideString(prefix, suffix string, midStarRepeatCnt int) string {
	var b strings.Builder

	b.Grow(len(prefix) + len(suffix) + midStarRepeatCnt)
	b.WriteString(prefix)
	for b.Len() < len(prefix)+midStarRepeatCnt {
		b.WriteString("*")
	}
	b.WriteString(suffix)
	return b.String()
}
