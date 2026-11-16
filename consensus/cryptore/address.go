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

import "github.com/ylxai/ylex-core/v2/common"

// YLEX protocol addresses.
//
// An ICAN address is 22 bytes: a one byte network prefix, a one byte
// checksum and a 20 byte payload. YLEX uses the "ce" prefix, the same
// private-network prefix the previous chain used for networks outside its
// first two chain ids.
//
// The two addresses below are derived from the payload rather than written
// out by hand, so the checksum is always correct and the addresses can be
// regenerated from the payload alone.

// treasuryPayload and liquidityPayload are the address payloads of the two
// protocol addresses. They are derived rather than written out, so they
// always have the right length for the current prefix, and so it is
// obvious which of the two addresses is which. Changing them changes the
// addresses, and therefore the chain's genesis, so they must never change
// after launch.
var (
	treasuryPayload  = derivePayload(0x01)
	liquidityPayload = derivePayload(0x02)
)

// TreasuryAddress holds the 20% treasury share of every block reward.
//
// It is a normal account like any other, so the balance is public and
// verifiable on chain. The key that controls it is held by the YLEX
// Foundation's five-seat council; see GOVERNANCE.md. The council is bound by
// the same rules the chain is: every spend is published, and the multisig
// threshold is a property of the account, not a promise.
var TreasuryAddress = mustAddress(treasuryPayload)

// LiquidityAddress holds the 10% liquidity share of every block reward.
//
// It funds the on-chain AMM that provides the first market for YLEX, so the
// chain has liquidity from genesis without depending on any central exchange
// to list it. The address key is held by the same council, under a separate
// threshold, and it can only move into AMM positions.
var LiquidityAddress = mustAddress(liquidityPayload)

// derivePayload builds the 18 byte address payload for a YLEX protocol
// address.
//
// The payload is derived rather than written out, so it always has the
// right length for the current prefix length and it is obvious which of the
// two protocol addresses is which. The tag distinguishes them; the rest is
// filled so the payload cannot be confused with a user account.
func derivePayload(tag byte) []byte {
	n := common.AddressLength - len(common.DefaultNetworkID.Bytes()) - 1
	if n < 1 {
		panic("ylex: address prefix leaves no room for a payload")
	}
	payload := make([]byte, n)
	payload[0] = tag
	return payload
}

// mustAddress builds a valid YLEX ICAN address from a 20 byte payload and
// panics if the result does not verify. These are compile-time constants,
// so a failure here is a programming error rather than a runtime condition.
func mustAddress(payload []byte) common.Address {
	// YLEX addresses use the three ASCII bytes "ylx" as their network
	// prefix, so every address on the network starts with "ylx".
	prefix := common.DefaultNetworkID.Bytes()
	if len(prefix) != 3 {
		panic("ylex: expected the three byte ylx network prefix")
	}

	// address layout: prefix(3) || checksum(1) || payload(18)
	raw := make([]byte, 0, common.AddressLength)
	raw = append(raw, prefix...)
	raw = append(raw, common.Hex2Bytes(common.CalculateChecksum(payload, prefix))...)
	raw = append(raw, payload...)

	addr, err := common.HexToAddress(common.Bytes2Hex(raw))
	if err != nil {
		panic("ylex: invalid address constant: " + err.Error())
	}
	return addr
}
