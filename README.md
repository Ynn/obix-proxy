1. ./compile.sh
2.  set env : OBIX_SERVER_URL=http://localhost:7777/ and OBIX_NAME=my-obix
3. ./proxy
4. Example call :
```shell
yo@ada:~$ curl -l -H "Accept:application/json" 'http://localhost:7777/obix/org/lon/@XBureau922/@XLocations/@XLSI/@XSunSensor/@XnvoSunLux_2/$@CSNVT_lux/?extract=lux2'
{"lux2" : 1734}
yo@ada:~$ curl -l -H "Accept:line" 'http://localhost:7777/obix/org/lon/@XBureau922/@XLocations/@XLSI/@XSunSensor/@XnvoSunLux_2/$@CSNVT_lux/?extract=lux2'
my-obix lux2=1735 1510093803253292690
```

Example docker-compose-sample.yml


