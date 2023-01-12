package modules

import "errors"

type IReqLimitCode uint64

type ICodeReqLimit map[IReqLimitCode]*IReqLimit

var code IReqLimitCode = 2023

var GlobalCodeReqLimit ICodeReqLimit

func GetNewReqLimitCode() IReqLimitCode {
	code = code + 1
	return code
}

func CacheCodeReqLimit(code IReqLimitCode, reqLimit *IReqLimit) error {
	if GlobalCodeReqLimit == nil {
		GlobalCodeReqLimit = make(ICodeReqLimit)
	}
	if GlobalCodeReqLimit[code] != nil {
		return errors.New("duplicate code")
	}
	GlobalCodeReqLimit[code] = reqLimit
	return nil
}
