package oauth

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lavinas-science/learn-utils-go/rest_errors"
)

const (
	headerXPublic    = "X-Public"
	headerXClientId  = "X-Client-Id"
	headerXCallerId  = "X-Caller-Id"
	paramAccessToken = "access_token"
	UserContentType  = "application/json"
	UserBaseURI      = "http://127.0.0.1:9090"
	UserURI          = "/oauth/access_token"
	timeoutSeconds   = 1
)

var (
	client = resty.New()
)

func init() {
	client.SetTimeout(timeoutSeconds * time.Second)
}

type accessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func IsPublic(r *http.Request) bool {
	if r == nil {
		return true
	}
	return r.Header.Get(headerXPublic) == "true"
}

func GetCallerId(r *http.Request) int64 {
	if r == nil {
		return 0
	}
	c, err := strconv.ParseInt(r.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}
	return c
}

func GetClientId(r *http.Request) int64 {
	if r == nil {
		return 0
	}
	c, err := strconv.ParseInt(r.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return c
}

func AuthenticateRequest(r *http.Request) *rest_errors.RestErr {
	if r == nil {
		return nil
	}
	cleanRequest(r)
	ats := strings.TrimSpace(r.URL.Query().Get(paramAccessToken))
	if ats == "" {
		return nil
	}
	at, err := getAccessToken(ats)
	if err != nil {
		return err
	}
	r.Header.Add(headerXCallerId, strconv.FormatInt(at.UserId, 10))
	r.Header.Add(headerXClientId, strconv.FormatInt(at.ClientId, 10))
	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}
	request.Header.Del(headerXClientId)
	request.Header.Del(headerXCallerId)
}

func getAccessToken(ats string) (*accessToken, *rest_errors.RestErr) {
	re, err := client.R().SetHeader("Content-Type", UserContentType).Get(UserBaseURI + UserURI + "/" + ats)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("Authentication Service off")
	}
	if re.RawResponse.StatusCode > 299 {
		var rErr rest_errors.RestErr
		if err := json.Unmarshal(re.Body(), &rErr); err != nil {
			return nil, rest_errors.NewInternalServerError("Invalid rest-client error unmarshall client")
		}
		return nil, &rErr
	}
	var at accessToken
	if err := json.Unmarshal(re.Body(), &at); err != nil {
		return nil, rest_errors.NewInternalServerError("Invalid rest-client access token unmarshall client")
	}
	return &at, nil
}
