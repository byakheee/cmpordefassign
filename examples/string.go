// test pattern string
package main

import "log"

func String() {
	// This code should error.
	shouldErr := "one"
	if v := retStringPointer(); v != nil {
		shouldErr = *v
	}

	// This code should error.
	w := retStringPointer()
	if w != nil {
		shouldErr = *w
	}

	// This code should error.
	if x := retString(); x != "" {
		shouldErr = x
	}

	// This code should error.
	if x := retStringPointer(); x != nil {
		if y := retString(); y != "" {
			shouldErr = *x + y
		}
	}

	// This code should not error. ブロック内で宣言した変数を再代入しているため
	if y := retString(); y != "" {
		var shouldNotErr string //lint:ignore S1021 This is a test
		shouldNotErr = y
		log.Println(shouldNotErr)
	}

	// This code should not error. cmp.Orに置き換えられないため
	shouldNotErr := "four"
	if x := retStringPointer(); x != nil {
		if isTrue() {
			shouldNotErr = *x
		}
	}

	log.Println(shouldErr)
	log.Println(shouldNotErr)
}

func retStringPointer() *string {
	res := "two"
	return &res
}

func retString() string {
	res := "three"
	return res
}
