# scontract

**scontract** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

This blockchain includes CosmWasm support for deploying and executing smart contracts.

## 📚 Documentation

- **[QUICK_REFERENCE.md](./QUICK_REFERENCE.md)** - 빠른 참조 가이드 (자주 사용하는 명령어)
- **[CONTRACT_DEPLOYMENT_GUIDE.md](./CONTRACT_DEPLOYMENT_GUIDE.md)** - 스마트 컨트랙트 빌드 및 배포 가이드
- **[SIMPLE_CRUD_GUIDE.md](./SIMPLE_CRUD_GUIDE.md)** - Simple CRUD 컨트랙트 사용법

## 🚀 Quick Start

### 블록체인 시작

```bash
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### 스마트 컨트랙트 배포

```bash
# 1. 블록체인 바이너리 빌드
make install

# 2. 컨트랙트 빌드 (네이티브 Linux 디렉토리 사용)
mkdir -p ~/temp-build
cp -r contracts/simple-crud ~/temp-build/
cd ~/temp-build/simple-crud
cargo build --release --target wasm32-unknown-unknown
cp target/wasm32-unknown-unknown/release/simple_crud.wasm /mnt/c/blockpj/custom-contract/

# 3. 컨트랙트 업로드
scontractd tx wasm store simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes

# 4. 컨트랙트 인스턴스화
scontractd tx wasm instantiate <CODE_ID> '{}' \
  --from alice \
  --label "simple-crud-v1" \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

자세한 내용은 [CONTRACT_DEPLOYMENT_GUIDE.md](./CONTRACT_DEPLOYMENT_GUIDE.md)를 참조하세요.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers a frontend scaffolding feature (based on Vue) to help you quickly build a web frontend for your blockchain:

Use: `ignite scaffold vue`
This command can be run within your scaffolded blockchain project.


For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/scontract@latest! | sudo bash
```
`username/scontract` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/ignite/installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.com/invite/ignitecli)
