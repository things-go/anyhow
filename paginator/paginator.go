package paginator

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Param 分页查询参数
type Param struct {
	PageIndex int64 `form:"page_index" json:"page_index"`
	PageSize  int64 `form:"page_size" json:"page_size"`
}

// New 创建分页实例并校验分页查询参数有效性,非法值将设置为默认值,
// 默认分页索引: 第1页
// 默认分页大小: DefaultPageSize
// 默认最大分页大小: DefaultMaxPageSize
// 可修改 DefaultPageSize, DefaultMaxPageSize 改变全局默认分页大小或传入 limitPageSize 参数
func New(index, size int64, limitPageSize ...int64) Param {
	return Param{
		PageIndex: index,
		PageSize:  size,
	}.Inspect(limitPageSize...)
}

// Inspect 校验分页查询参数有效性,非法值将设置为默认值,
// 默认分页索引: 第1页
// 默认分页大小: DefaultPageSize
// 默认最大分页大小: DefaultMaxPageSize
// 可修改 DefaultPageSize, DefaultMaxPageSize 改变全局默认分页大小或传入参数
func (sf Param) Inspect(pageSize ...int64) Param {
	if sf.PageIndex <= 0 {
		sf.PageIndex = 1
	}

	if sf.PageSize <= 0 || sf.PageSize > DefaultMaxPageSize {
		size := int64(-1)
		if len(pageSize) > 0 && pageSize[0] > 0 && pageSize[0] < DefaultMaxPageSize {
			size = pageSize[0]
		}
		switch {
		case size > 0:
			sf.PageSize = size
		case sf.PageSize <= 0:
			sf.PageSize = DefaultPageSize
		default:
			sf.PageSize = DefaultMaxPageSize
		}
	}
	return sf
}

func (sf Param) Value() (limit int64, offset int64) {
	if sf.PageSize > 0 {
		limit = sf.PageSize
		if sf.PageIndex > 0 {
			offset = sf.PageSize * (sf.PageIndex - 1)
		}
	}
	return limit, offset
}

func (sf Param) Limit() clause.Expression {
	limit, offset := sf.Value()
	return clause.Limit{
		Limit:  int(limit),
		Offset: int(offset),
	}
}

func (sf Param) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(sf.Limit())
	}
}
