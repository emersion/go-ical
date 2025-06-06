package ical

import (
	"reflect"
	"testing"
	"time"

	"github.com/teambition/rrule-go"
)

func TestRecurrenceSet(t *testing.T) {
	events := exampleCalendar.Events()
	if len(events) != 1 {
		t.Fatalf("len(Calendar.Events()) = %v, want 1", len(events))
	}
	event := events[0]

	wantRecurrenceSet := &rrule.Set{}
	rrule, err := rrule.NewRRule(rrule.ROption{
		Freq:      rrule.YEARLY,
		Bymonth:   []int{3},
		Byweekday: []rrule.Weekday{rrule.SU.Nth(3)},
	})
	if err != nil {
		t.Errorf("Could not build rrule: %v", err) // Should never really happen.
	}
	wantRecurrenceSet.DTStart(time.Date(1996, 9, 18, 14, 30, 0, 0, time.UTC))
	wantRecurrenceSet.RRule(rrule)

	if gotRecurrenceSet, err := event.RecurrenceSet(nil); err != nil {
		t.Errorf("Props.RecurrenceSet() = %v", err)
	} else if !reflect.DeepEqual(gotRecurrenceSet, wantRecurrenceSet) {
		t.Errorf("Props.RecurrenceSet() = %v, want %v", gotRecurrenceSet, wantRecurrenceSet)
	}
}

func TestRecurrenceSetIsAbsent(t *testing.T) {
	event := Component{}
	gotRecurrenceSet, err := event.RecurrenceSet(nil)

	if gotRecurrenceSet != nil || err != nil {
		t.Errorf("Component.RecurrenceSet() = %v, %v, want nil, nil", gotRecurrenceSet, err)
	}
}

func TestRecurrenceSetWithRDate(t *testing.T) {
	// It creates an event with a daily recurrence rule for 2 days, but also
	// adds a single, separate recurrence date (RDATE).
	event := &Component{
		Name: CompEvent,
		Props: Props{
			PropDateTimeStart: []Prop{{
				Name:  PropDateTimeStart,
				Value: "20230101T100000Z",
			}},
			PropRecurrenceRule: []Prop{{
				Name:  PropRecurrenceRule,
				Value: "FREQ=DAILY;COUNT=2",
			}},
			PropRecurrenceDates: []Prop{{
				Name:  PropRecurrenceDates,
				Value: "20230110T100000Z",
			}},
		},
	}

	// 1. Get the recurrence set from the component
	gotRecurrenceSet, err := event.RecurrenceSet(time.UTC)
	if err != nil {
		t.Fatalf("Component.RecurrenceSet() returned an unexpected error: %v", err)
	}
	if gotRecurrenceSet == nil {
		t.Fatal("Component.RecurrenceSet() returned nil, but a set was expected")
	}

	// 2. Define the expected occurrences
	// The RRULE generates Jan 1 and Jan 2. The RDATE adds Jan 10.
	wantOccurrences := []time.Time{
		time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 10, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 10, 10, 0, 0, 0, time.UTC),
	}

	// 3. Get the actual occurrences from the generated set
	gotOccurrences := gotRecurrenceSet.All()

	// 4. Compare the results
	if !reflect.DeepEqual(gotOccurrences, wantOccurrences) {
		t.Errorf("RecurrenceSet did not process RDATE correctly.\n got: %v\nwant: %v", gotOccurrences, wantOccurrences)
	}
}
