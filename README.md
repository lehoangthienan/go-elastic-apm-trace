# EXAMPLE APM TRACE & MONITORING

![jaeger-go-trace](https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt136299afed87309d/5c98d669da7491ed59827e6c/APM_BlogThumbnail.png)

## SETUP

```
make setup
make init
```

## HOW TO START
```
make dev-svc-a
make dev-svc-b
curl http://localhost:3000/go
```


## APM CLIENT
```
http://localhost:5601/app/apm/services?rangeFrom=now-10m&rangeTo=now&environment=
```