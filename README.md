# livestream_dl_containerized
A docker container and scripts for using livestream_dl conveniently!


Run with:
` sudo docker run -it --rm -v "$(pwd):/out" --user [UID]:[GID] zeppelinsforever/livestream_dl_containerized:latest [URL] ` or ` docker run -it --rm -v "$(pwd):/out" --user [UID]:[GID] zeppelinsforever/livestream_dl_containerized:latest [URL] ` depending on if you have a docker user set up. `-it` makes the container interactive, and `--rm` removes the container after being run. You wouldn't want a ton of containers hanging around, taking up space on your device, right?
If you don't know your UID and GID, you can quickly find that on a Linux system with the `id` command. The full command may look something like this:
` docker run  -it --rm -v "$(pwd):/out" --user 1000:1000 zeppelinsforever/livestream_dl_containerized:latest [URL] `

If you are on a system that doesn't have this information or doesn't respect Unix file permissions (like Windows, for example), the `--user [UID]:[GID]` section may be unnecessary, and can be removed. Further testing on DOcker running in WSL on Windows is needed...
All container operations run as root (as such, you shouldn't use this in any security-critical environment), and we run the container with the `--user` permissions so that file output is marked as being owned by yourself, rather than the "root" user - this could cause issues.

The image itself is at:
https://hub.docker.com/r/zeppelinsforever/livestream_dl_containerized

If you are using the version tagged `:latest`, you may need to occasionally update the package to the newest version. You can do this by running `docker pull zeppelinsforever/livestream_dl-container:latest`.

A Script to make all of these processes easier will be available soon-ish.

---

## To-Do:
- Pass cookies file for membership streams.
- Make script to interface with docker image. Clean up container after finished.
