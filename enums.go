package ical

// Components as defined in RFC 5545 section 3.6.
const (
	CompCalendar = "VCALENDAR"
	CompEvent    = "VEVENT"
	CompToDo     = "VTODO"
	CompJournal  = "VJOURNAL"
	CompFreeBusy = "VFREEBUSY"
	CompTimezone = "VTIMEZONE"
	CompAlarm    = "VALARM"
)

// Timezone components.
const (
	CompTimezoneStandard = "STANDARD"
	CompTimezoneDaylight = "DAYLIGHT"
)

// Properties as defined in RFC 5545 section 3.7 and section 3.8.
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

// Property parameters as defined in RFC 5545 section 3.2.
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

// ValueType is the type of a property.
type ValueType string

// Value types as defined in RFC 5545 section 3.3.
const (
	ValueDefault         ValueType = ""
	ValueBinary          ValueType = "BINARY"
	ValueBool            ValueType = "BOOLEAN"
	ValueCalendarAddress ValueType = "CAL-ADDRESS"
	ValueDate            ValueType = "DATE"
	ValueDateTime        ValueType = "DATE-TIME"
	ValueDuration        ValueType = "DURATION"
	ValueFloat           ValueType = "FLOAT"
	ValueInt             ValueType = "INTEGER"
	ValuePeriod          ValueType = "PERIOD"
	ValueRecurrence      ValueType = "RECUR"
	ValueText            ValueType = "TEXT"
	ValueTime            ValueType = "TIME"
	ValueURI             ValueType = "URI"
	ValueUTCOffset       ValueType = "UTC-OFFSET"
)

var defaultValueTypes = map[string]ValueType{
	PropCalendarScale:      ValueText,
	PropMethod:             ValueText,
	PropProductID:          ValueText,
	PropVersion:            ValueText,
	PropAttach:             ValueURI, // can be binary
	PropCategories:         ValueText,
	PropClass:              ValueText,
	PropComment:            ValueText,
	PropDescription:        ValueText,
	PropGeo:                ValueFloat,
	PropLocation:           ValueText,
	PropPercentComplete:    ValueInt,
	PropPriority:           ValueInt,
	PropResources:          ValueText,
	PropStatus:             ValueText,
	PropSummary:            ValueText,
	PropCompleted:          ValueDateTime,
	PropDateTimeEnd:        ValueDateTime, // can be date
	PropDue:                ValueDateTime, // can be date
	PropDateTimeStart:      ValueDateTime, // can be date
	PropDuration:           ValueDuration,
	PropFreeBusy:           ValuePeriod,
	PropTransparency:       ValueText,
	PropTimezoneID:         ValueText,
	PropTimezoneName:       ValueText,
	PropTimezoneOffsetFrom: ValueUTCOffset,
	PropTimezoneOffsetTo:   ValueUTCOffset,
	PropTimezoneURL:        ValueURI,
	PropAttendee:           ValueCalendarAddress,
	PropContact:            ValueText,
	PropOrganizer:          ValueCalendarAddress,
	PropRecurrenceID:       ValueDateTime, // can be date
	PropRelatedTo:          ValueText,
	PropURL:                ValueURI,
	PropUID:                ValueText,
	PropExceptionDates:     ValueDateTime, // can be date
	PropRecurrenceDates:    ValueDateTime, // can be date or period
	PropRecurrenceRule:     ValueRecurrence,
	PropAction:             ValueText,
	PropRepeat:             ValueInt,
	PropTrigger:            ValueDuration, // can be date-time
	PropCreated:            ValueDateTime,
	PropDateTimeStamp:      ValueDateTime,
	PropLastModified:       ValueDateTime,
	PropSequence:           ValueInt,
	PropRequestStatus:      ValueText,
}

type EventStatus string

const (
	EventTentative EventStatus = "TENTATIVE"
	EventConfirmed EventStatus = "CONFIRMED"
	EventCancelled EventStatus = "CANCELLED"
)
