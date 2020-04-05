package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/go-shiori/go-readability"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
)

func getArticle(ctx context.Context, bucket *blob.Bucket, url string) (*string, error) {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hashKey := hex.EncodeToString(hasher.Sum(nil))
	ok, err := bucket.Exists(ctx, hashKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		var content string
		article, err := readability.FromURL(url, 1*time.Second)
		if err != nil {
			content = err.Error()
		} else {
			content = article.Content
		}
		err = bucket.WriteAll(ctx, hashKey, []byte(content), nil)
		if err != nil {
			return nil, err
		}
		return &content, nil
	}
	content, err := bucket.ReadAll(ctx, hashKey)
	if err != nil {
		return nil, err
	}
	str := string(content)
	return &str, nil
}
