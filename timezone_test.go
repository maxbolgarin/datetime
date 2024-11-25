package datetime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/maxbolgarin/datetime"
)

func TestNewTimezone(t *testing.T) {
	loc := time.FixedZone("TestZone", 3600) // +01:00
	tz := datetime.NewTimezone(loc)

	if tz.Loc().String() != "UTC+1" {
		t.Errorf("NewTimezone location mismatch: got %s, want UTC+1", tz.Loc().String())
	}
	if offset := tz.Offset(); offset != 3600 {
		t.Errorf("NewTimezone offset mismatch: got %d, want 3600", offset)
	}
	if offsetHours := tz.OffsetHours(); offsetHours != 1 {
		t.Errorf("NewTimezone offsetHours mismatch: got %d, want 1", offsetHours)
	}

	tz = datetime.NewTimezone(nil)
	if tz.Loc().String() != "UTC" {
		t.Errorf("NewTimezone location mismatch: got %s, want UTC", tz.Loc().String())
	}
	if offset := tz.Offset(); offset != 0 {
		t.Errorf("NewTimezone offset mismatch: got %d, want 0", offset)
	}
}

func TestNewTimezoneFromTime(t *testing.T) {
	loc := time.FixedZone("TestZone", -3600) // -01:00
	tm := time.Now().In(loc)
	tz := datetime.NewTimezoneFromTime(tm)

	if tz.Loc().String() != "UTC-1" {
		t.Errorf("Expected location UTC-1, got %s", tz.Loc().String())
	}
	if offset := tz.Offset(); offset != -3600 {
		t.Errorf("Expected offset -3600, got %d", offset)
	}

	loc = time.FixedZone("TestZone", -3600-15*60) // -01:00
	tm = time.Now().In(loc)
	tz = datetime.NewTimezoneFromTime(tm)

	if tz.Loc().String() != "UTC-1:15" {
		t.Errorf("Expected location UTC-1:15, got %s", tz.Loc().String())
	}
	if offset := tz.Offset(); offset != -3600-15*60 {
		t.Errorf("Expected offset -3600-15*60, got %d", offset)
	}
}

func TestParseTimezone(t *testing.T) {
	cases := []struct {
		input     string
		expected  string // expected string representation
		expectErr bool
	}{
		{"UTC+02:00", "UTC+2", false},
		{"Europe/London", "UTC", false},
		{"Europe/Moscow", "UTC+3", false},
		{"Invalid/Zone", "", true},
		{"UTC-15:00", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		tz, err := datetime.ParseTimezone(c.input)
		if (err != nil) != c.expectErr {
			t.Errorf("ParseTimezone(%s) error = %v, wantErr %v", c.input, err, c.expectErr)
			continue
		}
		if !c.expectErr && tz.String() != c.expected {
			t.Errorf("ParseTimezone(%s) = %s, expected %s", c.input, tz.String(), c.expected)
		}
	}
}

func TestTimezoneMarshalJSON(t *testing.T) {
	loc := time.FixedZone("TestZone", 3600)
	tz := datetime.NewTimezone(loc)
	data, err := json.Marshal(tz)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	expected := `"UTC+1"`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func TestTimezoneUnmarshalJSON(t *testing.T) {
	jsonData := `"UTC+2"`
	var tz datetime.Timezone
	err := json.Unmarshal([]byte(jsonData), &tz)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if tz.String() != "UTC+2" {
		t.Errorf("UnmarshalJSON = %s, want UTC+2", tz.String())
	}
}

func TestParseUTCOffset(t *testing.T) {
	utcTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		id     string
		input  string
		result time.Time
		isErr  bool
	}{
		{
			id:    "1",
			input: "",
			isErr: true,
		},
		{
			id:     "2",
			input:  "0",
			result: utcTime,
		},
		{
			id:     "3",
			input:  "UTC+0",
			result: utcTime,
		},
		{
			id:     "4",
			input:  "UTC-0",
			result: utcTime,
		},
		{
			id:     "5",
			input:  "+0 0",
			result: utcTime,
		},
		{
			id:     "6",
			input:  "-0 0",
			result: utcTime,
		},
		{
			id:    "7",
			input: "+",
			isErr: true,
		},
		{
			id:     "8",
			input:  "1",
			result: utcTime.In(time.FixedZone("", getOffset(1, 0, 1))),
		},
		{
			id:     "9",
			input:  "+3",
			result: utcTime.In(time.FixedZone("", getOffset(3, 0, 1))),
		},
		{
			id:     "10",
			input:  "-3",
			result: utcTime.In(time.FixedZone("", getOffset(3, 0, -1))),
		},
		{
			id:     "11",
			input:  "3 30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, 1))),
		},
		{
			id:     "12",
			input:  "+3 30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, 1))),
		},
		{
			id:     "13",
			input:  "-3 30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, -1))),
		},
		{
			id:     "14",
			input:  "-3:30",
			result: utcTime.In(time.FixedZone("", getOffset(3, 30, -1))),
		},
		{
			id:     "15",
			input:  "+14",
			result: utcTime.In(time.FixedZone("", getOffset(14, 0, 1))),
		},
		{
			id:    "16",
			input: "-14",
			isErr: true,
		},
		{
			id:    "17",
			input: "15",
			isErr: true,
		},
		{
			id:    "18",
			input: "3 15",
			isErr: true,
		},
		{
			id:    "19",
			input: "2 30",
			isErr: true,
		},
		{
			id:    "20",
			input: "-4 30",
			isErr: true,
		},
		{
			id:    "21",
			input: "+3 31",
			isErr: true,
		},
		{
			id:    "20",
			input: "-12 45",
			isErr: true,
		},
		{
			id:     "21",
			input:  "+12 45",
			result: utcTime.In(time.FixedZone("", getOffset(12, 45, 1))),
		},
		{
			id:    "22",
			input: "+13 45",
			isErr: true,
		},
		{
			id:    "23",
			input: "+13 45 33",
			isErr: true,
		},
		{
			id:    "24",
			input: "+d",
			isErr: true,
		},
		{
			id:    "25",
			input: "13 a",
			isErr: true,
		},
	}

	for _, test := range testCases {
		tz, err := datetime.ParseUTCOffset(test.input)
		if err != nil {
			if !test.isErr {
				t.Errorf("%s -> unexpected error %s", test.id, err)
			}
			continue
		}

		if !test.result.Equal(utcTime.In(tz)) {
			t.Errorf("%s -> expected %v, got %v", test.id, test.result, utcTime.In(tz))
		}
	}
}

func getOffset(hours, minutes, sign int) int {
	return sign*hours*60*60 + sign*minutes*60
}
