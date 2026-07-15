package constant

const (
	AnswerFrequencyOneTime                   = "One Time"
	AnswerFrequencyMultipleTimesOneDay       = "Multiple Times, One Day"
	AnswerFrequencyMultipleTimesDifferentDay = "Multiple Times, Different Day"
	AnswerFrequencyLegacyMultiple            = "Multiple"
)

func IsValidSurveyAnswerFrequencyForWrite(value string) bool {
	switch value {
	case AnswerFrequencyOneTime, AnswerFrequencyMultipleTimesOneDay, AnswerFrequencyMultipleTimesDifferentDay:
		return true
	default:
		return false
	}
}
