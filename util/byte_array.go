package util

import (
	"encoding/binary"
	"errors"
	"math"
)

var ErrTooLarge = errors.New("bytes.Buffer: too large")
var ErrOutOfRange = errors.New("out of range")
var ErrNotEnough = errors.New("Not enough momery")

type ByteArray struct {
	buf       []byte
	bootstrap [64]byte
	position  int
	length    int
	byteOrder binary.ByteOrder
}

func NewByteArray() *ByteArray {
	obj := &ByteArray{
		byteOrder: binary.LittleEndian,
	}

	return obj
}
func (this *ByteArray) SetPosition(position int) (err error) {
	if position >= len(this.buf) {
		err = ErrOutOfRange
	}
	this.position = position
	return
}
func (this *ByteArray) Position() int {
	return this.position
}
func (this *ByteArray) SetOrder(byteOrder binary.ByteOrder) {
	this.byteOrder = byteOrder
}
func (this *ByteArray) Length() int {
	return this.length
}
func (this *ByteArray) Bytes() []byte {
	return this.buf
}
func (this *ByteArray) checkRange(size int) (err error) {
	if (this.position + size) > len(this.buf) {
		err = ErrOutOfRange
	}
	return
}
func (this *ByteArray) ReadBoolean() (v bool, err error) {
	err = this.checkRange(1)
	if err != nil {
		return
	}
	if this.buf[this.position] > 0 {
		v = true
	} else {
		v = false
	}
	this.position++
	return
}
func (this *ByteArray) ReadByte() (v byte, err error) {
	err = this.checkRange(1)
	if err != nil {
		return
	}
	if this.position < len(this.buf) {
		v = this.buf[this.position]
	}
	this.position++
	return
}

//从this中读数据到other中
//offset this的偏移值
//length 最出的最大长度，== 0 读this从offset开始的所有数据
func (this *ByteArray) ReadByteArray(other *ByteArray, offset int, length int) (err error) {
	l := len(this.buf)
	if l == 0 {
		err = ErrOutOfRange
		return
	}
	if offset >= len(this.buf) {
		err = ErrOutOfRange
		return
	}

	if length == 0 {
		length = l - offset
	}
	this.position = offset
	for this.position < l && length > 0 {
		other.WriteByte(this.buf[this.position])
		length--
		this.position++
	}
	return
}
func (this *ByteArray) ReadBytes(buf []byte, offset, length int) (err error) {
	l := len(this.buf)
	if l == 0 {
		err = ErrOutOfRange
		return
	}
	if offset >= len(this.buf) {
		err = ErrOutOfRange
		return
	}

	if length == 0 ||
		length > l-offset {
		length = l - offset
	}

	if len(buf) < length {
		err = ErrNotEnough
		return
	}
	this.position = offset
	offset = 0
	for offset < l && length > 0 {
		buf[offset] = this.buf[this.position]
		this.position++
		offset++
		length--
	}
	this.position += length
	return
}
func (this *ByteArray) ReadDouble() (v float64, err error) {
	err = this.checkRange(8)
	if err != nil {
		return
	}
	u64 := this.byteOrder.Uint64(this.buf[this.position:])
	v = math.Float64frombits(u64)
	this.position += 8
	return
}

func (this *ByteArray) ReadFloat() (v float32, err error) {
	err = this.checkRange(4)
	if err != nil {
		return
	}
	u32 := this.byteOrder.Uint32(this.buf[this.position : this.position+4])
	v = math.Float32frombits(u32)
	this.position += 4
	return
}
func (this *ByteArray) ReadInt() (v int32, err error) {
	err = this.checkRange(4)
	if err != nil {
		return
	}
	v = int32(this.byteOrder.Uint32(this.buf[this.position : this.position+4]))
	this.position += 4
	return
}
func (this *ByteArray) ReadShort() (v int16, err error) {
	err = this.checkRange(2)
	if err != nil {
		return
	}
	v = int16(this.byteOrder.Uint16(this.buf[this.position : this.position+2]))
	this.position += 2
	return
}

func (this *ByteArray) ReadUnsignedByte() (v uint8, err error) {
	err = this.checkRange(1)
	if err != nil {
		return
	}
	v = uint8(this.buf[this.position])
	this.position++
	return
}

func (this *ByteArray) ReadUnsignedInt() (v uint32, err error) {
	vint, err := this.ReadInt()
	if err != nil {
		return
	}
	v = uint32(vint)
	return
}

func (this *ByteArray) ReadUnsignedShort() (v uint16, err error) {
	vshort, err := this.ReadShort()
	if err != nil {
		return
	}
	v = uint16(vshort)
	return
}

func (this *ByteArray) WriteBoolean(v bool) {
	var b byte
	if v {
		b = 1
	} else {
		b = 0
	}

	this.grow(1)
	this.buf[this.position] = b
	this.position++
	this.length = this.position
	return
}
func (this *ByteArray) WriteByte(v byte) {
	this.grow(1)
	this.buf[this.position] = v
	this.position++
	this.length = this.position
	return
}
func (this *ByteArray) WriteByteArray(other *ByteArray, offset int, length int) {
	ol := len(other.buf)
	if offset >= ol {
		offset = 0
	}

	if length == 0 {
		length = ol
	}
	other.position = offset
	for other.position < ol && length > 0 {
		this.grow(1)
		this.buf[this.position] = other.buf[other.position]
		this.position++
		other.position++
		length--
	}
	this.length = this.position
	return
}
func (this *ByteArray) WriteBytes(buf []byte, offset int, length int) {
	ol := len(buf)
	if offset >= ol {
		offset = 0
	}

	if length == 0 {
		length = ol
	}

	for offset < ol && length > 0 {
		this.grow(1)
		this.buf[this.position] = buf[offset]
		this.position++
		offset++
		length--
	}
	this.length = this.position
	return
}
func (this *ByteArray) WriteDouble(v float64) {
	u64 := math.Float64bits(v)
	this.grow(8)
	this.byteOrder.PutUint64(this.buf[this.position:], u64)
	this.position += 8
	this.length = this.position
	return
}

func (this *ByteArray) WriteFloat(v float32) {
	u32 := math.Float32bits(v)
	this.grow(4)
	this.byteOrder.PutUint32(this.buf[this.position:], u32)
	this.position += 4
	this.length = this.position
	return
}
func (this *ByteArray) WriteInt(v int32) {
	this.grow(4)
	this.byteOrder.PutUint32(this.buf[this.position:], uint32(v))
	this.position += 4
	this.length = this.position
	return
}
func (this *ByteArray) WriteShort(v int16) {
	this.grow(2)
	this.byteOrder.PutUint16(this.buf[this.position:], uint16(v))
	this.position += 2
	this.length = this.position
	return
}
func (this *ByteArray) WriteUnsignedInt(v uint32) {
	this.WriteInt(int32(v))
	return
}

//清除数据
func (this *ByteArray) Clear() {
	this.buf = this.buf[:0]
	this.position = 0
	this.length = 0
}

//剩余数据长度
func (this *ByteArray) BytesAvailable() int {
	n := this.length - this.position
	if n > 0 {
		return n
	}
	return 0
}

//将指定位置和长度的数据往头部挪
//lenght:== 0 时，长度是从offset开始到尾部的所有数据
//返回实际移动的长度，错误码
func (this *ByteArray) MoveToHead(offset, lenght int) (n int, err error) {
	if offset >= this.length {
		return 0, ErrOutOfRange
	}
	if offset+lenght >= this.length || lenght == 0 {
		n = this.length - offset
	} else {
		n = lenght
	}
	copy(this.buf, this.buf[offset:offset+n])
	this.position = 0
	this.length = n
	this.buf = this.buf[:n]

	return
}
func (b *ByteArray) grow(n int) {
	if b.position+n > cap(b.buf) {
		var buf []byte
		if b.buf == nil && n <= len(b.bootstrap) {
			buf = b.bootstrap[:n]
		} else {
			// not enough space anywhere
			buf = makeSlice(2*cap(b.buf) + n)
			copy(buf, b.buf[:b.position])
		}
		b.buf = buf
	}
	m := len(b.buf)
	if b.position+n >= m {
		b.buf = b.buf[:b.position+n]
	}
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}
