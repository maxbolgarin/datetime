package datetime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/maxbolgarin/datetime"
)

func TestNewDate(t *testing.T) {
	date := datetime.NewDate(2023, 4, 15)
	if date.Year() != 2023 || date.Month() != time.April || date.Day() != 15 {
		t.Error("NewDate did not initialize the date correctly")
	}
}

func TestNewDateFromString(t *testing.T) {
	dateStr := "2023-04-15"
	date, err := datetime.NewDateFromString(dateStr)
	if err != nil || date.String() != dateStr {
		t.Error("NewDateFromString failed to parse the valid date string")
	}

	invalidDateStr := "invalid-date"
	_, err = datetime.NewDateFromString(invalidDateStr)
	if err == nil {
		t.Error("NewDateFromString did not return an error for an invalid date string")
	}
}

func TestNewDateFromTime(t *testing.T) {
	tm := time.Date(2023, time.April, 15, 10, 0, 0, 0, time.UTC)
	date := datetime.NewDateFromTime(tm)
	if !date.EqualDate(datetime.NewDate(2023, 4, 15)) {
		t.Error("NewDateFromTime did not create an equivalent Date")
	}
}

func TestNowDate(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	date := datetime.NowDate(loc)
	expected := datetime.NewDate(now.Year(), int(now.Month()), now.Day())
	if !date.EqualDate(expected) {
		t.Error("NowDate did not return the current date in UTC")
	}
}

func TestToday(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	date := datetime.Today(datetime.EmptyTime, loc)
	now := time.Now().In(loc)
	expected := datetime.NewDate(now.Year(), int(now.Month()), now.Day())
	if !date.EqualDate(expected) {
		t.Error("Today did not return the current date in UTC")
	}
	if !date.IsToday(datetime.EmptyTime, loc) {
		t.Error("IsToday should return false for the current day")
	}
}

func TestParseDate(t *testing.T) {
	validDates := []string{"2023-04-15", "2023.04.15", "2023 04 15", "2023_04_15", "2023-04-15"}
	for _, dateStr := range validDates {
		date, err := datetime.ParseDate(dateStr)
		if err != nil || !date.EqualDate(datetime.NewDate(2023, 4, 15)) {
			t.Errorf("ParseDate failed for valid date string: %s", dateStr)
		}
	}

	_, err := datetime.ParseDate("")
	if err == nil {
		t.Error("ParseDate should fail for empty string")
	}

	_, err = datetime.ParseDate("2023'04'15")
	if err == nil {
		t.Error("ParseDate should fail for invalid separators")
	}

	_, err = datetime.ParseDate("invalid-date")
	if err == nil {
		t.Error("ParseDate should fail for completely invalid date string")
	}

	_, err = datetime.ParseDate("da-04-05")
	if err == nil {
		t.Error("ParseDate should fail for invalid date string")
	}

	_, err = datetime.ParseDate("2024-aa-05")
	if err == nil {
		t.Error("ParseDate should fail for invalid date string")
	}

	_, err = datetime.ParseDate("2024-04-bb")
	if err == nil {
		t.Error("ParseDate should fail for invalid date string")
	}
}

func TestSortDates(t *testing.T) {
	dates := []datetime.Date{
		datetime.NewDate(2023, 4, 15),
		datetime.NewDate(2022, 4, 15),
		datetime.NewDate(2024, 4, 15),
	}
	datetime.SortDates(dates, false)
	if !(dates[0].EqualDate(datetime.NewDate(2022, 4, 15)) && dates[2].EqualDate(datetime.NewDate(2024, 4, 15))) {
		t.Error("SortDates failed to sort in ascending order")
	}

	datetime.SortDates(dates, true)
	if !(dates[0].EqualDate(datetime.NewDate(2024, 4, 15)) && dates[2].EqualDate(datetime.NewDate(2022, 4, 15))) {
		t.Error("SortDates failed to sort in descending order")
	}
}

func TestDateMethods(t *testing.T) {
	date := datetime.NewDate(2023, 4, 15)
	if date.String() != "2023-04-15" {
		t.Error("String method did not return expected format")
	}

	if !date.NextDay().EqualDate(datetime.NewDate(2023, 4, 16)) {
		t.Error("NextDay method did not return the correct next day")
	}

	if !date.PrevDay().EqualDate(datetime.NewDate(2023, 4, 14)) {
		t.Error("PrevDay method did not return the correct previous day")
	}

	if date.IsZero() {
		t.Error("IsZero should return false for a valid date")
	}

	zeroDate := datetime.Date{}
	if !zeroDate.IsZero() {
		t.Error("IsZero should return true for a zero date")
	}
}

func TestEqualDate(t *testing.T) {
	date1 := datetime.NewDate(2023, 4, 15)
	date2 := datetime.NewDate(2023, 4, 15)
	if !date1.EqualDate(date2) {
		t.Error("EqualDate should return true for equivalent dates")
	}

	date3 := datetime.NewDate(2023, 4, 16)
	if date1.EqualDate(date3) {
		t.Error("EqualDate should return false for different dates")
	}
}

func TestMarshalJSON(t *testing.T) {
	date := datetime.NewDate(2023, 4, 15)
	jsonData, err := json.Marshal(date)
	if err != nil || string(jsonData) != "\"2023-04-15\"" {
		t.Error("MarshalJSON failed")
	}

	var newDate datetime.Date
	err = json.Unmarshal(jsonData, &newDate)
	if err != nil || !newDate.EqualDate(date) {
		t.Error("UnmarshalJSON failed")
	}
}

func TestTransformDatesToString(t *testing.T) {
	dates := []datetime.Date{
		datetime.NewDate(2023, 4, 15),
		datetime.NewDate(2022, 4, 15),
		datetime.NewDate(2024, 4, 15),
	}
	expected := []string{"2023-04-15", "2022-04-15", "2024-04-15"}
	result := datetime.TransformDatesToString(dates)
	for i, dateStr := range expected {
		if result[i] != dateStr {
			t.Errorf("TransformDatesToString failed, expected %s but got %s", dateStr, result[i])
		}
	}
}

func TestIsArgNextDay(t *testing.T) {
	date := datetime.NewDate(2023, 4, 15)
	if !date.IsArgNextDay(datetime.NewDate(2023, 4, 16)) {
		t.Error("IsArgNextDay should return true for next day")
	}

	if date.IsArgNextDay(datetime.NewDate(2023, 4, 14)) {
		t.Error("IsArgNextDay should return false for previous day")
	}

	if date.IsArgNextDay(datetime.NewDate(2023, 4, 15)) {
		t.Error("IsArgNextDay should return false for same day")
	}

	if date.IsArgNextDay(datetime.NewDate(2022, 4, 15)) {
		t.Error("IsArgNextDay should return false for different year")
	}

	if date.IsArgNextDay(datetime.NewDate(2023, 3, 15)) {
		t.Error("IsArgNextDay should return false for different month")
	}

	if !date.IsArgNextDay(datetime.NewDate(2025, 4, 15)) {
		t.Error("IsArgNextDay should return true for next year")
	}

	if !date.IsArgNextDay(datetime.NewDate(2023, 7, 15)) {
		t.Error("IsArgNextDay should return true for next month")
	}
}

func TestDateRange(t *testing.T) {
	testCases := []struct {
		id     string
		input1 string
		input2 string
		res    int
	}{
		{
			id:     "1",
			input1: "2020-01-01",
			input2: "2020-01-01",
			res:    0,
		},
		{
			id:     "2",
			input1: "2020-01-01",
			input2: "2020-01-02",
			res:    1,
		},
		{
			id:     "3",
			input1: "2020-01-02",
			input2: "2020-01-01",
			res:    1,
		},
		{
			id:     "4",
			input1: "2020-01-01",
			input2: "2019-01-01",
			res:    365,
		},
		{
			id:     "5",
			input1: "2020-02-01",
			input2: "2020-01-02",
			res:    30,
		},
		{
			id:     "6",
			input1: "2020-04-02",
			input2: "2020-01-01",
			res:    92,
		},
	}

	for _, test := range testCases {
		i1, err := datetime.ParseDate(test.input1)
		if err != nil {
			t.Fatal(err)
		}
		i2, err := datetime.ParseDate(test.input2)
		if err != nil {
			t.Fatal(err)
		}

		result := i1.Range(i2)
		if result != test.res {
			t.Errorf("%s -> expected %v, got %v", test.id, test.res, result)
		}
	}
}
