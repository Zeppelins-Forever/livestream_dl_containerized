echo "Run this in the same directory as the Go livestream_dl_containerized project for building archive-helper."
mkdir Releases

read -n 1 -s -r -p "Make sure the .go file is set to output INFO logs.\nPress any key to continue..."
echo
# INFO Builds - change the
GOOS=linux GOARCH=amd64 go build -o Releases/archive-helper-INFO-Linux-x64
GOOS=linux GOARCH=arm64 go build -o Releases/archive-helper-INFO-Linux-ARM64
GOOS=windows GOARCH=amd64 go build -o Releases/archive-helper-INFO-Windows-x64
GOOS=windows GOARCH=arm64 go build -o Releases/archive-helper-INFO-Windows-ARM64.exe
GOOS=darwin GOARCH=arm64 go build -o Releases/archive-helper-INFO-MacOS-AppleSilicon

read -n 1 -s -r -p "Make sure the .go file is set to output DEBUG logs.\nPress any key to continue..."
echo
# DEBUG Builds
GOOS=linux GOARCH=amd64 go build -o Releases/archive-helper-DEBUG-Linux-x64
GOOS=linux GOARCH=arm64 go build -o Releases/archive-helper-DEBUG-Linux-ARM64
GOOS=windows GOARCH=amd64 go build -o Releases/archive-helper-DEBUG-Windows-x64
GOOS=windows GOARCH=arm64 go build -o Releases/archive-helper-DEBUG-Windows-ARM64.exe
GOOS=darwin GOARCH=arm64 go build -o Releases/archive-helper-DEBUG-MacOS-AppleSilicon
