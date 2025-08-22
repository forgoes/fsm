# Introduction
## How to run
### Clone the repository
```shell
git clone git@github.com:forgoes/fsm.git
cd fsm
```
### Initialize Go module
```shell
go mod tidy -v
```
### Run tests
Run all tests with verbose output:
```shell
make test
```
Run tests with the race detector:
```shell
make race
```
### Check coverage
```shell
make cover
```
### Run benchmarks
```shell
make bench
```
### Clean up
```shell
make clean
```

## What’s inside
```markdown
fsm/
├── fsm.go             # Core FSM implementation (states, transitions, guards, actions)
├── fsm_test.go        # Unit tests for the FSM core
├── mod_three.go        # Example FSM: binary string modulo 3 calculator
├── mod_three_test.go   # Unit tests for the modThreeFSM example
├── go.mod             # Go module definition
└── Makefile           # Helper commands for testing, coverage, and benchmarks
```

### Key components

**FSM (fsm.go)**
Provides a thread-safe finite state machine implementation with support for guards and actions.

**modThreeFSM (mod_three.go)**
Demonstrates how to use the FSM to compute the remainder of a binary string modulo 3.

**Tests (fsm_test.go, mod_three_test.go)**
Cover normal transitions, guards, actions, concurrency safety, and example FSM behavior.

**Makefile**
Shortcuts for running tests, checking race conditions, generating coverage reports, and running benchmarks.

# Design
The core FSM (fsm.go) is placed at the root for simple third-party import.
The mod_three implementation is only an example; it normally wouldn’t belong in a library, but is kept here and tested for demonstration.

## Why concurrency matters here

When working in the Go environment, I believe concurrency is something we always need to keep in mind.
Go makes it easy to run things in parallel, so without proper concurrency handling, even a simple FSM could behave incorrectly when accessed by multiple goroutines.

In a distributed environment, the local sync.RWMutex could be replaced or avoided depending on requirements:

1. **Distributed lock**: e.g. use etcd or Consul to coordinate FSM state changes across nodes.

2. **Database lock**: rely on row-level locks or advisory locks in a relational database (e.g. SELECT … FOR UPDATE in PostgreSQL, or MySQL’s GET_LOCK()), ensuring serialized access.

3. **Hash-based ordering**: map a unique ID (e.g. Kafka partition key) to a deterministic processing order, avoiding explicit locks.

# Motivation
The reason behind this FSM exercise is actually connected to distributed systems. Here’s the big picture:

### CAP theorem
In distributed systems, partition tolerance is a must.
That means we always have to choose between consistency (C) and availability (A).

### What web services usually do
Most web services don’t go for strong consistency, because it hurts availability.
Unless we shrink the scale with something like Raft + hashing, strong consistency is often not practical.
Instead, we normally accept eventual consistency, often through message queues.

### The duplicate message problem
Message queues guarantee that messages aren’t lost, but retries mean duplicates will happen.
That’s where unique IDs and FSMs come in: they help us make sure that even with duplicate events, the system state remains correct.


# Future Work
1. Scalability – Support larger or hierarchical FSMs beyond small fixed states.
2. Persistence – Store and recover FSM state (e.g., database, Kafka) for reliability.
3. Observability – Add logging, metrics, and visualization to aid debugging.

