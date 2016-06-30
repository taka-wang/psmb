#goclient

PSMB test cases in golang

## Docker

### From the scratch
```bash
# build docker image 
docker build -t takawang/psmb-goclient .

# build arm version image 
#docker build -t takawang/arm-psmb-goclient -f Dockerfile.arm .


# mount file system
docker run -v /tmp:/tmp -it takawang/psmb-goclient /bin/bash

# run go test
go test -v

# Print app output
docker logs <container id>
# Enter the container
docker exec -it <container id> /bin/bash
```

### Pull pre-built docker image
```bash
docker pull takawang/psmb-goclient

# arm version
#docker pull takawang/arm-psmb-goclient
```
