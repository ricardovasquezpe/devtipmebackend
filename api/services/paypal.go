package services

import (
	"fmt"
	_ "fmt"
	"os"

	//paypalsdk "github.com/netlify/PayPal-Go-SDK"
	"github.com/plutov/paypal"
)

func NewClient() (*paypal.Client, error) {
	c, err := paypal.NewClient("AaLifL9xZQYOxIeqUVxYTrGIm_bWY1m9KWPKaRt_4PptuQLNZm74V9jLC8ZlKFS53wvP-_7VZm8hm1zz", "EIN0wM9LOAMp97DN5epav-6Iy59xH2GoK2YeHsaXhKFQFPzMeGg1ADYPTyqrmEBRWlRyAmlKajjIRMBE", paypal.APIBaseSandBox)

	if err != nil {
		return nil, err
	}

	c.SetLog(os.Stdout)
	_, err = c.GetAccessToken()

	if err != nil {
		return nil, err
	}

	return c, nil
}

func Authorize(client *paypal.Client, orderId string, amount string) error {
	capture, err := client.CaptureOrder(orderId, paypal.CaptureOrderRequest{})
	if err != nil {
		return err
	}

	fmt.Print("CAPTURA")
	fmt.Print(capture.Status)

	return nil
	/*
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

		return nil*/
}
