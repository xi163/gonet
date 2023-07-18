package result

import "github.com/cwloo/gonet/utils/validator"

type R struct {
	Code   int    `json:"code" form:"code"`
	ErrMsg string `json:"errmsg" form:"errmsg"`
	Req    any    `json:"request,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func (s *R) Ok() bool {
	return s.Code == 0
}

func (s *R) Empty() bool {
	return validator.Empty(s.Data)
}
