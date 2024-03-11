/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-22 23:35:30
 * @LastEditTime: 2024-03-11 21:34:12
 * @LastEditors: yujiajie
 */
package util

import "time"

// 获取某一周的周一时间
func GetMonday(d time.Time) time.Time {
	offset := time.Monday - d.Weekday()
	if offset > 0 {
		offset = -6
	}
	return GetZeroTime(d).AddDate(0, 0, int(offset))
}

// 获取某一月的最后一天
func GetMonthLastDay(d time.Time) time.Time {
	return GetMonthFirstDay(d).AddDate(0, 1, -1)
}

// 获取某一月的第一天
func GetMonthFirstDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 格式化日期
func FormatTime(d time.Time) string {
	return d.Format(time.DateTime)
}

// 日期转时间戳
func DateToTime(date string) (int64, error) {
	d, err := time.Parse(time.DateTime, date)
	if err != nil {
		return 0, err
	}
	return d.Unix(), nil
}

// 时间戳转日期
func TimeToDate(timestamp int64) string {
	d := time.Unix(timestamp, 0)
	return FormatTime(d)
}
