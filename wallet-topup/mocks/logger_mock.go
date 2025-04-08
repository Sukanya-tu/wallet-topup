package mocks

import "github.com/stretchr/testify/mock"

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Error(args ...interface{}) {
	l.Called(args...)
}

func (l *LoggerMock) Warn(args ...interface{}) {
	l.Called(args...)
}

func (l *LoggerMock) Info(args ...interface{}) {
	l.Called(args...)
}

func (l *LoggerMock) Errorf(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *LoggerMock) Warnf(format string, args ...interface{}) {
	l.Called(format, args)
}

func (l *LoggerMock) Infof(format string, args ...interface{}) {
	l.Called(format, args)
}
