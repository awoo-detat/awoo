package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuxSeerTestSuite is a suite of unit tests for the Seer role.
type AuxSeerTestSuite struct {
	suite.Suite
	auxseer *Role
}

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *AuxSeerTestSuite) SetupTest() {
	suite.auxseer = AuxSeer()
}

func (suite *AuxSeerTestSuite) TestAttributes() {
	// things an aux seer is
	assert.True(suite.T(), suite.auxseer.IsSeer())
	assert.True(suite.T(), suite.auxseer.ViewForSeer())

	// things an aux seer is not
	assert.False(suite.T(), suite.auxseer.IsMaxEvil())
	assert.False(suite.T(), suite.auxseer.ViewForMaxEvil())
	assert.False(suite.T(), suite.auxseer.IsAuxEvil())
	assert.False(suite.T(), suite.auxseer.ViewForAuxEvil())
}

// TestActions tests that an aux seer has the start/night actions we expect. The functionality of those
// actions is covered in the game.actions suite.
func (suite *AuxSeerTestSuite) TestActions() {
	// things an aux seer does
	assert.True(suite.T(), suite.auxseer.HasRandomN0Clear())
	assert.True(suite.T(), suite.auxseer.ViewsForAux())

	// this an aux seer does not do
	assert.False(suite.T(), suite.auxseer.ViewsForMax())
	assert.False(suite.T(), suite.auxseer.HasNightKill())
	assert.False(suite.T(), suite.auxseer.ViewsForSeer())
	assert.False(suite.T(), suite.auxseer.KnowsMaxes())
}

// TestKill ensures that when you kill an aux seer, they stay down.
func (suite *AuxSeerTestSuite) TestKill() {
	assert.False(suite.T(), suite.auxseer.Kill())
}

func TestAuxSeerTestSuite(t *testing.T) {
	suite.Run(t, new(AuxSeerTestSuite))
}
