package api_error

import (
	"persons/error_handler"
)

var GORMBDErr = error_handler.AlterErrorType{Code: 12, ToUserCode: 6000, Object: "GORM",
	CommentTemplate: `Ошибка в БД "%s"`}

var MergoErr = error_handler.AlterErrorType{Code: 16, ToUserCode: 6000, Object: "mergo",
	CommentTemplate: `Ошибка в mergo "%s"`}

var CopyErr = error_handler.AlterErrorType{Code: 103, ToUserCode: 6000, Object: "copy",
	CommentTemplate: `Ошибка в copy "%s"`}

//var JsonUnmarshalError = error_handler.AlterErrorType{Code: 159, ToUserCode: 5000, Object: "Header",
//	CommentTemplate: `Ошибка при распаковке json части header. В json Header тип поля IDJWT - беззнаковое число, остальные строки`}
////
var HeaderJsonUnmarshalErr = error_handler.AlterErrorType{Code: 159, ToUserCode: 499, Object: "Header",
	CommentTemplate: `Ошибка при распаковки json содержащий Header в структуру. Тип IDJWT - число, остальные поля строки`}

var ValidXmlErr = error_handler.AlterErrorType{
	Code:            1,
	ToUserCode:      403,
	Object:          "XML",
	CommentTemplate: `Ошибка валидации XML в структуру "%s".`,
}

var UnmarshalXmlErr = error_handler.AlterErrorType{
	Code:            1,
	ToUserCode:      403,
	Object:          "XML",
	CommentTemplate: `Ошибка  Unmarshal XML в структуру %s. Ошибка "%s".`,
}
var CampaignExist = error_handler.AlterErrorType{
	Code:            2,
	ToUserCode:      306,
	Object:          "campaign",
	CommentTemplate: `Приемная кампания с данным UID "%d" уже существует.`,
}
var CampaignTypeNotExist = error_handler.AlterErrorType{
	Code:            4,
	ToUserCode:      423,
	Object:          "campaign",
	CommentTemplate: `Приемной камании с ID "%d" не найден.`,
}
var CampaignStatusNotExist = error_handler.AlterErrorType{
	Code:            5,
	ToUserCode:      424,
	Object:          "campaign",
	CommentTemplate: `Тип статуса приемной камании с ID "%d" не найден.`,
}

var CampaignEducationFormsRepeat = error_handler.AlterErrorType{
	Code:            6,
	ToUserCode:      001,
	Object:          "campaign",
	CommentTemplate: `У приемной кампании UID "%s" повторяющиеся значения форм образования.`,
}

var CampaignEducationLevelRepeat = error_handler.AlterErrorType{
	Code:            7,
	ToUserCode:      002,
	Object:          "campaign",
	CommentTemplate: `У приемной кампании UID "%s" повторяющиеся значения уровней образования.`,
}
var EducationFormNotFound = error_handler.AlterErrorType{Code: 8, ToUserCode: 4010, Object: "EducationForm",
	CommentTemplate: `Форма образования с ID "%d" не найдена.`,
}

var EducationLevelNotFound = error_handler.AlterErrorType{Code: 9, ToUserCode: 4011, Object: "EducationLevel",
	CommentTemplate: `Форма образования с ID "%d" не найдена.`,
}
var EducationLevelByCampainTypeNotFound = error_handler.AlterErrorType{Code: 10, ToUserCode: 4008, Object: "EducationLevel",
	CommentTemplate: `Уровень образования "%s" отсутствует в типе приемной кампании "%s".`,
}
var CampaignByTypeExist = error_handler.AlterErrorType{Code: 11, ToUserCode: 3022, Object: "campaign",
	CommentTemplate: `Приемная кампания с типом ID  "%d" уже существует UID "%s".`}

var CampaignNotExist = error_handler.AlterErrorType{Code: 13, ToUserCode: 4018, Object: "campaign",
	CommentTemplate: `Приемной капании с данным UID "%s" не существует.`}

var CampaignEducFormHaveCompetitiveGroup = error_handler.AlterErrorType{Code: 14, ToUserCode: 1003, Object: "campaign",
	CommentTemplate: `Ошибка при удалении формы образвания приемной компании UID "%s" , у формы образования ID "%d" присутствует конкурсная группа UID "%s".`}

var CampaignEducLevelHaveAdmissionVolume = error_handler.AlterErrorType{Code: 15, ToUserCode: 1004, Object: "campaign",
	CommentTemplate: `Ошибка при удалении уровней образвания приемной компании UID "%s" , у уровня образования  ID "%d" присутствует КЦП UID "%s".`}

var EntrantExist = error_handler.AlterErrorType{Code: 16, ToUserCode: 3007, Object: "entrant",
	CommentTemplate: `Абитуриент со СНИЛС  "%s" уже существует.`}

var GenderNotExist = error_handler.AlterErrorType{Code: 17, ToUserCode: 4040, Object: "Gender",
	CommentTemplate: `Пол с ID "%d" не найден.`}

var OKSMNotExist = error_handler.AlterErrorType{Code: 18, ToUserCode: 4041, Object: "OKSM",
	CommentTemplate: `Запись в Okcm с ID "%d" не найдена.`}

var EntrantNoTExist = error_handler.AlterErrorType{Code: 19, ToUserCode: 4042, Object: "entrant",
	CommentTemplate: `Абитуриент со СНИЛС  "%s" не найден.`}

var EntrantRemove = error_handler.AlterErrorType{Code: 20, ToUserCode: 1004, Object: "entrant",
	CommentTemplate: `Удаление абитуриентов запрещено.`}

var DirectionNotExist = error_handler.AlterErrorType{Code: 21, ToUserCode: 4001, Object: "Direction",
	CommentTemplate: `Специальность с ID "%d" не найдена.`}

var AdmissionVolumeExist = error_handler.AlterErrorType{Code: 22, ToUserCode: 3003, Object: "admission_volume",
	CommentTemplate: `КЦП с UID "%s" уже существует.`}

var AdmissionVolumeByEduclevelExist = error_handler.AlterErrorType{Code: 23, ToUserCode: 3004, Object: "admission_volume",
	CommentTemplate: `КЦП с уровнем образования с ID "%d" уже существует.`}

var DirectionBYOrgNotExist = error_handler.AlterErrorType{Code: 24, ToUserCode: 3004, Object: "admission_volume",
	CommentTemplate: `Направление подготовки ID "%d" не найдено у текущей организации.`}

var AdmissionVolumeEditByDirectionNotExist = error_handler.AlterErrorType{Code: 25, ToUserCode: 4013, Object: "admission_volume",
	CommentTemplate: `Отсутствует редактируемые КЦП UID "%s" по направлению подготовки ID "%d".`}

var AdmissionVolumeSpaceLessCompetitiveGroupSpaceSum = error_handler.AlterErrorType{Code: 26, ToUserCode: 2004, Object: "admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах превышает количество мест в КЦП UID "%s" по направлению подготовки ID "%d".`}

var AdmissionVolumeSpaceLessBudgetO = error_handler.AlterErrorType{Code: 27, ToUserCode: 2005, Object: "BudgetO",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет очной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessBudgetOz = error_handler.AlterErrorType{Code: 28, ToUserCode: 2006, Object: "BudgetOz",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет очно-заочной формы обучения обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessBudgetZ = error_handler.AlterErrorType{Code: 29, ToUserCode: 2007, Object: "BudgetZ",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessPaidO = error_handler.AlterErrorType{Code: 30, ToUserCode: 2008, Object: "PaidO",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе очной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessPaidOz = error_handler.AlterErrorType{Code: 31, ToUserCode: 2009, Object: "PaidOz",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе очно-заочной формы обучения обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessPaidZ = error_handler.AlterErrorType{Code: 32, ToUserCode: 2010, Object: "PaidZ",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessTargetO = error_handler.AlterErrorType{Code: 33, ToUserCode: 2011, Object: "TargetO",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления очной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessTargetOz = error_handler.AlterErrorType{Code: 34, ToUserCode: 2012, Object: "TargetOz",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления очной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessTargetZ = error_handler.AlterErrorType{Code: 35, ToUserCode: 2013, Object: "TargetZ",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessQuotaO = error_handler.AlterErrorType{Code: 36, ToUserCode: 2014, Object: "QuotaO",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте очной формы обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessQuotaOz = error_handler.AlterErrorType{Code: 37, ToUserCode: 2015, Object: "QuotaOz",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте очно-заочной формы обучения обучения превышает количество мест в КЦП  UID "%s" "%d".`}

var AdmissionVolumeSpaceLessQuotaZ = error_handler.AlterErrorType{Code: 38, ToUserCode: 2016, Object: "QuotaZ",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeSpaceSum = error_handler.AlterErrorType{Code: 39, ToUserCode: 2017, Object: "admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах превышает количество мест в КЦП UID "%s" по направлению подготовки ID "%d".`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeBudgetO = error_handler.AlterErrorType{Code: 40, ToUserCode: 2018, Object: "BudgetO",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема, по всем уровням бюджета, на бюджет очной формы обучения превышает количество мест в КЦП UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeBudgetOz = error_handler.AlterErrorType{Code: 41, ToUserCode: 2019, Object: "BudgetOz",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема, по всем уровням бюджета, на бюджет очно-заочной формы обучения обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeBudgetZ = error_handler.AlterErrorType{Code: 42, ToUserCode: 2020, Object: "BudgetZ",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета на бюджет заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumePaidO = error_handler.AlterErrorType{Code: 43, ToUserCode: 2021, Object: "PaidO",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета на бюджет заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumePaidOz = error_handler.AlterErrorType{Code: 44, ToUserCode: 2022, Object: "PaidOz",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета платной основе очно-заочной формы обучения обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumePaidZ = error_handler.AlterErrorType{Code: 45, ToUserCode: 2023, Object: "PaidZ",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета платной основе заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeTargetO = error_handler.AlterErrorType{Code: 46, ToUserCode: 2024, Object: "TargetO",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета целевого направления очной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeTargetOz = error_handler.AlterErrorType{Code: 46, ToUserCode: 2025, Object: "TargetOz",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета целевого направления очной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeTargetZ = error_handler.AlterErrorType{Code: 47, ToUserCode: 2026, Object: "TargetZ",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета целевого направления заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeQuotaO = error_handler.AlterErrorType{Code: 48, ToUserCode: 2027, Object: "QuotaO",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета целевого направления заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeQuotaOz = error_handler.AlterErrorType{Code: 49, ToUserCode: 2028, Object: "QuotaOz",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета по квоте очно-заочной формы обучения обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeSpaceLessDistributedAdmissionVolumeQuotaZ = error_handler.AlterErrorType{Code: 50, ToUserCode: 2029, Object: "QuotaZ",
	CommentTemplate: `Сумма по распределенным контрольным цифрам приема по всем уровням бюджета целевого направления заочной формы обучения превышает количество мест в КЦП  UID "%s" "%d"`}

var AdmissionVolumeHaveCompetitiveGroup = error_handler.AlterErrorType{Code: 51, ToUserCode: 1005, Object: "admission_volume",
	CommentTemplate: `У КЦП UID "%s" Существуют конкурсные группы UIDs "%s".`}

var EducationLevelInCampaignNotExist = error_handler.AlterErrorType{Code: 52, ToUserCode: 4044, Object: "CampaignEducLevel",
	CommentTemplate: `Уровень образования "%s" отсутствует в кампании "%s".`}

var CampaignEducationLevelNotExist = error_handler.AlterErrorType{Code: 53, ToUserCode: 4027, Object: "campaign",
	CommentTemplate: `Уровень образования с ID "%d" у кампании UID "%s" не найден.`}

var EducationSourceNotExist = error_handler.AlterErrorType{Code: 54, ToUserCode: 4002, Object: "competitive_group",
	CommentTemplate: `Источник финансирования с ID "%d" не найден.`}

var CompetitiveGroupExist = error_handler.AlterErrorType{Code: 55, ToUserCode: 3013, Object: "competitive_group",
	CommentTemplate: `Конкурсная группа с UID "%s" уже существует.`}

var AdmissionVolumeByDirectionNotExist = error_handler.AlterErrorType{Code: 56, ToUserCode: 4012, Object: "admission_volume",
	CommentTemplate: `КЦП с уровнем образования ID "%d" не найдено.`}

var CompetitiveGroupSpaceSumMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 57, ToUserCode: 4, Object: "competitive_group",
	CommentTemplate: `Сумма мест в конкурсных группах превышает количесво мест заданых в КЦП.`}

var CompetitiveGroupBudgetOMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 58, ToUserCode: 5, Object: "competitive_group_BudgetO",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет очной формы обучения превышает количество мест в КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupBudgetOzMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 59, ToUserCode: 6, Object: "competitive_group_BudgetOz",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет очно-заочной формы обучения обучения превышает значение в  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupBudgetZMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 60, ToUserCode: 7, Object: "competitive_group_BudgetZ",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет заочной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupPaidOMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 61, ToUserCode: 8, Object: "competitive_group_PaidO",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе очной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupPaidOzMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 62, ToUserCode: 9, Object: "competitive_group_PaidOz",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе очной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupPaidZMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 63, ToUserCode: 10, Object: "competitive_group_PaidZ",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе заочной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupTargetOMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 64, ToUserCode: 11, Object: "competitive_group_TargetO",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления очной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupTargetOzMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 65, ToUserCode: 12, Object: "competitive_group_TargetOz",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления очной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupTargetZMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 66, ToUserCode: 13, Object: "competitive_group_TargetZ",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления очной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupQuotaOMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 67, ToUserCode: 14, Object: "competitive_group_QuotaO",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте очной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupQuotaOzMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 68, ToUserCode: 15, Object: "competitive_group_QuotaOz",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте очно-заочной формы обучения обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var CompetitiveGroupQuotaZMoreThenAdmissionVolume = error_handler.AlterErrorType{Code: 69, ToUserCode: 16, Object: "competitive_group_QuotaZ",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте заочной формы обучения превышает значение  КЦП UID "%s" Идентификатор направления ID "%d".`}

var LevelBudgetNotExist = error_handler.AlterErrorType{Code: 70, ToUserCode: 4004, Object: "LevelBudget",
	CommentTemplate: `Уровень бюджета с ID "%d" не найден.`}

var CompetitiveGroupSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 71, ToUserCode: 17, Object: "competitive_group_Distributed",
	CommentTemplate: `Сумма мест в конкурсных группах превышает значение в распределенном КЦП.`}

var CompetitiveGroupBudgetOSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 72, ToUserCode: 18, Object: "competitive_group_Distributed_BudgetO",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет очной формы обучения превышает значение распределенного КЦП UID "%s".`}

var CompetitiveGroupBudgetOzSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 73, ToUserCode: 19, Object: "competitive_group_Distributed_BudgetOz",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет очно-заочной формы обучения обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupBudgetZSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 74, ToUserCode: 20, Object: "competitive_group_Distributed_BudgetZ",
	CommentTemplate: `Сумма мест в конкурсных группах на бюджет заочной формы обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupPaidOSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 75, ToUserCode: 21, Object: "competitive_group_Distributed_PaidO",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе очной формы обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupPaidOzSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 76, ToUserCode: 22, Object: "competitive_group_Distributed_PaidOz",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе очно-заочной формы обучения обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupPaidZSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 77, ToUserCode: 23, Object: "competitive_group_Distributed_PaidZ",
	CommentTemplate: `Сумма мест в конкурсных группах платной основе заочной формы обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupTargetOSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 78, ToUserCode: 24, Object: "competitive_group_Distributed_TargetO",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления очной формы обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupTargetOzSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 79, ToUserCode: 25, Object: "competitive_group_Distributed_TargetOz",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления очно-заочной формы обучения обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupTargetZSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 80, ToUserCode: 26, Object: "competitive_group_Distributed_TargetZ",
	CommentTemplate: `Сумма мест в конкурсных группах целевого направления заочной формы обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupQuotaOSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 81, ToUserCode: 27, Object: "competitive_group_Distributed_QuotaO",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте очной формы обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupQuotaOzSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 82, ToUserCode: 28, Object: "competitive_group_Distributed_QuotaOz",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте очно-заочной формы обучения обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupQuotaZSpaceSumFromDistributedAdmissionVolume = error_handler.AlterErrorType{Code: 83, ToUserCode: 29, Object: "competitive_group_Distributed_QuotaZ",
	CommentTemplate: `Сумма мест в конкурсных группах по квоте заочной формы обучения превышает значение  распределенного КЦП UID "%s".`}

var CompetitiveGroupNotExist = error_handler.AlterErrorType{Code: 84, ToUserCode: 4020, Object: "competitive_group_Distributed_QuotaZ",
	CommentTemplate: `Конкурсная группа с UID "%s" отсутствует.`}

var CompetitiveGroupSpaceLoverApplicationBudgetO = error_handler.AlterErrorType{Code: 85, ToUserCode: 2030, Object: "competitive_group_Application_BudgetO",
	CommentTemplate: `Колличество поданых заявлений на бюджет очной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationBudgetOz = error_handler.AlterErrorType{Code: 86, ToUserCode: 2031, Object: "competitive_group_Application_BudgetOz",
	CommentTemplate: `Колличество поданых заявлений на бюджет очной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationBudgetZ = error_handler.AlterErrorType{Code: 87, ToUserCode: 2032, Object: "competitive_group_Application_BudgetZ",
	CommentTemplate: `Колличество поданых заявлений на бюджет заочной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationQuotaO = error_handler.AlterErrorType{Code: 88, ToUserCode: 2033, Object: "competitive_group_Application_QuotaO",
	CommentTemplate: `Колличество поданых заявлений по квоте очной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationQuotaOz = error_handler.AlterErrorType{Code: 89, ToUserCode: 2034, Object: "competitive_group_Application_QuotaOz",
	CommentTemplate: `Колличество поданых заявлений по квоте очной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationQuotaZ = error_handler.AlterErrorType{Code: 90, ToUserCode: 2035, Object: "competitive_group_Application_QuotaZ",
	CommentTemplate: `Колличество поданых заявлений по квоте заочной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationPaidO = error_handler.AlterErrorType{Code: 91, ToUserCode: 2036, Object: "competitive_group_Application_PaidO",
	CommentTemplate: `Колличество поданых заявлений на платной основе очной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationPaidOz = error_handler.AlterErrorType{Code: 92, ToUserCode: 2037, Object: "competitive_group_Application_PaidOz",
	CommentTemplate: `Колличество поданых заявлений на платной основе очно-заочной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationPaidZ = error_handler.AlterErrorType{Code: 93, ToUserCode: 2038, Object: "competitive_group_Application_PaidZ",
	CommentTemplate: `Колличество поданых заявлений на платной основе заочной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationTargetO = error_handler.AlterErrorType{Code: 94, ToUserCode: 2039, Object: "competitive_group_Application_TargetO",
	CommentTemplate: `Колличество поданых заявлений целевого направления очной формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationTargetOz = error_handler.AlterErrorType{Code: 95, ToUserCode: 2040, Object: "competitive_group_Application_TargetOz",
	CommentTemplate: `Колличество поданых заявлений целевого направления очно-заочной  формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupSpaceLoverApplicationTargetZ = error_handler.AlterErrorType{Code: 96, ToUserCode: 2041, Object: "competitive_group_Application_TargetZ",
	CommentTemplate: `Колличество поданых заявлений целевого направления заочной  формы обучения превышает значение в редактируемой конкурсной группе UID "%s".`}

var CompetitiveGroupHaveApplication = error_handler.AlterErrorType{Code: 97, ToUserCode: 1005, Object: "competitive_group",
	CommentTemplate: `Невозможно удалить. У конкурсной группы с UID "%s" присутствуют заявления.`}

var AchievementExist = error_handler.AlterErrorType{Code: 98, ToUserCode: 3001, Object: "achievement",
	CommentTemplate: `Идивидуальное достижение с UID "%s" уже существует.`}

var AchievementCategoryNotExist = error_handler.AlterErrorType{Code: 99, ToUserCode: 4009, Object: "achievement",
	CommentTemplate: `Категория идивидуальных достижений с ID "%d" не найдена.`}

var AchievementByCategoryExist = error_handler.AlterErrorType{Code: 100, ToUserCode: 4045, Object: "achievement",
	CommentTemplate: `Идивидуальное достижение с категорией ID "%d" уже существует у кампании с UID "%s".`}

var AchievementNotExist = error_handler.AlterErrorType{Code: 101, ToUserCode: 4046, Object: "achievement",
	CommentTemplate: `Идивидуальное достижение с UID "%s" отсутствует.`}

var AchievementCampaignHaveCompetitiveGroup = error_handler.AlterErrorType{Code: 102, ToUserCode: 2030, Object: "achievement",
	CommentTemplate: `Редактирование не доступно. У приемной кампании с UID "%s" Существуют конкурсные группы UIDs "%s".`}

var CampaignIDNotExist = error_handler.AlterErrorType{Code: 104, ToUserCode: 4085, Object: "campaign",
	CommentTemplate: `Приемная комапания с ID "%d" не найдена.`}

var CompetitiveBenefitsExist = error_handler.AlterErrorType{Code: 105, ToUserCode: 3020, Object: "competitive_benefits",
	CommentTemplate: `Льгота на конкурс с UID "%s" уже существует.`}

var OlympicDiplomaTypeNotExist = error_handler.AlterErrorType{Code: 106, ToUserCode: 4005, Object: "competitive_benefits",
	CommentTemplate: `Тип диплома олимпиад с ID "%d" не найден.`}

var CompetitiveBenefitsAllreadyExist = error_handler.AlterErrorType{Code: 107, ToUserCode: 3021, Object: "competitive_benefits",
	CommentTemplate: `Льгота с  ID  "%d" уже существует у конкурсной группы с UID "%s".`}

var IDOlympicProfilesIsNotUnique = error_handler.AlterErrorType{Code: 108, ToUserCode: 3022, Object: "competitive_benefits",
	CommentTemplate: `У добавляемой льготы UID "%s" повторяющиеся значения профилей олимпиад.`}

var IDOlympicProfilesIsNotExist = error_handler.AlterErrorType{Code: 109, ToUserCode: 4047, Object: "IDOlympicProfiles",
	CommentTemplate: `Сведения об олимпиаде и ее организаторах с ID "%d" не найдены.`}

var CompetitiveBenefitsNotExist = error_handler.AlterErrorType{Code: 110, ToUserCode: 4048, Object: "competitive_benefits",
	CommentTemplate: `Льгота на конкурс с UID "%s" не существует.`}

var CompetitiveBenefitsCompetitiveGroupNotExist = error_handler.AlterErrorType{Code: 111, ToUserCode: 4049, Object: "competitive_benefits",
	CommentTemplate: `Конкурсная группа льготы с UID "%s" не найдена`}

var CompetitiveGroupHaveApplicationEdit = error_handler.AlterErrorType{Code: 112, ToUserCode: 1005, Object: "competitive_group",
	CommentTemplate: `Невозможно редактировать. У конкурсной группы с UID "%s" присутствуют заявления.`}

var CompetitiveGroupProgramExist = error_handler.AlterErrorType{Code: 113, ToUserCode: 3011, Object: "competitive_group_programs",
	CommentTemplate: `Образовательная программа с UID "%s" уже существует.`}

var SubdivisionOrgNotExist = error_handler.AlterErrorType{Code: 114, ToUserCode: 3022, Object: "SubdivisionOrg",
	CommentTemplate: `Кафедра с UID "%s" не существует.`}

var CompetitiveGroupProgramNotExist = error_handler.AlterErrorType{Code: 115, ToUserCode: 4050, Object: "competitive_group_programs",
	CommentTemplate: `Образовательная программы с UID "%s" не существует.`}

var AdmissionVolumeNotExist = error_handler.AlterErrorType{Code: 116, ToUserCode: 4051, Object: "admission_volume",
	CommentTemplate: `Отсутствует КЦП UID "%s" по направлению подготовки ID "%d".`}

var DistributedAdmissionVolumeExist = error_handler.AlterErrorType{Code: 117, ToUserCode: 3023, Object: "distributed_admission_volume",
	CommentTemplate: `Обьем приема с UID "%s"  уже существует.`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpace = error_handler.AlterErrorType{Code: 118, ToUserCode: 30,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема превышает количество мест в КЦП UID "%s".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceBudgetO = error_handler.AlterErrorType{Code: 119, ToUserCode: 31,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" на бюджет очной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceBudgetOz = error_handler.AlterErrorType{Code: 120, ToUserCode: 32,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s"  на бюджет очно-заочной формы обучения обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceBudgetZ = error_handler.AlterErrorType{Code: 121, ToUserCode: 33,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" на бюджет заочной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpacePaidO = error_handler.AlterErrorType{Code: 122, ToUserCode: 34,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s"на платной основе очной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpacePaidOz = error_handler.AlterErrorType{Code: 123, ToUserCode: 35,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" на платной основе очно-заочной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpacePaidZ = error_handler.AlterErrorType{Code: 124, ToUserCode: 36,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" на платной основе заочной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceTargetO = error_handler.AlterErrorType{Code: 125, ToUserCode: 37,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" целевого направления очной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceTargetOz = error_handler.AlterErrorType{Code: 126, ToUserCode: 38,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" целевого направления очно-заочной формы обучения  превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceTargetZ = error_handler.AlterErrorType{Code: 127, ToUserCode: 39,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" целевого направления заочной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceQuotaO = error_handler.AlterErrorType{Code: 128, ToUserCode: 40,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" по квоте очной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceQuotaOz = error_handler.AlterErrorType{Code: 129, ToUserCode: 41,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" по квоте очно-заочной формы обучения обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeSpaceSumMoreDistributedAdmissionVolumeSpaceQuotaZ = error_handler.AlterErrorType{Code: 130, ToUserCode: 42,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест по всем уровням бюджетов обьема приема UIDs:"%s" по квоте заочной формы обучения превышает количество мест в КЦП UID "%s" c идентификатором направления ID "%d".`}

var DistributedAdmissionVolumeNotExist = error_handler.AlterErrorType{Code: 131, ToUserCode: 4052,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Обьем приема с UID "%s" не существует.`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparison = error_handler.AlterErrorType{Code: 132, ToUserCode: 2042,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест, в конкурсных группах, по уровню бюджета обьема приема, превышает значение в изменяемом объеме приема UID "%s".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonBudgetO = error_handler.AlterErrorType{Code: 133, ToUserCode: 2043,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на бюджет очной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonBudgetOz = error_handler.AlterErrorType{Code: 134, ToUserCode: 2044,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s"  на бюджет очно-заочной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonBudgetZ = error_handler.AlterErrorType{Code: 135, ToUserCode: 2045,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на бюджет заочной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonPaidO = error_handler.AlterErrorType{Code: 136, ToUserCode: 2046,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s"на платной основе очной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonPaidOz = error_handler.AlterErrorType{Code: 137, ToUserCode: 2047,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на платной основе очно-заочной формы обучения обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonPaidZ = error_handler.AlterErrorType{Code: 138, ToUserCode: 2048,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на платной основе заочной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonTargetO = error_handler.AlterErrorType{Code: 139, ToUserCode: 2049,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" целевого направления очной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonTargetOz = error_handler.AlterErrorType{Code: 140, ToUserCode: 2050,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" целевого направления очно-заочной формы обучения обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonTargetZ = error_handler.AlterErrorType{Code: 141, ToUserCode: 2051,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" целевого направления заочной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonQuotaO = error_handler.AlterErrorType{Code: 142, ToUserCode: 2052,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" по квоте очной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonQuotaOz = error_handler.AlterErrorType{Code: 143, ToUserCode: 2053,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" по квоте очно-заочной формы обучения обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupAndDistributedAdmissionVolumeSumComparisonQuotaZ = error_handler.AlterErrorType{Code: 144, ToUserCode: 2054,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" по квоте заочной формы обучения превышает значение в редактируемом объеме приема UID  "%s" идентификатор направления ID "%d".`}

var CompetitiveGroupNoZeroComparison = error_handler.AlterErrorType{Code: 145, ToUserCode: 1005,
	Object:          "distributed_admission_volume",
	CommentTemplate: `В конкурсных группах UIDs: "%s" есть распределенные места по уровню бюджета.`}

var CompetitiveGroupNoZeroComparisonBudgetO = error_handler.AlterErrorType{Code: 146, ToUserCode: 1006,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на бюджет очной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonBudgetOz = error_handler.AlterErrorType{Code: 147, ToUserCode: 1007,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на бюджет очно-заочной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonBudgetZ = error_handler.AlterErrorType{Code: 148, ToUserCode: 1008,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на бюджет заочной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonPaidO = error_handler.AlterErrorType{Code: 149, ToUserCode: 1009,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на платной основе очной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonPaidOz = error_handler.AlterErrorType{Code: 150, ToUserCode: 1010,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на платной основе очно-заочной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonPaidZ = error_handler.AlterErrorType{Code: 151, ToUserCode: 1011,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" на платной основе заочной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonTargetO = error_handler.AlterErrorType{Code: 152, ToUserCode: 1012,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" целевого направления очной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonTargetOz = error_handler.AlterErrorType{Code: 153, ToUserCode: 1013,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" целевого направления очно-заочной формы обучения обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonTargetZ = error_handler.AlterErrorType{Code: 154, ToUserCode: 1014,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" целевого направления заочной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonQuotaO = error_handler.AlterErrorType{Code: 155, ToUserCode: 1015,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" по квоте очной формы обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonQuotaOz = error_handler.AlterErrorType{Code: 156, ToUserCode: 1016,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" по квоте очно-заочной формы обучения обучения превышает нулевое значение.`}

var CompetitiveGroupNoZeroComparisonQuotaZ = error_handler.AlterErrorType{Code: 157, ToUserCode: 1017,
	Object:          "distributed_admission_volume",
	CommentTemplate: `Сумма мест в конкурсных группах UIDs: "%s" по квоте заочной формы обучения превышает нулевое значение.`}

var DocNotHaveUID = error_handler.AlterErrorType{Code: 158, ToUserCode: 43,
	Object:          "identification",
	CommentTemplate: `При добалении документа необходимо присвоить UID`}

var EntranceTestBenefitsExist = error_handler.AlterErrorType{Code: 159, ToUserCode: 3015,
	Object:          "entrance_test_benefits",
	CommentTemplate: `Льгота на вступительное испытание с UID "%s" уже существует`}

var EntranceTestNotExist = error_handler.AlterErrorType{Code: 160, ToUserCode: 4035,
	Object:          "entrance_test",
	CommentTemplate: `Вступительное испытание с UID "%s" не найдено`}

var EntranceTestOlympicProfilesIsNotUnique = error_handler.AlterErrorType{Code: 161, ToUserCode: 3,
	Object:          "entrance_test_benefits",
	CommentTemplate: `У добавляемой льготы UID "%s"  повторяющиеся значения профилей олимпиад`}

var EntranceTestBenefitNotExist = error_handler.AlterErrorType{Code: 162, ToUserCode: 4033,
	Object:          "entrance_test_benefits",
	CommentTemplate: `Льгота на вступительное испытание с UID  "%s"  не существует`}

var EntranceTestOlympicProfilesIsNotUniqueEdit = error_handler.AlterErrorType{Code: 163, ToUserCode: 3,
	Object:          "entrance_test_benefits",
	CommentTemplate: `У редактируемой льготы UID "%s"  повторяющиеся значения профилей олимпиад`}

var EntranceTestByBenefitNotExist = error_handler.AlterErrorType{Code: 164, ToUserCode: 4020,
	Object:          "entrance_test_benefits",
	CommentTemplate: `Вступительное испытание льготы с UID "%s" не найдено`}

var CompetitiveGroupByEntranceTestBenefitNotExist = error_handler.AlterErrorType{Code: 165, ToUserCode: 4053,
	Object:          "entrance_test_benefits",
	CommentTemplate: `Конкурсная группа вступительного испытания UID "%s" не найдена`}

var EntranceTestExist = error_handler.AlterErrorType{Code: 166, ToUserCode: 3017,
	Object:          "entrance_test",
	CommentTemplate: `Вступительное испытание с UID "%s" уже существует`}

var ReplaceEntranceTestExist = error_handler.AlterErrorType{Code: 167, ToUserCode: 3023,
	Object:          "entrance_test",
	CommentTemplate: `Заменяемое испытание с  UID  "%s" уже существует у конкурсной группы с UID "%s".`}

var EntranceTestTypeNotExist = error_handler.AlterErrorType{Code: 168, ToUserCode: 4002,
	Object:          "EntranceTestType",
	CommentTemplate: `Тип вступительного испытания с ID "%d" не найден`}

var SubjectNotExist = error_handler.AlterErrorType{Code: 169, ToUserCode: 4006,
	Object:          "entrance_test",
	CommentTemplate: `Предмет с ID "%d" не найден`}

var PriorityCheckFail = error_handler.AlterErrorType{Code: 170, ToUserCode: 4054,
	Object:          "entrance_test",
	CommentTemplate: `Значение поля Priority находится вне диапазона от 1 до 10.`}

var ReplaceEntranceTestNotExist = error_handler.AlterErrorType{Code: 171, ToUserCode: 4054,
	Object:          "entrance_test",
	CommentTemplate: `Испытание Замена с UID  "%s"  не найдено.`}

var EntranceTestWithIDSubjectNotExist = error_handler.AlterErrorType{Code: 172, ToUserCode: 3018,
	Object:          "entrance_test",
	CommentTemplate: `Предмет ЕГЭ с  ID  "%d" уже существует у конкурсной группы с UID  "%s".`}

var EntranceTestSCompetitiveGroupNotExist = error_handler.AlterErrorType{Code: 173, ToUserCode: 3023,
	Object:          "entrance_test",
	CommentTemplate: `Конкурсная группа льготы с UID "%s" не найдена.`}

var ApplicationStatusAdd = error_handler.AlterErrorType{Code: 174, ToUserCode: 1018, Object: "application_status",
	CommentTemplate: `Действие добавление для изменения статуса запрещено.`}

var ApplicationStatusRemove = error_handler.AlterErrorType{Code: 175, ToUserCode: 2055, Object: "application_status",
	CommentTemplate: `Действие удаление для изменения статуса запрещено.`}

var ApplicationStatusNotExist = error_handler.AlterErrorType{Code: 176, ToUserCode: 4009, Object: "achievement",
	CommentTemplate: `Категория идивидуальных достижений с ID "%d" не найдена.`}

var ApplicationNotExist = error_handler.AlterErrorType{Code: 177, ToUserCode: 4056, Object: "application",
	CommentTemplate: `Заявление с %s "%s" не найдено.`}

var UIDReqired = error_handler.AlterErrorType{Code: 178, ToUserCode: 43, Object: "identification",
	CommentTemplate: `Для добавления нового элемента с %s "%s" необходимо присвоить UID.`}

var DocumentTypesNotExist = error_handler.AlterErrorType{Code: 179, ToUserCode: 4057, Object: "DocumentTypes",
	CommentTemplate: `Статус Документа с ID "%d" не найден.`}

var EntrantNotExist = error_handler.AlterErrorType{Code: 180, ToUserCode: 4058, Object: "entrant",
	CommentTemplate: `Абитуриент со СНИЛС "%s" не найден.`}

var IdentificationExist = error_handler.AlterErrorType{Code: 181, ToUserCode: 4059, Object: "identification",
	CommentTemplate: `Документ удостоверяюший личность c "%s" "%s" уже существует.`}

var BadSignTutorial = error_handler.AlterErrorType{Code: 182, ToUserCode: 5001, Object: "jwt",
	CommentTemplate: `Проверка ЭЦП завершилась с результатом %s. Signarure представляет собой подписанный хэш данных, кодированный в base64.Расчёт хеш-суммы, 256 бит - ГОСТ Р 34.11-2012. Расчет подписи, 256 бит - ГОСТ Р 34.10-2012. Подпись передается без контейнера. Подписываемая последовательность представляет собой конкатенацию части Header и Payload разделенных символом точка. Разделитель "." входит в последовательность. Для JWT отправляемых в метод "/api/token/info"  "Header..signature" "." входит в последовательность пример : "Header." . Также, следуют учитывать особенность работы ГОСТ алгоритмов: При использовании API КриптоПро, для формирования хеша используется функция CryptHashData, после формирования подписи функцией CryptSignHash, полученный набор байт следует инвертировать перед кодированием в base64. Т.е. в base64 кодируется массив байт "задом-наперед" от полученного функцией CryptSignHash.`}

var BadSign = error_handler.AlterErrorType{Code: 183, ToUserCode: 5002, Object: "jwt",
	CommentTemplate: `Проверка ЭЦП завершилась с результатом %s. `}
