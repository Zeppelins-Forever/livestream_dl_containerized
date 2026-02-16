## This is here because I'm forgetful. Nothing here is useful to end users.


### Push to Docker Hub via:
```
docker buildx build --no-cache --push \
            --platform linux/amd64 \
            --tag zeppelinsforever/livestream_dl_containerized:0.0.13 \
            --tag zeppelinsforever/livestream_dl_containerized:latest .
```
### Due to issues with APK tool on Alpine ARM64 builds, make sure `qemu-user-static` and `qemu-user-static-binfmt` are installed.
Run `sudo pacman -S qemu-user-static qemu-user-static-binfmt` to fix on Arch.

Check https://docs.docker.com/build/building/multi-platform/ for tools needed to release multiple architectures from one device.
Run ` docker run --privileged --rm tonistiigi/binfmt --install all ` to install Docker's QEMU tools for multi-arch builds.

### If you can't push, run
` docker login `
