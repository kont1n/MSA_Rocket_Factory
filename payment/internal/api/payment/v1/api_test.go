package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	// Arrange
	mockService := new(MockPaymentService)

	// Act
	api := NewAPI(mockService)

	// Assert
	assert.NotNil(t, api)
	assert.Equal(t, mockService, api.paymentService)
}

func TestNewAPI_WithNilService(t *testing.T) {
	// Act
	api := NewAPI(nil)

	// Assert
	assert.NotNil(t, api)
	assert.Nil(t, api.paymentService)
}
