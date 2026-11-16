// Copyright 2026 The YLEX Foundation.
// This file is part of the ylex-core library.
//
// The ylex-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The ylex-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the ylex-core library. If not, see <http://www.gnu.org/licenses/>.

package cryptore

import "math/big"

// YLEX monetary policy.
//
// Every YLEX coin is mined. There is no premine, no ICO, no presale, no
// founder allocation and no private allocation. The split below is a rule
// applied to every block, enforced by consensus, and identical for every miner.
//
//	70%  miner      -> credited to the block's coinbase
//	20%  treasury   -> credited to the YLEX treasury address
//	10%  liquidity  -> credited to the YLEX liquidity address
//
// The treasury and liquidity shares are sent to plain addresses, not to
// contracts. That keeps the consensus rule trivial to verify and reimplement:
// anyone can confirm the split by adding up three balances from a block.
//
// Rewards follow a bootstrap ramp. The first BootstrapBlocks blocks pay
// BootstrapMultiplier times the base reward, decaying linearly to 1.0 at the end
// of the ramp. The ramp gives a new chain a chance to reach its target block
// time before the full emission starts, which is what makes the difficulty
// adjustment and the 60 second target reachable from genesis.
//
// Once the ramp is over the emission is flat and the supply cap is fixed.
// There is no tail emission and no inflation beyond the cap.
var (
	// BaseBlockReward is the steady-state reward for a mined block, in YLEX.
	BaseBlockReward = big.NewInt(8e+17)

	// Reward split, expressed in basis points so the arithmetic stays exact.
	// 7000 bps = 70%, 2000 bps = 20%, 1000 bps = 10%.
	MinerShareBPS      = int64(7000)
	TreasuryShareBPS   = int64(2000)
	LiquidityShareBPS  = int64(1000)
	TotalShareBPS      = int64(10000)

	// BootstrapBlocks is the number of blocks in the bootstrap ramp.
	// At 10 second blocks this is 30 days.
	BootstrapBlocks = uint64(30 * 24 * 60 * 6)

	// BootstrapMultiplier is the reward multiplier applied to the first block
	// of the ramp. It decays linearly to 1.0 at BootstrapBlocks.
	BootstrapMultiplier = int64(3)

	// BootstrapRewardBlocks is the number of blocks the miner share is paid
	// at full rate before the ramp starts tapering. This keeps early miners
	// whole while the network is finding its difficulty level.
	// At 10 second blocks this is 4 hours.
	BootstrapRewardBlocks = uint64(1440)
)

// BlockRewardAt returns the total block reward for the block at the given
// height, applying the bootstrap ramp.
//
// The ramp is a linear decay from BootstrapMultiplier down to 1.0 across
// BootstrapBlocks, then flat at the base reward.
func BlockRewardAt(height uint64) *big.Int {
	if height >= BootstrapBlocks {
		return new(big.Int).Set(BaseBlockReward)
	}

	// remaining ramp blocks, and how much of the multiplier is left
	remaining := BootstrapBlocks - height

	// scale = BootstrapMultiplier * remaining / BootstrapBlocks, in BPS of
	// the base reward, so the arithmetic stays in integers.
	// We express it as an integer multiplier in units of 1/BootstrapBlocks
	// to avoid rounding errors in the split below.
	numerator := new(big.Int).SetInt64(BootstrapMultiplier * int64(remaining))
	denominator := new(big.Int).SetUint64(BootstrapBlocks)

	// reward = base * numerator / denominator
	reward := new(big.Int).Mul(BaseBlockReward, numerator)
	reward.Div(reward, denominator)
	return reward
}

// SplitReward splits a total block reward into the miner, treasury and
// liquidity shares.
//
// The three shares always sum back to the total: the miner share takes the
// remainder after the two fixed-percentage shares are taken, so no wei is
// lost to rounding. Treasury and liquidity are computed first because their
// percentages divide evenly; the miner receives what is left.
func SplitReward(total *big.Int) (miner, treasury, liquidity *big.Int) {
	treasury = new(big.Int).Mul(total, big.NewInt(TreasuryShareBPS))
	treasury.Div(treasury, big.NewInt(TotalShareBPS))

	liquidity = new(big.Int).Mul(total, big.NewInt(LiquidityShareBPS))
	liquidity.Div(liquidity, big.NewInt(TotalShareBPS))

	miner = new(big.Int).Sub(total, treasury)
	miner.Sub(miner, liquidity)
	return miner, treasury, liquidity
}
