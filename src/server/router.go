// thresh-wallet
//
// Copyright (c) 2019 KeyFuse Labs
//
// GPLv3 License

package server

import (
	"xlog"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

// APIMux --
type APIMux struct {
	*chi.Mux
	handler *Handler
}

// NewAPIRouter -- create new apiMux.
func NewAPIRouter(log *xlog.Log, conf *Config) APIMux {
	router := chi.NewRouter()
	router.Use(middleware.DefaultCompress)
	router.Use(middleware.DefaultLogger)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})
	router.Use(cors.Handler)

	handler := NewHandler(log, conf)
	router.Group(func(r chi.Router) {
		// Limiter.
		lmt := tollbooth.NewLimiter(0.1, nil)
		lmt.SetMessage("You have reached maximum request limit.")
		r.Use(tollbooth_chi.LimitHandler(lmt))

		r.Get("/api/server/info", handler.serverInfo)
		r.Post("/api/login/vcode", handler.loginVCode)
		r.Post("/api/login/token", handler.loginToken)
	})

	router.Group(func(r chi.Router) {
		// Limiter.
		lmt := tollbooth.NewLimiter(5, nil)
		lmt.SetMessage("You have reached maximum request limit.")
		r.Use(tollbooth_chi.LimitHandler(lmt))

		r.Use(jwtauth.Verifier(handler.tokenAuth))
		r.Use(jwtauth.Authenticator)

		// Wallet.
		r.Post("/api/wallet/txs", handler.walletTxs)
		r.Post("/api/wallet/check", handler.walletCheck)
		r.Post("/api/wallet/create", handler.walletCreate)
		r.Post("/api/wallet/pushtx", handler.walletPushTx)
		r.Post("/api/wallet/balance", handler.walletBalance)
		r.Post("/api/wallet/unspent", handler.walletUnspent)
		r.Post("/api/wallet/sendfees", handler.walletSendFees)
		r.Post("/api/wallet/portfolio", handler.walletPortfolio)
		r.Post("/api/wallet/addresses", handler.walletAddresses)

		// ECDSA.
		r.Post("/api/ecdsa/r2", handler.ecdsaR2)
		r.Post("/api/ecdsa/s2", handler.ecdsaS2)
		r.Post("/api/ecdsa/newaddress", handler.ecdsaNewAddress)

		// Backup.
		r.Post("/api/backup/vcode", handler.backupVCode)
		r.Post("/api/backup/verify", handler.backupVerify)
		r.Post("/api/backup/store", handler.backupStore)
		r.Post("/api/backup/restore", handler.backupRestore)
	})
	return APIMux{router, handler}
}

// Init -- used init the mux.
func (a *APIMux) Init() error {
	return a.handler.Init()
}

// Close -- used to close the mux.
func (a *APIMux) Close() {
	a.handler.Close()
}
