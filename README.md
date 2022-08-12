# glogv

glogv is a zerolog log viewer.  It converts zerologs standard json tags (time, level, message and error) and makes them more readable in the console.  The output is color coded depending on the log level of the message.

Installation:

```bash
go install github.com/cwbriscoe/glogv@latest
```

Usage:

```bash
# make sure your $GOPATH/bin is in your path
tail --follow=name /path/to/file.log | glogv
# or
cat /path/to/file.log | glogv
# etc
```

Here are a couple of useful bash function shortcuts:

```bash
tl() {
  tail --follow=name $1 | glogv
}

cl() {
  cat $1 | glogv
}
```

Then simply use them like this:

```bash
tl /path/to/file.log 
cl /path/to/file.log 
```
