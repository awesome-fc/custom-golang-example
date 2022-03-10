package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	gr "github.com/awesome-fc/golang-runtime"
)

func initialize(ctx *gr.FCContext) error {
	ctx.GetLogger().Infoln("init golang!")
	return nil
}

func handler(ctx *gr.FCContext, event []byte) (bbb []byte, err error) {
	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)

	_, err = json.Marshal(ctx)
	if err != nil {
		fcLogger.Error("error:", err)
	}
	fcLogger.Infof("got oss event!")

	// see event format in: https://help.aliyun.com/document_detail/70140.html?spm=a2c4g.11186623.6.578.5eb8cc74AJCA9p#OSS
	//     context in: https://help.aliyun.com/document_detail/56316.html#using-context
	ossEvent := OssEvent{}
	if err := json.Unmarshal(event, &ossEvent); err != nil {
		return []byte(""), fmt.Errorf("parse event err: %s", err.Error())
	}

	fcLogger.Infof("event is :%+v", ossEvent)
	creds := ctx.Credentials
	event0 := ossEvent.Events[0]
	bucketName := event0.Oss.Bucket.Name
	endpoint := "oss-" + event0.Region + "-internal.aliyuncs.com"

	client, err := oss.New(endpoint, creds.AccessKeyID, creds.AccessKeySecret,
		oss.SecurityToken(creds.SecurityToken))
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return []byte(""), err
	}
	objectName := event0.Oss.Object.Key

	if "ObjectCreated:PutSymlink" == event0.EventName {
		var header http.Header
		if header, err = bucket.GetSymlink(objectName); err != nil {
			return []byte(""), err
		}
		if objectName = header.Get("X-Oss-Symlink-Target"); objectName == "" {
			return []byte(""), fmt.Errorf("invalid symlink %s", event0.Oss.Object.Key)
		}
	}
	if fileType := filepath.Ext(objectName); fileType != ".zip" {
		return []byte(""), fmt.Errorf("not a zip file")
	}
	_, zipName := filepath.Split(objectName)
	newKey := strings.Replace(zipName, ".zip", "/", -1)

	fcLogger.Infof("get obj file: %s", objectName)

	err = bucket.GetObjectToFile(objectName, "/tmp/tmp.zip")
	if err != nil {
		fcLogger.Errorf("get obj file err: %s", err.Error())
		return []byte(""), err
	}
	archive, err := zip.OpenReader("/tmp/tmp.zip")
	if err != nil {
		fcLogger.Errorf("open zip file err: %s", err.Error())
		return []byte(""), err
	}
	for _, f := range archive.File {
		if fd, err := f.Open(); err != nil {
			fcLogger.Errorf("open err %s", err.Error())
			fd.Close()
			return []byte(""), err
		} else {
			if err := bucket.PutObject(newKey + f.Name, fd); err != nil {
				fcLogger.Errorf("put err %s", err.Error())
			}
			fd.Close()
		}
	}
	fcLogger.Infof("put ok")

	return []byte("ok"), nil
}

func main() {
	gr.Start(handler, initialize)
}
