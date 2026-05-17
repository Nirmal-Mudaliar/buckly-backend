package twilio

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioClient struct {
	Client     *twilio.RestClient
	ServiceSID string
}

func NewTwilioClient(accountSID string, authToken string, serviceSID string) *TwilioClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioClient{
		Client:     client,
		ServiceSID: serviceSID,
	}
}

func (t *TwilioClient) SendOTP(phoneNo string) error {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNo)
	params.SetChannel("sms")

	_, err := t.Client.VerifyV2.CreateVerification(t.ServiceSID, params)
	return err
}

func (t *TwilioClient) VerifyOTP(phoneNo string, code string) (bool, error) {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNo)
	params.SetCode(code)

	resp, err := t.Client.VerifyV2.CreateVerificationCheck(t.ServiceSID, params)
	if err != nil {
		return false, err
	}
	return *resp.Status == "approved", nil
}
