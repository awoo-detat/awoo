package tally

import (
	"github.com/awoo-detat/awoo/player"
)

type ShortTally struct {
	Candidate string   `json:"candidate"`
	Votes     []string `json:"votes"`
}

func Short(verbose []*TallyItem) []*ShortTally {
	var short []*ShortTally
	for _, item := range verbose {
		c := item.Candidate.Identifier()
		votes := []string{}
		for _, v := range item.Votes {
			votes = append(votes, v.Identifier())
		}
		short = append(short, &ShortTally{
			Candidate: c,
			Votes:     votes,
		})
	}

	return short
}

type TallyItem struct {
	Candidate player.Player   `json:"candidate"`
	Votes     []player.Player `json:"votes"`
}

func Item(c player.Player, v []player.Player) *TallyItem {
	return &TallyItem{
		Candidate: c,
		Votes:     v,
	}
}
