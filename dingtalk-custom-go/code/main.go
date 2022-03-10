package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	gr "github.com/awesome-fc/golang-runtime"
)

func initialize(ctx *gr.FCContext) error {
	ctx.GetLogger().Infoln("init golang!")
	return nil
}

func sign() (timeStamp string, signed string) {
	timeStamp = fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	// todo: fill your secret here
	secret := "SECXXXXXX"
	stringToSign := timeStamp + "\n" + secret

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))
	// Write Data to it
	io.WriteString(h, stringToSign)
	// base64 encode
	signed = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return timeStamp, signed
}

func handler(ctx *gr.FCContext, event []byte) ([]byte, error) {
	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	_, err := json.Marshal(ctx)
	if err != nil {
		fcLogger.Error("error:", err)
	}
	fcLogger.Infof("hello robot!")

	// from here is our robot code
	// see https://open.dingtalk.com/document/robots/custom-robot-access for robot api
	robotUrl, err := ioutil.ReadFile("robot")
	if err != nil {
		return []byte(""), fmt.Errorf("no dingtalk robot url found in file robot")
	}
	robotUrlStr := string(robotUrl)
	if robotUrlStr = strings.Trim(robotUrlStr, " "); robotUrlStr == "" {
		return []byte(""), fmt.Errorf("no dingtalk robot url found in file robot")
	}
	timeStamp, signed := sign()
	robotUrlStr = fmt.Sprintf("%s&timestamp=%s&sign=%s", robotUrlStr, timeStamp, signed)
	jsonObj := map[string]interface{}{}
	jsonObj["msgtype"] = "text"
	jsonObj["text"] = map[string]string{
		"content": "It's time to drink water :-)",
	}

	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return []byte(""), fmt.Errorf("failed to make json request of dingtalk")
	}
	resp, err := http.Post(robotUrlStr, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return []byte(""), err
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	fcLogger.Infof("%+v", msg)

	return []byte("bye bot"), nil
}

func main() {
	gr.Start(handler, initialize)
}
