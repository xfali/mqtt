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

//4 - MQTT
var MqttString = []byte{0, 4, 77, 81, 84, 84}

type String struct {
    length uint16
    data   []byte
}

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
        return n+n2, err2
    }

    return n+n2, nil
}

func ParseString(r io.ByteReader) (*String, error) {
    var buf []byte
    var s, size uint16
    for i := 0; ; i++ {
        b, err := r.ReadByte()
        if err != nil {
            return &String{length: size, data: buf}, err
        }
        if i == 0 {
            size = uint16(b << 8)
            continue
        } else if i == 1 {
            size |= uint16(b)
            buf = make([]byte, size)
            continue
        }

        if size <= s {
            break
        }
        buf[s] = b
        s++
    }
    return &String{length: size, data: buf}, nil
}

func (s *String) Length() uint16 {
    return s.length
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
