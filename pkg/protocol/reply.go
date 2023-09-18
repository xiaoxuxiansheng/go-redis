package protocol

import (
	"strconv"
	"strings"
)

// CRLF is the line separator of redis serialization protocol
const CRLF = "\r\n"

type Reply interface {
	ToBytes() []byte
}

type SimpleStringReply struct {
	Str string
}

func NewSimpleStringReply(str string) *SimpleStringReply {
	return &SimpleStringReply{
		Str: str,
	}
}

func (s *SimpleStringReply) ToBytes() []byte {
	return []byte("+" + s.Str + CRLF)
}

type IntReply struct {
	Code int64
}

func NewIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

type ErrReply struct {
	ErrStr string
}

func NewErrReply(errStr string) *ErrReply {
	return &ErrReply{
		ErrStr: errStr,
	}
}

func (e *ErrReply) ToBytes() []byte {
	return []byte("-" + e.ErrStr + CRLF)
}

var (
	nullReply     = &NullReply{}
	nullBulkBytes = []byte("$-1\r\n")
)

type NullReply struct {
}

func NewNullReply() *NullReply {
	return nullReply
}

func (n *NullReply) ToBytes() []byte {
	return nullBulkBytes
}

type BulkReply struct {
	Arg []byte
}

func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

func (b *BulkReply) ToBytes() []byte {
	if b.Arg == nil {
		return nullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

type MultiBulkReply struct {
	Args [][]byte
}

func NewMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

func (m *MultiBulkReply) ToBytes() []byte {
	var strBuf strings.Builder
	strBuf.WriteString("*" + strconv.Itoa(len(m.Args)) + CRLF)
	for _, arg := range m.Args {
		if arg == nil {
			strBuf.WriteString("$-1" + CRLF)
			continue
		}
		strBuf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
	}

	return nil
}

var emptyMultiBulkBytes = []byte("*0\r\n")

// EmptyMultiBulkReply is a empty list
type EmptyMultiBulkReply struct{}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// ToBytes marshal redis.Reply
func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}
