package main

type ConfigMissingError struct {
	config string
}

type UnknownSectionTypeError struct {
	sectiontype string
}

func (e *ConfigMissingError) Error() string {
	return "configuration is missing parameter '" + e.config + "'."
}

func (e *UnknownSectionTypeError) Error() string {
	return "unknown section type '" + e.sectiontype + "'."
}

func NewConfigMissingError(config string) error {
	e := &ConfigMissingError{}
	e.config = config
	return e
}

func NewUnknownSectionTypeError(sectiontype string) error {
	e := &UnknownSectionTypeError{}
	e.sectiontype = sectiontype
	return e
}
