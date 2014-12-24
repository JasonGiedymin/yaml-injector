# yaml-injector TODO

## v0.1.0a

- [x] Add License file
- [x] Add Notice file
- [x] Add Godeps file
- [x] Add DataReader interface
- [x] Add YamlReader implementation
- [x] Add NewYamlReader constructor
- [x] Add JsonReader implementation #future
- [x] Add BaseData struct for embedded base data
- [x] Apply YamlReader implementation to tests and core
- [x] Introduce "go.crypto" for stdin terminal testing
- [x] Add debug output to test for stdin data #alpha
- [x] Modify README
  - [x] add LICENSE information
  - [x] add GoDeps information
- [x] Add gitignore
- [x] Add TODO file

## v0.2.0

- [x] Update docs
- [x] Modify MapData to map[string]interface[] from map[interface{}]interface to
   x  make using JSON easier
- [x] Add JsonData methods and Map() implementation
- [x] Rename inject param from `yaml_input` to `dest_file` to make it obvious
      the param is the destination data file.
- [x] Use JsonData when stdin is passed in to the command
- [ ] Add JSON stdin tests
- [ ] Add dot notation parser

## v0.5.0

- [ ] tag and build binary
