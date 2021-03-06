//note: this entire file is a fuck fest of type aliases and nested types
//it makes no sense at all, but thats just how psd files work.

package descriptor

import (
	"os"

	util "github.com/fabulousduck/pisdui/pisdui/util/file"
)

type Reference struct {
	ItemCount      uint32
	ReferenceItems []*ReferenceItem
}

type ReferenceItem struct {
	OsTypeKey  string
	OsKeyBlock referenceOsKeyBlock
}

func (reference Reference) getOsKeyBlockID() string {
	return "obj "
}

func NewReference() *Reference {
	return new(Reference)
}

func (reference *Reference) Parse(file *os.File) {
	reference.ItemCount = util.ReadBytesLong(file)
	var i uint32
	for i = 0; i < reference.ItemCount; i++ {
		referenceItem := new(ReferenceItem)
		referenceItem.Parse(file)
		reference.ReferenceItems = append(reference.ReferenceItems, referenceItem)
	}
}

func (referenceItem *ReferenceItem) Parse(file *os.File) {
	referenceItem.OsTypeKey = util.ReadBytesString(file, 4)
	referenceItem.OsKeyBlock = parseReferenceOsKeyBlock(file, referenceItem.OsTypeKey)
}

type Bool struct {
	Value bool
}

func (Bool Bool) getOsKeyBlockID() string {
	return "bool"
}

func NewBool() *Bool {
	return new(Bool)
}

func (Bool *Bool) Parse(file *os.File) {
	Bool.Value = int(util.ReadSingleByte(file)) == 1
}

type Enum struct {
	Type  string
	Value string
}

func (enum Enum) getOsKeyBlockID() string {
	return "enum"
}

func NewEnum() *Enum {
	return new(Enum)
}

func (enum *Enum) Parse(file *os.File) {
	typeLength := util.ReadBytesLong(file)
	if typeLength < 1 {
		enum.Type = util.ReadBytesString(file, 4)
	} else {
		enum.Type = util.ReadBytesString(file, int(typeLength))
	}

	enumLength := util.ReadBytesLong(file)
	if enumLength < 1 {
		enum.Value = util.ReadBytesString(file, 4)
	} else {
		enum.Value = util.ReadBytesString(file, int(enumLength))
	}

}

type Text struct {
	Value string
}

func (text Text) getOsKeyBlockID() string {
	return "TEXT"
}

func NewText() *Text {
	return new(Text)
}

func (text *Text) Parse(file *os.File) {
	text.Value = util.ParseUnicodeString(file)
}

type Double struct {
	Value float64
}

func (double Double) getOsKeyBlockID() string {
	return "doub"
}

func NewDouble() *Double {
	return new(Double)
}

func (double *Double) Parse(file *os.File) error {
	value, err := util.ReadDouble(file)
	if err != nil {
		return err
	}
	double.Value = value
	return nil

}

type Unitfloat struct {
	UnitType string
	value    float64
}

func (unitFloat Unitfloat) getOsKeyBlockID() string {
	return "UntF"
}

func NewUnitFloat() *Unitfloat {
	return new(Unitfloat)
}

func (unitFloat *Unitfloat) Parse(file *os.File) error {
	unitFloat.UnitType = util.ReadBytesString(file, 4)
	double, err := util.ReadDouble(file)
	if err == nil {
		return err
	}
	unitFloat.value = double
	return nil
}

type Integer struct {
	Value uint32
}

func (integer Integer) getOsKeyBlockID() string {
	return "long"
}

func NewInteger() *Integer {
	return new(Integer)
}

func (integer *Integer) Parse(file *os.File) {
	integer.Value = util.ReadBytesLong(file)
}
