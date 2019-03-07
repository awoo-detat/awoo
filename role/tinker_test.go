package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TinkerTestSuite is a suite of unit tests for the Tinker role.
type TinkerTestSuite struct {
	suite.Suite
}

// TestTinkerSeer asserts that a tinker seer will correctly be viewed incorrectly.
func (suite *TinkerTestSuite) TestTinkerSeer() {
	tinker := Seer()
	tinker.SetTinker()

	assert.False(suite.T(), tinker.ViewForSeer())
	assert.True(suite.T(), tinker.IsSeer())
	assert.True(suite.T(), tinker.ViewForMaxEvil())
	assert.False(suite.T(), tinker.IsMaxEvil())
	assert.True(suite.T(), tinker.ViewForAuxEvil())
	assert.False(suite.T(), tinker.IsAuxEvil())
}

// TestTinkerWerewolf asserts that a tinker werewolf will correctly be viewed incorrectly.
func (suite *TinkerTestSuite) TestTinkerWerewolf() {
	tinker := Werewolf()
	tinker.SetTinker()

	assert.True(suite.T(), tinker.ViewForSeer())
	assert.False(suite.T(), tinker.IsSeer())
	assert.False(suite.T(), tinker.ViewForMaxEvil())
	assert.True(suite.T(), tinker.IsMaxEvil())
	assert.True(suite.T(), tinker.ViewForAuxEvil())
	assert.False(suite.T(), tinker.IsAuxEvil())
}

// TestTinkerAux asserts that a tinker aux evil will correctly be viewed incorrectly.
func (suite *TinkerTestSuite) TestTinkerAux() {
	tinker := Cultist()
	tinker.SetTinker()

	assert.True(suite.T(), tinker.ViewForSeer())
	assert.False(suite.T(), tinker.IsSeer())
	assert.True(suite.T(), tinker.ViewForMaxEvil())
	assert.False(suite.T(), tinker.IsMaxEvil())
	assert.False(suite.T(), tinker.ViewForAuxEvil())
	assert.True(suite.T(), tinker.IsAuxEvil())
}

func TestTinkerTestSuite(t *testing.T) {
	suite.Run(t, new(TinkerTestSuite))
}
