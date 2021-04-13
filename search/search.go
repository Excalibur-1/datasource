package search

import (
	"fmt"
	"strings"

	"xorm.io/builder"
	"xorm.io/xorm"
)

type Operator int

const (
	EQ Operator = iota + 1
	NE
	LIKE
	GT
	LT
	GTE
	LTE
	IN
	NI
	IsNull
	NotNull
)

func OperatorValueOf(operator string) Operator {
	switch operator {
	case "EQ":
		return EQ
	case "NE":
		return NE
	case "LIKE":
		return LIKE
	case "GT":
		return GT
	case "LT":
		return LT
	case "GTE":
		return GTE
	case "LTE":
		return LTE
	case "IN":
		return IN
	case "NI":
		return NI
	case "IsNull":
		return IsNull
	case "NotNull":
		return NotNull
	}
	return 0
}

type Filter struct {
	FieldName string
	Value     interface{}
	Operator  Operator
}

// Parse searchParams中key的格式为OPERATOR_FIELDNAME
func Parse(searchParams map[string]interface{}) (filters []Filter) {
	filters = make([]Filter, len(searchParams))
	index := 0
	for k, v := range searchParams {
		// 过滤掉空值
		if v == nil {
			continue
		}
		// 拆分operator与filedAttribute
		names := strings.Split(k, "_")
		if len(names) != 2 {
			fmt.Printf("%s is not a valid search filter name\n", k)
			continue
		}
		filedName := names[1]
		operator := OperatorValueOf(names[0])
		filters[index] = Filter{FieldName: filedName, Value: v, Operator: operator}
		index++
	}
	return
}

func OrmFilter(filters []Filter, orm *xorm.Session) *xorm.Session {
	for _, filter := range filters {
		switch filter.Operator {
		case EQ:
			orm.Where(builder.Eq{filter.FieldName: filter.Value})
		case NE:
			orm.Where(builder.Neq{filter.FieldName: filter.Value})
		case LIKE:
			orm.Where(builder.Like{filter.FieldName, filter.Value.(string)})
		case GT:
			orm.Where(builder.Gt{filter.FieldName: filter.Value})
		case LT:
			orm.Where(builder.Lt{filter.FieldName: filter.Value})
		case GTE:
			orm.Where(builder.Gte{filter.FieldName: filter.Value})
		case LTE:
			orm.Where(builder.Lte{filter.FieldName: filter.Value})
		case IN:
			orm.Where(builder.In(filter.FieldName, filter.Value))
		case NI:
			orm.Where(builder.NotIn(filter.FieldName, filter.Value))
		case IsNull:
			orm.Where(builder.IsNull{filter.FieldName})
		case NotNull:
			orm.Where(builder.NotNull{filter.FieldName})
		}
	}
	return orm
}

func BuilderFilter(filters []Filter, bu *builder.Builder) {
	for _, filter := range filters {
		switch filter.Operator {
		case EQ:
			bu.Where(builder.Eq{filter.FieldName: filter.Value})
		case NE:
			bu.Where(builder.Neq{filter.FieldName: filter.Value})
		case LIKE:
			bu.Where(builder.Like{filter.FieldName, filter.Value.(string)})
		case GT:
			bu.Where(builder.Gt{filter.FieldName: filter.Value})
		case LT:
			bu.Where(builder.Lt{filter.FieldName: filter.Value})
		case GTE:
			bu.Where(builder.Gte{filter.FieldName: filter.Value})
		case LTE:
			bu.Where(builder.Lte{filter.FieldName: filter.Value})
		case IN:
			bu.Where(builder.In(filter.FieldName, filter.Value))
		case NI:
			bu.Where(builder.NotIn(filter.FieldName, filter.Value))
		case IsNull:
			bu.Where(builder.IsNull{filter.FieldName})
		case NotNull:
			bu.Where(builder.NotNull{filter.FieldName})
		}
	}
}
