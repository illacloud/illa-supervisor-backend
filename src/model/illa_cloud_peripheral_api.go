package model

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	EMAIL_DELIVER_BASEURL                = "https://email.illasoft.com/v1/"
	EMAIL_DELIVER_USAGE_SUBSCRIBE        = "subscribe"
	EMAIL_DELIVER_USAGE_VERIFICATIONCODE = "code"
	EMAIL_DELIVER_INVITE_EMAIL           = "invite"
	EMAIL_DELIVER_SHARE_APP_EMAIL        = "shareApp"
)

func SendSubscriptionEmail(email string) error {
	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{"email": email}).
		Post(EMAIL_DELIVER_BASEURL + EMAIL_DELIVER_USAGE_SUBSCRIBE)
	if resp.StatusCode() != http.StatusOK || err != nil {
		return errors.New("failed to send subscription email")
	}
	fmt.Printf("response: %+v, err: %+v", resp, err)
	return nil
}

func SendVerificationEmail(email, code, usage string) error {
	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{"email": email, "code": code, "usage": usage}).
		Post(EMAIL_DELIVER_BASEURL + EMAIL_DELIVER_USAGE_VERIFICATIONCODE)
	if resp.StatusCode() != http.StatusOK || err != nil {
		return errors.New("failed to send verification code email")
	}
	fmt.Printf("response: %+v, err: %+v", resp, err)
	return nil
}

func SendInviteEmail(m *EmailInviteMessage) error {
	fmt.Printf("[CALL] SendInviteEmail() m: %v\n", m)
	payload := m.Export()
	client := resty.New()
	resp, err := client.R().
		SetBody(payload).
		Post(EMAIL_DELIVER_BASEURL + EMAIL_DELIVER_INVITE_EMAIL)
	fmt.Printf("response: %+v, err: %+v", resp, err)
	if resp.StatusCode() != http.StatusOK || err != nil {
		return errors.New("failed to send invite email")
	}
	return nil
}

func SendShareAppEmail(m *EmailShareAppMessage) error {
	fmt.Printf("%v\n", m)
	payload := m.Export()
	client := resty.New()
	resp, err := client.R().
		SetBody(payload).
		Post(EMAIL_DELIVER_BASEURL + EMAIL_DELIVER_SHARE_APP_EMAIL)
	if resp.StatusCode() != http.StatusOK || err != nil {
		fmt.Printf("response: %+v, err: %+v", resp, err)
		return errors.New("failed to send share app email")
	}
	fmt.Printf("response: %+v, err: %+v", resp, err)
	return nil
}
