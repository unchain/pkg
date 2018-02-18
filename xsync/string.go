package xsync

import "sync"

type String struct {
	stringMut sync.RWMutex
	string    string
}

func (t *String) Get() string {
	t.stringMut.RLock()
	defer t.stringMut.RUnlock()

	return t.string
}

func (t *String) Set(string string) {
	t.stringMut.Lock()
	defer t.stringMut.Unlock()

	t.string = string
}

func NewString() *String {
	return &String{
		string: "",
	}
}
