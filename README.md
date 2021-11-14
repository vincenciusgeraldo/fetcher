# Fetcher
Simple web page saver using Go.

## How to use it
### Using executable
#### Prerequisite
- Go 1.16
####
```bash
make dep build
```
```bash
./bin/fetcher <OPTIONS> <url> <url> ...
```
### Using Docker image
#### Prerequisite
- Docker
####
```bash
docker build -t fetcher .
```
```bash
sudo docker run -v <path_to_folder>:/go/src/app fetcher <OPTIONS> <url> <url>
```

## Command Options
| Option      | Description |
| ----------- | ----------- |
| -l, --local      | save assets locally       |
| -m, --metadata   | print metadata        |