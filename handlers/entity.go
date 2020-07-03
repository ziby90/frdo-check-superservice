package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"time"
)

type UpdateUid struct {
	Changed  time.Time
	Uid      *string
	IdAuthor uint
}

func (result *ResultInfo) EditUid(data digest.EditUid) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var updateUid UpdateUid
	ent := data.Entity
	p, err := ent.GetById(data.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	if p.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Несоответствие организаций`)
		tx.Rollback()
		return
	}
	fmt.Println(data.Uid)
	if data.Uid != nil {
		err = ent.CheckUid(*data.Uid, *p)
		if err != nil {
			result.SetErrorResult(err.Error())
			tx.Rollback()
			return
		}
		updateUid.Uid = data.Uid
	} else {
		updateUid.Uid = nil
	}

	updateUid.Changed = time.Now()
	updateUid.IdAuthor = result.User.Id
	db := tx.Table(p.TableName).Where(`id=?`, data.Id).Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Updates(map[string]interface{}{"changed": updateUid.Changed, "id_author": updateUid.IdAuthor, "uid": updateUid.Uid})
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id`:      data.Id,
		`uid`:     updateUid.Uid,
		`changed`: updateUid.Changed.Format("2006-01-02 15:04:05"),
	}
	tx.Commit()
	return
}
