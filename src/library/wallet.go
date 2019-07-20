// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package library

import (
	"fmt"
	"net/http"
	"strings"

	"proto"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore"
	"github.com/tokublock/tokucore/xcore/bip32"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xcrypto/secp256k1"
)

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

// EcdsaAddressResponse --
type EcdsaAddressResponse struct {
	Status
	Pos     uint32 `json:"pos"`
	Address string `json:"address"`
}

// APIEcdsaNewAddress -- ecdsa new address api.
func APIEcdsaNewAddress(url string, token string) string {
	rsp := &EcdsaAddressResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/ecdsa/newaddress", url)

	req := &proto.EcdsaAddressRequest{}
	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	address := &proto.EcdsaAddressResponse{}
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

func APIWalletSend(url string, token string, chainnet string, masterPrvKeyStr string, toAddress string, amount uint64, fees uint64) string {
	var err error
	var to xcore.Address
	var change xcore.Address
	var shareR1 *secp256k1.Scalar
	var masterPrvKey *bip32.HDKey
	var unspents []proto.WalletUnspentResponse

	rsp := &WalletSendResponse{}
	rsp.Code = http.StatusOK

	// Net.
	net := network.TestNet
	switch chainnet {
	case MainNet:
		net = network.MainNet
	}

	// Master pravite key.
	{
		masterPrvKey, err = bip32.NewHDKeyFromString(masterPrvKeyStr)
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
		tx, err := txBuilder.BuildTransaction()
		if err != nil {
			rsp.Code = http.StatusInternalServerError
			rsp.Message = err.Error()
			return marshal(rsp)
		}

		for i, unspent := range unspents {
			var sighash []byte
			if strings.HasPrefix(unspent.Address, net.Bech32HRPSegwit) {
				sighash = tx.WitnessSignatureHash(i, xcore.SigHashAll)
			} else {
				sighash = tx.RawSignatureHash(i, xcore.SigHashAll)
			}

			cliPrvKey, err := masterPrvKey.Derive(unspent.Pos)
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

			aliceParty := xcrypto.NewEcdsaParty(cliPrvKey.PrivateKey())
			// Phase1.
			sharepub := aliceParty.Phase1(svrPubKey.PublicKey())
			// Phase2.
			encpk1, encpub1, scalarR1 := aliceParty.Phase2(sighash)

			// Get R2.
			{
				r2req := &proto.EcdsaR2Request{
					Pos:  unspent.Pos,
					Hash: sighash,
					R1:   scalarR1,
				}

				path := fmt.Sprintf("%s/api/ecdsa/r2", url)
				httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, r2req)
				if err != nil {
					rsp.Code = http.StatusInternalServerError
					rsp.Message = err.Error()
					return marshal(rsp)
				}
				r2rsp := &proto.EcdsaR2Response{}
				if err := httpRsp.Json(&r2rsp); err != nil {
					rsp.Code = httpRsp.StatusCode()
					rsp.Message = err.Error()
					return marshal(rsp)
				}

				// Check two party Share R is same or not.
				shareR1 = aliceParty.Phase3(r2rsp.R2)
				if r2rsp.ShareR.X.Cmp(shareR1.X) != 0 || r2rsp.ShareR.Y.Cmp(shareR1.Y) != 0 {
					rsp.Code = http.StatusInternalServerError
					rsp.Message = fmt.Sprintf("shareR.not.equal")
					return marshal(rsp)
				}
			}

			// Get S2.
			{
				s2req := &proto.EcdsaS2Request{
					Pos:     unspent.Pos,
					Hash:    sighash,
					R1:      scalarR1,
					EncPK1:  encpk1,
					EncPub1: encpub1,
					ShareR:  shareR1,
				}

				path := fmt.Sprintf("%s/api/ecdsa/s2", url)
				httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, s2req)
				if err != nil {
					rsp.Code = http.StatusInternalServerError
					rsp.Message = err.Error()
					return marshal(rsp)
				}
				s2rsp := &proto.EcdsaS2Response{}
				if err := httpRsp.Json(&s2rsp); err != nil {
					rsp.Code = httpRsp.StatusCode()
					rsp.Message = err.Error()
					return marshal(rsp)
				}

				// Phase5.
				sharesig, err := aliceParty.Phase5(shareR1, s2rsp.S2)
				if err != nil {
					rsp.Code = http.StatusInternalServerError
					rsp.Message = err.Error()
					return marshal(rsp)
				}

				// Verify.
				if err := xcrypto.EcdsaVerify(sharepub, sighash, sharesig); err != nil {
					rsp.Code = http.StatusInternalServerError
					rsp.Message = err.Error()
					return marshal(rsp)
				}

				// Embed IdxSignature.
				tx.EmbedIdxEcdsaSignature(i, sharepub, sharesig, xcore.SigHashAll)
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
