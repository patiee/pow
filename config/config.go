package config

import "math/big"

// BlockDifficultyEpoch is the number of blocks when network adjusts the difficulty
const BlockDifficultyEpoch = 2016

var MaxDifficultyTarget, _ = new(big.Int).SetString("00000000FFFF0000000000000000000000000000000000000000000000000000", 16) // 2^(256-32)
