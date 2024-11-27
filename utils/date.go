package utils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	DateTimeFormat        = "2006-01-02 15:04:05"
	DateFormat            = "2006-01-02"
	DefaultDatetimeFormat = "2006-01-02T15:04:05.000Z"
	DefaultDatetimeTZ     = "2006-01-02T15:04:05Z"

	DD_MM_YYYY = "02/01/2006"
)

var timeZone string

func SetSystemDateTimeZone(tz string) {
	timeZone = tz
}

func getSystemTimeZone() string {

	if timeZone == "" {
		log.Fatal("required load system timeZone config")
	}

	return timeZone
}

// ParseStringToDateTime 字符串转time对象
func ParseStringToDateTime(str string, dateFormat string) (date time.Time) {
	// required LoadLocation
	loc, err := time.LoadLocation(getSystemTimeZone())
	if err != nil {
		log.Fatalf("ParseStringToDateTime Error ", err)
	}

	date, err = time.ParseInLocation(dateFormat, str, loc)
	if err != nil {
		log.Fatalf("ParseStringToDateTime Error ", err)
	}
	return date
}

// ParseStringToTimestamp 字符串转时间戳
func ParseStringToTimestamp(dateStr string, dateFormat string) int64 {
	return ParseStringToDateTime(dateStr, dateFormat).Unix()
}

// ParseTimeToTimestamp time对象转时间戳
func ParseTimeToTimestamp(date *time.Time) (timestamp int64) {
	return date.UnixNano() / 1e6
}

// ParseTimestampToDateTime 时间戳转time对象
func ParseTimestampToDateTime(timestamp int64) *time.Time {
	dateTime := time.Unix(timestamp, 0)
	return &dateTime
}

// ParseUnixMillToDateTime returns the local Time corresponding to the given Unix millisecond time
func ParseUnixMillToDateTime(timestamp int64) *time.Time {
	t := time.Unix(0, timestamp*int64(time.Millisecond))
	return &t
}

// FormatTimestampToString 格式化时间字符串
func FormatTimestampToString(timeUnix int64, dateFormat string) (str string) {
	// required LoadLocation
	loc, _ := time.LoadLocation(getSystemTimeZone())
	return time.Unix(timeUnix, 0).In(loc).Format(dateFormat)
}

// FormatDateTimeToString 格式化时间字符串
func FormatDateTimeToString(dateTime *time.Time, dateFormat string) (str string) {
	// required LoadLocation
	loc, _ := time.LoadLocation(getSystemTimeZone())
	return dateTime.In(loc).Format(dateFormat)
}

// FormatTimestampToDays 格式化时间到天  eg: 2022-06-27 00:00:00
func FormatTimestampToDays(timestamp int64) time.Time {
	return ParseStringToDateTime(FormatTimestampToString(timestamp, DateFormat), DateFormat)
}

// FormatTimestampToDaysEnd 格式化时间到一天的结束 eg: 2022-06-27 23:59:59
func FormatTimestampToDaysEnd(timestamp int64) time.Time {
	return ParseStringToDateTime(fmt.Sprintf(FormatTimestampToString(timestamp, DateFormat)+" %s", "23:59:59"), DateTimeFormat)
}

// AddTime 时间相加
func AddTime(dateTime time.Time, duration time.Duration) time.Time {
	return dateTime.Add(duration)
}

func Duration2Float(elapsed time.Duration) float64 {
	return float64(elapsed.Nanoseconds()) / 1e6
}

// Since returns the time elapsed since t.
// It is shorthand for time.Now().Sub(t).
func Since(begin time.Time) float64 {
	elapsed := time.Since(begin)
	return Duration2Float(elapsed)
}

func GetCurrentTime() *time.Time {
	currentTime := time.Now()
	return &currentTime
}

// ParseTimeStringToTZ 将不同的时间处理成TZ格式后解析
func ParseTimeStringToTZ(str string) (time.Time, error) {
	loc, err := time.LoadLocation(getSystemTimeZone())
	if err != nil {
		return time.Time{}, err
	}
	newStr := parseTimeStringToTZString(str)

	return time.ParseInLocation(DefaultDatetimeFormat, newStr, loc)
}

// ParseTimeStringToTZNoLoc 将不同的时间处理成TZ格式后解析
func ParseTimeStringToTZNoLoc(str string) (time.Time, error) {
	newStr := parseTimeStringToTZString(str)
	return time.Parse(DefaultDatetimeFormat, newStr)
}

func parseTimeStringToTZString(str string) string {
	var (
		newStr, strDate, strTime string
		barTimes                 int
	)

	if len(strings.Split(str, " ")) == 2 {
		strDate = strings.Split(str, " ")[0]
		strTime = strings.Split(str, " ")[1]
	} else if len(strings.Split(str, "T")) == 2 {
		strDate = strings.Split(str, "T")[0]
		strTime = strings.Split(str, "T")[1]
	} else {
		strDate = strings.Split(str, " ")[0]
	}

	// parse date
	for i, s := range []rune(strDate) {
		if string(s) == "-" {
			barTimes += 1
			newStr += string(s)
			continue
		}
		switch barTimes {
		case 1:
			if string([]rune(strDate)[i+1]) == "-" && string([]rune(strDate)[i-1]) == "-" {
				newStr += "0"
			}
		case 2:
			if string([]rune(strDate)[i-1]) == "-" && len([]rune(strDate))-1 == i {
				newStr += "0"
			}
		}
		newStr += string(s)
	}

	if strTime != "" {
		newStr += "T"
	}

	// parse time
	for i, item := range strings.Split(strTime, ":") {
		switch i {
		case 0, 1:
			if i == 0 && item == "" {
				break
			}
			if len([]rune(item)) == 1 {
				newStr += "0"
			}
			newStr += item
			newStr += ":"
		case 2:
			if len(strings.Split(item, ".")) == 2 {
				for ii, ss := range strings.Split(item, ".") {
					switch ii {
					case 0:
						if len([]rune(ss)) == 1 {
							newStr += "0"
						}
						newStr += ss
						newStr += "."
					case 1:
						if len([]rune(ss)) > 4 {
							break
						}
						newStr += ss
					}
				}
			} else {
				if len([]rune(item)) == 1 {
					newStr += "0"
				}
				newStr += item
			}
		}
	}
	baseTime := "2006-01-02T00:00:00.000Z"
	if len([]rune(newStr)) < len([]rune(baseTime)) {
		newStr += baseTime[len([]rune(newStr)):]
	}
	return newStr
}

func ParseTime(ts string) (time.Duration, error) {
	if ts == "" {
		return 0, errors.New("ts is nil")
	}
	l := len(ts)
	total := time.Duration(0)
	j := 0
	for i := 0; i < l; i++ {
		switch ts[i] {
		case 's':
			t, err := strconv.ParseInt(ts[j:i], 10, 64)
			if err != nil {
				return 0, err
			}
			total += time.Second * (time.Duration(t))
			j = i + 1
		case 'm':
			t, err := strconv.ParseInt(ts[j:i], 10, 64)
			if err != nil {
				return 0, err
			}
			total += time.Minute * (time.Duration(t))
			j = i + 1
		case 'h':
			t, err := strconv.ParseInt(ts[j:i], 10, 64)
			if err != nil {
				return 0, err
			}
			total += time.Hour * (time.Duration(t))
			j = i + 1
		case 'd':
			t, err := strconv.ParseInt(ts[j:i], 10, 64)
			if err != nil {
				return 0, err
			}
			total += time.Hour * 24 * (time.Duration(t))
			j = i + 1
		}
	}
	return total, nil
}

// ConvertBuddhistCalendar Convert BuddhistCalendar
func ConvertBuddhistCalendar(t *time.Time) (*time.Time, error) {

	dateTimeStrs := strings.Split(t.Format(DateTimeFormat), "-")
	year, err := strconv.Atoi(dateTimeStrs[0])
	if err != nil {
		return nil, err
	}

	format := strconv.Itoa(year-543) + "-" + dateTimeStrs[1] + "-" + dateTimeStrs[2]
	parse, err := time.Parse(DateTimeFormat, format)
	if err != nil {
		return nil, err
	}

	return &parse, nil
}

func IctTime() (time.Time, error) {
	// 获取泰国时区
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Now(), err
	}

	// 获取当前时间（泰国时区）
	now := time.Now().In(location)
	return now, nil
}
