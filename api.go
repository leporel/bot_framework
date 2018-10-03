package bfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/leporel/bot_framework/bfmodels"
)

func SendMessage(activity bfmodels.Activity, authorizationToken string) error {
	replyUrl := fmt.Sprintf(SendMessageTemplate, activity.ServiceURL, activity.Conversation.ID)
	return SendActivityRequest(&activity, replyUrl, authorizationToken)
}

func SendReplyMessage(activity *bfmodels.Activity, message, authorizationToken string) error {
	responseActivity := &bfmodels.Activity{
		Type:         activity.Type,
		From:         activity.Recipient,
		Conversation: activity.Conversation,
		Recipient:    activity.From,
		Text:         message,
		ReplyToID:    activity.ID,
	}
	replyUrl := fmt.Sprintf(ReplyMessageTemplate, activity.ServiceURL, activity.Conversation.ID, activity.ID)
	return SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}

func SendActivityRequest(activity *bfmodels.Activity, replyUrl, authorizationToken string) error {
	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	if jsonEncoded, err := json.Marshal(*activity); err != nil {
		return err
	} else {
		req, err := http.NewRequest(
			http.MethodPost,
			replyUrl,
			bytes.NewBuffer(*&jsonEncoded),
		)
		if err == nil {
			req.Header.Set(AuthorizationHeaderKey, AuthorizationHeaderValuePrefix+authorizationToken)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(*&req)
			if err == nil {
				defer resp.Body.Close()
				var statusCode int = resp.StatusCode
				if statusCode == http.StatusOK || statusCode == http.StatusCreated ||
					statusCode == http.StatusAccepted || statusCode == http.StatusNoContent {
					return nil
				} else {
					switch statusCode {
					case 400:
						err = fmt.Errorf("The request was malformed or otherwise incorrect.", resp)
					case 401:
						err = fmt.Errorf("The bot is not authorized to make the request.", resp)
					case 403:
						err = fmt.Errorf("The bot is not allowed to perform the requested operation.", resp)
					case 404:
						err = fmt.Errorf("The requested resource was not found.", resp)
					case 500:
						err = fmt.Errorf("An internal server error occurred.", resp)
					case 503:
						err = fmt.Errorf("The service is unavailable.", resp)
					default:
						err = fmt.Errorf(UnexpectedHttpStatusCodeTemplate, statusCode, resp)
					}
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}
}
