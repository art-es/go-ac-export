package activecollab

import (
	"encoding/json"
	"github.com/art-es/ac-export/src/logging"
	"github.com/art-es/ac-export/src/utils"
	"io/ioutil"
	"log"
)

type AuthenticationPayload struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientName   string `json:"client_name"`
	ClientVendor string `json:"client_vendor"`
}

type AuthenticationBody struct {
	IsOk  bool   `json:"is_ok"`
	Token string `json:"token"`
}

func (ac *ActiveCollab) Authenticate() {
	resp := utils.SendRequestWithJsonBody("POST", AuthenticateUri, &ac.AuthenticationPayload, nil)
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		log.Fatal("authentication status not successful\nstatus code: \n", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("authentication request body reading error\nmessage: %s", err)
	}
	logging.RequestInfo(resp, content)

	var body AuthenticationBody
	err = json.Unmarshal(content, &body)
	if err != nil {
		log.Fatalf("error in decoding authentication request body\nmessage: %s", err)
	}

	if !body.IsOk {
		log.Fatalf("\"is_ok\" field is not true in authentication request body")
	}

	ac.Token = body.Token
}
