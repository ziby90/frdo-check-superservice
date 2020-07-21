package handlers_admin

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

var NewsSearchArray = []string{
	`title`,
}

func (result *ResultInfo) AddNews(data digest.News, files []*multipart.FileHeader) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var news digest.News
	news = data
	news.Created = time.Now()
	news.IdAuthor = result.User.Id

	if news.DateNews.Year() < 1901 {
		news.DateNews = time.Now()
	}

	db := tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&news)
	if db.Error != nil || news.Id == 0 {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}

	var idsFiles []uint
	if len(files) > 0 {
		for _, value := range files {
			var file digest.FileNew
			path := getPath(result.User.Id, file.TableName(), time.Now())
			file.Title = value.Filename
			ext := filepath.Ext(path + `/` + value.Filename)
			file.Mime = ext
			file.Size = value.Size
			file.IdNews = news.Id
			file.IdAuthor = result.User.Id
			file.Created = time.Now()
			sha1FileName := sha1.Sum([]byte(news.TableName() + time.Now().String()))
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
			fileContent, err := value.Open()
			if err != nil {
				result.SetErrorResult(`Невозможно прочитать файл`)
				tx.Rollback()
				return
			} else {
				defer fileContent.Close()
			}

			_, err = io.Copy(dst, fileContent)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
			file.PathFile = name

			db := tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&file)
			if db.Error != nil {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
			idsFiles = append(idsFiles, file.Id)
		}
	}

	result.Items = map[string]interface{}{
		`id_new`: news.Id,
		`files`:  idsFiles,
	}
	result.Done = true
	tx.Commit()
	return
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
	db := conn.Where(`id > 0`)
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
		if db.Error.Error() == `record not found` {
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
	db := conn.Where(`id=?`, idNew).Find(&news)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
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

func (result *ResultInfo) AddFileNew(ID uint, files []*multipart.FileHeader) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var doc digest.News
	db := conn.Where(`id=?`, ID).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Новость не найдена."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}

	var idsFiles []uint
	if len(files) > 0 {
		for _, value := range files {
			var file digest.FileNew
			path := getPath(result.User.Id, file.TableName(), time.Now())
			file.Title = value.Filename
			ext := filepath.Ext(path + `/` + value.Filename)
			file.Mime = ext
			file.Size = value.Size
			file.IdNews = doc.Id
			file.IdAuthor = result.User.Id
			file.Created = time.Now()
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
			fileContent, err := value.Open()
			if err != nil {
				result.SetErrorResult(`Невозможно прочитать файл`)
				tx.Rollback()
				return
			} else {
				defer fileContent.Close()
			}

			_, err = io.Copy(dst, fileContent)
			if err != nil {
				result.SetErrorResult(err.Error())
				tx.Rollback()
				return
			}
			file.PathFile = name

			db := tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&file)
			if db.Error != nil {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
			idsFiles = append(idsFiles, file.Id)
		}
	} else {
		result.SetErrorResult(`Где файлы, Лебовски?`)
		tx.Rollback()
		return
	}
	result.Items = map[string]interface{}{
		`ids_files`: idsFiles,
	}

	result.Done = true
	tx.Commit()
	return
}

func (result *ResultInfo) RemoveFileNew(ID uint) {
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

	db = conn.Exec(`DELETE FROM info.files WHERE id=?`, doc.Id)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		return
	}

	result.Done = true
	return
}

func (result *ResultInfo) EditNew(data digest.News) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var old digest.News
	var new digest.News
	db := tx.Where(`id=?`, data.Id).Find(&old)
	if old.Id <= 0 {
		result.SetErrorResult(`Новость не найдена`)
		tx.Rollback()
		return
	}
	new = old
	if data.Title != nil {
		s := strings.TrimSpace(*data.Title)
		new.Title = &s
	}

	new.Created = time.Now()
	new.IdAuthor = result.User.Id
	new.Published = data.Published
	new.Deleted = data.Deleted
	new.Content = data.Content
	if data.DateNews.Year() < 1901 {
		new.DateNews = time.Now()
	} else {
		new.DateNews = data.DateNews
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&new)
	if db.Error != nil {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
	result.Done = true
	result.Items = map[string]interface{}{
		`id_new`: new.Id,
	}
	tx.Commit()
	return

}
