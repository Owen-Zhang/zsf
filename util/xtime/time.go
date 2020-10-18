package xtime

import "time"

const timeFormat = "2006-01-02 15:04:05"

//Time 自定义时间
type Time struct {
	time.Time
}

//Now 返回当前时间
func Now() *Time {
	return &Time{
		time.Now(),
	}
}

//Unix converted from timestamp
func Unix(sec, nsec int64) *Time {
	return &Time{
		time.Unix(sec, nsec),
	}
}

//Today 返回今天的开始时间
func Today() *Time {
	return Now().BeginOfDay()
}

//BeginOfYear 一年的开始时间
func (t *Time) BeginOfYear() *Time {
	y, _, _ := t.Date()
	return &Time{
		time.Date(y, time.January, 1, 0, 0, 0, 0, t.Time.Location()),
	}
}

//EndOfYear 一年的结束时间
func (t *Time) EndOfYear() *Time {
	return &Time{
		t.BeginOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond),
	}
}

//BeginOfMonth 一个月的开始时间
func (t *Time) BeginOfMonth() *Time {
	y, m, _ := t.Date()
	return &Time{
		time.Date(y, m, 1, 0, 0, 0, 0, t.Location()),
	}
}

//EndOfMonth 一个月的结束时间
func (t *Time) EndOfMonth() *Time {
	return &Time{
		t.BeginOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond),
	}
}

//BeginOfWeek 一周的开始时间(星期天为第一天)
func (t *Time) BeginOfWeek() *Time {
	beginDate := t.BeginOfDay()
	return &Time{
		beginDate.Add(-time.Duration(beginDate.Weekday())),
	}
}

//EndOfWeek 一周的结束时间(星期六为结束的一天)
func (t *Time) EndOfWeek() *Time {
	y, m, d := t.BeginOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond).Date()
	return &Time{
		time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location()),
	}
}

//BeginOfDay 一天的开始时间
func (t *Time) BeginOfDay() *Time {
	y, m, d := t.Date()
	return &Time{
		time.Date(y, m, d, 0, 0, 0, 0, t.Location()),
	}
}

//EndOfDay 一天的结时间
func (t *Time) EndOfDay() *Time {
	return &Time{
		t.BeginOfDay().Add(-time.Nanosecond),
	}
}

//UnmarshalJSON 反序列化时间
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	parseTime, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	if err != nil {
		return err
	}
	*t = Time{parseTime}
	return nil
}

//MarshalJSON 序列化time
func (t *Time) MarshalJSON() ([]byte, error) {
	tmp := make([]byte, 0, len(timeFormat)+2)
	tmp = append(tmp, '"')
	tmp = time.Time(t.Time).AppendFormat(tmp, timeFormat)
	tmp = append(tmp, '"')
	return tmp, nil
}

func (t *Time) String() string {
	return t.Format(timeFormat)
}
