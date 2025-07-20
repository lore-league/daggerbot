package commands

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func roll(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return MessageSend(s, m, parseRoll(c.Args(), m.Author.DisplayName()))
}

func parseRoll(args []string, roller string) string {
	if len(args) < 1 {
		return rollDuality(roller)
	}

	var results []string

	for _, roll := range args {
		if strings.EqualFold(roll, "duality") || strings.EqualFold(roll, "duelity") {
			results = append(results, rollDuality(roller))
			continue
		}

		rollRegex := regexp.MustCompile("^[0-9]*d[0-9]+.*$")
		isMultiDiceRoll := rollRegex.MatchString(roll)

		diceNum, err := strconv.ParseFloat(roll, 64)
		if err != nil && !isMultiDiceRoll {
			results = append(results, fmt.Sprintf("%s is not a valid roll. A roll is a number, duality, or dice abbreviation", roll))
			continue
		}
		if isMultiDiceRoll {
			results = append(results, rollMultiDice(roll, roller))
		} else {
			diceRollResultString := strconv.FormatFloat(rollDice(diceNum), 'f', -1, 64)

			if diceRollResultString == "1" {
				results = append(results, fmt.Sprintf("%s d%s result is %s :cry:\n", roller, roll, diceRollResultString))
				continue
			}

			results = append(results, fmt.Sprintf("%s d%s result is %s\n", roller, roll, diceRollResultString))
		}
	}

	return strings.Join(results, "\n")
}

func rollDuality(roller string) string {

	hope := rollDice(12)
	hopeString := strconv.FormatFloat(hope, 'f', -1, 64)

	fear := rollDice(12)
	fearString := strconv.FormatFloat(fear, 'f', -1, 64)

	result := hope + fear
	resultString := strconv.FormatFloat(result, 'f', -1, 64)

	if hope == fear {
		return fmt.Sprintf("# %s CRIT!!! :dagger: :heart:\n> with double %s", strings.ToUpper(roller), hopeString)
	}

	dualityResult := fmt.Sprintf("%s rolled %s ", roller, resultString)
	if hope > fear {
		dualityResult += "with Hope :heart:"
	} else {
		dualityResult += "with Fear :dagger:"
	}

	return fmt.Sprintf("%s\n> _Hope_ was %s and _Fear_ was %s", dualityResult, hopeString, fearString)
}

func rollDice(diceSides float64) float64 {
	return math.Ceil(rand.Float64() * diceSides)
}

func rollMultiDice(diceDesignation string, roller string) string {
	var positiveModifier float64 = 0
	var negativeModifier float64 = 0
	var diceTotal float64 = 0
	var diceDesignationString = diceDesignation

	if strings.Contains(diceDesignation, "+") {
		modifierSplitStrings := strings.Split(diceDesignation, "+")
		positiveModifier, _ = strconv.ParseFloat(modifierSplitStrings[1], 64)
		diceDesignation = modifierSplitStrings[0]
	} else if strings.Contains(diceDesignation, "-") {
		modifierSplitStrings := strings.Split(diceDesignation, "-")
		negativeModifier, _ = strconv.ParseFloat(modifierSplitStrings[1], 64)
		diceDesignation = modifierSplitStrings[0]
	}

	splitString := strings.Split(diceDesignation, "d")
	numRolls, err := strconv.Atoi(splitString[0])
	if err != nil {
		numRolls = 1
	}
	rollValue, err := strconv.ParseFloat(splitString[1], 64)
	if err != nil {
		return fmt.Sprintf("%s is not a valid roll. A roll is a number, duality, or dice abbreviation", diceDesignationString)
	}
	for i := 0; i < numRolls; i++ {
		dicevalue := rollDice(rollValue)
		diceTotal += dicevalue
	}

	diceTotal += (positiveModifier - negativeModifier)
	diceTotalString := strconv.FormatFloat(diceTotal, 'f', -1, 64)

	return fmt.Sprintf("%s %s result is %s\n", roller, diceDesignationString, diceTotalString)
}

func init() {
	RegisterCommand(NewCommand("Roll", "Replies with Roll!", roll))
}
