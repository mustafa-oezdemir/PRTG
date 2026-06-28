package prtgtime

import (
	"testing"
	"time"
)

func TestParseDateTimeParsesEuropeanPRTGFormat(t *testing.T) {
	SetDefaultTimezone("UTC")

	parsed, unixString, err := ParseDateTime("06.03.2025 15:11:00")
	if err != nil {
		t.Fatalf("ParseDateTime returned error: %v", err)
	}

	expected := time.Date(2025, 3, 6, 15, 11, 0, 0, time.UTC)
	if parsed.Unix() != expected.Unix() {
		t.Fatalf("expected unix %d, got %d", expected.Unix(), parsed.Unix())
	}
	if unixString != "1741273860" {
		t.Fatalf("expected unix string 1741273860, got %q", unixString)
	}
}

func TestParseDateTimeHandlesRangesCrossingMidnight(t *testing.T) {
	SetDefaultTimezone("UTC")

	parsed, unixString, err := ParseDateTime("06.03.2025 23:59:00 - 00:01:00")
	if err != nil {
		t.Fatalf("ParseDateTime returned error: %v", err)
	}

	expected := time.Date(2025, 3, 7, 0, 1, 0, 0, time.UTC)
	if parsed.Unix() != expected.Unix() {
		t.Fatalf("expected range end unix %d, got %d", expected.Unix(), parsed.Unix())
	}
	if unixString != "1741305660" {
		t.Fatalf("expected unix string 1741305660, got %q", unixString)
	}
}

func TestParseDateTimeReturnsErrorForUnsupportedFormat(t *testing.T) {
	SetDefaultTimezone("UTC")

	_, _, err := ParseDateTime("not-a-date")
	if err == nil {
		t.Fatal("expected unsupported datetime to return an error")
	}
}
