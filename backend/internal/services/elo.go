package services

import "math"

// K-factor: 32 for new artists (<100 votes), 16 for established (100+ votes)
func KFactor(voteCount int) int {
	if voteCount < 100 {
		return 32
	}
	return 16
}

func CalculateELO(winnerScore, loserScore, kFactor int) (newWinner, newLoser int) {
	expectedWinner := 1.0 / (1.0 + math.Pow(10, float64(loserScore-winnerScore)/400))
	expectedLoser := 1.0 - expectedWinner
	newWinner = winnerScore + int(float64(kFactor)*(1-expectedWinner))
	newLoser = loserScore + int(float64(kFactor)*(0-expectedLoser))
	return
}
