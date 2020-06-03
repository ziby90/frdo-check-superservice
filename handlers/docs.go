package handlers

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gabriel-vasile/mimetype"
	"github.com/jinzhu/gorm"
	sendToEpguPath "gitlab.com/unkal/sendtoepgu/path_files"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"persons/config"
	"persons/digest"
	"time"
)

//type DocResponseGeneral struct {
//	Id                  		uint      			"json:"id""   // Идентификатор
//	NameDocumentType		    string    			"json:"name_document_type""
//	IdDocumentType            	uint                "json:"id_document_type""
//	IdIdentDocument    			uint                "json:"id_ident_document""
//	IdEntrant            		uint                "json:"id_entrant""
//	Created             		time.Time 			"json:"created""    // Дата создания
//	Checked						bool				"json:"checked""
//}

func getPath(idEntrant uint, category string, t time.Time) string {
	path, _ := sendToEpguPath.GetPath(idEntrant, t, category)
	//path := `./uploads/docs/` + fmt.Sprintf(`%v`, idEntrant) + `/` + category + `/` + t.Format(`02-01-2006`)
	return path
}

func getIdentName(id uint) string {
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var identification digest.Identifications
	db := conn.Preload(`DocumentType`).Where(`id=?`, id).Find(&identification)
	if db.Error == nil && identification.Id > 0 {
		issueDate := identification.IssueDate.Format(`2006-01-02`)
		name := identification.DocumentType.Name + ` ` + identification.DocSeries + ` ` + identification.DocNumber + ` от ` + issueDate
		return name
	} else {
		return ``
	}

}

func (result *ResultInfo) AddCompatriot(idEntrant uint, cmp digest.Compatriot, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `compatriot`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Compatriot
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.IdEntrant = idEntrant
	doc.Id = 0
	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}

	db = tx.Find(&doc.CompatriotCategory, doc.IdCompatriotCategory)
	if db.Error != nil || doc.CompatriotCategory.Id <= 0 {
		result.SetErrorResult(`Не найдена категория`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Compatriot
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddComposition(idEntrant uint, cmp digest.Composition, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `composition`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Composition
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.IdEntrant = idEntrant
	doc.Id = 0
	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}

	db = tx.Find(&doc.CompositionThemes, doc.IdCompositionTheme)
	if db.Error != nil || doc.CompositionThemes.Id <= 0 {
		result.SetErrorResult(`Не найдена тема сочинения`)
		tx.Rollback()
		return
	}

	db = tx.Find(&doc.AppealStatuses, doc.IdAppealStatus)
	if db.Error != nil || doc.AppealStatuses.Id <= 0 {
		result.SetErrorResult(`Не найден статус апе`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Composition
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddDisability(idEntrant uint, cmp digest.Disability, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `disability`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Disability
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.IdEntrant = idEntrant
	doc.Id = 0
	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}

	db = tx.Find(&doc.DisabilityType, doc.IdDisabilityType)
	if db.Error != nil || doc.DisabilityType.Id <= 0 {
		result.SetErrorResult(`Не найден тип инвалидности`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Disability
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}
	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddEge(idEntrant uint, cmp digest.Ege, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `ege`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Ege
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.IdEntrant = idEntrant
	doc.Id = 0
	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}

	db = tx.Find(&doc.Region, doc.IdRegion)
	if db.Error != nil || doc.Region.Id <= 0 {
		result.SetErrorResult(`Не найден регион`)
		tx.Rollback()
		return
	}
	db = tx.Find(&doc.Subject, doc.IdSubject)
	if db.Error != nil || doc.Subject.Id <= 0 {
		result.SetErrorResult(`Не найден субъект`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Ege
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddEducations(idEntrant uint, cmp digest.Educations, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `educations`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Educations
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.IdEntrant = idEntrant
	doc.Id = 0
	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}
	if doc.IdDirection != nil {
		db = tx.Find(&doc.Direction, *doc.IdDirection)
		if db.Error != nil || doc.Direction.Id <= 0 {
			result.SetErrorResult(`Не найдено направление`)
			tx.Rollback()
			return
		}
	}

	db = tx.Find(&doc.EducationLevel, doc.IdEducationLevel)
	if db.Error != nil || doc.EducationLevel.Id <= 0 {
		result.SetErrorResult(`Не найден уровень образования`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Educations
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddIdentifications(idEntrant uint, cmp digest.Identifications, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `identification`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Identifications
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Find(&doc.Okcm, doc.IdOkcm)
	if db.Error != nil || doc.Okcm.Id <= 0 {
		result.SetErrorResult(`Не найдено оксм`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Identifications
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddMilitaries(idEntrant uint, cmp digest.Militaries, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `militaries`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Militaries
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}
	db = tx.Find(&doc.MilitaryCategories, doc.IdCategory)
	if db.Error != nil || doc.MilitaryCategories.Id <= 0 {
		result.SetErrorResult(`Не найдена категория чего то`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Militaries
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddOlympicsDocs(idEntrant uint, cmp digest.OlympicsDocs, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `olympics`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.OlympicsDocs
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}
	db = tx.Find(&doc.Olympics, doc.IdOlympic)
	if db.Error != nil || doc.Olympics.Id <= 0 {
		result.SetErrorResult(`Не найдена категория чего то`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.OlympicsDocs
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddOrphans(idEntrant uint, cmp digest.Orphans, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `orphans`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Orphans
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}
	db = tx.Find(&doc.OrphanCategories, doc.IdCategory)
	if db.Error != nil || doc.OrphanCategories.Id <= 0 {
		result.SetErrorResult(`Не найдена категория чего то`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Orphans
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddOther(idEntrant uint, cmp digest.Other, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `other`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Other
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.Other
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddParentsLost(idEntrant uint, cmp digest.ParentsLost, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `parents_lost`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.ParentsLost
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}
	db = tx.Find(&doc.ParentsLostCategory, doc.IdCategory)
	if db.Error != nil || doc.ParentsLostCategory.Id <= 0 {
		result.SetErrorResult(`Не найдена категория чего то`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.ParentsLost
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddRadiationWork(idEntrant uint, cmp digest.RadiationWork, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `radiation_work`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.RadiationWork
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}
	db = tx.Find(&doc.RadiationWorkCategory, doc.IdCategory)
	if db.Error != nil || doc.RadiationWorkCategory.Id <= 0 {
		result.SetErrorResult(`Не найдена категория чего то`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.RadiationWork
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}
func (result *ResultInfo) AddVeteran(idEntrant uint, cmp digest.Veteran, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var category digest.DocumentSysCategories
	_ = conn.Where(`name_table=?`, `veteran`).Find(&category)
	if !category.Actual {
		result.SetErrorResult(`Ошибка категории`)
		tx.Rollback()
		return
	}
	var entrant digest.Entrants
	db := conn.Find(&entrant, idEntrant)

	if entrant.Id == 0 {
		result.SetErrorResult(`Абитуриент не найден`)
		tx.Rollback()
		return
	}
	var doc digest.Veteran
	path := getPath(idEntrant, doc.TableName(), time.Now())
	doc = cmp
	doc.IdOrganization = result.User.CurrentOrganization.Id
	doc.Created = time.Now()
	doc.EntrantsId = idEntrant
	doc.Id = 0

	db = tx.Where(`id_entrant=?`, idEntrant).Find(&doc.DocumentIdentification, doc.IdIdentDocument)
	if db.Error != nil || doc.DocumentIdentification.Id <= 0 {
		result.SetErrorResult(`Не найден удостоверяющий документ`)
		tx.Rollback()
		return
	}
	db = tx.Find(&doc.VeteranCategory, doc.IdCategory)
	if db.Error != nil || doc.VeteranCategory.Id <= 0 {
		result.SetErrorResult(`Не найдена категория чего то`)
		tx.Rollback()
		return
	}

	db = tx.Where(`id_sys_category=?`, category.Id).Find(&doc.DocumentType, doc.IdDocumentType)
	if db.Error != nil || doc.DocumentType.Id <= 0 {
		result.SetErrorResult(`Не найден тип документа`)
		tx.Rollback()
		return
	}

	if cmp.Uid != nil {
		var exist digest.RadiationWork
		tx.Where(`id_organization=? AND uid=?`, result.User.CurrentOrganization.Id, *doc.Uid).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Документ с данным uid уже существует у выбранной организации`)
			tx.Rollback()
			return
		}
		doc.Uid = cmp.Uid
	}
	if f != nil {
		ext := filepath.Ext(path + `/` + f.Header.Filename)
		sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
		name := hex.EncodeToString(sha1FileName[:]) + ext
		if _, err := os.Stat(path); err != nil {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
		}
		dst, err := os.Create(filepath.Join(path, name))
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, f.MultFile)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		doc.PathFile = &name
	} else {
		doc.PathFile = nil
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&doc)
	if db.Error != nil || doc.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		if f != nil {
			_ = os.Remove(filepath.Join(path, f.Header.Filename))
		}
		tx.Rollback()
		return
	}

	result.Items = map[string]interface{}{
		`id_entrant`:  idEntrant,
		`id_document`: doc.Id,
	}
	result.Done = true
	tx.Commit()
}

func (result *ResultInfo) GetInfoEDocs(ID uint, tableName string) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	var sysCategory digest.DocumentSysCategories
	db = conn.Where("name_table=?", tableName).Find(&sysCategory)
	switch tableName {
	case "compatriot":
		var r digest.Compatriot
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.CompatriotCategory, "IdCompatriotCategory")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                       r.Id,
				"id_ident_document":        r.IdIdentDocument,
				"id_entrant":               r.IdEntrant,
				"id_document_type":         r.DocumentType.Id,
				"name_document_type":       r.DocumentType.Name,
				"doc_name":                 r.DocName,
				"doc_org":                  r.DocOrg,
				"id_compatriot_category":   r.CompatriotCategory.Id,
				"name_compatriot_category": r.CompatriotCategory.Name,
				"created":                  r.Created,
				"name_sys_category":        sysCategory.Name,
				"code_sys_category":        sysCategory.NameTable,
				"uid":                      r.Uid,
				"file":                     file,
				"name_ident_document":      getIdentName(r.IdIdentDocument),
			}
		}
		break
	case "composition":
		var r digest.Composition
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.CompositionThemes, "IdCompositionTheme")
			db = conn.Model(&r).Related(&r.AppealStatuses, "IdAppealStatus")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                     r.Id,
				"id_ident_document":      r.IdIdentDocument,
				"id_entrant":             r.IdEntrant,
				"id_document_type":       r.DocumentType.Id,
				"name_document_type":     r.DocumentType.Name,
				"doc_name":               r.DocName,
				"doc_org":                r.DocOrg,
				"doc_year":               r.DocYear,
				"id_composition_theme":   r.CompositionThemes.Id,
				"name_composition_theme": r.CompositionThemes.Name,
				"id_appeal_status":       r.AppealStatuses.Id,
				"name_appeal_status":     r.AppealStatuses.Name,
				"has_appealed":           r.HasAppealed,
				"created":                r.Created,
				"issue_date":             issueDate,
				"result":                 r.Result,
				"name_sys_category":      sysCategory.Name,
				"uid":                    r.Uid,
				"file":                   file,
				"name_ident_document":    getIdentName(r.IdIdentDocument),
				"code_sys_category":      sysCategory.NameTable,
			}
		}
		break
	case "ege":
		var r digest.Ege
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.Region, "IdRegion")
			db = conn.Model(&r).Related(&r.Subject, "IdSubject")
			issueDate := r.IssueDate.Format("2006-01-02")
			resultDate := r.ResultDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.IdEntrant,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"register_number":     r.RegisterNumber,
				"doc_number":          r.DocNumber,
				"mark":                r.Mark,
				"issue_date":          issueDate,
				"result_date":         resultDate,
				"id_region":           r.Region.Id,
				"name_region":         r.Region.Name,
				"id_subject":          r.Subject.Id,
				"name_subject":        r.Subject.Name,
				"checked":             r.Checked,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	case "educations":
		var r digest.Educations
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.Direction, "IdDirection")
			db = conn.Model(&r).Related(&r.EducationLevel, "IdEducationLevel")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                   r.Id,
				"id_ident_document":    r.IdIdentDocument,
				"id_entrant":           r.IdEntrant,
				"id_document_type":     r.DocumentType.Id,
				"name_document_type":   r.DocumentType.Name,
				"doc_name":             r.DocName,
				"doc_org":              r.DocOrg,
				"register_number":      r.RegisterNumber,
				"doc_number":           r.DocNumber,
				"doc_series":           r.DocSeries,
				"issue_date":           issueDate,
				"id_direction":         r.Direction.Id,
				"name_direction":       r.Direction.Name,
				"id_education_level":   r.EducationLevel.Id,
				"name_education_level": r.EducationLevel.Name,
				"checked":              r.Checked,
				"created":              r.Created,
				"name_sys_category":    sysCategory.Name,
				"uid":                  r.Uid,
				"file":                 file,
				"name_ident_document":  getIdentName(r.IdIdentDocument),
				"code_sys_category":    sysCategory.NameTable,
			}
		}
		break
	case "disability":
		var r digest.Disability
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.DisabilityType, "IdDisabilityType")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                   r.Id,
				"id_ident_document":    r.IdIdentDocument,
				"id_entrant":           r.IdEntrant,
				"id_document_type":     r.DocumentType.Id,
				"name_document_type":   r.DocumentType.Name,
				"id_disability_type":   r.DisabilityType.Id,
				"name_disability_type": r.DisabilityType.Name,
				"doc_name":             r.DocName,
				"doc_org":              r.DocOrg,
				"doc_number":           r.DocNumber,
				"issue_date":           issueDate,
				"checked":              r.Checked,
				"created":              r.Created,
				"name_sys_category":    sysCategory.Name,
				"uid":                  r.Uid,
				"file":                 file,
				"name_ident_document":  getIdentName(r.IdIdentDocument),
				"code_sys_category":    sysCategory.NameTable,
			}
		}
		break
	case "identification":
		var r digest.Identifications
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.Okcm, "IdOkcm")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                 r.Id,
				"id_entrant":         r.EntrantsId,
				"id_document_type":   r.DocumentType.Id,
				"name_document_type": r.DocumentType.Name,
				"surname":            r.Surname,
				"name":               r.Name,
				"patronymic":         r.Patronymic,
				"doc_series":         r.DocSeries,
				"doc_number":         r.DocNumber,
				"doc_organization":   r.DocOrganization,
				"id_okcm":            r.IdOkcm,
				"name_okcm":          r.Okcm.ShortName,
				"issue_date":         issueDate,
				"subdivision_code":   r.SubdivisionCode,
				"checked":            r.Checked,
				"created":            r.Created,
				"name_sys_category":  sysCategory.Name,
				"uid":                r.Uid,
				"file":               file,
				"code_sys_category":  sysCategory.NameTable,
			}
		}
		break
	case "militaries":
		var r digest.Militaries
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.MilitaryCategories, "IdCategory")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.EntrantsId,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"id_category":         r.MilitaryCategories.Id,
				"name_category":       r.MilitaryCategories.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"doc_number":          r.DocNumber,
				"doc_series":          r.DocSeries,
				"issue_date":          issueDate,
				"checked":             r.Checked,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	case "olympics":
		var r digest.OlympicsDocs
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.Olympics, "IdOlympic")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.EntrantsId,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"id_olympic":          r.Olympics.Id,
				"name_olympic":        r.Olympics.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"doc_number":          r.DocNumber,
				"doc_series":          r.DocSeries,
				"issue_date":          issueDate,
				"checked":             r.Checked,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	case "orphans":
		var r digest.Orphans
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.OrphanCategories, "IdCategory")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.EntrantsId,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"id_category":         r.OrphanCategories.Id,
				"name_category":       r.OrphanCategories.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"doc_series":          r.DocSeries,
				"doc_number":          r.DocNumber,
				"issue_date":          issueDate,
				"checked":             r.Checked,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	case "other":
		var r digest.Other
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.EntrantsId,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"doc_number":          r.DocNumber,
				"doc_series":          r.DocSeries,
				"issue_date":          issueDate,
				"checked":             r.Checked,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	case "parents_lost":
		var r digest.ParentsLost
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.ParentsLostCategory, "IdCategory")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.EntrantsId,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"doc_number":          r.DocNumber,
				"doc_series":          r.DocSeries,
				"issue_date":          issueDate,
				"checked":             r.Checked,
				"id_category":         r.ParentsLostCategory.Id,
				"name_category":       r.ParentsLostCategory.Name,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	case "radiation_work":
		var r digest.RadiationWork
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.RadiationWorkCategory, "IdCategory")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.EntrantsId,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"doc_number":          r.DocNumber,
				"doc_series":          r.DocSeries,
				"issue_date":          issueDate,
				"checked":             r.Checked,
				"id_category":         r.RadiationWorkCategory.Id,
				"name_category":       r.RadiationWorkCategory.Name,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	case "veteran":
		var r digest.Veteran
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, "IdDocumentType")
			db = conn.Model(&r).Related(&r.VeteranCategory, "IdCategory")
			issueDate := r.IssueDate.Format("2006-01-02")
			file := false
			if r.PathFile != nil {
				file = true
			}
			result.Items = map[string]interface{}{
				"id":                  r.Id,
				"id_ident_document":   r.IdIdentDocument,
				"id_entrant":          r.EntrantsId,
				"id_document_type":    r.DocumentType.Id,
				"name_document_type":  r.DocumentType.Name,
				"doc_name":            r.DocName,
				"doc_org":             r.DocOrg,
				"doc_number":          r.DocNumber,
				"doc_series":          r.DocSeries,
				"issue_date":          issueDate,
				"checked":             r.Checked,
				"id_category":         r.VeteranCategory.Id,
				"name_category":       r.VeteranCategory.Name,
				"created":             r.Created,
				"name_sys_category":   sysCategory.Name,
				"uid":                 r.Uid,
				"file":                file,
				"name_ident_document": getIdentName(r.IdIdentDocument),
				"code_sys_category":   sysCategory.NameTable,
			}
		}
		break
	default:
		message := "Неизвестный справочник."
		result.Message = &message
		return
	}

	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Документ не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}
	result.Done = true

	return
}

func (result *ResultInfo) GetFileDoc(ID uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.VDocuments
	db := conn.Where(`id_document=?`, ID).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Документ не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}
	if doc.PathFile != nil {
		filename := *doc.PathFile
		path := getPath(doc.EntrantsId, `documents.`+doc.NameTable, doc.Created) + `/` + filename
		f, err := os.Open(path)
		if err != nil {
			result.SetErrorResult(err.Error())
			return
		} else {
			defer f.Close()
			reader := bufio.NewReader(f)
			content, _ := ioutil.ReadAll(reader)
			ext := mimetype.Detect(content)
			file := digest.FileC{
				Content: content,
				Size:    int64(len(content)),
				Title:   filename,
				Type:    ext.Extension(),
			}
			result.Items = file
		}
	} else {
		message := "Файл не найден."
		result.Message = &message
		return
	}
	result.Done = true
	return
}

func (result *ResultInfo) RemoveFileDoc(ID uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.VDocuments
	db := conn.Where(`id_document=?`, ID).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Документ не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}

	if doc.UidEpgu != nil {
		message := "Документы, загруженные через ЕПГУ не подлежат редактированию"
		result.Message = &message
		return
	}
	if (doc.IdOrganization != nil && *doc.IdOrganization != result.User.CurrentOrganization.Id) || doc.IdOrganization == nil {
		message := "Выбранная организация не соответствует организации , создавшей документ."
		result.Message = &message
		return
	}
	if doc.PathFile != nil {
		db = conn.Exec(`UPDATE documents.`+doc.NameTable+` SET path_file=null WHERE id=?`, doc.IdDocument)
		if db.Error != nil {
			result.SetErrorResult(db.Error.Error())
			return
		}
	} else {
		message := "Файл не найден."
		result.Message = &message
		return
	}
	result.Done = true
	return
}

func (result *ResultInfo) AddFileDoc(ID uint, f *digest.File) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.VDocuments
	db := conn.Where(`id_document=?`, ID).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Документ не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}
	if doc.UidEpgu != nil {
		message := "Документы, загруженные через ЕПГУ не подлежат редактированию"
		result.Message = &message
		return
	}
	if (doc.IdOrganization != nil && *doc.IdOrganization != result.User.CurrentOrganization.Id) || doc.IdOrganization == nil {
		message := "Выбранная организация не соответствует организации , создавшей документ."
		result.Message = &message
		return
	}
	path := getPath(doc.EntrantsId, `documents.`+doc.NameTable, time.Now())
	ext := filepath.Ext(path + `/` + f.Header.Filename)
	sha1FileName := sha1.Sum([]byte(doc.TableName() + time.Now().String()))
	name := hex.EncodeToString(sha1FileName[:]) + ext
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			result.SetErrorResult(err.Error())
			return
		}
	}
	dst, err := os.Create(filepath.Join(path, name))
	if err != nil {
		result.SetErrorResult(err.Error())
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, f.MultFile)
	if err != nil {
		result.SetErrorResult(err.Error())
		return
	}

	db = conn.Exec(`UPDATE documents.`+doc.NameTable+` SET path_file=? WHERE id=?`, &name, doc.IdDocument)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		return
	}
	result.Done = true
	return
}
