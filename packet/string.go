// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package packet

import (
    "io"
    "math"
    "mqtt/errcode"
)

type String struct {
    length uint16
    data   []byte
}

type StringPair [2]String

type Bytes String

func NewString(s string) (*String, error) {
    if len(s) > math.MaxUint16 {
        return nil, errcode.StringOutOfRange
    }
    ret := String{
        length: uint16(len(s)),
        data:   []byte(s),
    }
    return &ret, nil
}

func FromString(s string) (String, error) {
    if len(s) > math.MaxUint16 {
        return String{}, errcode.StringOutOfRange
    }
    ret := String{
        length: uint16(len(s)),
        data:   []byte(s),
    }
    return ret, nil
}

func NewStringPair(s1, s2 string) (StringPair, error) {
    if len(s1) > math.MaxUint16 || len(s2) > math.MaxUint16 {
        return [2]String{}, errcode.StringOutOfRange
    }
    ret := [2]String{
        String{
            length: uint16(len(s1)),
            data:   []byte(s1),
        },
        String{
            length: uint16(len(s2)),
            data:   []byte(s2),
        },
    }
    return ret, nil
}

func (p StringPair) Equals(other StringPair) bool {
    return p[0].String() == other[0].String() && p[1].String() == other[1].String()
}

func EncodeString(w io.Writer, s string) (int, error) {
    if len(s) > math.MaxUint16 {
        return 0, errcode.StringOutOfRange
    }
    b := make([]byte, 2)
    l := len(s)
    b[0] = byte(l >> 8)
    b[1] = byte(0xFF & l)
    n, err := w.Write(b)
    if err != nil {
        return n, err
    }
    n2, err2 := w.Write([]byte(s))
    if err != nil {
        return n + n2, err2
    }

    return n + n2, nil
}

func WriteBytes(w io.Writer, s Bytes) (int, error) {
    return WriteString(w, String(s))
}

func WriteString(w io.Writer, s String) (int, error) {
    b := make([]byte, 2)
    l := s.length
    b[0] = byte(l >> 8)
    b[1] = byte(0xFF & l)
    n, err := w.Write(b)
    if err != nil {
        return n, err
    }
    n2, err2 := w.Write([]byte(s.data))
    if err != nil {
        return n + n2, err2
    }

    return n + n2, nil
}

func ParseBytes(r io.Reader) (ret *Bytes, n int, err error) {
    s, n, e := ParseString(r)
    return (*Bytes)(s), n, e
}

func ParseString(r io.Reader) (ret *String, n int, err error) {
    header := make([]byte, 1)
    var size uint16
    readSize := 0
    n, err = r.Read(header)
    if err != nil {
        return nil, n, err
    }
    readSize += n
    size = uint16(header[0] << 8)
    n, err = r.Read(header)
    if err != nil {
        return nil, readSize + n, err
    }
    readSize += n
    size |= uint16(header[0])

    buf := make([]byte, size)
    n, err = r.Read(buf)
    if err != nil {
        return nil, readSize + n, err
    }
    readSize += n

    return &String{length: size, data: buf}, readSize, nil
}

func (s *Bytes) Reset(v []byte) error {
    if len(v) > math.MaxUint16 {
        return errcode.StringOutOfRange
    }
    s.length = uint16(len(v))
    s.data = v
    return nil
}

func (s *String) Reset(v string) error {
    if len(v) > math.MaxUint16 {
        return errcode.StringOutOfRange
    }
    s.length = uint16(len(v))
    s.data = []byte(v)
    return nil
}

func (s *String) Length() uint16 {
    return s.length
}

func (s *String) AllLength() uint16 {
    return s.length + 2
}

func (s *StringPair) String() string {
    return "[0]: " + s[0].String() + " [1]: " + s[1].String()
}

func (s *String) String() string {
    return string(s.data)
}

type StringReader struct {
    cur int
    S   String
}

func (r *StringReader) Reset(s String) {
    r.cur = 0
    r.S = s
}

func (r *StringReader) Read(d []byte) (int, error) {
    n := 0
    length := len(d)
    for i := 0; i < length; i++ {
        if r.cur == 0 {
            d[i] = byte(r.S.length >> 8)
        } else if r.cur == 1 {
            d[i] = byte(r.S.length & 0xFF)
        } else {
            read := r.cur - 2
            if read == int(r.S.length) {
                return n, io.EOF
            }
            d[i] = r.S.data[r.cur-2]
        }
        r.cur++
        n++
    }
    return n, nil
}
