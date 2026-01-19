# livestream_dl_containerized
A docker container and scripts for using livestream_dl (semi)conveniently!

# Easiest Usage:
Make sure Docker Engine is running on your machine (and that the user you're running as can directly run Docker containers, otherwise this may not work properly) and run `archive-helper.sh`. It will guide you through a set of "sensible defaults" for downloading a stream. (this is not done yet, but will be soon)

## Running livestream_dl_containerized directly:
### Quirks:
- If you want to pass cookies to the container, I recommend using `-v /full/path/to/my_cookies.txt:/cookies/cookies.txt` as an argument when directly launching the container via "docker run". Replace "/full/path/to/my_cookies.txt" with your actual (not relative) system path to your cookies file. The container has a folder to put it in (`/cookies`) and the above arguments will place it in there as `cookies.txt`. Also, pass the argument `--cookies /cookies/cookies.txt` after the container name, so the container knows where to find the mounted cookies file.
- You currently cannot use `--output`, as this is relied upon for certain functions within the container itself. Functionality which lets you cusomize the output name, without adjusting where it's output to, may come later.

Refer to https://github.com/CanOfSocks/livestream_dl?tab=readme-ov-file#modification-of-yt-dlp for a full list of commands.

### Recommended commands:

(Linux/MacOS)> Download video to current directory:
> `docker run -it --rm -v "$(pwd):/out" -e MY_UID=$(id -u) -e MY_GID=$(id -g) zeppelinsforever/livestream_dl_containerized:latest --log-level DEBUG --wait-for-video 60 --live-chat --resolution best [URL]`

(Windows)> Download video to current directory (UID and GID are arbitrary, since Windows doesn't use Unix-like permissions):
> `docker run -it --rm -v "C:\full\system\path\to-folder:/out" -e MY_UID=1000 -e MY_GID=1000 zeppelinsforever/livestream_dl_containerized:latest --log-level DEBUG --wait-for-video 60 --live-chat --resolution best [URL]`

(Linux/MacOS)> Download video to current directory, using cookies - for accessing membership content:
> `docker run -it --rm -v "$(pwd):/out" -e MY_UID=$(id -u) -e MY_GID=$(id -g) -v /FULL/path/to/your_cookies.txt:/cookies/cookies.txt zeppelinsforever/livestream_dl_containerized:latest --log-level DEBUG --cookies /cookies/cookies.txt --wait-for-video 60 --live-chat --resolution best [URL]`

(Windows)> Download video to current directory, using cookies - for accessing membership content:
> Not yet tested

## Required arguments:
You MUST run your container with the following arguments:
`docker run -it --rm -v "$(pwd):/out" -e MY_UID=$(id -u) -e MY_GID=$(id -g) [image name] <arguments> [URL]`

---

The image itself is at:
https://hub.docker.com/r/zeppelinsforever/livestream_dl_containerized

If you are using the version tagged `:latest`, you may need to occasionally update the package to the newest version. You can do this by running `docker pull zeppelinsforever/livestream_dl-container:latest`.

A Script to make all of these processes easier will be available soon-ish.

---

## To-Do:
- Add output name cusomization, without allowing user to change directory (as this would break container functionality)
- Make script to interface with docker image.
