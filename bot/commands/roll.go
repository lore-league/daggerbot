package commands

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func roll(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		args = c.Args()
	)
	args = nil

	if len(args) < 1 {
		MessageSend(s, m, rollDuality())
		return nil
	}

	diceNum, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		MessageSend(s, m, "Need to roll a number or duality")
		return nil
	}
	diceRollResult := int(math.Ceil(rand.Float64() * diceNum))
	MessageSend(s, m, fmt.Sprintf("%o", diceRollResult))
	return nil

}

func init() {
	RegisterCommand(NewCommand("Roll", "Replies with Roll!", roll))
}

func rollDuality() string {

	hope := int(math.Ceil(rand.Float64() * 12))
	fear := int(math.Ceil(rand.Float64() * 12))
	var dualityResult string

	if hope > fear {
		dualityResult = "You rolled with Hope :heart:"
	} else if fear > hope {
		dualityResult = "You rolled with Fear :dagger:"
	} else {
		dualityResult = "YOU CRIT!!!! :dagger: :heart:"
	}

	return fmt.Sprintf(">>> %s \nYour Hope roll was %o \nYour Fear roll was %o \nYour roll was %o\n %f", dualityResult, hope, fear, hope+fear, (rand.Float64() * 12))

}

func rollDice(num int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return (r.Intn(num) + 1)
}
