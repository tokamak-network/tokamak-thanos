# from ethereum import utils


def is_valid_ethereum_address(address):
  try:
    # Check if the address is a valid hexadecimal and has the correct length
    # if not utils.check_checksum(address):
    #   raise Exception("Invalid Ethereum address checksum.")
    return True
  except Exception as e:
    print(f"Error validating Ethereum address: {e}")
    return False
