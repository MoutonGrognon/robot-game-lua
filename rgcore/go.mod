module github.com/MoutonGrognon/robot-game-lua/rgcore

go 1.21.4

require (
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgconst v0.0.0
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgdebug v0.0.0
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgentities v0.0.0
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgerrors v0.0.0
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgutils v0.0.0
)

replace (
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgconst v0.0.0 => ./rgconst
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgdebug v0.0.0 => ./rgdebug
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgentities v0.0.0 => ./rgentities
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgerrors v0.0.0 => ./rgerrors
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgutils v0.0.0 => ./rgutils
)
