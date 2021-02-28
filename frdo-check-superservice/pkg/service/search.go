package service

import (
	"errors"
	"fmt"
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/repository"
)

type SearchService struct {
	repo repository.SearchInDb
}

func NewSearchService(repo repository.SearchInDb) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) Search(person model.Person, document model.EduDocument) ([]model.ResultDb, error) {
	var result []model.ResultDb
	if _, ok := s.repo.GetEducationLevel()[int(document.EduLevel)]; !ok {
		return nil, errors.New(`not found module for education level`)
	}
	if _, ok := s.repo.GetDbConnections()[s.repo.GetEducationLevel()[int(document.EduLevel)]]; !ok {
		return nil, errors.New(fmt.Sprintf(`module name %s not found in connects`, s.repo.GetEducationLevel()[int(document.EduLevel)]))
	}
	nameModule := s.repo.GetEducationLevel()[int(document.EduLevel)]

	idRecipients, err := s.repo.SearchPerson(person, nameModule)
	if err != nil {
		return nil, err
	}
	for _, idRecipient := range idRecipients {
		idDocuments, err := s.repo.SearchDocument(idRecipient, document)
		if err != nil {
			return nil, err
		}
		result = append(result, model.ResultDb{
			IdRecipient: idRecipient,
			IdDocuments: idDocuments,
			NameModule:  nameModule,
		})
	}
	return result, nil
}

func (s *SearchService) GetDocuments(person model.Person) ([]model.EduDocument, error) {
	var result []model.EduDocument
	var personsCount int = 0
	for nameModule, _ := range s.repo.GetDbConnections() {
		if nameModule == `CHECK_SMEV` {
			continue
		}
		idRecipients, err := s.repo.SearchPerson(person, nameModule)
		if err != nil {
			return nil, err
		}
		switch len(idRecipients) {
		case 0:
			break
		case 1:
			//if len(idRecipients) == 1{
			personsCount++
			for _, idRecipient := range idRecipients {
				var docs []model.EduDocument
				docs, err = s.repo.SearchAllDocument(idRecipient, nameModule)
				if err != nil {
					return nil, err
				}
				result = append(result, docs...)
			}
		//}
		default:
			return result, errors.New(`any more persons`)
		}
	}
	if personsCount == 0 {
		return result, errors.New(`not found a person`)
	}

	return result, nil
}

func (s *SearchService) AddNewQueueLog(smevMsg model.SMEVMessage, typeQueue string) (uint, error) {
	id, err := s.repo.AddNewQueue(smevMsg, typeQueue)
	return id, err
}

func (s *SearchService) UpdateQueueLog(smevMsg model.SMEVMessage, typeQueue string) error {
	err := s.repo.UpdateQueueLog(smevMsg, typeQueue)
	return err
}

func GetResponse(results []model.ResultDb) (uint, string) {
	switch len(results) {
	case 0:
		return uint(2), `документ по реквизитам не найден`
	case 1:
		if len(results[0].IdDocuments) > 1 {
			return uint(3), `необходимо дополнить запрос дополнительными данными`
		}
		return uint(1), `документ по  реквизитам найден`
	default:
		return uint(3), `необходимо дополнить запрос дополнительными данными`
	}
}
func GetResponseAllDocs(results []model.EduDocument) (uint, string) {
	if len(results) > 0 {
		return uint(1), `документы найдены`
	} else {
		return uint(2), `документы не найдены`
	}
}
