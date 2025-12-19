# Simple CRUD Contract 실행 가이드

이 문서는 `contracts/simple-crud` 컨트랙트를 빌드하고 실행하는 전체 과정을 설명합니다.

## 사전 요구사항

- Rust 및 Cargo 설치
- wasm32-unknown-unknown 타겟 설치
- Ignite CLI 설치
- scontractd 바이너리 빌드 완료

```bash
# Rust wasm32 타겟 설치
rustup target add wasm32-unknown-unknown

# Ignite CLI 설치 (필요시)
curl https://get.ignite.com/cli! | bash
```

## 전체 실행 순서

### 1단계: 블록체인 노드 시작

```bash
# Ignite를 사용한 개발 모드 실행
ignite chain serve

# 또는 직접 실행
scontractd start
```

이 명령은 다음을 자동으로 수행합니다:
- 의존성 설치
- 바이너리 빌드
- 블록체인 초기화
- 로컬 노드 시작

### 2단계: 컨트랙트 빌드

새 터미널 창을 열고 다음 명령을 실행합니다:

#### 방법 1: 일반 빌드 (개발용)

```bash
# contracts/simple-crud 디렉토리로 이동
cd contracts/simple-crud

# 최적화된 wasm 빌드
cargo build --release --target wasm32-unknown-unknown

# 빌드 결과 확인
ls -lh target/wasm32-unknown-unknown/release/simple_crud.wasm
```

빌드가 성공하면 `simple_crud.wasm` 파일이 생성됩니다.

#### 방법 2: Docker 최적화 빌드 (프로덕션 권장)

Docker를 사용한 최적화 빌드는 다음 장점을 제공합니다:
- **파일 크기 최소화** - 일반 빌드 대비 50-70% 작은 크기
- **Deterministic 빌드** - 동일한 코드에서 항상 동일한 바이너리 생성
- **프로덕션 배포에 적합** - 메인넷 배포 시 필수

**단일 컨트랙트 최적화:**

```bash
# 프로젝트 루트에서 실행
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/optimizer:0.16.0 ./contracts/simple-crud

# 결과 확인 - artifacts 디렉토리에 생성됨
ls -lh artifacts/
```

**워크스페이스 전체 최적화 (여러 컨트랙트):**

```bash
# 프로젝트 루트에서 실행
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/workspace-optimizer:0.16.0

# 결과 확인
ls -lh artifacts/
```

최적화 빌드 결과:
- `artifacts/simple_crud.wasm` - 최적화된 wasm 파일
- `artifacts/checksums.txt` - 체크섬 파일 (검증용)

**빌드 크기 비교:**

```bash
# 일반 빌드 크기
ls -lh contracts/simple-crud/target/wasm32-unknown-unknown/release/simple_crud.wasm

# 최적화 빌드 크기 (보통 훨씬 작음)
ls -lh artifacts/simple_crud.wasm
```

**Docker 최적화 사용 시 3단계 명령어 수정:**

```bash
# 최적화된 파일로 저장
scontractd tx wasm store artifacts/simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json
```

### 3단계: 컨트랙트 저장 (Store)

컨트랙트 코드를 블록체인에 업로드합니다:

```bash
scontractd tx wasm store contracts/simple-crud/target/wasm32-unknown-unknown/release/simple_crud.wasm \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json
```

**출력에서 `code_id`를 확인하세요.** 이 값은 다음 단계에서 사용됩니다.

예시 출력:
```json
{
  "logs": [{
    "events": [{
      "type": "store_code",
      "attributes": [{
        "key": "code_id",
        "value": "1"
      }]
    }]
  }]
}
```

### 4단계: 컨트랙트 인스턴스화 (Instantiate)

저장된 코드로부터 컨트랙트 인스턴스를 생성합니다:

```bash
# CODE_ID를 이전 단계에서 받은 값으로 교체
scontractd tx wasm instantiate <CODE_ID> '{}' \
  --from alice \
  --label "simple-crud" \
  --no-admin \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes \
  --output json
```

**출력에서 `_contract_address`를 확인하세요.** 이 주소로 컨트랙트와 상호작용합니다.

### 5단계: CRUD 작업 테스트

#### Create (생성)

```bash
# CONTRACT_ADDRESS를 이전 단계에서 받은 값으로 교체
scontractd tx wasm execute <CONTRACT_ADDRESS> \
  '{"create":{"key":"name","value":"Alice"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

#### Read (조회)

```bash
scontractd query wasm contract-state smart <CONTRACT_ADDRESS> \
  '{"read":{"key":"name"}}' \
  --chain-id scontract \
  --output json
```

예상 출력:
```json
{
  "data": {
    "value": "Alice"
  }
}
```

#### Update (수정)

```bash
scontractd tx wasm execute <CONTRACT_ADDRESS> \
  '{"update":{"key":"name","value":"Bob"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

조회로 확인:
```bash
scontractd query wasm contract-state smart <CONTRACT_ADDRESS> \
  '{"read":{"key":"name"}}' \
  --chain-id scontract \
  --output json
```

#### Delete (삭제)

```bash
scontractd tx wasm execute <CONTRACT_ADDRESS> \
  '{"delete":{"key":"name"}}' \
  --from alice \
  --gas auto \
  --gas-adjustment 1.3 \
  --chain-id scontract \
  --yes
```

삭제 후 조회:
```bash
scontractd query wasm contract-state smart <CONTRACT_ADDRESS> \
  '{"read":{"key":"name"}}' \
  --chain-id scontract \
  --output json
```

삭제된 키는 `null` 값을 반환합니다.

## 자동 테스트 스크립트 실행

위의 모든 단계를 자동으로 실행하는 스크립트가 제공됩니다:

```bash
# 프로젝트 루트에서 실행
chmod +x test-contract.sh
./test-contract.sh
```

이 스크립트는 다음을 자동으로 수행합니다:
1. 컨트랙트 저장
2. 컨트랙트 인스턴스화
3. 데이터 생성 (Create)
4. 데이터 조회 (Read)
5. 데이터 수정 (Update)
6. 수정된 데이터 조회
7. 데이터 삭제 (Delete)
8. 삭제 후 조회

## 컨트랙트 구조

```
contracts/simple-crud/
├── Cargo.toml          # Rust 프로젝트 설정
├── src/
│   ├── contract.rs     # 메인 컨트랙트 로직
│   ├── msg.rs          # 메시지 정의 (ExecuteMsg, QueryMsg)
│   ├── state.rs        # 상태 저장소 정의
│   ├── error.rs        # 에러 타입 정의
│   └── lib.rs          # 라이브러리 엔트리포인트
└── target/             # 빌드 출력 디렉토리
```

## 주요 메시지 타입

### ExecuteMsg (실행 메시지)
- `Create { key: String, value: String }` - 새 데이터 생성
- `Update { key: String, value: String }` - 기존 데이터 수정
- `Delete { key: String }` - 데이터 삭제

### QueryMsg (조회 메시지)
- `Read { key: String }` - 키로 값 조회

## 트러블슈팅

### 1. 컨트랙트 빌드 실패
```bash
# 의존성 업데이트
cd contracts/simple-crud
cargo update
cargo clean
cargo build --release --target wasm32-unknown-unknown
```

### 2. Docker 최적화 빌드 에러

**"permission denied" 에러:**
```bash
# Linux/Mac: 현재 사용자 권한으로 실행
docker run --rm -v "$(pwd)":/code \
  --user $(id -u):$(id -g) \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/optimizer:0.16.0 ./contracts/simple-crud
```

**Windows (WSL) 파일시스템 문제:**

WSL의 `/mnt/c` 경로에서 빌드 시 `target` 디렉토리가 제대로 생성되지 않는 문제가 있습니다.

**해결 방법: 네이티브 Linux 디렉토리 사용**

```bash
# 1. 임시 빌드 디렉토리 생성
mkdir -p ~/temp-build
cp -r contracts/simple-crud ~/temp-build/

# 2. 네이티브 Linux 경로에서 빌드
cd ~/temp-build/simple-crud
cargo build --release --target wasm32-unknown-unknown

# 3. 생성된 wasm 파일 확인
ls -lh target/wasm32-unknown-unknown/release/simple_crud.wasm

# 4. 프로젝트 디렉토리로 복사
cp target/wasm32-unknown-unknown/release/simple_crud.wasm \
   /mnt/c/blockpj/custom-contract/simple_crud.wasm

# 5. 복사된 파일로 업로드
cd /mnt/c/blockpj/custom-contract
scontractd tx wasm store simple_crud.wasm --from alice --gas auto --gas-adjustment 1.3 --chain-id scontract --yes
```

**자세한 내용**: [CONTRACT_DEPLOYMENT_GUIDE.md](./CONTRACT_DEPLOYMENT_GUIDE.md) 참조

**Docker 최적화 시도 (WSL에서 제한적):**
```bash
# WSL에서 네이티브 디렉토리 사용
cd ~/temp-build
docker run --rm -v "$(pwd)":/code \
  cosmwasm/optimizer:0.16.0 ./simple-crud
```

**Docker 캐시 초기화:**
```bash
# 볼륨 캐시가 손상된 경우
docker volume rm $(basename "$(pwd)")_cache registry_cache
# 또는 모든 Docker 볼륨 삭제
docker volume prune
```

**최신 optimizer 이미지 사용:**
```bash
# 이미지 업데이트
docker pull cosmwasm/optimizer:0.16.0
docker pull cosmwasm/workspace-optimizer:0.16.0
```

### 3. 가스 부족 에러
`--gas-adjustment` 값을 늘려보세요 (예: 1.3 → 1.5)

### 4. 계정(alice) 없음
Ignite CLI로 체인을 시작하면 자동으로 테스트 계정이 생성됩니다.
수동으로 생성하려면:
```bash
scontractd keys add alice
```

### 5. 블록체인 초기화 재설정
```bash
ignite chain serve --reset-once
```

### 6. artifacts 디렉토리가 생성되지 않음
Docker 최적화 빌드는 반드시 **프로젝트 루트**에서 실행해야 합니다:
```bash
# 잘못된 예 (contracts/simple-crud 안에서 실행)
cd contracts/simple-crud
docker run ... # ❌ 작동하지 않음

# 올바른 예 (프로젝트 루트에서 실행)
cd /mnt/c/blockpj/custom-contract
docker run ... # ✅ artifacts/ 생성됨
```

## 추가 정보

- CosmWasm 버전: 2.1.4
- 최적화 빌드 설정: `Cargo.toml`의 `[profile.release]` 참조
- 체인 ID: `scontract`
- 기본 노드 RPC: `http://localhost:26657`

## 참고 자료

- [CosmWasm 공식 문서](https://docs.cosmwasm.com/)
- [Cosmos SDK 문서](https://docs.cosmos.network/)
- [Ignite CLI 문서](https://docs.ignite.com/)
