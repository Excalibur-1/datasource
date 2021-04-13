package search

import (
	"strings"

	"xorm.io/builder"
	"xorm.io/xorm"

	"github.com/Excalibur-1/datasource/common"
	"github.com/Excalibur-1/datasource/sort"
)

// Query {"pageSize":10,"page":0,"sorted":[{"id":"firstName","desc":false}],"filtered":[{"id":"firstName","value":"3"}]}
type Query struct {
	common.Query
}

func NewQuery(query common.Query) Query {
	return Query{query}
}

func (sp Query) MarkOrder(column map[string]Filter) (sorted *sort.Sort) {
	if len(sp.Sorted) > 0 {
		sorted = sort.Sorted()
		for _, v := range sp.Sorted {
			if v.Desc {
				sorted.Desc(column[v.Id].FieldName)
			} else {
				sorted.Asc(column[v.Id].FieldName)
			}
		}
	}
	return
}

func (sp Query) MarkOrmFiltered(column map[string]Filter, orm *xorm.Session) {
	if len(sp.Filtered) > 0 {
		for _, v := range sp.Filtered {
			if k, ok := column[v.Id]; ok {
				switch k.Operator {
				case NE:
					orm.Where(builder.Neq{k.FieldName: v.Value})
				case LIKE:
					orm.Where(builder.Like{k.FieldName, v.Value.(string)})
				case GT:
					orm.Where(builder.Gt{k.FieldName: v.Value})
				case LT:
					orm.Where(builder.Lt{k.FieldName: v.Value})
				case GTE:
					orm.Where(builder.Gte{k.FieldName: v.Value})
				case LTE:
					orm.Where(builder.Lte{k.FieldName: v.Value})
				case IN:
					orm.Where(markIn(true, k.FieldName, v.Value))
				case NI:
					orm.Where(markIn(false, k.FieldName, v.Value))
				case IsNull:
					orm.Where(builder.IsNull{k.FieldName})
				case NotNull:
					orm.Where(builder.NotNull{k.FieldName})
				default:
					orm.Where(builder.Eq{k.FieldName: v.Value})
				}
			}
		}
	}
}

func (sp Query) MarkSqlFiltered(column map[string]Filter, bu *builder.Builder) {
	if len(sp.Filtered) > 0 {
		for _, v := range sp.Filtered {
			if k, ok := column[v.Id]; ok {
				switch k.Operator {
				case NE:
					bu.Where(builder.Neq{k.FieldName: v.Value})
				case LIKE:
					bu.Where(builder.Like{k.FieldName, v.Value.(string)})
				case GT:
					bu.Where(builder.Gt{k.FieldName: v.Value})
				case LT:
					bu.Where(builder.Lt{k.FieldName: v.Value})
				case GTE:
					bu.Where(builder.Gte{k.FieldName: v.Value})
				case LTE:
					bu.Where(builder.Lte{k.FieldName: v.Value})
				case IN:
					bu.Where(markIn(true, k.FieldName, v.Value))
				case NI:
					bu.Where(markIn(false, k.FieldName, v.Value))
				case IsNull:
					bu.Where(builder.IsNull{k.FieldName})
				case NotNull:
					bu.Where(builder.NotNull{k.FieldName})
				default:
					bu.Where(builder.Eq{k.FieldName: v.Value})
				}
			}
		}
	}
}

func markIn(isIn bool, fieldName string, value interface{}) builder.Cond {
	var cond builder.Cond
	if vs, ok := value.(string); ok {
		value = strings.Split(vs, ",")
	}
	if isIn {
		cond = builder.In(fieldName, value)
	} else {
		cond = builder.NotIn(fieldName, value)
	}
	return cond
}
