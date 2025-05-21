package config

import (
	"os"
)

// FileSystem is an interface for file system operations
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	Create(name string) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	GetEnv(key string) string
	UserHomeDir() (string, error)
	IsNotExist(err error) bool
}

// Parser is an interface for config parser operations
type Parser interface {
	AddSection(section string)
	HasSection(section string) bool
	Set(section, option, value string)
	Get(section, option string) (string, error)
	RemoveSection(section string)
	Sections() []string
	SaveWithDelimiter(filename, delimiter string) error
}

// ParserFactory is an interface for creating config parsers
type ParserFactory interface {
	NewConfigParserFromFile(filename string) (Parser, error)
}
