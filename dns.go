package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"math/rand"
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
		bitSize := bits.Len(uint(len(part)))
		byteSize := (bitSize + 7) / 8
		for i := 0; i < byteSize; i++ {
			output = append(output, byte(uint(len(part))>>(i*8)))
		}
		output = append(output, []byte(part)...)
	}
	output = append(output, byte(0))
	return output
}

func (d *DNSQuestion) toBytes() []byte {
	output := d.encodeName()
	output = binary.BigEndian.AppendUint16(output, d.type_)
	output = binary.BigEndian.AppendUint16(output, d.class)
	return output
}

func buildQuery(domain string, recordType uint16) []byte {
	id := uint16(rand.Int())
	RECURSION_DESIRED := uint16(1 << 8)
	header := DNSHeader{id: id, numQuestions: 1, flags: RECURSION_DESIRED}
	question := DNSQuestion{name: domain, type_: recordType, class: CLASS_IN}
	output := header.toBytes()
	output = append(output, question.toBytes()...)
	return output
}

func main() {
	h := DNSQuestion{"google.com", 1, 513}
	// query := h.toBytes()
	query := h.encodeName()

	// query := buildQuery("example.com", TYPE_A)

	for _, b := range query {
		fmt.Printf("0x%x ", b)
	}

	// conn, err := net.Dial("udp", "8.8.8.8:53")
	// if err != nil {
	// 	fmt.Println("error writing")
	// 	return
	// }
	// n, err := conn.Write(query)
	// if err != nil {
	// 	fmt.Println("error writing")
	// 	return
	// }
	// fmt.Printf("wrote %d bytes\n", n)
	//
	// buffer := make([]byte, 1024)
	// n, err = conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("error reading")
	// 	return
	// }
	// fmt.Printf("read %d bytes\n", n)
	//
	// fmt.Println(buffer[:n])
	return
}
