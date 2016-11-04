# Dev 

## Build
```
$ docker-compose run app bash
$ env GOOS=linux go build -v github.com/lukin0110/push/cmd/push/
$ ./push
```

or:

```
$ docker-compose run app bash
$ env GOOS=linux go install -v github.com/lukin0110/push/cmd/push/
$ push
```

or :

```
$ docker-compose run app bash
$ env GOOS=darwin go install -v github.com/lukin0110/push/cmd/push/
$ push
```
