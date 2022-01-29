package services

import (
	_ "fmt"
	"os"

	//paypalsdk "github.com/netlify/PayPal-Go-SDK"
	"github.com/plutov/paypal"
)

func NewClient() (*paypal.Client, error) {
	c, err := paypal.NewClient(os.Getenv("PAYPAL_CLIENT_ID_PRODUCTION"), os.Getenv("PAYPAL_SECRET_PRODUCTION"), paypal.APIBaseLive)

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

func Authorize(client *paypal.Client, orderId string) error {
	_, err := client.CaptureOrder(orderId, paypal.CaptureOrderRequest{})
	if err != nil {
		return err
	}

	/*fmt.Print("CAPTURA")
	fmt.Print(capture.Status)*/

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
