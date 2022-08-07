# glogv

glogv is a zerolog log viewer.  It converts zerologs standard json tags (time, level, message and error) and makes them more readable in the console.  The output is color coded depending on the log level of the message.

glogv does not suport the custom json that is created when using the .Str(k string, v string) method of zerolog.  It would be pretty easy to add your logs custom json if you wanted to fork the repo.

Installation:

```bash
go install github.com/cwbriscoe/glogv@latest
```

Usage:

```bash
# make sure your GOPATH/bin is in your path
tail -f /path/to/file.log | glogv
# or
cat /path/to/file.log | glogv
# etc
```

Here are a couple of useful bash function shortcuts:

```bash
tl() {
  tail -f $1 | glogv
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
