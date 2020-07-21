package handlers

import (
	"bufio"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io/ioutil"
	"os"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

var NewsSearchArray = []string{
	`title`,
}

func (result *Result) GetListNews() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var news []digest.News
	sortField := `date_news`
	sortOrder := `desc`
	if result.Sort.Field != `` {
		sortField = result.Sort.Field
	} else {
		result.Sort.Field = sortField
	}

	fmt.Print(result.Sort.Field, sortField)
	db := conn.Where(`published IS true AND deleted IS NOT true`)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], NewsSearchArray) >= 0 {
			db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}
	dbCount := db.Model(&news).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Order(sortField + ` ` + sortOrder).Find(&news)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Новости не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, new := range news {
			c := map[string]interface{}{
				`id`:        new.Id,
				`title`:     new.Title,
				`content`:   new.Content,
				`date_news`: new.DateNews,
				`published`: new.Published,
				`deleted`:   new.Deleted,
				`id_author`: new.IdAuthor,
				`created`:   new.Created,
			}
			var files []digest.FileNew
			db = conn.Where(`id_news=?`, new.Id).Find(&files)
			var f []interface{}
			for _, file := range files {
				f = append(f, map[string]interface{}{
					`id`:    file.Id,
					`title`: file.Title,
					`size`:  file.Size,
					`type`:  file.Mime,
				})
			}
			c[`files`] = f
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Новости не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) GetInfoNew(idNew uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var news digest.News
	db := conn.Where(`id=? AND published IS true AND deleted IS NOT true`, idNew).Find(&news)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Новость не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		c := map[string]interface{}{
			`id`:        news.Id,
			`title`:     news.Title,
			`content`:   news.Content,
			`date_news`: news.DateNews,
			`published`: news.Published,
			`deleted`:   news.Deleted,
			`id_author`: news.IdAuthor,
			`created`:   news.Created,
		}
		var files []digest.FileNew
		db = conn.Where(`id_news=?`, idNew).Find(&files)
		var f []interface{}
		for _, file := range files {
			f = append(f, map[string]interface{}{
				`id`:    file.Id,
				`title`: file.Title,
				`size`:  file.Size,
				`type`:  file.Mime,
			})
		}
		c[`files`] = f
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Новость не найдена.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) GetFileNew(ID uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.FileNew
	db := conn.Where(`id=?`, ID).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Файл не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}

	var news digest.News
	db = conn.Where(`id = ?  AND published IS true AND deleted IS NOT true`, doc.IdNews).Find(&news)
	if news.Id <= 0 {
		message := "Файл не найден"
		result.Message = &message
		return
	}

	filename := doc.PathFile

	path := `./` + getPath(doc.IdAuthor, doc.TableName(), doc.Created) + filename
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

	result.Done = true
	return
}
