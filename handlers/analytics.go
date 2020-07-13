package handlers

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"persons/config"
)

// записи
type CountRecord struct {
	FullTitle        string `json:"full_title"`
	SumRecordApi     int64  `json:"sum_record_api"`
	SumRecordCabinet int64  `json:"sum_record_cabinet"`
}

// количество отправленных пакетов через АПИ
type CountPackage struct {
	FullTitle       string `json:"full_title"`
	CountAllAdd     int64  `json:"count_all_add"  gorm:"column:count_all_add"`
	CountAllRemove  int64  `json:"count_all_remove"  gorm:"column:count_all_remove"`
	CountTrueAdd    int64  `json:"count_true_add"  gorm:"column:count_true_add"`
	CountTrueRemove int64  `json:"count_true_remove"  gorm:"column:count_true_remove"`
}

//заявления и абитуриенты по организациям
type CountApp struct {
	FullTitle string `json:"full_title"`
	Count     int64  `json:"count"`
	Count3    int64  `json:"count_3" gorm:"column:count_3"`
	Count4    int64  `json:"count_4" gorm:"column:count_4"`
	Count5    int64  `json:"count_5" gorm:"column:count_5"`
	Count10   int64  `json:"count_10" gorm:"column:count_10"`
	CountAll  int64  `json:"count_all"`
}

func (result *ResultInfo) GetAnalytics() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	f := excelize.NewFile()
	// Количество записей 1 запрос
	var countRecords []CountRecord
	cmdCountRecords := `with a as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Приемные кампании'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api
		, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.campaigns c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title),
		a1 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title,
		'Конкурсные группы'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api,
		count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.competitive_groups c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), 
		a2 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title,
		'Идивидуальные достижения'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api
		, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.achievements c
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a3 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Объем приема'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api
		, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.admission_volume c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a4 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Формы приемных кампаний'::text as data_type, 0 as count_api, count(c.id) as count_cabinet from cmp.campaigns_educ_form c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a5 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Уровни приемных кампаний'::text as data_type, 0 as count_api, count(c.id) as count_cabinet from cmp.campaigns_educ_level c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a6 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, 'Образовательные программы в конкурсах'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.competitive_group_programs c join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), a7 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, 'Вступительные испытания'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.entrance_test c join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), a8 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, 'Объем приема в соответствии с уровнем бюджета'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.distributed_admission_volume c join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), z as (select * from (select * from a union all select * from a1 union all select * from a2 union all select * from a3 union all select * from a4 union all select * from a5 union all select * from a6 union all select * from a7 union all select * from a8) aa where full_title <> 'Полная тестовая')
		, z1 as (select full_title, sum(count_api) as sum_record_api
		, sum(count_cabinet) as sum_record_cabinet from z group by full_title order by full_title)
		, z2 as (select 'Общее количество'::character varying(1000), sum(sum_record_api) as sum1, sum(sum_record_cabinet) as sum2 from z1 ) 
		select * from z1 UNION ALL select * from z2;`
	db := conn.Raw(cmdCountRecords).Scan(&countRecords)
	if db.Error != nil {
		m := db.Error.Error()
		result.Message = &m
	} else {
		// Create a new sheet.
		sheet := "Количество записей"
		index := f.NewSheet(sheet)
		styleTitle, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#CCFFFF"],"pattern":1}}`)

		f.SetColWidth(sheet, "A", "A", 100)
		f.SetColWidth(sheet, "B", "B", 50)
		f.SetColWidth(sheet, "C", "C", 50)
		// Set title
		f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", 1), fmt.Sprintf(`%v%d`, "C", 1), styleTitle)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", 1), `Организация`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", 1), `Количество записей внесенных через API`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", 1), `Количество записей внесенных через кабинет`)
		for i, value := range countRecords {
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", i+2), value.FullTitle)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", i+2), value.SumRecordApi)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", i+2), value.SumRecordCabinet)
			if i+1 == len(countRecords) {
				styleSum, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["##C1F0B6"],"pattern":1}}`)
				f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", i+2), fmt.Sprintf(`%v%d`, "C", i+2), styleSum)
			}
		}
		// Set active sheet of the workbook.
		f.SetActiveSheet(index)
	}

	// количество отправленных пакетов через АПИ 2 запрос
	var countPackages []CountPackage
	cmdCountPackages := `with a as
(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, count(rh.action) filter (where rh.action = 'Add') as count_all_add,
count(rh.action) filter (where rh.action = 'Remove') as count_all_remove,
count(rh.action) filter (where rh.action = 'Add' and rq.id_status = 3) as count_true_add,
count(rh.action) filter (where rh.action = 'Remove' and rq.id_status = 3) as count_true_remove
from admin.jwt_request_header rh
join admin.organizations o ON o.ogrn = rh.ogrn and o.kpp = rh.kpp
join admin.jwt_request_queue rq ON rq.id = rh.id_jwt 
where rh.created_at >= ? and rh.action in ('Add', 'Remove') 
group by o.full_title order by o.full_title),
b as (select 'Общее количество'::character varying(1000), sum(count_all_add), sum(count_all_remove),
sum(count_true_add), sum(count_true_remove)
from a)
select * from a UNION ALL select * from b;`
	db2 := conn.Raw(cmdCountPackages, `2020-06-02`).Scan(&countPackages)
	if db2.Error != nil {
		m := db2.Error.Error()
		result.Message = &m
	} else {
		// Create a new sheet.
		sheet := "Количество пакетов"
		_ = f.NewSheet(sheet)
		styleTitle, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#fbd1ff"],"pattern":1}}`)
		f.SetColWidth(sheet, "A", "A", 100)
		f.SetColWidth(sheet, "B", "B", 30)
		f.SetColWidth(sheet, "C", "C", 30)
		f.SetColWidth(sheet, "D", "D", 30)
		f.SetColWidth(sheet, "E", "E", 30)
		// Set title
		f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", 1), fmt.Sprintf(`%v%d`, "E", 1), styleTitle)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", 1), `Организация`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", 1), `Общее количество пакетов на добавление`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", 1), `Общее количество пакетов на удаление`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", 1), `Количество успешных пакетов на добавление`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", 1), `Количество успешных пакетов на удаление`)
		for i, value := range countPackages {
			if value.FullTitle == `Полная тестовая` {
				fmt.Println(`**********`)
				fmt.Println(value.CountAllAdd)
				fmt.Println(value.CountAllRemove)
				fmt.Println(value.CountTrueAdd)
				fmt.Println(value.CountTrueRemove)
				fmt.Println(`**********`)
			}
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", i+2), value.FullTitle)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", i+2), value.CountAllAdd)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", i+2), value.CountAllRemove)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", i+2), value.CountTrueAdd)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", i+2), value.CountTrueRemove)
			if i+1 == len(countPackages) {
				styleSum, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["##C1F0B6"],"pattern":1}}`)
				f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", i+2), fmt.Sprintf(`%v%d`, "E", i+2), styleSum)
			}
		}
	}
	// количество заявлений 3 запрос
	var countApps []CountApp
	cmdCountApps := `with a as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title,-- название ОО
count(distinct id_entrant),  -- Количество абитуриентов
count(ap.id) as count_all, -- Общее количество заявлений
count(ap.id) filter (where ap.id_status = 3) as count_3, -- Отправка в ООВО
count(ap.id) filter (where ap.id_status = 4) as count_4, -- Доставлено в ООВО
count(ap.id) filter (where ap.id_status = 5) as count_5, -- Заявление принято в ООВО
count(ap.id) filter (where ap.id_status = 10) as count_10 -- Заявление отклонено
 from admin.organizations o
join app.applications ap ON ap.id_organization = o.id
group by o.full_title order by o.full_title),
a1 as (select 'Общее количество'::text AS full_title,
count(distinct id_entrant),
count(ap.id) as count_all, -- Общее количество
count(ap.id) filter (where ap.id_status = 3) as count_3, -- Отправка в ООВО
count(ap.id) filter (where ap.id_status = 4) as count_4, -- Доставлено в ООВО
count(ap.id) filter (where ap.id_status = 5) as count_5, -- Заявление принято в ООВО
count(ap.id) filter (where ap.id_status = 10) as count_10 -- Заявление отклонено
 from admin.organizations o
join app.applications ap ON ap.id_organization = o.id)
select * from a union all select * from a1;`
	db3 := conn.Raw(cmdCountApps).Scan(&countApps)
	if db3.Error != nil {
		m := db3.Error.Error()
		result.Message = &m
	} else {
		// Create a new sheet.
		sheet := "Количество заявлений"
		_ = f.NewSheet(sheet)
		styleTitle, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#FFFBD1"],"pattern":1}}`)
		f.SetColWidth(sheet, "A", "A", 100)
		f.SetColWidth(sheet, "B", "B", 30)
		f.SetColWidth(sheet, "C", "C", 30)
		f.SetColWidth(sheet, "D", "D", 30)
		f.SetColWidth(sheet, "E", "E", 30)
		f.SetColWidth(sheet, "F", "F", 30)
		f.SetColWidth(sheet, "G", "G", 30)
		// Set title
		f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", 1), fmt.Sprintf(`%v%d`, "G", 1), styleTitle)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", 1), `Организация`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", 1), `Количество абитуриентов`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", 1), `Общее количество заявлений`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", 1), `Отправка в ООВО`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", 1), `Доставлено в ООВО`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", 1), `Заявление принято в ООВО`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", 1), `Заявление отклонено`)
		for i, value := range countApps {

			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", i+2), value.FullTitle)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", i+2), value.Count)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", i+2), value.CountAll)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", i+2), value.Count3)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", i+2), value.Count4)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", i+2), value.Count5)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", i+2), value.Count10)
			if i+1 == len(countApps) {
				styleSum, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["##C1F0B6"],"pattern":1}}`)
				f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", i+2), fmt.Sprintf(`%v%d`, "G", i+2), styleSum)
			}
		}
	}

	// Save xlsx file by the given path.
	path := `uploads/Book1.xlsx`
	if err := f.SaveAs(path); err != nil {
		m := err.Error()
		result.Message = &m
		return
	}
	result.Done = true
	result.Items = path
	return

}
func (result *ResultInfo) GetAnalyticsTest() {
	result.Done = false
	items := make(map[string]interface{})
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	// Количество записей 1 запрос
	var countRecords []CountRecord
	cmdCountRecords := `with a as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Приемные кампании'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api
		, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.campaigns c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title),
		a1 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title,
		'Конкурсные группы'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api,
		count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.competitive_groups c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), 
		a2 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title,
		'Идивидуальные достижения'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api
		, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.achievements c
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a3 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Объем приема'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api
		, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.admission_volume c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a4 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Формы приемных кампаний'::text as data_type, 0 as count_api, count(c.id) as count_cabinet from cmp.campaigns_educ_form c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a5 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title
		, 'Уровни приемных кампаний'::text as data_type, 0 as count_api, count(c.id) as count_cabinet from cmp.campaigns_educ_level c 
		join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title)
		, a6 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, 'Образовательные программы в конкурсах'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.competitive_group_programs c join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), a7 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, 'Вступительные испытания'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.entrance_test c join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), a8 as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, 'Объем приема в соответствии с уровнем бюджета'::text as data_type, count(c.id) filter (where c.id_author is null OR c.id_author = 1317) as count_api, count(c.id) filter (where c.id_author is not null AND c.id_author <> 1317) as count_cabinet from cmp.distributed_admission_volume c join admin.organizations o ON o.id = c.id_organization group by o.full_title order by o.full_title), z as (select * from (select * from a union all select * from a1 union all select * from a2 union all select * from a3 union all select * from a4 union all select * from a5 union all select * from a6 union all select * from a7 union all select * from a8) aa where full_title <> 'Полная тестовая')
		, z1 as (select full_title, sum(count_api) as sum_record_api
		, sum(count_cabinet) as sum_record_cabinet from z group by full_title order by full_title)
		, z2 as (select 'Общее количество'::character varying(1000), sum(sum_record_api) as sum1, sum(sum_record_cabinet) as sum2 from z1 ) 
		select * from z1 UNION ALL select * from z2;`
	db := conn.Raw(cmdCountRecords).Scan(&countRecords)
	if db.Error != nil {
		m := db.Error.Error()
		result.Message = &m
	} else {
		var r []interface{}
		for _, value := range countRecords {
			r = append(r, map[string]interface{}{
				`full_title`:       value.FullTitle,
				`SumRecordApi`:     value.SumRecordApi,
				`SumRecordCabinet`: value.SumRecordCabinet,
			})
		}
		items[`countRecords`] = r
	}

	// количество отправленных пакетов через АПИ 2 запрос
	var countPackages []CountPackage
	cmdCountPackages := `with a as
(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, count(rh.action) filter (where rh.action = 'Add') as count_all_add,
count(rh.action) filter (where rh.action = 'Remove') as count_all_remove,
count(rh.action) filter (where rh.action = 'Add' and rq.id_status = 3) as count_true_add,
count(rh.action) filter (where rh.action = 'Remove' and rq.id_status = 3) as count_true_remove
from admin.jwt_request_header rh
join admin.organizations o ON o.ogrn = rh.ogrn and o.kpp = rh.kpp
join admin.jwt_request_queue rq ON rq.id = rh.id_jwt
where rh.created_at >= '2.06.2020' and action in ('Add', 'Remove')
group by o.full_title order by o.full_title),
b as (select 'Общее количество'::character varying(1000), sum(count_all_add), sum(count_all_remove),
sum(count_true_add), sum(count_true_remove)
from a)
select * from a UNION ALL select * from b;`
	db2 := conn.Raw(cmdCountPackages).Scan(&countPackages)
	if db2.Error != nil {
		m := db2.Error.Error()
		result.Message = &m
	} else {
		// Create a new sheet.
		var r []interface{}
		for _, value := range countPackages {
			r = append(r, map[string]interface{}{
				`full_title`:      value.FullTitle,
				`CountAllAdd`:     value.CountAllAdd,
				`CountAllRemove`:  value.CountAllRemove,
				`CountTrueAdd`:    value.CountTrueAdd,
				`CountTrueRemove`: value.CountTrueRemove,
			})
		}
		items[`countPackages`] = r
	}
	// количество заявлений 3 запрос
	var countApps []CountApp
	cmdCountApps := `with a as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title,-- название ОО
count(distinct id_entrant),  -- Количество абитуриентов
count(ap.id) as count_all, -- Общее количество заявлений
count(ap.id) filter (where ap.id_status = 3) as count_3, -- Отправка в ООВО
count(ap.id) filter (where ap.id_status = 4) as count_4, -- Доставлено в ООВО
count(ap.id) filter (where ap.id_status = 5) as count_5, -- Заявление принято в ООВО
count(ap.id) filter (where ap.id_status = 10) as count_10 -- Заявление отклонено
 from admin.organizations o
join app.applications ap ON ap.id_organization = o.id
group by o.full_title order by o.full_title),
a1 as (select 'Общее количество'::text AS full_title,
count(distinct id_entrant),
count(ap.id) as count_all, -- Общее количество
count(ap.id) filter (where ap.id_status = 3) as count_3, -- Отправка в ООВО
count(ap.id) filter (where ap.id_status = 4) as count_4, -- Доставлено в ООВО
count(ap.id) filter (where ap.id_status = 5) as count_5, -- Заявление принято в ООВО
count(ap.id) filter (where ap.id_status = 10) as count_10 -- Заявление отклонено
 from admin.organizations o
join app.applications ap ON ap.id_organization = o.id)
select * from a union all select * from a1;`
	db3 := conn.Raw(cmdCountApps).Scan(&countApps)
	if db3.Error != nil {
		m := db3.Error.Error()
		result.Message = &m
	} else {
		var r []interface{}
		for _, value := range countApps {
			r = append(r, map[string]interface{}{
				`full_title`: value.FullTitle,
				`Count`:      value.Count,
				`CountAll`:   value.CountAll,
				`Count3`:     value.Count3,
				`Count4`:     value.Count4,
				`Count5`:     value.Count5,
				`Count10`:    value.Count10,
			})
		}
		items[`countApps`] = r
	}

	// Save xlsx file by the given path.

	result.Done = true
	result.Items = items
	return

}
