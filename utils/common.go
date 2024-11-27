package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"github.com/syyongx/php2go"
	"strconv"
	"strings"
)

// 唯一id
func GetGuid() string {
	guid := xid.New()
	return guid.String()
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha256(key, data []byte) string {
	//h := hmac.New(sha256.New, secret)
	//h.Write([]byte(data))
	//sha := hex.EncodeToString(h.Sum(nil))
	//return base64.StdEncoding.EncodeToString([]byte(sha))

	h := hmac.New(sha256.New, key)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func SessionID(projectID int, uid string) string {
	return strconv.Itoa(projectID) + "_" + uid
}

func SessionID2PidUid(sessionID string) (projectID int, UID string, err error) {
	if sessionID == "" {
		err = fmt.Errorf("session id cannot be empty")
		return
	}
	sessionIDSlice := strings.Split(sessionID, "_")
	if len(sessionIDSlice) != 2 {
		err = fmt.Errorf("session id format error")
		return
	}

	projectID, err = strconv.Atoi(sessionIDSlice[0])
	if err != nil {
		err = fmt.Errorf("project id to int err:%s", err.Error())
		return
	}

	UID = sessionIDSlice[1]

	return
}

func Int64ZoomString(num int64, zoom uint8) string {
	strZero := "0000000000"
	str := strconv.FormatInt(num, 10)

	var rt string
	lZoom := int(zoom)
	lStr := len(str)

	diffLen := lZoom - lStr

	if diffLen >= 0 {
		rt = "0." + strZero[:diffLen] + str
	} else {
		rt = str[:lStr-lZoom] + "." + str[lStr-lZoom:]
	}
	return rt
}

func Int64ToDecimalString(num int64, zoom, keep uint8) string {

	if keep > 10 || zoom > 10 {
		return ""
	}

	strZero := "0000000000"
	strNum := strconv.FormatInt(num, 10)

	var sign string
	var str string
	if num < 0 {
		sign = "-"
		strNum := strconv.FormatInt(num, 10)
		str = strNum[1:]
	} else {
		sign = ""
		str = strNum
	}

	var rt string
	lZoom := int(zoom)
	lStr := len(str)

	diffLen := lZoom - lStr
	pos := 0
	if diffLen >= 0 {
		rt = "0." + strZero[:diffLen] + str
		pos = 2
	} else {
		rt = str[:lStr-lZoom] + "." + str[lStr-lZoom:]
		pos = lStr - lZoom + 1
	}

	if keep < zoom {
		v := rt[pos+int(keep) : pos+int(keep)+1]
		if v > "5" {
			tmp, _ := strconv.ParseInt(strings.Replace(rt[:pos+int(keep)], ".", "", 1), 10, 64)

			tmpStr := strconv.FormatInt(tmp+1, 10)

			var rtTmp string
			lKeep := int(keep)
			lTmp := len(tmpStr)

			diffLen := lKeep - lTmp

			if diffLen >= 0 {
				rtTmp = "0." + strZero[:diffLen] + tmpStr
			} else {
				rtTmp = tmpStr[:lStr-lZoom] + "." + tmpStr[lStr-lZoom:]
			}

			return sign + rtTmp

		} else {

			return sign + rt[0:len(rt)-int(zoom-keep)]

		}
	}

	return sign + rt
}

func YmdToStr(year, month, day int, join string) string {
	var sM string
	var sD string
	if month < 10 {
		sM = "0" + strconv.Itoa(month)
	} else {
		sM = strconv.Itoa(month)
	}
	if day < 10 {
		sD = "0" + strconv.Itoa(day)
	} else {
		sD = strconv.Itoa(day)
	}

	return strconv.Itoa(year) + join + sM + join + sD
}

func CompareVersion(version1, version2, operator string) (res bool, err error) {
	operators := []string{
		">",
		"<",
		">=",
		"<=",
		"==",
		"!=",
	}
	if !InArray(operator, operators) {
		err = errors.New("比较符错误")
		return
	}

	return php2go.VersionCompare(version1, version2, operator), nil
}
