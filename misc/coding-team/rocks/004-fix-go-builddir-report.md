# Report: Remove invalid go-builddir field

## Summary
- Removed the `go-builddir: cmd/fancni` line from the `fancni-binary` part in `rockcraft.yaml`
- The Go plugin does not support this field; it runs `go install ./...` by default, which builds all packages including `cmd/fancni`
- The existing `organize` mapping (`bin/fancni: fancni`) remains correct

## Files changed
- `rockcraft.yaml`

## Notable tradeoffs or risks
None. The Go plugin's default behavior (`go install ./...`) builds the `cmd/fancni` binary and places it at `bin/fancni`, so no functional change results from removing the invalid field.
