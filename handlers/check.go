package handlers

import (
	"errors"
	"fmt"
	"persons/config"
	"persons/digest"
)

func CheckCampaignByUser(idCampaign uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`cmp.campaigns`).Select(`id`).Where(`id_organization=? AND id=? AND actual is true`, user.CurrentOrganization.Id, idCampaign).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данной компании `)
	}
}
func CheckEditCampaign(idCampaign uint) error {
	var campaign digest.Campaign
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	conn.Where(`id=? AND actual is true`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}
	var endDate []digest.VEndApplication
	// дата подачи заявления
	conn.Where(`id=? AND id_app_accept_phase IS NULL AND end_date IS NOT NULL `, campaign.Id).Find(&endDate)
	if len(endDate) > 0 {
		return errors.New(`Редактирование приемных компаний невозможно, когда проставлена дата окончания приема заявлений `)
	}
	var existAdmission []digest.AdmissionVolume
	conn.Where(`id_campaign=? AND actual is true`, campaign.Id).Find(&existAdmission)
	if len(existAdmission) > 0 {
		return errors.New(`Редактирование приемных компаний невозможно, когда добавлены КЦП `)
	}
	var existCompetitive []digest.CompetitiveGroup
	conn.Where(`id_campaign=? AND actual is true`, campaign.Id).Find(&existCompetitive)
	if len(existCompetitive) > 0 {
		return errors.New(`Редактирование приемных компаний невозможно, когда добавлены конкусрные группы `)
	}
	return nil
}

func CheckApplicationByUser(idApplication uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`app.applications`).Select(`id`).Where(`id_organization=? AND id=? AND actual is true`, user.CurrentOrganization.Id, idApplication).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данному заявлению `)
	}
}

func CheckAdmissionVolumeByUser(idAdmissionVolume uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`cmp.admission_volume`).Select(`id`).Where(`id_organization=? AND id=? AND actual is true`, user.CurrentOrganization.Id, idAdmissionVolume).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данным кцп `)
	}
}
func CheckEditAdmission(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // Статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}
	return nil
}
func CheckRemoveAdmission(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // Статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}

	return nil
}
func CheckAddAdmission(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // Статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}
	return nil
}

func CheckEditAchievements(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true `, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}
	return nil
}
func CheckAddAchievements(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // 3 - Статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}
	return nil
}

func CheckEditEndDate(idCampaign uint) error {
	var campaign digest.Campaign
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	conn.Where(`id=? AND actual is true`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // 3 - Статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}
	return nil
}

func CheckEditCompetitiveGroup(idCompetitive uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup
	conn.Where(`id=? AND actual is true`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`Конкурсная группа не найдена `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, competitive.IdCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Компания конкурсной группы не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // 3 - статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}

	var applications []digest.Application
	conn.Where(`id_competitive_group=? AND actual IS TRUE`, idCompetitive).Find(&applications)
	if len(applications) > 0 {
		return errors.New(`Найдены заявления с данной конкурсной группой. Редактирование невозможно `)
	}
	return nil
}
func CheckCompetitiveGroupByUser(idCompetitiveGroup uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`cmp.competitive_groups`).Select(`id`).Where(`id_organization=? AND id=? AND actual is true`, user.CurrentOrganization.Id, idCompetitiveGroup).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данной конкурсной группе `)
	}
}
func CheckEditProgramsCompetitiveGroup(idCompetitive uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup
	conn.Where(`id=? AND actual is true`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`Конкурсная группа не найдена `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, competitive.IdCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Компания конкурсной группы не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // 3 - статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}
	return nil
}
func CheckEditEntranceCompetitiveGroup(idCompetitive uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup
	conn.Where(`id=? AND actual is true`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`Конкурсная группа не найдена `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, competitive.IdCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Компания конкурсной группы не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}
	return nil
}
func CheckAddCompetitive(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}
	return nil
}
func CheckRemoveCompetitive(idCompetitive uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	conn.Where(`id=? AND actual is true`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`Конкурсная группа не найдена `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual is true`, competitive.IdCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Компания конкурсной группы не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}

	var applications []digest.Application
	conn.Where(`id_competitive_group=? AND actual IS TRUE`, idCompetitive).Find(&applications)
	if len(applications) > 0 {
		return errors.New(`Найдены заявления с данной конкурсной группой. Редактирование невозможно `)
	}
	return nil
}
func CheckNumberCompetitive(competitive CompetitiveGroup, number int64) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var admission digest.AdmissionVolume
	conn.Table(`cmp.admission_volume`).Where(`id_campaign=? AND id_direction=? AND actual is true`, competitive.IdCampaign, competitive.IdDirection).Find(&admission)
	if admission.Id <= 0 {
		if number == 0 {
			return nil
		} else {
			return errors.New(`Пока не создано КЦП максимальное количество мест - 0 `)
		}
	}
	var distributed digest.DistributedAdmissionVolume
	conn.Where(`id_admission_volume=? AND id_level_budget=? AND actual is true`, admission.Id, competitive.IdLevelBudget).Find(&distributed)
	if distributed.Id <= 0 && competitive.IdEducationSource != 3 { // 3 - платка
		if number == 0 {
			return nil
		} else {
			return errors.New(`Пока не создано КЦП максимальное количество мест - 0 `)
		}
	}
	numField := ``
	var admNumber int64
	switch competitive.IdEducationSource {
	case 1: // Бюджетные места
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `budget_o`
			admNumber = distributed.BudgetO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `budget_oz`
			admNumber = distributed.BudgetOz
			break
		case 3: // Заочная
			numField = `budget_z`
			admNumber = distributed.BudgetZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	case 2: // Квота приема лиц, имеющих особое право
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `quota_o`
			admNumber = distributed.QuotaO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `quota_oz`
			admNumber = distributed.QuotaOz
			break
		case 3: // Заочная
			numField = `quota_z`
			admNumber = distributed.QuotaZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	case 3: // С оплатой обучения
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `paid_o`
			admNumber = admission.PaidO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `paid_oz`
			admNumber = admission.PaidOz
			break
		case 3: // Заочная
			numField = `paid_z`
			admNumber = admission.PaidZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	case 4: // Целевой прием
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `target_o`
			admNumber = distributed.TargetO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `target_oz`
			admNumber = distributed.TargetOz
			break
		case 3: // Заочная
			numField = `target_z`
			admNumber = distributed.TargetZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	default:
		return errors.New(`Ошибка `)
	}
	var sumNumber struct {
		Sum int64
	}
	db := conn.Raw(`SELECT sum(`+numField+`) 
					FROM cmp.competitive_groups 
					WHERE id!=? AND actual is true
					AND id_education_form=? AND id_education_source=? AND id_level_budget=?
					AND id_campaign=? AND id_direction=?`, competitive.Id, competitive.IdEducationForm, competitive.IdEducationSource, competitive.IdLevelBudget, competitive.IdCampaign, competitive.IdDirection).Scan(&sumNumber)
	if db.Error != nil {
		return db.Error
	}
	fmt.Println(`***********************************`)
	fmt.Println(numField)
	fmt.Println(`admNumber`, admNumber)
	fmt.Println(`sumNumber`, sumNumber.Sum)
	fmt.Println(`number`, number)

	if sumNumber.Sum+number > admNumber {
		m := ``
		if competitive.IdEducationSource != 3 {
			m = fmt.Sprintf(`Максимальное количество мест в данном уровне бюджета: %v. Из них распределено по конкурсам: %v. Осталось мест: %v`, admNumber, sumNumber.Sum, admNumber-sumNumber.Sum)
		} else {
			m = fmt.Sprintf(`Максимальное количество мест КЦП: %v. Из них распределено по конкурсам: %v. Осталось мест: %v`, admNumber, sumNumber.Sum, admNumber-sumNumber.Sum)
		}

		fmt.Println(m)
		return errors.New(m)
	}
	return nil
}
func CheckNumberCompetitiveById(idCompetitive uint, number int64) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var competitive digest.CompetitiveGroup
	conn.Where(`id=? AND actual is true`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`КГ не найдены`)
	}
	//if competitive.IdEducationSource == 3 { // 3 - платка
	//
	//	return nil
	//}
	var admission digest.AdmissionVolume
	conn.Table(`cmp.admission_volume`).Where(`id_campaign=? AND id_direction=? AND actual is true`, competitive.IdCampaign, competitive.IdDirection).Find(&admission)
	if admission.Id <= 0 {
		if number == 0 {
			return nil
		} else {
			return errors.New(`Пока не создано КЦП максимальное количество мест - 0 `)
		}
	}
	var distributed digest.DistributedAdmissionVolume
	conn.Where(`id_admission_volume=? AND id_level_budget=?`, admission.Id, competitive.IdLevelBudget).Find(&distributed)
	if distributed.Id <= 0 && competitive.IdEducationSource != 3 { // 3 - платка
		if number == 0 {
			return nil
		} else {
			return errors.New(`Пока не создано КЦП максимальное количество мест - 0 `)
		}
	}
	numField := ``
	var admNumber int64
	switch competitive.IdEducationSource {
	case 1: // Бюджетные места
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `budget_o`
			admNumber = distributed.BudgetO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `budget_oz`
			admNumber = distributed.BudgetOz
			break
		case 3: // Заочная
			numField = `budget_z`
			admNumber = distributed.BudgetZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	case 2: // Квота приема лиц, имеющих особое право
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `quota_o`
			admNumber = distributed.QuotaO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `quota_oz`
			admNumber = distributed.QuotaOz
			break
		case 3: // Заочная
			numField = `quota_z`
			admNumber = distributed.QuotaZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	case 3: // С оплатой обучения
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `paid_o`
			admNumber = admission.PaidO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `paid_oz`
			admNumber = admission.PaidOz
			break
		case 3: // Заочная
			numField = `paid_z`
			admNumber = admission.PaidZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	case 4: // Целевой прием
		switch competitive.IdEducationForm {
		case 1: // Очная форма
			numField = `target_o`
			admNumber = distributed.TargetO
			break
		case 2: // Очно-заочная(вечерняя)
			numField = `target_oz`
			admNumber = distributed.TargetOz
			break
		case 3: // Заочная
			numField = `target_z`
			admNumber = distributed.TargetZ
			break
		default:
			return errors.New(`Ошибка `)
		}
		break
	default:
		return errors.New(`Ошибка `)
	}
	var sumNumber struct {
		Sum int64
	}
	db := conn.Raw(`SELECT sum(`+numField+`) 
					FROM cmp.competitive_groups 
					WHERE id!=? AND actual is true
					AND id_education_form=? AND id_education_source=? AND id_level_budget=?
					AND id_campaign=? AND id_direction=?`, competitive.Id, competitive.IdEducationForm, competitive.IdEducationSource, competitive.IdLevelBudget, competitive.IdCampaign, competitive.IdDirection).Scan(&sumNumber)
	if db.Error != nil {
		return db.Error
	}
	fmt.Println(numField)
	fmt.Println(`admNumber`, admNumber)
	fmt.Println(`sumNumber`, sumNumber.Sum)
	fmt.Println(`number`, number)

	if sumNumber.Sum+number > admNumber {
		m := ``
		if competitive.IdEducationSource != 3 {
			m = fmt.Sprintf(`Максимальное количество мест в данном уровне бюджета: %v. Из них распределено по конкурсам: %v. Осталось мест: %v`, admNumber, sumNumber.Sum, admNumber-sumNumber.Sum)
		} else {
			m = fmt.Sprintf(`Максимальное количество мест КЦП: %v. Из них распределено по конкурсам: %v. Осталось мест: %v`, admNumber, sumNumber.Sum, admNumber-sumNumber.Sum)
		}

		fmt.Println(m)
		return errors.New(m)
	}
	return nil
}
