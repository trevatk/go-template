package logging

import (
	"fmt"

	"go.uber.org/zap"
)

// New
func New() (*zap.Logger, error) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("unable to create uber/zap development logger %v", err)
	}

	return logger, nil
}
