package modules

type ServiceContentCode int

type LocationDetail struct {
	Code         ServiceContentCode
	Reg          string
	Pass         []string
	Balance      string
	ReqLimit     []ReqLimitInfo
	RedirectCode int
	RedirectUrl  string
}

type ServerDetail struct {
	Code     ServiceContentCode
	Id       string
	Listen   uint16
	Name     []string
	SSLKey   string
	SSLCert  string
	ReqLimit []ReqLimitInfo
}

type ServiceDetail struct {
	Location LocationDetail
	Server   ServerDetail
}

type ReqLimitType string

type ReqLimitInfo struct {
	Code         ServiceContentCode
	Type         ReqLimitType
	Time         uint64
	CountPerTime uint64
}

type NameLocation map[string][]ServiceDetail
type ServerLocation map[uint16]NameLocation

type CodeReqLimitInfoMap map[ServiceContentCode][]ReqLimitInfo
