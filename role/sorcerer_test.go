package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SorcererTestSuite is a suite of unit tests for the Sorcerer role.
type SorcererTestSuite struct {
	suite.Suite
	sorcerer *Role
}

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *SorcererTestSuite) SetupTest() {
	suite.sorcerer = Sorcerer()
}

func (suite *SorcererTestSuite) TestAttributes() {
	// things a sorcerer is
	assert.True(suite.T(), suite.sorcerer.IsAuxEvil())
	assert.True(suite.T(), suite.sorcerer.ViewForAuxEvil())

	// things a sorcerer is not
	assert.False(suite.T(), suite.sorcerer.IsSeer())
	assert.False(suite.T(), suite.sorcerer.ViewForSeer())
	assert.False(suite.T(), suite.sorcerer.IsMaxEvil())
	assert.False(suite.T(), suite.sorcerer.ViewForMaxEvil())
}

// TestActions tests that a sorcerer has the start/night actions we expect. The functionality of those
// actions is covered in the game.actions suite.
func (suite *SorcererTestSuite) TestActions() {
	// things a sorcerer does
	assert.True(suite.T(), suite.sorcerer.ViewsForSeer())
	assert.True(suite.T(), suite.sorcerer.HasRandomN0Clear())

	// this a sorcerer does not do
	assert.False(suite.T(), suite.sorcerer.ViewsForMax())
	assert.False(suite.T(), suite.sorcerer.HasNightKill())
	assert.False(suite.T(), suite.sorcerer.ViewsForAux())
	assert.False(suite.T(), suite.sorcerer.KnowsMaxes())

}

// TestKill ensures that when you kill a sorcerer, they stay down.
func (suite *SorcererTestSuite) TestKill() {
	assert.False(suite.T(), suite.sorcerer.Kill())
}

func TestSorcererTestSuite(t *testing.T) {
	suite.Run(t, new(SorcererTestSuite))
}
