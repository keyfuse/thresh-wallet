// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"crypto/rand"
	"fmt"
	"net/mail"
	"strings"

	"xlog"

	"github.com/txthinking/mailx"
)

// Smtp --
type Smtp struct {
	log  *xlog.Log
	conf *Config
}

// NewSmtp -- creates new Smtp.
func NewSmtp(log *xlog.Log, conf *Config) *Smtp {
	return &Smtp{
		log:  log,
		conf: conf,
	}
}

// Backup -- used to backup the user wallet json file via smtp.
func (smtp *Smtp) Backup(uid string, name string) error {
	log := smtp.log
	conf := smtp.conf

	if conf.Smtp != nil {
		go func(conf *Config) {
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
					Name: name,
				},
				To:      to,
				Subject: fmt.Sprintf("%v-%v", conf.ChainNet, uid),
				Body:    "",
				Attachment: []string{
					attachment,
				},
			}
			if err := server.Send(message); err != nil {
				log.Error("smtp.backup.send[%v].error:%+v", uid, err)
			}
		}(conf)
	}
	return nil
}

// VCode -- send vcode to email.
func (smtp *Smtp) VCode(uid string, name string, vcode string) error {
	log := smtp.log
	conf := smtp.conf

	if conf.Smtp != nil {
		seed := make([]byte, 16)
		rand.Read(seed)

		go func(conf *Config) {
			server := &mailx.SMTP{
				Server:   conf.Smtp.Server,
				Port:     conf.Smtp.Port,
				UserName: conf.Smtp.UserName,
				Password: conf.Smtp.Password,
			}

			message := &mailx.Message{
				From: &mail.Address{
					Name: fmt.Sprintf("%s-No-Reply-%x", name, seed),
				},
				To: []*mail.Address{
					&mail.Address{Address: uid},
				},
				Subject: "KeyFuse ID Verification Code",
				Body:    fmt.Sprintf("Your KeyFuse ID Verification Code is: <b>%v</b>", vcode),
			}
			if err := server.Send(message); err != nil {
				log.Error("smtp.vcode[%v].send[%v].error:%+v", vcode, uid, err)
			}
		}(conf)
	}
	return nil
}
