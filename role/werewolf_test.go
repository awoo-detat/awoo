package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WerewolfTestSuite is a suite of unit tests for the Seer role.
type WerewolfTestSuite struct {
	suite.Suite
	werewolf *Role
}

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *WerewolfTestSuite) SetupTest() {
	suite.werewolf = Werewolf()
}

func (suite *WerewolfTestSuite) TestAttributes() {
	// things a werewolf is
	assert.True(suite.T(), suite.werewolf.IsMaxEvil())
	assert.True(suite.T(), suite.werewolf.ViewForMaxEvil())

	// things a werewolf is not
	assert.False(suite.T(), suite.werewolf.IsSeer())
	assert.False(suite.T(), suite.werewolf.ViewForSeer())
	assert.False(suite.T(), suite.werewolf.IsAuxEvil())
	assert.False(suite.T(), suite.werewolf.ViewForAuxEvil())
}

// TestActions tests that a werewolf has the start/night actions we expect. The functionality of
// those actions is covered in the game.actions suite.
func (suite *WerewolfTestSuite) TestActions() {
	// things a werewolf does
	assert.True(suite.T(), suite.werewolf.HasNightKill())
	assert.True(suite.T(), suite.werewolf.KnowsMaxes())

	// this a werewolf does not do
	assert.False(suite.T(), suite.werewolf.ViewsForMax())
	assert.False(suite.T(), suite.werewolf.HasRandomN0Clear())
	assert.False(suite.T(), suite.werewolf.ViewsForSeer())
	assert.False(suite.T(), suite.werewolf.ViewsForAux())

}

// TestKill ensures that when you kill a werewolf, they stay down.
func (suite *WerewolfTestSuite) TestKill() {
	assert.False(suite.T(), suite.werewolf.Kill())
}

func TestWerewolfTestSuite(t *testing.T) {
	suite.Run(t, new(WerewolfTestSuite))
}
