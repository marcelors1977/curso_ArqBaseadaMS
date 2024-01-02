package create_client

import (
	"testing"

	mocks "github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateClientUseCase_Execute(t *testing.T) {
	m := &mocks.ClientGatewayMock{}
	m.On("Save", mock.Anything).Return(nil)
	uc := NewCreateClientUseCase(m)

	output, err := uc.Execute(CreateClientInputDto{
		Name:  "John Doe",
		Email: "john@doe",
	})
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "john@doe", output.Email)
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 1)
}
