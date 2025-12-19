# 빠른 참조 가이드 (Quick Reference)

## 자주 사용하는 명령어

### 블록체인 관련

```bash
# 블록체인 시작 (개발 모드)
ignite chain serve

# 블록체인 초기화 후 재시작
ignite chain serve --reset-once

# 바이너리 빌드
make install

# PATH 설정 (영구적)
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### 계정 관리

```bash
# 계정 목록 조회
scontractd keys list

# 새 계정 생성
scontractd keys add <account-name>

# 계정 주소 확인
scontractd keys show alice -a

# 계정 잔액 조회
scontractd query bank balances $(scontractd keys show alice -a)
```

### 컨트랙트 빌드

```bash
# WSL 환경: 네이티브 Linux 디렉토리 사용
mkdir -p ~/temp-build
cp -r contracts/simple-crud ~/temp-build/
cd ~/temp-build/simple-crud
cargo build --release --target wasm32-unknown-unknown
cp target/wasm32-unknown-unknown/release/simple_crud.wasm \
   /mnt/c/blockpj/custom-contract/simple_crud.wasm

# 일반 환경: 프로젝트 디렉토리에서 직접 빌드
cd contracts/simple-crud
cargo build --release --target wasm32-unknown-unknown
```

### 컨트랙트 배포

```bash
# 1. 컨트랙트 업로드
scontractd tx wasm store simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes

# 2. 업로드된 코드 목록 확인
scontractd query wasm list-code

# 3. 컨트랙트 인스턴스화 (CODE_ID는 위에서 확인한 값 사용)
scontractd tx wasm instantiate <CODE_ID> '{}' \
  --from alice \
  --label "simple-crud-v1" \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes

 scontractd tx wasm instantiate 2 '{}' \
    --from alice \
    --label "simple-crud-v1" \
    --no-admin \
    --gas auto \
    --gas-adjustment 1.3 \
    --chain-id scontract \
    --yes


# 4. 인스턴스 목록 조회
scontractd query wasm list-contract-by-code <CODE_ID>
```

### CRUD 작업

**변수 설정 (편의를 위해)**
```bash
CONTRACT=<계약주소>
```

#### Create (생성)
```bash
scontractd tx wasm execute $CONTRACT \
  '{"create":{"key":"name","value":"Alice"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

#### Read (조회)
```bash
scontractd query wasm contract-state smart $CONTRACT \
  '{"read":{"key":"name"}}' \
  --output json
```

#### Update (수정)
```bash
scontractd tx wasm execute $CONTRACT \
  '{"update":{"key":"name","value":"Bob"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

#### Delete (삭제)
```bash
scontractd tx wasm execute $CONTRACT \
  '{"delete":{"key":"name"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

### 트랜잭션 조회

```bash
# 트랜잭션 해시로 조회
scontractd query tx <TX_HASH>

# JSON 형식으로 출력
scontractd query tx <TX_HASH> --output json

# 로그만 확인
scontractd query tx <TX_HASH> --output json | jq '.raw_log'
```

### 컨트랙트 상태 조회

```bash
# 컨트랙트 정보
scontractd query wasm contract $CONTRACT

# 컨트랙트 전체 상태
scontractd query wasm contract-state all $CONTRACT

# 특정 키 조회 (raw)
scontractd query wasm contract-state raw $CONTRACT <hex-encoded-key>
```

### 디버깅

```bash
# 블록체인 로그 확인
tail -f ~/.scontract/logs/scontract.log

# 트랜잭션 실패 원인 확인
scontractd query tx <TX_HASH> --output json | jq '.raw_log'

# 가스 예측 (실제 실행하지 않음)
scontractd tx wasm execute $CONTRACT '{"create":{"key":"test","value":"test"}}' \
  --from alice \
  --gas auto \
  --dry-run

# 체인 상태 확인
scontractd status
```

### 파일 크기 확인

```bash
# 빌드된 wasm 파일 크기
ls -lh simple_crud.wasm

# 사람이 읽기 쉬운 형식으로
du -h simple_crud.wasm
```

## 환경 변수 설정

`.bashrc` 또는 `.zshrc`에 추가하면 편리합니다:

```bash
# scontractd PATH 추가
export PATH="$HOME/go/bin:$PATH"

# 자주 사용하는 변수
export CHAIN_ID="scontract"
export NODE="tcp://localhost:26657"
export KEYRING_BACKEND="test"

# 별칭 (alias) 설정
alias scd='scontractd'
alias scdq='scontractd query'
alias scdt='scontractd tx'
alias scdqw='scontractd query wasm'
alias scdtw='scontractd tx wasm'
```

설정 후 적용:
```bash
source ~/.bashrc  # 또는 source ~/.zshrc
```

## 스크립트 예제

### 전체 배포 스크립트

```bash
#!/bin/bash

# 빌드
cd ~/temp-build/simple-crud
cargo build --release --target wasm32-unknown-unknown
cp target/wasm32-unknown-unknown/release/simple_crud.wasm /mnt/c/blockpj/custom-contract/

# 프로젝트 디렉토리로 이동
cd /mnt/c/blockpj/custom-contract

# 업로드
UPLOAD_TX=$(scontractd tx wasm store simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json)

# CODE_ID 추출
CODE_ID=$(echo $UPLOAD_TX | jq -r '.logs[0].events[] | select(.type=="store_code") | .attributes[] | select(.key=="code_id") | .value')

echo "Uploaded with CODE_ID: $CODE_ID"

# 인스턴스화
INIT_TX=$(scontractd tx wasm instantiate $CODE_ID '{}' \
  --from alice \
  --label "simple-crud-v1" \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json)

# 컨트랙트 주소 추출
CONTRACT=$(echo $INIT_TX | jq -r '.logs[0].events[] | select(.type=="instantiate") | .attributes[] | select(.key=="_contract_address") | .value')

echo "Contract instantiated at: $CONTRACT"
echo "Export this for later use:"
echo "export CONTRACT=$CONTRACT"
```

### CRUD 테스트 스크립트

```bash
#!/bin/bash

CONTRACT=$1

if [ -z "$CONTRACT" ]; then
  echo "Usage: $0 <contract-address>"
  exit 1
fi

echo "=== CREATE ==="
scontractd tx wasm execute $CONTRACT \
  '{"create":{"key":"test","value":"hello"}}' \
  --from alice --gas auto --gas-adjustment 1.3 --chain-id scontract --yes

sleep 1

echo "=== READ ==="
scontractd query wasm contract-state smart $CONTRACT \
  '{"read":{"key":"test"}}' --output json | jq

echo "=== UPDATE ==="
scontractd tx wasm execute $CONTRACT \
  '{"update":{"key":"test","value":"world"}}' \
  --from alice --gas auto --gas-adjustment 1.3 --chain-id scontract --yes

sleep 1

echo "=== READ AFTER UPDATE ==="
scontractd query wasm contract-state smart $CONTRACT \
  '{"read":{"key":"test"}}' --output json | jq

echo "=== DELETE ==="
scontractd tx wasm execute $CONTRACT \
  '{"delete":{"key":"test"}}' \
  --from alice --gas auto --gas-adjustment 1.3 --chain-id scontract --yes

sleep 1

echo "=== READ AFTER DELETE ==="
scontractd query wasm contract-state smart $CONTRACT \
  '{"read":{"key":"test"}}' --output json | jq
```

## 문제 해결 체크리스트

- [ ] 블록체인이 실행 중인가? (`scontractd status`)
- [ ] 계정에 잔액이 있는가? (`scontractd query bank balances $(scontractd keys show alice -a)`)
- [ ] wasm 파일이 존재하는가? (`ls -lh simple_crud.wasm`)
- [ ] 올바른 CODE_ID를 사용했는가? (`scontractd query wasm list-code`)
- [ ] 올바른 CONTRACT 주소를 사용했는가? (`scontractd query wasm list-contract-by-code <CODE_ID>`)
- [ ] JSON 형식이 올바른가? (작은따옴표로 감싸고, 키는 큰따옴표 사용)
- [ ] WSL 환경에서 네이티브 Linux 디렉토리를 사용했는가?

## 유용한 리소스

| 문서 | 설명 |
|------|------|
| [CONTRACT_DEPLOYMENT_GUIDE.md](./CONTRACT_DEPLOYMENT_GUIDE.md) | 상세한 배포 가이드 |
| [SIMPLE_CRUD_GUIDE.md](./SIMPLE_CRUD_GUIDE.md) | CRUD 컨트랙트 사용법 |
| [readme.md](./readme.md) | 프로젝트 개요 |

## JSON 메시지 템플릿

### Execute 메시지
```json
// Create
{"create":{"key":"KEY","value":"VALUE"}}

// Update
{"update":{"key":"KEY","value":"VALUE"}}

// Delete
{"delete":{"key":"KEY"}}
```

### Query 메시지
```json
// Read
{"read":{"key":"KEY"}}
```

## 자주 발생하는 오류

| 오류 | 원인 | 해결 방법 |
|------|------|----------|
| `command not found: scontractd` | PATH 설정 안됨 | `export PATH="$HOME/go/bin:$PATH"` |
| `account not found` | 계정 미생성 | `scontractd keys add alice` |
| `insufficient funds` | 잔액 부족 | 개발 모드에서 자동으로 충전됨 (체인 재시작 필요시 `ignite chain serve --reset-once`) |
| `out of gas` | 가스 부족 | `--gas-adjustment` 값을 1.5로 증가 |
| `target directory not found` | WSL 파일시스템 문제 | 네이티브 Linux 디렉토리(`~/temp-build`)에서 빌드 |
| `code not found` | 잘못된 CODE_ID | `scontractd query wasm list-code`로 확인 |
| `contract not found` | 잘못된 CONTRACT 주소 | `scontractd query wasm list-contract-by-code <CODE_ID>`로 확인 |
