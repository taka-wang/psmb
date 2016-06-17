# mb
modbus reactor: [modbusd + modbus slave simulator](https://github.com/taka-wang/modbusd)

# Docker

## build image
```bash
docker build -t takawang/mbd .
```

## run image
```bash
docker run -v /tmp:/tmp -itd --name=mbd takawang/mbd
```
