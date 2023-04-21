### Unique ID Generation

#### How do I run this?

```bash
# Build and install.
go install .

# Run tester.
maelstrom test -w unique-ids \
  --bin bin/maelstrom-uidgen \
  --time-limit 30 \
  --rate 1000 \
  --node-count 3 \
  --availability total \
  --nemesis partition
```
