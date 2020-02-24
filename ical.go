// Package ical implements the iCalendar file format.
//
// iCalendar is defined in RFC 5545.
package ical

type Params map[string][]string

type Property struct {
	Name   string
	Params Params
	Value  string
}

type Properties map[string][]Property

type Component struct {
	Name       string
	Properties Properties
	Children   []Component
}

const (
	CompCalendar = "VCALENDAR"
	CompEvent    = "VEVENT"
	CompToDo     = "VTODO"
	CompJournal  = "VJOURNAL"
	CompFreeBusy = "VFREEBUSY"
	CompTimezone = "VTIMEZONE"
	CompAlarm    = "VALARM"
)

const (
	CompTimezoneStandard = "STANDARD"
	CompTimezoneDaylight = "DAYLIGHT"
)

const (
	// Calendar properties
	PropCalendarScale = "CALSCALE"
	PropMethod        = "METHOD"
	PropProductID     = "PRODID"
	PropVersion       = "VERSION"

	// Component properties
	PropAttach          = "ATTACH"
	PropCategories      = "CATEGORIES"
	PropClass           = "CLASS"
	PropComment         = "COMMENT"
	PropDescription     = "DESCRIPTION"
	PropGeo             = "GEO"
	PropLocation        = "LOCATION"
	PropPercentComplete = "PERCENT-COMPLETE"
	PropPriority        = "PRIORITY"
	PropResources       = "RESOURCES"
	PropStatus          = "STATUS"
	PropSummary         = "SUMMARY"

	// Date and time component properties
	PropCompleted     = "COMPLETED"
	PropDateTimeEnd   = "DTEND"
	PropDue           = "DUE"
	PropDateTimeStart = "DTSTART"
	PropDuration      = "DURATION"
	PropFreeBusy      = "FREEBUSY"
	PropTransparency  = "TRANSP"

	// Timezone component properties
	PropTimezoneID         = "TZID"
	PropTimezoneName       = "TZNAME"
	PropTimezoneOffsetFrom = "TZOFFSETFROM"
	PropTimezoneOffsetTo   = "TZOFFSETTO"
	PropTimezoneURL        = "TZURL"

	// Relationship component properties
	PropAttendee     = "ATTENDEE"
	PropContact      = "CONTACT"
	PropOrganizer    = "ORGANIZER"
	PropRecurrenceID = "RECURRENCE-ID"
	PropRelatedTo    = "RELATED-TO"
	PropURL          = "URL"
	PropUID          = "UID"

	// Recurrence component properties
	PropExceptionDates  = "EXDATE"
	PropRecurrenceDates = "RDATE"
	PropRecurrenceRule  = "RRULE"

	// Alarm component properties
	PropAction  = "ACTION"
	PropRepeat  = "REPEAT"
	PropTrigger = "TRIGGER"

	// Change management component properties
	PropCreated       = "CREATED"
	PropDateTimeStamp = "DTSTAMP"
	PropLastModified  = "LAST-MODIFIED"
	PropSequence      = "SEQUENCE"

	// Miscellaneous component properties
	PropRequestStatus = "REQUEST-STATUS"
)

const (
	ParamAltRep              = "ALTREP"
	ParamCommonName          = "CN"
	ParamCalendarUserType    = "CUTYPE"
	ParamDelegatedFrom       = "DELEGATED-FROM"
	ParamDelegatedTo         = "DELEGATED-TO"
	ParamDir                 = "DIR"
	ParamEncoding            = "ENCODING"
	ParamFormatType          = "FMTTYPE"
	ParamFreeBusyType        = "FBTYPE"
	ParamLanguage            = "LANGUAGE"
	ParamMember              = "MEMBER"
	ParamParticipationStatus = "PARTSTAT"
	ParamRange               = "RANGE"
	ParamRelated             = "RELATED"
	ParamRelationshipType    = "RELTYPE"
	ParamRole                = "ROLE"
	ParamRSVP                = "RSVP"
	ParamSentBy              = "SENT-BY"
	ParamTimezoneID          = "TZID"
	ParamValue               = "VALUE"
)

type ValueType string

const (
	ValueBinary          ValueType = "BINARY"
	ValueBoolean         ValueType = "BOOLEAN"
	ValueCalendarAddress ValueType = "CAL-ADDRESS"
	ValueDate            ValueType = "DATE"
	ValueDateTime        ValueType = "DATE-TIME"
	ValueDuration        ValueType = "DURATION"
	ValueFloat           ValueType = "FLOAT"
	ValueInteger         ValueType = "INTEGER"
	ValuePeriod          ValueType = "PERIOD"
	ValueRecurrence      ValueType = "RECUR"
	ValueText            ValueType = "TEXT"
	ValueTime            ValueType = "TIME"
	ValueURI             ValueType = "URI"
	ValueUTCOffset       ValueType = "UTC-OFFSET"
)

type Calendar struct {
	Component
}
