// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"xlog"
)

type vcode struct {
	code string
	uid  string
	then time.Time
}

// Vcode --
type Vcode struct {
	mu     sync.Mutex
	log    *xlog.Log
	conf   *Config
	vcodes map[string]*vcode
}

// NewVcode -- creates new Vcode.
func NewVcode(log *xlog.Log, conf *Config) *Vcode {
	return &Vcode{
		log:    log,
		conf:   conf,
		vcodes: make(map[string]*vcode),
	}
}

// Add -- used to add a new <uid, code> pair to vcode pool.
func (vc *Vcode) Add(uid string, code string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	vcode := &vcode{
		code: code,
		uid:  uid,
		then: time.Now(),
	}
	vc.vcodes[uid] = vcode
}

// Remove -- used to remove the code of the uid from the vcode pool.
func (vc *Vcode) Remove(uid string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	delete(vc.vcodes, uid)
}

// Check -- used to check the code valid or not with the vcode in the pool.
func (vc *Vcode) Check(uid string, code string) error {
	log := vc.log

	vc.mu.Lock()
	defer vc.mu.Unlock()
	expired := vc.conf.VCodeExpired

	vcode, ok := vc.vcodes[uid]
	if !ok {
		return fmt.Errorf("vcode.check[%s,%s].does.not.exists", uid, code)
	}

	if strings.Compare(vcode.code, code) != 0 {
		log.Error("vcode.check[%s].vcode.[%s]!=req.[%s]", uid, vcode.code, code)
		return fmt.Errorf("vcode.uid[%s].vcode[%s].invalid", uid, code)
	}

	dur := time.Since(vcode.then)
	if dur > time.Duration(expired)*time.Second {
		return fmt.Errorf("vcode.check[%s].expired[%+v]", uid, dur)
	}
	return nil
}
