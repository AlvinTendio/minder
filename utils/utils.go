package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	layoutSlash         = "02/01/2006"
	layoutSlashWithTime = "02/01/2006 15:04:05"
)

func ConvertStringToTime(data string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}, err
	}
	layout := "2006-01-02 15:04:05.0"
	t, err := time.ParseInLocation(layout, data, loc)
	if err != nil {
		t, err = time.ParseInLocation(layoutSlashWithTime, data, loc)

	}
	return t, err
}

func RemoveDetailTime(t time.Time) time.Time {
	res := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return res
}

func ConvertTimeToString(t time.Time) string {
	return t.Format(layoutSlashWithTime)
}

func ConvertTimeToDateString(t time.Time) string {
	layout := "02-01-2006"
	return t.Format(layout)
}

func TrimPhoneNumberZeroPrefix(phone string) (string, error) {
	if len(phone) > 0 && string(phone[0]) == "0" {
		res, err := strconv.Atoi(phone)
		if err != nil {
			log.Println("Failed to parse phone number:", err)
			return "", err
		}

		return "62" + strconv.Itoa(res), err
	}
	return phone, nil
}

func ConvertStringToInt(num string) (int, error) {
	reg := regexp.MustCompile("[^0-9.]")
	processedString := reg.ReplaceAllString(num, "")
	res, err := strconv.Atoi(processedString)
	return res, err
}

func ValidationTransactionDate(date string) bool {
	layout := layoutSlash
	_, err := time.Parse(layout, date)
	return err == nil
}

func ConvertDateTimeGray(date string) string {
	layout := "2006-01-02 15:04:05 +0700 WIB"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Println(err)
	}

	return t.Format(layoutSlashWithTime)
}

func ConvertInt64ToString(data int64) string {
	return strconv.Itoa(int(data))
}

func ConvertStringToDate(date string) time.Time {
	layout := layoutSlash
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Println(err)
		return time.Now()
	}

	return t
}

func ConvertTimeToDateStringSlash(t time.Time) string {
	layout := layoutSlash
	return t.Format(layout)
}

func GenerateRandom() string {
	n := 2
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("%X", b)
}

func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

func CheckStringInSlice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func CheckSliceInSlice(s []string, str []string) bool {
	for _, v := range s {
		for _, x := range str {
			if v == x {
				return true
			}
		}
	}

	return false
}
