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
					// TODO Refactor error types, return json response
					switch statusCode {
					case 400:
						err = ErrStatusIncorrect
					case 401:
						err = ErrStatusAuthorization
					case 403:
						err = ErrStatusBadRequest
					case 404:
						err = ErrStatusNotFound
					case 500:
						err = ErrStatusServer
					case 503:
						err = ErrStatusUnavailable
					default:
						err = ErrUnexpectedHttpStatus
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
