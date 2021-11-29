package main

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/mantil-io/go-mantil-template/api/chat"
)

func TestChat(t *testing.T) {
	api := httpexpect.New(t, apiURL)

	req := chat.DefaultRequest{
		// TODO add attributes
	}
	api.POST("/chat").
		WithJSON(req).
		Expect().
		ContentType("application/json").
		Status(http.StatusOK).
		JSON().Object().
		Value("TODO")

}
