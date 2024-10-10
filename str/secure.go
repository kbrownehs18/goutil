package str

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	httpurl "net/url"
	"strconv"
	"time"

	"github.com/kbrownehs18/goutil/arr"
)

// AuthCodeType auth code type
type AuthCodeType int

const (
	// ENCODE encode str
	ENCODE AuthCodeType = iota
	// DECODE decode str
	DECODE
)

// Md5Sum md5
func Md5Sum(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	return hex.EncodeToString(h.Sum(nil))
}

// CertType Certificate type
type CertType int

const (
	// PKCS1 CertType
	PKCS1 CertType = iota
	// PKCS8 CertType
	PKCS8
)

// RsaEncode rsa encode
func RsaEncode(b, rsaKey []byte, t ...CertType) ([]byte, error) {
	block, _ := pem.Decode(rsaKey)
	if block == nil {
		return b, errors.New("key error")
	}
	certType := PKCS8
	if len(t) > 0 {
		certType = t[0]
	}
	var pub interface{}
	var err error
	switch certType {
	case PKCS1:
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
	}

	if err != nil {
		return b, err
	}
	return rsa.EncryptPKCS1v15(crand.Reader, pub.(*rsa.PublicKey), b)
}

// RsaDecode rsa decode
func RsaDecode(b, rsaKey []byte, t ...CertType) ([]byte, error) {
	block, _ := pem.Decode(rsaKey)
	if block == nil {
		return b, errors.New("key error")
	}
	certType := PKCS8
	if len(t) > 0 {
		certType = t[0]
	}
	var priv *rsa.PrivateKey
	var privTemp interface{}
	var err error
	switch certType {
	case PKCS1:
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		privTemp, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	if err != nil {
		return b, err
	}
	if privTemp != nil {
		priv = privTemp.(*rsa.PrivateKey)
	}

	return rsa.DecryptPKCS1v15(crand.Reader, priv, b)
}

// Base64Encode string encode
func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Base64Decode string decode
func Base64Decode(str string) ([]byte, error) {
	x := len(str) * 3 % 4
	switch {
	case x == 2:
		str += "=="
	case x == 1:
		str += "="
	}
	return base64.StdEncoding.DecodeString(str)
}

// Authcode Discuz Authcode golang version
// params[0] encrypt/decrypt bool true：encrypt false：decrypt, default: false
// params[1] key
// params[2] expires time(second)
// params[3] dynamic key length
func Authcode(text string, params ...interface{}) (str string, err error) {
	l := len(params)

	isEncode := DECODE
	key := "abcdefghijklmnopqrstuvwxyz0123456789"
	expiry := 0
	cKeyLen := 8

	if l > 0 {
		isEncode = params[0].(AuthCodeType)
	}

	if l > 1 {
		key = params[1].(string)
	}

	if l > 2 {
		expiry = params[2].(int)
		if expiry < 0 {
			expiry = 0
		}
	}

	if l > 3 {
		cKeyLen = params[3].(int)
		if cKeyLen < 0 {
			cKeyLen = 0
		}
	}
	if cKeyLen > 32 {
		cKeyLen = 32
	}

	timestamp := time.Now().Unix()

	// md5sum key
	mKey := Md5Sum(key)

	// keyA encrypt
	keyA := Md5Sum(mKey[0:16])
	// keyB validate
	keyB := Md5Sum(mKey[16:])
	// keyC dynamic key
	var keyC string
	if cKeyLen > 0 {
		if isEncode == ENCODE {
			// encrypt generate a key
			keyC = Md5Sum(fmt.Sprint(timestamp))[32-cKeyLen:]
		} else {
			// decrypt get key from header of string
			keyC = text[0:cKeyLen]
		}
	}

	// generate encrypt/decrypt key
	cryptKey := keyA + Md5Sum(keyA+keyC)
	// key length
	keyLen := len(cryptKey)
	if isEncode == ENCODE {
		// The first 10 strings is expires time
		// 10-26 strings is validator strings
		var d int64
		if expiry > 0 {
			d = timestamp + int64(expiry)
		}
		text = fmt.Sprintf("%010d%s%s", d, Md5Sum(text + keyB)[0:16], text)
	} else {
		// get strings except dynamic key
		b, e := Base64Decode(text[cKeyLen:])
		if e != nil {
			return "", e
		}
		text = string(b)
	}

	// text length
	textLen := len(text)
	if textLen <= 0 {
		err = fmt.Errorf("auth [%s] textLen <= 0", text)
		return
	}

	// keys
	box := arr.RangeArray(0, 256)
	//
	rndKey := make([]int, 0, 256)
	cryptKeyB := []byte(cryptKey)
	for i := 0; i < 256; i++ {
		pos := i % keyLen
		rndKey = append(rndKey, int(cryptKeyB[pos]))
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + box[i] + rndKey[i]) % 256
		box[i], box[j] = box[j], box[i]
	}

	textB := []byte(text)
	a := 0
	j = 0
	result := make([]byte, 0, textLen)
	for i := 0; i < textLen; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		box[a], box[j] = box[j], box[a]
		result = append(result, byte(int(textB[i])^(box[(box[a]+box[j])%256])))
	}
	fmt.Println(result)
	if isEncode == ENCODE {
		// trim equal
		return keyC + Base64Encode(result), nil
	}

	// check expire time
	d, e := strconv.ParseInt(string(result[0:10]), 10, 0)
	if e != nil {
		err = fmt.Errorf("expires time error: %s", e.Error())
		return
	}

	if (d == 0 || d-timestamp > 0) && string(result[10:26]) == Md5Sum(string(result[26:]) + keyB)[0:16] {
		return string(result[26:]), nil
	}

	err = fmt.Errorf("authcode text [%s] error", text)
	return
}

// URLEncode urlencode
func URLEncode(params interface{}) string {
	q, ok := params.(string)
	if ok {
		return httpurl.QueryEscape(q)
	}
	m, ok := params.(map[string]string)
	if ok {
		val := httpurl.Values{}
		for k, v := range m {
			val.Set(k, v)
		}

		return val.Encode()
	}

	return ""
}

// URLDecode urldecode
func URLDecode(str string) (string, error) {
	return httpurl.QueryUnescape(str)
}

// SHA256 return string
func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// HMacSHA256 hmac sha256
func HMacSHA256(s, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// MaxEncryptBlock rsa encode max length
var MaxEncryptBlock = 117

// MaxDecryptBlock rsa decode max length
var MaxDecryptBlock = 128

// RSAEncode rsa
func RSAEncode(b, key []byte, t ...CertType) ([]byte, error) {
	l := len(b)
	offset := 0
	var data bytes.Buffer
	var i int
	for l-offset > 0 {
		var cache []byte
		var err error
		if l-offset > MaxEncryptBlock {
			cache, err = RsaEncode(b[offset:offset+MaxEncryptBlock], key, t...)
		} else {
			cache, err = RsaEncode(b[offset:], key, t...)
		}
		if err != nil {
			log.Print("RSA Encode error: ", err)
			return nil, err
		}
		data.Write(cache)
		i++
		offset = i * MaxEncryptBlock
	}

	return data.Bytes(), nil
}

// RSADecode rsa decode
func RSADecode(b, key []byte, t ...CertType) ([]byte, error) {
	l := len(b)
	offset := 0
	var data bytes.Buffer
	var i int
	for l-offset > 0 {
		var cache []byte
		var err error
		if l-offset > MaxDecryptBlock {
			cache, err = RsaDecode(b[offset:offset+MaxDecryptBlock], key, t...)
		} else {
			cache, err = RsaDecode(b[offset:], key, t...)
		}
		if err != nil {
			log.Print("RSA Decode error: ", err)
			return nil, err
		}
		data.Write(cache)
		i++
		offset = i * MaxDecryptBlock
	}

	return data.Bytes(), nil
}
