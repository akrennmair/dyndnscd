include $(GOROOT)/src/Make.inc

TARG=dyndnscd
GOFILES=main.go poll.go fetch.go error.go update.go log.go

include $(GOROOT)/src/Make.cmd
