package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SeerTestSuite is a suite of unit tests for the Seer role.
type SeerTestSuite struct {
	suite.Suite
	seer *Role
}

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *SeerTestSuite) SetupTest() {
	suite.seer = Seer()
}

func (suite *SeerTestSuite) TestAttributes() {
	// things a seer is
	assert.True(suite.T(), suite.seer.IsSeer())
	assert.True(suite.T(), suite.seer.ViewForSeer())

	// things a seer is not
	assert.False(suite.T(), suite.seer.IsMaxEvil())
	assert.False(suite.T(), suite.seer.ViewForMaxEvil())
	assert.False(suite.T(), suite.seer.IsAuxEvil())
	assert.False(suite.T(), suite.seer.ViewForAuxEvil())
}

// TestActions tests that a seer has the start/night actions we expect. The functionality of those
// actions is covered in the game.actions suite.
func (suite *SeerTestSuite) TestActions() {
	// things a seer does
	assert.True(suite.T(), suite.seer.ViewsForMax())
	assert.True(suite.T(), suite.seer.HasRandomN0Clear())

	// this a seer does not do
	assert.False(suite.T(), suite.seer.HasNightKill())
	assert.False(suite.T(), suite.seer.ViewsForSeer())
	assert.False(suite.T(), suite.seer.ViewsForAux())
	assert.False(suite.T(), suite.seer.KnowsMaxes())

}

// TestKill ensures that when you kill a seer, they stay down.
func (suite *SeerTestSuite) TestKill() {
	assert.False(suite.T(), suite.seer.Kill())
}

func TestSeerTestSuite(t *testing.T) {
	suite.Run(t, new(SeerTestSuite))
}
