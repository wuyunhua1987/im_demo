package enigma_im

import (
	"backend/http_client"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	appKey    = "imk_fda4d0992833e8a9fd40f90918d814a189fb0f8e3b08abef"
	appSecret = "ims_522a4d3ad1f45152436aad8e2a9d51d27037c20a81662fe0d4771090ce4556e1"
	baseHost  = "https://api.miaowankeji.com/api/v1"
)

func SendChannelMsg(ctx context.Context, userId, roomId, content string) error {
	timestamp := time.Now().Unix()
	signature := calculateSignature(map[string][]string{
		"app_id":    {appKey},
		"timestamp": {fmt.Sprintf("%d", timestamp)},
		"user_id":   {userId},
	}, appSecret)
	m := map[string]interface{}{
		"channel_id": roomId,
		"payload":    content,
	}
	url := fmt.Sprintf("%s/send-channel-message?app_id=%s&timestamp=%d&signature=%s&user_id=%s", baseHost, appKey, timestamp, signature, userId)
	res := struct {
		Error string `json:"error"`
	}{}
	err := http_client.HttpDo(http.MethodPost, map[string]string{}, url, m, &res)
	if err != nil {
		return fmt.Errorf("failed to send channel msg: %w", err)
	}
	if res.Error != "" {
		return fmt.Errorf("failed to send channel msg: %s", res.Error)
	}
	return nil
}

func GetUserToken(ctx context.Context, userId string) string {
	timestamp := time.Now().Unix()
	signature := calculateSignature(map[string][]string{
		"app_id":    {appKey},
		"timestamp": {fmt.Sprintf("%d", timestamp)},
		"user_id":   {userId},
	}, appSecret)
	m := map[string]interface{}{}
	url := fmt.Sprintf("%s/get-user?app_id=%s&timestamp=%d&signature=%s&user_id=%s", baseHost, appKey, timestamp, signature, userId)
	res := struct {
		Error  string `json:"error"`
		UserId string `json:"user_id"`
		Token  string `json:"token"`
	}{}
	err := http_client.HttpDo(http.MethodGet, map[string]string{}, url, m, &res)
	if err != nil {
		log.Printf("Failed to get user token: %v", err)
		return ""
	}
	if res.Error != "" {
		if res.Error == "User not found" {
			return CreateUser(ctx, userId)
		} else {
			log.Printf("Failed to get user token: %v", res.Error)
			return ""
		}
	}
	return res.Token
}

func CreateUser(ctx context.Context, userId string) string {
	timestamp := time.Now().Unix()
	signature := calculateSignature(map[string][]string{
		"app_id":    {appKey},
		"timestamp": {fmt.Sprintf("%d", timestamp)},
	}, appSecret)
	m := map[string]interface{}{
		"user_id": userId,
	}
	url := fmt.Sprintf("%s/create-user?app_id=%s&timestamp=%d&signature=%s", baseHost, appKey, timestamp, signature)
	res := struct {
		Error string `json:"error"`
		Token string `json:"token"`
	}{}
	err := http_client.HttpDo(http.MethodPost, map[string]string{}, url, m, &res)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return ""
	}
	if res.Error != "" {
		log.Printf("Failed to create user: %v", res.Error)
		return ""
	}
	return res.Token
}

func calculateSignature(params map[string][]string, secret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var toSign strings.Builder
	for i, k := range keys {
		if i > 0 {
			toSign.WriteString("&")
		}
		fmt.Fprintf(&toSign, "%s=%s", k, params[k][0])
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(toSign.String()))
	return hex.EncodeToString(h.Sum(nil))
}
