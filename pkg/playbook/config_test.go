package playbook_test

import (
	"testing"

	
	"github.com/DanielPickens/pkg/playbook"
	"github.com/stretchr/testify/assert"
)

var configs = playbook.NewConfigService(mock.NewConfigRepository(),
	playbook.NewPlaybookService(mock.NewPlaybookRepository()))

func TestGenerateOk(t *testing.T) {
	inventories := mock.NewInventoryRepository()

	inv, _ := inventories.Get("test1")

	assert.Nil(t, configs.Generate(inv))
}

func TestGenerateEmptyNamespace(t *testing.T) {

	inventories := mock.NewInventoryRepository()

	inv, _ := inventories.Get("Test2")

	assert.Error(t, configs.Generate(inv))
}

func TestDeleteOk(t *testing.T) {
	assert.Nil(t, configs.Delete("test"))
}
