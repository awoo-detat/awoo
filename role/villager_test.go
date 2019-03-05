package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VillagerTestSuite is a suite of unit tests for the Seer role.
type VillagerTestSuite struct {
	suite.Suite
	villager *Role
}

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *VillagerTestSuite) SetupTest() {
	suite.villager = Villager()
}

func (suite *VillagerTestSuite) TestAttributes() {
	// things a villager is not
	assert.False(suite.T(), suite.villager.IsSeer())
	assert.False(suite.T(), suite.villager.ViewForSeer())
	assert.False(suite.T(), suite.villager.IsMaxEvil())
	assert.False(suite.T(), suite.villager.ViewForMaxEvil())
	assert.False(suite.T(), suite.villager.IsAuxEvil())
	assert.False(suite.T(), suite.villager.ViewForAuxEvil())
}

// TestActions tests that a villager has the start/night actions we expect. The functionality of those
// actions is covered in the game.actions suite.
func (suite *VillagerTestSuite) TestActions() {
	// this a villager does not do
	assert.False(suite.T(), suite.villager.ViewsForMax())
	assert.False(suite.T(), suite.villager.HasRandomN0Clear())
	assert.False(suite.T(), suite.villager.HasNightKill())
	assert.False(suite.T(), suite.villager.ViewsForSeer())
	assert.False(suite.T(), suite.villager.ViewsForAux())
	assert.False(suite.T(), suite.villager.KnowsMaxes())
}

// TestKill ensures that when you kill a villager, they stay down.
func (suite *VillagerTestSuite) TestKill() {
	assert.False(suite.T(), suite.villager.Kill())
}

func TestVillagerTestSuite(t *testing.T) {
	suite.Run(t, new(VillagerTestSuite))
}
