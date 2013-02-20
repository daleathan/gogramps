package xml

import (
	"encoding/xml"
	"fmt"
)

type raw struct {
	XMLName  xml.Name
	Contents string `xml:",innerxml"`
}

// A DBObj is a referenceable object. It is any element with a handle.
type DBObj interface {
	GetHandle() string
}

type Tag struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ tag"`

	Handle   string `xml:"handle,attr"`
	Name     string `xml:"name,attr"`
	Color    string `xml:"color,attr"`
	Priority string `xml:"priority,attr"`
	Change   string `xml:"change,attr"`

	Unparsed []raw `xml:",any"`
}

func (o Tag) GetHandle() string { return o.Handle }

type dbObj struct {
	XMLName xml.Name

	ID     string `xml:"id,attr,omitempty"`
	Handle string `xml:"handle,attr"`
	Priv   int    `xml:"priv,attr,omitempty"`
	Change string `xml:"change,attr"`

	Unparsed []raw `xml:",any"`
}

func (o dbObj) GetHandle() string { return o.Handle }

type dateCommon struct {
	Quality   string `xml:"quality,attr,omitempty"`
	CFormat   string `xml:"cformat,attr,omitempty"`
	DualDated string `xml:"dualdated,attr,omitempty"`
	NewYear   string `xml:"newyear,attr,omitempty"`

	Unparsed []raw `xml:",any"`
}

type DateVal struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ dateval"`
	Val     string   `xml:"val,attr"`
	Type    string   `xml:"type,attr,omitempty"`

	dateCommon
}

type DateRange struct {
	XMLName xml.Name
	Start   string `xml:"start,attr"`
	Stop    string `xml:"stop,attr"`

	dateCommon
}

type hasDate struct {
	DateRange *DateRange `xml:"daterange"`
	DateSpan  *DateRange `xml:"datespan"`
	DateVal   *DateVal   `xml:"dateval"`
}

// Get a string representing the date (value, range or span).
func (v hasDate) GetDateString() string {
	if v.DateVal != nil {
		return v.DateVal.Val
	}
	if v.DateSpan != nil {
		return fmt.Sprintf("from %s to %s", v.DateSpan.Start, v.DateSpan.Stop)
	}
	if v.DateRange != nil {
		return fmt.Sprintf("between %s and %s", v.DateSpan.Start, v.DateSpan.Stop)
	}
	return ""
}

// A link is a reference to a DBObj.  It is any element with an hlink attr.
type Link interface {
	GetHLink() string
}

type GenericLink struct {
	XMLName  xml.Name
	HLink    string `xml:"hlink,attr"`
	Unparsed []raw  `xml:",any"`
}

func (l GenericLink) GetHLink() string {
	return l.HLink
}

type Event struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ event"`

	dbObj
	hasDate

	Type         *string       `xml:"type"`
	Place        *GenericLink  `xml:"place"`
	Description  *string       `xml:"description"`
	NoteRefs     []GenericLink `xml:"noteref"`
	CitationRefs []GenericLink `xml:"citationref"`
}

type Attribute struct {
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ attribute"`
	Priv     int      `xml:"priv,attr,omitempty"`
	Type     string   `xml:"type,attr"`
	Value    string   `xml:"value,attr"`
	Unparsed []raw    `xml:",any"`
}

type EventRef struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ eventref"`

	GenericLink

	Role       string        `xml:"role,attr"`
	Attributes []Attribute   `xml:"attribute"`
	NoteRefs   []GenericLink `xml:"noteref"`
	Unparsed   []raw         `xml:",any"`
}

// Get attribute value with type t.
func (e EventRef) GetAttribute(t string) string {
	for _, v := range e.Attributes {
		if v.Type == t {
			return v.Value
		}
	}
	return ""
}

type Region struct {
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ region"`
	Corner1X int      `xml:"corner1_x,attr"`
	Corner1Y int      `xml:"corner1_y,attr"`
	Corner2X int      `xml:"corner2_x,attr"`
	Corner2Y int      `xml:"corner2_y,attr"`
	Unparsed []raw    `xml:",any"`
}

type ObjRef struct {
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ objref"`
	GenericLink
	Region   Region   `xml:"region"`
	Unparsed []raw    `xml:",any"`
}

type Name struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ name"`
	Alt     int      `xml:"alt,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Priv    int      `xml:"priv,attr,omitempty"`
	// The following two fields are actually CDATA in the DTD
	Sort    int `xml:"sort,attr,omitempty"`
	Display int `xml:"display,attr,omitempty"`

	First        *string       `xml:"first"`
	Call         *string       `xml:"call"`
	Surnames     []string      `xml:"surname"`
	Suffix       *string       `xml:"suffix"`
	Nick         *string       `xml:"nick"`
	CitationRefs []GenericLink `xml:"citationref"`
	Unparsed     []raw         `xml:",any"`
}

func (n Name) GetSurname() string {
	if len(n.Surnames) > 0 {
		return n.Surnames[0]
	}
	return ""
}

func (n Name) GetFirstName() string {
	if n.First == nil {
		return ""
	}
	return *n.First
}

func (n Name) String() string { return n.GetSurname() + ", " + n.GetFirstName() }

type Person struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ person"`

	dbObj

	Gender       string        `xml:"gender"`
	Names        []Name        `xml:"name"`
	EventRefs    []EventRef    `xml:"eventref"`
	ObjRefs      []ObjRef      `xml:"objref"`
	ChildOfs     []GenericLink `xml:"childof"`
	ParentIns    []GenericLink `xml:"parentin"`
	NoteRefs     []GenericLink `xml:"noteref"`
	CitationRefs []GenericLink `xml:"citationref"`
	TagRefs      []GenericLink `xml:"tagref"`
}

// The preferred name is the first name that's not marked as an alternate.
func (p Person) GetPreferredName() *Name {
	for _, n := range p.Names {
		if n.Alt == 0 {
			return &n
		}
	}
	return nil
}

// Find the EventRef referencing event e.
func (p Person) FindEventRef(e Event) *EventRef {
	for _, v := range p.EventRefs {
		if v.GetHLink() == e.GetHandle() {
			return &v
		}
	}
	return nil
}

type People struct {
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ people"`
	Home     string   `xml:"home,attr"`
	Persons  []Person `xml:"person"`
	Unparsed []raw    `xml:",any"`
}

type ChildRef struct {
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ childref"`

  GenericLink

	Priv    int      `xml:"priv,attr,omitempty"`
  MRel    string   `xml:"mrel,attr,omitempty"`
  FRel    string   `xml:"frel,attr,omitempty"`
}

type Family struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ family"`

	dbObj

	Rel          *string       `xml:"rel"`
	Father       *GenericLink  `xml:"father"`
	Mother       *GenericLink  `xml:"mother"`
	EventRefs    []EventRef    `xml:"eventref"`
	ChildRefs    []ChildRef    `xml:"childref"`
	NoteRefs     []GenericLink `xml:"noteref"`
	CitationRefs []GenericLink `xml:"citationref"`
	TagRefs      []GenericLink `xml:"tagref"`
}

type Citation struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ citation"`

	dbObj
	hasDate

	Page       *string `xml:"page"`
	Confidence *string `xml:"confidence"`

	NoteRefs  []GenericLink `xml:"noteref"`
	ObjRefs   []ObjRef      `xml:"objref"`
	SourceRef GenericLink   `xml:"sourceref"`
}

type DataItem struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ data_item"`

	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`

	Unparsed []raw `xml:",any"`
}

type RepoRef struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ reporef"`

	GenericLink
	Priv   int    `xml:"priv,attr,omitempty"`
	CallNo string `xml:"callno,attr,omitempty"`
	Medium string `xml:"medium,attr,omitempty"`

	Unparsed []raw `xml:",any"`
}

type Source struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ source"`

	dbObj

	STitle    *string       `xml:"stitle"`
	SAuthor   *string       `xml:"sauthor"`
	SPubInfo  *string       `xml:"spubinfo"`
	SAbbrev   *string       `xml:"sabbrev"`
	NoteRefs  []GenericLink `xml:"noteref"`
	ObjRefs   []ObjRef      `xml:"objref"`
	DataItems []DataItem    `xml:"data_item"`
	RepoRefs  []RepoRef     `xml:"reporef"`
}

type Coord struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ coord"`

	Long string `xml:"long,attr"`
	Lat  string `xml:"lat,attr"`

	Unparsed []raw `xml:",any"`
}

type Location struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ location"`

	Street string `xml:"street,attr,omitempty"`
	Locality string `xml:"locality,attr,omitempty"`
	City     string `xml:"city,attr,omitempty"`
	Parish   string `xml:"parish,attr,omitempty"`
	County   string `xml:"county,attr,omitempty"`
	State    string `xml:"state,attr,omitempty"`
	Country  string `xml:"country,attr,omitempty"`
	Postal   string `xml:"postal,attr,omitempty"`
	Phone    string `xml:"phone,attr,omitempty"`

	Unparsed []raw `xml:",any"`
}

type URL struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ url"`

	Priv        int    `xml:"priv,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	HRef        string `xml:"href,attr"`
	Description string `xml:"description,attr,omitempty"`

	Unparsed []raw `xml:",any"`
}

type PlaceObj struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ placeobj"`

	dbObj

	PTitle       *string       `xml:"ptitle"`
	Coord        *Coord        `xml:"coord"`
	Locations    []Location    `xml:"location"`
	ObjRefs      []ObjRef      `xml:"objref"`
	URLs         []URL         `xml:"url"`
	NoteRefs     []GenericLink `xml:"noteref"`
	CitationRefs []GenericLink `xml:"citationref"`
}

type File struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ file"`

	Src         string `xml:"src,attr"`
	Mime        string `xml:"mime,attr"`
	Description string `xml:"description,attr"`

	Unparsed []raw `xml:",any"`
}

type Object struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ object"`

	dbObj

	File     File          `xml:"file"`
	NoteRefs []GenericLink `xml:"noteref"`
}

type Repository struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ repository"`

	dbObj

	RName string `xml:"rname"`
	Type  string `xml:"type"`
	URL   *URL   `xml:"url"`
}

type Note struct {
	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ note"`

	dbObj

	Format string `xml:"format,attr,omitempty"`
	Type   string `xml:"type,attr"`

	Text    string        `xml:"text"`
	TagRefs []GenericLink `xml:"tagref"`
}

// A Database represents an entire Gramps XML file.  
type Database struct {
	XMLName      xml.Name     `xml:"http://gramps-project.org/xml/1.5.0/ database"`
	Header       raw          `xml:"header"`
	Tags         []Tag        `xml:"tags>tag"`
	Events       []Event      `xml:"events>event"`
	People       People       `xml:"people"`
	Families     []Family     `xml:"families>family"`
	Citations    []Citation   `xml:"citations>citation"`
	Sources      []Source     `xml:"sources>source"`
	Places       []PlaceObj   `xml:"places>placeobj"`
	Objects      []Object     `xml:"objects>object"`
	Repositories []Repository `xml:"repositories>repository"`
	Notes        []Note       `xml:"notes>note"`

	Unparsed []raw `xml:",any"`
}
