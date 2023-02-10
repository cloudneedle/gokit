package main

type ErrCode int

func (i ErrCode) Int() int {
	return int(i)
}

//go:generate stringer -type ErrCode -linecomment
const (
	Ok          ErrCode = iota // OK
	Sys                        // 系统错误
	UserNameErr                // 用户名错误
)
