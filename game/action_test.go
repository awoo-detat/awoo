package game

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/awoo-detat/awoo/player"
	"github.com/awoo-detat/awoo/player/communicator"
	"github.com/awoo-detat/awoo/role"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ActionTestSuite is a suite of unit tests for the Seer role.
type ActionTestSuite struct {
	suite.Suite
	from    *player.Player
	to      *player.Player
	player3 *player.Player
	game    *Game
}

const (
	FromName    = "From"
	ToName      = "To"
	Player3Name = "Player 3"
)

// SetupTest performs any necessary actions to get ready for individual tests.
func (suite *ActionTestSuite) SetupTest() {
	devNull := &communicator.Communicator{
		Logger: log.New(ioutil.Discard, "", 0),
	}
	c := make(chan *player.Player)
	suite.from = player.New(devNull, c)
	suite.from.Name = FromName
	suite.to = player.New(devNull, c)
	suite.to.Name = ToName
	suite.player3 = player.New(devNull, c)
	suite.player3.Name = Player3Name
	suite.game = New(c)

	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

// TestNoNightAction ensures that a role with no night actions will have nothing run.
func (suite *ActionTestSuite) TestNoNightAction() {
	suite.from.Role = role.Cultist()
	suite.to.Role = role.Werewolf()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), "", result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestSorcererPositiveView tests that a sorcerer viewing a seer will be a positive hit.
func (suite *ActionTestSuite) TestSorcererPositiveView() {
	suite.from.Role = role.Sorcerer()
	suite.to.Role = role.Seer()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), fmt.Sprintf("%s %s", suite.to.Name, IsSeer), result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestSorcererNegativeView tests that a sorcerer viewing a non-seer will be a negative hit.
func (suite *ActionTestSuite) TestSorcererNegativeView() {
	suite.from.Role = role.Sorcerer()
	suite.to.Role = role.Cultist()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), fmt.Sprintf("%s %s", suite.to.Name, IsNotSeer), result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestSeerPositiveView tests that a seer viewing a wolf will be a positive hit.
func (suite *ActionTestSuite) TestSeerPositiveView() {
	suite.from.Role = role.Seer()
	suite.to.Role = role.Werewolf()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), fmt.Sprintf("%s %s", suite.to.Name, IsWerewolf), result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestSeerNegativeView tests that a seer viewing a non-wolf will be a negative hit.
func (suite *ActionTestSuite) TestSeerNegativeView() {
	suite.from.Role = role.Seer()
	suite.to.Role = role.Cultist()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), fmt.Sprintf("%s %s", suite.to.Name, IsNotWerewolf), result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestAuxSeerPositiveView tests that an aux seer viewing an auxwill be a positive hit.
func (suite *ActionTestSuite) TestAuxSeerPositiveView() {
	suite.from.Role = role.AuxSeer()
	suite.to.Role = role.Cultist()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), fmt.Sprintf("%s %s", suite.to.Name, IsAuxEvil), result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestAuxSeerNegativeView tests that an aux seer viewing a non-aux will be a negative hit.
func (suite *ActionTestSuite) TestAuxSeerNegativeView() {
	suite.from.Role = role.AuxSeer()
	suite.to.Role = role.Werewolf()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), fmt.Sprintf("%s %s", suite.to.Name, IsNotAuxEvil), result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestNightKill ensures that a wolf eating someone will be reported correctly.
func (suite *ActionTestSuite) TestNightKill() {
	suite.from.Role = role.Werewolf()
	suite.to.Role = role.Villager()

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), "", result.PlayerMessage)
	assert.Equal(suite.T(), suite.to, result.Killed)
}

// TestToughNightKill ensures that a wolf eating someone who doesn't die
// won't reveal the target
func (suite *ActionTestSuite) TestToughNightKill() {
	suite.from.Role = role.Werewolf()
	suite.to.Role = role.Villager()
	suite.to.Role.Health = 2

	result := suite.game.NightAction(&FingerPoint{suite.from, suite.to})
	assert.Equal(suite.T(), "", result.PlayerMessage)
	assert.Nil(suite.T(), result.Killed)
}

// TestKnowsMaxes tests that a role that is told who the wolves are will be told.
func (suite *ActionTestSuite) TestKnowsMaxes() {
	suite.from.Role = role.Cultist()
	suite.to.Role = role.Werewolf()
	suite.player3.Role = role.Werewolf()

	suite.game.AddPlayer(suite.from)
	suite.game.AddPlayer(suite.to)
	suite.game.AddPlayer(suite.player3)

	result := suite.game.StartAction(suite.from)
	assert.True(suite.T(), strings.HasPrefix(result.PlayerMessage, WolfListMultiple))
	result = suite.game.StartAction(suite.to)
	assert.True(suite.T(), strings.HasPrefix(result.PlayerMessage, WolfListMultiple))
}

// TestKnowsMax tests that a role that is told who the wolf is will be told.
func (suite *ActionTestSuite) TestKnowsMax() {
	suite.from.Role = role.Cultist()
	suite.to.Role = role.Werewolf()

	suite.game.AddPlayer(suite.from)
	suite.game.AddPlayer(suite.to)

	result := suite.game.StartAction(suite.from)
	assert.Equal(suite.T(), fmt.Sprintf("%s %s", WolfListSingle, suite.to.Name), result.PlayerMessage)
}

// TestKnowsMaxIsMax tests that a max evil won't be told they're the only one.
func (suite *ActionTestSuite) TestKnowsMaxIsMax() {
	suite.from.Role = role.Werewolf()
	suite.game.AddPlayer(suite.from)

	result := suite.game.StartAction(suite.from)
	assert.Equal(suite.T(), "", result.PlayerMessage)
}

// TestSeerN0Clear tests that a seer's N0 clear functions properly.
func (suite *ActionTestSuite) TestSeerN0Clear() {
	suite.from.Role = role.Seer()
	suite.to.Role = role.Werewolf()

	suite.game.AddPlayer(suite.from)
	suite.game.AddPlayer(suite.to)

	// this should never happen, but does test that a seer won't get a view if there is none
	result := suite.game.StartAction(suite.from)
	assert.Equal(suite.T(), "", result.PlayerMessage)

	// but then adding a non-wolf provides them as the clear
	suite.player3.Role = role.Villager()
	suite.game.AddPlayer(suite.player3)
	result = suite.game.StartAction(suite.from)
	assert.True(suite.T(), strings.HasSuffix(result.PlayerMessage, IsNotWerewolf))
}

// TestSorcererN0Clear tests that a sorcerer's N0 clear functions properly.
func (suite *ActionTestSuite) TestSorcererN0Clear() {
	suite.from.Role = role.Sorcerer()
	suite.to.Role = role.Seer()
	suite.player3.Role = role.Villager()

	suite.game.AddPlayer(suite.from)
	suite.game.AddPlayer(suite.to)
	suite.game.AddPlayer(suite.player3)

	result := suite.game.StartAction(suite.from)
	assert.True(suite.T(), strings.HasSuffix(result.PlayerMessage, IsNotSeer))
}

// TestAuxSeerN0Clear tests that an aux seer's N0 clear functions properly.
func (suite *ActionTestSuite) TestAuxSeerN0Clear() {
	suite.from.Role = role.AuxSeer()
	suite.to.Role = role.Cultist()
	suite.player3.Role = role.Villager()

	suite.game.AddPlayer(suite.from)
	suite.game.AddPlayer(suite.to)
	suite.game.AddPlayer(suite.player3)

	result := suite.game.StartAction(suite.from)
	assert.True(suite.T(), strings.HasSuffix(result.PlayerMessage, IsNotAuxEvil))
}

func TestActionTestSuite(t *testing.T) {
	suite.Run(t, new(ActionTestSuite))
}
