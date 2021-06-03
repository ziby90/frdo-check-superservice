package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"persons-daemon/config"
	"persons-daemon/model"
	"strings"
	"time"
)

func RatingCompetitiveApplicationStart(countWorker int, timeSleep time.Duration) {
	m := `Найден необработанный пакет со списками рейтинга или как там`
	packageRatingCompetitiveApplication := Package{
		Task:                make(chan uint),
		CountWorker:         countWorker,
		Job:                 RatingCompetitiveApplicationJob,
		ConditionCreateTask: ConditionRatingCompetitiveApplication,
		MessageWorker:       &m,
		TimeSleep:           timeSleep,
	}
	packageRatingCompetitiveApplication.Start()
}

func ConditionRatingCompetitiveApplication() *uint {
	var pack model.RatingCompetitiveApplicationPackages
	Conn.Where(`id_status=1`).Find(&pack).Limit(1)
	if pack.Id > 0 {
		Conn.Model(&pack).Where(`id=?`, pack.Id).Updates(map[string]interface{}{"error": nil, "count_all": 0, "count_add": 0, "id_status": 2, "duration": 0})
		return &pack.Id
	}
	return nil
}

func RatingCompetitiveApplicationJob(id uint) {
	var p model.RatingCompetitiveApplicationPackages
	Conn.Find(&p, id)
	Conn.Exec(`DELETE FROM rating.completitive_groups_applications WHERE id_competitive_group in (Select id_competitive_group FROM packages.rating_competitive_applications_element WHERE id_package=?)`, id)
	Conn.Exec(`DELETE FROM packages.rating_competitive_applications_element WHERE id_package=?`, id)
	pathFile := RestorePath(p.PathFile, p.IdAuthor, p.TableName(), p.CreatedAt)
	ext := strings.Split(p.PathFile, ".")
	p.NewLog(pathFile)
	file, err := os.Open(pathFile)
	fmt.Println(pathFile)
	defer file.Close()
	start := time.Now()
	p.SetLog(fmt.Sprintf(`Начинаю обработку пакета - %d`, id))
	if err != nil {
		fmt.Println(err)
		m := `Ошибка открытия файла`
		p.SetLog(m)
		Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"count_all": 0, "count_add": 0, "error": m, "id_status": 4})
		return
	}
	switch strings.ToUpper(ext[len(ext)-1]) {
	//case `CSV`:
	//	RatingApplicationsParseCsvFile(path, p)
	//	break
	case `XML`:
		RatingApplicationsParseXmlFile(file, p)
		break
	default:
		m := `Неверный формат файла`
		p.SetLog(m)
		Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"count_all": 0, "count_add": 0, "error": m, "id_status": 4})
		return
	}
	p.SetLog(fmt.Sprintf(`Завершил обработку пакета - %d`, id))
	duration := time.Since(start)
	p.SetLog(fmt.Sprintf(`Время обработки пакета составило - %f`, duration.Seconds()))
	Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"duration": duration.Seconds()})
}

func RatingApplicationsParseXmlFile(file *os.File, p model.RatingCompetitiveApplicationPackages) {
	id := p.Id
	countAll := 0
	countAdd := 0
	byteValue, _ := ioutil.ReadAll(file)
	var ratings model.RatingCompetitiveApplicationXml
	err := xml.Unmarshal(byteValue, &ratings)
	if err != nil {
		Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"count_all": 0, "count_add": 0, "error": err.Error(), "id_status": 4})
		return
	}
	var idCompetitiveGroup *uint
	var nameCompetitiveGroup *string
	var publicElements []model.RatingCompetitiveApplicationRow
	var elements []model.RatingCompetitiveApplicationElement
	//Conn.Exec(`DELETE FROM rating.applications WHERE id_competitive_group=?`, idCompetitiveGroup)
	idCompetitiveGroup, nameCompetitiveGroup = FindCompetitiveGroup(ratings.CompetitiveGroupApplicationsList.UIDCompetitiveGroup)
	nameOrganization := FindOrganization(p.IdOrganization)
	idOrganization := p.IdOrganization

	if idCompetitiveGroup != nil {
		errLink := CheckOrganizationCompetitive(*idCompetitiveGroup, p.IdOrganization)
		if errLink != nil {
			Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"count_all": 0, "count_add": 0, "error": errLink.Error(), "id_status": 4})
			return
		}
	}

	ratingRequest := model.RatingCompetitiveRequest{
		IdCompetitiveGroup: idCompetitiveGroup,
		CompetitiveGroup:   nameCompetitiveGroup,
		IdOrganization:     idOrganization,
		Organization:       nameOrganization,
		AdmissionVolume:    ratings.CompetitiveGroupApplicationsList.AdmissionVolume,
		CountFirstStep:     ratings.CompetitiveGroupApplicationsList.CountFirstStep,
		CountSecondStep:    ratings.CompetitiveGroupApplicationsList.CountSecondStep,
		UpdatedAt:          ratings.CompetitiveGroupApplicationsList.UpdatedAt,
	}

	for _, value := range ratings.CompetitiveGroupApplicationsList.Applications.Application {
		countAll++
		var orderId interface{}
		var idApplication *uint
		if value.IDApplicationChoice.UIDEpgu != nil {
			idApplication = FindApplicationsEpgu(*value.IDApplicationChoice.UIDEpgu)
			orderId = value.IDApplicationChoice.UIDEpgu
		}
		if value.IDApplicationChoice.UID != nil {
			idApplication = FindApplications(*value.IDApplicationChoice.UID)
			//orderId = value.IDApplicationChoice.UID
		}

		element := model.RatingCompetitiveApplicationElement{
			RatingCompetitiveRequest: ratingRequest,
			RatingCompetitiveApplication: model.RatingCompetitiveApplication{
				IdApplication:      idApplication,
				Orderid:            orderId,
				Rating:             value.Rating,
				WithoutTests:       value.WithoutTests,
				ReasonWithoutTests: CheckEmptyString(value.ReasonWithoutTests),
				EntranceTest1:      value.EntranceTest1,
				Result1:            value.Result1,
				EntranceTest2:      value.EntranceTest2,
				Result2:            value.Result2,
				EntranceTest3:      value.EntranceTest3,
				Result3:            value.Result3,
				EntranceTest4:      value.EntranceTest4,
				Result4:            value.Result4,
				EntranceTest5:      value.EntranceTest5,
				Result5:            value.Result5,
				Mark:               value.Mark,
				Benefit:            value.Benefit,
				ReasonBenefit:      CheckEmptyString(value.ReasonBenefit),
				SumMark:            value.SumMark,
				Agreed:             value.Agreed,
				Original:           value.Original,
				Addition:           CheckEmptyString(value.Addition),
				Enlisted:           value.Enlisted,
			},
			IdPackage: id,
			Checked:   false,
			Error:     nil,
			CreatedAt: time.Now(),
		}
		if idCompetitiveGroup == nil {
			m := `Не найдена конкурсная группа`
			element.Error = &m
		}

		if element.IdCompetitiveGroup != nil {
			element.Checked = true
			publicElem := model.RatingCompetitiveApplicationRow{
				RatingCompetitiveRequest: element.RatingCompetitiveRequest,
				RatingCompetitiveApplication: model.RatingCompetitiveApplication{
					Orderid:            element.Orderid,
					Rating:             element.Rating,
					WithoutTests:       element.WithoutTests,
					ReasonWithoutTests: element.ReasonWithoutTests,
					EntranceTest1:      element.EntranceTest1,
					Result1:            element.Result1,
					EntranceTest2:      element.EntranceTest2,
					Result2:            element.Result2,
					EntranceTest3:      element.EntranceTest3,
					Result3:            element.Result3,
					EntranceTest4:      element.EntranceTest4,
					Result4:            element.Result4,
					EntranceTest5:      element.EntranceTest5,
					Result5:            element.Result5,
					Mark:               element.Mark,
					Benefit:            element.Benefit,
					ReasonBenefit:      element.ReasonBenefit,
					SumMark:            element.SumMark,
					Agreed:             element.Agreed,
					Original:           element.Original,
					Addition:           element.Addition,
					Enlisted:           element.Enlisted,
					IdApplication:      element.RatingCompetitiveApplication.IdApplication,
				},
			}
			if idApplication != nil {
				errLink := CheckOrganizationApplication(*idApplication, p.IdOrganization)
				if errLink != nil {
					m := errLink.Error()
					element.Error = &m
					element.Checked = false
				} else {
					publicElements = append(publicElements, publicElem)
					element.Checked = true
					Conn.Exec(`UPDATE app.applications SET rating=?, updated_at=? WHERE id=?`, element.Rating, time.Now(), element.IdApplication)
				}
			}

		}

		if idCompetitiveGroup == nil {
			m := `Конкурсная группа не найдена`
			element.Error = &m
		}
		if idOrganization <= 0 {
			m := `Организация не найдена`
			element.Error = &m
		}
		elements = append(elements, element)
		//Conn.Create(&element)
	}
	if len(publicElements) > 0 {
		err = BulkCreatePublicElementsApplications(publicElements)
		if err == nil {
			countAdd = len(publicElements)
		}
	}
	if len(elements) > 0 {
		err = BulkCreateElementsApplications(elements)
		if err == nil {
			countAll = len(elements)
		}
	}
	//send_to_epgu_xml.PrepareSendRatingCompetitiveFroup(config.Conf.EnvType, idCompetitiveGroup)
	if countAll == 0 {
		Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"error": `Данные не обнаружены`, "count_all": countAll, "count_add": countAdd, "id_status": 4})
	}
	Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"error": nil, "count_all": countAll, "count_add": countAdd, "id_status": 3})

	if idCompetitiveGroup!=nil && len(publicElements)>0{
		queryRating := model.Rating{
			IdPackage:          p.Id,
			IdCompetitiveGroup: *idCompetitiveGroup,
			Type:               "ratingcompetitive",
			IdAuthor:           p.IdAuthor,
			Source:             "web",
		}
		err = SendRatingToRabbit(queryRating)
		if err != nil {
			Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"error": err.Error(), "count_all": countAll, "count_add": countAdd, "id_status": 5})
		}
	}


}

func SendRatingToRabbit(rating model.Rating) error {
	code := config.Conf.Code
	if code == nil {
		return errors.New(`not found code in config`)
	}
	query := model.AddQueueRating{
		Code:   *code,
		Rating: rating,
	}
	msg, err := json.Marshal(query)
	if err != nil {
		return err
	}
	var urlPath string
	if config.Conf.UrlRabbit != nil {
		urlPath = *config.Conf.UrlRabbit
	}
	resp, err := http.Post(urlPath+"/api/rating/new", "application/json", bytes.NewBuffer(msg))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(`nbvnvb`)
	}
	return nil
}

func RatingCompetitiveApplicationParseCsvFile(path string, p model.RatingCompetitiveApplicationPackages) {
	//file, err := os.Open(path)
	//if err != nil {
	//	m := `Ошибка открытия файла`
	//	Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"count_all": 0, "count_add": 0, "error": m, "id_status": 4})
	//	return
	//}
	//defer file.Close()
	//countAll := 0
	//countAdd := 0
	//reader := csv.NewReader(file)
	//reader.Comma = '%'
	//
	//
	//send_to_epgu_xml.PrepareSendRatingCompetitiveFroup(config.Conf.EnvType, idCompetitiveGroup)
	//Conn.Model(&p).Where(`id=?`, id).Updates(map[string]interface{}{"error": nil, "count_all": countAll, "count_add": countAdd, "id_status": 3})
}
func FindCompetitiveGroup(uid string) (id *uint, name *string) {
	var item struct {
		Name *string `json:"name"`
		Id   *uint   `json:"id"`
	}
	db := Conn.Where(`uid=? AND actual is true`, uid).Table(`cmp.competitive_groups`).Limit(1).Scan(&item)
	if db.Error != nil && db.Error.Error() != `record not found` {
		//panic(db.Error.Error())
		return nil, nil
	}
	return item.Id, item.Name
}
func CheckOrganizationCompetitive(idCompetitive uint, idOrganization uint) error {
	var item struct {
		Id *uint `json:"id"`
	}
	db := Conn.Where(`id=? AND actual is true AND id_organization=?`, idCompetitive, idOrganization).Table(`cmp.competitive_groups`).Limit(1).Scan(&item)
	if db.Error != nil && db.Error.Error() != `record not found` {
		//panic(db.Error.Error())
		return db.Error
	}
	if item.Id == nil {
		return errors.New(`Не найдена связь организации и конкурса `)
	}
	return nil
}
func CheckOrganizationApplication(idApplication uint, idOrganization uint) error {
	var item struct {
		Id *uint `json:"id"`
	}
	db := Conn.Where(`id=? AND actual is true AND id_organization=?`, idApplication, idOrganization).Table(`app.applications`).Limit(1).Scan(&item)
	if db.Error != nil && db.Error.Error() != `record not found` {
		//panic(db.Error.Error())
		return db.Error
	}
	if item.Id == nil {
		return errors.New(`Не найдена связь организации и заявления `)
	}
	return nil
}

func FindOrganization(id uint) string {
	var item struct {
		ShortTitle string `json:"name"`
	}
	db := Conn.Where(`id=? AND actual is true`, id).Table(`admin.organizations`).Scan(&item)
	if db.Error != nil && db.Error.Error() != `record not found` {
		//panic(db.Error.Error())
		return ``
	}
	return item.ShortTitle
}
func FindApplicationsEpgu(uidEpgu int64) *uint {
	var item struct {
		Id *uint `json:"id"`
	}
	db := Conn.Where(`uid_epgu=? AND actual is true`, uidEpgu).Table(`app.applications`).Scan(&item)
	if db.Error != nil && db.Error.Error() != `record not found` {
		//panic(db.Error.Error())
		return nil
	}
	return item.Id
}
func FindApplications(uid string) *uint {
	var item struct {
		Id *uint `json:"id"`
	}
	db := Conn.Where(`uid=? AND actual is true`, uid).Table(`app.applications`).Scan(&item)
	if db.Error != nil && db.Error.Error() != `record not found` {
		//panic(db.Error.Error())
		return nil
	}
	return item.Id
}
func CheckEmptyString(s *string) *string {
	var newString string
	if s == nil {
		return nil
	}
	newString = strings.TrimSpace(*s)
	if newString == `` {
		return nil
	} else {
		return &newString
	}
}

func BulkCreatePublicElementsApplications(rs []model.RatingCompetitiveApplicationRow) error {
	tx := Conn.Begin()
	countThousand := len(rs) / 1000
	for i := 0; i <= countThousand; i++ {
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for j := i * 1000; j < len(rs) && j < (i*1000+1000); j++ {
			f := rs[j]
			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?"+
				",?,?,?,?,?,?,?,?,?,?"+
				",?,?,?,?,?,?,?,?,?,?,?)")
			valueArgs = append(valueArgs, f.IdCompetitiveGroup)
			valueArgs = append(valueArgs, f.CompetitiveGroup)
			valueArgs = append(valueArgs, f.IdOrganization)
			valueArgs = append(valueArgs, f.Organization)
			valueArgs = append(valueArgs, f.AdmissionVolume)
			valueArgs = append(valueArgs, f.CountFirstStep)
			valueArgs = append(valueArgs, f.CountSecondStep)
			valueArgs = append(valueArgs, f.UpdatedAt)
			valueArgs = append(valueArgs, f.IdApplication)
			valueArgs = append(valueArgs, f.Orderid)
			valueArgs = append(valueArgs, f.WithoutTests)
			valueArgs = append(valueArgs, f.ReasonWithoutTests)
			valueArgs = append(valueArgs, f.EntranceTest1)
			valueArgs = append(valueArgs, f.Result1)
			valueArgs = append(valueArgs, f.EntranceTest2)
			valueArgs = append(valueArgs, f.Result2)
			valueArgs = append(valueArgs, f.EntranceTest3)
			valueArgs = append(valueArgs, f.Result3)
			valueArgs = append(valueArgs, f.EntranceTest4)
			valueArgs = append(valueArgs, f.Result4)
			valueArgs = append(valueArgs, f.EntranceTest5)
			valueArgs = append(valueArgs, f.Result5)
			valueArgs = append(valueArgs, f.Mark)
			valueArgs = append(valueArgs, f.Rating)
			valueArgs = append(valueArgs, f.Benefit)
			valueArgs = append(valueArgs, f.SumMark)
			valueArgs = append(valueArgs, f.Agreed)
			valueArgs = append(valueArgs, f.Original)
			valueArgs = append(valueArgs, f.Addition)
			valueArgs = append(valueArgs, f.Enlisted)
		}
		if len(valueArgs) <= 0 {
			continue
		}
		smt := `INSERT INTO rating.completitive_groups_applications(id_competitive_group, competitive_group, id_organization, organization, admission_volume, count_first_step, count_second_step, changed, id_application, orderid
,  without_tests, reason_without_tests, entrance_test1, result1, entrance_test2, result2, entrance_test3, result3, entrance_test4
, result4, entrance_test5, result5, mark,rating, benefit, sum_mark, agreed, original, addition, enlisted) VALUES %s`
		smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
		db := tx.Exec(smt, valueArgs...)
		if db.Error != nil {
			fmt.Println(fmt.Sprintf(`Поломалася на тысяче № %d`, i))
			tx.Rollback()
			return db.Error
		}
	}
	tx.Commit()
	return nil
}
func BulkCreateElementsApplications(rs []model.RatingCompetitiveApplicationElement) error {
	tx := Conn.Begin()
	countThousand := len(rs) / 1000
	for i := 0; i <= countThousand; i++ {
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for j := i * 1000; j < len(rs) && j < (i*1000+1000); j++ {
			f := rs[j]
			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			valueArgs = append(valueArgs, f.IdPackage)
			valueArgs = append(valueArgs, f.Checked)
			valueArgs = append(valueArgs, f.Error)
			valueArgs = append(valueArgs, f.CreatedAt)
			valueArgs = append(valueArgs, f.IdRatingApplication)
			valueArgs = append(valueArgs, f.IdCompetitiveGroup)
			valueArgs = append(valueArgs, f.CompetitiveGroup)
			valueArgs = append(valueArgs, f.IdOrganization)
			valueArgs = append(valueArgs, f.Organization)
			valueArgs = append(valueArgs, f.AdmissionVolume)
			valueArgs = append(valueArgs, f.CountFirstStep)
			valueArgs = append(valueArgs, f.CountSecondStep)
			valueArgs = append(valueArgs, f.UpdatedAt)
			valueArgs = append(valueArgs, f.IdApplication)
			valueArgs = append(valueArgs, f.Orderid)
			valueArgs = append(valueArgs, f.WithoutTests)
			valueArgs = append(valueArgs, f.ReasonWithoutTests)
			valueArgs = append(valueArgs, f.EntranceTest1)
			valueArgs = append(valueArgs, f.Result1)
			valueArgs = append(valueArgs, f.EntranceTest2)
			valueArgs = append(valueArgs, f.Result2)
			valueArgs = append(valueArgs, f.EntranceTest3)
			valueArgs = append(valueArgs, f.Result3)
			valueArgs = append(valueArgs, f.EntranceTest4)
			valueArgs = append(valueArgs, f.Result4)
			valueArgs = append(valueArgs, f.EntranceTest5)
			valueArgs = append(valueArgs, f.Result5)
			valueArgs = append(valueArgs, f.Mark)
			valueArgs = append(valueArgs, f.Rating)
			valueArgs = append(valueArgs, f.Benefit)
			valueArgs = append(valueArgs, f.SumMark)
			valueArgs = append(valueArgs, f.Agreed)
			valueArgs = append(valueArgs, f.Original)
			valueArgs = append(valueArgs, f.Addition)
			valueArgs = append(valueArgs, f.Enlisted)
		}
		if len(valueArgs) <= 0 {
			continue
		}
		smt := `INSERT INTO packages.rating_competitive_applications_element(id_package,checked, error, created_at, id_rating_application, id_competitive_group, competitive_group, id_organization, organization, admission_volume, count_first_step, count_second_step, changed, id_application
, orderid,  without_tests, reason_without_tests, entrance_test1, result1, entrance_test2, result2, entrance_test3, result3, entrance_test4, result4, entrance_test5, result5
, mark,rating, benefit, sum_mark, agreed, original, addition, enlisted) VALUES %s`
		smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

		db := tx.Exec(smt, valueArgs...)
		if db.Error != nil {
			fmt.Println(fmt.Sprintf(`Поломалася на тысяче № %d`, i))
			tx.Rollback()
			return db.Error
		}
	}
	tx.Commit()
	return nil
}
