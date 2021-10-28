package ical

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/teambition/rrule-go"
)

func toCRLF(s string) string {
	return strings.ReplaceAll(s, "\n", "\r\n")
}

var exampleCalendarStr = toCRLF(`BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VEVENT
CATEGORIES:CONFERENCE
DESCRIPTION;ALTREP="cid:part1.0001@example.org":Networld+Interop Conference and Exhibit\nAtlanta World Congress Center\n Atlanta\, Georgia
DTEND:19960920T220000Z
DTSTAMP:19960704T120000Z
DTSTART:19960918T143000Z
ORGANIZER:mailto:jsmith@example.com
RRULE:FREQ=YEARLY;BYDAY=3SU;BYMONTH=3
STATUS:CONFIRMED
SUMMARY;FOO=bar,"b:az":Networld+Interop Conference
UID:uid1@example.com
END:VEVENT
END:VCALENDAR
`)

var exampleCalendar = &Calendar{&Component{
	Name: "VCALENDAR",
	Props: Props{
		"PRODID": []Prop{{
			Name:   "PRODID",
			Params: Params{},
			Value:  "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN",
		}},
		"VERSION": []Prop{{
			Name:   "VERSION",
			Params: Params{},
			Value:  "2.0",
		}},
	},
	Children: []*Component{
		{
			Name: "VEVENT",
			Props: Props{
				"DTSTAMP": []Prop{{
					Name:   "DTSTAMP",
					Params: Params{},
					Value:  "19960704T120000Z",
				}},
				"UID": []Prop{{
					Name:   "UID",
					Params: Params{},
					Value:  "uid1@example.com",
				}},
				"RRULE": []Prop{{
					Name:   "RRULE",
					Params: Params{},
					Value:  "FREQ=YEARLY;BYDAY=3SU;BYMONTH=3",
				}},
				"ORGANIZER": []Prop{{
					Name:   "ORGANIZER",
					Params: Params{},
					Value:  "mailto:jsmith@example.com",
				}},
				"DTSTART": []Prop{{
					Name:   "DTSTART",
					Params: Params{},
					Value:  "19960918T143000Z",
				}},
				"DTEND": []Prop{{
					Name:   "DTEND",
					Params: Params{},
					Value:  "19960920T220000Z",
				}},
				"STATUS": []Prop{{
					Name:   "STATUS",
					Params: Params{},
					Value:  "CONFIRMED",
				}},
				"CATEGORIES": []Prop{{
					Name:   "CATEGORIES",
					Params: Params{},
					Value:  "CONFERENCE",
				}},
				"SUMMARY": []Prop{{
					Name: "SUMMARY",
					Params: Params{
						"FOO": []string{"bar", "b:az"},
					},
					Value: "Networld+Interop Conference",
				}},
				"DESCRIPTION": []Prop{{
					Name: "DESCRIPTION",
					Params: Params{
						"ALTREP": []string{"cid:part1.0001@example.org"},
					},
					Value: `Networld+Interop Conference and Exhibit\nAtlanta World Congress Center\n Atlanta\, Georgia`,
				}},
			},
		},
	},
}}

func TestCalendar(t *testing.T) {
	events := exampleCalendar.Events()
	if len(events) != 1 {
		t.Fatalf("len(Calendar.Events()) = %v, want 1", len(events))
	}
	event := events[0]

	wantSummary := "Networld+Interop Conference"
	if summary, err := event.Props.Text(PropSummary); err != nil {
		t.Errorf("Event.Props.Text(PropSummary) = %v", err)
	} else if summary != wantSummary {
		t.Errorf("Event.Props.Text(PropSummary) = %v, want %v", summary, wantSummary)
	}

	wantDesc := "Networld+Interop Conference and Exhibit\nAtlanta World Congress Center\n Atlanta, Georgia"
	if desc, err := event.Props.Text(PropDescription); err != nil {
		t.Errorf("Event.Props.Text(PropDescription) = %v", err)
	} else if desc != wantDesc {
		t.Errorf("Event.Props.Text(PropDescription) = %v, want %v", desc, wantDesc)
	}

	wantDTStamp := time.Date(1996, 07, 04, 12, 0, 0, 0, time.UTC)
	if dtStamp, err := event.Props.DateTime(PropDateTimeStamp, nil); err != nil {
		t.Errorf("Event.Props.DateTime(PropDateTimeStamp) = %v", err)
	} else if dtStamp != wantDTStamp {
		t.Errorf("Event.Props.DateTime(PropDateTimeStamp) = %v, want %v", dtStamp, wantDTStamp)
	}

	wantDTStart := time.Date(1996, 9, 18, 14, 30, 0, 0, time.UTC)
	if dtStart, err := event.DateTimeStart(nil); err != nil {
		t.Errorf("Event.DateTimeStart() = %v", err)
	} else if dtStart != wantDTStart {
		t.Errorf("Event.DateTimeStart() = %v, want %v", dtStart, wantDTStart)
	}

	wantDTEnd := time.Date(1996, 9, 20, 22, 0, 0, 0, time.UTC)
	if dtEnd, err := event.DateTimeEnd(nil); err != nil {
		t.Errorf("Event.DateTimeEnd() = %v", err)
	} else if dtEnd != wantDTEnd {
		t.Errorf("Event.DateTimeEnd() = %v, want %v", dtEnd, wantDTEnd)
	}
}

func TestGetDate(t *testing.T) {
	localTimezone, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Alias        string
		Value        string
		ValueType    ValueType
		TZID         string
		Location     *time.Location
		ExpectedDate time.Time
	}{
		{
			Alias:        "datetime-local-nil",
			Value:        "20200923T195100",
			ValueType:    ValueDateTime,
			TZID:         "Europe/Paris",
			Location:     nil,
			ExpectedDate: time.Date(2020, time.September, 23, 19, 51, 0, 0, localTimezone),
		},
		{
			Alias:        "datetime-nil-local",
			Value:        "20200923T195100",
			ValueType:    ValueDateTime,
			TZID:         "",
			Location:     localTimezone,
			ExpectedDate: time.Date(2020, time.September, 23, 19, 51, 0, 0, localTimezone),
		},
		{
			Alias:        "datetime-nil-nil",
			Value:        "20200923T195100",
			ValueType:    ValueDateTime,
			TZID:         "",
			Location:     nil,
			ExpectedDate: time.Date(2020, time.September, 23, 19, 51, 0, 0, time.UTC),
		},
		{
			Alias:        "datetime-Z-no-location",
			Value:        "20200923T195100Z",
			ValueType:    ValueDateTime,
			TZID:         "Europe/Paris",
			Location:     nil,
			ExpectedDate: time.Date(2020, time.September, 23, 19, 51, 0, 0, time.UTC),
		},
		{
			Alias:        "datetime-Z-no-tzid",
			Value:        "20200923T195100Z",
			ValueType:    ValueDateTime,
			TZID:         "",
			Location:     localTimezone,
			ExpectedDate: time.Date(2020, time.September, 23, 19, 51, 0, 0, time.UTC),
		},
		{
			Alias:        "date-local-nil",
			Value:        "20200923",
			ValueType:    ValueDate,
			TZID:         "Europe/Paris",
			Location:     nil,
			ExpectedDate: time.Date(2020, time.September, 23, 0, 0, 0, 0, localTimezone),
		},
		{
			Alias:        "date-nil-local",
			Value:        "20200923",
			ValueType:    ValueDate,
			TZID:         "",
			Location:     localTimezone,
			ExpectedDate: time.Date(2020, time.September, 23, 0, 0, 0, 0, localTimezone),
		},
		{
			Alias:        "date-nil-nil",
			Value:        "20200923",
			ValueType:    ValueDate,
			TZID:         "",
			Location:     nil,
			ExpectedDate: time.Date(2020, time.September, 23, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.Alias, func(t *testing.T) {
			p := NewProp("FakeProp")
			p.Value = tCase.Value
			p.SetValueType(tCase.ValueType)
			if tCase.TZID != "" {
				p.Params.Set(PropTimezoneID, tCase.TZID)
			}
			value, err := p.DateTime(tCase.Location)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := value, tCase.ExpectedDate; value.String() != tCase.ExpectedDate.String() {
				t.Errorf("bad date: %s, expected: %s", got, want)
			}
		})
	}
}

func TestSetDate(t *testing.T) {
	localTimezone, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Alias        string
		Date         time.Time
		ExpectedTZID string
		ExpectedDate string
	}{
		{
			Alias:        "UTC",
			Date:         time.Date(2020, time.September, 20, 15, 7, 0, 0, time.UTC),
			ExpectedTZID: "",
			ExpectedDate: "20200920T150700Z",
		},
		{
			Alias:        "local",
			Date:         time.Date(2020, time.September, 20, 17, 7, 0, 0, localTimezone),
			ExpectedTZID: "Europe/Paris",
			ExpectedDate: "20200920T170700",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.Alias, func(t *testing.T) {
			p := NewProp("FakeProp")
			p.SetDateTime(tCase.Date)
			if got, want := p.Params.Get(PropTimezoneID), tCase.ExpectedTZID; got != want {
				t.Errorf("bad tzid: %s, expected: %s", got, want)
			}
			if got, want := p.Value, tCase.ExpectedDate; got != want {
				t.Errorf("bad date: %s, expected: %s", got, want)
			}
		})
	}
}

func TestRoundtripURI(t *testing.T) {
	testCases := []struct {
		Alias    string
		Expected string
	}{
		{
			Alias:    "empty_url",
			Expected: "",
		},
		{
			Alias:    "scheme_and_port",
			Expected: "https://google.com:8080",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.Alias, func(t *testing.T) {
			ue, err := url.Parse(tCase.Expected)
			if err != nil {
				t.Fatalf("%#v", err)
			}
			probs := make(Props)
			probs.SetURI("asdf", ue)
			ug, err := probs.URI("asdf")
			if err != nil {
				t.Errorf("%#v", err)
			}
			if got, want := ug.String(), ue.String(); got != want {
				t.Errorf("bad url: %s, expected: %s", got, want)
			}
		})
	}
}

func TestRecurrenceRule(t *testing.T) {
	events := exampleCalendar.Events()
	if len(events) != 1 {
		t.Fatalf("len(Calendar.Events()) = %v, want 1", len(events))
	}
	props := events[0].Props

	wantRecurrenceRule := &rrule.ROption{
		Freq:      rrule.YEARLY,
		Bymonth:   []int{3},
		Byweekday: []rrule.Weekday{rrule.SU.Nth(3)},
	}
	if roption, err := props.RecurrenceRule(); err != nil {
		t.Errorf("Props.RecurrenceRule() = %v", err)
	} else if !reflect.DeepEqual(roption, wantRecurrenceRule) {
		t.Errorf("Props.RecurrenceRule() = %v, want %v", roption, wantRecurrenceRule)
	}
}

func TestRecurrenceRuleIsAbsent(t *testing.T) {
	props := Props{}

	roption, err := props.RecurrenceRule()
	if roption != nil || err != nil {
		t.Errorf("Props.RecurrenceRule() = %v, %v, want nil, nil", roption, err)
	}
}

func TestRecurrenceRuleSetToNil(t *testing.T) {
	props := Props{
		"RRULE": []Prop{{
			Name:   "RRULE",
			Params: Params{},
			Value:  "FREQ=YEARLY;BYDAY=3SU;BYMONTH=3",
		}},
	}

	props.SetRecurrenceRule(nil)

	roption, err := props.RecurrenceRule()
	if roption != nil || err != nil {
		t.Errorf("Props.RecurrenceRule() = %v, %v, want nil, nil", roption, err)
	}
}

func TestRecurrenceRuleRoundTrip(t *testing.T) {
	recurrenceRule := &rrule.ROption{
		Freq:      rrule.YEARLY,
		Bymonth:   []int{3},
		Byweekday: []rrule.Weekday{rrule.SU.Nth(3)},
	}

	props := Props{}
	props.SetRecurrenceRule(recurrenceRule)

	if roption, err := props.RecurrenceRule(); err != nil {
		t.Errorf("Props.RecurrenceRule() = %v", err)
	} else if !reflect.DeepEqual(roption, recurrenceRule) {
		t.Errorf("Props.RecurrenceRule() = %v, want %v", roption, recurrenceRule)
	}
}
