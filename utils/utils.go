package utils

import (
	"os"
	"strconv"
	"strings"
	"net/url"
)

func CurrentPath() (path string) {
	path, _ = os.Getwd()
	return
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func ConvertToRawTumblrAddress(u string) string {
	rawPhotoURL, _ := url.Parse(u)
	rawPhotoURL.Scheme = "http"
	rawPhotoURL.Host = "data.tumblr.com"
	photoFileExt := rawPhotoURL.Path[strings.LastIndex(rawPhotoURL.Path, "."):]
	photoTumblrSuffix := rawPhotoURL.Path[:strings.LastIndex(rawPhotoURL.Path, "_")] + "_raw" + photoFileExt;
	rawPhotoURL.Path = photoTumblrSuffix
	return rawPhotoURL.String()
}

func IntToBool(i int) bool {
	return !(i == 0)
}