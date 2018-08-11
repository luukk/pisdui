package imageresource

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabulousduck/pisdui/pisdui/util"
)

/*Data contains the resource blocks
used by the photoshop file and the length of the
section in the photoshop file*/
type Data struct {
	Length         uint32
	ResourceBlocks []ResourceBlock
}

/*ResourceBlock contains the raw unparsed data from
a resource block in the photoshop file*/
type ResourceBlock struct {
	byteSize            uint32
	Signature           string
	ID                  uint16
	PascalString        string
	DataSize            uint32
	DataBlock           []byte
	ParsedResourceBlock parsedResourceBlock
}

/*NewData creates a new ImageResources struct
and returns a pointer to it.
This exists so the top level pisdui struct can create one
to prevent import cycles*/
func NewData() *Data {
	return new(Data)
}

/*Parse will read all image resources located in
the photoshop file and will read them into the ImageResources struct*/
func (ir *Data) Parse(file *os.File) {
	ir.Length = util.ReadBytesLong(file)
	var i uint32
	for i = 0; i < ir.Length; {
		block := ir.parseResourceBlock(file)
		ir.ResourceBlocks = append(ir.ResourceBlocks, *block)
		spew.Dump(ir)
		i += block.byteSize
	}
}

func (ir *Data) parseResourceBlock(file *os.File) *ResourceBlock {
	readByteCount := 0

	block := new(ResourceBlock)
	block.Signature = util.ReadBytesString(file, 4)
	readByteCount += 4

	block.ID = util.ReadBytesShort(file)
	readByteCount += 2

	pascalString, stringLength := ir.parsePascalString(file)
	readByteCount += stringLength

	block.PascalString = pascalString
	block.DataSize = util.ReadBytesLong(file)
	readByteCount += 4

	block.DataBlock = util.ReadBytesNInt(file, block.DataSize)
	readByteCount += int(block.DataSize)

	if block.DataSize%2 != 0 {
		util.ReadSingleByte(file)
		readByteCount += 1
	}

	block.byteSize = uint32(readByteCount)
	return block
}

func (ir *Data) parsePascalString(file *os.File) (string, int) {
	b := util.ReadSingleByte(file)
	if b == 0 {
		util.ReadSingleByte(file)
		return "", 1
	}

	s := util.ReadBytesString(file, b)

	if b%2 != 0 {
		util.ReadSingleByte(file)
	}
	return s, len(s)
}
