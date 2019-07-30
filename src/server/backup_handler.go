// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"proto"
)

func (h *Handler) backupVCode(w http.ResponseWriter, r *http.Request) {
	log := h.log
	vcode := h.backupCode
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("backupVCode", r)
	if err != nil {
		log.Error("api.backup.vcode.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.BackupVCodeRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.backup.vcode[%v].decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.backup.vcode.req:%+v", req)

	seed := make([]byte, 256)
	if _, err := rand.Read(seed); err != nil {
		log.Error("api.backup.vcode[%v].rand.seed.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	sha256 := sha256.Sum256(seed)
	code := fmt.Sprintf("%x", sha256)
	rsp := proto.BackupVCodeResponse{
		VCode: code,
	}
	vcode.Add(uid, code)
	resp.writeJSON(rsp)
}

func (h *Handler) backupStore(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	smtp := h.smtp
	vcode := h.backupCode
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("backupStore", r)
	if err != nil {
		log.Error("api.backup.store.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.BackupStoreRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.backup.store[%v].decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.backup.store.req:%+v", req)

	// Check.
	backup, err := wdb.GetBackup(uid)
	if err != nil {
		log.Error("api.backup.store[%v].pubkey.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	if backup.EncryptionPubKey != "" {
		log.Error("api.backup[%v].check.exists", uid)
		resp.writeErrorWithStatus(400, fmt.Errorf("backup.exists"))
		return
	}

	// Verify.
	{
		// vcode.
		if err := vcode.Check(uid, req.VCode); err != nil {
			log.Error("api.backup.store.vcode.error:%+v", err)
			resp.writeErrorWithStatus(400, err)
			return
		}

		// signature.
		if err := rsaVerify(req.EncryptionPubKey, req.VCode, req.Signature); err != nil {
			log.Error("api.backup.store.rsa.verify.error:%+v", err)
			resp.writeErrorWithStatus(400, err)
			return
		}
		vcode.Remove(uid)
	}

	// wdb Backup.
	if err := wdb.StoreBackup(uid, req.Email, req.DeviceID, req.CloudService, req.EncryptedPrvKey, req.EncryptionPubKey); err != nil {
		log.Error("api.backup.wdb.store.backup.error:%+v", err)
		resp.writeErrorWithStatus(500, err)
		return
	}

	// smtp backup.
	if err := smtp.Backup(uid, "KeyFuse Server Store Backup"); err != nil {
		log.Error("api.backup.wdb.store.backup.smtp.error:%+v", err)
		resp.writeErrorWithStatus(500, nil)
		return
	}
	rsp := &proto.BackupStoreResponse{}
	log.Info("api.backup.store.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) backupRestore(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	vcode := h.backupCode
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("backupRestore", r)
	if err != nil {
		log.Error("api.backup.store.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.BackupRestoreRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.backup.restore[%v].decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.backup.restore.req:%+v", req)

	// Check.
	backup, err := wdb.GetBackup(uid)
	if err != nil {
		log.Error("api.backup.restore[%v].get.backup.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	if backup.EncryptionPubKey == "" {
		log.Error("api.backup.restore[%v].encryption.pubkey.not.exists", uid)
		resp.writeErrorWithStatus(400, fmt.Errorf("backup.encryption.pubkey.not.exists"))
		return
	}

	// Verify.
	{
		// vcode.
		if err := vcode.Check(uid, req.VCode); err != nil {
			log.Error("api.backup.restore.vcode.error:%+v", err)
			resp.writeErrorWithStatus(400, err)
			return
		}

		// signature.
		if err := rsaVerify(backup.EncryptionPubKey, req.VCode, req.Signature); err != nil {
			log.Error("api.backup.restore.rsa.verify.error:%+v", err)
			resp.writeErrorWithStatus(400, err)
			return
		}
	}

	// OK.
	vcode.Remove(uid)
	rsp := &proto.BackupRestoreResponse{
		Time:            backup.Time,
		EncryptedPrvKey: backup.EncryptedPrvKey,
	}
	log.Info("api.backup.restore.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}
