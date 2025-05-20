package config

import (
	"os"

	"github.com/bigkevmcd/go-configparser"
)

// DefaultFileSystem implements FileSystem using actual OS functions
type DefaultFileSystem struct{}

func (fs *DefaultFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *DefaultFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fs *DefaultFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *DefaultFileSystem) GetEnv(key string) string {
	return os.Getenv(key)
}

func (fs *DefaultFileSystem) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (fs *DefaultFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// DefaultParserFactory implements ParserFactory using actual configparser
type DefaultParserFactory struct{}

func (f *DefaultParserFactory) NewConfigParserFromFile(filename string) (Parser, error) {
	parser, err := configparser.NewConfigParserFromFile(filename)
	if err != nil {
		return nil, err
	}
	return &ConfigParserWrapper{parser: parser}, nil
}

// ConfigParserWrapper wraps configparser.ConfigParser to implement Parser
type ConfigParserWrapper struct {
	parser *configparser.ConfigParser
}

func (w *ConfigParserWrapper) AddSection(section string) {
	w.parser.AddSection(section)
}

func (w *ConfigParserWrapper) HasSection(section string) bool {
	return w.parser.HasSection(section)
}

func (w *ConfigParserWrapper) Set(section, option, value string) {
	w.parser.Set(section, option, value)
}

func (w *ConfigParserWrapper) Get(section, option string) (string, error) {
	return w.parser.Get(section, option)
}

func (w *ConfigParserWrapper) RemoveSection(section string) {
	w.parser.RemoveSection(section)
}

func (w *ConfigParserWrapper) Sections() []string {
	return w.parser.Sections()
}

func (w *ConfigParserWrapper) SaveWithDelimiter(filename, delimiter string) error {
	return w.parser.SaveWithDelimiter(filename, delimiter)
}
