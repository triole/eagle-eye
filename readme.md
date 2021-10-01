# Eagle Eye ![build](https://github.com/triole/eagle-eye/actions/workflows/build.yaml/badge.svg)

<!--- mdtoc: toc begin -->

1. [Synopsis](#synopsis)
2. [Usage](#usage)
3. [Variables](#variables)
4. [Disclaimer](#disclaimer)<!--- mdtoc: toc end -->

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
ee -p
```

## Variables

```go mdox-exec="r -p"

The following vars are available:

  {{.dir}}	folder of the file that triggered the event
  {{.file}}	file that triggered the event

```

## Disclaimer

Warning. Use this software at your own risk. I may not be hold responsible for any data loss, starving your kittens or losing the bling bling powerpoint presentation you made to impress human resources with the efficiency of your employee's performance.
