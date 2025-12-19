# CosmWasm 스마트 컨트랙트 배포 가이드

## 개요

이 문서는 `simple-crud` CosmWasm 스마트 컨트랙트를 빌드하고 `scontract` 블록체인에 배포하는 전체 과정을 설명합니다.

## 환경 정보

- **작업 디렉토리**: `/mnt/c/blockpj/custom-contract`
- **플랫폼**: WSL2 (Linux 6.6.87.2-microsoft-standard-WSL2)
- **블록체인**: scontract (Cosmos SDK 기반)
- **컨트랙트**: simple-crud

## 발생했던 문제

### 1. cosmwasm-optimizer 오류

```bash
docker run --rm \
  -v "$(pwd)":/code \
  cosmwasm/optimizer:0.16.0
```

**오류 내용**:
- `scontractd: command not found`
- WSL/Windows 파일시스템 마운트 문제로 wasm 파일 생성 실패

### 2. 타겟 디렉토리 문제

WSL의 `/mnt/c` 경로에서 Rust 빌드 시 `target` 디렉토리가 제대로 생성되지 않는 문제 발생

## 해결 과정

### 1단계: 블록체인 바이너리 빌드

```bash
make install
```

**결과**:
- 바이너리 위치: `/home/jun/go/bin/scontractd`
- PATH에 추가 필요

### 2단계: 네이티브 Linux 디렉토리에서 컨트랙트 빌드

WSL 파일시스템 문제를 우회하기 위해 네이티브 Linux 디렉토리 사용:

```bash
# 임시 디렉토리 생성 및 소스 복사
mkdir -p ~/temp-build
cp -r contracts/simple-crud ~/temp-build/

# 빌드
cd ~/temp-build/simple-crud
cargo build --release --target wasm32-unknown-unknown

# 빌드된 wasm 파일 프로젝트로 복사
cp ~/temp-build/simple-crud/target/wasm32-unknown-unknown/release/simple_crud.wasm \
   /mnt/c/blockpj/custom-contract/simple_crud.wasm
```

**빌드 결과**:
- 파일: `simple_crud.wasm`
- 크기: 241KB
- 위치: `/mnt/c/blockpj/custom-contract/simple_crud.wasm`

### 3단계: 블록체인에 업로드

```bash
/home/jun/go/bin/scontractd tx wasm store simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json
```

**업로드 결과**:
- Transaction Hash: `74957DD9717827D0B1FF0921C462A639E41559E885848493228D6901244CA1CB`
- Code ID: **1**
- Gas Used: 1,601,364
- Status: 성공

## 최종 결과

### 생성된 파일

| 파일 | 경로 | 설명 |
|------|------|------|
| scontractd | `/home/jun/go/bin/scontractd` | 블록체인 바이너리 |
| simple_crud.wasm | `/mnt/c/blockpj/custom-contract/simple_crud.wasm` | 컴파일된 스마트 컨트랙트 (241KB) |

### 업로드된 컨트랙트 정보

```bash
# 업로드된 코드 확인
/home/jun/go/bin/scontractd query wasm list-code
```

**출력**:
```yaml
code_infos:
- code_id: "1"
  creator: cosmos1hdx677fxqw5ec4qmppg0pk9ct9fevv4wvpw829
  data_hash: D673B89CE13BC5353B7C6381F6FB18318DDEA58E111DCC7D6D95164DFA98DE06
  instantiate_permission:
    addresses: []
    permission: Everybody
```

## 다음 단계: 컨트랙트 인스턴스화

### 1. 컨트랙트 인스턴스 생성

```bash
/home/jun/go/bin/scontractd tx wasm instantiate 1 '{}' \
  --from alice \
  --label "simple-crud-v1" \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json
```

### 2. 생성된 컨트랙트 주소 확인

```bash
# 인스턴스 목록 조회
/home/jun/go/bin/scontractd query wasm list-contract-by-code 1

# 특정 컨트랙트 정보 조회
/home/jun/go/bin/scontractd query wasm contract <CONTRACT_ADDRESS>
```

### 3. 컨트랙트와 상호작용

#### Store 메시지 실행 (데이터 저장)

```bash
/home/jun/go/bin/scontractd tx wasm execute <CONTRACT_ADDRESS> \
  '{"store":{"key":"test_key","value":"test_value"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

#### Query 실행 (데이터 조회)

```bash
/home/jun/go/bin/scontractd query wasm contract-state smart <CONTRACT_ADDRESS> \
  '{"get":{"key":"test_key"}}'
```

#### Delete 메시지 실행 (데이터 삭제)

```bash
/home/jun/go/bin/scontractd tx wasm execute <CONTRACT_ADDRESS> \
  '{"delete":{"key":"test_key"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

## 유용한 팁

### PATH 설정

매번 전체 경로를 입력하지 않으려면 `~/.bashrc`에 추가:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

이후 `scontractd` 명령어만으로 실행 가능합니다.

### 컨트랙트 재빌드 시

```bash
# 1. 네이티브 Linux 디렉토리에서 빌드
cd ~/temp-build/simple-crud
cargo build --release --target wasm32-unknown-unknown

# 2. wasm 파일 복사
cp target/wasm32-unknown-unknown/release/simple_crud.wasm \
   /mnt/c/blockpj/custom-contract/simple_crud.wasm

# 3. 새 버전으로 업로드 (새로운 code_id 생성)
scontractd tx wasm store simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

### 최적화된 빌드

프로덕션 환경에서는 wasm 파일 크기를 최소화하는 것이 중요합니다.

**방법 1**: 네이티브 Linux 환경에서 optimizer 사용

```bash
# optimizer를 Linux 디렉토리에서 실행
cd ~/temp-build
docker run --rm -v "$(pwd)":/code \
  cosmwasm/optimizer:0.16.0 ./simple-crud
```

**방법 2**: wasm-opt 직접 사용 (설치 필요)

```bash
# wasm-opt 설치
sudo apt-get install binaryen

# wasm 파일 최적화
wasm-opt -Oz \
  ~/temp-build/simple-crud/target/wasm32-unknown-unknown/release/simple_crud.wasm \
  -o simple_crud_optimized.wasm
```

## 문제 해결

### WSL 파일시스템 관련 문제

WSL의 `/mnt/c` 마운트 경로에서 빌드 시 문제가 발생할 수 있습니다. 이 경우 네이티브 Linux 디렉토리(`~/temp-build`)를 사용하세요.

### 트랜잭션 실패

```bash
# 계정 잔액 확인
scontractd query bank balances $(scontractd keys show alice -a)

# 가스 설정 조정
--gas auto --gas-adjustment 1.5
```

### 컨트랙트 에러 디버깅

```bash
# 트랜잭션 상세 조회
scontractd query tx <TX_HASH>

# 로그 확인
scontractd query tx <TX_HASH> --output json | jq '.raw_log'
```

## 참고 자료

- [CosmWasm Documentation](https://docs.cosmwasm.com/)
- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Rust Cargo Book](https://doc.rust-lang.org/cargo/)

## 변경 이력

| 날짜 | 버전 | 설명 |
|------|------|------|
| 2025-12-19 | 1.0 | 초기 배포 (Code ID: 1) |
