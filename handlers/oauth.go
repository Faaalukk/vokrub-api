package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

// ── Google ────────────────────────────────────────────────────────────────────

// GET /api/customer/auth/oauth/google
func GoogleInit(c *fiber.Ctx) error {
	params := url.Values{
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"redirect_uri":  {os.Getenv("GOOGLE_REDIRECT_URL")},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"access_type":   {"offline"},
	}
	return c.Redirect("https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode())
}

// GET /api/customer/auth/oauth/google/callback
func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=oauth_denied")
	}

	accessToken, err := exchangeGoogleCode(code)
	if err != nil {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=oauth_failed")
	}

	info, err := fetchGoogleUserInfo(accessToken)
	if err != nil {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=oauth_failed")
	}

	customer, err := findOrCreateOAuthCustomer("google", info["sub"].(string), info["name"].(string), strOrEmpty(info["email"]), strOrEmpty(info["picture"]))
	if err != nil {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=server_error")
	}

	token := makeCustomerJWT(customer)
	return c.Redirect(fmt.Sprintf("%s/auth/callback?token=%s", os.Getenv("APP_URL"), token))
}

func exchangeGoogleCode(code string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")},
		"redirect_uri":  {os.Getenv("GOOGLE_REDIRECT_URL")},
		"grant_type":    {"authorization_code"},
	}
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access_token in google response")
	}
	return token, nil
}

func fetchGoogleUserInfo(accessToken string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return info, nil
}

// ── Facebook ──────────────────────────────────────────────────────────────────

// GET /api/customer/auth/oauth/facebook
func FacebookInit(c *fiber.Ctx) error {
	params := url.Values{
		"client_id":     {os.Getenv("FACEBOOK_APP_ID")},
		"redirect_uri":  {os.Getenv("FACEBOOK_REDIRECT_URL")},
		"response_type": {"code"},
		"scope":         {"email,public_profile"},
	}
	return c.Redirect("https://www.facebook.com/v18.0/dialog/oauth?" + params.Encode())
}

// GET /api/customer/auth/oauth/facebook/callback
func FacebookCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=oauth_denied")
	}

	accessToken, err := exchangeFacebookCode(code)
	if err != nil {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=oauth_failed")
	}

	info, err := fetchFacebookUserInfo(accessToken)
	if err != nil {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=oauth_failed")
	}

	picture := ""
	if pic, ok := info["picture"].(map[string]interface{}); ok {
		if data, ok := pic["data"].(map[string]interface{}); ok {
			picture = strOrEmpty(data["url"])
		}
	}

	customer, err := findOrCreateOAuthCustomer("facebook", info["id"].(string), info["name"].(string), strOrEmpty(info["email"]), picture)
	if err != nil {
		return c.Redirect(os.Getenv("APP_URL") + "/auth?error=server_error")
	}

	token := makeCustomerJWT(customer)
	return c.Redirect(fmt.Sprintf("%s/auth/callback?token=%s", os.Getenv("APP_URL"), token))
}

func exchangeFacebookCode(code string) (string, error) {
	params := url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("FACEBOOK_APP_ID")},
		"client_secret": {os.Getenv("FACEBOOK_APP_SECRET")},
		"redirect_uri":  {os.Getenv("FACEBOOK_REDIRECT_URL")},
	}
	resp, err := http.Get("https://graph.facebook.com/v18.0/oauth/access_token?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access_token in facebook response")
	}
	return token, nil
}

func fetchFacebookUserInfo(accessToken string) (map[string]interface{}, error) {
	endpoint := "https://graph.facebook.com/me?fields=id,name,email,picture&access_token=" + url.QueryEscape(accessToken)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return info, nil
}

// ── Util ──────────────────────────────────────────────────────────────────────

func strOrEmpty(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
