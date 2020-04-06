package service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"persons/config"
	"strconv"
)

type Paginator struct {
	TotalCount    int                    `json:"total"`
	AllPage       int                    `json:"count_page"`
	CurrentPage   int                    `json:"page"`
	Offset        int                    `json:"offset"`
	Limit         int                    `json:"limit"`
	Field         string                 `json:"-"`
	Cmd           string                 `json:"-"`
	CmdPag        string                 `json:"-"`
	Args          map[string]interface{} `json:"-"`
	Configuration *config.Configuration  `json:"-"`
}

func (paginator *Paginator) UsePaginator() {
	db, err := sqlx.Connect("postgres", paginator.Configuration.DbString)
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println(`ошибка ` + err.Error())
		}
	}()
	if err != nil {
		fmt.Println(err)
	}
	row, err := db.NamedQuery("SELECT count("+paginator.Field+") "+paginator.Cmd, paginator.Args)
	if err == nil {
		defer func() {
			err := row.Close()
			if err != nil {
				fmt.Println(`ошибка ` + err.Error())
			}
		}()
	}
	if err != nil {
		fmt.Println(err)
	}

	if row.Next() {
		err = row.Scan(&paginator.TotalCount)
		if err != nil {
			fmt.Println(err)
		}
	}

	paginator.AllPage = paginator.TotalCount / paginator.Limit
	if (paginator.TotalCount % paginator.Limit) > 0 {
		paginator.AllPage++
	}
	paginator.Offset = (paginator.CurrentPage - 1) * paginator.Limit
	paginator.CmdPag = " LIMIT " + strconv.Itoa(paginator.Limit) + " OFFSET " + strconv.Itoa(paginator.Offset)
}
