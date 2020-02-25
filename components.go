package ical

import (
	"fmt"
	"strings"
	"time"
)

// Calendar is the top-level iCalendar object.
type Calendar struct {
	*Component
}

// NewCalendar creates a new calendar object.
func NewCalendar() *Calendar {
	return &Calendar{NewComponent(CompCalendar)}
}

// Events extracts the list of events contained in the calendar.
func (cal *Calendar) Events() []Event {
	l := make([]Event, 0, len(cal.Children))
	for _, child := range cal.Children {
		if child.Name == CompEvent {
			l = append(l, Event{child})
		}
	}
	return l
}

// Event represents a scheduled amount of time on a calendar.
type Event struct {
	*Component
}

// NewEvent creates a new event.
func NewEvent() *Event {
	return &Event{NewComponent(CompEvent)}
}

// DateTimeStart returns the inclusive start of the event.
func (e *Event) DateTimeStart(loc *time.Location) (time.Time, error) {
	return e.Props.DateTime(PropDateTimeStart, loc)
}

// DateTimeEnd returns the non-inclusive end of the event.
func (e *Event) DateTimeEnd(loc *time.Location) (time.Time, error) {
	if prop := e.Props.Get(PropDateTimeEnd); prop != nil {
		return prop.DateTime(loc)
	}

	startProp := e.Props.Get(PropDateTimeStart)
	if startProp == nil {
		return time.Time{}, nil
	}

	start, err := startProp.DateTime(loc)
	if err != nil {
		return time.Time{}, err
	}

	var dur time.Duration
	if durProp := e.Props.Get(PropDuration); durProp != nil {
		dur, err = durProp.Duration()
		if err != nil {
			return time.Time{}, err
		}
	} else if startProp.ValueType() == ValueDate {
		dur = 24 * time.Hour
	}

	return start.Add(dur), nil
}

func (e *Event) Status() (EventStatus, error) {
	s, err := e.Props.Text(PropStatus)
	if err != nil {
		return "", err
	}

	switch status := EventStatus(strings.ToUpper(s)); status {
	case "", EventTentative, EventConfirmed, EventCancelled:
		return status, nil
	default:
		return "", fmt.Errorf("ical: invalid VEVENT STATUS: %q", status)
	}
}

func (e *Event) SetStatus(status EventStatus) {
	if status == "" {
		e.Props.Del(PropStatus)
	} else {
		e.Props.SetText(PropStatus, string(status))
	}
}
