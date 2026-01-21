#!/bin/sh
source myenv/bin/activate

#########################
#python3 runner.py --log-level DEBUG --output "/out/[%(upload_date)s] %(title)s (%(id)s).%(ext)s"  --resolution best $1
#python3 runner.py --log-level DEBUG --embed-thumbnail --wait-for-video 60 --output "/out/[%(upload_date)s] %(title)s (%(id)s)"  --resolution best $1
#python3 runner.py --log-level DEBUG --embed-thumbnail --wait-for-video 60 --live-chat --output "/out/[%(upload_date)s] %(title)s (%(id)s"  --resolution best $1

#  If I want the container to automatically use cookies if /cookies/cookies.txt exists:
#COOKIE_FILE="/cookies/cookies.txt"
#if [ -f "$COOKIE_FILE" ]; then
#    # If the cookies file exists, inject the --cookies flag before the rest of the arguments
#    echo "Found /cookies/cookies.txt, applying to runner..."
#    python3 runner.py --cookies "$COOKIE_FILE" "$@"
#else
#    # Otherwise, run normally
#    python3 runner.py "$@"
#fi
#########################

# Automatically passes all arguments to runner.py (not including arguments for Docker itself, such as -v "/mounted/filesystem/")
echo "$@"
python3 runner.py --output "/temp/[%(upload_date)s] %(title)s (%(id)s)" "$@"
# At some point, add custom user argument to allow customizing file name (though filter certain characters, preventing user from changing file location accidentally)

# Necessary currently since certain functions in the container need Root to run, but this conflicts with
# permissions outside the contamv iner. Therefore, need intermediary storage location inside the container.
#chmod +x /temp/*
chown -R ${MY_UID}:${MY_GID} /temp/*
mv /temp/* /out/
