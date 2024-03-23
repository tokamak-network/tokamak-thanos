import argparse
import time
import http.client
import json
import urllib

parser = argparse.ArgumentParser(description="Your script description")
parser.add_argument("--l1-rpc", type=str, required=True, help="L1 RPC")

args = parser.parse_args()


def read_data(filename):
  # Open the file and read its content
  with open(filename, 'r') as f:
    data = json.load(f)
  return data

def make_req(url, payload):
  try:
    parsed_url = urllib.parse.urlparse(url)
    conn = http.client.HTTPConnection(parsed_url.hostname, parsed_url.port)

    headers = {
      'Content-Type': 'application/json'
    }

    # Send the request
    conn.request("POST", parsed_url.path, payload, headers)

    # Get the response
    response = conn.getresponse()

    # Check for successful response
    if response.status == 200:
      # Close the connection
      conn.close()
      return True
    else:
      print(f"Error: {response.status}")
      # Close the connection
      conn.close()
      return False
  except:
    return False


def main():
  l1_rpc = args.l1_rpc

  print(f"Connected to RPC node: {l1_rpc}")

  data_path = "./genesis.json"
  data = read_data(data_path)

  for addr, acc in data["alloc"].items():
    balance = acc.get("balance")
    code = acc.get("code")
    nonce = acc.get("nonce")
    storage = acc.get("storage")
    if balance is not None:
      payload = json.dumps({
        "jsonrpc": "2.0",
        "id": 1,
        "method": "anvil_setBalance",
        "params": [
          addr,
          balance
        ]
      })
      is_success = make_req(l1_rpc, payload)
      if not is_success:
        print(f"failed to set balance, addr: {addr}, balance: {balance}")

    if code is not None:
      payload = json.dumps({
        "jsonrpc": "2.0",
        "id": 1,
        "method": "anvil_setCode",
        "params": [
          addr,
          code
        ]
      })
      is_success = make_req(l1_rpc, payload)
      if not is_success:
        print(f"failed to set code, addr: {addr}, code: {code}")

    if nonce is not None:
      payload = json.dumps({
        "jsonrpc": "2.0",
        "id": 1,
        "method": "anvil_setNonce",
        "params": [
          addr,
          nonce
        ]
      })
      is_success = make_req(l1_rpc, payload)
      if not is_success:
        print(f"failed to set nonce, addr: {addr}, nonce: {nonce}")

    if storage is not None:
      for slot, val in storage.items():
        payload = json.dumps({
          "jsonrpc": "2.0",
          "id": 1,
          "method": "anvil_setStorageAt",
          "params": [
            addr,
            slot,
            val
          ]
        })
        is_success = make_req(l1_rpc, payload)
        if not is_success:
          print(f"failed to set storage, addr: {addr}, slot: {slot}, val: {val}")



if __name__ == "__main__":
  main()
