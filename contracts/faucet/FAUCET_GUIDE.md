# Faucet Contract 사용 가이드

이 문서는 Faucet Contract를 빌드, 배포하고 사용하는 전체 과정을 설명합니다.

## 1. 빌드 및 배포 (Build & Deploy)

### 1단계: 빌드 (Build)

WSL 환경의 파일 시스템 호환성 문제를 피하기 위해 Linux 네이티브 디렉토리에서 빌드하는 것을 권장합니다.

```bash
# 1. 임시 디렉토리 생성 및 소스 복사 (프로젝트 루트에서 실행)
mkdir -p ~/temp-build/faucet
cp -r contracts/faucet/* ~/temp-build/faucet/

# 2. 빌드 실행
cd ~/temp-build/faucet
cargo build --release --target wasm32-unknown-unknown

# 3. wasm 파일 프로젝트 경로로 복사
# cp target/wasm32-unknown-unknown/release/faucet.wasm /mnt/c/blockpj/custom-contract/faucet.wasm
cp ../../target/wasm32-unknown-unknown/release/faucet.wasm ../../faucet.wasm


# 4. 프로젝트 경로로 복귀
cd /mnt/c/blockpj/custom-contract
```

### 2단계: 컨트랙트 업로드 (Store)

**중요: `cp` 명령어가 성공했는지 반드시 확인하세요!**
만약 `No such file or directory` 에러가 났다면 `faucet.wasm` 파일이 제대로 복사되지 않은 것입니다. 이 상태로 진행하면 이전 파일(Simple Crud)이 올라갑니다.

올바른 복사 명령어 (프로젝트 루트 `custom-contract/contracts/faucet` 폴더 기준):
```bash
cp ../../target/wasm32-unknown-unknown/release/faucet.wasm ../../faucet.wasm
```

빌드된 `faucet.wasm` 파일을 블록체인에 저장합니다.

```bash
scontractd tx wasm store faucet.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json
```

*   **결과 확인**: 출력 JSON에서 `code_id`를 확인하세요. (이후 단계에서 사용)

### 3단계: 컨트랙트 생성 (Instantiate)

확인한 `code_id`를 사용하여 컨트랙트를 생성합니다. (예: `2`)

```bash
# CODE_ID를 실제 값으로 변경하세요
CODE_ID=2

scontractd tx wasm instantiate $CODE_ID \
  '{"cooldown_hours":24, "amount":{"denom":"uteam","amount":"1000000"}}' \
  --from alice \
  --label "faucet-v1" \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --admin alice \
  --yes
```

*   `cooldown_hours`: 24시간
*   `amount`: 청구 시 지급할 토큰 (1,000,000 uteam)
*   **결과 확인**: 트랜잭션 로그나 조회를 통해 생성된 **Contract Address**를 확인하세요.

**Contract Address 찾는 방법:**

만약 트랜잭션 로그에서 주소를 못 찾으셨다면, 다음 명령어로 생성된 컨트랙트 목록을 확인할 수 있습니다. 가장 최근(마지막)에 있는 주소가 방금 생성한 것입니다.

`$CODE_ID` 변수가 설정되어 있지 않다면, 실제 숫자(예: 2)를 직접 입력하세요.

```bash
# 예시: Code ID가 2인 경우
scontractd query wasm list-contract-by-code 2 --output json
```

### 4단계: Faucet에 자금 충전 (Funding)

Faucet이 사용자에게 토큰을 나눠주려면, 컨트랙트 자체가 먼저 토큰을 가지고 있어야 합니다.

`$CONTRACT_ADDR` {"contracts":["cosmos1wug8sewp6cedgkmrmvhl3lf3tulagm9hnvy8p0rppz9yjw0g4wtqdxfzff","cosmos1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqez7la9"],"pagination":{"next_key":null,"total":"0"}}

```bash
# CONTRACT_ADDR를 실제 컨트랙트 주소로 변경하세요
CONTRACT_ADDR=cosmos1wug8sewp6cedgkmrmvhl3lf3tulagm9hnvy8p0rppz9yjw0g4wtqdxfzff

scontractd tx bank send alice $CONTRACT_ADDR 100000000uteam \
  --chain-id scontract \
  --yes
```

## 2. 사용하기 (Interact)

### 토큰 청구 (Execute: Claim)

사용자(예: `bob`)가 토큰을 요청합니다.

```bash
scontractd tx wasm execute $CONTRACT_ADDR \
  '{"claim":{}}' \
  --from bob \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

*   성공 시: `uteam` 토큰이 `bob`에게 전송됩니다.
*   실패 시(쿨타임 중): `ClaimCooldownNotExpired` 에러가 발생합니다.

### 상태 조회 (Query)

**설정 조회:**

```bash
scontractd query wasm contract-state smart $CONTRACT_ADDR '{"get_config":{}}'
```

**마지막 청구 시간 조회:**

```bash
scontractd query wasm contract-state smart $CONTRACT_ADDR \
  '{"get_last_claim":{"address":"<USER_ADDRESS>"}}'
```
