package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const TYPE_A = 1
const CLASS_IN = 1

type DNSHeader struct {
	id             uint16
	flags          uint16
	numQuestions   uint16
	numAnswers     uint16
	numAuthorities uint16
	numAdditionals uint16
}

func (d *DNSHeader) toBytes() []byte {
	output := make([]byte, 0)
	output = binary.BigEndian.AppendUint16(output, d.id)
	output = binary.BigEndian.AppendUint16(output, d.flags)
	output = binary.BigEndian.AppendUint16(output, d.numQuestions)
	output = binary.BigEndian.AppendUint16(output, d.numAnswers)
	output = binary.BigEndian.AppendUint16(output, d.numAuthorities)
	output = binary.BigEndian.AppendUint16(output, d.numAdditionals)
	return output
}

type DNSQuestion struct {
	name  string
	type_ uint16
	class uint16
}

func (d *DNSQuestion) encodeName() []byte {
	output := make([]byte, 0)
	parts := strings.Split(d.name, ".")
	for _, part := range parts {
		output = binary.BigEndian.AppendUint16(output, uint16(len(part)))
		output = append(output, []byte(part)...)
	}
	output = binary.BigEndian.AppendUint16(output, 0)
	return output
}

func (d *DNSQuestion) toBytes() []byte {
	output := d.encodeName()
	output = binary.BigEndian.AppendUint16(output, d.type_)
	output = binary.BigEndian.AppendUint16(output, d.class)
	return output
}

func buildQuery(domain string, recordType uint16) []byte {
	output := make([]byte, 0)
	return output
}

func main() {
	h := DNSQuestion{"google.com", 1, 513}
	fmt.Println(h.toBytes())
	return
}
