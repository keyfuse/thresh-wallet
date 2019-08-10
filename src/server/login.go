// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"regexp"
)

type LoginType int

const (
	Unknow LoginType = iota
	Mobile
	Email
)

func loginType(uid string) LoginType {
	mobileRegexp := regexp.MustCompile(`(?:^|[^0-9])(1[34578][0-9]{9})(?:$|[^0-9])`)
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if mobileRegexp.MatchString(uid) {
		return Mobile
	} else if emailRegexp.MatchString(uid) {
		return Email
	}
	return Unknow
}
