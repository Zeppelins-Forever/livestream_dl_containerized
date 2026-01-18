# livestream_dl_containerized
A docker container and scripts for using livestream_dl conveniently!


Run with:
` sudo docker run -v "$(pwd):/out" --user [UID]:[GID] zeppelinsforever/livestream_dl_containerized:latest [URL] ` or ` docker run -v "$(pwd):/out" --user [UID]:[GID] zeppelinsforever/livestream_dl_containerized:latest [URL] ` depending on if you have a docker user set up.
If you don't know your UID and GID, you can quickly find that on a Linux system with the `id` command. The full command may look something like this:
` docker run -v "$(pwd):/out" --user 1000:1000 zeppelinsforever/livestream_dl_containerized:latest [URL] `

If you are on a system that doesn't have this information or doesn't respect Unix file permissions (like Windows, for example), the `--user [UID]:[GID]` section may be unnecessary, and can be removed. Further testing on DOcker running in WSL on Windows is needed...
All container operations run as root (as such, you shouldn't use this in any security-critical environment), and we run the container with the `--user` permissions so that file output is marked as being owned by yourself, rather than the "root" user - this could cause issues.

The image itself is at:
https://hub.docker.com/r/zeppelinsforever/livestream_dl_containerized

A Script to make this process easier will be available soon-ish.

---

## To-Do:
- Pass cookies file for membership streams.
- Make script to interface with docker image. Clean up container after finished.
