# livestream_dl_containerized
A docker container and scripts for using livestream_dl conveniently!


Run with:
` sudo docker run -v --user [UID]:[GID] "$(pwd):/out" zeppelinsforever/livestream_dl_containerized:latest [URL] ` or ` docker run -v --user [UID]:[GID] "$(pwd):/out" zeppelinsforever/livestream_dl_containerized:latest [URL] ` depending on if you have a docker user set up.
If you don't know your UID and GID, you can quickly find that on a Linux system with the `id` command. The full command may look something like this:
` docker run -v --user 1000:1000 "$(pwd):/out" zeppelinsforever/livestream_dl_containerized:latest [URL] `

If you are on a system that doesn't have this information (like Windows, for example), the `--user [UID]:[GID]` section may be unnecessary, and can be removed.

The image itself is at:
https://hub.docker.com/r/zeppelinsforever/livestream_dl_containerized

A Script to make this process easier will be available soon-ish.

---

## To-Do:
- Make script to interface with docker image. Clean up container after finished.
