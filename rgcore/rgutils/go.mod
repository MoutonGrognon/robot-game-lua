module github.com/MoutonGrognon/robot-game-lua/rgcore/rgutils

go 1.21.4

require (
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgconst v0.0.0
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgentities v0.0.0
)

replace (
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgconst v0.0.0 => ../rgconst
	github.com/MoutonGrognon/robot-game-lua/rgcore/rgentities v0.0.0 => ../rgentities
)
