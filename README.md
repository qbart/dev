# dev

## utils

```
# generate UUID v4
dev rand uuid4 

# generate random bytes encoded with base64
dev rand bytes           
dev rand bytes --size 32 (default is 64)

# manage known hosts
dev ssh known-hosts reset
dev ssh known-hosts diff
```

## go

```
# generate empty go app
dev go main > main.go
```
