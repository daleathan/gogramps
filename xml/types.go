package xml

import (
	"encoding/xml"
	"fmt"
  "strings"
)

type raw struct {
	XMLName  xml.Name
	Contents string `xml:",innerxml"`
}

type Created struct {
	Date    string `xml:"date,attr"`
	Version string `xml:"version,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ created"`
	Unparsed []*raw   `xml:",any"`
}

type Researcher struct {
	ResName     *string `xml:"resname"`
	ResAddr     *string `xml:"resaddr"`
	ResLocality *string `xml:"reslocality"`
	ResCity     *string `xml:"rescity"`
	ResState    *string `xml:"resstate"`
	ResCountry  *string `xml:"rescountry"`
	ResPostal   *string `xml:"respostal"`
	ResPhone    *string `xml:"resphone"`
	ResEMail    *string `xml:"resemail"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ researcher"`
	Unparsed []*raw   `xml:",any"`
}

type Header struct {
	Created    Created     `xml:"created"`
	Researcher *Researcher `xml:"researcher"`
	MediaPath  *string     `xml:"mediapath"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ header"`
	Unparsed []*raw   `xml:",any"`
}

type NameFormat struct {
	Number string `xml:"number,attr"`
	Name   string `xml:"name,attr"`
	FmtStr string `xml:"fmt_str,attr"`
	Active int    `xml:"active,attr,omitempty"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ format"`
	Unparsed []*raw   `xml:",any"`
}

// A DBObj is a referenceable object. It is any element with a handle.
type DBObj interface {
	GetHandle() string
  // The name of the XML element.
  GetName() string
}

type Tag struct {
	Handle   string `xml:"handle,attr"`
	Name     string `xml:"name,attr"`
	Color    string `xml:"color,attr"`
	Priority string `xml:"priority,attr"`
	Change   string `xml:"change,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ tag"`
	Unparsed []*raw   `xml:",any"`
}

func (o Tag) GetHandle() string { return o.Handle }

type dbObj struct {
	ID     string `xml:"id,attr,omitempty"`
	Handle string `xml:"handle,attr"`
	Priv   int    `xml:"priv,attr,omitempty"`
	Change string `xml:"change,attr"`

	XMLName  xml.Name
	Unparsed []*raw `xml:",any"`
}

func (o dbObj) GetHandle() string { return o.Handle }
func (o dbObj) GetName() string { return o.XMLName.Local }

type dateCommon struct {
	Quality   string `xml:"quality,attr,omitempty"`
	CFormat   string `xml:"cformat,attr,omitempty"`
	DualDated string `xml:"dualdated,attr,omitempty"`
	NewYear   string `xml:"newyear,attr,omitempty"`

	Unparsed []*raw `xml:",any"`
}

type DateStr struct {
	Val      string   `xml:"val,attr"`
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ datestr"`
	Unparsed []*raw   `xml:",any"`
}

type DateVal struct {
	dateCommon
	Val  string `xml:"val,attr"`
	Type string `xml:"type,attr,omitempty"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ dateval"`
}

type DateRange struct {
	dateCommon
	Start string `xml:"start,attr"`
	Stop  string `xml:"stop,attr"`

	XMLName xml.Name
}

type hasDate struct {
	DateRange *DateRange `xml:"daterange"`
	DateSpan  *DateRange `xml:"datespan"`
	DateVal   *DateVal   `xml:"dateval"`
	DateStr   *DateStr   `xml:"datestr"`
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
	if v.DateStr != nil {
		return v.DateStr.Val
	}
	return ""
}

// A link is a reference to a DBObj.  It is any element with an hlink attr.
type Link interface {
	GetHLink() string
}

type GenericLink struct {
	HLink    string `xml:"hlink,attr"`
	XMLName  xml.Name
	Unparsed []*raw `xml:",any"`
}

func (l GenericLink) GetHLink() string {
	return l.HLink
}

type Attribute struct {
	Priv  int    `xml:"priv,attr,omitempty"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`

	CitationRefs []*GenericLink `xml:"citationref"`
	NoteRefs     []*GenericLink `xml:"noteref"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ attribute"`
	Unparsed []*raw   `xml:",any"`
}

type Event struct {
	dbObj
	hasDate

	Type         *string        `xml:"type"`
	Place        *GenericLink   `xml:"place"`
	Description  *string        `xml:"description"`
	Attributes   []*Attribute   `xml:"attribute"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`
	ObjRefs      []*ObjRef      `xml:"objref"`
	XMLName      xml.Name       `xml:"http://gramps-project.org/xml/1.5.0/ event"`
}

type EventRef struct {
	GenericLink

	Role       string         `xml:"role,attr"`
	Attributes []*Attribute   `xml:"attribute"`
	NoteRefs   []*GenericLink `xml:"noteref"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ eventref"`
	Unparsed []*raw   `xml:",any"`
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
	Corner1X int      `xml:"corner1_x,attr"`
	Corner1Y int      `xml:"corner1_y,attr"`
	Corner2X int      `xml:"corner2_x,attr"`
	Corner2Y int      `xml:"corner2_y,attr"`
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ region"`
	Unparsed []*raw   `xml:",any"`
}

type ObjRef struct {
	GenericLink

	Region       *Region        `xml:"region"`
	Attributes   []*Attribute   `xml:"attribute"`
	CitationRefs []*GenericLink `xml:"citationref"`
	NoteRefs     []*GenericLink `xml:"noteref"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ objref"`
	Unparsed []*raw   `xml:",any"`
}

type Surname struct {
	Prefix     string `xml:"prefix,attr,omitempty"`
	Prim       int    `xml:"prim,attr,omitempty"`
	Derivation string `xml:"derivation,attr,omitempty"`
	Connector  string `xml:"connector,attr,omitempty"`
	Value      string `xml:",chardata"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ surname"`
	Unparsed []*raw   `xml:",any"`
}

type Name struct {
	Alt  int    `xml:"alt,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	Priv int    `xml:"priv,attr,omitempty"`
	// The following two fields are actually CDATA in the DTD
	Sort    int `xml:"sort,attr,omitempty"`
	Display int `xml:"display,attr,omitempty"`

	hasDate
	First        *string        `xml:"first"`
	Call         *string        `xml:"call"`
	Surnames     []*Surname     `xml:"surname"`
	Suffix       *string        `xml:"suffix"`
	Title        *string        `xml:"title"`
	Nick         *string        `xml:"nick"`
	FamilyNick   *string        `xml:"familynick"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ name"`
	Unparsed []*raw   `xml:",any"`
}

func (n Name) GetSurname() string {
	if len(n.Surnames) > 0 {
		return n.Surnames[0].Value
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

// A short form of the name: surname plus the first part of the firstname.
func (n Name) Short() string {
  return n.GetSurname() + ", " + strings.Split(n.GetFirstName(), " ")[0]
}

type Temple struct {
	Val      string   `xml:"val,attr"`
	Unparsed []*raw   `xml:",any"`
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ temple"`
}

type Status struct {
	Val      string   `xml:"val,attr"`
	Unparsed []*raw   `xml:",any"`
	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ status"`
}

type LDSOrd struct {
	Priv int    `xml:"priv,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`

	hasDate
	Temple       *Temple        `xml:"temple"`
	Place        *GenericLink   `xml:"place"`
	Status       *Status        `xml:"status"`
	SealedTo     *GenericLink   `xml:"sealed_to"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ lds_ord"`
	Unparsed []*raw   `xml:",any"`
}

type Address struct {
	Priv int `xml:"priv,attr,omitempty"`

	hasDate
	Street       *string        `xml:"street"`
	Locality     *string        `xml:"locality"`
	City         *string        `xml:"city"`
	County       *string        `xml:"county"`
	State        *string        `xml:"state"`
	Country      *string        `xml:"country"`
	Postal       *string        `xml:"postal"`
	Phone        *string        `xml:"phone"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ address"`
	Unparsed []*raw   `xml:",any"`
}

type PersonRef struct {
	GenericLink
	Priv int    `xml:"priv,attr,omitempty"`
	Rel  string `xml:"rel,attr"`

	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ personref"`
}

type Person struct {
	dbObj

	Gender       string         `xml:"gender"`
	Names        []*Name        `xml:"name"`
	EventRefs    []*EventRef    `xml:"eventref"`
	LDSOrds      []*LDSOrd      `xml:"lds_ord"`
	ObjRefs      []*ObjRef      `xml:"objref"`
	Addresses    []*Address     `xml:"address"`
	Attributes   []*Attribute   `xml:"attribute"`
	URLs         []*URL         `xml:"url"`
	ChildOfs     []*GenericLink `xml:"childof"`
	ParentIns    []*GenericLink `xml:"parentin"`
	PersonRefs   []*PersonRef   `xml:"personref"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`
	TagRefs      []*GenericLink `xml:"tagref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ person"`
}

// The preferred name is the first name that's not marked as an alternate.
func (p Person) GetPreferredName() *Name {
	for _, n := range p.Names {
		if n.Alt == 0 {
			return n
		}
	}
	return nil
}

// Find the EventRef referencing event e.
func (p Person) FindEventRef(e Event) *EventRef {
	for _, v := range p.EventRefs {
		if v.GetHLink() == e.GetHandle() {
			return v
		}
	}
	return nil
}

type People struct {
	Home     string    `xml:"home,attr"`
	Persons  []*Person `xml:"person"`
	XMLName  xml.Name  `xml:"http://gramps-project.org/xml/1.5.0/ people"`
	Unparsed []*raw    `xml:",any"`
}

type ChildRef struct {
	GenericLink

	Priv int    `xml:"priv,attr,omitempty"`
	MRel string `xml:"mrel,attr,omitempty"`
	FRel string `xml:"frel,attr,omitempty"`

	CitationRefs []*GenericLink `xml:"citationref"`
	NoteRefs     []*GenericLink `xml:"noteref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ childref"`
}

type Family struct {
	dbObj
	Rel          *string        `xml:"rel"`
	Father       *GenericLink   `xml:"father"`
	Mother       *GenericLink   `xml:"mother"`
	EventRefs    []*EventRef    `xml:"eventref"`
	LDSOrds      []*LDSOrd      `xml:"lds_ord"`
	ObjRefs      []*ObjRef      `xml:"objref"`
	ChildRefs    []*ChildRef    `xml:"childref"`
	Attributes   []*Attribute   `xml:"attribute"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`
	TagRefs      []*GenericLink `xml:"tagref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ family"`
}

type Citation struct {
	dbObj
	hasDate

	Page       *string `xml:"page"`
	Confidence *string `xml:"confidence"`

	NoteRefs  []*GenericLink `xml:"noteref"`
	ObjRefs   []*ObjRef      `xml:"objref"`
	SourceRef GenericLink    `xml:"sourceref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ citation"`
}

type DataItem struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ data_item"`
	Unparsed []*raw   `xml:",any"`
}

type RepoRef struct {
	GenericLink
	Priv   int    `xml:"priv,attr,omitempty"`
	CallNo string `xml:"callno,attr,omitempty"`
	Medium string `xml:"medium,attr,omitempty"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ reporef"`
	Unparsed []*raw   `xml:",any"`
}

type Source struct {
	dbObj

	STitle    *string        `xml:"stitle"`
	SAuthor   *string        `xml:"sauthor"`
	SPubInfo  *string        `xml:"spubinfo"`
	SAbbrev   *string        `xml:"sabbrev"`
	NoteRefs  []*GenericLink `xml:"noteref"`
	ObjRefs   []*ObjRef      `xml:"objref"`
	DataItems []*DataItem    `xml:"data_item"`
	RepoRefs  []*RepoRef     `xml:"reporef"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ source"`
}

type Coord struct {
	Long string `xml:"long,attr"`
	Lat  string `xml:"lat,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ coord"`
	Unparsed []*raw   `xml:",any"`
}

type Location struct {
	Street   string `xml:"street,attr,omitempty"`
	Locality string `xml:"locality,attr,omitempty"`
	City     string `xml:"city,attr,omitempty"`
	Parish   string `xml:"parish,attr,omitempty"`
	County   string `xml:"county,attr,omitempty"`
	State    string `xml:"state,attr,omitempty"`
	Country  string `xml:"country,attr,omitempty"`
	Postal   string `xml:"postal,attr,omitempty"`
	Phone    string `xml:"phone,attr,omitempty"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ location"`
	Unparsed []*raw   `xml:",any"`
}

type URL struct {
	Priv        int    `xml:"priv,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	HRef        string `xml:"href,attr"`
	Description string `xml:"description,attr,omitempty"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ url"`
	Unparsed []*raw   `xml:",any"`
}

type PlaceObj struct {
	dbObj

	PTitle       *string        `xml:"ptitle"`
	Coord        *Coord         `xml:"coord"`
	Locations    []*Location    `xml:"location"`
	ObjRefs      []*ObjRef      `xml:"objref"`
	URLs         []*URL         `xml:"url"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ placeobj"`
}

type File struct {
	Src         string `xml:"src,attr"`
	Mime        string `xml:"mime,attr"`
	Description string `xml:"description,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ file"`
	Unparsed []*raw   `xml:",any"`
}

type Object struct {
	dbObj
	hasDate

	File         File           `xml:"file"`
	Attributes   []*Attribute   `xml:"attribute"`
	NoteRefs     []*GenericLink `xml:"noteref"`
	CitationRefs []*GenericLink `xml:"citationref"`
	TagRefs      []*GenericLink `xml:"tagref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ object"`
}

type Repository struct {
	dbObj

	RName     string         `xml:"rname"`
	Type      string         `xml:"type"`
	Addresses []*Address     `xml:"address"`
	URL       *URL           `xml:"url"`
	NoteRefs  []*GenericLink `xml:"noteref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ repository"`
}

type Range struct {
	Start string `xml:"start,attr"`
	End   string `xml:"end,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ range"`
	Unparsed []*raw   `xml:",any"`
}

type Style struct {
	Name   string   `xml:"name,attr"`
	Value  string   `xml:"value,attr,omitempty"`
	Ranges []*Range `xml:"range"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ style"`
	Unparsed []*raw   `xml:",any"`
}

type Note struct {
	dbObj

	Format string `xml:"format,attr,omitempty"`
	Type   string `xml:"type,attr"`

	Text    string         `xml:"text"`
	Styles  []*Style       `xml:"style"`
	TagRefs []*GenericLink `xml:"tagref"`

	XMLName xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ note"`
}

type Bookmark struct {
	GenericLink
	Target string `xml:"target,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ bookmark"`
	Unparsed []*raw   `xml:",any"`
}

type NameMap struct {
	Type  string `xml:"type,attr"`
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ map"`
	Unparsed []*raw   `xml:",any"`
}

// A Database represents an entire Gramps XML file.  
type Database struct {
	Header       Header        `xml:"header"`
	NameFormats  []*NameFormat `xml:"name-formats>format"`
	Tags         []*Tag        `xml:"tags>tag"`
	Events       []*Event      `xml:"events>event"`
	People       People        `xml:"people"`
	Families     []*Family     `xml:"families>family"`
	Citations    []*Citation   `xml:"citations>citation"`
	Sources      []*Source     `xml:"sources>source"`
	Places       []*PlaceObj   `xml:"places>placeobj"`
	Objects      []*Object     `xml:"objects>object"`
	Repositories []*Repository `xml:"repositories>repository"`
	Notes        []*Note       `xml:"notes>note"`
	Bookmarks    []*Bookmark   `xml:"bookmarks>bookmark"`
	NameMaps     []*NameMap    `xml:"namemaps>map"`

	XMLName  xml.Name `xml:"http://gramps-project.org/xml/1.5.0/ database"`
	Unparsed []*raw   `xml:",any"`
}
