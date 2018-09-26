# JumpCloud 

**PWHASH** - The Password Hashing Service

## Installing

- `go get github.com/connormckelvey/jumpcloud/cmd/pwhash`

## Running

- `pwhash` (Start with defaults)
- `pwhash -h`

```
Usage of pwhash:
  -d int
    	Amount to delay hash response in seconds (default 5)
  -p int
    	Port to listen on (default 8080)
```

## Testing

- `cd $GOPATH/src/github.com/connormckelvey/jumpcloud`
- `make test`

