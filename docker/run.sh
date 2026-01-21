#!/bin/sh

# 1. Defaults to UID & GID of 1000 if not provided by user.
TARGET_UID=${MY_UID:-1000}
TARGET_GID=${MY_GID:-1000}

# 2. Setup User and Group
if ! getent group appgroup >/dev/null 2>&1; then
    groupadd -g "$TARGET_GID" appgroup
fi

if ! id -u appuser >/dev/null 2>&1; then
    useradd -u "$TARGET_UID" -g "$TARGET_GID" -m -s /bin/sh appuser
fi

# 3. Fix Permissions
# We chown /cookies so the app can create lockfiles there if needed
chown -R "$TARGET_UID:$TARGET_GID" /livestream_dl-main /cookies

# 4. Check for User-Provided Output Flag
# We loop through arguments to see if '--output' was passed
HAS_CUSTOM_OUTPUT="false"
for arg in "$@"; do
    if [ "$arg" = "--output" ]; then
        HAS_CUSTOM_OUTPUT="true"
        break
    fi
done

# 5. Execute
# If user provided --output, we pass "$@" directly (User MUST include "/out/" at the start of their file path.
# Ex. "--output /out/filename" (do not include file extension) or "--output /out/dir1/dir2/filename"
# If user does not use "--output" flag, we prepend our default --output template
if [ "$HAS_CUSTOM_OUTPUT" = "true" ]; then
    echo "Custom output detected. Using user arguments..."
    exec su-exec "$TARGET_UID:$TARGET_GID" \
        /livestream_dl-main/myenv/bin/python3 runner.py \
        "$@"
else
    echo "No output flag detected. Using default /out/ template..."
    exec su-exec "$TARGET_UID:$TARGET_GID" \
        /livestream_dl-main/myenv/bin/python3 runner.py \
        --output "/out/[%(upload_date)s] %(title)s (%(id)s)" \
        "$@"
fi
