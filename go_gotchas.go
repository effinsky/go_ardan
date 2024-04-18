package main

// GO Gotchas:
// 1 nil in Go actually has a type. and it assumes the type of pointer it needs
// when assigned (internally). nil in Go can be a "concretely typed value",
// var i *int, for instance. You can imagine it as if every value in the language
// instead of being just the value is actually a tuple tuple (Type, Value).
// and then sometimes nil is just (nil, nil)

// 2 calling methods on nil pointers is fine, since they get promoted to a value
// of the type to which they point (they carry a pointer to the type)

// 3 types based on other types and values of those types passed to funcs/meths
// as base and derived types -- so a type based on another type is derived from it

// 4 implementing interfaces is only about the implementers having the right
// member funcs on them. compare vs 3.

// 5 returning concrete error types from func will create bugs when checking err
// vs nil!

// 6 we get npds when doing dot access on nil pointers because go automatically
// dereferences these pointers to values, and they do not point to any valid memory.

// 7 limits of immutability with value parameters: if we pass a value struct that
// contains a ptr to a struct to a func, we can mutate that inner struct thru
// the ptr. maybe we need ptr semantics on the outer too to make mutability clear.

// concurrency: executing code instructions out of order
// parallelism: doing multiple tasks at the same time

// thread levels for actual parallelism:
// hardware thread -> os thread -> goroutine

// opportunities for go scheduler context switching among goroutines :
// go keyword used; GC happening; syscalls; clocking calls in the program (mutexes etc)

// manage data races with : atomics for simple values like counters; mutexes for more
// complex data. consider rwmutex when necessary
