# livestream_dl_containerized
A docker container and scripts for using CanOfSocks's [livestream_dl](https://github.com/CanOfSocks/livestream_dl) (semi)conveniently! All you need is Docker!
This is a tool to help download actively running YouTube livestreams

# Easiest Usage (optional):
Make sure Docker Engine is running on your machine, then download [archive-helper](https://github.com/Zeppelins-Forever/livestream_dl_containerized/releases) (and optionally add it to your PATH) and run `archive-helper` (Linux/MacOS) or `archive-helper.exe` (Windows) in the terminal. It will run a set of "sensible defaults" for downloading a stream. It will ask you for a URL and the full path to a cookies file (only use cookies if you are downloading membership content).

Archive-helper checks if Docker is installed and running, pulls the newest version of [zeppelinsforever/livestream_dl_containerized](https://hub.docker.com/r/zeppelinsforever/livestream_dl_containerized), elevates the docker call if needed (i.e. if your user must run docker containers as sudo), and runs the container with commansd similar to those in the "Recommended Commands" section below.

It can be ran the following ways:
- `archive-helper [--silent] [--cookies "/full/path/to/cookies.txt"] [URL]`
- `archive-helper [URL]`
- `archive-helper`

| Option | Configuration |
| --- | --- |
| [--silent] | Redirects Stdout and Stderr (all terminal output) to a file called nogup.{date}.out |
| [--cookies "/path/to/cookies.txt"] | Accepts a cookies file so you can download a members-only stream. If you need help exporting cookies from your browser, I recommend downloading [this browser extension](https://github.com/rotemdan/ExportCookies), logging into YouTube in an incognito tab (to avoid cookie reuse), and exporting them via the extension. <b>Important:</b> Make sure you use the FULL system path, no user-specific path or relative path. |
| [URL | Just paste the URL of the livestream you want to download. |


# Running livestream_dl_containerized directly:
### Quirks:
- If you want to pass cookies to the container, I recommend using `-v /full/path/to/my_cookies.txt:/cookies/cookies.txt` as an argument when directly launching the container via "docker run". Replace "/full/path/to/my_cookies.txt" with your actual (not relative) system path to your cookies file. The container has a folder to put it in (`/cookies`) and the above arguments will place it in there as `cookies.txt`. Also, pass the argument `--cookies /cookies/cookies.txt` after the container name, so the container knows where to find the mounted cookies file within the container.
- You currently cannot use `--output`, as this is relied upon for certain functions within the container itself (well, you can use it, but it will break stuff). Functionality which lets you customize the output name, without adjusting where it's output to, may come later. For instance, for the container to work properly with varying permissions on different systems, the intermediary files are written to a directory inside the container, and are moved out when done.
- Note: If you are running these docker commands directly within Linux's `nohup` (i.e. `nohup docker run ... --resolution best [URL] &`), always exclude the `-it` command. Since nohup runs the command without interactivity (and streams the output to a "nohup.out" file by default) and in another process, you will not have the interacitvity that `-it` requires, and it will fail.
Otherwise, it runs almost exactly the same as traditional livestream_dl.

Refer to https://github.com/CanOfSocks/livestream_dl?tab=readme-ov-file#modification-of-yt-dlp for a full list of commands.

### Recommended commands:

(Linux/MacOS)> Download video to current directory:
> `docker run -it --rm -v "$(pwd):/out" -e MY_UID=$(id -u) -e MY_GID=$(id -g) zeppelinsforever/livestream_dl_containerized:latest --log-level DEBUG --wait-for-video 60 --live-chat --resolution best [URL]`

(Windows)> Download video to current directory (UID and GID are arbitrary, since Windows doesn't use Unix-like permissions):
> `docker run -it --rm -v "C:\full\system\path\to-folder:/out" -e MY_UID=1000 -e MY_GID=1000 zeppelinsforever/livestream_dl_containerized:latest --log-level DEBUG --wait-for-video 60 --live-chat --resolution best [URL]`

(Linux/MacOS)> Download video to current directory, using cookies - for accessing membership content:
> `docker run -it --rm -v "$(pwd):/out" -e MY_UID=$(id -u) -e MY_GID=$(id -g) -v /FULL/path/to/your_cookies.txt:/cookies/cookies.txt zeppelinsforever/livestream_dl_containerized:latest --log-level DEBUG --cookies /cookies/cookies.txt --wait-for-video 60 --live-chat --resolution best [URL]`

(Windows)> Download video to current directory, using cookies - for accessing membership content (UID and GID are arbitrary, since Windows doesn't use Unix-like permissions):
> `docker run -it --rm -v "C:\full\system\path\to-folder:/out" -e MY_UID=1000 -e MY_GID=1000 -v C:\full\path\to\your_cookies.txt:/cookies/cookies.txt zeppelinsforever/livestream_dl_containerized:latest --log-level DEBUG --cookies /cookies/cookies.txt --wait-for-video 60 --live-chat --resolution best [URL]`

## Required Arguments:
You MUST run your container with the following arguments:
`docker run -v "$(pwd):/out" -e MY_UID=$(id -u) -e MY_GID=$(id -g) <image name> [arguments] <URL>`
If you're going to run the container in such a way that you will manually feed in parameters that it asks for, make sure you run docker with the `-it` flag for interactivity.

---

The image itself is at:
https://hub.docker.com/r/zeppelinsforever/livestream_dl_containerized

If you are using the version tagged `:latest`, you may need to occasionally update the package to the newest version. You can do this by running `docker pull zeppelinsforever/livestream_dl-container:latest`.

## To-Do:
- Add output name cusomization, without allowing user to change directory (as this would break container functionality)
- Make script to interface with docker image.
