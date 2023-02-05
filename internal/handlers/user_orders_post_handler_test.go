package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/helpers"
)

func TestUserOrdersPOSTHandler_ServeHTTP_Negative(t *testing.T) {

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}

	testEnvironment := NewTestEnvironment()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderNumber := helpers.RandomString(10)
			orderDTO := request.NewUserOrdersPOSTDTO(orderNumber)

			body, errMarshaling := json.Marshal(orderDTO)
			require.NoError(t, errMarshaling)

			response := testEnvironment.ServeHandler(t, testEnvironment.userOrderPOSTHandler, body)
			assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		})
	}
}
