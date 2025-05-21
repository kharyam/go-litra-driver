package config

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileSystem is a mock implementation of FileSystem
type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockFileSystem) Create(name string) (*os.File, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockFileSystem) GetEnv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockFileSystem) UserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockFileSystem) IsNotExist(err error) bool {
	args := m.Called(err)
	return args.Bool(0)
}

// MockParser is a mock implementation of Parser
type MockParser struct {
	mock.Mock
}

func (m *MockParser) AddSection(section string) {
	m.Called(section)
}

func (m *MockParser) HasSection(section string) bool {
	args := m.Called(section)
	return args.Bool(0)
}

func (m *MockParser) Set(section, option, value string) {
	m.Called(section, option, value)
}

func (m *MockParser) Get(section, option string) (string, error) {
	args := m.Called(section, option)
	return args.String(0), args.Error(1)
}

func (m *MockParser) RemoveSection(section string) {
	m.Called(section)
}

func (m *MockParser) Sections() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockParser) SaveWithDelimiter(filename, delimiter string) error {
	args := m.Called(filename, delimiter)
	return args.Error(0)
}

// MockParserFactory is a mock implementation of ParserFactory
type MockParserFactory struct {
	mock.Mock
}

func (m *MockParserFactory) NewConfigParserFromFile(filename string) (Parser, error) {
	args := m.Called(filename)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(Parser), args.Error(1)
}

// MockFileInfo is a mock implementation of os.FileInfo
type MockFileInfo struct {
	mock.Mock
}

func (m *MockFileInfo) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockFileInfo) Size() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockFileInfo) Mode() os.FileMode {
	args := m.Called()
	return args.Get(0).(os.FileMode)
}

func (m *MockFileInfo) ModTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockFileInfo) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockFileInfo) Sys() interface{} {
	args := m.Called()
	return args.Get(0)
}

// TestExists tests the exists function
func TestExists(t *testing.T) {
	mockFS := &MockFileSystem{}

	// Test existing file
	mockFS.On("Stat", "/path/to/existing/file").Return(&MockFileInfo{}, nil).Once()
	existsResult, err := exists(mockFS, "/path/to/existing/file")
	assert.NoError(t, err)
	assert.True(t, existsResult)

	// Test non-existing file
	notExistErr := errors.New("file does not exist")
	mockFS.On("Stat", "/path/to/nonexistent/file").Return(nil, notExistErr).Once()
	mockFS.On("IsNotExist", notExistErr).Return(true).Once()
	existsResult, err = exists(mockFS, "/path/to/nonexistent/file")
	assert.NoError(t, err)
	assert.False(t, existsResult)

	// Test error case
	otherErr := errors.New("permission denied")
	mockFS.On("Stat", "/path/to/error/file").Return(nil, otherErr).Once()
	mockFS.On("IsNotExist", otherErr).Return(false).Once()
	existsResult, err = exists(mockFS, "/path/to/error/file")
	assert.Equal(t, otherErr, err)
	assert.False(t, existsResult)

	mockFS.AssertExpectations(t)
}

// TestGetConfig tests the getConfig function
func TestGetConfig(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockParserFactory := &MockParserFactory{}
	mockParser := &MockParser{}
	mockFile := &os.File{}

	// Test with XDG_CONFIG_HOME set
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()

	// First check if the directory exists
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(nil, os.ErrNotExist).Once()
	mockFS.On("IsNotExist", os.ErrNotExist).Return(true).Once()
	mockFS.On("MkdirAll", "/xdg/config/home/llgd", os.ModeDir).Return(nil).Once()

	// Then check if the config file exists
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(nil, os.ErrNotExist).Once()
	mockFS.On("Create", "/xdg/config/home/llgd/config").Return(mockFile, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	parser, configFile := getConfig(mockFS, mockParserFactory)
	assert.Equal(t, mockParser, parser)
	assert.Equal(t, "/xdg/config/home/llgd/config", configFile)

	// Test with XDG_CONFIG_HOME not set
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("").Once()
	mockFS.On("UserHomeDir").Return("/home/user", nil).Once()
	mockFS.On("Stat", "/home/user/.llgd_config").Return(nil, os.ErrNotExist).Once()
	mockFS.On("Create", "/home/user/.llgd_config").Return(mockFile, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/home/user/.llgd_config").Return(mockParser, nil).Once()

	parser, configFile = getConfig(mockFS, mockParserFactory)
	assert.Equal(t, mockParser, parser)
	assert.Equal(t, "/home/user/.llgd_config", configFile)

	mockFS.AssertExpectations(t)
	mockParserFactory.AssertExpectations(t)
}

// TestAddOrUpdateProfile tests the AddOrUpdateProfile function
func TestAddOrUpdateProfile(t *testing.T) {
	// Save original implementation
	originalFS := defaultFS
	originalParserFactory := defaultParserFactory

	// Create mocks
	mockFS := &MockFileSystem{}
	mockParserFactory := &MockParserFactory{}
	mockParser := &MockParser{}

	// Replace with mocks
	defaultFS = mockFS
	defaultParserFactory = mockParserFactory

	// Restore after test
	defer func() {
		defaultFS = originalFS
		defaultParserFactory = originalParserFactory
	}()

	// Setup mock behavior
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	// Test creating a new profile
	mockParser.On("HasSection", "test_profile").Return(false).Once()
	mockParser.On("AddSection", "test_profile").Once()
	mockParser.On("Set", "test_profile", Bright, "50").Once()
	mockParser.On("Set", "test_profile", Temp, "4000").Once()
	mockParser.On("Set", "test_profile", Power, "1").Once()
	mockParser.On("SaveWithDelimiter", "/xdg/config/home/llgd/config", "=").Return(nil).Once()

	AddOrUpdateProfile("test_profile", 50, 4000, 1)

	// Test updating an existing profile
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()
	mockParser.On("HasSection", "test_profile").Return(true).Once()
	mockParser.On("Set", "test_profile", Bright, "75").Once()
	mockParser.On("SaveWithDelimiter", "/xdg/config/home/llgd/config", "=").Return(nil).Once()

	AddOrUpdateProfile("test_profile", 75, -1, -1)

	// Test updating with -1 values (should not change)
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()
	mockParser.On("HasSection", "test_profile").Return(true).Once()
	mockParser.On("SaveWithDelimiter", "/xdg/config/home/llgd/config", "=").Return(nil).Once()

	AddOrUpdateProfile("test_profile", -1, -1, -1)

	mockFS.AssertExpectations(t)
	mockParserFactory.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}

// TestUpdateCurrentState tests the UpdateCurrentState function
func TestUpdateCurrentState(t *testing.T) {
	// Save original implementation
	originalFS := defaultFS
	originalParserFactory := defaultParserFactory

	// Create mocks
	mockFS := &MockFileSystem{}
	mockParserFactory := &MockParserFactory{}
	mockParser := &MockParser{}

	// Replace with mocks
	defaultFS = mockFS
	defaultParserFactory = mockParserFactory

	// Restore after test
	defer func() {
		defaultFS = originalFS
		defaultParserFactory = originalParserFactory
	}()

	// Setup mock behavior
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	// Test creating a new profile
	mockParser.On("HasSection", CurrentProfileName).Return(false).Once()
	mockParser.On("AddSection", CurrentProfileName).Once()
	mockParser.On("Set", CurrentProfileName, Bright, "50").Once()
	mockParser.On("Set", CurrentProfileName, Temp, "4000").Once()
	mockParser.On("Set", CurrentProfileName, Power, "1").Once()
	mockParser.On("SaveWithDelimiter", "/xdg/config/home/llgd/config", "=").Return(nil).Once()

	UpdateCurrentState(50, 4000, 1)

	mockFS.AssertExpectations(t)
	mockParserFactory.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}

// TestDeleteProfile tests the DeleteProfile function
func TestDeleteProfile(t *testing.T) {
	// Save original implementation
	originalFS := defaultFS
	originalParserFactory := defaultParserFactory

	// Create mocks
	mockFS := &MockFileSystem{}
	mockParserFactory := &MockParserFactory{}
	mockParser := &MockParser{}

	// Replace with mocks
	defaultFS = mockFS
	defaultParserFactory = mockParserFactory

	// Restore after test
	defer func() {
		defaultFS = originalFS
		defaultParserFactory = originalParserFactory
	}()

	// Setup mock behavior
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	// Test deleting an existing profile
	mockParser.On("HasSection", "test_profile").Return(true).Once()
	mockParser.On("RemoveSection", "test_profile").Once()
	mockParser.On("SaveWithDelimiter", "/xdg/config/home/llgd/config", "=").Return(nil).Once()

	DeleteProfile("test_profile")

	mockFS.AssertExpectations(t)
	mockParserFactory.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}

// TestReadProfile tests the ReadProfile function
func TestReadProfile(t *testing.T) {
	// Save original implementation
	originalFS := defaultFS
	originalParserFactory := defaultParserFactory

	// Create mocks
	mockFS := &MockFileSystem{}
	mockParserFactory := &MockParserFactory{}
	mockParser := &MockParser{}

	// Replace with mocks
	defaultFS = mockFS
	defaultParserFactory = mockParserFactory

	// Restore after test
	defer func() {
		defaultFS = originalFS
		defaultParserFactory = originalParserFactory
	}()

	// Setup mock behavior
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	// Test reading an existing profile
	mockParser.On("Get", "test_profile", Bright).Return("50", nil).Once()
	mockParser.On("Get", "test_profile", Temp).Return("4000", nil).Once()
	mockParser.On("Get", "test_profile", Power).Return("1", nil).Once()

	brightness, temperature, power := ReadProfile("test_profile")
	assert.Equal(t, 50, brightness)
	assert.Equal(t, 4000, temperature)
	assert.Equal(t, 1, power)

	// Test reading a non-existent profile
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	mockParser.On("Get", "nonexistent_profile", Bright).Return("", errors.New("section not found")).Once()
	mockParser.On("Get", "nonexistent_profile", Temp).Return("", errors.New("section not found")).Once()
	mockParser.On("Get", "nonexistent_profile", Power).Return("", errors.New("section not found")).Once()

	brightness, temperature, power = ReadProfile("nonexistent_profile")
	assert.Equal(t, -1, brightness)
	assert.Equal(t, -1, temperature)
	assert.Equal(t, -1, power)

	mockFS.AssertExpectations(t)
	mockParserFactory.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}

// TestReadCurrentState tests the ReadCurrentState function
func TestReadCurrentState(t *testing.T) {
	// Save original implementation
	originalFS := defaultFS
	originalParserFactory := defaultParserFactory

	// Create mocks
	mockFS := &MockFileSystem{}
	mockParserFactory := &MockParserFactory{}
	mockParser := &MockParser{}

	// Replace with mocks
	defaultFS = mockFS
	defaultParserFactory = mockParserFactory

	// Restore after test
	defer func() {
		defaultFS = originalFS
		defaultParserFactory = originalParserFactory
	}()

	// Setup mock behavior
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	// Test reading the current profile
	mockParser.On("Get", CurrentProfileName, Bright).Return("50", nil).Once()
	mockParser.On("Get", CurrentProfileName, Temp).Return("4000", nil).Once()
	mockParser.On("Get", CurrentProfileName, Power).Return("1", nil).Once()

	brightness, temperature, power := ReadCurrentState()
	assert.Equal(t, 50, brightness)
	assert.Equal(t, 4000, temperature)
	assert.Equal(t, 1, power)

	mockFS.AssertExpectations(t)
	mockParserFactory.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}

// TestGetProfileNames tests the GetProfileNames function
func TestGetProfileNames(t *testing.T) {
	// Save original implementation
	originalFS := defaultFS
	originalParserFactory := defaultParserFactory

	// Create mocks
	mockFS := &MockFileSystem{}
	mockParserFactory := &MockParserFactory{}
	mockParser := &MockParser{}

	// Replace with mocks
	defaultFS = mockFS
	defaultParserFactory = mockParserFactory

	// Restore after test
	defer func() {
		defaultFS = originalFS
		defaultParserFactory = originalParserFactory
	}()

	// Setup mock behavior
	mockFS.On("GetEnv", "XDG_CONFIG_HOME").Return("/xdg/config/home").Once()
	mockFS.On("Stat", "/xdg/config/home/llgd").Return(&MockFileInfo{}, nil).Once()
	mockFS.On("Stat", "/xdg/config/home/llgd/config").Return(&MockFileInfo{}, nil).Once()
	mockParserFactory.On("NewConfigParserFromFile", "/xdg/config/home/llgd/config").Return(mockParser, nil).Once()

	// Test getting profile names
	mockParser.On("Sections").Return([]string{CurrentProfileName, "profile1", "profile2"}).Once()

	profiles := GetProfileNames()

	// Verify current is first
	assert.Equal(t, CurrentProfileName, profiles[0])

	// Verify all profiles are included
	assert.Contains(t, profiles, "profile1")
	assert.Contains(t, profiles, "profile2")
	assert.Len(t, profiles, 3)

	mockFS.AssertExpectations(t)
	mockParserFactory.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}
