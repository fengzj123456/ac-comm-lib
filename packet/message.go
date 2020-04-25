package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	fmt "fmt"
	"hash/crc32"
	"io"
	math "math"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/ironzhang/pearls/endian"
)

const (
	checksumLen = 4
)

var (
	MaxHeadLen = math.MaxUint8
	MaxBodyLen = math.MaxUint16
)

var (
	ErrInvalidChecksum  = errors.New("invalid checksum")
	ErrInvalidHeadLen   = errors.New("invalid head length")
	ErrInvalidBodyLen   = errors.New("invalid body length")
	ErrInvalidPacketLen = errors.New("invalid packet length")
)

/*
Packet format
+----------+
| Length   | Varint64, 1 - 10 bytes
+----------+
| Payload  | [Length-4]byte, Length-4 bytes
+----------+
| Checksum | uint32, BigEndian, 4 bytes
+----------+

Payload format
+---------+
| HeadLen | Varint64, 1 - 10 bytes
+---------+
| Head    | pb.Head
+---------+
| Body    | proto.Message
+---------+
*/

type Message struct {
	Head Head
	Body proto.Message
}

func NewMessage(body proto.Message) *Message {
	return &Message{
		Head: Head{Name: proto.MessageName(body)},
		Body: body,
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("Head: { %s}, Body: { %s}", m.Head.String(), m.Body.String())
}

func (m *Message) Encode() ([]byte, error) {
	head, err := proto.Marshal(&m.Head)
	if err != nil {
		return nil, err
	}
	hlen := len(head)
	if hlen > MaxHeadLen {
		return nil, ErrInvalidHeadLen
	}
	body, err := proto.Marshal(m.Body)
	if err != nil {
		return nil, err
	}
	if len(body) > MaxBodyLen {
		return nil, ErrInvalidBodyLen
	}

	var buf bytes.Buffer
	if err = endian.EncodeVarint(&buf, int64(hlen)); err != nil {
		return nil, err
	}
	if _, err = buf.Write(head); err != nil {
		return nil, err
	}
	if _, err = buf.Write(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) Decode(data []byte) (err error) {
	hlen, n := binary.Varint(data)
	if n <= 0 || hlen <= 0 || int(hlen) > MaxHeadLen || int(hlen) > len(data)-n {
		return ErrInvalidHeadLen
	}

	hend := n + int(hlen)
	if err = proto.Unmarshal(data[n:hend], &m.Head); err != nil {
		return err
	}
	if m.Body, err = newBody(m.Head.Name); err != nil {
		return err
	}
	if err = proto.Unmarshal(data[hend:], m.Body); err != nil {
		return err
	}
	return nil
}

func newBody(name string) (proto.Message, error) {
	typ := proto.MessageType(name)
	if typ == nil || typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("unknown message name: %s", name)
	}
	return reflect.New(typ.Elem()).Interface().(proto.Message), nil
}

func maxPacketLen() int {
	return binary.MaxVarintLen64 + MaxHeadLen + MaxBodyLen + checksumLen
}

func Write(w io.Writer, m *Message) error {
	payload, err := m.Encode()
	if err != nil {
		return err
	}

	length := len(payload) + checksumLen
	if length > maxPacketLen() {
		return ErrInvalidPacketLen
	}
	checksum := crc32.ChecksumIEEE(payload)

	if err = endian.EncodeVarint(w, int64(length)); err != nil {
		return err
	}
	if _, err = w.Write(payload); err != nil {
		return err
	}
	if err = endian.BigEndian.EncodeUint32(w, checksum); err != nil {
		return err
	}
	return nil
}

func Read(r io.Reader) (*Message, error) {
	length, err := endian.DecodeVarint(r)
	if err != nil {
		return nil, err
	}
	if length-checksumLen <= 0 || int(length) > maxPacketLen() {
		return nil, ErrInvalidPacketLen
	}
	buf := make([]byte, length)
	if _, err = io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	payload := buf[:length-checksumLen]
	checksum := binary.BigEndian.Uint32(buf[length-checksumLen:])
	if checksum != crc32.ChecksumIEEE(payload) {
		return nil, ErrInvalidChecksum
	}

	var m Message
	if err = m.Decode(payload); err != nil {
		return nil, err
	}
	return &m, nil
}
