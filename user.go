package dtrack

import (
	"context"
	"net/http"
	"net/url"
)

func (c Client) ForceChangePassword(ctx context.Context, username, password, newPassword string) error {
	body := url.Values{}
	body.Set("username", username)
	body.Set("password", password)
	body.Set("newPassword", newPassword)
	body.Set("confirmPassword", newPassword)

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/user/forceChangePassword", withBody(body))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")

	_, err = c.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
