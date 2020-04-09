package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/digest"
)

//type DocResponseGeneral struct {
//	Id                  		uint      			`json:"id"`   // Идентификатор
//	NameDocumentType		    string    			`json:"name_document_type"`
//	IdDocumentType            	uint                `json:"id_document_type"`
//	IdIdentDocument    			uint                `json:"id_ident_document"`
//	IdEntrant            		uint                `json:"id_entrant"`
//	Created             		time.Time 			`json:"created"`    // Дата создания
//	Checked						bool				`json:"checked"`
//}

func (result *ResultInfo) GetInfoEDocs(ID uint, tableName string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	switch tableName {
	case `compatriot`:
		var r digest.Compatriot
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.CompatriotCategory, `IdCompatriotCategory`)
			result.Items = map[string]interface{}{
				`id`:                       r.Id,
				`id_ident_document`:        r.IdIdentDocument,
				`id_document_type`:         r.DocumentType.Id,
				`name_document_type`:       r.DocumentType.Name,
				`doc_name`:                 r.DocName,
				`doc_org`:                  r.DocOrg,
				`id_compatriot_category`:   r.CompatriotCategory.Id,
				`name_compatriot_category`: r.CompatriotCategory.Name,
				`path_files`:               r.PathFiles,
				`created`:                  r.Created,
			}
		}
		break
	case `composition`:
		var r digest.Composition
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.CompositionThemes, `IdCompositionThemes`)
			db = conn.Model(&r).Related(&r.AppealStatuses, `IdAppealStatus`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                     r.Id,
				`id_ident_document`:      r.IdIdentDocument,
				`id_document_type`:       r.DocumentType.Id,
				`name_document_type`:     r.DocumentType.Name,
				`doc_name`:               r.DocName,
				`doc_org`:                r.DocOrg,
				`id_composition_theme`:   r.CompositionThemes.Id,
				`name_composition_theme`: r.CompositionThemes.Name,
				`id_appeal_status`:       r.AppealStatuses.Id,
				`name_appeal_status`:     r.AppealStatuses.Name,
				`path_files`:             r.PathFiles,
				`has_appealed`:           r.HasAppealed,
				`created`:                r.Created,
				`issue_date`:             issueDate,
				`result`:                 r.Result,
			}
		}
		break
	case `ege`:
		var r digest.Ege
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.Region, `IdRegion`)
			db = conn.Model(&r).Related(&r.Subject, `IdSubject`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			resultDate := r.ResultDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`register_number`:    r.RegisterNumber,
				`doc_number`:         r.DocNumber,
				`mark`:               r.Mark,
				`issue_date`:         issueDate,
				`result_date`:        resultDate,
				`id_region`:          r.Region.Id,
				`name_region`:        r.Region.Name,
				`id_subject`:         r.Subject.Id,
				`name_subject`:       r.Subject.Name,
				`checked`:            r.Checked,
				`created`:            r.Created,
			}
		}
		break
	case `educations`:
		var r digest.Educations
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.Direction, `IdDirection`)
			db = conn.Model(&r).Related(&r.EducationLevel, `IdEducationLevel`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                   r.Id,
				`id_ident_document`:    r.IdIdentDocument,
				`id_document_type`:     r.DocumentType.Id,
				`name_document_type`:   r.DocumentType.Name,
				`doc_name`:             r.DocName,
				`doc_org`:              r.DocOrg,
				`register_number`:      r.RegisterNumber,
				`doc_number`:           r.DocNumber,
				`doc_series`:           r.DocSeries,
				`issue_date`:           issueDate,
				`id_direction`:         r.Direction.Id,
				`name_direction`:       r.Direction.Name,
				`id_education_level`:   r.EducationLevel.Id,
				`name_education_level`: r.EducationLevel.Name,
				`checked`:              r.Checked,
				`created`:              r.Created,
			}
		}
		break
	case `disability`:
		var r digest.Disability
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.DisabilityType, `IdDisabilityType`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                   r.Id,
				`id_ident_document`:    r.IdIdentDocument,
				`id_document_type`:     r.DocumentType.Id,
				`name_document_type`:   r.DocumentType.Name,
				`id_disability_type`:   r.DisabilityType.Id,
				`name_disability_type`: r.DisabilityType.Name,
				`doc_name`:             r.DocName,
				`doc_org`:              r.DocOrg,
				`doc_number`:           r.DocNumber,
				`issue_date`:           issueDate,
				`checked`:              r.Checked,
				`path_file`:            r.PathFiles,
				`created`:              r.Created,
			}
		}
		break
	case `militaries`:
		var r digest.Militaries
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.MilitaryCategories, `IdCategory`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`id_category`:        r.MilitaryCategories.Id,
				`name_category`:      r.MilitaryCategories.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`doc_number`:         r.DocNumber,
				`issue_date`:         issueDate,
				`checked`:            r.Checked,
				`path_file`:          r.PathFiles,
				`created`:            r.Created,
			}
		}
		break
	case `olympics`:
		var r digest.OlympicsDocs
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.Olympics, `IdOlympic`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`id_olympic`:         r.Olympics.Id,
				`name_olympic`:       r.Olympics.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`doc_number`:         r.DocNumber,
				`issue_date`:         issueDate,
				`checked`:            r.Checked,
				`path_file`:          r.PathFiles,
				`created`:            r.Created,
			}
		}
		break
	case `orphans`:
		var r digest.Orphans
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.OrphanCategories, `IdCategory`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`id_category`:        r.OrphanCategories.Id,
				`name_category`:      r.OrphanCategories.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`doc_number`:         r.DocNumber,
				`issue_date`:         issueDate,
				`checked`:            r.Checked,
				`path_file`:          r.PathFiles,
				`created`:            r.Created,
			}
		}
		break
	case `others`:
		var r digest.Other
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`doc_number`:         r.DocNumber,
				`issue_date`:         issueDate,
				`checked`:            r.Checked,
				`path_file`:          r.PathFiles,
				`created`:            r.Created,
			}
		}
		break
	case `parents_lost`:
		var r digest.ParentsLost
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.ParentsLostCategory, `IdCategory`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`doc_number`:         r.DocNumber,
				`issue_date`:         issueDate,
				`checked`:            r.Checked,
				`id_category`:        r.ParentsLostCategory.Id,
				`name_category`:      r.ParentsLostCategory.Name,
				`path_file`:          r.PathFiles,
				`created`:            r.Created,
			}
		}
		break
	case `radiation_work`:
		var r digest.RadiationWork
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.RadiationWorkCategory, `IdCategory`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`doc_number`:         r.DocNumber,
				`issue_date`:         issueDate,
				`checked`:            r.Checked,
				`id_category`:        r.RadiationWorkCategory.Id,
				`name_category`:      r.RadiationWorkCategory.Name,
				`path_file`:          r.PathFiles,
				`created`:            r.Created,
			}
		}
		break
	case `veteran`:
		var r digest.Veteran
		db = conn.Find(&r, ID)
		if r.Id > 0 {
			db = conn.Model(&r).Related(&r.DocumentType, `IdDocumentType`)
			db = conn.Model(&r).Related(&r.VeteranCategory, `IdCategory`)
			issueDate := r.IssueDate.Format(`2006-01-02`)
			result.Items = map[string]interface{}{
				`id`:                 r.Id,
				`id_ident_document`:  r.IdIdentDocument,
				`id_document_type`:   r.DocumentType.Id,
				`name_document_type`: r.DocumentType.Name,
				`doc_name`:           r.DocName,
				`doc_org`:            r.DocOrg,
				`doc_number`:         r.DocNumber,
				`issue_date`:         issueDate,
				`checked`:            r.Checked,
				`id_category`:        r.VeteranCategory.Id,
				`name_category`:      r.VeteranCategory.Name,
				`path_file`:          r.PathFiles,
				`created`:            r.Created,
			}
		}
		break
	default:
		message := `Неизвестный справочник.`
		result.Message = &message
		return
	}
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Документ не найден.`
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	result.Done = true
	return
}
