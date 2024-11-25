package datetime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/maxbolgarin/datetime"
)

func TestNewTime(t *testing.T) {
	cases := []struct {
		hour, minute int
		expected     string
	}{
		{10, 30, "10:30"},
		{23, 59, "23:59"},
		{0, 0, "00:00"},
	}

	for _, c := range cases {
		tm := datetime.NewTime(c.hour, c.minute)
		if tm.String() != c.expected {
			t.Errorf("NewTime(%d, %d) = %s; want %s", c.hour, c.minute, tm.String(), c.expected)
		}
	}
}

func TestNewTimeFromString(t *testing.T) {
	cases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"10:30", "10:30", false},
		{"23:59", "23:59", false},
		{"00:00", "00:00", false},
		{"25:00", "", true},
		{"invalid", "", true},
	}

	for _, c := range cases {
		tm, err := datetime.NewTimeFromString(c.input)
		if (err != nil) != c.expectErr || (!c.expectErr && tm.String() != c.expected) {
			t.Errorf("NewTimeFromString(%s) = %v, %v; want %v, %v", c.input, tm, err, c.expected, c.expectErr)
		}
	}
}

func TestNewFromTime(t *testing.T) {
	now := time.Now()
	tm := datetime.NewFromTime(now)
	expected := now.Format("15:04")
	if tm.String() != expected {
		t.Errorf("NewFromTime(%v) = %s; want %s", now, tm.String(), expected)
	}
}

func TestNowTime(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	tm := datetime.NowTime(loc)
	expected := time.Now().In(loc).Format("15:04")
	if tm.String() != expected {
		t.Errorf("NowTime() = %s; want %s", tm.String(), expected)
	}
}

func TestParseTime(t *testing.T) {
	cases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"10:30", "10:30", false},
		{"10-30", "10:30", false},
		{"1030", "10:30", false},
		{"9999", "", true},
		{"abcd", "", true},
		{"", "", true},
		{"23:61", "", true},
		{"25:00", "", true},
		{"23:dd", "", true},
		{"1/1/1", "", true},
	}

	for _, c := range cases {
		tm, err := datetime.ParseTime(c.input)
		if (err != nil) != c.expectErr || (!c.expectErr && tm.String() != c.expected) {
			t.Errorf("ParseTime(%s) = %v, %v; want %v, %v", c.input, tm, err, c.expected, c.expectErr)
		}
	}
}

func TestTimeRange(t *testing.T) {
	low := datetime.NewTime(10, 15)
	high := datetime.NewTime(15, 45)

	diff := low.Range(high)
	expected := time.Hour*5 + time.Minute*30
	if diff != expected {
		t.Errorf("Range() = %v; want %v", diff, expected)
	}

	low = datetime.NewTime(10, 15)
	high = datetime.NewTime(10, 15)

	diff = low.Range(high)
	expected = 0
	if diff != expected {
		t.Errorf("Range() = %v; want %v", diff, expected)
	}

	low = datetime.NewTime(10, 15)
	high = datetime.NewTime(10, 14)

	diff = low.Range(high)
	expected = -time.Minute
	if diff != expected {
		t.Errorf("Range() = %v; want %v", diff, expected)
	}

	low = datetime.NewTime(10, 15)
	high = datetime.NewTime(10, 16)

	diff = low.Range(high)
	expected = time.Minute
	if diff != expected {
		t.Errorf("Range() = %v; want %v", diff, expected)
	}
}

func TestRangeUp(t *testing.T) {
	low := datetime.NewTime(22, 30)
	high := datetime.NewTime(1, 45)

	diff := low.RangeUp(high)
	// 22:30 to 1:45 should be 3 hours and 15 minutes
	expected := time.Hour*3 + time.Minute*15
	if diff != expected {
		t.Errorf("RangeUp() = %v; want %v", diff, expected)
	}

	low = datetime.NewTime(10, 30)
	high = datetime.NewTime(15, 15)

	diff = low.RangeUp(high)
	// 10:30 to 15:15 should be 4 hours and 45 minutes
	expected = time.Hour*4 + time.Minute*45
	if diff != expected {
		t.Errorf("RangeUp() = %v; want %v", diff, expected)
	}

	low = datetime.NewTime(0, 15)
	high = datetime.NewTime(0, 0)

	diff = low.RangeUp(high)
	expected = time.Hour*23 + time.Minute*45
	if diff != expected {
		t.Errorf("RangeUp() = %v; want %v", diff, expected)
	}
}

func TestAddTime(t *testing.T) {
	start := datetime.NewTime(10, 30)
	cases := []struct {
		duration time.Duration
		expected string
	}{
		{time.Hour*3 + time.Minute*45, "14:15"},
		{time.Hour*24 + time.Hour*3 + time.Minute*45, "14:15"},
		{time.Hour*23 + time.Minute*45, "10:15"},
		{time.Hour * 24, "10:30"},
		{time.Hour * 25, "11:30"},
		{time.Minute * 30, "11:00"},
		{time.Minute * 60, "11:30"},
	}

	for _, c := range cases {
		result := start.AddTime(c.duration)
		if result.String() != c.expected {
			t.Errorf("AddTime(%v) = %s; want %s", c.duration, result.String(), c.expected)
		}
	}
}

func TestSubTime(t *testing.T) {
	start := datetime.NewTime(10, 30)
	cases := []struct {
		duration time.Duration
		expected string
	}{
		{time.Hour*24 + time.Hour*2 + time.Minute*45, "07:45"},
		{time.Hour*23 + time.Minute*45, "10:45"},
		{time.Hour * 24, "10:30"},
		{time.Hour * 25, "09:30"},
		{time.Minute * 30, "10:00"},
		{time.Minute * 60, "09:30"},
	}

	for _, c := range cases {
		result := start.SubTime(c.duration)
		if result.String() != c.expected {
			t.Errorf("SubTime(%v) = %s; want %s", c.duration, result.String(), c.expected)
		}
	}
}

func TestMinutesFromDayBegin(t *testing.T) {
	cases := []struct {
		hour, minute int
		expected     int
	}{
		{10, 30, 9*60 + 30},
		{23, 59, 22*60 + 59},
		{0, 0, 23*60},
		{1, 0, 0},
	}

	for _, c := range cases {
		result := datetime.NewTime(c.hour, c.minute).MinutesFromDayBegin(datetime.NewTime(1, 0))
		if result != c.expected {
			t.Errorf("MinutesFromDayBegin(%d, %d) = %d; want %d", c.hour, c.minute, result, c.expected)
		}
	}
}

func TestMinutesTillDayEnd(t *testing.T) {
	cases := []struct {
		hour, minute int
		expected     int
	}{
		{10, 30, 24*60 - 10*60 - 30},
		{23, 59, 1},
		{0, 0, 1440},
	}

	for _, c := range cases {
		result := datetime.NewTime(c.hour, c.minute).MinutesTillDayEnd(datetime.EmptyTime)
		if result != c.expected {
			t.Errorf("MinutesTillDayEnd(%d, %d) = %d; want %d", c.hour, c.minute, result, c.expected)
		}
	}
}

func TestEqualTime(t *testing.T) {
	time1 := datetime.NewTime(8, 15)
	time2 := datetime.NewTime(8, 15)
	time3 := datetime.NewTime(9, 30)

	if !time1.EqualTime(time2) || time1.EqualTime(time3) {
		t.Errorf("EqualTime() failed")
	}
}

func TestComparisonMethods(t *testing.T) {
	earlier := datetime.NewTime(8, 15)
	later := datetime.NewTime(9, 30)

	if !earlier.IsBefore(later) || earlier.IsAfter(later) {
		t.Errorf("IsBefore/IsAfter comparison failed")
	}

	if earlier.IsBeforeStrict(earlier) || later.IsBeforeStrict(earlier) {
		t.Errorf("IsBeforeStrict comparison failed")
	}

	if later.IsAfterStrict(earlier) != true {
		t.Errorf("IsAfterStrict comparison failed")
	}
}

func TestSmartDiff(t *testing.T) {
	start := datetime.NewTime(22, 30)
	end := datetime.NewTime(1, 45)
	expected := time.Hour*3 + time.Minute*15

	if diff := start.SmartDiff(end); diff != expected {
		t.Errorf("SmartDiff() = %v; want %v", diff, expected)
	}
}

func TestRoundDownToFives(t *testing.T) {
	cases := []struct {
		input    datetime.Time
		expected string
	}{
		{datetime.NewTime(10, 2), "10:00"},
		{datetime.NewTime(10, 3), "10:00"},
		{datetime.NewTime(10, 8), "10:05"},
		{datetime.NewTime(10, 15), "10:10"},
	}

	for _, c := range cases {
		if result := c.input.RoundDownToFives(); result.String() != c.expected {
			t.Errorf("RoundDownToFives() = %s; want %s", result.String(), c.expected)
		}
	}
}

func TestRoundUpToFives(t *testing.T) {
	cases := []struct {
		input    datetime.Time
		expected string
	}{
		{datetime.NewTime(10, 2), "10:05"},
		{datetime.NewTime(10, 3), "10:05"},
		{datetime.NewTime(10, 8), "10:10"},
		{datetime.NewTime(10, 15), "10:15"},
	}

	for _, c := range cases {
		if result := c.input.RoundUpToFives(); result.String() != c.expected {
			t.Errorf("RoundUpToFives() = %s; want %s", result.String(), c.expected)
		}
	}
}

func TestIsZero(t *testing.T) {
	if !datetime.EmptyTime.IsZero() {
		t.Error("EmptyTime should be zero")
	}
	if datetime.NewTime(10, 0).IsZero() {
		t.Error("NewTime(10, 0) should not be zero")
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	timeStruct := datetime.NewTime(10, 15)
	expected := `"10:15"`

	result, err := json.Marshal(timeStruct)
	if err != nil || string(result) != expected {
		t.Errorf("MarshalJSON() = %s, %v; want %s", string(result), err, expected)
	}

	empty := datetime.EmptyTime
	expectedEmpty := `null`

	result, err = json.Marshal(empty)
	if err != nil || string(result) != expectedEmpty {
		t.Errorf("MarshalJSON(datetime.EmptyTime) = %s, %v; want %s", string(result), err, expectedEmpty)
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	data := `"10:15"`

	var timeStruct datetime.Time
	err := json.Unmarshal([]byte(data), &timeStruct)
	if err != nil || timeStruct.String() != "10:15" {
		t.Errorf("UnmarshalJSON() = %s, %v; want %s", timeStruct.String(), err, "10:15")
	}

	// Test empty JSON
	timeStruct = datetime.Time{}
	data = `null`
	err = json.Unmarshal([]byte(data), &timeStruct)
	if err != nil || !timeStruct.IsZero() {
		t.Errorf("UnmarshalJSON(null) = %v, %v; want zero value", timeStruct, err)
	}
}

// текущее время (допустим 00 15) меньше времени начала дня (допустим 04 00)
// 1. 00:10 -- меньше текущего, меньше начала -- было недавно -- 2 приоритет
// 2. ??:?? -- меньше текущего, больше начала -- невозможно
// 3. 00:20 -- больше текущего, меньше начала -- скоро будет -- 3 приоритет
// 4. 04:10 -- больше текущего, больше начала -- было давно, в самом верху -- 1 приоритет

// текущее время (допустим 18 00) больше времени начала дня (допустим 04 00)
// 1. 02:00 -- меньше текущего, меньше начала -- сегодня еще будет, в самом низу -- 4 приоритет
// 2. 14:00 -- меньше текущего, больше начала -- было недавно --  2 приоритет
// 3. ??:?? -- больше текущего, меньше начала -- невозможно
// 4. 18:30 -- больше текущего, больше начала -- через 30 минут - 3 приоритет

func TestGetTimeSortingPriority(t *testing.T) {
	var (
		dayStart  = datetime.NewTime(4, 0)
		nowBefore = datetime.NewTime(0, 15)
		nowAfter  = datetime.NewTime(18, 0)
	)

	testCases := []struct {
		id           string
		entered, now datetime.Time
		result       datetime.SortingPriority
	}{
		{
			id:      "b1",
			entered: datetime.NewTime(0, 10),
			now:     nowBefore,
			result:  datetime.BeforePriority,
		},
		{
			id:     "b2",
			now:    nowBefore,
			result: datetime.BeforePriority,
		},
		{
			id:      "b3",
			entered: datetime.NewTime(0, 20),
			now:     nowBefore,
			result:  datetime.AfterPriority,
		},
		{
			id:      "b4",
			entered: datetime.NewTime(4, 10),
			now:     nowBefore,
			result:  datetime.LongAgoPriority,
		},
		{
			id:      "a1",
			entered: datetime.NewTime(2, 0),
			now:     nowAfter,
			result:  datetime.NotSoonPriority,
		},
		{
			id:      "a2",
			entered: datetime.NewTime(14, 0),
			now:     nowAfter,
			result:  datetime.BeforePriority,
		},
		{
			id:     "a3",
			now:    nowAfter,
			result: datetime.NotSoonPriority,
		},
		{
			id:      "a4",
			entered: datetime.NewTime(18, 30),
			now:     nowAfter,
			result:  datetime.AfterPriority,
		},
		{
			id:      "eq",
			entered: dayStart,
			now:     nowAfter,
			result:  datetime.BeforePriority,
		},
		{
			id:      "eq2",
			entered: dayStart,
			now:     dayStart,
			result:  datetime.BeforePriority,
		},
	}

	for _, test := range testCases {
		result := datetime.GetTimeSortingPriority(test.entered, test.now, dayStart)
		if test.result != result {
			t.Errorf("%s -> expected %d, got %d", test.id, test.result, result)
		}
	}
}
