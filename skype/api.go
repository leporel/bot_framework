package skypeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leporel/bot_framework/models"
)

func SendReplyMessage(activity *Activity, message, authorizationToken string) error {
	responseActivity := &Activity{
		Type:         activity.Type,
		From:         activity.Recipient,
		Conversation: activity.Conversation,
		Recipient:    activity.From,
		Text:         message,
		ReplyToID:    activity.ID,
	}
	replyUrl := fmt.Sprintf(replyMessageTemplate, activity.ServiceURL, activity.Conversation.ID, activity.ID)
	return SendActivityRequest(responseActivity, replyUrl, authorizationToken)
}

func SendActivityRequest(activity *models.Activity, replyUrl, authorizationToken string) error {
	client := &http.Client{}
	if jsonEncoded, err := json.Marshal(*activity); err != nil {
		return err
	} else {
		req, err := http.NewRequest(
			http.MethodPost,
			replyUrl,
			bytes.NewBuffer(*&jsonEncoded),
		)
		if err == nil {
			req.Header.Set(authorizationHeaderKey, authorizationHeaderValuePrefix+authorizationToken)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(*&req)
			if err == nil {
				defer resp.Body.Close()
				var statusCode int = resp.StatusCode
				if statusCode == http.StatusOK || statusCode == http.StatusCreated ||
					statusCode == http.StatusAccepted || statusCode == http.StatusNoContent {
					return nil
				} else {
					return fmt.Errorf(unexpectedHttpStatusCodeTemplate, statusCode)
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}
}
