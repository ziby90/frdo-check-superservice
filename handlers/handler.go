package handlers

import "strconv"

type Paginator struct {
	TotalCount    int                    `json:"total"`
	AllPage       int                    `json:"count_page"`
	CurrentPage   int                    `json:"page"`
	Offset        int                    `json:"offset"`
	Limit         int                    `json:"limit"`
}

type Sort struct {
	Field 	string 		`json:"field_sort"`
	Order 	string 		`json:"field_order"`
}


type Params struct {
	Paginator 	Paginator
	Sort		Sort
	Search 		string `json:"search"`
}

type Result struct {
	Done 		bool			`json:"done"`
	Message 	*string			`json:"message, omitempty"`
	Items 		interface{}		`json:"data"`
	Params		Params 			`json:"params"`
}

type ResultCls struct {
	Done 		bool			`json:"done"`
	Message 	*string			`json:"message, omitempty"`
	Items 		interface{}		`json:"data"`
}

func NewResult() Result{
	res := Result{
		Done:    false,
		Message: nil,
		Items:   nil,
		Params:  Params{
			Paginator: Paginator{
				TotalCount:  0,
				AllPage:     0,
				CurrentPage: 0,
				Offset:      0,
				Limit:       10,
			},
			Sort:      Sort{
				Field: "created",
				Order: "ASC",
			},
			Search:  "",
		},
	}
	return res
}

func (result *Result) MakeUrlParams(keys map[string][]string){
		if len(keys[`limit`])>0 {
			if v, ok := strconv.Atoi(keys[`limit`][0]); ok ==nil{
				if v>0 && v<50 {
					result.Params.Paginator.Limit = v
				}
			}
		}
}
