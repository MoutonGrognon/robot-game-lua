package rgentities

/*
#include "rgentities.h"
*/
import "C"

type (
	ActionType int
	Action     struct {
		ActionType ActionType `json:"actionType"`
		X          int        `json:"x"`
		Y          int        `json:"y"`
	}

	Location struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	LocationType int

	Bot struct {
		X        int `json:"x"`
		Y        int `json:"y"`
		Hp       int `json:"hp"`
		Id       int `json:"id"`
		PlayerId int `json:"playerId"`
	}
	BotState struct {
		Bot    Bot    `json:"bot"`
		Action Action `json:"action"`
	}
)
