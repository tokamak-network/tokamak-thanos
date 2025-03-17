// Semgrep tests for Solidity rules are defined in this file.
// Semgrep tests do not need to be valid Solidity code but should be syntactically correct so that
// Semgrep can parse them. You don't need to be able to *run* the code here but it should look like
// the code that you expect to catch with the rule.
//
// Semgrep testing 101
// Use comments like "ruleid: <rule-id>" to assert that the rule catches the code.
// Use comments like "ok: <rule-id>" to assert that the rule does not catch the code.

/// NOTE: Semgrep limitations mean that the rule for this check is defined as a relatively loose regex that searches the
/// remainder of the file after the `@custom:proxied` natspec tag is detected.
/// This means that we must test the case without this natspec tag BEFORE the case with the tag or the rule will apply
/// to the remainder of the file.

// If no predeploy natspec, disableInitializers can or cannot be called in constructor
contract SemgrepTest__sol_safety_use_disable_initializer {
    // ok: sol-safety-use-disable-initializer
    constructor() {
        // ...
        _disableInitializers();
        // ...
    }

    // ok: sol-safety-use-disable-initializer
    constructor() {
        // ...
    }
}

// if no predeploy natspec, disableInitializers must be called in constructor
/// @custom:proxied true
contract SemgrepTest__sol_safety_use_disable_initializer {
    // ok: sol-safety-use-disable-initializer
    constructor() {
        // ...
        _disableInitializers();
        // ...
    }

    // ruleid: sol-safety-use-disable-initializer
    constructor() {
        // ...
    }
}

/// NOTE: the predeploy natspec below is valid for all contracts after this one
/// @custom:predeploy
// if predeploy natspec, disableInitializers may or may not be called in constructor
contract SemgrepTest__sol_safety_use_disable_initializer {
    // ok: sol-safety-use-disable-initializer
    constructor() {
        // ...
    }

    // ok: sol-safety-use-disable-initializer
    constructor() {
        // ...
        _disableInitializers();
        // ...
    }
}
