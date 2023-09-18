package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type lineHandler func(header []byte, reader *bufio.Reader) *Payload

type Payload struct {
	Reply Reply
	Err   error
}

type Parser struct {
	lineHandlers map[byte]lineHandler
}

func NewParser() *Parser {
	p := Parser{}
	p.lineHandlers = map[byte]lineHandler{
		'+': p.parseSimpleString,
		'-': p.parseError,
		':': p.parseInt,
		'$': p.parseBulkString,
		'*': p.parseArray,
	}
	return &p
}

func (p *Parser) ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go p.parse(reader, ch)
	return ch
}

func (p *Parser) parse(rawReader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			// ...
		}
	}()

	reader := bufio.NewReader(rawReader)
	for {
		firstLine, err := reader.ReadBytes('\n')
		if err != nil {
			ch <- &Payload{
				Err: err,
			}
			close(ch)
			return
		}

		length := len(firstLine)
		if length <= 2 || firstLine[length-1] != '\n' || firstLine[length-2] != '\r' {
			continue
		}

		firstLine = bytes.TrimSuffix(firstLine, []byte{'\r', '\n'})
		lineHandler, ok := p.lineHandlers[firstLine[0]]
		if !ok {
			fmt.Println("not found line handler...")
			continue
		}

		payload := lineHandler(firstLine, reader)
		ch <- payload
		if payload.Err != nil {
			close(ch)
			return
		}
	}
}

func (p *Parser) parseSimpleString(header []byte, reader *bufio.Reader) *Payload {
	content := header[1:]
	return &Payload{
		Reply: NewSimpleStringReply(string(content)),
	}
}

func (p *Parser) parseInt(header []byte, reader *bufio.Reader) (payload *Payload) {
	var (
		err error
	)
	defer func() {
		if err != nil {
			payload = &Payload{
				Err: err,
			}
		}
	}()

	var i int64
	if i, err = strconv.ParseInt(string(header[1:]), 10, 64); err != nil {
		return
	}

	return &Payload{
		Reply: NewIntReply(i),
	}
}

func (p *Parser) parseError(header []byte, reader *bufio.Reader) *Payload {
	return &Payload{
		Reply: NewErrReply(string(header[1:])),
	}
}

func (p *Parser) parseBulkString(header []byte, reader *bufio.Reader) (payload *Payload) {
	body, err := p.parseBulkStringBody(header, reader)
	return &Payload{
		Reply: NewBulkReply(body),
		Err:   err,
	}
}

func (p *Parser) parseBulkStringBody(header []byte, reader *bufio.Reader) ([]byte, error) {
	strLen, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil {
		return nil, err
	}

	body := make([]byte, strLen+2)
	if _, err = io.ReadFull(reader, body); err != nil {
		return nil, err
	}

	return body, nil
}

func (p *Parser) parseArray(header []byte, reader *bufio.Reader) (payload *Payload) {
	var (
		err error
	)
	defer func() {
		if err != nil {
			payload = &Payload{
				Err: err,
			}
		}
	}()

	var nStrs int64
	if nStrs, err = strconv.ParseInt(string(header[1:]), 10, 64); err != nil {
		return
	}

	if nStrs <= 0 {
		return &Payload{
			Reply: NewEmptyMultiBulkReply(),
		}
	}

	lines := make([][]byte, 0, nStrs)
	for i := int64(0); i < nStrs; i++ {
		var firstLine []byte
		if firstLine, err = reader.ReadBytes('\n'); err != nil {
			return
		}

		length := len(firstLine)
		if length < 4 || firstLine[length-2] != '\r' || firstLine[length-1] != '\n' || firstLine[0] != '$' {
			continue
		}

		var batchStringBody []byte
		if batchStringBody, err = p.parseBulkStringBody(firstLine[:length-2], reader); err != nil {
			return
		}

		lines = append(lines, batchStringBody)
	}

	return &Payload{
		Reply: NewMultiBulkReply(lines),
	}
}
