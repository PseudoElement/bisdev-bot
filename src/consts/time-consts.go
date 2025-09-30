package consts

type TimeRangeToMin struct {
	Mins_10  string
	Mins_30  string
	Hour_1   string
	Hours_3  string
	Hours_6  string
	Hours_12 string
	Day_1    string
	Days_3   string
	Week_1   string
	Weeks_2  string
	Month_1  string
	Months_3 string
	AllTime  string
}

var TIME_TO_MIN = TimeRangeToMin{
	Mins_10:  "10",
	Mins_30:  "30",
	Hour_1:   "60",
	Hours_3:  "180",
	Hours_6:  "360",
	Hours_12: "720",
	Day_1:    "1440",
	Days_3:   "4320",
	Week_1:   "10080",
	Weeks_2:  "20160",
	Month_1:  "43200",
	Months_3: "129600",
	// 100 years
	AllTime: "52560000",
}
