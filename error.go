package main

import (
	"os"
)

type ConfigMissingError struct {
	config string
}

type UnknownSectionTypeError struct {
	sectiontype string
}

func (e *ConfigMissingError) String() string {
	return "configuration is missing parameter '" + e.config + "'."
}

func (e *UnknownSectionTypeError) String() string {
	return "unknown section type '" + e.sectiontype + "'."
}

func NewConfigMissingError(config string) os.Error {
	e := &ConfigMissingError{}
	e.config = config
	return e
}

func NewUnknownSectionTypeError(sectiontype string) os.Error {
	e := &UnknownSectionTypeError{}
	e.sectiontype = sectiontype
	return e
}

