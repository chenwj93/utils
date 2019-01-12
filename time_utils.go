package utils

import "time"

var LOCATIONZONE, _ = time.LoadLocation(LOCATION)

//create by cwj on 2017-10-17
// return now time by string
func Now() string {
	return time.Now().In(LOCATIONZONE).Format(TIME_FORMAT_1)
}

func ZeroTime() time.Time {
	a, _ := time.Parse(TIME_FORMAT_3, "0001-01-01")
	return a
}

func TodayWithout() string {
	return time.Now().In(LOCATIONZONE).Format("20060102")
}

func Today() string {
	return time.Now().In(LOCATIONZONE).Format("2006-01-02")
}

func Year() string {
	return time.Now().In(LOCATIONZONE).Format("2006")
}

func GetToday() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, ZERO, ZERO, ZERO, ZERO, LOCATIONZONE)
}

func GetToday24() time.Time {
	return GetToday().Add(24 * time.Hour)
}


type Time time.Time


func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TIME_FORMAT_1+`"`, string(data), LOCATIONZONE)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TIME_FORMAT_1)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TIME_FORMAT_1)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).In(LOCATIONZONE).Format(TIME_FORMAT_1)
}

func NowTime() Time{
	return Time(time.Now())
}
