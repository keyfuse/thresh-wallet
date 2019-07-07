package server

import (
	"net/http/httptest"

	"xlog"
)

const (
	mockToken           = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkaWQiOiIiLCJtcGsiOiJ0cHViRDZOelZia3JZaFo0WDdDbjFxR1E3WFJldW1ONHlGdmdQM21zOGRQdFRpTEQ3d3BQOTVjcW1iYUFrazVXU1phU3JncGd0bVBRaHBOR21reFZSZXpQM1dONDg2eEVkZHNXSFUyMmE2Rjd5SloiLCJ0IjoxNTYyMjEyOTI1LCJ1aWQiOiJib2h1In0.9eXMwIiIFkGvkh8pL8dnSs_oRbAGemvkowmxlwK5zlc"
	mockCliMasterPubKey = "tpubD6NzVbkrYhZ4X7Cn1qGQ7XReumN4yFvgP3ms8dPtTiLD7wpP95cqmbaAkk5WSZaSrgpgtmPQhpNGmkxVRezP3WN486xEddsWHU22a6F7yJZ"
)

func MockConfig() *Config {
	conf := DefaultConfig()
	conf.DisableVCode = true
	conf.DataDir = "/tmp/tss"
	return conf
}

func MockServer() *httptest.Server {
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	conf := MockConfig()
	router := NewAPIRouter(log, conf)
	return httptest.NewServer(router)
}
