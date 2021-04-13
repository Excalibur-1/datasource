package common

import (
	"encoding/json"
	"github.com/Excalibur-1/datasource/pb"
)

// Query {"pageSize":10,"page":0,"sorted":[{"id":"firstName","desc":false}],"filtered":[{"id":"firstName","value":"3"}]}
type Query struct {
	PageSize int32 `json:"pageSize"`
	Page     int32 `json:"page"`
	Sorted   []struct {
		Id   string `json:"id"`
		Desc bool   `json:"desc"`
	} `json:"sorted"`
	Filtered []struct {
		Id    string      `json:"id"`
		Value interface{} `json:"value"`
	} `json:"filtered"`
}

func (sp *Query) SetPageSize(pageSize int32) {
	sp.PageSize = pageSize
}
func (sp *Query) SetPage(page int32) {
	sp.Page = page
}
func (sp *Query) SetSorted(id string, desc bool) {
	sp.Sorted = append(sp.Sorted, struct {
		Id   string `json:"id"`
		Desc bool   `json:"desc"`
	}{Id: id, Desc: desc})
}
func (sp *Query) SetFiltered(id string, value interface{}) {
	sp.Filtered = append(sp.Filtered, struct {
		Id    string      `json:"id"`
		Value interface{} `json:"value"`
	}{Id: id, Value: value})
}
func (sp Query) MarkPage() *Pagination {
	if sp.PageSize <= 0 {
		sp.PageSize = 20
	}
	if sp.Page <= 0 {
		sp.Page = 1
	}
	return NewPagination(int(sp.Page), int(sp.PageSize), 0)
}
func (sp *Query) ToPb() (pq *pb.Query) {
	pq = &pb.Query{PageSize: sp.PageSize, Page: sp.Page}
	if len(sp.Sorted) > 0 {
		pq.Sorted, _ = json.Marshal(sp.Sorted)
	}
	if len(sp.Filtered) > 0 {
		pq.Filtered, _ = json.Marshal(sp.Filtered)
	}
	return
}
func (sp *Query) ForPb(pq *pb.Query) *Query {
	sp.PageSize = pq.PageSize
	sp.Page = pq.Page
	var sort []struct {
		Id   string `json:"id"`
		Desc bool   `json:"desc"`
	}
	if err := json.Unmarshal(pq.Sorted, &sort); err == nil {
		sp.Sorted = sort
	}
	var filter []struct {
		Id    string      `json:"id"`
		Value interface{} `json:"value"`
	}
	if err := json.Unmarshal(pq.Filtered, &filter); err == nil {
		sp.Filtered = filter
	}
	return sp
}
