package oauth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/lavinas-science/learn-oauth-go/utils/errors"
)

const (
	headerXPublic   = "X-Public"
	headerXClientId = "X-Client-Id"
	headreXCallerId = "X-Caller-Id"

	paramAccessToken = "access_token"

	UserContentType = "application/json"
	UserBaseURI = "http://127.0.0.1:9090"
	UserURI = "/oauth/access_token"

)

var (
	client = resty.New()
)

type accessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

type oauthInterface interface {
}

type oauthClient struct {
}

func IsPublic(r *http.Request) bool {
	if r == nil {
		return true
	}
	return r.Header.Get(headerXPublic) == "true"
}

func AuthenticateRequest(r *http.Request) *errors.RestErr {
	if r == nil {
		return nil
	}
	at := strings.TrimSpace(r.URL.Query().Get(paramAccessToken))
	if at == "" {
		return nil
	}
	return nil
}

func getAccessToken(at string) (*accessToken, *errors.RestErr) {
	resp, err := client.R().SetHeader("Content-Type",UserContentType).Get(UserBaseURI + UserURI + "/" + at)
	if err != nil {
		return nil, errors.NewInternalServerError("Authentication Service off")
	}
	return nil, nil

}