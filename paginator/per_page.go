package paginator

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Param2 分页查询参数
type Param2 struct {
	Page    int64 `form:"page" json:"page"`
	PerPage int64 `form:"per_page" json:"per_page"`
}

// New2 创建分页实例并校验分页查询参数有效性,非法值将设置为默认值,
// 默认分页索引: 第1页
// 默认分页大小: DefaultPageSize
// 默认最大分页大小: DefaultMaxPageSize
// 可修改 DefaultPageSize, DefaultMaxPageSize 改变全局默认分页大小或传入 limitPageSize 参数
func New2(page, perPage int64, limitPerPage ...int64) Param2 {
	return Param2{
		Page:    page,
		PerPage: perPage,
	}.Inspect(limitPerPage...)
}

// Inspect 校验分页查询参数有效性,非法值将设置为默认值,
// 默认分页索引: 第1页
// 默认分页大小: DefaultPageSize
// 默认最大分页大小: DefaultMaxPageSize
// 可修改 DefaultPageSize, DefaultMaxPageSize 改变全局默认分页大小或传入参数
func (p Param2) Inspect(pageSize ...int64) Param2 {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PerPage <= 0 || p.PerPage > DefaultMaxPageSize {
		size := int64(-1)
		if len(pageSize) > 0 && pageSize[0] > 0 && pageSize[0] < DefaultMaxPageSize {
			size = pageSize[0]
		}
		switch {
		case size > 0:
			p.PerPage = size
		case p.PerPage <= 0:
			p.PerPage = DefaultPageSize
		default:
			p.PerPage = DefaultMaxPageSize
		}
	}
	return p
}

func (sf Param2) Value() (limit int64, offset int64) {
	if sf.PerPage > 0 {
		limit = sf.PerPage
		if sf.Page > 0 {
			offset = sf.PerPage * (sf.Page - 1)
		}
	}
	return limit, offset
}

func (sf Param2) Limit() clause.Expression {
	limit, offset := sf.Value()
	return clause.Limit{
		Limit:  int(limit),
		Offset: int(offset),
	}
}

func (sf Param2) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(sf.Limit())
	}
}
