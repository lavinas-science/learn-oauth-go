package oauth

import (
	"net/http"
	"testing"
	"os"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/lavinas-science/learn-utils-go/rest_errors"
)

const (
	mock = true
)

func TestMain(m *testing.M) {
	if mock {
		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()
	}
	os.Exit(m.Run())
}


func TestConstants(t *testing.T) {
	assert.EqualValues(t, headerXPublic, "X-Public")
	assert.EqualValues(t, headerXClientId, "X-Client-Id")
	assert.EqualValues(t, headerXCallerId, "X-Caller-Id")
	assert.EqualValues(t, paramAccessToken, "access_token")
	assert.EqualValues(t, UserContentType, "application/json")
	assert.EqualValues(t, UserBaseURI, "http://127.0.0.1:9090")
	assert.EqualValues(t, UserURI, "/oauth/access_token")
	assert.EqualValues(t, timeoutSeconds, 1)
}

func TestIsPublicNil(t *testing.T) {
	assert.True(t, IsPublic(nil))
}

func TestIsPublic(t *testing.T) {
	assert.True(t, IsPublic(nil))
}

func TestGetAccessTokenOk(t *testing.T) {
	httpmock.RegisterResponder("GET", UserBaseURI + UserURI + "/" + "52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649",
		func(req *http.Request) (*http.Response, error) {
			at := accessToken{
				AccessToken: "52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649",
				UserId: 10,
				ClientId: 10,
				Expires: 100000,
			}
			resp, err := httpmock.NewJsonResponse(http.StatusOK, at)
			if err != nil {
				return httpmock.NewJsonResponse(http.StatusInternalServerError, "Internal error")
			}
			return resp, nil
		},
	)
	at, rErr := getAccessToken("52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649")
	assert.Nil(t, rErr)
	assert.NotNil(t, at)
}

func TestGetAccessTokenError(t *testing.T) {
	httpmock.RegisterResponder("GET", UserBaseURI + UserURI + "/" + "xxxxx",
		func(req *http.Request) (*http.Response, error) {
			rErr := rest_errors.NewNotFoundError("no access token found with given id")
			resp, err := httpmock.NewJsonResponse(rErr.Status(), rErr)
			if err != nil {
				return httpmock.NewJsonResponse(http.StatusInternalServerError, "Internal error")
			}
			return resp, nil
		},
	)
	at, rErr := getAccessToken("xxxxx")
	assert.Nil(t, at)
	assert.NotNil(t, rErr)
	assert.EqualValues(t, rErr, rest_errors.NewNotFoundError("no access token found with given id"))
}