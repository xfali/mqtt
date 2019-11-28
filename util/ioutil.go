// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package util

type CountWriter int64

func (w *CountWriter) Write(v []byte) (int, error) {
    l := len(v)
    *w = *w + CountWriter(l)
    return l, nil
}

func (w *CountWriter) Reset() {
    *w = 0
}

func (w *CountWriter) Count() int64 {
    return int64(*w)
}
