package services

import (
	"encoding/json"
	"os"
	"net/http"
	"errors"
	paypalsdk "github.com/netlify/PayPal-Go-SDK"
)

func NewClient() (*paypalsdk.TokenResponse, error) {
	c, err := paypalsdk.NewClient("AaLifL9xZQYOxIeqUVxYTrGIm_bWY1m9KWPKaRt_4PptuQLNZm74V9jLC8ZlKFS53wvP-_7VZm8hm1zz", "EO1666r1TXzBER4yO7BuWhAxYqtz_R9zL-HF1ejsIX7CVPhjCG3aH11vPuLn5ELUgjAXi2frnVrpTMjC", paypalsdk.APIBaseSandBox)

	if err != nil {
		return nil, err
	}

	c.SetLog(os.Stdout)
	accessToken, err := c.GetAccessToken()

	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func Authorize(accessToken string, orderId string) error {
	urlPaypal := "https://api-m.sandbox.paypal.com"
	url := urlPaypal + "/v2/checkout/orders/" + orderId + "/capture"

	client := &http.Client {}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + accessToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	returnedCluster := map[string]interface{}{}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&returnedCluster)

	if returnedCluster["status"] != "COMPLETED" {
		return errors.New(returnedCluster["message"].(string))
	}

	return nil
}
