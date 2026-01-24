// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title JsonUtils
/// @notice JSON parsing and string manipulation utilities for shutdown scripts
/// @dev These utilities are designed for pattern counting in JSON strings,
///      not full JSON parsing. For full parsing, use vm.parseJson*() functions.
library JsonUtils {
    // ========== String Manipulation ==========

    /// @notice Remove 0x prefix from hex string
    /// @param value Input string (may or may not have 0x prefix)
    /// @return String without 0x prefix
    function strip0x(
        string memory value
    ) internal pure returns (string memory) {
        bytes memory data = bytes(value);
        if (
            data.length >= 2 &&
            data[0] == bytes1("0") &&
            (data[1] == bytes1("x") || data[1] == bytes1("X"))
        ) {
            bytes memory trimmed = new bytes(data.length - 2);
            for (uint256 i = 2; i < data.length; i++) {
                trimmed[i - 2] = data[i];
            }
            return string(trimmed);
        }
        return value;
    }

    /// @notice Convert hash string to function selector
    /// @dev Creates selector from keccak256 of "_" + stripped_hash + "()"
    /// @param hashString Hash string (with or without 0x prefix)
    /// @return Function selector (bytes4)
    function hashToSelector(
        string memory hashString
    ) internal pure returns (bytes4) {
        string memory stripped = strip0x(hashString);
        return bytes4(keccak256(abi.encodePacked("_", stripped, "()")));
    }

    // ========== JSON Pattern Counting ==========

    /// @notice Count occurrences of pattern in JSON string
    /// @param json JSON string to search in
    /// @param pattern Pattern to search for
    /// @return Number of occurrences
    function countOccurrences(
        string memory json,
        string memory pattern
    ) internal pure returns (uint256) {
        bytes memory bJson = bytes(json);
        bytes memory bPattern = bytes(pattern);
        uint256 count = 0;

        if (bJson.length < bPattern.length) return 0;

        for (uint256 i = 0; i <= bJson.length - bPattern.length; i++) {
            if (_isMatch(bJson, i, bPattern)) {
                count++;
            }
        }
        return count;
    }

    /// @notice Count tokens in JSON array (counts "l1Token" occurrences)
    /// @param json JSON string to search in
    /// @return Number of tokens
    function countTokens(string memory json) internal pure returns (uint256) {
        return countOccurrences(json, '"l1Token"');
    }

    /// @notice Count claims in a specific token entry (counts "hash" in token's data section)
    /// @param json JSON string to search in
    /// @param tokenIdx Token index (0-based)
    /// @return Number of claims in the token
    function countClaimsInToken(
        string memory json,
        uint256 tokenIdx
    ) internal pure returns (uint256) {
        return countOccurrencesInTokenData(json, tokenIdx, '"hash"');
    }

    /// @notice Count occurrences within a specific token's data section
    /// @param json JSON string to search in
    /// @param tokenIdx Token index (0-based)
    /// @param pattern Pattern to search for
    /// @return Number of occurrences within the token's section
    function countOccurrencesInTokenData(
        string memory json,
        uint256 tokenIdx,
        string memory pattern
    ) internal pure returns (uint256) {
        bytes memory bJson = bytes(json);
        bytes memory bPattern = bytes(pattern);
        bytes memory bL1Token = bytes('"l1Token"');
        uint256 count = 0;
        uint256 tokenFoundCount = 0;

        for (uint256 i = 0; i <= bJson.length - bL1Token.length; i++) {
            if (_isMatch(bJson, i, bL1Token)) {
                if (tokenFoundCount == tokenIdx) {
                    // Found the target token, count patterns until next token
                    for (
                        uint256 j = i + bL1Token.length;
                        j <= bJson.length - bPattern.length;
                        j++
                    ) {
                        // Stop if we hit the next token
                        if (_isMatch(bJson, j, bL1Token)) break;

                        if (_isMatch(bJson, j, bPattern)) {
                            count++;
                        }
                    }
                    return count;
                }
                tokenFoundCount++;
            }
        }
        return 0;
    }

    /// @notice Check if bytes match pattern at given index
    /// @param data Byte array to search in
    /// @param index Starting index
    /// @param pattern Pattern to match
    /// @return True if pattern matches at index
    function _isMatch(
        bytes memory data,
        uint256 index,
        bytes memory pattern
    ) internal pure returns (bool) {
        if (index + pattern.length > data.length) return false;
        for (uint256 i = 0; i < pattern.length; i++) {
            if (data[index + i] != pattern[i]) return false;
        }
        return true;
    }

    // ========== JSON Validation ==========

    /// @notice Check if JSON is an empty array
    /// @param json JSON string to check
    /// @return True if JSON represents an empty array []
    function isEmptyJsonArray(
        string memory json
    ) internal pure returns (bool) {
        bytes memory data = bytes(json);
        uint256 i = 0;

        // Skip leading whitespace
        while (i < data.length && _isWhitespace(data[i])) {
            i++;
        }
        if (i >= data.length || data[i] != "[") {
            return false;
        }
        i++;

        // Skip whitespace after [
        while (i < data.length && _isWhitespace(data[i])) {
            i++;
        }
        if (i >= data.length || data[i] != "]") {
            return false;
        }
        i++;

        // Check only whitespace remains
        while (i < data.length) {
            if (!_isWhitespace(data[i])) {
                return false;
            }
            i++;
        }
        return true;
    }

    // ========== Number Parsing ==========

    /// @notice Parse uint from string (supports hex and decimal formats)
    /// @dev Handles formats: "123", "0x7b", "DEC:123"
    /// @param value String to parse
    /// @return Parsed uint256 value
    function parseUint(string memory value) internal pure returns (uint256) {
        bytes memory data = bytes(value);
        require(data.length > 0, "JsonUtils: empty string");

        uint256 start = 0;
        uint256 end = data.length;

        // Trim whitespace
        while (start < end && _isWhitespace(data[start])) {
            start++;
        }
        while (end > start && _isWhitespace(data[end - 1])) {
            end--;
        }
        require(end > start, "JsonUtils: whitespace only");

        uint256 result = 0;

        // Handle hex format (0x or 0X prefix)
        if (
            end - start >= 2 &&
            data[start] == "0" &&
            (data[start + 1] == "x" || data[start + 1] == "X")
        ) {
            start += 2;
            for (uint256 i = start; i < end; i++) {
                uint8 c = uint8(data[i]);
                if (c >= 48 && c <= 57) {
                    // 0-9
                    result = (result << 4) + (c - 48);
                } else if (c >= 65 && c <= 70) {
                    // A-F
                    result = (result << 4) + (c - 55);
                } else if (c >= 97 && c <= 102) {
                    // a-f
                    result = (result << 4) + (c - 87);
                } else {
                    revert("JsonUtils: invalid hex char");
                }
            }
            return result;
        }

        // Handle decimal (skip non-digit prefix like "DEC:")
        uint256 i = start;
        while (i < end && (uint8(data[i]) < 48 || uint8(data[i]) > 57)) {
            i++;
        }
        require(i < end, "JsonUtils: no digits found");

        while (i < end && uint8(data[i]) >= 48 && uint8(data[i]) <= 57) {
            result = result * 10 + (uint8(data[i]) - 48);
            i++;
        }
        return result;
    }

    // ========== JSON String Escaping ==========

    /// @notice Escape special characters for JSON output
    /// @param str String to escape
    /// @return Escaped string safe for JSON
    function escapeJson(
        string memory str
    ) internal pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        uint256 escapeCount = 0;

        // Count characters that need escaping
        for (uint256 i = 0; i < strBytes.length; i++) {
            if (strBytes[i] == '"' || strBytes[i] == "\\") {
                escapeCount++;
            }
        }

        if (escapeCount == 0) return str;

        // Build escaped string
        bytes memory result = new bytes(strBytes.length + escapeCount);
        uint256 j = 0;
        for (uint256 i = 0; i < strBytes.length; i++) {
            if (strBytes[i] == '"' || strBytes[i] == "\\") {
                result[j++] = "\\";
            }
            result[j++] = strBytes[i];
        }

        return string(result);
    }

    // ========== Internal Helpers ==========

    /// @notice Check if byte is whitespace
    /// @param c Byte to check
    /// @return True if whitespace (space, newline, carriage return, tab)
    function _isWhitespace(bytes1 c) internal pure returns (bool) {
        return c == 0x20 || c == 0x0a || c == 0x0d || c == 0x09;
    }

    /// @notice Public wrapper for whitespace check
    /// @param c Byte to check
    /// @return True if whitespace
    function isWhitespace(bytes1 c) internal pure returns (bool) {
        return _isWhitespace(c);
    }
}
