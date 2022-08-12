# glogv

glogv is a zerolog log viewer.  It converts zerologs standard json tags (`time`, `level`, `message` and `error`) and makes them more pleasant to view in the console.  The output is color coded depending on the log level of the message.

glogv currently only works on `linux`.

Installation:

```bash
go install github.com/cwbriscoe/glogv@latest
```

Basic Usage:

```bash
# this willl 'cat' the file and display it in the console
glogv /path/to/file.log

# this will 'tail' the file
glogv -tail /path/to/file.log

# filter through grep to show only certain levels
glogv -tail /path/to/file.log | grep -e ERR -e WRN --color=never -a
# or
glogv -tail /path/to/file.log | awk 'match($2,/WRN|ERR/)'
```

Also works will rolled log files:

```bash
glogv /path/to/file.log.gz
```

Supports more than one file at a time:

```bash
glogv -tail /path/to/file1.log /path/to/file2.log
glogv /path/to/file1.log.gz /path/to/file2.log.gz
```

Can also be used as a STDIN reader:

```bash
# make sure your $GOPATH/bin is in your path
tail --follow=name /path/to/file.log | glogv
# or
cat /path/to/file.log | glogv
# etc
```
