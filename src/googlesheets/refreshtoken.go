package googlesheets

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/art-es/ac-export/src/logging"
	"github.com/art-es/ac-export/src/utils"
)

type RefreshAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (gs *GoogleSheets) RefreshAccessToken() {
	payload := map[string]string{
		"refresh_token": gs.RefreshToken,
		"client_id":     gs.ClientID,
		"client_secret": gs.ClientSecret,
		"grant_type":    "refresh_token",
	}

	resp := utils.SendRequestWithQueryBody(
		"POST",
		RefreshAccessTokenUri,
		payload,
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		log.Fatal("authentication status not successful\nstatus code: \n", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("authentication request body reading error\nmessage: %s", err)
	}
	logging.RequestInfo(resp, content)

	var respData RefreshAccessTokenResponse
	err = json.Unmarshal(content, &respData)

	gs.AccessToken = respData.AccessToken
}
