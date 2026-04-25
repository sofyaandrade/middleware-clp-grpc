package connection

import (
	"bytes"
	"encoding/json"
	"io"
	"middleware/internal/domain/models"
	"middleware/internal/domain/security"
	"net/http"

	"github.com/spf13/viper"
)

var (
	Tokens   models.ResponseLogin
	urlLogin = "/login/"
)

func RequestToken() (*models.ResponseLogin, error) {

	urlMedicoes := retornarUrl()

	isAuthorized, err := security.IsAuthorized(Tokens.AccessToken, viper.GetString("TOKEN_SECRET"))
	if isAuthorized {
		return &Tokens, err
	}

	isAuthorized, _ = security.IsAuthorized(Tokens.RefreshToken, viper.GetString("TOKEN_REFRESH"))
	if isAuthorized {
		return RequestAccessAndRefreshToken(urlMedicoes+"/refresh/", true)
	}
	return RequestAccessAndRefreshToken(urlMedicoes+urlLogin, false)
}

func RequestAccessAndRefreshToken(url string, refresh bool) (*models.ResponseLogin, error) {
	var corpoRequisicao []byte
	if refresh {
		var request models.RefreshTokenRequest
		request.RefreshToken = Tokens.RefreshToken

		corpoRequisicao, _ = json.Marshal(&request)

	} else {
		var login models.Login

		login.Email = viper.GetString("API_LOGIN")
		login.Password = viper.GetString("API_PASSWORD")

		corpoRequisicao, _ = json.Marshal(&login)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(corpoRequisicao))
	if err != nil {
		return &Tokens, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, deveRetornar, err := clienteHttp(req)
	if deveRetornar {
		return &Tokens, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &Tokens)
	if err != nil {
		return &Tokens, err
	}
	return &Tokens, err
}

func RequestRestWithTokenGET(caminho string) ([]byte, error) {

	url := retornarUrl()

	req, err := http.NewRequest("GET", url+caminho, nil)
	if err != nil {
		return nil, err
	}
	token, err := RequestToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, deveRetornar, err := clienteHttp(req)
	if deveRetornar {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func RequestRestGET(caminho string) ([]byte, error) {

	url := retornarUrl()

	req, err := http.NewRequest("GET", url+caminho, nil)
	if err != nil {
		return nil, err
	}

	resp, deveRetornar, err := clienteHttp(req)
	if deveRetornar {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func RequestRestWithTokenPOST(dados interface{}, caminho string) ([]byte, error) {
	url := retornarUrl()

	corpoRequisicao, err := json.Marshal(dados)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url+caminho, bytes.NewBuffer(corpoRequisicao))
	if err != nil {
		return nil, err
	}
	token, err := RequestToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, deveRetornar, err := clienteHttp(req)
	if deveRetornar {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func retornarUrl() string {
	host := viper.GetString("DB_HOST")
	porta := viper.GetString("API_PORT_SMB_MEDICOES")
	url := "http://" + host + ":" + porta
	return url
}

func clienteHttp(req *http.Request) (*http.Response, bool, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, true, err
	}
	return resp, false, nil
}
