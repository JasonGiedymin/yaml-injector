# Yaml Injector

Yaml Injector v0.1.0 - Jason Giedymin


## Installation

To Install:
```bash
# TODO
```

## Usage

Running (if installed):

```bash
yaml-injector \
  --debug \
  --test \
  --file test/input.yaml \
  --using test/data.yaml \
  --key "b[0]" \
  inject into "a"
```

See the [Makefile](Makefile) for expanded example usage.


## Development

Depencencies:

```bash
gpm install
```

### Testing:

Uses a makefile. Call upon the unit-tests make task:

```bash
make unit-tests
```

### Commands

Current commands are:

```
-----------------------Commands:----------------------
 help            : Test help listing 
 benchmark       : Run benchmark 
 unit-tests      : Test app 
 test-datafile   : Test with datafile to stdout 
 test-preview    : Test preview functionality to stdout 
 test-stdin      : Test with JSON stdin to stdout 
------------------------------------------------------
```

## Benchmarks

```
Comments:
  - 1ms == 1000000ns
```

### 12/31/2014
```
            Test               Runs           Time/Op         Mem/Op          Mem/Stats
BenchmarkMapInPlaceSelector   500000        3431 ns/op       546 B/op       19 allocs/op
BenchmarkGetKey              5000000         532 ns/op        66 B/op        3 allocs/op

Comments:
  - 3431ns per operation == 0.003431 ms
  -  532ns per operation == 0.000532 ms
```

### 1/1/2015
```
            Test               Runs           Time/Op         Mem/Op          Mem/Stats
BenchmarkMapInPlaceSelector   500000        3013 ns/op       550 B/op       16 allocs/op
BenchmarkGetValue            5000000         532 ns/op        67 B/op        3 allocs/op
Comments:
  - 3013ns per operation == 0.003013 ms
```

## License

yaml-injector is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.


## Todo

See [TODO](TODO) file for details.
