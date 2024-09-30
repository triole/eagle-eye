# Eagle Eye ![build](https://github.com/triole/eagle-eye/actions/workflows/build.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Usage](#usage)
- [Variables](#variables)
- [Disclaimer](#disclaimer)

<!-- /toc -->

## Synopsis

A folder watcher that runs commands on change of files or folders. Inspired and based on [watcher](https://github.com/radovskyb/watcher.git) but with more intuitive command line flags and logging support.

## Usage

```shell
# simple usage, watch current dir and run 'ls -la'
ee ls -la

# pass events path to the command
ee cat {{.file}}

# spectate mode, just prints changes, does not execute a command
ee -s

# filter and watch a specific folder
# NOTE: flags always have to be in front of the command
ee -r .md$ -f /home/user/my_markdowns ls -lah

# get a list of available variables
ee -n
```

## Variables

```go mdox-exec="r -n"
[36mINFO   [0m[2024-09-30 10:17:25.348 CEST] Watch folder                                  [36maction[0m=just spectate [36mfolder[0m=/home/olaf/rolling/golang/projects/eagle-eye/src [36mlog-file[0m=/dev/stdout [36mlog-json[0m=false [36mlog-level[0m=info [36mlog-no-colours[0m=false

The following vars are available:

  {{.file}}	file that triggered the event
  {{.folder}}	folder of the file that triggered the event

```

## Disclaimer

Warning. Use this software at your own risk. I may not be hold responsible for any data loss, starving your kittens or losing the bling bling powerpoint presentation you made to impress human resources with the efficiency of your employee's performance.
