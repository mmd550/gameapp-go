package entity

import "time"

type game struct{
	Id uint
	CategoryId uint
	QuestionIds []uint
	PlayerIds []uint
	StartTime time.Time
}

type Player struct{
	Id uint
    GameId uint
	UserId uint
	Score uint
	Answers []PlayerAnswer
}

type PlayerAnswer struct{
	Id uint
	PlayerId uint
	QuestionId uint
	AnswerId uint
}