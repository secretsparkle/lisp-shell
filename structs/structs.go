package structs

import (
	"fmt"
)

type Control struct {
	Args []string
	Body []string
}

type Function struct {
	Name     string
	Args     []string
	Bindings map[string]string
	Body     List
}

type Element struct {
	next *Element
	prev *Element
	Data interface{}
}

type List struct {
	Head   *Element
	Tail   *Element
	Length int
}

func (e *Element) Next() *Element {
	return e.next
}

func (e *Element) Prev() *Element {
	return e.prev
}

func (l *List) Len() int {
	return l.Length
}

func (l *List) PushFront(v interface{}) *List {
	e := new(Element)
	e.Data = v
	l.Length++
	if l.Head == nil && l.Tail == nil {
		l.Head = e
	} else if l.Head != nil && l.Tail == nil {
		l.Tail = l.Head
		l.Head = e
		l.Head.next = l.Tail
		l.Tail.prev = l.Head
	} else if l.Head == nil && l.Tail != nil {
		l.Head = e
		l.Head.next = l.Tail
		l.Tail.prev = l.Head
	} else {
		e.next = l.Head
		l.Head.prev = e
		l.Head = e
	}
	return l
}

func (l *List) PushBack(v interface{}) *List {
	e := new(Element)
	e.Data = v
	l.Length++
	if l.Head == nil && l.Tail == nil {
		l.Head = e
	} else if l.Head != nil && l.Tail == nil {
		l.Tail = e
		l.Tail.prev = l.Head
		l.Head.next = l.Tail
	} else if l.Head == nil && l.Tail != nil {
		l.Head = l.Tail
		l.Tail = e
		l.Tail.prev = l.Head
		l.Head.next = l.Tail
	} else if l.Tail == nil {
		l.Head.next = e
		e.prev = l.Head
		l.Tail = e
	} else {
		l.Tail.next = e
		e.prev = l.Tail
		l.Tail = e
	}
	return l
}

func PrintList(l List) {
	for e := l.Head; e != nil; e = e.Next() {
		switch e.Data.(type) {
		case float64:
			if e.Prev() == nil {
				fmt.Print(e.Data)
			} else {
				fmt.Print(" ", e.Data)
			}
		case string:
			if e.Prev() == nil {
				fmt.Print(e.Data)
			} else {
				fmt.Print(" ", e.Data)
			}
		default:
			if e.Prev() != nil {
				fmt.Print(" ")
			}
			fmt.Print("(")
			PrintList(e.Data.(List))
			fmt.Print(")")
		}
	}
}
