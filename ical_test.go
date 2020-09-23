package ical

import (
	"strings"
	"testing"
	"time"
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

func TestSetDate(b *testing.T) {
	localTimezone, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		b.Fatal(err)
	}

	testCases := []struct {
		Alias          string
		Date           time.Time
		ExpectedResult string
	}{
		{
			Alias:          "UTC",
			Date:           time.Date(2020, time.September, 20, 17, 7, 0, 0, time.UTC),
			ExpectedResult: "20200920T170700Z",
		},
		{
			Alias:          "local-tz",
			Date:           time.Date(2020, time.September, 20, 17, 7, 0, 0, localTimezone),
			ExpectedResult: "20200920T150700Z",
		},
	}

	//nolint:scopelint
	for _, tCase := range testCases {
		testFn := func(t *testing.T) {
			p := NewProp("FakeProp")
			p.SetDateTime(tCase.Date)
			if got, want := p.Value, tCase.ExpectedResult; got != want {
				t.Errorf("bad result: %s, expected: %s", got, want)
			}
		}
		b.Run(tCase.Alias, testFn)
	}
}
