package entities

import (
	"errors"
	"frdo-check-superservice/model"
)

type RunCheck interface {
	Search(request model.FRDORequest) (string, uint, error)
}

func (e *Entities) NewKindActionType(kindAction int) (RunCheck, error) {
	switch kindAction {
	case 10:
		return NewCheckDoc(e.Service), nil
	case 11:
		return NewGetAllDocs(e.Service), nil
	default:
		return nil, errors.New(`unknown kind action type`)
	}
}
