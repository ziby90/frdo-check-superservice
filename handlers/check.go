package handlers

import (
	"errors"
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
	db := conn.Table(`cmp.campaigns`).Select(`id`).Where(`id_organization=? AND id=?`, user.CurrentOrganization.Id, idCampaign).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данной компании `)
	}
}
func CheckCompetitiveGroupByUser(idCompetitiveGroup uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`cmp.competitive_groups`).Select(`id`).Where(`id_organization=? AND id=?`, user.CurrentOrganization.Id, idCompetitiveGroup).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данной конкурсной группе `)
	}
}
func CheckAdmissionVolumeByUser(idAdmissionVolume uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`cmp.admission_volume`).Select(`id`).Where(`id_organization=? AND id=?`, user.CurrentOrganization.Id, idAdmissionVolume).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данным кцп `)
	}
}
func CheckApplicationByUser(idApplication uint, user digest.User) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var count int
	if user.Role.Code == `administrator` {
		return nil
	}
	db := conn.Table(`app.applications`).Select(`id`).Where(`id_organization=? AND id=?`, user.CurrentOrganization.Id, idApplication).Count(&count)
	if db.Error != nil {
		return db.Error
	}
	if count > 0 {
		return nil
	} else {
		return errors.New(`У пользователя нет доступа к данному заявлению `)
	}
}

func CheckEditAchievements(idAchievement uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var achievement digest.IndividualAchievements
	conn.Where(`id=? AND actual`, idAchievement).Find(&achievement)
	if achievement.Id <= 0 {
		return errors.New(`Достижение не найдено `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual`, achievement.IdCampaign).Find(&campaign)
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
	conn.Where(`id=? AND actual`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // 3 - Статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}
	return nil
}
func CheckEditAdmission(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual`, idCampaign).Find(&campaign)
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
	conn.Where(`id=? AND actual`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // Статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}
	return nil
}
func CheckEditCampaign(idCampaign uint) error {
	var campaign digest.Campaign
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	conn.Where(`id=? AND actual`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}
	var endDate []digest.VEndApplication
	conn.Where(`id=? AND id_app_accept_phase IS NULL AND end_date IS NOT NULL`, campaign.Id).Find(&endDate)
	if len(endDate) > 0 {
		return errors.New(`Редактирование приемных компаний невозможно, когда проставлена дата окончания приема заявлений `)
	}
	var existAdmission []digest.AdmissionVolume
	conn.Where(`id_campaign=?`, campaign.Id).Find(&existAdmission)
	if len(existAdmission) > 0 {
		return errors.New(`Редактирование приемных компаний невозможно, когда добавлены КЦП `)
	}
	var existCompetitive []digest.CompetitiveGroup
	conn.Where(`id_campaign=?`, campaign.Id).Find(&existCompetitive)
	if len(existCompetitive) > 0 {
		return errors.New(`Редактирование приемных компаний невозможно, когда добавлены конкусрные группы `)
	}
	return nil
}
func CheckEditEndDate(idCampaign uint) error {
	var campaign digest.Campaign
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	conn.Where(`id=? AND actual`, idCampaign).Find(&campaign)
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
	conn.Where(`id=? AND actual`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`Конкурсная группа не найдена `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual`, competitive.IdCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Компания конкурсной группы не найдена `)
	}
	if campaign.IdCampaignStatus == 3 { // 3 - статус завершена
		return errors.New(`Редактирование невозможно. Приемная компания завершена. `)
	}

	var applications []digest.Application
	conn.Where(`id_competitive_groups=? `, idCompetitive).Find(&applications)
	if len(applications) > 0 {
		return errors.New(`Найдены заявления с данной конкурсной группой. Редактирование невозможно `)
	}
	return nil
}
func CheckEditProgramsCompetitiveGroup(idCompetitive uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)

	var competitive digest.CompetitiveGroup
	conn.Where(`id=? AND actual`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`Конкурсная группа не найдена `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual`, competitive.IdCampaign).Find(&campaign)
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
	conn.Where(`id=? AND actual`, idCompetitive).Find(&competitive)
	if competitive.Id <= 0 {
		return errors.New(`Конкурсная группа не найдена `)
	}
	var campaign digest.Campaign
	conn.Where(`id=? AND actual`, competitive.IdCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Компания конкурсной группы не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}
	return nil
}
func CheckAddRemoveCompetitive(idCampaign uint) error {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	conn.Where(`id=? AND actual`, idCampaign).Find(&campaign)
	if campaign.Id <= 0 {
		return errors.New(`Приемная компания не найдена `)
	}
	if campaign.IdCampaignStatus != 1 { // 1 - Статус набор не начался
		return errors.New(`Редактирование возможно только в статусе набор не начался. `)
	}
	return nil
}
