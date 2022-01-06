package acharset

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io/ioutil"
)

var charsetList = map[string]string{
	"HZGB2312": "HZ-GB-2312",
	"GB2312":   "HZ-GB-2312",
	"hzgb2312": "HZ-GB-2312",
	"gb2312":   "HZ-GB-2312",
}
func Decode(value string, charset string) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(value)), getEncoding(charset).NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// getEncoding returns the encoding.Encoding interface object for `charset`.
func getEncoding(charset string) encoding.Encoding {
	if c, ok := charsetList[charset]; ok {
		charset = c
	}
	enc, err := ianaindex.MIB.Encoding(charset)
	if err != nil {
		panic(err)
	}
	return enc
}
