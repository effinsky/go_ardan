package main

// GO Gotchas:
// 1 nil in Go actually has a type. and it assumes the type of pointer it needs
// when assigned (internally). nil in Go can be a "concretely typed value",
// var i *int, for instance. You can imagine it as if every value in the language
// instead of being just the singular value is actually a tuple (Type, Value).
// and then sometimes nil is just (nil, nil)

// 2 calling methods on nil pointers is fine, since they get promoted to a value
// of the type to which they point (they carry a pointer to the type). also
// methods are just sugar over functions that that the receiver as the first
// arg, really. so it's like a function that accepts a nil as the first arg --
// acceptable.

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

// 8 package-oriented design: packages should PROVIDE, not CONTAIN. one
// indication a package is providing is that a file inside the pkg should make
// sense to be called the name of the package. basically, pgks should not be
// dumps like common, utils, helpers, types.

// concurrency: executing code instructions out of order
// parallelism: doing multiple tasks at the same time

// thread levels for actual parallelism:
// hardware thread -> os thread -> goroutine

// opportunities for go scheduler context switching among goroutines :
// go keyword used; GC happening; syscalls; clocking calls in the program
// (mutexes etc)

// manage data races with :
// -- atomics for simple values like counters;
// -- mutexes for more complex data. consider rwmutex when necessary

// PANIC vs FATAL
// log.Panic allows deferred functions to execute and can be recovered from,
// while log.Fatal calls os.Exit and crashes the program irrecoverably and does
// not execute deferred functions.

// HANDLING ERRORS WITH DEFER FOR CLEANUPS, ROLLBACKS, AND SUCH:
//
// When a function encounters a return statement, the defer functions are
// called before the function actually returns but after the return values are
// computed. Ok, so you have the values computed for the returns. But here's
// the catch: the defer funcs cannot see that state of the values.
// if you're using a locally scoped error variable, the defer funcs will see
// the value of err at the point in the code where the defer func is defined,
// not the value of err that is being returned.
//
// On the other hand, if err is a named return, it will be in the scope of the
// entire function. So, the value of err inside the defer would be the value that
// is going to be returned by the function, because the defer function executes
// after the return values (named return variables) have been computed but
// before the function actually returns.

// GO flags
// Build and Installation:
//
// -o filename: Sets the output filename for the compiled program.
// -v: Enables verbose mode, printing additional information about the build process.
// -x: Prints the commands executed during the build process.
// -tags tag1 tag2...: Specifies build tags to include or exclude during compilation.
// Tags are defined in the source code and control which parts are compiled based on specific conditions.
// -buildmode mode: Sets the build mode (e.g., exe for executable, cgo for programs using C code).

// Optimization:
//
// -gcflags flags: Passes flags directly to the Go compiler. Some common flags for optimization include:
// -O: Enables various optimizations (levels from -O1 to -O3).
// -inline: Attempts to inline functions for performance improvement.
// -l: Disables inlining (useful for debugging).
// -linkflags flags: Passes flags directly to the linker. Some common flags include:
// -s: Omits symbol table and debug information for a smaller binary size.

// Testing and Debugging:
//
// -race: Enables data race detection to identify potential concurrency issues.
// -msan: Enables memory sanitizer to detect memory access errors.
// -cover: Enables coverage profiling to measure how much of your code is executed by tests.
// -debug: Enables debugging information in the executable for use with debuggers.

// Advanced Flags:
//
// -asmflags flags: Passes flags directly to the assembly generation stage.
// -workdir dir: Sets the working directory for the compiled program.
// -trimpath: Removes all information about the Go package path from the compiled binary.
// Escape Analysis and Benchmarking:
//
// -m: Used with go run or go build -gcflags to print escape analysis information.
// -bench=regexp: Selects benchmarks to run based on a regular expression with go run.
// -benchmem: Includes memory allocation statistics in benchmark output with go run.

// GO CONCURRENCY LIGHWEIGHT AND EFFICIENT HOW:
// Golang's concurrency model is lightweight compared to OS thread concurrency
// due to its use of goroutines and a highly optimized runtime scheduler. Here
// are the key factors that contribute to this:
//
// Goroutines:
//
// Lightweight: Goroutines (sized in kbs) are much smaller than OS threads (sized in mbs).
// A goroutine typically starts with a few kilobytes of stack space, which grows
// and shrinks as needed. In contrast, an OS thread usually starts with a fixed
// stack size, often in the megabytes.

// Fast Creation: Creating and destroying goroutines is much faster than
// creating and destroying OS threads. Goroutines are managed by the Go
// runtime, which has been optimized for performance. Goroutines are managed
// entirely in user space by the Go runtime, while OS threads involve kernel
// space operations. User space operations are generally faster because they
// don't require context switches to kernel mode.

// Scheduling

// The Go runtime uses a lightweight, cooperative scheduler for goroutines. This
// scheduler operates in user space and can make quick decisions about goroutine
// execution. OS thread scheduling involves the kernel scheduler, which has more
// overhead.

// Allocation Strategy

// The Go runtime pre-allocates and reuses goroutines from a pool, reducing the
// need for frequent memory allocation. OS thread creation often involves more
// complex resource allocation.
//

// Context Switching

// Switching between goroutines is faster than switching between OS threads
// because it doesn't involve a full context switch at the OS level (do not
// require kernel-mode transitions). While OS thread creation has been optimized
// over time, it still involves more heavyweight operations like allocating
// kernel resources, updating kernel data structures, and potentially
// interacting with hardware.

// M:N Scheduling:
//
// User-Space Scheduling: The Go runtime uses an M:N scheduler, where M
// goroutines are multiplexed onto N OS threads. This allows for efficient use
// of system resources and avoids the overhead associated with OS-level context
// switching.
//
// Non-blocking I/O and Syscalls:
//
// Asynchronous Operations: The Go runtime is designed to handle blocking
// operations efficiently. When a goroutine performs a blocking I/O operation
// or system call, the runtime can schedule other goroutines on the same
// thread, thus avoiding the need for a separate thread per blocking operation.
// While go runtime will often use a blocking op on the goroutine, it will use
// non-blocking syscalls for network ops etc under the hood.
//
// Garbage Collection and Memory Management (Language Integration):
//
// Go's garbage collector works in harmony with
// the goroutine scheduler, minimizing the impact of garbage collection pauses
// on goroutine scheduling. Goroutines are tightly integrated with Go's garbage
// collector and memory model, allowing for optimizations that might not be
// feasible with OS threads used across different languages and runtimes
//
// Work Stealing:
//
// Load Balancing: The Go scheduler uses a work-stealing algorithm to balance
// the load among available OS threads. If one OS thread becomes idle, it can
// steal goroutines from other threads, ensuring efficient CPU
// utilization.
//
// Scalability:
//
// High Concurrency: Because goroutines are so lightweight, a Go program can
// efficiently handle hundreds of thousands or even millions of concurrent
// tasks, which would be impractical with OS threads due to memory and
// performance constraints. By using these techniques, Go's concurrency model
// achieves high efficiency and scalability, making it a powerful tool for
// building concurrent applications.

// The performance advantage of goroutines is most noticeable when you're
// dealing with a large number of concurrent operations. For a small number of
// long-running, computationally intensive tasks, the difference might be less
// significant.

// Concurrency (out of order execution) vs Parallelism (simultaneous execution):
// Go achieves concurrency through goroutines, but true parallelism
// (simultaneous execution) only occurs when goroutines run on different OS
// threads on multiple CPU cores.

/// How to manage a read/write heavy database situation? ///
//
// -- read/write replicas
// -- indexing on the most often used fields
// -- partitioning
