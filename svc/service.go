package svc

type Service interface {
	UpperCase(string2 string)(string,error)
	Count(string2 string) int
}


