// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"
	"net/mail"
	"strings"

	"xlog"

	"github.com/txthinking/mailx"
)

type Smtp struct {
	log  *xlog.Log
	conf *Config
}

func NewSmtp(log *xlog.Log, conf *Config) *Smtp {
	return &Smtp{
		log:  log,
		conf: conf,
	}
}

func (smtp *Smtp) Backup(uid string, subject string) error {
	conf := smtp.conf

	if conf.Smtp != nil {
		attachment := fmt.Sprintf("%v/%v.json", conf.DataDir, uid)
		tos := strings.Split(conf.Smtp.BackupTo, ",")

		server := &mailx.SMTP{
			Server:   conf.Smtp.Server,
			Port:     conf.Smtp.Port,
			UserName: conf.Smtp.UserName,
			Password: conf.Smtp.Password,
		}

		to := make([]*mail.Address, 0)
		for _, email := range tos {
			to = append(to, &mail.Address{Address: email})
		}
		message := &mailx.Message{
			From: &mail.Address{
				Name: subject,
			},
			To:      to,
			Subject: fmt.Sprintf("%v-%v", conf.ChainNet, uid),
			Body:    "",
			Attachment: []string{
				attachment,
			},
		}
		return server.Send(message)
	}
	return nil
}
