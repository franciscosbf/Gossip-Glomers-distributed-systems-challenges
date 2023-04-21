### Echo

#### How do I run this?

```bash
# Build and install.
go install .

# Run tester.
maelstrom test -w echo \
  --bin bin/maelstrom-echo \
  --node-count 1 \
  --time-limit 10
```
