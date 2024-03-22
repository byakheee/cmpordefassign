// test pattern struct
package main

import "log"

type Form struct {
	String string
	Int    int
}

type Entity struct {
	String string
	Int    int
}

func Struct() {
	// This code should not error. cmp.Orに置き換えられないため
	form := &Form{
		String: "one",
		Int:    1,
	}
	if a := retStructPointer(); a != nil {
		form.String = a.String
		form.Int = a.Int
	}

	log.Println(form)
}

func retStructPointer() *Entity {
	res := &Entity{
		String: "two",
		Int:    2,
	}
	return res
}
