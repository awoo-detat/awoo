package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HunterTestSuite is a suite of unit tests for the Seer role.
type HunterTestSuite struct {
	suite.Suite
	hunter *Role
}

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *HunterTestSuite) SetupTest() {
	suite.hunter = Hunter()
}

func (suite *HunterTestSuite) TestAttributes() {
	// things a hunter is not
	assert.False(suite.T(), suite.hunter.IsSeer())
	assert.False(suite.T(), suite.hunter.ViewForSeer())
	assert.False(suite.T(), suite.hunter.IsMaxEvil())
	assert.False(suite.T(), suite.hunter.ViewForMaxEvil())
	assert.False(suite.T(), suite.hunter.IsAuxEvil())
	assert.False(suite.T(), suite.hunter.ViewForAuxEvil())
	assert.Equal(suite.T(), suite.hunter.Parity, 2)
}

// TestActions tests that a hunter has the start/night actions we expect. The functionality of those
// actions is covered in the game.actions suite.
func (suite *HunterTestSuite) TestActions() {
	// this a hunter does not do
	assert.False(suite.T(), suite.hunter.ViewsForMax())
	assert.False(suite.T(), suite.hunter.HasRandomN0Clear())
	assert.False(suite.T(), suite.hunter.HasNightKill())
	assert.False(suite.T(), suite.hunter.ViewsForSeer())
	assert.False(suite.T(), suite.hunter.ViewsForAux())
	assert.False(suite.T(), suite.hunter.KnowsMaxes())
}

// TestKill ensures that when you kill a hunter, they stay down.
func (suite *HunterTestSuite) TestKill() {
	assert.False(suite.T(), suite.hunter.Kill())
}

func TestHunterTestSuite(t *testing.T) {
	suite.Run(t, new(HunterTestSuite))
}
