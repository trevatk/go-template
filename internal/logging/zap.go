// Package logging service logger
package logging

import (
	"fmt"

	"go.uber.org/zap"
)

// New create new uber/zap logger instance
func New() (*zap.Logger, error) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("unable to create uber/zap development logger %v", err)
	}

	return logger, nil
}
