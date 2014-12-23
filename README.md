# Yaml Injector

Yaml Injector v0.1.0 - Jason Giedymin


## Installation

To Install:
```bash
# TODO
```

## Usage


## Development

Depencencies:

```bash
gpm install
```

Testing:

```bash
go test -v
```

Running:

```bash
go run yaml-injector.go \
  --debug \
  --file test/input.yaml \
  --using test/data.yaml \
  --key "b[0]" \
  inject into "a"
```


## License

yaml-injector is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.


## Todo

- [x] stdin/out/err
- [ ] key parser
- [ ] *nix data pipe in
