package sql

import (
	"errors"
	"fmt"
	"frdo-check-superservice/model"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

type Search struct {
	dbs map[string]*sqlx.DB
}

func NewSearch(dbs map[string]*sqlx.DB) *Search {
	return &Search{dbs: dbs}
}

func (s *Search) GetDbConnections() map[string]*sqlx.DB {
	return s.dbs
}

func (s *Search) SearchPerson(person model.Person, nameModule string) ([]uint, error) {
	var ids []uint
	if _, ok := s.dbs[nameModule]; !ok {
		logrus.Println(nameModule, s.dbs[nameModule])
		return ids, errors.New(fmt.Sprintf(`module name %s not found in connects`, nameModule))
	}
	connect := s.dbs[nameModule]

	query, args := GetStringPerson(person, nameModule)

	err := connect.Select(&ids, query, args...)
	if err != nil {
		return ids, err
	}
	return ids, nil
}

func (s *Search) SearchDocument(recipientId uint, document model.EduDocument) ([]uint, error) {
	var ids []uint
	if _, ok := s.GetEducationLevel()[int(document.EduLevel)]; !ok {
		return ids, errors.New(`not found module for education level`)
	}
	if _, ok := s.dbs[s.GetEducationLevel()[int(document.EduLevel)]]; !ok {
		return ids, errors.New(fmt.Sprintf(`module name %s not found in connects`, s.GetEducationLevel()[int(document.EduLevel)]))
	}
	nameModule := s.GetEducationLevel()[int(document.EduLevel)]
	connect := s.dbs[nameModule]
	query, args := GetStringDocument(document, recipientId, nameModule)
	err := connect.Select(&ids, query, args...)
	if err != nil {
		return ids, err
	}
	return ids, nil
}

func (s *Search) SearchAllDocument(recipientId uint, nameModule string) ([]model.EduDocument, error) {
	var docs []model.EduDocument
	if _, ok := s.dbs[nameModule]; !ok {
		return docs, errors.New(fmt.Sprintf(`module name %s not found in connects`, nameModule))
	}
	connect := s.dbs[nameModule]
	query := fmt.Sprintf(` SELECT  doc.series , doc.number, org.ogrn, org.full_title, org.short_title, date(doc.issue_date) FROM persons.documents doc JOIN admin.organizations org ON org.id=doc.organization_id WHERE recipient_id=$1`)
	//query, args := GetStringDocument(document, recipientId, nameModule)
	logrus.Println(query)
	rows, err := connect.Queryx(query, recipientId)

	if err != nil {
		return docs, err
	}
	for rows.Next() {
		var newDocument model.EduDocument
		for key, value := range s.GetEducationLevel() {
			if value == nameModule {
				newDocument.EduLevel = uint(key)
			}
		}
		err = rows.Scan(&newDocument.Series, &newDocument.Number, &newDocument.OGRN, &newDocument.FullOrgName, &newDocument.ShortOrgName, &newDocument.DocumentIssueDate)
		if err != nil {
			return nil, err
		}
		docs = append(docs, newDocument)
	}

	return docs, nil
}

func (s *Search) AddNewQueue(smevMsg model.SMEVMessage, typeQueue string) (uint, error) {
	//var newId uint
	connect, err := s.dbs[`CHECK_SMEV`]
	if !err {
		return smevMsg.Id, errors.New(`not found connect to DB_CHECK_SMEV`)
	}
	fields := []string{`package_id`, `smev_id`, `message_id`, `"type"`, `"xml"`, `error`, `status`, `kind_action_id`}
	values := []string{`$1`, `$2`, `$3`, `$4`, `$5`, `$6`, `$7`, `$8`}
	args := []interface{}{smevMsg.PackageId, smevMsg.SmevId, smevMsg.MessageId, smevMsg.Type, smevMsg.Xml, smevMsg.Error, smevMsg.Status, smevMsg.KindActionId}
	nameTable := `public.queue_in`
	if typeQueue == `OUT` {
		nameTable = `public.queue_out`
		fields = append(fields, `status_response_id`)
		values = append(values, `$9`)
		args = append(args, smevMsg.StatusResponseId)
	}
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RETURNING id`, nameTable, strings.Join(fields, `,`), strings.Join(values, `,`))
	//logrus.Printf("updateQuery: %v", query)
	//logrus.Printf("args: %v", args)
	row := connect.QueryRow(query, args...)
	if errInsert := row.Scan(&smevMsg.Id); errInsert != nil {
		return smevMsg.Id, errInsert
	}
	return smevMsg.Id, nil
}
func (s *Search) UpdateQueueLog(smevMsg model.SMEVMessage, typeQueue string) error {
	connect, err := s.dbs[`CHECK_SMEV`]
	if !err {
		return errors.New(`not found connect to DB_CHECK_SMEV`)
	}

	nameTable := `public.queue_in`
	if typeQueue == `OUT` {
		nameTable = `public.queue_out`
	}
	query := fmt.Sprintf(`UPDATE %s set updated_at=$1, error=$2 WHERE id=$3`, nameTable)
	args := []interface{}{time.Now(), smevMsg.Error, smevMsg.Id}

	_, errUpdate := connect.Exec(query, args...)
	return errUpdate
}

func GetStringPerson(p model.Person, nameModule string) (string, []interface{}) {
	var args []interface{}
	var argId = 3
	query := fmt.Sprintf(`Select id FROM persons.recipients WHERE upper(surname)=$1 AND upper(name)=$2`)
	args = append(args, strings.ToUpper(p.Surename))
	args = append(args, strings.ToUpper(p.Name))
	if p.Birthday != nil {
		birthday, err := time.Parse("2006-01-02", *p.Birthday)
		if err == nil {
			query += fmt.Sprintf(` AND birthday=$%d`, argId)
			args = append(args, birthday.Format("2006-01-02"))
			argId++
		}
	}
	if p.Patronymic != nil {
		query += fmt.Sprintf(` AND upper(patronymic)=$%d`, argId)
		args = append(args, strings.ToUpper(*p.Patronymic))
		argId++
	}
	if p.Snils != nil {
		reg := regexp.MustCompile("[^0-9]+")
		replaceStr := reg.ReplaceAllString(*p.Snils, "")
		query += fmt.Sprintf(` AND regexp_replace(snils, '[^0-9]+'::text, ''::text, 'g'::text)=$%d`, argId)
		args = append(args, replaceStr)
		argId++
	}
	if p.Gender != nil {
		query += fmt.Sprintf(` AND sex=$%d`, argId)
		args = append(args, *p.Gender)
		argId++
	}
	//if p.RFCitizen != nil {
	//	query += fmt.Sprintf(` AND nationality_id=$%d`, argId)
	//	args = append(args, *p.RFCitizen)
	//	argId++
	//}

	return query, args
}

func GetStringDocument(d model.EduDocument, recipientId uint, nameModule string) (string, []interface{}) {
	var args []interface{}
	var argId = 2
	query := fmt.Sprintf(`Select a.id FROM persons.documents a JOIN admin.organizations org ON org.id=a.organization_id WHERE recipient_id=$1`)
	args = append(args, recipientId)
	if d.Number != `` {
		query += fmt.Sprintf(` AND a.number=$%d`, argId)
		args = append(args, d.Number)
		argId++
	}
	if d.DocumentIssueDate != nil {
		issueDate, err := time.Parse("2006-01-02", *d.DocumentIssueDate)
		if err != nil {
			query += fmt.Sprintf(` AND date(a.issue_date)=$%d`, argId)
			args = append(args, issueDate)
			argId++
		}
	}
	if d.Series != nil {
		query += fmt.Sprintf(` AND a.series=$%d`, argId)
		args = append(args, *d.Series)
		argId++
	}
	if d.OGRN != nil {
		query += fmt.Sprintf(` AND org.ogrn=$%d`, argId)
		args = append(args, *d.OGRN)
		argId++
	}
	if d.KPP != nil {
		query += fmt.Sprintf(` AND org.kpp=$%d`, argId)
		args = append(args, *d.KPP)
		argId++
	}
	return query, args
}

func (s *Search) GetEducationLevel() map[int]string {
	r := make(map[int]string)
	r[1] = `VPO`
	r[2] = `SCHOOL`
	r[3] = `SPO`
	r[4] = `MIG`
	r[5] = `PO`
	r[6] = `DPO`
	return r
}
