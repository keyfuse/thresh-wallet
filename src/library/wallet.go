// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package library

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"proto"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xbase"
	"github.com/keyfuse/tokucore/xcore"
	"github.com/keyfuse/tokucore/xcore/bip32"
	"github.com/keyfuse/tokucore/xcrypto"
	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
)

// WalletCheckResponse --
type WalletCheckResponse struct {
	Status
	WalletExists       bool   `json:"wallet_exists"`
	BackupExists       bool   `json:"backup_exists"`
	ForceRecover       bool   `json:"force_recover"`
	BackupTimestamp    int64  `json:"backup_timestamp"`
	BackupCloudService string `json:"backup_cloudservice"`
}

// APIWalletCheck -- check the user/backup exists.
func APIWalletCheck(url string, token string) string {
	rsp := &WalletCheckResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/check", url)

	req := &proto.WalletCheckRequest{}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	ret := &proto.WalletCheckResponse{}
	if err := httpRsp.Json(ret); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.WalletExists = ret.WalletExists
	rsp.BackupExists = ret.BackupExists
	rsp.ForceRecover = ret.ForceRecover
	rsp.BackupTimestamp = ret.BackupTimestamp
	rsp.BackupCloudService = ret.BackupCloudService
	return marshal(rsp)
}

// WalletCreateResponse --
type WalletCreateResponse struct {
	Status
}

// APIWalletCreate -- create the wallet.
func APIWalletCreate(url string, token string, masterPrvKey string) string {
	var signature string
	var masterPubKey string

	rsp := &WalletCreateResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/create", url)

	// Master pravite key.
	{
		masterkey, err := bip32.NewHDKeyFromString(masterPrvKey)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}

		// Pubkey.
		var net *network.Network
		switch pre := masterPrvKey[:4]; pre {
		case "xprv":
			net = network.MainNet
		case "tprv":
			net = network.TestNet
		}
		masterPubKey = masterkey.HDPublicKey().ToString(net)

		hash := sha256.Sum256([]byte(masterPubKey))
		sig, err := xcrypto.EcdsaSign(masterkey.PrivateKey(), hash[:])
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		signature = fmt.Sprintf("%x", sig)
	}

	req := &proto.WalletCreateRequest{
		Signature:    signature,
		MasterPubKey: masterPubKey,
	}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	ret := &proto.WalletCreateResponse{}
	if err := httpRsp.Json(ret); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	return marshal(rsp)
}

// WalletPortfolioResponse --
type WalletPortfolioResponse struct {
	Status
	CoinSymbol   string  `json:"coin_symbol"`
	FiatSymbol   string  `json:"fiat_symbol"`
	CurrentPrice float64 `json:"current_price"`
}

// APIWalletPortfolio -- portfolio api.
func APIWalletPortfolio(url string, token string) string {
	rsp := &WalletPortfolioResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/portfolio", url)

	req := &proto.WalletPortfolioRequest{}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	ret := &proto.WalletPortfolioResponse{}
	if err := httpRsp.Json(ret); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.CoinSymbol = ret.CoinSymbol
	rsp.FiatSymbol = ret.FiatSymbol
	rsp.CurrentPrice = ret.CurrentPrice
	return marshal(rsp)
}

// WalletBalanceResponse --
type WalletBalanceResponse struct {
	Status
	CoinValue uint64 `json:"coin_value"`
}

// APIWalletBalance -- Wallet balance api.
func APIWalletBalance(url string, token string) string {
	rsp := &WalletBalanceResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/balance", url)

	req := &proto.WalletBalanceRequest{}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	balance := &proto.WalletBalanceResponse{}
	if err := httpRsp.Json(balance); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.CoinValue = balance.CoinValue
	return marshal(rsp)
}

// WalletNewAddressResponse --
type WalletNewAddressResponse struct {
	Status
	Pos     uint32 `json:"pos"`
	Address string `json:"address"`
}

// APIWalletNewAddress -- new address api.
func APIWalletNewAddress(url string, token string) string {
	rsp := &WalletNewAddressResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/newaddress", url)

	req := &proto.WalletNewAddressRequest{}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	address := &proto.WalletNewAddressResponse{}
	if err := httpRsp.Json(address); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.Pos = address.Pos
	rsp.Address = address.Address
	return marshal(rsp)
}

// WalletTxsResponse --
type WalletTxsResponse struct {
	Status
	Txs []proto.WalletTxsResponse `json:"txs"`
}

// APIWalletTxs -- get the txs.
func APIWalletTxs(url string, token string, offset int, limit int) string {
	rsp := &WalletTxsResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/txs", url)

	req := &proto.WalletTxsRequest{
		Offset: offset,
		Limit:  limit,
	}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	var txsRsp []proto.WalletTxsResponse
	if err := httpRsp.Json(&txsRsp); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.Txs = txsRsp
	return marshal(rsp)
}

// WalletAddressesResponse --
type WalletAddressesResponse struct {
	Status
	Addresses []proto.WalletAddressesResponse `json:"addresses"`
}

// APIWalletAddresses -- get the address list.
func APIWalletAddresses(url string, token string, offset int, limit int) string {
	rsp := &WalletAddressesResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/addresses", url)

	req := &proto.WalletAddressesRequest{
		Offset: offset,
		Limit:  limit,
	}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	var addrsRsp []proto.WalletAddressesResponse
	if err := httpRsp.Json(&addrsRsp); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.Addresses = addrsRsp
	return marshal(rsp)
}

// WalletPrepareSendResponse --
type WalletSendFeesResponse struct {
	Status
	Fees          uint64 `json:"fees"`
	FeeMode       string `json:"feemode"`
	TotalValue    uint64 `json:"total_value"`
	SendableValue uint64 `json:"sendable_value"`
}

// APIWalletSendFees -- used to prepare the fees before the txn build.
func APIWalletSendFees(url string, token string, sendValue uint64) string {
	feemode := "fast"

	rsp := &WalletSendFeesResponse{}
	rsp.Code = http.StatusOK

	// Get sendfees.
	{
		req := &proto.WalletSendFeesRequest{
			Priority:  feemode,
			SendValue: sendValue,
		}
		path := fmt.Sprintf("%s/api/wallet/sendfees", url)
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}

		feesRsp := &proto.WalletSendFeesResponse{}
		if err := httpRsp.Json(feesRsp); err != nil {
			rsp.Code = httpRsp.StatusCode()
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		rsp.Fees = feesRsp.Fees
		rsp.FeeMode = feemode
		rsp.TotalValue = feesRsp.TotalValue
		rsp.SendableValue = feesRsp.SendableValue
	}
	return marshal(rsp)
}

// WalletSendResponse --
type WalletSendResponse struct {
	Status
	TxID string `json:"txid"`
}

func APIWalletSend(url string, token string, chainnet string, masterPrvKey string, toAddress string, amount uint64, fees uint64, msg string) string {
	var err error
	var to xcore.Address
	var change xcore.Address
	var masterkey *bip32.HDKey
	var unspents []proto.WalletUnspentResponse

	rsp := &WalletSendResponse{}
	rsp.Code = http.StatusOK

	// Net.
	net := network.TestNet
	switch chainnet {
	case MainNet:
		net = network.MainNet
	}

	// Check msg.
	{
		if msg != "" {
			if len(msg) > 64 {
				rsp.Code = http.StatusInternalServerError
				rsp.Message = fmt.Sprintf("message.too.long[%v].max[%v]", len(msg), 64)
				return marshal(rsp)
			}
		}
	}

	// Master pravite key.
	{
		masterkey, err = bip32.NewHDKeyFromString(masterPrvKey)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
	}

	// To address.
	{
		to, err = xcore.DecodeAddress(toAddress, net)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
	}

	// Get unspents.
	{
		req := &proto.WalletUnspentRequest{
			Amount: amount + fees,
		}

		path := fmt.Sprintf("%s/api/wallet/unspent", url)
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}

		if err := httpRsp.Json(&unspents); err != nil {
			rsp.Code = httpRsp.StatusCode()
			rsp.Message = err.Error()
			return marshal(rsp)
		}
	}

	// Change address.
	{
		changeAddress := unspents[0].Address
		change, err = xcore.DecodeAddress(changeAddress, net)
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
	}

	// Transaction build.
	{
		// Coins.
		coinBuilder := xcore.NewCoinBuilder()
		for _, unspent := range unspents {
			coinBuilder.AddOutput(
				unspent.Txid,
				unspent.Vout,
				unspent.Value,
				unspent.Scriptpubkey)
		}
		coins := coinBuilder.ToCoins()

		// Transaction builder.
		txBuilder := xcore.NewTransactionBuilder()
		for _, coin := range coins {
			txBuilder.AddCoin(coin).Then()
		}
		txBuilder.To(to, amount)
		txBuilder.SetChange(change).SendFees(fees)
		if msg != "" {
			txBuilder.AddPushData([]byte(msg))
		}
		tx, err := txBuilder.BuildTransaction()
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}

		for i, unspent := range unspents {
			var sighash []byte
			if strings.HasPrefix(unspent.Address, net.Bech32HRPSegwit) {
				_, version, _, err := xbase.WitnessDecode(unspent.Address)
				if err != nil {
					rsp.Code = http.StatusInternalServerError
					rsp.Message = err.Error()
					return marshal(rsp)
				}

				switch version {
				case 0x00:
					sighash = tx.WitnessV0SignatureHash(i, xcore.SigHashAll)
				default:
					rsp.Code = http.StatusInternalServerError
					rsp.Message = fmt.Sprintf("bench32.address[%v].version.unknow", unspent.Address)
					return marshal(rsp)
				}
			} else {
				sighash = tx.RawSignatureHash(i, xcore.SigHashAll)
			}

			cliPrvKey, err := masterkey.Derive(unspent.Pos)
			if err != nil {
				rsp.Code = http.StatusInternalServerError
				rsp.Message = err.Error()
				return marshal(rsp)
			}
			svrPubKey, err := bip32.NewHDKeyFromString(unspent.SvrPubKey)
			if err != nil {
				rsp.Code = http.StatusInternalServerError
				rsp.Message = err.Error()
				return marshal(rsp)
			}

			// Signature.
			scriptHex, err := hex.DecodeString(unspent.Scriptpubkey)
			if err != nil {
				rsp.Code = http.StatusInternalServerError
				rsp.Message = err.Error()
				return marshal(rsp)
			}

			// Signature version.
			script, err := xcore.ParseLockingScript(scriptHex)
			if err != nil {
				rsp.Code = http.StatusInternalServerError
				rsp.Message = err.Error()
				return marshal(rsp)
			}
			scriptVersion := script.GetScriptVersion()
			switch scriptVersion {
			case xcore.BASE, xcore.WITNESS_V0:
				if err := signECDSA(url, token, unspent.Pos, i, sighash, tx, cliPrvKey, svrPubKey); err != nil {
					rsp.Code = http.StatusInternalServerError
					rsp.Message = err.Error()
					return marshal(rsp)
				}
			default:
				rsp.Code = http.StatusInternalServerError
				rsp.Message = fmt.Sprintf("script.version[%v].unsupport", scriptVersion)
				return marshal(rsp)
			}
		}

		// Verify Tx.
		if err := tx.Verify(); err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}
		localtxid := tx.ID()

		// Push tx.
		{
			path := fmt.Sprintf("%s/api/wallet/pushtx", url)

			req := &proto.TxPushRequest{
				TxHex: fmt.Sprintf("%x", tx.Serialize()),
			}
			httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
			if err != nil {
				rsp.Code = http.StatusInternalServerError
				rsp.Message = err.Error()
				return marshal(rsp)
			}

			pushrsp := &proto.TxPushResponse{}
			if err := httpRsp.Json(pushrsp); err != nil {
				rsp.Code = httpRsp.StatusCode()
				rsp.Message = err.Error()
				return marshal(rsp)
			}
			rsp.TxID = pushrsp.TxID
			if localtxid != pushrsp.TxID {
				rsp.Code = http.StatusInternalServerError
				rsp.Message = fmt.Sprintf("library.send.to.address[%v].push.tx.txid[local:%v, remote:%v].error", toAddress, localtxid, pushrsp.TxID)
				return marshal(rsp)
			}
		}
	}
	return marshal(rsp)
}

func signECDSA(url string, token string, pos uint32, txInIdx int, sighash []byte, tx *xcore.Transaction, cliPrvKey *bip32.HDKey, svrPubKey *bip32.HDKey) error {
	var shareR1 *secp256k1.Scalar

	aliceParty := xcrypto.NewEcdsaParty(cliPrvKey.PrivateKey())
	// Phase1.
	sharepub := aliceParty.Phase1(svrPubKey.PublicKey())
	// Phase2.
	encpk1, encpub1, scalarR1 := aliceParty.Phase2(sighash)

	// Get R2.
	{
		r2req := &proto.EcdsaR2Request{
			Pos:  pos,
			Hash: sighash,
			R1:   scalarR1,
		}

		path := fmt.Sprintf("%s/api/ecdsa/r2", url)
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, r2req)
		if err != nil {
			return err
		}
		r2rsp := &proto.EcdsaR2Response{}
		if err := httpRsp.Json(&r2rsp); err != nil {
			return err
		}

		// Check two party Share R is same or not.
		shareR1 = aliceParty.Phase3(r2rsp.R2)
		if r2rsp.ShareR.X.Cmp(shareR1.X) != 0 || r2rsp.ShareR.Y.Cmp(shareR1.Y) != 0 {
			return fmt.Errorf("shareR.not.equal")
		}
	}

	// Get S2.
	{
		s2req := &proto.EcdsaS2Request{
			Pos:     pos,
			Hash:    sighash,
			R1:      scalarR1,
			EncPK1:  encpk1,
			EncPub1: encpub1,
			ShareR:  shareR1,
		}

		path := fmt.Sprintf("%s/api/ecdsa/s2", url)
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, s2req)
		if err != nil {
			return err
		}
		s2rsp := &proto.EcdsaS2Response{}
		if err := httpRsp.Json(&s2rsp); err != nil {
			return err
		}

		// Phase5.
		sharesig, err := aliceParty.Phase5(shareR1, s2rsp.S2)
		if err != nil {
			return err
		}

		// Verify.
		if err := xcrypto.EcdsaVerify(sharepub, sighash, sharesig); err != nil {
			return err
		}

		// Embed IdxSignature.
		tx.EmbedIdxEcdsaSignature(txInIdx, sharepub, sharesig, xcore.SigHashAll)
	}
	return nil
}
