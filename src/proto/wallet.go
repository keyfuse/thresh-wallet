// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package proto

// WalletCheckRequest --
type WalletCheckRequest struct {
}

// WalletCheckResponse --
type WalletCheckResponse struct {
	UserExists         bool   `json:"user_exists"`
	BackupExists       bool   `json:"backup_exists"`
	BackupTimestamp    int64  `json:"backup_timestamp"`
	BackupCloudService string `json:"backup_cloudservice"`
}

// WalletCreateRequest --
type WalletCreateRequest struct {
	Signature    string `json:"signature"`
	MasterPubKey string `json:"masterpubkey"`
}

// WalletCreateResponse --
type WalletCreateResponse struct {
}

// WalletPortfolioRequest --
type WalletPortfolioRequest struct {
	Code string `json:"code"`
}

// WalletPortfolioResponse --
type WalletPortfolioResponse struct {
	CoinSymbol   string  `json:"coin_symbol"`
	FiatSymbol   string  `json:"fiat_symbol"`
	CurrentPrice float64 `json:"current_price"`
}

// WalletBalanceRequest --
type WalletBalanceRequest struct {
}

// WalletBalanceResponse --
type WalletBalanceResponse struct {
	CoinValue uint64 `json:"coin_value"`
}

// WalletUnspentRequest --
type WalletUnspentRequest struct {
	Amount uint64 `json:"amount"`
}

// WalletUnspentResponse --
type WalletUnspentResponse struct {
	Pos          uint32 `json:"pos"`
	Txid         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	Value        uint64 `json:"value"`
	Address      string `json:"address"`
	Confirmed    bool   `json:"confirmed"`
	SvrPubKey    string `json:"svrpubkey"`
	Scriptpubkey string `json:"scriptpubkey"`
}

// TxPushRequest --
type TxPushRequest struct {
	TxHex string `json:"txhex"`
}

// TxPushResponse --
type TxPushResponse struct {
	TxID string `json:"txid"`
}

// WalletTxsRequest --
type WalletTxsRequest struct {
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	OrderBy string `json:"orderby"`
}

// WalletTxsResponse --
type WalletTxsResponse struct {
	Txid        string `json:"txid"`
	Fee         int64  `json:"fee"`
	Link        string `json:"link"`
	Value       int64  `json:"value"`
	Confirmed   bool   `json:"confirmed"`
	BlockTime   int64  `json:"block_time"`
	BlockHeight int64  `json:"block_height"`
}

// WalletSendFeesRequest --
type WalletSendFeesRequest struct {
	Priority  string `json:"priority"`
	SendValue uint64 `json:"send_value"`
}

// WalletSendFeesResponse --
type WalletSendFeesResponse struct {
	Fees          uint64 `json:"fees"`
	TotalValue    uint64 `json:"total_value"`
	SendableValue uint64 `json:"sendable_value"`
}
