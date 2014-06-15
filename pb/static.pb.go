// Code generated by protoc-gen-gogo.
// source: static.proto
// DO NOT EDIT!

/*
	Package pb is a generated protocol buffer package.

	It is generated from these files:
		static.proto

	It has these top-level messages:
		StaticId
*/
package pb

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// discarding unused import gogoproto "gogoproto/gogo.pb"

import io "io"
import code_google_com_p_gogoprotobuf_proto "code.google.com/p/gogoprotobuf/proto"

import fmt "fmt"
import strings "strings"
import reflect "reflect"

import fmt1 "fmt"
import strings1 "strings"
import code_google_com_p_gogoprotobuf_proto1 "code.google.com/p/gogoprotobuf/proto"
import sort "sort"
import strconv "strconv"
import reflect1 "reflect"

import bytes "bytes"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type StaticId struct {
	Hash             []byte `protobuf:"bytes,1,req,name=hash" json:"hash"`
	Length           int64  `protobuf:"varint,2,req,name=length" json:"length"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *StaticId) Reset()      { *m = StaticId{} }
func (*StaticId) ProtoMessage() {}

func (m *StaticId) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *StaticId) GetLength() int64 {
	if m != nil {
		return m.Length
	}
	return 0
}

func init() {
}
func (m *StaticId) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return code_google_com_p_gogoprotobuf_proto.ErrWrongType
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash, data[index:postIndex]...)
			index = postIndex
		case 2:
			if wireType != 0 {
				return code_google_com_p_gogoprotobuf_proto.ErrWrongType
			}
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				m.Length |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := code_google_com_p_gogoprotobuf_proto.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (this *StaticId) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&StaticId{`,
		`Hash:` + fmt.Sprintf("%v", this.Hash) + `,`,
		`Length:` + fmt.Sprintf("%v", this.Length) + `,`,
		`XXX_unrecognized:` + fmt.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringStatic(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *StaticId) Size() (n int) {
	var l int
	_ = l
	l = len(m.Hash)
	n += 1 + l + sovStatic(uint64(l))
	n += 1 + sovStatic(uint64(m.Length))
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovStatic(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozStatic(x uint64) (n int) {
	return sovStatic(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StaticId) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *StaticId) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	data[i] = 0xa
	i++
	i = encodeVarintStatic(data, i, uint64(len(m.Hash)))
	i += copy(data[i:], m.Hash)
	data[i] = 0x10
	i++
	i = encodeVarintStatic(data, i, uint64(m.Length))
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func encodeFixed64Static(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Static(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintStatic(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (this *StaticId) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings1.Join([]string{`&pb.StaticId{` + `Hash:` + fmt1.Sprintf("%#v", this.Hash), `Length:` + fmt1.Sprintf("%#v", this.Length), `XXX_unrecognized:` + fmt1.Sprintf("%#v", this.XXX_unrecognized) + `}`}, ", ")
	return s
}
func valueToGoStringStatic(v interface{}, typ string) string {
	rv := reflect1.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect1.Indirect(rv).Interface()
	return fmt1.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func extensionToGoStringStatic(e map[int32]code_google_com_p_gogoprotobuf_proto1.Extension) string {
	if e == nil {
		return "nil"
	}
	s := "map[int32]proto.Extension{"
	keys := make([]int, 0, len(e))
	for k := range e {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	ss := []string{}
	for _, k := range keys {
		ss = append(ss, strconv.Itoa(k)+": "+e[int32(k)].GoString())
	}
	s += strings1.Join(ss, ",") + "}"
	return s
}
func (this *StaticId) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*StaticId)
	if !ok {
		return false
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Hash, that1.Hash) {
		return false
	}
	if this.Length != that1.Length {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
