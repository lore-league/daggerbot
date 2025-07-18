package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func dmroll(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		args = c.Args()
	)

	if len(args) < 1 {
		MessagePrivateSend(s, m, rollDuality())
		return nil
	}

	var results []string

	for _, roll := range args {

		if roll == "duality" || roll == "duelity" {
			results = append(results, rollDuality())
			continue
		}

		isMultiDiceRoll, _ := regexp.MatchString("^[0-9]*d[0-9]+.*$", roll)

		diceNum, err := strconv.ParseFloat(roll, 64)
		if err != nil && !isMultiDiceRoll {
			results = append(results, fmt.Sprintf("%s is not a valid roll. A roll is a number, duality, or dice abbreviation", roll))
			continue
		}
		if isMultiDiceRoll {
			results = append(results, multiDiceRoll(roll))
		} else {
			_, diceRollResultString := rollDice(diceNum)

			if diceRollResultString == "1" {
				results = append(results, fmt.Sprintf("your d%s result is %s :cry:\n", roll, diceRollResultString))
				continue
			}

			results = append(results, fmt.Sprintf("your d%s result is %s\n", roll, diceRollResultString))
		}
	}

	MessagePrivateSend(s, m, strings.Join(results, "\n"))
	return nil

}

func init() {
	RegisterCommand(NewCommand("DMRoll", "Replies with Roll!", dmroll))
}
