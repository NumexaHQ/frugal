package costcalculation

import "strings"

func CalculateOpenAICost(promptToken, completionToken int, model string) (cost float64) {
	if strings.Contains(model, "gpt-3.5-turbo-16k") {
		cost = 0.003*float64(promptToken) + 0.004*float64(completionToken)
	} else if strings.Contains(model, "gpt-3.5-turbo") {
		cost = 0.0015*float64(promptToken) + 0.002*float64(completionToken)
	} else if strings.Contains(model, "gpt-4-32k") {
		cost = 0.06*float64(promptToken) + 0.12*float64(completionToken)
	} else if strings.Contains(model, "gpt-4") {
		cost = 0.03*float64(promptToken) + 0.06*float64(completionToken)
	} else if strings.Contains(model, "ada") {
		cost = 0.0016*float64(promptToken) + 0.0016*float64(completionToken)
	} else if strings.Contains(model, "davinci") {
		cost = 0.1200*float64(promptToken) + 0.1200*float64(completionToken)
	} else if strings.Contains(model, "babbage") {
		cost = 0.0024*float64(promptToken) + 0.0024*float64(completionToken)
	} else if strings.Contains(model, "curie") {
		cost = 0.0030*float64(promptToken) + 0.0030*float64(completionToken)
	} else {
		cost = 0.0000
	}

	return cost / 1000
}
