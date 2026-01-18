# This is here because I'm forgetful. Nothing here is useful to end users.


### Push to Docker Hub via:
```
sudo docker buildx build --push \
            --platform linux/amd64 \
            --tag zeppelinsforever/livestream_dl_containerized:0.0.1 \
            --tag zeppelinsforever/livestream_dl_containerized:latest .
```
### Can't use other architectures because they won't build. Ignore for now, investigate later.

### If you can't push, run
` docker login `
