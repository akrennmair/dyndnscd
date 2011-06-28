include $(GOROOT)/src/Make.inc

TARG=dyndnscd
GOFILES=main.go poll.go

include $(GOROOT)/src/Make.cmd
