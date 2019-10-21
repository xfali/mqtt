// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package binlist

import (
    "bytes"
    "container/list"
)

type BinNode struct {
    k string
    //start
    s int
    //end
    e int
}

type BinMetaList list.List

func NewMetaList() *BinMetaList {
    ret := &BinMetaList{}
    (*list.List)(ret).Init()
    return ret
}

type BinList struct {
    D    []byte
    list *BinMetaList
}

func New(l *BinMetaList) *BinList {
    return &BinList{
        list: l,
    }
}

func (bm *BinMetaList) Put(key string, s, size int) *BinMetaList {
    n := &BinNode{
        k: key,
        s: s,
        e: s + size,
    }

    (*list.List)(bm).PushBack(n)
    return bm
}

func (bm *BinMetaList) Get(key string) *BinNode {
    for f := (*list.List)(bm).Front(); f != nil; f = f.Next() {
        n := f.Value.(*BinNode)
        if n.k == key {
            return n
        }
    }
    return nil
}

func (bm *BinList) Get(key string) []byte {
    n := bm.list.Get(key)
    if n != nil {
        return bm.D[n.s: n.e]
    }
    return nil
}

type BinListBuilder struct {
    buffer bytes.Buffer
    l      *BinMetaList
}

func NewBuilder() *BinListBuilder {
    return &BinListBuilder{
        l: NewMetaList(),
    }
}

func (b *BinListBuilder) Append(key string, d []byte) *BinListBuilder {
    cur := b.buffer.Len()
    b.l.Put(key, cur, len(d))

    b.buffer.Write(d)
    return b
}

func (b *BinListBuilder) Build() *BinList {
    l := BinList{
        D:    b.buffer.Bytes(),
        list: b.l,
    }

    return &l
}
