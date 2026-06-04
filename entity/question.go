package entity

type Question struct {
	Id              uint
	Text            string
	PossibleAnswers []PossibleAnswer
	CorrectAnswer   uint
	Difficulty      QuestionDifficulty
	CategoryId      uint
}

type PossibleAnswer struct {
	Id     uint
	Text   string
	Choice PossibleAnswerChoice
}

type PossibleAnswerChoice uint

const (
	FirstAnswerChoice PossibleAnswerChoice = iota + 1
	SecondAnswerChoice
	ThirdAnswerChoice
	FourthAnswerChoice
)

type QuestionDifficulty string

const (
	QuestionDifficultyEasy   = "question-difficulty-easy"
	QuestionDifficultyMedium = "question-difficulty-medium"
	QuestionDifficultyHard = "question-difficulty-hard"
)
