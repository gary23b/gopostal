package local

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gary23b/gopostal"
	"github.com/stretchr/testify/require"
)

const baseURL string = "http://localhost:8080/"

func Test_getSession(t *testing.T) {
	requestURL := baseURL + "api/sessionCreate"

	// secrets, err := gopostal.ReadSecrets()
	// require.NoError(t, err)

	// type sessionCreateInJson struct {

	// }

	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	require.NoError(t, err)

	res, _, err := gopostal.MakeRequestWithoutRedirects(req, 10*time.Second)
	require.NoError(t, err)

	err = res.SaveToJson("./temp/sessionCreate.json")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.Status)
}

type Session struct {
	CSRFToken string `json:"csrfToken"`
}

func parseBodySession(body []byte) (*Session, error) {
	ret := &Session{}
	err := gopostal.DecodeJson(body, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func Test_login(t *testing.T) {
	requestURL := baseURL + "api/login"

	// secrets, err := gopostal.ReadSecrets()
	// require.NoError(t, err)

	sessionResp, err := gopostal.ReadResponseFromJson("./temp/sessionCreate.json")
	require.NoError(t, err)
	session, err := parseBodySession(sessionResp.Body)
	require.NoError(t, err)

	type inT struct {
		CSRFToken string `json:"csrfToken"`
		UserName  string `json:"userName"`
		Password  string `json:"password"`
	}
	in := &inT{
		CSRFToken: session.CSRFToken,
		UserName:  "test2",
		Password:  "bacon12345",
	}

	body, err := gopostal.EncodeJson(in)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
	require.NoError(t, err)

	req.AddCookie(sessionResp.Cookies["session"])

	res, _, err := gopostal.MakeRequestWithoutRedirects(req, 10*time.Second)
	require.NoError(t, err)

	err = res.SaveToJson("./temp/login.json")
	require.NoError(t, err)

	fmt.Println(res.BodyString)

	require.Equal(t, http.StatusOK, res.Status)
}
