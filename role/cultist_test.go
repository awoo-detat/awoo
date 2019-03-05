package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CultistTestSuite is a suite of unit tests for the Seer role.
type CultistTestSuite struct {
	suite.Suite
	cultist *Role
}

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *CultistTestSuite) SetupTest() {
	suite.cultist = Cultist()
}

func (suite *CultistTestSuite) TestAttributes() {
	// things a cultist is
	assert.True(suite.T(), suite.cultist.IsAuxEvil())
	assert.True(suite.T(), suite.cultist.ViewForAuxEvil())

	// things a cultist is not
	assert.False(suite.T(), suite.cultist.IsSeer())
	assert.False(suite.T(), suite.cultist.ViewForSeer())
	assert.False(suite.T(), suite.cultist.IsMaxEvil())
	assert.False(suite.T(), suite.cultist.ViewForMaxEvil())
}

// TestActions tests that a cultist has the start/night actions we expect. The functionality of those
// actions is covered in the game.actions suite.
func (suite *CultistTestSuite) TestActions() {
	// things a cultist does
	assert.True(suite.T(), suite.cultist.KnowsMaxes())

	// this a cultist does not do
	assert.False(suite.T(), suite.cultist.ViewsForMax())
	assert.False(suite.T(), suite.cultist.HasRandomN0Clear())
	assert.False(suite.T(), suite.cultist.HasNightKill())
	assert.False(suite.T(), suite.cultist.ViewsForSeer())
	assert.False(suite.T(), suite.cultist.ViewsForAux())

}

// TestKill ensures that when you kill a seer, they stay down.
func (suite *CultistTestSuite) TestKill() {
	assert.False(suite.T(), suite.cultist.Kill())
}

func TestCultistTestSuite(t *testing.T) {
	suite.Run(t, new(CultistTestSuite))
}
