package models

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func StrToMD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func PwdHash(str string) string {
	return StrToMD5(str)
}

func StringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

func FormatUnixstamp(stamp int64) (out string) {
	return time.Unix(stamp, 0).Format("2006-01-02 15:04:05")
}

func FormatSize(v interface{}) string {
	var size int64
	if fval, ok := v.(int64); ok {
		size = int64(fval)
	} else if sval, ok := v.(uint64); ok {
		size = int64(sval)
	}

	ord := []string{"K", "M", "G", "T", "P", "E"}
	o := 0
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)

	if size < 973 {
		fmt.Fprintf(w, "%3d ", size)
		w.Flush()
		return buf.String()
	}

	for {
		remain := size & 1023
		size >>= 10

		if size >= 973 {
			o++
			continue
		}

		if size < 9 || (size == 9 && remain < 973) {
			remain = ((remain * 5) + 256) / 512
			if remain >= 10 {
				size++
				remain = 0
			}

			fmt.Fprintf(w, "%d.%d%s", size, remain, ord[o])
			break
		}

		if remain >= 512 {
			size++
		}

		fmt.Fprintf(w, "%3d%s", size, ord[o])
		break
	}

	w.Flush()
	return buf.String()
}
