// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

/// @title Packed round requested state library
/// @notice Stores a packed mapping of round to its requested state
/// @notice The mapping uses uint248 for keys since rounds are represented as uint256 and there are 256 (2^8) values per word.
/**
 * @notice This code is adapted from Uniswap's TickBitmap implementation
 * @dev Original source: https://github.com/Uniswap/v4-core/blob/main/src/libraries/TickBitmap.sol
 * Modified to fit our specific use case
 */
library Bitmap {
    function position(uint256 round) internal pure returns (uint248 wordPos, uint8 bitPos) {
        assembly ("memory-safe") {
            wordPos := shr(8, round)
            bitPos := and(round, 0xff)
        }
    }

    function flipBit(mapping(uint248 => uint256) storage self, uint256 round) internal {
        assembly ("memory-safe") {
            // calculate the storage slot corresponding to the round
            // wordPos = round >> 8
            mstore(0, shr(8, round))
            mstore(0x20, self.slot)
            // the slot of self[wordPos] is keccak256(abi.encode(wordPos, self.slot))
            let slot := keccak256(0, 0x40)
            // mask = 1 << bitPos = 1 << (round & 0xff)
            // self[wordPos] ^= mask
            sstore(slot, xor(sload(slot), shl(and(round, 0xff), 1)))
        }
    }

    /// @notice includes the round itself
    function nextRequestedRound(mapping(uint248 => uint256) storage self, uint256 round)
        internal
        view
        returns (uint256 next, bool requested)
    {
        unchecked {
            (uint248 wordPos, uint8 bitPos) = position(round);
            // all the 1s at or to the left of the bitPos
            uint256 mask = ~((1 << bitPos) - 1);
            uint256 masked = self[wordPos] & mask;

            // if there are no requested rounds to the left of the current round, return leftmost in the word
            requested = masked != 0;
            next = requested ? round + leastSignificantBit(masked) - bitPos : round + type(uint8).max - bitPos;
        }
    }

    /// @notice Solady (https://github.com/Vectorized/solady/blob/8200a70e8dc2a77ecb074fc2e99a2a0d36547522/src/utils/LibBit.sol)
    /// @notice Returns the index of the least significant bit of the number,
    ///     where the least significant bit is at index 0 and the most significant bit is at index 255
    /// @param x the value for which to compute the least significant bit, must be greater than 0
    /// @return r the index of the least significant bit
    function leastSignificantBit(uint256 x) internal pure returns (uint8 r) {
        require(x > 0);

        assembly ("memory-safe") {
            // Isolate the least significant bit.
            x := and(x, sub(0, x))
            // For the upper 3 bits of the result, use a De Bruijn-like lookup.
            // Credit to adhusson: https://blog.adhusson.com/cheap-find-first-set-evm/
            // forgefmt: disable-next-item
            r := shl(
                5,
                shr(
                    252,
                    shl(
                        shl(
                            2,
                            shr(
                                250,
                                mul(
                                    x,
                                    0xb6db6db6ddddddddd34d34d349249249210842108c6318c639ce739cffffffff
                                )
                            )
                        ),
                        0x8040405543005266443200005020610674053026020000107506200176117077
                    )
                )
            )
            // For the lower 5 bits of the result, use a De Bruijn lookup.
            // forgefmt: disable-next-item
            r := or(
                r,
                byte(
                    and(div(0xd76453e0, shr(r, x)), 0x1f),
                    0x001f0d1e100c1d070f090b19131c1706010e11080a1a141802121b1503160405
                )
            )
        }
    }
}
