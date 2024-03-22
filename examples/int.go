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

	// This code should not error. ifの条件式によって、cmp.Orに置き換えられないため
	shouldNotErr := 4
	if g := retIntPointer(); g != nil {
		if isTrue() {
			shouldNotErr = *g
		}
	}

	// This code should not error. cmp.Orでの置き換え先もゼロ値であり、置き換えられないため
	var shouldNotErr2 *int
	if h := retIntPointer(); h != nil {
		shouldNotErr2 = h
	}

	// This code should not error. cmp.Orでの置き換え先もゼロ値であり、置き換えられないため
	var shouldNotErr3 int
	if i := retInt(); i != 0 {
		shouldNotErr3 = i
	}

	// This code should not error. cmp.Orでの置き換え先もゼロ値であり、置き換えられないため
	shouldNotErr4 := 0
	if j := retInt(); j != 0 {
		shouldNotErr4 = j
	}

	log.Println(shouldErr)
	log.Println(shouldNotErr)
	log.Println(shouldNotErr2)
	log.Println(shouldNotErr3)
	log.Println(shouldNotErr4)
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
