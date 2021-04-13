package sort

import (
	"strings"
)

type Direction int

const (
	ASC Direction = iota + 1
	DESC
	DefaultDirection = ASC
)

func (d Direction) Ascending() bool {
	return d == ASC
}
func (d Direction) Descending() bool {
	return d == DESC
}

func (d Direction) ToString() string {
	direction := []string{"", "ASC", "DESC"}
	return direction[d]
}

var Ordered = Order{}

type Order struct {
	direction Direction
	property  string
}

func (o Order) By(property string) Order {
	return Order{direction: DefaultDirection, property: property}
}

func (o Order) ByProperties(direction Direction, properties ...string) (orders []Order) {
	orders = make([]Order, len(properties))
	for _, v := range properties {
		orders = append(orders, Order{direction: direction, property: v})
	}
	return
}

func (o Order) Asc(property string) Order {
	return Order{direction: ASC, property: property}
}

func (o Order) Desc(property string) Order {
	return Order{direction: DESC, property: property}
}

type Sort struct {
	orders []Order
}

func Sorted() *Sort {
	return &Sort{orders: make([]Order, 0)}
}

func (s *Sort) Reset() {
	s.orders = s.orders[:0]
}

func (s *Sort) By(property string) *Sort {
	s.orders = append(s.orders, Ordered.By(property))
	return s
}

func (s *Sort) ByProperties(direction Direction, properties ...string) *Sort {
	for _, v := range properties {
		s.orders = append(s.orders, Order{direction: direction, property: v})
	}
	return s
}

func (s *Sort) ByOrder(orders ...Order) *Sort {
	s.orders = append(s.orders, orders...)
	return s
}

func (s *Sort) Asc(property string) *Sort {
	s.orders = append(s.orders, Ordered.Asc(property))
	return s
}

func (s *Sort) Desc(property string) *Sort {
	s.orders = append(s.orders, Ordered.Desc(property))
	return s
}

func (s *Sort) ToString() (str string) {
	// faster than bytes.Buffer
	var asc, desc strings.Builder
	for _, v := range s.orders {
		if v.direction == DESC {
			desc.WriteString(v.property)
			desc.WriteString(",")
		} else {
			asc.WriteString(v.property)
			asc.WriteString(",")
		}
	}
	if desc.Len() > 1 {
		str += desc.String()
		str = str[:len(str)-1]
		str += " DESC"
	}
	if asc.Len() > 1 {
		if desc.Len() > 1 {
			str += ","
		}
		str += asc.String()
		str = str[:len(str)-1]
		str += " ASC"
	}
	return
}

func (s *Sort) FirstAscString() (str string) {
	// faster than bytes.Buffer
	var asc, desc strings.Builder
	for _, v := range s.orders {
		if v.direction == DESC {
			desc.WriteString(v.property)
			desc.WriteString(",")
		} else {
			asc.WriteString(v.property)
			asc.WriteString(",")
		}
	}
	if asc.Len() > 1 {
		str += asc.String()
		str = str[:len(str)-1]
		str += " ASC"
	}
	if desc.Len() > 1 {
		if asc.Len() > 1 {
			str += ","
		}
		str += desc.String()
		str = str[:len(str)-1]
		str += " DESC"
	}
	return
}
