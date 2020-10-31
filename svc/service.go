package svc

import "strings"

type Service interface {
	UpperCase(s string)(string,error)
	Count(s  string) int
}

type svcImpl struct {
	username string
}


var _ Service = (*svcImpl)(nil)

func New(username string) Service {
	return &svcImpl{username: username}
}


func (svc *svcImpl) UpperCase(s string) (string, error) {

	if s == ""{
		return "",ErrStringEmpty
	}


	return strings.ToUpper(s),nil


}

func (svc *svcImpl) Count(s string) int {
	return len(s)
}

