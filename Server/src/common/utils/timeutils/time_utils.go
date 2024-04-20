package timeutils

import (
	"fmt"
	"time"
)

// 时段表
var TimeTable []string

func InitTimetable() {
	layout := "15:04"
	beginStr := "00:00"
	t, err := time.ParseInLocation(layout, beginStr, time.Local)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 48; i++ {
		str := t.Format(layout)
		TimeTable = append(TimeTable, str)
		t = t.Add(time.Minute * 30)
	}
	// timeTable = []string{"07:00", "07:30", "08:00", "08:30", "09:00", "09:30", "10:00", "10:30", "11:00", "11:30", "12:00", "12:30", "13:00", "13:30"}
}

func IsLegal(start, end string) (bool, error) {
	layout := "15:04"
	startTime, err := time.ParseInLocation(layout, start, time.Local)
	if err != nil {
		return false, err
	}
	endTime, err := time.ParseInLocation(layout, end, time.Local)
	if err != nil {
		return false, err
	}
	return endTime.After(startTime), nil
}

func IsNowBefore(tStr string) (bool, error) {
	layout := "2006-01-02 15:04"
	today := time.Now().Local().Format("2006-01-02")
	t, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", today, tStr), time.Local)
	if err != nil {
		return false, err
	}
	return time.Now().Local().Before(t), nil
}