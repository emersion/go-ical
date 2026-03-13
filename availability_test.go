package ical

import (
	"strings"
	"testing"
	"time"
)

// Example from RFC 7953 section 3.1.
var exampleAvailabilityStr = toCRLF(`BEGIN:VCALENDAR
PRODID:-//Example Corp.//CalDAV Client//EN
VERSION:2.0
BEGIN:VAVAILABILITY
DTSTAMP:20111005T133225Z
DTEND:20111009T000000Z
DTSTART:20111002T000000Z
UID:20111005T133225Z-00001@example.com
BEGIN:AVAILABLE
DTSTAMP:20111005T133225Z
DTEND:20111002T180000Z
DTSTART:20111002T090000Z
SUMMARY:Weekday morning hours
UID:20111005T133225Z-00001-A@example.com
END:AVAILABLE
END:VAVAILABILITY
END:VCALENDAR
`)

func TestDecodeAvailability(t *testing.T) {
	dec := NewDecoder(strings.NewReader(exampleAvailabilityStr))
	cal, err := dec.Decode()
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	vas := cal.VAvailabilities()
	if len(vas) != 1 {
		t.Fatalf("want 1 VAVAILABILITY, got %d", len(vas))
	}
	va := vas[0]

	uid, err := va.Props.Text(PropUID)
	if err != nil {
		t.Fatalf("UID error: %v", err)
	}
	if uid != "20111005T133225Z-00001@example.com" {
		t.Errorf("want UID %q, got %q", "20111005T133225Z-00001@example.com", uid)
	}

	// Default BUSYTYPE is BUSY-UNAVAILABLE.
	bt, err := va.BusyType()
	if err != nil {
		t.Fatalf("BusyType() error: %v", err)
	}
	if bt != BusyTypeBusyUnavailable {
		t.Errorf("want BusyType %q, got %q", BusyTypeBusyUnavailable, bt)
	}

	dtstart, err := va.DateTimeStart(time.UTC)
	if err != nil {
		t.Fatalf("DateTimeStart() error: %v", err)
	}
	wantStart := time.Date(2011, 10, 2, 0, 0, 0, 0, time.UTC)
	if !dtstart.Equal(wantStart) {
		t.Errorf("want DTSTART %v, got %v", wantStart, dtstart)
	}

	dtend, err := va.DateTimeEnd(time.UTC)
	if err != nil {
		t.Fatalf("DateTimeEnd() error: %v", err)
	}
	wantEnd := time.Date(2011, 10, 9, 0, 0, 0, 0, time.UTC)
	if !dtend.Equal(wantEnd) {
		t.Errorf("want DTEND %v, got %v", wantEnd, dtend)
	}

	avail := va.Available()
	if len(avail) != 1 {
		t.Fatalf("want 1 AVAILABLE, got %d", len(avail))
	}
	a := avail[0]

	summary, err := a.Props.Text(PropSummary)
	if err != nil {
		t.Fatalf("AVAILABLE SUMMARY error: %v", err)
	}
	if summary != "Weekday morning hours" {
		t.Errorf("want SUMMARY %q, got %q", "Weekday morning hours", summary)
	}

	aStart, err := a.DateTimeStart(time.UTC)
	if err != nil {
		t.Fatalf("AVAILABLE DateTimeStart() error: %v", err)
	}
	wantAStart := time.Date(2011, 10, 2, 9, 0, 0, 0, time.UTC)
	if !aStart.Equal(wantAStart) {
		t.Errorf("want AVAILABLE DTSTART %v, got %v", wantAStart, aStart)
	}
}

func TestEncodeAvailability(t *testing.T) {
	cal := NewCalendar()
	cal.Props.SetText(PropProductID, "-//Example Corp.//CalDAV Client//EN")
	cal.Props.SetText(PropVersion, "2.0")

	va := NewVAvailability()
	va.Props.SetDateTime(PropDateTimeStamp, time.Date(2011, 10, 5, 13, 32, 25, 0, time.UTC))
	va.Props.SetDateTime(PropDateTimeStart, time.Date(2011, 10, 2, 0, 0, 0, 0, time.UTC))
	va.Props.SetDateTime(PropDateTimeEnd, time.Date(2011, 10, 9, 0, 0, 0, 0, time.UTC))
	va.Props.SetText(PropUID, "20111005T133225Z-00001@example.com")

	a := NewAvailable()
	a.Props.SetDateTime(PropDateTimeStamp, time.Date(2011, 10, 5, 13, 32, 25, 0, time.UTC))
	a.Props.SetDateTime(PropDateTimeStart, time.Date(2011, 10, 2, 9, 0, 0, 0, time.UTC))
	a.Props.SetDateTime(PropDateTimeEnd, time.Date(2011, 10, 2, 18, 0, 0, 0, time.UTC))
	a.Props.SetText(PropSummary, "Weekday morning hours")
	a.Props.SetText(PropUID, "20111005T133225Z-00001-A@example.com")

	va.Children = append(va.Children, a.Component)
	cal.Children = append(cal.Children, va.Component)

	var buf strings.Builder
	enc := NewEncoder(&buf)
	if err := enc.Encode(cal); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "BEGIN:VAVAILABILITY") {
		t.Errorf("encoded output missing BEGIN:VAVAILABILITY:\n%s", got)
	}
	if !strings.Contains(got, "BEGIN:AVAILABLE") {
		t.Errorf("encoded output missing BEGIN:AVAILABLE:\n%s", got)
	}
}

func TestVAvailabilityBusyType(t *testing.T) {
	va := NewVAvailability()

	// Explicit BUSYTYPE.
	va.SetBusyType(BusyTypeBusy)
	bt, err := va.BusyType()
	if err != nil {
		t.Fatalf("BusyType() error: %v", err)
	}
	if bt != BusyTypeBusy {
		t.Errorf("want %q, got %q", BusyTypeBusy, bt)
	}

	// BUSY-TENTATIVE.
	va.SetBusyType(BusyTypeBusyTentative)
	bt, err = va.BusyType()
	if err != nil {
		t.Fatalf("BusyType() error: %v", err)
	}
	if bt != BusyTypeBusyTentative {
		t.Errorf("want %q, got %q", BusyTypeBusyTentative, bt)
	}
}

func TestEncodeAvailabilityValidation(t *testing.T) {
	t.Run("DTEND and DURATION mutual exclusion in VAVAILABILITY", func(t *testing.T) {
		cal := NewCalendar()
		cal.Props.SetText(PropProductID, "-//Test//EN")
		cal.Props.SetText(PropVersion, "2.0")

		va := NewVAvailability()
		va.Props.SetDateTime(PropDateTimeStamp, time.Now().UTC())
		va.Props.SetText(PropUID, "test-uid@example.com")
		va.Props.SetDateTime(PropDateTimeEnd, time.Now().UTC())
		va.Props.SetText(PropDuration, "PT1H")
		cal.Children = append(cal.Children, va.Component)

		var buf strings.Builder
		enc := NewEncoder(&buf)
		if err := enc.Encode(cal); err == nil {
			t.Error("expected error for DTEND+DURATION in VAVAILABILITY, got nil")
		}
	})

	t.Run("DTEND and DURATION mutual exclusion in AVAILABLE", func(t *testing.T) {
		cal := NewCalendar()
		cal.Props.SetText(PropProductID, "-//Test//EN")
		cal.Props.SetText(PropVersion, "2.0")

		va := NewVAvailability()
		va.Props.SetDateTime(PropDateTimeStamp, time.Now().UTC())
		va.Props.SetText(PropUID, "test-uid@example.com")

		a := NewAvailable()
		a.Props.SetDateTime(PropDateTimeStamp, time.Now().UTC())
		a.Props.SetDateTime(PropDateTimeStart, time.Now().UTC())
		a.Props.SetText(PropUID, "avail-uid@example.com")
		a.Props.SetDateTime(PropDateTimeEnd, time.Now().UTC())
		a.Props.SetText(PropDuration, "PT1H")

		va.Children = append(va.Children, a.Component)
		cal.Children = append(cal.Children, va.Component)

		var buf strings.Builder
		enc := NewEncoder(&buf)
		if err := enc.Encode(cal); err == nil {
			t.Error("expected error for DTEND+DURATION in AVAILABLE, got nil")
		}
	})

	t.Run("forbidden nested component in VAVAILABILITY", func(t *testing.T) {
		cal := NewCalendar()
		cal.Props.SetText(PropProductID, "-//Test//EN")
		cal.Props.SetText(PropVersion, "2.0")

		va := NewVAvailability()
		va.Props.SetDateTime(PropDateTimeStamp, time.Now().UTC())
		va.Props.SetText(PropUID, "test-uid@example.com")
		va.Children = append(va.Children, NewComponent(CompAlarm))
		cal.Children = append(cal.Children, va.Component)

		var buf strings.Builder
		enc := NewEncoder(&buf)
		if err := enc.Encode(cal); err == nil {
			t.Error("expected error for non-AVAILABLE child in VAVAILABILITY, got nil")
		}
	})

	t.Run("forbidden nested component in AVAILABLE", func(t *testing.T) {
		cal := NewCalendar()
		cal.Props.SetText(PropProductID, "-//Test//EN")
		cal.Props.SetText(PropVersion, "2.0")

		va := NewVAvailability()
		va.Props.SetDateTime(PropDateTimeStamp, time.Now().UTC())
		va.Props.SetText(PropUID, "test-uid@example.com")

		a := NewAvailable()
		a.Props.SetDateTime(PropDateTimeStamp, time.Now().UTC())
		a.Props.SetDateTime(PropDateTimeStart, time.Now().UTC())
		a.Props.SetText(PropUID, "avail-uid@example.com")
		a.Props.SetDateTime(PropDateTimeEnd, time.Now().UTC())
		a.Children = append(a.Children, NewComponent(CompAlarm))

		va.Children = append(va.Children, a.Component)
		cal.Children = append(cal.Children, va.Component)

		var buf strings.Builder
		enc := NewEncoder(&buf)
		if err := enc.Encode(cal); err == nil {
			t.Error("expected error for nested component in AVAILABLE, got nil")
		}
	})
}
