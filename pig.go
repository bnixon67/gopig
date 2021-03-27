/*
Copyright 2021 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// displayRules displays the rules of the game
func displayRules() {

	const rules = `
Pig is a simple dice game.

Each turn, a player repeatedly rolls a die until either a 1 is rolled or the
player decides to "hold":

- If the player rolls a 1, they score nothing and it becomes the next player's
  turn.
- If the player rolls any other number, it is added to their turn total and the
  player's turn continues.
- If a player chooses to "hold", their turn total is added to their score, and
  it becomes the next player's turn.

The first player to score 100 or more points wins.`

	fmt.Println(rules)
}

// getPlayerName prompts for and returns the name for player num.
//   num is a zero-based integer
func getPlayerName(num int) string {

	var (
		name string // name of the player
		err  error  // reader error
	)

	// loop until we get a valid name
	for valid := false; !valid; {

		// prompt for player name
		fmt.Printf("Enter name for player %d: ", num+1)

		// read player name
		reader := bufio.NewReader(os.Stdin)
		name, err = reader.ReadString('\n')

		// display message if name cannot be read and restart loop
		// used continue rather than else to improve readability
		if err != nil {
			fmt.Printf("Could not read player name. Error: %v\n",
				err)
			continue
		}

		// trim leading and trailing space from name
		name = strings.TrimSpace(name)

		// display message if name is empty and restart loop
		// used continue rather than else to improve readability
		if name == "" {
			fmt.Println("Name cannot be empty")
			continue
		}

		//  name is valid so set exit condition
		valid = true
	}

	// return name
	return name
}

// Roll die with sides and return result
// sides is number of sides for the die, should be > 0 else panic
func roll(sides int) int {
	// need to increment by 1 since Intn returns >= 0 and < sides
	return rand.Intn(sides) + 1
}

// askHold returns true if the player wants to hold
func askHold() bool {

	var (
		hold     bool   // return value
		response string // string input
		err      error  // error from reader
	)

	// loop until we get a valid response
	for valid := false; !valid; {

		fmt.Println("\nWould you like to [h]old or [r]oll?")

		// read response
		reader := bufio.NewReader(os.Stdin)
		response, err = reader.ReadString('\n')

		// display message if response cannot be read and restart loop
		if err != nil {
			fmt.Println("Could not read response. Error:", err)
			fmt.Println("Please enter h for hold or r for roll")
			continue
		}

		// trim leading and trailing space from name
		response = strings.TrimSpace(response)

		switch response {
		case "h":
			hold = true
			valid = true
		case "r":
			hold = false
			valid = true
		default:
			fmt.Println("Invalid response.",
				"Please enter h for hold or r for roll.\n")
		}

	}

	return hold
}

// playTurn for player name with current score and return new score
func playTurn(name string, score int) int {

	// score for the current turn
	turnTotal := 0

	fmt.Println("=========================================================")
	fmt.Printf("%s's turn\n", name)

	// loop under play holds (or busts)
	for hold := false; !hold; {
		// roll six sided die
		roll := roll(6)

		// if roll is 1, then player busted, so exit loop
		if roll == 1 {
			fmt.Printf("%s rolled a 1 and busted\n", name)
			turnTotal = 0
			break
		}

		// add roll to current turnTotal
		turnTotal += roll

		// display current score information
		fmt.Printf("%s rolled a %d", name, roll)
		fmt.Printf(", turn total of %d", turnTotal)
		fmt.Printf(",  potential score of %d\n", score+turnTotal)

		// ask player to hold or roll
		hold = askHold()
	}

	return turnTotal
}

// use init to seed random number generator with current time
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// struct to store player information
	type Player struct {
		name  string // name of player
		score int    // current score of player
	}

	// create two players
	var player [2]Player

	// display rules
	displayRules()

	// get player names
	for num := 0; num < len(player); num++ {
		player[num].name = getPlayerName(num)
		fmt.Printf("Player %d name is %s.\n\n", num+1, player[num].name)
	}

	// loop until there is a winner
	for winner := false; !winner; {

		// play turn for each players
		for num := 0; num < len(player); num++ {

			// play turn for a player and display current score
			player[num].score +=
				playTurn(player[num].name, player[num].score)
			fmt.Printf("%s's current score is %d\n\n",
				player[num].name, player[num].score)

			// if player has a score of 100 or more they win
			if player[num].score >= 100 {
				fmt.Printf("%s wins with a score of %d.\n",
					player[num].name, player[num].score)
				winner = true
				break // skip other players when winner
			}
		}
	}

}
