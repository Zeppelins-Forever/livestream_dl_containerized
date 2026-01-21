## This is here because I'm forgetful. Nothing here is useful to end users.


### Push to Docker Hub via:
```
docker buildx build --push \
            --platform linux/arm64/v8,linux/amd64 \
            --tag zeppelinsforever/livestream_dl_containerized:0.0.8 \
            --tag zeppelinsforever/livestream_dl_containerized:latest .
```
### Deno does not support 32bit ARM architectures or RISC-V in Alpine's repos, and thus cannot build to other architectures currently.
Check https://docs.docker.com/build/building/multi-platform/ for tools needed to release multiple architectures from one device.
Run ` docker run --privileged --rm tonistiigi/binfmt --install all ` to install Docker's QEMU tools for multi-arch builds.

### If you can't push, run
` docker login `
