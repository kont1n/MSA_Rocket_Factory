package converter

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConverterSuite общий suite для тестов конвертеров
type ConverterSuite struct {
	suite.Suite
}

func (s *ConverterSuite) SetupSuite() {
}

func (s *ConverterSuite) SetupTest() {
}

func (s *ConverterSuite) TearDownSuite() {
}

// TestConverterSuite запускает все тесты конвертеров
func TestConverterSuite(t *testing.T) {
	suite.Run(t, new(ConverterSuite))
}
