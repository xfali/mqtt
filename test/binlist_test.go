// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "mqtt/container/binlist"
    "testing"
)

func TestBinList(t *testing.T) {
    ml := binlist.NewMetaList()
    l := binlist.New(ml)
    l.D = []byte{
        1, 2, 3, 4,
    }

    ml.Put("1", 0, 1).Put("2", 1, 1).Put("3", 2, 1).Put("4", 3, 1)
    t.Log(l.Get("1"))
    t.Log(l.Get("2"))
    t.Log(l.Get("3"))
    t.Log(l.Get("4"))
    t.Log(l.Get("5"))
    ml.Put("5", 0, 4)
    t.Log(l.Get("5"))
}

func TestBinList2(t *testing.T) {
    builder := binlist.NewBuilder()
    builder.Append("1", []byte{1})
    builder.Append("2", []byte{2})
    builder.Append("3", []byte{3})
    builder.Append("4", []byte{4})
    l := builder.Build()

    t.Log(l.Get("1"))
    t.Log(l.Get("2"))
    t.Log(l.Get("3"))
    t.Log(l.Get("4"))
    t.Log(l.Get("5"))
}
