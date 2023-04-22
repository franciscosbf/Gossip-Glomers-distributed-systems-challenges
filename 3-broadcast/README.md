### Unique ID Generation

#### How do I run this?

```bash
# Build and install.
go install .

# Start tester 3a.
maelstrom test -w broadcast \
  --bin bin/maelstrom-broadcast \
  --node-count 1 \
  --time-limit 20 \
  --rate 10

# Start tester 3b.
maelstrom test -w broadcast \
  --bin bin/maelstrom-broadcast \
  --node-count 5 \
  --time-limit 20 \
  --rate 10

# Start tester 3c.
maelstrom test -w broadcast \
  --bin bin/maelstrom-broadcast \
  --node-count 5 \
  --time-limit 20 \
  --rate 10 \
  --nemesis partition

# Start tester 3d Part I and II.
maelstrom test -w broadcast \
  --bin bin/maelstrom-broadcast \
  --node-count 25 \
  --time-limit 20 \
  --rate 100 \
  --latency 100
```
