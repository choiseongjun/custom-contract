# Escrow Contract 사용 가이드

구매자(Buyer)와 판매자(Seller) 간의 안전한 거래를 위한 에스크로 컨트랙트입니다.

## 1. 개요 (Workflow)

1.  **Instantiate**: 구매자가 컨트랙트를 생성하고 판매자, 금액, 만료시간을 설정합니다.
2.  **Deposit**: 구매자가 설정한 금액을 컨트랙트에 예치합니다. (상태: `Funded`)
3.  **Release**: 물품 수령 후 구매자가 자금 지급을 승인하면 판매자에게 돈이 갑니다. (조건: 만료 전)
4.  **Refund**: 만약 거래가 성사되지 않고 시간이 지나면(만료), 구매자가 돈을 돌려받습니다.

## 2. 빌드 및 배포

### 빌드 (Linux 환경 권장)

```bash
cd /mnt/c/blockpj/custom-contract
# 임시 폴더에서 빌드
mkdir -p ~/temp-build/escrow
cp -r contracts/escrow/* ~/temp-build/escrow/
cd ~/temp-build/escrow

cargo build --release --target wasm32-unknown-unknown

# 결과물 복사
cp target/wasm32-unknown-unknown/release/escrow.wasm /mnt/c/blockpj/custom-contract/escrow.wasm
```

### 업로드 (Store)

```bash
cd /mnt/c/blockpj/custom-contract

scontractd tx wasm store escrow.wasm \
  --from alice \
  --gas auto --gas-adjustment 1.3 --chain-id scontract --yes \
  --output json | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value'
```

*   출력된 **Code ID**를 기억하세요. (예: `6`)

## 3. 사용 시나리오

### 상황 설정
*   **구매자**: `alice`
*   **판매자**: `bob`
*   **금액**: `1000000uteam`
*   **잠금 시간**: `60초` (테스트를 위해 짧게 설정)

### Step 1. 컨트랙트 생성 (Instantiate) - Alice 실행

```bash
# Code ID를 실제 값으로 변경
export CODE_ID=6
# cosmos1wug8sewp6cedgkmrmvhl3lf3tulagm9hnvy8p0rppz9yjw0g4wtqdxfzff
scontractd tx wasm instantiate $CODE_ID \
  '{"seller":"cosmos1wug8sewp6cedgkmrmvhl3lf3tulagm9hnvy8p0rppz9yjw0g4wtqdxfzff", "amount":{"denom":"uteam","amount":"1000000"}, "lock_time":60}' \
  --from alice \
  --label "escrow-test" \
  --gas auto --gas-adjustment 1.3 --chain-id scontract --admin alice --yes
```

*   **Contract Address**를 찾아서 환경변수에 저장하세요.
    ```bash
    export CONTRACT_ADDR=cosmos1wug8sewp6cedgkmrmvhl3lf3tulagm9hnvy8p0rppz9yjw0g4wtqdxfzff
    ```

### Step 2. 자금 예치 (Deposit) - Alice 실행

생성할 때 약속한 금액(`1000000uteam`)을 정확히 보냅니다.

```bash
scontractd tx wasm execute $CONTRACT_ADDR '{"deposit":{}}' \
  --amount 1000000uteam \
  --from alice \
  --gas auto --gas-adjustment 1.3 --chain-id scontract --yes
```

*   이제 상태는 `Funded`가 됩니다. 돈은 Alice 지갑에서 빠져나가 컨트랙트에 보관됩니다.

### Step 3-A. 정상 거래 완료 (Release) - Alice 실행

물건을 잘 받았다면, 만료 시간(60초) 내에 Alice가 판매자에게 돈을 보내줍니다.

```bash
scontractd tx wasm execute $CONTRACT_ADDR '{"release":{}}' \
  --from alice \
  --gas auto --gas-adjustment 1.3 --chain-id scontract --yes
```

*   성공 시: Bob에게 돈이 입금됩니다. (상태: `Released`)

### Step 3-B. 거래 취소/환불 (Refund) - Alice 실행

만약 60초가 지날 때까지 물건이 안 왔다면, Alice는 돈을 돌려받을 수 있습니다.

```bash
# 60초가 지난 후 실행 가능
scontractd tx wasm execute $CONTRACT_ADDR '{"refund":{}}' \
  --from alice \
  --gas auto --gas-adjustment 1.3 --chain-id scontract --yes
```

*   성공 시: Alice에게 돈이 다시 돌아옵니다. (상태: `Refunded`)
*   주의: 60초가 지나기 전에는 에러(`NotExpired`)가 납니다.

## 4. 상태 조회

```bash
scontractd query wasm contract-state smart $CONTRACT_ADDR '{"get_config":{}}'
```
