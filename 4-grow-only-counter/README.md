### Grow-only counter 

#### How do I run this?

```bash
# Build and install.
go install .

# Start tester.
maelstrom test -w g-counter \
  --bin bin/maelstrom-gocounter \
  --node-count 3 \
  --rate 100 \
  --time-limit 20 \
  --nemesis partition
```
