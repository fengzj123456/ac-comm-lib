package packet

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	proto "github.com/golang/protobuf/proto"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		m    *Message
		name string
	}{
		{
			m:    NewMessage(&TheTestMsg{}),
			name: "packet.TheTestMsg",
		},
		{
			m:    NewMessage(&Head{}),
			name: "packet.Head",
		},
	}
	for i, tt := range tests {
		if got, want := tt.m.Head.Name, tt.name; got != want {
			t.Errorf("%d: message name: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: message name: got %v", i, got)
		}
	}
}

func TestNewBody(t *testing.T) {
	tests := []struct {
		name string
		msg  proto.Message
		err  bool
	}{
		{"Head", nil, true},
		{"packet.Head", &Head{}, false},
		{"packet.TheTestMsg", &TheTestMsg{}, false},
	}
	for i, tt := range tests {
		msg, err := newBody(tt.name)
		if tt.err {
			if err == nil {
				t.Fatalf("%d: new body: err is nil", i)
			} else {
				t.Logf("%d: new body: err(%v) is not nil", i, err)
			}
		} else {
			if err != nil {
				t.Fatalf("%d: new body: %v", i, err)
			}
			if got, want := reflect.TypeOf(msg), reflect.TypeOf(tt.msg); got != want {
				t.Fatalf("%d: message type: got %v, want %v", i, got, want)
			} else {
				t.Logf("%d: message type: got %v", i, got)
			}
		}
	}
}

func TestMessageString(t *testing.T) {
	var w io.Writer
	w = os.Stdout
	w = ioutil.Discard
	msgs := []Message{
		{
			Head: Head{Name: "packet.TheTestMsg"},
			Body: &TheTestMsg{},
		},
		{
			Head: Head{Name: "packet.TheTestMsg"},
			Body: &TheTestMsg{A: 1, B: 2, C: "hello"},
		},
	}
	for i, m := range msgs {
		fmt.Fprintf(w, "%d: message string: %v\n", i, m.String())
	}
}

func messageEqual(m1, m2 *Message) bool {
	if !proto.Equal(&m1.Head, &m2.Head) {
		return false
	}
	if !proto.Equal(m1.Body, m2.Body) {
		return false
	}
	return true
}

func TestMessage(t *testing.T) {
	tests := []struct {
		msg Message
		err bool
	}{
		{
			msg: Message{
				Head: Head{Name: "packet.TheTestMsg"},
				Body: &TheTestMsg{},
			},
			err: false,
		},
		{
			msg: Message{
				Head: Head{Name: "packet.TheTestMsg"},
				Body: &TheTestMsg{A: 1, B: 2, C: "hello"},
			},
			err: false,
		},
		{
			msg: Message{
				Head: Head{Name: "packet.TheTestMsgXX"},
				Body: &TheTestMsg{},
			},
			err: true,
		},
	}
	for i, tt := range tests {
		data, err := tt.msg.Encode()
		if err != nil {
			t.Fatalf("%d: message encode: %v", i, err)
		}

		var m Message
		err = m.Decode(data)
		if tt.err {
			if err == nil {
				t.Fatalf("%d: decode: err is nil", i)
			} else {
				t.Logf("%d: decode: err(%v) is not nil", i, err)
			}
		} else {
			if err != nil {
				t.Fatalf("%d: decode: %v", i, err)
			}
			if got, want := &m, &tt.msg; !messageEqual(got, want) {
				t.Fatalf("%d: message: got %v, want %v", i, got, want)
			} else {
				t.Logf("%d: message: got %v", i, got)
			}
		}
	}
}

func TestReadWrite(t *testing.T) {
	msgs := []Message{
		{
			Head: Head{Name: "packet.TheTestMsg"},
			Body: &TheTestMsg{},
		},
		{
			Head: Head{Name: "packet.TheTestMsg"},
			Body: &TheTestMsg{A: 1, B: 2, C: "hello"},
		},
	}

	var buf bytes.Buffer
	for i, msg := range msgs {
		if err := Write(&buf, &msg); err != nil {
			t.Fatalf("%d: write: %v", i, err)
		}
	}
	t.Logf("bytes: %d: %x", buf.Len(), buf.Bytes())
	for i, msg := range msgs {
		m, err := Read(&buf)
		if err != nil {
			t.Fatalf("%d: read: %v", i, err)
		}
		if got, want := m, &msg; !messageEqual(got, want) {
			t.Fatalf("%d: message: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: message: got %v", i, got)
		}
	}
}
