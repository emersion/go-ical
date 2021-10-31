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
