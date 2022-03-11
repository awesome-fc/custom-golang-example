package main

import (
	"time"
)

type OssEvent struct {
	Events []struct {
		EventName    string    `json:"eventName"`
		EventSource  string    `json:"eventSource"`
		EventTime    time.Time `json:"eventTime"`
		EventVersion string    `json:"eventVersion"`
		Oss          struct {
			Bucket struct {
				Arn           string `json:"arn"`
				Name          string `json:"name"`
				OwnerIdentity string `json:"ownerIdentity"`
			} `json:"bucket"`
			Object struct {
				DeltaSize int    `json:"deltaSize"`
				ETag      string `json:"eTag"`
				Key       string `json:"key"`
				Size      int    `json:"size"`
			} `json:"object"`
			OssSchemaVersion string `json:"ossSchemaVersion"`
			RuleId           string `json:"ruleId"`
		} `json:"oss"`
		Region            string `json:"region"`
		RequestParameters struct {
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			RequestId string `json:"requestId"`
		} `json:"responseElements"`
		UserIdentity struct {
			PrincipalId string `json:"principalId"`
		} `json:"userIdentity"`
	} `json:"events"`
}