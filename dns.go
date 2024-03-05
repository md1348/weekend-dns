package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"math/rand"
	"net"
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

func buildHeader(buf *bytes.Buffer) (*DNSHeader, error) {
	header := make([]byte, 12)
	n, err := buf.Read(header)
	if n != 12 || err != nil {
		return nil, errors.New("invalid dns header - must be 12 bytes")
	}

	return &DNSHeader{
		id:             uint16(header[0])<<8 | uint16(header[1]),
		flags:          uint16(header[2])<<8 | uint16(header[3]),
		numQuestions:   uint16(header[4])<<8 | uint16(header[5]),
		numAnswers:     uint16(header[6])<<8 | uint16(header[7]),
		numAuthorities: uint16(header[8])<<8 | uint16(header[9]),
		numAdditionals: uint16(header[10])<<8 | uint16(header[11]),
	}, nil
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

type DNSRecord struct {
	name  string
	type_ uint16
	class uint16
	ttl   int
	data  string
}

func parseDomainSimple(buf *bytes.Buffer) (string, error) {
	parts := make([]byte, 0)

	for {
		length, err := buf.ReadByte()
		if err != nil {
			return "", errors.New("error reading domain name 1")
		}

		if length == 0 {
			break
		}

		if len(parts) > 0 {
			parts = append(parts, '.')
		}

		part := make([]byte, length)
		_, err = buf.Read(part)
		if err != nil {
			return "", errors.New("error reading domain name 2")
		}
		parts = append(parts, part...)
	}

	return string(parts), nil
}

// Assumes header has been parsed
func parseQuestion(buf *bytes.Buffer) (*DNSQuestion, error) {
	domain, err := parseDomainSimple(buf)
	if err != nil {
		return nil, errors.New("error parsing domain")
	}

	data := make([]byte, 4)
	buf.Read(data)
	type_ := uint16(data[0])<<8 | uint16(data[1])
	class := uint16(data[2])<<8 | uint16(data[3])

	return &DNSQuestion{
		name:  domain,
		type_: type_,
		class: class,
	}, nil
}

func main() {
	query := buildQuery("example.com", TYPE_A)

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("error writing")
		return
	}
	n, err := conn.Write(query)
	if err != nil {
		fmt.Println("error writing")
		return
	}
	fmt.Printf("wrote %d bytes\n", n)

	readBytes := make([]byte, 1024)
	n, err = conn.Read(readBytes)
	if err != nil {
		fmt.Println("error reading")
		return
	}
	fmt.Printf("read %d bytes\n", n)

	buf := bytes.NewBuffer(readBytes)
	header, _ := buildHeader(buf)
	question, _ := parseQuestion(buf)
	fmt.Println(header)
	fmt.Println(question)
}
