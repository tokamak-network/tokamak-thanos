// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice Error returns when a non-depositor account tries to set L1 block values.
error NotDepositor();

/// @notice Error returns when a chain is not in the interop dependency set.
error NotDependency();

/// @notice Error returns when the interop dependency set size is too large.
error DependencySetSizeTooLarge();

/// @notice Error returns when adding a chain that is already in the interop dependency set.
error AlreadyDependency();

/// @notice Error returns when attempting to remove the chain's own chain ID from the dependency set.
error CantRemovedDependency();
