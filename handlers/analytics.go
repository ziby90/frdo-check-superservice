package handlers

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"persons/config"
	"time"
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

// список заявлений
type AnalApplications struct { // analytics applications, конечно же
	AppNumber            string    `json:"app_number" gorm:"column:app_number"`
	RegistrationDate     time.Time `json:"registration_date"  gorm:"column:registration_date"`
	Surname              string    `json:"surname"  gorm:"column:surname"`
	Name                 string    `json:"name"  gorm:"column:name"`
	Patronymic           string    `json:"patronymic"  gorm:"column:patronymic"`
	DocSeries            string    `json:"doc_series"  gorm:"column:doc_series"`
	DocNumber            string    `json:"doc_number"  gorm:"column:doc_number"`
	CompetitiveGroupName string    `json:"competitive_group_name"  gorm:"column:competitive_group_name"`
	CompetitiveGroupUid  string    `json:"competitive_group_uid"  gorm:"column:competitive_group_uid"`
	UidEpgu              string    `json:"uid_epgu"  gorm:"column:uid_epgu"`
	Uid                  string    `json:"uid"  gorm:"column:uid"`
}

// список конкурсных групп
type AnalCompetitiveGroups struct { // analytics competitive groups, конечно же
	CompetitiveGroupName string `json:"competitive_group_name" gorm:"column:competitive_group_name"`
	Uid                  string `json:"uid"  gorm:"column:uid"`
	CampaignName         string `json:"campaign_name" gorm:"column:campaign_name"`
	EducationForm        string `json:"education_form" gorm:"column:education_form"`
	EducationSource      string `json:"education_source" gorm:"column:education_source"`
	EducationLevel       string `json:"education_level" gorm:"column:education_level"`
	LevelBudget          string `json:"level_budget" gorm:"column:level_budget"`
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
join app.applications ap ON ap.id_organization = o.id WHERE ap.uid_epgu IS NOT NULL
group by o.full_title order by o.full_title),
a1 as (select 'Общее количество'::text AS full_title,
count(distinct id_entrant),
count(ap.id) as count_all, -- Общее количество
count(ap.id) filter (where ap.id_status = 3) as count_3, -- Отправка в ООВО
count(ap.id) filter (where ap.id_status = 4) as count_4, -- Доставлено в ООВО
count(ap.id) filter (where ap.id_status = 5) as count_5, -- Заявление принято в ООВО
count(ap.id) filter (where ap.id_status = 10) as count_10 -- Заявление отклонено
 from admin.organizations o
join app.applications ap ON ap.id_organization = o.id WHERE ap.uid_epgu IS NOT NULL)
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
func (result *ResultInfo) GetAnalyticsListApplications() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	f := excelize.NewFile()
	// Количество записей 1 запрос
	var items []AnalApplications
	cmd := `SELECT a.app_number
			, a.registration_date
			, a.uid
			, a.id ,i.surname
			, i.name
			, i.patronymic
			, i.doc_series
			, i.doc_number
			, a.uid_epgu
			, cg.name as competitive_group_name
			, cg.uid as competitive_group_uid
			FROM app.applications a
			JOIN documents.identification i ON i.id_entrant = a.id_entrant
			JOIN cmp.competitive_groups cg ON cg.id = a.id_competitive_group
			WHERE  a.id_organization = ? and cg.actual IS true AND a.actual IS true
			ORDER BY a.registration_date;`
	db := conn.Raw(cmd, result.User.CurrentOrganization.Id).Scan(&items)
	if db.Error != nil {
		m := db.Error.Error()
		result.Message = &m
	} else {
		// Create a new sheet.
		sheet := "Список заявлений"
		index := f.NewSheet(sheet)
		styleTitle, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#CCFFFF"],"pattern":1}}`)

		f.SetColWidth(sheet, "A", "K", 50)
		// Set title
		f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", 1), fmt.Sprintf(`%v%d`, "K", 1), styleTitle)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", 1), `Номер заявления`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", 1), `Дата регистрации`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", 1), `Фамилия`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", 1), `Имя`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", 1), `Отчество`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", 1), `Серия документа`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", 1), `Номер документа`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "H", 1), `UID ЕПГУ заявления`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "I", 1), `UID заявления`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "J", 1), `Название конкурсной группы`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "K", 1), `UID конкурсной группы`)
		for i, value := range items {
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", i+2), value.AppNumber)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", i+2), value.RegistrationDate)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", i+2), value.Surname)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", i+2), value.Name)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", i+2), value.Patronymic)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", i+2), value.DocSeries)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", i+2), value.DocNumber)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "H", i+2), value.UidEpgu)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "I", i+2), value.Uid)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "J", i+2), value.CompetitiveGroupName)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "K", i+2), value.CompetitiveGroupUid)
			//if i+1 == len(applications) {
			//	styleSum, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["##C1F0B6"],"pattern":1}}`)
			//	f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", i+2), fmt.Sprintf(`%v%d`, "I", i+2), styleSum)
			//}
		}
		// Set active sheet of the workbook.
		f.SetActiveSheet(index)
	}

	// Save xlsx file by the given path.

	path := fmt.Sprintf(`uploads/%v_anal_applications.xlsx`, result.User.CurrentOrganization.Id)
	if err := f.SaveAs(path); err != nil {
		m := err.Error()
		result.Message = &m
		return
	}
	result.Done = true
	result.Items = path
	return

}
func (result *ResultInfo) GetAnalyticsListCompetitiveGroup() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	f := excelize.NewFile()
	// Количество записей 1 запрос
	var items []AnalCompetitiveGroups
	cmd := `SELECT 
		cg.name as competitive_group_name
		, ef.name as education_form
		, cg.uid as uid
		, el.name as education_level
		, es.name as education_source
		, lb.name as level_budget
		, c.name as campaign_name 
		FROM cmp.competitive_groups cg
		JOIN cmp.campaigns c ON c.id = cg.id_campaign
		join cls.education_forms ef on ef.id= cg.id_education_form 
		join cls.education_levels el on el.id= cg.id_education_level 
		join cls.education_sources es on es.id= cg.id_education_source 
		join cls.level_budget lb on lb.id= cg.id_level_budget 
		 WHERE cg.id_organization = ?
		AND cg.actual is true AND c.actual IS true ORDER BY cg.name, c.name`
	db := conn.Raw(cmd, result.User.CurrentOrganization.Id).Scan(&items)
	if db.Error != nil {
		m := db.Error.Error()
		result.Message = &m
	} else {
		// Create a new sheet.
		sheet := "Список конкурсных групп"
		index := f.NewSheet(sheet)
		styleTitle, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#CCFFFF"],"pattern":1}}`)

		f.SetColWidth(sheet, "A", "F", 50)
		// Set title
		f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", 1), fmt.Sprintf(`%v%d`, "F", 1), styleTitle)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", 1), `Название конкурсной группы`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", 1), `UID`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", 1), `Название приемной компании`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", 1), `Уровень образования`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", 1), `Источник финансирования`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", 1), `Форма образования`)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", 1), `Уроовень бюджета`)
		for i, value := range items {
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", i+2), value.CompetitiveGroupName)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", i+2), value.Uid)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", i+2), value.CampaignName)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", i+2), value.EducationLevel)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", i+2), value.EducationSource)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", i+2), value.EducationForm)
			f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", i+2), value.LevelBudget)
			//if i+1 == len(applications) {
			//	styleSum, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["##C1F0B6"],"pattern":1}}`)
			//	f.SetCellStyle(sheet, fmt.Sprintf(`%v%d`, "A", i+2), fmt.Sprintf(`%v%d`, "I", i+2), styleSum)
			//}
		}
		// Set active sheet of the workbook.
		f.SetActiveSheet(index)
	}

	// Save xlsx file by the given path.

	path := fmt.Sprintf(`uploads/%v_anal_competitive_groups.xlsx`, result.User.CurrentOrganization.Id)
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
