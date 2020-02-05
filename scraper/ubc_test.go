package scraper

import (
	"testing"
)

func Test_preprocessTimestring(t *testing.T) {
	type args struct {
		ts string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"sample time string test",
			args{"Updated: Jan. 31, 2020 – 11:27 a.m. PST"},
			"Jan. 31, 2020 11:27 am PST"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := preprocessTimestring(tt.args.ts); got != tt.want {
				t.Errorf("preprocessTimestring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"msg string 1",
			args{"WEATHER ADVISORY: Campus is OPEN and day and evening classes are in session as normal. Due to current weather conditions, members of the UBC community are reminded of UBC’s Winter Weather Conditions Protocol. Drive safely and wear appropriate footwear when travelling around campus. For information on transit, visit Translink Alerts & Advisories."},
			"WEATHER ADVISORY: Campus is OPEN and day and evening classes are in session as normal."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processMessage(tt.args.message); got != tt.want {
				t.Errorf("processMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
