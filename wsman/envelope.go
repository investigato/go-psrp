package wsman

import (
	"encoding/xml"
)

// Envelope represents a SOAP 1.2 envelope for WS-Management messages.
type Envelope struct {
	XMLName xml.Name `xml:"s:Envelope"`

	// Namespace declarations
	NsSoap    string `xml:"xmlns:s,attr"`
	NsAddr    string `xml:"xmlns:wsa,attr"`
	NsWsman   string `xml:"xmlns:wsman,attr"`
	NsMsWsman string `xml:"xmlns:wsmv,attr"`
	NsShellNs string `xml:"xmlns:rsp,attr,omitempty"`
	NsXsiAttr string `xml:"xmlns:xsi,attr,omitempty"`

	Header *Header `xml:"s:Header"`
	Body   *Body   `xml:"s:Body"`
}

// Header represents the SOAP header containing WS-Addressing and WS-Management headers.
type Header struct {
	// WS-Addressing headers
	Action    *ActionHeader `xml:"wsa:Action,omitempty"`
	To        string        `xml:"wsa:To,omitempty"`
	MessageID string        `xml:"wsa:MessageID,omitempty"`
	ReplyTo   *ReplyTo      `xml:"wsa:ReplyTo,omitempty"`

	// WS-Management headers
	ResourceURI      *ResourceURIHeader     `xml:"wsman:ResourceURI,omitempty"`
	MaxEnvelopeSize  *MaxEnvelopeSizeHeader `xml:"wsman:MaxEnvelopeSize,omitempty"`
	OperationTimeout string                 `xml:"wsman:OperationTimeout,omitempty"`
	Locale           *Locale                `xml:"wsman:Locale,omitempty"`
	DataLocale       *DataLocale            `xml:"wsmv:DataLocale,omitempty"`
	SessionID        *SessionIDHeader       `xml:"wsmv:SessionId,omitempty"`

	// Shell-specific headers
	SelectorSet *SelectorSet `xml:"wsmv:SelectorSet,omitempty"`
	OptionSet   *OptionSet   `xml:"wsmv:OptionSet,omitempty"`
}

// ActionHeader represents Action element with mustUnderstand attribute.
type ActionHeader struct {
	MustUnderstand string `xml:"s:mustUnderstand,attr,omitempty"`
	Value          string `xml:",chardata"`
}

// SessionIDHeader represents SessionId element with mustUnderstand attribute.
type SessionIDHeader struct {
	MustUnderstand string `xml:"s:mustUnderstand,attr,omitempty"`
	Value          string `xml:",chardata"`
}

// ResourceURIHeader represents ResourceURI element with mustUnderstand attribute.
type ResourceURIHeader struct {
	MustUnderstand string `xml:"s:mustUnderstand,attr,omitempty"`
	Value          string `xml:",chardata"`
}

// MaxEnvelopeSizeHeader represents MaxEnvelopeSize element with mustUnderstand attribute.
type MaxEnvelopeSizeHeader struct {
	MustUnderstand string `xml:"s:mustUnderstand,attr,omitempty"`
	Value          int    `xml:",chardata"`
}

// Locale representing xml:lang attribute
type Locale struct {
	MustUnderstand bool   `xml:"s:mustUnderstand,attr,omitempty"`
	Lang           string `xml:"xml:lang,attr,omitempty"`
}

// DataLocale representing xml:lang attribute
type DataLocale struct {
	MustUnderstand bool   `xml:"s:mustUnderstand,attr,omitempty"`
	Lang           string `xml:"xml:lang,attr,omitempty"`
}

// ReplyTo represents the WS-Addressing ReplyTo element.
type ReplyTo struct {
	Address *AddressHeader `xml:"wsa:Address"`
}

// AddressHeader represents Address element with mustUnderstand attribute.
type AddressHeader struct {
	MustUnderstand string `xml:"s:mustUnderstand,attr,omitempty"`
	Value          string `xml:",chardata"`
}

// SelectorSet contains selectors for targeting specific resources.
type SelectorSet struct {
	Selectors []Selector `xml:"wsmv:Selector"`
}

// OptionSet contains options for the operation.
type OptionSet struct {
	MustUnderstand string   `xml:"s:mustUnderstand,attr,omitempty"`
	Options        []Option `xml:"wsmv:Option"`
}

// Option represents a single option.
type Option struct {
	MustComply string `xml:"MustComply,attr,omitempty"`
	Name       string `xml:"Name,attr"`
	Value      string `xml:",chardata"`
}

// Body represents the SOAP body.
type Body struct {
	Content []byte `xml:",innerxml"`
}

// NewEnvelope creates a new SOAP envelope with required namespace declarations.
func NewEnvelope() *Envelope {
	return &Envelope{
		NsSoap:    NsSoap,
		NsAddr:    NsAddressing,
		NsWsman:   NsWsman,
		NsMsWsman: NsWsmanMicrosoft,
		Header:    &Header{},
		Body:      &Body{},
	}
}

// WithAction sets the WS-Addressing Action header.
func (e *Envelope) WithAction(action string) *Envelope {
	e.Header.Action = &ActionHeader{
		MustUnderstand: "true",
		Value:          action,
	}
	return e
}

// WithTo sets the WS-Addressing To header (the endpoint URL).
func (e *Envelope) WithTo(to string) *Envelope {
	e.Header.To = to
	return e
}

// WithMessageID sets the WS-Addressing MessageID header.
func (e *Envelope) WithMessageID(messageID string) *Envelope {
	e.Header.MessageID = messageID
	return e
}

// WithReplyTo sets the WS-Addressing ReplyTo header.
func (e *Envelope) WithReplyTo(address string) *Envelope {
	e.Header.ReplyTo = &ReplyTo{
		Address: &AddressHeader{
			MustUnderstand: "true",
			Value:          address,
		},
	}
	return e
}

// WithResourceURI sets the WS-Management ResourceURI header.
func (e *Envelope) WithResourceURI(uri string) *Envelope {
	e.Header.ResourceURI = &ResourceURIHeader{
		MustUnderstand: "true",
		Value:          uri,
	}
	return e
}

// WithMaxEnvelopeSize sets the WS-Management MaxEnvelopeSize header.
func (e *Envelope) WithMaxEnvelopeSize(size int) *Envelope {
	e.Header.MaxEnvelopeSize = &MaxEnvelopeSizeHeader{
		MustUnderstand: "true",
		Value:          size,
	}
	return e
}

// WithOperationTimeout sets the WS-Management OperationTimeout header.
// The timeout should be in ISO 8601 duration format (e.g., "PT60S" for 60 seconds).
func (e *Envelope) WithOperationTimeout(timeout string) *Envelope {
	e.Header.OperationTimeout = timeout
	return e
}

// WithShellNamespace adds the Windows Shell namespace to the envelope.
func (e *Envelope) WithShellNamespace() *Envelope {
	e.NsShellNs = NsShell
	return e
}

// WithSelector adds a selector to the SelectorSet.
func (e *Envelope) WithSelector(name, value string) *Envelope {
	if e.Header.SelectorSet == nil {
		e.Header.SelectorSet = &SelectorSet{}
	}
	e.Header.SelectorSet.Selectors = append(e.Header.SelectorSet.Selectors,
		Selector{Name: name, Value: value})
	return e
}

// WithOption adds an option to the OptionSet.
func (e *Envelope) WithOption(name, value string) *Envelope {
	if e.Header.OptionSet == nil {
		e.Header.OptionSet = &OptionSet{
			MustUnderstand: "true",
		}
	}
	e.Header.OptionSet.Options = append(e.Header.OptionSet.Options,
		Option{Name: name, Value: value})
	return e
}

// WithOptionMustComply adds an option with MustComply="true" to the OptionSet.
func (e *Envelope) WithOptionMustComply(name, value string) *Envelope {
	if e.Header.OptionSet == nil {
		e.Header.OptionSet = &OptionSet{
			MustUnderstand: "true",
		}
	}
	e.Header.OptionSet.Options = append(e.Header.OptionSet.Options,
		Option{MustComply: "true", Name: name, Value: value})
	return e
}

// WithBody sets the SOAP body content.
func (e *Envelope) WithBody(content []byte) *Envelope {
	e.Body.Content = content
	return e
}

// WithSessionID sets the WS-Management SessionId header.
func (e *Envelope) WithSessionID(sessionID string) *Envelope {
	e.Header.SessionID = &SessionIDHeader{
		MustUnderstand: "true",
		Value:          sessionID,
	}
	return e
}

// WithLocale sets the WS-Management Locale header.
func (e *Envelope) WithLocale(lang string) *Envelope {
	e.Header.Locale = &Locale{
		Lang:           lang,
		MustUnderstand: false,
	}
	return e
}

// WithDataLocale sets the WS-Management DataLocale header.
func (e *Envelope) WithDataLocale(lang string) *Envelope {
	e.Header.DataLocale = &DataLocale{
		Lang:           lang,
		MustUnderstand: false,
	}
	return e
}

// Marshal serializes the envelope to XML.
func (e *Envelope) Marshal() ([]byte, error) {
	return xml.Marshal(e)
}

// MarshalIndent serializes the envelope to indented XML.
func (e *Envelope) MarshalIndent(prefix, indent string) ([]byte, error) {
	return xml.MarshalIndent(e, prefix, indent)
}
