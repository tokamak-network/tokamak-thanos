export enum Chains {
  Mainnet = 1,
  Sepolia = 11155111,
  LocalDevnet = 900
}

/**
 * Get DAO members with chain id
 *
 * @param chainid - Chain ID
 * @returns [foundation, dao] Pre-defined owner addresses
 * @throws Throw an error if the chain is not supported.
 */
export const getDAOMembers = (chainid: number): [string, string] => {
  if (chainid === Chains.Mainnet) {
    // Mainnet
    return [
      '0x0Fd5632f3b52458C31A2C3eE1F4b447001872Be9',
      '0x61dc95E5f27266b94805ED23D95B4C9553A3D049',
      // '0xDD9f0cCc044B0781289Ee318e5971b0139602C26', // TokamakDAO
    ]
  } else if (chainid === Chains.Sepolia) {
    // Sepolia
    return [
      '0x0Fd5632f3b52458C31A2C3eE1F4b447001872Be9',
      '0x61dc95E5f27266b94805ED23D95B4C9553A3D049',
      // '0xA2101482b28E3D99ff6ced517bA41EFf4971a386', // TokamakDAO
    ]
  } else if (chainid === Chains.LocalDevnet) {
    // LocalDevnet (for testing)
    return [
      '0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65', // devAccount-4
      '0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc' // devAccount-5
    ]
  } else {
    throw new Error('It is not a supported chain')
  }
}
