// test pattern int
package main

import "log"

func Int() {
	// This code should error.
	shouldErr := 1
	if a := retIntPointer(); a != nil {
		shouldErr = *a
	}

	// This code should error.
	b := retIntPointer()
	if b != nil {
		shouldErr = *b
	}

	// This code should error.
	if c := retInt(); c != 0 {
		shouldErr = c
	}

	// This code should error.
	if d := retIntPointer(); d != nil {
		if e := retInt(); e != 0 {
			shouldErr = *d + e
		}
	}

	// This code should not error. ブロック内で宣言した変数を再代入しているため
	if f := retInt(); f != 0 {
		var shouldNotErr int //lint:ignore S1021 This is a test
		shouldNotErr = f
		log.Println(shouldNotErr)
	}

	// This code should not error. cmp.Orに置き換えられないため
	shouldNotErr := 4
	if g := retIntPointer(); g != nil {
		if isTrue() {
			shouldNotErr = *g
		}
	}

	log.Println(shouldErr)
	log.Println(shouldNotErr)
}

func retIntPointer() *int {
	res := 2
	return &res
}

func retInt() int {
	res := 3
	return res
}

func isTrue() bool {
	return true
}
