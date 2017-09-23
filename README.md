# Novel

```
docker run --rm --name groovy-worker \
 -v $PWD:/workspace \
 hyper.cd/core/groovy:2.4.5 \
 bash -c 'cd /workspace; groovy -c utf-8 bot.groovy 6'
```

```
vi zh.sh
./zh.sh
```

```console
go get -v github.com/gregjones/httpcache
# modify github.com/gregjones/httpcache/httpcache.go, make `getFreshness` always returns `fresh`
go get -v github.com/peterbourgon/diskv
go build
./novel 1
```