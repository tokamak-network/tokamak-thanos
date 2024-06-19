import argparse
import logging
import os
import subprocess
import json
import socket
import datetime
import time
import shutil
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

    # Insert state data in genesis-l1.json
    file_path = paths.genesis_l1_path

    state_deposit_contract = '''
    "0x4242424242424242424242424242424242424242": {
      "balance": "0",
      "code": "0x60806040526004361061003f5760003560e01c806301ffc9a71461004457806322895118146100a4578063621fd130146101ba578063c5f2892f14610244575b600080fd5b34801561005057600080fd5b506100906004803603602081101561006757600080fd5b50357fffffffff000000000000000000000000000000000000000000000000000000001661026b565b604080519115158252519081900360200190f35b6101b8600480360360808110156100ba57600080fd5b8101906020810181356401000000008111156100d557600080fd5b8201836020820111156100e757600080fd5b8035906020019184600183028401116401000000008311171561010957600080fd5b91939092909160208101903564010000000081111561012757600080fd5b82018360208201111561013957600080fd5b8035906020019184600183028401116401000000008311171561015b57600080fd5b91939092909160208101903564010000000081111561017957600080fd5b82018360208201111561018b57600080fd5b803590602001918460018302840111640100000000831117156101ad57600080fd5b919350915035610304565b005b3480156101c657600080fd5b506101cf6110b5565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102095781810151838201526020016101f1565b50505050905090810190601f1680156102365780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561025057600080fd5b506102596110c7565b60408051918252519081900360200190f35b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a70000000000000000000000000000000000000000000000000000000014806102fe57507fffffffff0000000000000000000000000000000000000000000000000000000082167f8564090700000000000000000000000000000000000000000000000000000000145b92915050565b6030861461035d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806118056026913960400191505060405180910390fd5b602084146103b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603681526020018061179c6036913960400191505060405180910390fd5b6060821461040f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806118786029913960400191505060405180910390fd5b670de0b6b3a7640000341015610470576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806118526026913960400191505060405180910390fd5b633b9aca003406156104cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260338152602001806117d26033913960400191505060405180910390fd5b633b9aca00340467ffffffffffffffff811115610535576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602781526020018061182b6027913960400191505060405180910390fd5b6060610540826114ba565b90507f649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c589898989858a8a6105756020546114ba565b6040805160a0808252810189905290819060208201908201606083016080840160c085018e8e80828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910187810386528c815260200190508c8c808284376000838201819052601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690920188810386528c5181528c51602091820193918e019250908190849084905b83811015610648578181015183820152602001610630565b50505050905090810190601f1680156106755780820380516001836020036101000a031916815260200191505b5086810383528881526020018989808284376000838201819052601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018881038452895181528951602091820193918b019250908190849084905b838110156106ef5781810151838201526020016106d7565b50505050905090810190601f16801561071c5780820380516001836020036101000a031916815260200191505b509d505050505050505050505050505060405180910390a1600060028a8a600060801b604051602001808484808284377fffffffffffffffffffffffffffffffff0000000000000000000000000000000090941691909301908152604080517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0818403018152601090920190819052815191955093508392506020850191508083835b602083106107fc57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016107bf565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610859573d6000803e3d6000fd5b5050506040513d602081101561086e57600080fd5b5051905060006002806108846040848a8c6116fe565b6040516020018083838082843780830192505050925050506040516020818303038152906040526040518082805190602001908083835b602083106108f857805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016108bb565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610955573d6000803e3d6000fd5b5050506040513d602081101561096a57600080fd5b5051600261097b896040818d6116fe565b60405160009060200180848480828437919091019283525050604080518083038152602092830191829052805190945090925082918401908083835b602083106109f457805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016109b7565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610a51573d6000803e3d6000fd5b5050506040513d6020811015610a6657600080fd5b5051604080516020818101949094528082019290925280518083038201815260609092019081905281519192909182918401908083835b60208310610ada57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610a9d565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610b37573d6000803e3d6000fd5b5050506040513d6020811015610b4c57600080fd5b50516040805160208101858152929350600092600292839287928f928f92018383808284378083019250505093505050506040516020818303038152906040526040518082805190602001908083835b60208310610bd957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610b9c565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610c36573d6000803e3d6000fd5b5050506040513d6020811015610c4b57600080fd5b50516040518651600291889160009188916020918201918291908601908083835b60208310610ca957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610c6c565b6001836020036101000a0380198251168184511680821785525050505050509050018367ffffffffffffffff191667ffffffffffffffff1916815260180182815260200193505050506040516020818303038152906040526040518082805190602001908083835b60208310610d4e57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610d11565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610dab573d6000803e3d6000fd5b5050506040513d6020811015610dc057600080fd5b5051604080516020818101949094528082019290925280518083038201815260609092019081905281519192909182918401908083835b60208310610e3457805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610df7565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610e91573d6000803e3d6000fd5b5050506040513d6020811015610ea657600080fd5b50519050858114610f02576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260548152602001806117486054913960600191505060405180910390fd5b60205463ffffffff11610f60576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806117276021913960400191505060405180910390fd5b602080546001019081905560005b60208110156110a9578160011660011415610fa0578260008260208110610f9157fe5b0155506110ac95505050505050565b600260008260208110610faf57fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061102557805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610fe8565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015611082573d6000803e3d6000fd5b5050506040513d602081101561109757600080fd5b50519250600282049150600101610f6e565b50fe5b50505050505050565b60606110c26020546114ba565b905090565b6020546000908190815b60208110156112f05781600116600114156111e6576002600082602081106110f557fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061116b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161112e565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa1580156111c8573d6000803e3d6000fd5b5050506040513d60208110156111dd57600080fd5b505192506112e2565b600283602183602081106111f657fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061126b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161122e565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa1580156112c8573d6000803e3d6000fd5b5050506040513d60208110156112dd57600080fd5b505192505b6002820491506001016110d1565b506002826112ff6020546114ba565b600060401b6040516020018084815260200183805190602001908083835b6020831061135a57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161131d565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790527fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000095909516920191825250604080518083037ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8018152601890920190819052815191955093508392850191508083835b6020831061143f57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101611402565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa15801561149c573d6000803e3d6000fd5b5050506040513d60208110156114b157600080fd5b50519250505090565b60408051600880825281830190925260609160208201818036833701905050905060c082901b8060071a60f81b826000815181106114f457fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060061a60f81b8260018151811061153757fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060051a60f81b8260028151811061157a57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060041a60f81b826003815181106115bd57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060031a60f81b8260048151811061160057fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060021a60f81b8260058151811061164357fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060011a60f81b8260068151811061168657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060001a60f81b826007815181106116c957fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050919050565b6000808585111561170d578182fd5b83861115611719578182fd5b505082019391909203915056fe4465706f736974436f6e74726163743a206d65726b6c6520747265652066756c6c4465706f736974436f6e74726163743a207265636f6e7374727563746564204465706f7369744461746120646f6573206e6f74206d6174636820737570706c696564206465706f7369745f646174615f726f6f744465706f736974436f6e74726163743a20696e76616c6964207769746864726177616c5f63726564656e7469616c73206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c7565206e6f74206d756c7469706c65206f6620677765694465706f736974436f6e74726163743a20696e76616c6964207075626b6579206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f20686967684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f206c6f774465706f736974436f6e74726163743a20696e76616c6964207369676e6174757265206c656e677468a26469706673582212201dd26f37a621703009abf16e77e69c93dc50c79db7f6cc37543e3e0e3decdc9764736f6c634300060b0033",
      "storage": {
        "0x0000000000000000000000000000000000000000000000000000000000000022": "0xf5a5fd42d16a20302798ef6ed309979b43003d2320d9f0e8ea9831a92759fb4b",
        "0x0000000000000000000000000000000000000000000000000000000000000023": "0xdb56114e00fdd4c1f85c892bf35ac9a89289aaecb1ebd0a96cde606a748b5d71",
        "0x0000000000000000000000000000000000000000000000000000000000000024": "0xc78009fdf07fc56a11f122370658a353aaa542ed63e44c4bc15ff4cd105ab33c",
        "0x0000000000000000000000000000000000000000000000000000000000000025": "0x536d98837f2dd165a55d5eeae91485954472d56f246df256bf3cae19352a123c",
        "0x0000000000000000000000000000000000000000000000000000000000000026": "0x9efde052aa15429fae05bad4d0b1d7c64da64d03d7a1854a588c2cb8430c0d30",
        "0x0000000000000000000000000000000000000000000000000000000000000027": "0xd88ddfeed400a8755596b21942c1497e114c302e6118290f91e6772976041fa1",
        "0x0000000000000000000000000000000000000000000000000000000000000028": "0x87eb0ddba57e35f6d286673802a4af5975e22506c7cf4c64bb6be5ee11527f2c",
        "0x0000000000000000000000000000000000000000000000000000000000000029": "0x26846476fd5fc54a5d43385167c95144f2643f533cc85bb9d16b782f8d7db193",
        "0x000000000000000000000000000000000000000000000000000000000000002a": "0x506d86582d252405b840018792cad2bf1259f1ef5aa5f887e13cb2f0094f51e1",
        "0x000000000000000000000000000000000000000000000000000000000000002b": "0xffff0ad7e659772f9534c195c815efc4014ef1e1daed4404c06385d11192e92b",
        "0x000000000000000000000000000000000000000000000000000000000000002c": "0x6cf04127db05441cd833107a52be852868890e4317e6a02ab47683aa75964220",
        "0x000000000000000000000000000000000000000000000000000000000000002d": "0xb7d05f875f140027ef5118a2247bbb84ce8f2f0f1123623085daf7960c329f5f",
        "0x000000000000000000000000000000000000000000000000000000000000002e": "0xdf6af5f5bbdb6be9ef8aa618e4bf8073960867171e29676f8b284dea6a08a85e",
        "0x000000000000000000000000000000000000000000000000000000000000002f": "0xb58d900f5e182e3c50ef74969ea16c7726c549757cc23523c369587da7293784",
        "0x0000000000000000000000000000000000000000000000000000000000000030": "0xd49a7502ffcfb0340b1d7885688500ca308161a7f96b62df9d083b71fcc8f2bb",
        "0x0000000000000000000000000000000000000000000000000000000000000031": "0x8fe6b1689256c0d385f42f5bbe2027a22c1996e110ba97c171d3e5948de92beb",
        "0x0000000000000000000000000000000000000000000000000000000000000032": "0x8d0d63c39ebade8509e0ae3c9c3876fb5fa112be18f905ecacfecb92057603ab",
        "0x0000000000000000000000000000000000000000000000000000000000000033": "0x95eec8b2e541cad4e91de38385f2e046619f54496c2382cb6cacd5b98c26f5a4",
        "0x0000000000000000000000000000000000000000000000000000000000000034": "0xf893e908917775b62bff23294dbbe3a1cd8e6cc1c35b4801887b646a6f81f17f",
        "0x0000000000000000000000000000000000000000000000000000000000000035": "0xcddba7b592e3133393c16194fac7431abf2f5485ed711db282183c819e08ebaa",
        "0x0000000000000000000000000000000000000000000000000000000000000036": "0x8a8d7fe3af8caa085a7639a832001457dfb9128a8061142ad0335629ff23ff9c",
        "0x0000000000000000000000000000000000000000000000000000000000000037": "0xfeb3c337d7a51a6fbf00b9e34c52e1c9195c969bd4e7a0bfd51d5c5bed9c1167",
        "0x0000000000000000000000000000000000000000000000000000000000000038": "0xe71f0aa83cc32edfbefa9f4d3e0174ca85182eec9f3a09f6a6c0df6377a510d7",
        "0x0000000000000000000000000000000000000000000000000000000000000039": "0x31206fa80a50bb6abe29085058f16212212a60eec8f049fecb92d8c8e0a84bc0",
        "0x000000000000000000000000000000000000000000000000000000000000003a": "0x21352bfecbeddde993839f614c3dac0a3ee37543f9b412b16199dc158e23b544",
        "0x000000000000000000000000000000000000000000000000000000000000003b": "0x619e312724bb6d7c3153ed9de791d764a366b389af13c58bf8a8d90481a46765",
        "0x000000000000000000000000000000000000000000000000000000000000003c": "0x7cdd2986268250628d0c10e385c58c6191e6fbe05191bcc04f133f2cea72c1c4",
        "0x000000000000000000000000000000000000000000000000000000000000003d": "0x848930bd7ba8cac54661072113fb278869e07bb8587f91392933374d017bcbe1",
        "0x000000000000000000000000000000000000000000000000000000000000003e": "0x8869ff2c22b28cc10510d9853292803328be4fb0e80495e8bb8d271f5b889636",
        "0x000000000000000000000000000000000000000000000000000000000000003f": "0xb5fe28e79f1b850f8658246ce9b6a1e7b49fc06db7143e8fe0b4f2b0c5523a5c",
        "0x0000000000000000000000000000000000000000000000000000000000000040": "0x985e929f70af28d0bdd1a90a808f977f597c7c778c489e98d3bd8910d31ac0f7"
      }
    },
    '''.strip('\n')

    # insert data with specific line
    insert_data_at_line(file_path, state_deposit_contract, 799)

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

    # sleep 10s
    time.sleep(10)

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
        run_command(['docker', 'compose', 'up', '-d', 'l1'], cwd=paths.ops_bedrock_dir, env={
            'PWD': paths.ops_bedrock_dir,
            'L1_RPC': paths.l1_rpc_url if paths.fork_public_network else '',
            'BLOCK_NUMBER': paths.block_number,
        }, check=True)

        wait_up(8545)
        wait_for_rpc_server('127.0.0.1:8545')

        print("L1 is restarted.")
    except subprocess.CalledProcessError as e:
        print(f"Cannot restart L1 with error: {e}")

def insert_data_at_line(file_path, new_data, line_number):
    with open(file_path, 'r', encoding='utf-8') as file:
        lines = file.readlines()

    lines.insert(line_number - 1, new_data + '\n')

    with open(file_path, 'w', encoding='utf-8') as file:
        file.writelines(lines)
