include $(GOROOT)/src/Make.inc

TARG=dyndnscd
GOFILES=main.go poll.go fetch.go error.go update.go

include $(GOROOT)/src/Make.cmd
