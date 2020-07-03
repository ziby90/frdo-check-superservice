package handlers

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"persons/config"
)

type CountRecord struct {
	FullTitle        string `json:"full_title"`
	SumRecordApi     int64  `json:"sum_record_api"`
	SumRecordCabinet int64  `json:"sum_record_cabinet"`
}
type CountPackage struct {
	FullTitle        string `json:"full_title"`
	SumRecordApi     int64  `json:"sum_record_api"`
	SumRecordCabinet int64  `json:"sum_record_cabinet"`
}
type CountApp struct {
	FullTitle string `json:"full_title"`
	CountApp  int64  `json:"count_app"`
}

func (result *ResultInfo) GetAnalytics() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	f := excelize.NewFile()
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

	var countApps []CountApp
	cmdCountApps := `with a as(select btrim(regexp_replace(o.full_title::text, '[^А-яA-zЁё\- \d\.]+'::text, ''::text, 'g'::text)) AS full_title, count(ap.*) as count_app from app.applications ap
join admin.organizations o ON o.id = ap.id_organization
group by o.full_title order by o.full_title),
b as (select 'Общее количество'::text, sum(count_app) from a)
select * from a UNION ALL select * from b;`
	db = conn.Raw(cmdCountApps).Scan(&countApps)
	if db.Error != nil {
		m := db.Error.Error()
		result.Message = &m
	} else {
		// Create a new sheet.
		sheet := "Количество заявлений"
		_ = f.NewSheet(sheet)
		styleTitle, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#FFFBD1"],"pattern":1}}`)
		f.SetColWidth(sheet, "A", "A", 100)
		f.SetColWidth(sheet, "B", "B", 50)
		// Set title
		f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", 1), fmt.Sprintf(`%v%d`, "B", 1), styleTitle)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", 1), `Организация`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", 1), `Количество заявлений`)
		for i, value := range countApps {
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", i+2), value.FullTitle)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", i+2), value.CountApp)
			if i+1 == len(countApps) {
				styleSum, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["##C1F0B6"],"pattern":1}}`)
				f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", i+2), fmt.Sprintf(`%v%d`, "B", i+2), styleSum)
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
