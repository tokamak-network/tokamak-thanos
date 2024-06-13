import argparse
import logging
import os
import subprocess
import json
import socket
import datetime
import time
import shutil
import re
import http.client
from multiprocessing import Process, Queue
import concurrent.futures
from collections import namedtuple

pjoin = os.path.join

parser = argparse.ArgumentParser(description='Bedrock devnet launcher')
parser.add_argument('--monorepo-dir', help='Directory of the monorepo', default=os.getcwd())
parser.add_argument('--allocs', help='Only create the allocs and exit', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--test', help='Tests the deployment, must already be deployed', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--fork-public-network',
                    help='Fork the public network',
                    type=bool,
                    default=os.environ.get('L1_FORK_PUBLIC_NETWORK').lower() == 'true' if os.environ.get('L1_FORK_PUBLIC_NETWORK') else False)
parser.add_argument('--l1-rpc-url', help='Public L1 RPC URL', type=str, default=os.environ.get('L1_RPC'))
parser.add_argument('--block-number', help='From block number', type=int, default=os.environ.get('BLOCK_NUMBER'))
parser.add_argument('--l2-native-token', help='L2 native token', type=str, default=os.environ.get('L2_NATIVE_TOKEN'))
parser.add_argument('--admin-key', help='The admin private key for upgrade contracts', type=str, default=os.environ.get('DEVNET_ADMIN_PRIVATE_KEY'))
parser.add_argument('--l2-image', help='Using local l2', type=str, default=os.environ.get('L2_IMAGE') if os.environ.get('L2_IMAGE') is not None else 'onthertech/thanos-op-geth:nightly')
parser.add_argument('--l1-beacon', help='Using beacon RPC', type=str, default=os.environ.get('L1_BEACON'))

log = logging.getLogger()

# Global environment variables
DEVNET_NO_BUILD = os.getenv('DEVNET_NO_BUILD') == "true"
DEVNET_FPAC = os.getenv('DEVNET_FPAC') == "true"

class Bunch:
    def __init__(self, **kwds):
        self.__dict__.update(kwds)

class ChildProcess:
    def __init__(self, func, *args):
        self.errq = Queue()
        self.process = Process(target=self._func, args=(func, args))

    def _func(self, func, args):
        try:
            func(*args)
        except Exception as e:
            self.errq.put(str(e))

    def start(self):
        self.process.start()

    def join(self):
        self.process.join()

    def get_error(self):
        return self.errq.get() if not self.errq.empty() else None


def main():
    args = parser.parse_args()

    validate_fork_public_network(args)

    monorepo_dir = os.path.abspath(args.monorepo_dir)
    devnet_dir = pjoin(monorepo_dir, '.devnet')
    contracts_bedrock_dir = pjoin(monorepo_dir, 'packages', 'tokamak', 'contracts-bedrock')
    deployment_dir = pjoin(contracts_bedrock_dir, 'deployments', 'devnetL1')
    forge_dump_path = pjoin(contracts_bedrock_dir, 'Deploy-900.json')
    op_node_dir = pjoin(args.monorepo_dir, 'op-node')
    ops_bedrock_dir = pjoin(monorepo_dir, 'ops-bedrock')
    deploy_config_dir = pjoin(contracts_bedrock_dir, 'deploy-config')
    devnet_config_path = pjoin(deploy_config_dir, 'devnetL1.json')
    devnet_config_template_path = pjoin(deploy_config_dir, 'devnetL1-template.json')
    ops_chain_ops = pjoin(monorepo_dir, 'op-chain-ops')
    sdk_dir = pjoin(monorepo_dir, 'packages', 'tokamak', 'sdk')
    bedrock_devnet_dir = pjoin(monorepo_dir, 'bedrock-devnet')

    if args.fork_public_network:
      if args.block_number is not None:
        block_number = str(args.block_number)
      else:
        response = json.loads(eth_new_head(args.l1_rpc_url))
        block_number = str(int(response['result'], 16))
    else:
      block_number = "0"

    paths = Bunch(
      mono_repo_dir=monorepo_dir,
      devnet_dir=devnet_dir,
      contracts_bedrock_dir=contracts_bedrock_dir,
      deployment_dir=deployment_dir,
      forge_dump_path=forge_dump_path,
      l1_deployments_path=pjoin(deployment_dir, '.deploy'),
      deploy_config_dir=deploy_config_dir,
      devnet_config_path=devnet_config_path,
      devnet_config_template_path=devnet_config_template_path,
      op_node_dir=op_node_dir,
      ops_bedrock_dir=ops_bedrock_dir,
      ops_chain_ops=ops_chain_ops,
      sdk_dir=sdk_dir,
      genesis_l1_path=pjoin(devnet_dir, 'genesis-l1.json'),
      genesis_l2_path=pjoin(devnet_dir, 'genesis-l2.json'),
      allocs_path=pjoin(devnet_dir, 'allocs-l1.json'),
      addresses_json_path=pjoin(devnet_dir, 'addresses.json'),
      sdk_addresses_json_path=pjoin(devnet_dir, 'sdk-addresses.json'),
      rollup_config_path=pjoin(devnet_dir, 'rollup.json'),
      fork_public_network=args.fork_public_network,
      l1_rpc_url=args.l1_rpc_url,
      block_number=block_number,
      l2_native_token=args.l2_native_token,
      bedrock_devnet_path=bedrock_devnet_dir,
      admin_key=args.admin_key,
      l1_beacon=args.l1_beacon,
    )

    if args.test:
      log.info('Testing deployed devnet')
      devnet_test(paths)
      return

    os.makedirs(devnet_dir, exist_ok=True)

    if args.allocs:
      devnet_l1_genesis(paths)
      return

    git_commit = subprocess.run(['git', 'rev-parse', 'HEAD'], capture_output=True, text=True).stdout.strip()
    git_date = subprocess.run(['git', 'show', '-s', "--format=%ct"], capture_output=True, text=True).stdout.strip()
    # CI loads the images from workspace, and does not otherwise know the images are good as-is
    if DEVNET_NO_BUILD:
        log.info('Skipping docker images build')
    else:
        log.info(f'Building docker images for git commit {git_commit} ({git_date})')
        run_command(['docker', 'compose', 'build', '--progress', 'plain',
                    '--build-arg', f'GIT_COMMIT={git_commit}', '--build-arg', f'GIT_DATE={git_date}'],
                cwd=paths.ops_bedrock_dir, env={
            'PWD': paths.ops_bedrock_dir,
            'DOCKER_BUILDKIT': '1', # (should be available by default in later versions, but explicitly enable it anyway)
            'COMPOSE_DOCKER_CLI_BUILD': '1',  # use the docker cache
            'L2_IMAGE': args.l2_image,
            'L1_DOCKER_FILE': 'Dockerfile.l1.fork' if paths.fork_public_network else 'Dockerfile.l1',
            'L1_RPC': paths.l1_rpc_url if paths.fork_public_network else '',
            'BLOCK_NUMBER': paths.block_number,
            'L1_FORK_PUBLIC_NETWORK': str(paths.fork_public_network),
            'L1_RPC_BEACON': paths.l1_beacon if paths.l1_beacon else ''
        })

    log.info('Devnet starting')
    devnet_deploy(paths, args)


def init_devnet_l1_deploy_config(paths, update_timestamp=False, temp=True):
    deploy_config = read_json(paths.devnet_config_template_path) if temp else read_json(paths.devnet_config_path)
    if update_timestamp:
        deploy_config['l1GenesisBlockTimestamp'] = '{:#x}'.format(int(time.time()))
    if DEVNET_FPAC:
        deploy_config['useFaultProofs'] = True
        deploy_config['faultGameMaxDuration'] = 10
    write_json(paths.devnet_config_path, deploy_config)

def init_admin_geth(paths):
    deploy_config = read_json(paths.devnet_config_template_path)
    admin_address = deploy_config['finalSystemOwner']

    f = open(pjoin(paths.bedrock_devnet_path, 'genesis.json'), "w+")
    run_command([
        'geth', '--dev', 'dumpgenesis'
    ], cwd=paths.bedrock_devnet_path, stdout=f)
    f.close()

    genesis = read_json(pjoin(paths.bedrock_devnet_path, 'genesis.json'))

    genesis["config"]["chainId"] = 900

    alloc = genesis['alloc']
    alloc[admin_address] = {
        "balance": "10000000000000000000"
    }
    genesis['alloc'] = alloc

    write_file(pjoin(paths.bedrock_devnet_path, 'keystore'), paths.admin_key[2:])
    write_file(pjoin(paths.bedrock_devnet_path, 'password'), '1234')

    write_json(pjoin(paths.bedrock_devnet_path, 'genesis.json'), genesis)
    geth_init(paths)
    os.environ.setdefault(
        key = 'GETH_DATADIR',
        value = 'data',
    )

    key_store = pjoin(paths.bedrock_devnet_path, 'data', 'keystore')
    key_file = glob.glob(f"{key_store}/*{admin_address[2:].lower()}*")
    os.environ.setdefault(
        key = 'ETH_KEYSTORE',
        value = pjoin(paths.bedrock_devnet_path, 'data', 'keystore', key_file[0])
    )


def geth_init(paths):
    run_command([
        'geth', '--dev', '--datadir', 'data', 'init', 'genesis.json'
    ], cwd=paths.bedrock_devnet_path)

    run_command([
        'geth', '--datadir', 'data', '--password', 'password', 'account', 'import', 'keystore'
    ], cwd=paths.bedrock_devnet_path)
    delete_file(pjoin(paths.bedrock_devnet_path, 'keystore'))
    delete_file(pjoin(paths.bedrock_devnet_path, 'password'))
    delete_file(pjoin(paths.bedrock_devnet_path, 'genesis.json'))

def devnet_l1_genesis(paths):
    log.info('Generating L1 genesis state')
    init_devnet_l1_deploy_config(paths)

    fqn = 'scripts/Deploy.s.sol:Deploy'
    run_command([
        'forge', 'script', '--chain-id', '900', fqn, "--sig", "runWithStateDump()"
    ], env={}, cwd=paths.contracts_bedrock_dir)

    forge_dump = read_json(paths.forge_dump_path)
    write_json(paths.allocs_path, { "accounts": forge_dump })
    os.remove(paths.forge_dump_path)

    shutil.copy(paths.l1_deployments_path, paths.addresses_json_path)

# Bring up the devnet where the contracts are deployed to L1
def devnet_deploy(paths, args):
    if os.path.exists(paths.genesis_l1_path):
        log.info('L1 genesis already generated.')
    else:
        log.info('Generating L1 genesis.')
        if os.path.exists(paths.allocs_path) == False or DEVNET_FPAC == True:
            # If this is the FPAC devnet then we need to generate the allocs
            # file here always. This is because CI will run devnet-allocs
            # without DEVNET_FPAC=true which means the allocs will be wrong.
            # Re-running this step means the allocs will be correct.
            devnet_l1_genesis(paths)

        # It's odd that we want to regenerate the devnetL1.json file with
        # an updated timestamp different than the one used in the devnet_l1_genesis
        # function.  But, without it, CI flakes on this test rather consistently.
        # If someone reads this comment and understands why this is being done, please
        # update this comment to explain.
        init_devnet_l1_deploy_config(paths, update_timestamp=True, temp=False)
        run_command([
            'go', 'run', 'cmd/main.go', 'genesis', 'l1',
            '--deploy-config', paths.devnet_config_path,
            '--l1-allocs', paths.allocs_path,
            '--l1-deployments', paths.addresses_json_path,
            '--outfile.l1', paths.genesis_l1_path,
        ], cwd=paths.op_node_dir)

    # Bring up L1
    log.info('Starting L1.')
    run_command(['docker', 'compose', 'up', '-d', 'l1'], cwd=paths.ops_bedrock_dir, env={
        'PWD': paths.ops_bedrock_dir,
        'L1_RPC': paths.l1_rpc_url if paths.fork_public_network else '',
        'BLOCK_NUMBER': paths.block_number,
    })
    wait_up(8545)
    wait_for_rpc_server('127.0.0.1:8545')

    if os.path.exists(paths.genesis_l2_path):
        log.info('L2 genesis and rollup configs already generated.')
    else:
        log.info('Generating L2 genesis and rollup configs.')
        run_command([
            'go', 'run', 'cmd/main.go', 'genesis', 'l2',
            '--l1-rpc', 'http://localhost:8545',
            '--deploy-config', paths.devnet_config_path,
            '--l1-deployments', paths.addresses_json_path,
            '--outfile.l2', paths.genesis_l2_path,
            '--outfile.rollup', paths.rollup_config_path
        ], cwd=paths.op_node_dir)

    rollup_config = read_json(paths.rollup_config_path)
    addresses = read_json(paths.addresses_json_path)

    # Setup the beacon path
    run_command(['docker', 'compose', 'up', '-d', 'setup'],
    cwd=paths.ops_bedrock_dir)


    # Restart l1
    restart_l1_with_docker_compose(paths)

    # Reset the `genesis-l1.json` config file fork times.
    with open(paths.genesis_l1_path, 'r') as file:
        file_content = file.read()
    file_content = re.sub(r'"shanghaiTime".*$', '"shanghaiTime": 0,', file_content, flags=re.MULTILINE)
    file_content = re.sub(r'"cancunTime".*$', '"cancunTime": 0,', file_content, flags=re.MULTILINE)
    with open(paths.genesis_l1_path, 'w') as file:
        file.write(file_content)

    # Bring up beacon node
    log.info('Bringing up consensus-node and validator-client')
    run_command(['docker', 'compose', 'up', '-d', 'consensus-node', 'validator-client'], cwd=paths.ops_bedrock_dir)

    # Start the L2.
    log.info('Bringing up L2.')
    run_command(['docker', 'compose', 'up', '-d', 'l2'], cwd=paths.ops_bedrock_dir, env={
        'PWD': paths.ops_bedrock_dir,
        'L2_IMAGE': args.l2_image
    })

    # Wait for the L2 to be available.
    wait_up(9545)
    wait_for_rpc_server('127.0.0.1:9545')

    # Print out the addresses being used for easier debugging.
    l2_output_oracle = addresses['L2OutputOracleProxy']
    dispute_game_factory = addresses['DisputeGameFactoryProxy']
    batch_inbox_address = rollup_config['batch_inbox_address']
    log.info(f'Using L2OutputOracle {l2_output_oracle}')
    log.info(f'Using DisputeGameFactory {dispute_game_factory}')
    log.info(f'Using batch inbox {batch_inbox_address}')
    # Set up the base docker environment.
    docker_env={
        'PWD': paths.ops_bedrock_dir,
        'SEQUENCER_BATCH_INBOX_ADDRESS': batch_inbox_address,
        'L1_RPC': paths.l1_rpc_url if paths.fork_public_network else '',
        'BLOCK_NUMBER': paths.block_number,
        'WAITING_L1_PORT': '9999' if paths.fork_public_network else '8545',
        'L1_BEACON': paths.l1_beacon if paths.l1_beacon else '',
        'L1_FORK_PUBLIC_NETWORK': str(paths.fork_public_network)
    }

    # Selectively set the L2OO_ADDRESS or DGF_ADDRESS if using FPAC.
    # Must be done selectively because op-proposer throws if both are set.
    if DEVNET_FPAC:
        docker_env['DGF_ADDRESS'] = dispute_game_factory
        docker_env['DG_TYPE'] = '0'
        docker_env['PROPOSAL_INTERVAL'] = '10s'
    else:
        docker_env['L2OO_ADDRESS'] = l2_output_oracle


    # Bring up op-node, op-proposer, op-batcher, artifact-server
    log.info('Bringing up op-node, op-proposer, op-batcher and artifact-server.')
    run_command(['docker', 'compose', 'up', '-d', 'op-node', 'op-proposer', 'op-batcher', 'artifact-server'], cwd=paths.ops_bedrock_dir, env=docker_env)

    # Optionally bring up op-challenger.
    if DEVNET_FPAC:
        log.info('Bringing up op-challenger.')
        run_command(['docker', 'compose', 'up', '-d', 'op-challenger'], cwd=paths.ops_bedrock_dir, env=docker_env)

    log.info('Devnet ready.')


def eth_accounts(url):
    log.info(f'Fetch eth_accounts {url}')
    conn = http.client.HTTPConnection(url)
    headers = {'Content-type': 'application/json'}
    body = '{"id":2, "jsonrpc":"2.0", "method": "eth_accounts", "params":[]}'
    conn.request('POST', '/', body, headers)
    response = conn.getresponse()
    data = response.read().decode()
    conn.close()
    return data

def eth_new_head(url):
    parsed_url = urllib.parse.urlparse(url.strip())
    hostname = parsed_url.hostname
    path = parsed_url.path if parsed_url.path else '/'
    # Create a context that does not verify SSL certificates
    context = ssl._create_unverified_context()
    conn = http.client.HTTPSConnection(hostname, context=context)
    headers = {'Content-type': 'application/json'}
    body = '{"id":2, "jsonrpc":"2.0", "method": "eth_blockNumber", "params":[]}'
    try:
        conn.request('POST', path, body, headers)
        response = conn.getresponse()
        data = response.read().decode()

    except Exception as e:
        log.error(f"Failed to fetch the latest block: {e}")
        data = None
    finally:
        conn.close()

    return data


def wait_for_rpc_server(url):
    log.info(f'Waiting for RPC server at {url}')

    headers = {'Content-type': 'application/json'}
    body = '{"id":1, "jsonrpc":"2.0", "method": "eth_chainId", "params":[]}'

    while True:
        try:
            conn = http.client.HTTPConnection(url)
            conn.request('POST', '/', body, headers)
            response = conn.getresponse()
            if response.status < 300:
                log.info(f'RPC server at {url} ready')
                return
        except Exception as e:
            log.info(f'Waiting for RPC server at {url}')
            time.sleep(1)
        finally:
            if conn:
                conn.close()


CommandPreset = namedtuple('Command', ['name', 'args', 'cwd', 'timeout'])


def devnet_test(paths):
    # Run the two commands with different signers, so the ethereum nonce management does not conflict
    # And do not use devnet system addresses, to avoid breaking fee-estimation or nonce values.
    run_commands([
        CommandPreset('erc20-test',
          ['npx', 'hardhat',  'deposit-erc20', '--network',  'devnetL1',
           '--l1-contracts-json-path', paths.addresses_json_path, '--signer-index', '14'],
          cwd=paths.sdk_dir, timeout=8*60),
        CommandPreset('eth-test',
          ['npx', 'hardhat',  'deposit-eth', '--network',  'devnetL1',
           '--l1-contracts-json-path', paths.addresses_json_path, '--signer-index', '15'],
          cwd=paths.sdk_dir, timeout=8*60)
    ], max_workers=2)


def run_commands(commands: list[CommandPreset], max_workers=2):
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(run_command_preset, cmd) for cmd in commands]

        for future in concurrent.futures.as_completed(futures):
            result = future.result()
            if result:
                print(result.stdout)


def run_command_preset(command: CommandPreset):
    with subprocess.Popen(command.args, cwd=command.cwd,
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True) as proc:
        try:
            # Live output processing
            for line in proc.stdout:
                # Annotate and print the line with timestamp and command name
                timestamp = datetime.datetime.utcnow().strftime('%H:%M:%S.%f')
                # Annotate and print the line with the timestamp
                print(f"[{timestamp}][{command.name}] {line}", end='')

            stdout, stderr = proc.communicate(timeout=command.timeout)

            if proc.returncode != 0:
                raise RuntimeError(f"Command '{' '.join(command.args)}' failed with return code {proc.returncode}: {stderr}")

        except subprocess.TimeoutExpired:
            raise RuntimeError(f"Command '{' '.join(command.args)}' timed out!")

        except Exception as e:
            raise RuntimeError(f"Error executing '{' '.join(command.args)}': {e}")

        finally:
            # Ensure process is terminated
            proc.kill()
    return proc.returncode


def run_command(args, check=True, shell=False, cwd=None, env=None, timeout=None, stdout=None):
    env = env if env else {}
    return subprocess.run(
        args,
        check=check,
        shell=shell,
        env={
            **os.environ,
            **env
        },
        cwd=cwd,
        timeout=timeout,
        stdout=stdout
    )


def wait_up(port, retries=10, wait_secs=1):
    for i in range(0, retries):
        log.info(f'Trying 127.0.0.1:{port}')
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            s.connect(('127.0.0.1', int(port)))
            s.shutdown(2)
            log.info(f'Connected 127.0.0.1:{port}')
            return True
        except Exception:
            time.sleep(wait_secs)

    raise Exception(f'Timed out waiting for port {port}.')


def write_json(path, data):
    with open(path, 'w+') as f:
        json.dump(data, f, indent='  ')


def read_json(path):
    with open(path, 'r') as f:
        return json.load(f)


def validate_fork_public_network(args):
    fork_public_network = args.fork_public_network
    l1_rpc_url = args.l1_rpc_url
    l2_native_token = args.l2_native_token
    block_number = args.block_number
    # If fork the public network, validate the required params related to
    if fork_public_network:
      if not l1_rpc_url:
        raise Exception("Please provide the L1_RPC URL for the forked network.")

      if not l2_native_token:
        raise Exception("Please provide the L2_NATIVE_TOKEN for the forked network.")

      log.info(f'Fork from RPC URL: {l1_rpc_url}, block number: {block_number}, l2 native token: {l2_native_token}')


def write_file(path, data):
    f = open(path, 'w+')
    f.write(data)
    f.close()

def delete_file(path):
    os.remove(path)

def restart_l1_with_docker_compose(paths):
    try:
        # Restart L1
        log.info('Re-starting L1.')
        subprocess.run(['docker', 'compose', 'up', '-d', 'l1'], cwd=paths.ops_bedrock_dir, env={
            'PWD': paths.ops_bedrock_dir,
            'L1_RPC': paths.l1_rpc_url if paths.fork_public_network else '',
            'BLOCK_NUMBER': paths.block_number,
        }, check=True)

        wait_up(8545)
        wait_for_rpc_server('127.0.0.1:8545')

        print("L1 is restarted.")
    except subprocess.CalledProcessError as e:
        print(f"Cannot restart L1 with error: {e}")