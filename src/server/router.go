// thresh-wallet
//
// Copyright (c) 2019 keyfuse.io
//
// GPLv3 License

package server

import (
	"xlog"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

// APIMux --
type APIMux struct {
	*chi.Mux
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
		r.Post("/api/vcode", handler.vcodefn)
		r.Post("/api/token", handler.tokenfn)
	})

	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(handler.tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/ecdsa/newaddress", handler.ecdsaNewAddress)
		r.Post("/api/wallet/balance", handler.walletBalance)
		//r.Post("/api/ecdsa/r2", handler.ecdsaR2)
		//r.Post("/api/ecdsa/s2", handler.ecdsaS2)
	})
	return APIMux{router}
}
