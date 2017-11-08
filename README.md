Example docker-compose.yml :
```yaml
version: '3'
services:
  obix-proxy:
    build: .
    image: nnynn/obix-proxy:latest
    ports:
      - 7777:8080
    environment:
      - OBIX_SERVER_URL=http://example.com:8080/
      - OBIX_NAME=my-obix
```

Example call :

```shell
yo@ada:~$ curl -l -H "Accept:application/json" 'http://localhost:7777/obix/org/lon/@XBureau922/@XLocations/@XLSI/@XSunSensor/@XnvoSunLux_2/$@CSNVT_lux/?extract=lux2'
{"lux2" : 1734}
yo@ada:~$ curl -l -H "Accept:line" 'http://localhost:7777/obix/org/lon/@XBureau922/@XLocations/@XLSI/@XSunSensor/@XnvoSunLux_2/$@CSNVT_lux/?extract=lux2'
my-obix lux2=1735 1510093803253292690
```