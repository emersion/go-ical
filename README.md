# go-ical

[![Go Reference](https://pkg.go.dev/badge/github.com/emersion/go-ical.svg)](https://pkg.go.dev/github.com/emersion/go-ical)

An [iCalendar] library for Go supporting both [RFC 5545] calendar events and [RFC 7953] availability.

## Installation

```sh
go get github.com/nobleach/go-ical
```

## Overview

The library models iCalendar data as a tree of `Component` objects, each holding a `Props` map and a list of `Children`. Higher-level types (`Calendar`, `Event`, `VAvailability`, `Available`) wrap a `*Component` and add typed accessor methods.

## Decoding

Parse one or more calendars from any `io.Reader`:

```go
dec := ical.NewDecoder(r)
for {
    cal, err := dec.Decode()
    if err == io.EOF {
        break
    } else if err != nil {
        log.Fatal(err)
    }

    for _, event := range cal.Events() {
        summary, err := event.Props.Text(ical.PropSummary)
        if err != nil {
            log.Fatal(err)
        }
        start, err := event.DateTimeStart(time.UTC)
        if err != nil {
            log.Fatal(err)
        }
        end, err := event.DateTimeEnd(time.UTC)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("Event: %s (%s – %s)", summary, start, end)
    }
}
```

## Encoding

Build a calendar and write it to any `io.Writer`:

```go
event := ical.NewEvent()
event.Props.SetText(ical.PropUID, "uid@example.org")
event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
event.Props.SetText(ical.PropSummary, "Team standup")
event.Props.SetDateTime(ical.PropDateTimeStart, time.Now())
event.Props.SetDateTime(ical.PropDateTimeEnd, time.Now().Add(30*time.Minute))

cal := ical.NewCalendar()
cal.Props.SetText(ical.PropVersion, "2.0")
cal.Props.SetText(ical.PropProductID, "-//My App//EN")
cal.Children = append(cal.Children, event.Component)

if err := ical.NewEncoder(os.Stdout).Encode(cal); err != nil {
    log.Fatal(err)
}
```

## Recurrence rules

Events with an `RRULE` can produce a full recurrence set including `RDATE` additions and `EXDATE` exclusions:

```go
set, err := event.RecurrenceSet(time.UTC)
if err != nil {
    log.Fatal(err)
}
if set != nil {
    // All occurrences between now and 30 days from now.
    occurrences := set.Between(time.Now(), time.Now().Add(30*24*time.Hour), true)
    for _, t := range occurrences {
        log.Println(t)
    }
}
```

## RFC 7953 Availability (VAVAILABILITY)

`VAVAILABILITY` describes when a calendar user is generally available. It contains one or more `AVAILABLE` subcomponents that mark specific available windows. Outside those windows the `BUSYTYPE` property indicates whether the user is `BUSY`, `BUSY-UNAVAILABLE` (default), or `BUSY-TENTATIVE`.

### Decoding availability

```go
dec := ical.NewDecoder(r)
cal, err := dec.Decode()
if err != nil {
    log.Fatal(err)
}

for _, va := range cal.VAvailabilities() {
    start, _ := va.DateTimeStart(time.UTC)
    end, _ := va.DateTimeEnd(time.UTC)
    busyType, _ := va.BusyType() // defaults to BusyTypeBusyUnavailable
    log.Printf("Availability window: %s – %s (busy type: %s)", start, end, busyType)

    for _, a := range va.Available() {
        aStart, _ := a.DateTimeStart(time.UTC)
        aEnd, _ := a.DateTimeEnd(time.UTC)
        summary, _ := a.Props.Text(ical.PropSummary)
        log.Printf("  Available: %s (%s – %s)", summary, aStart, aEnd)
    }
}
```

### Encoding availability

```go
cal := ical.NewCalendar()
cal.Props.SetText(ical.PropVersion, "2.0")
cal.Props.SetText(ical.PropProductID, "-//My App//EN")

va := ical.NewVAvailability()
va.Props.SetText(ical.PropUID, "avail-001@example.com")
va.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
va.Props.SetDateTime(ical.PropDateTimeStart, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
va.Props.SetDateTime(ical.PropDateTimeEnd, time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
va.SetBusyType(ical.BusyTypeBusyUnavailable)

// Mark Monday–Friday 09:00–17:00 as available.
a := ical.NewAvailable()
a.Props.SetText(ical.PropUID, "avail-001-A@example.com")
a.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
a.Props.SetDateTime(ical.PropDateTimeStart, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC))
a.Props.SetDateTime(ical.PropDateTimeEnd, time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC))
a.Props.SetText(ical.PropSummary, "Business hours")

va.Children = append(va.Children, a.Component)
cal.Children = append(cal.Children, va.Component)

if err := ical.NewEncoder(os.Stdout).Encode(cal); err != nil {
    log.Fatal(err)
}
```

## Key types and constants

| Concept               | Go symbol                                                                         |
| --------------------- | --------------------------------------------------------------------------------- |
| Calendar object       | `ical.Calendar` / `ical.NewCalendar()`                                            |
| Event                 | `ical.Event` / `ical.NewEvent()`                                                  |
| Availability block    | `ical.VAvailability` / `ical.NewVAvailability()`                                  |
| Available window      | `ical.Available` / `ical.NewAvailable()`                                          |
| Busy types            | `ical.BusyTypeBusy`, `ical.BusyTypeBusyUnavailable`, `ical.BusyTypeBusyTentative` |
| Event status          | `ical.EventTentative`, `ical.EventConfirmed`, `ical.EventCancelled`               |
| MIME type / extension | `ical.MIMEType` (`text/calendar`), `ical.Extension` (`ics`)                       |

Property name constants (`ical.PropSummary`, `ical.PropUID`, etc.) and parameter name constants (`ical.ParamTimezoneID`, etc.) are defined in `enums.go`.

## License

MIT

[iCalendar]: https://tools.ietf.org/html/rfc5545
[RFC 5545]: https://tools.ietf.org/html/rfc5545
[RFC 7953]: https://tools.ietf.org/html/rfc7953
