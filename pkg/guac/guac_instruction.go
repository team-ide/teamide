package guac

import (
	"fmt"
	"strconv"
)

//The Guacamole protocol consists of instructions. Each instruction is a comma-delimited list followed by a terminating semicolon, where the first element of the list is the instruction opcode, and all following elements are the arguments for that instruction:
//
//OPCODE,ARG1,ARG2,ARG3,...;
//Each element of the list has a positive decimal integer length prefix separated by the value of the element by a period. This length denotes the number of Unicode characters in the value of the element, which is encoded in UTF-8:
//
//LENGTH.VALUE
//Any number of complete instructions make up a message which is sent from client to server or from server to client. Client to server instructions are generally control instructions (for connecting or disconnecting) and events (mouse and keyboard). Server to client instructions are generally drawing instructions (caching, clipping, drawing images), using the client as a remote display.
//
//For example, a complete and valid instruction for setting the display size to 1024x768 would be:
//
//4.size,1.0,4.1024,3.768;
//Here, the instruction would be decoded into four elements: "size", the opcode of the size instruction, "0", the index of the default layer, "1024", the desired width in pixels, and "768", the desired height in pixels.
//
//The structure of the Guacamole protocol is important as it allows the protocol to be streamed while also being easily parsable by JavaScript. JavaScript does have native support for conceptually-similar structures like XML or JSON, but neither of those formats is natively supported in a way that can be streamed; JavaScript requires the entirety of the XML or JSON message to be available at the time of decoding. The Guacamole protocol, on the other hand, can be parsed as it is received, and the presence of length prefixes within each instruction element means that the parser can quickly skip around from instruction to instruction without having to iterate over every character.

// Instruction represents a Guacamole instruction
type Instruction struct {
	Opcode string
	Args   []string
	cache  string
}

// NewInstruction creates an instruction
func NewInstruction(opcode string, args ...string) *Instruction {
	return &Instruction{
		Opcode: opcode,
		Args:   args,
	}
}

// String returns the on-wire representation of the instruction
func (i *Instruction) String() string {
	if len(i.cache) > 0 {
		return i.cache
	}

	i.cache = fmt.Sprintf("%d.%s", len(i.Opcode), i.Opcode)
	for _, value := range i.Args {
		i.cache += fmt.Sprintf(",%d.%s", len(value), value)
	}
	i.cache += ";"

	return i.cache
}

func (i *Instruction) Byte() []byte {
	return []byte(i.String())
}

//Parse 解析data 到 guacd instruction  todo:: 优化这个算法可以 提高net.io
func Parse(data []byte) (*Instruction, error) {
	elementStart := 0

	// Build list of elements
	elements := make([]string, 0, 1)
	for elementStart < len(data) {
		// Find end of length
		lengthEnd := -1
		for i := elementStart; i < len(data); i++ {
			if data[i] == '.' {
				lengthEnd = i
				break
			}
		}
		// read() is required to return a complete instruction. If it does
		// not, this is a severe internal error.
		if lengthEnd == -1 {
			return nil, ErrServer.NewError("ReadSome returned incomplete instruction.")
		}

		// Parse length
		length, e := strconv.Atoi(string(data[elementStart:lengthEnd]))
		if e != nil {
			return nil, ErrServer.NewError("ReadSome returned wrong pattern instruction.", e.Error())
		}

		// Parse element from just after period
		elementStart = lengthEnd + 1
		element := string(data[elementStart : elementStart+length])

		// Append element to list of elements
		elements = append(elements, element)

		// ReadSome terminator after element
		elementStart += length
		terminator := data[elementStart]

		// Continue reading instructions after terminator
		elementStart++

		// If we've reached the end of the instruction
		if terminator == ';' {
			break
		}

	}

	return NewInstruction(elements[0], elements[1:]...), nil
}

// ReadOne takes an instruction from the stream and parses it into an Instruction
func ReadOne(stream *Stream) (instruction *Instruction, err error) {
	var instructionBuffer []byte
	instructionBuffer, err = stream.ReadSome()
	if err != nil {
		return
	}

	return Parse(instructionBuffer)
}
