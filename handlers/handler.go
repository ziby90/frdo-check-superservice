package handlers

import (
	"math"
	"strconv"
)

type Paginator struct {
	TotalCount  int `json:"total"`
	AllPage     int `json:"count_page"`
	CurrentPage int `json:"page"`
	Offset      int `json:"offset"`
	Limit       int `json:"limit"`
}

type Sort struct {
	Field string `json:"sortby"`
	Order string `json:"order"`
}

type Result struct {
	Done      bool        `json:"done"`
	Message   *string     `json:"message, omitempty"`
	Items     interface{} `json:"data"`
	Paginator Paginator   `json:"paginator"`
	Sort      Sort        `json:"sort"`
	Search    string      `json:"search"`
}

type ResultInfo struct {
	Done    bool        `json:"done"`
	Message *string     `json:"message, omitempty"`
	Items   interface{} `json:"data"`
}

type ResultCls struct {
	Done    bool        `json:"done"`
	Message *string     `json:"message, omitempty"`
	Items   interface{} `json:"data"`
	Search  string      `json:"search, 	omitempty"`
}

func NewResult() Result {
	res := Result{
		Done:    false,
		Message: nil,
		Items:   nil,
		Paginator: Paginator{
			TotalCount:  0,
			AllPage:     0,
			CurrentPage: 1,
			Offset:      0,
			Limit:       10,
		},
		Sort: Sort{
			Field: "created",
			Order: "ASC",
		},
		Search: "",
	}
	return res
}

func (result *Result) MakeUrlParams(keys map[string][]string) {
	if len(keys[`limit`]) > 0 {
		if v, ok := strconv.Atoi(keys[`limit`][0]); ok == nil {
			if v > 0 && v < 50 {
				result.Paginator.Limit = v
			}
		}
	}
	if len(keys[`page`]) > 0 {
		if v, ok := strconv.Atoi(keys[`page`][0]); ok == nil {
			if v > 0 && v < 50 {
				result.Paginator.CurrentPage = v
			}
		}
	}
	if len(keys[`search`]) > 0 {
		result.Search = keys[`search`][0]
	}
	if len(keys[`order`]) > 0 {
		result.Sort.Order = keys[`order`][0]
	}
	if len(keys[`sortby`]) > 0 {
		result.Sort.Field = keys[`sortby`][0]
	}
}

func (result *ResultCls) MakeUrlParams(keys map[string][]string) {
	if len(keys[`search`]) > 0 {
		result.Search = keys[`search`][0]
	}
}

func (paginator *Paginator) Make() {
	if paginator.TotalCount <= 0 {
		return
	}
	paginator.AllPage = int(math.Ceil(float64(paginator.TotalCount) / float64(paginator.Limit)))
	paginator.Offset = (paginator.CurrentPage - 1) * paginator.Limit
}
