# Points 모듈 개발 가이드

## 개요

`Points` 모듈은 기업형 포인트 및 정산 시스템을 구현한 Cosmos SDK 커스텀 모듈입니다.

### 주요 기능

- 포인트 발행, 사용, 전송
- 거래 내역 관리
- 정산 요청 및 처리
- 포인트 잔액 조회

## 목차

1. [모듈 설계](#모듈-설계)
2. [생성 과정](#생성-과정)
3. [모듈 구조](#모듈-구조)
4. [메시지 타입](#메시지-타입)
5. [쿼리](#쿼리)
6. [사용 방법](#사용-방법)
7. [테스트](#테스트)

---

## 모듈 설계

### 기능 요구사항

#### 1. 포인트 관리
- **발행 (Issue)**: 관리자가 사용자에게 포인트 발급
- **사용 (Spend)**: 포인트로 결제/사용
- **전송 (Transfer)**: 사용자간 포인트 이동

#### 2. 정산 기능
- **정산 요청 (RequestSettlement)**: 가맹점이 포인트를 현금으로 정산
- 정산 승인/거부 (향후 구현 예정)
- 정산 내역 조회

#### 3. 조회 기능
- 포인트 잔액 조회
- 거래 내역 조회
- 정산 내역 조회

### 데이터 모델

```
PointBalance (Map)
├── address: string    - 계정 주소 (key)
└── balance: uint64    - 포인트 잔액

Transaction (List)
├── id: uint64         - 거래 ID (자동 생성)
├── sender: string     - 보내는 사람
├── recipient: string  - 받는 사람
├── amount: uint64     - 금액
├── txType: string     - 거래 타입 (issue/spend/transfer)
└── timestamp: int64   - 타임스탬프

Settlement (List)
├── id: uint64         - 정산 ID (자동 생성)
├── requester: string  - 요청자
├── amount: uint64     - 금액
├── status: string     - 상태 (pending/approved/rejected)
└── timestamp: int64   - 타임스탬프
```

---

## 생성 과정

### 1단계: 모듈 생성

```bash
ignite scaffold module points --yes
```

**생성되는 파일:**
- `x/points/` - 모듈 디렉토리
- `proto/scontract/points/` - Protobuf 정의
- 기본 keeper, types, module 파일들

### 2단계: 메시지 타입 추가

#### 2-1. IssuePoints 메시지

```bash
ignite scaffold message issue-points \
  recipient:string \
  amount:uint64 \
  reason:string \
  --module points \
  --yes
```

**생성되는 파일:**
- `x/points/keeper/msg_server_issue_points.go`
- `x/points/types/message_issue_points.go`
- `proto/scontract/points/v1/tx.proto` (수정)

#### 2-2. SpendPoints 메시지

```bash
ignite scaffold message spend-points \
  amount:uint64 \
  description:string \
  --module points \
  --yes
```

**생성되는 파일:**
- `x/points/keeper/msg_server_spend_points.go`
- `x/points/types/message_spend_points.go`

#### 2-3. TransferPoints 메시지

```bash
ignite scaffold message transfer-points \
  recipient:string \
  amount:uint64 \
  --module points \
  --yes
```

**생성되는 파일:**
- `x/points/keeper/msg_server_transfer_points.go`
- `x/points/types/message_transfer_points.go`

#### 2-4. RequestSettlement 메시지

```bash
ignite scaffold message request-settlement \
  amount:uint64 \
  --module points \
  --yes
```

**생성되는 파일:**
- `x/points/keeper/msg_server_request_settlement.go`
- `x/points/types/message_request_settlement.go`

### 3단계: 저장소 타입 추가

#### 3-1. PointBalance (Map 타입)

```bash
ignite scaffold map point-balance \
  address:string \
  balance:uint64 \
  --module points \
  --no-message \
  --yes
```

**특징:**
- Map 타입: key-value 저장소
- `address`가 키로 사용됨
- 각 계정의 포인트 잔액 저장

**생성되는 파일:**
- `proto/scontract/points/v1/point_balance.proto`
- `x/points/keeper/query_point_balance.go`
- `x/points/types/key_point_balance.go`

#### 3-2. Transaction (List 타입)

```bash
ignite scaffold list transaction \
  sender:string \
  recipient:string \
  amount:uint64 \
  txType:string \
  timestamp:int64 \
  --module points \
  --no-message \
  --yes
```

**특징:**
- List 타입: 순차적으로 저장되는 리스트
- 자동 증가하는 ID 생성
- 모든 거래 내역 저장

**생성되는 파일:**
- `proto/scontract/points/v1/transaction.proto`
- `x/points/keeper/query_transaction.go`

#### 3-3. Settlement (List 타입)

```bash
ignite scaffold list settlement \
  requester:string \
  amount:uint64 \
  status:string \
  timestamp:int64 \
  --module points \
  --no-message \
  --yes
```

**특징:**
- List 타입
- 정산 요청 기록 저장

**생성되는 파일:**
- `proto/scontract/points/v1/settlement.proto`
- `x/points/keeper/query_settlement.go`

---

## 모듈 구조

### 디렉토리 구조

```
x/points/
├── keeper/
│   ├── genesis.go                         - Genesis 초기화
│   ├── keeper.go                          - Keeper 정의
│   ├── msg_server.go                      - 메시지 서버
│   ├── msg_server_issue_points.go         - 포인트 발행 핸들러
│   ├── msg_server_spend_points.go         - 포인트 사용 핸들러
│   ├── msg_server_transfer_points.go      - 포인트 전송 핸들러
│   ├── msg_server_request_settlement.go   - 정산 요청 핸들러
│   ├── query.go                           - 쿼리 서버
│   ├── query_point_balance.go             - 잔액 조회 핸들러
│   ├── query_transaction.go               - 거래 내역 조회 핸들러
│   └── query_settlement.go                - 정산 내역 조회 핸들러
├── types/
│   ├── codec.go                           - 코덱 등록
│   ├── errors.go                          - 에러 정의
│   ├── expected_keepers.go                - 다른 모듈 인터페이스
│   ├── genesis.go                         - Genesis 타입
│   ├── keys.go                            - 저장소 키
│   ├── message_issue_points.go            - IssuePoints 메시지
│   ├── message_spend_points.go            - SpendPoints 메시지
│   ├── message_transfer_points.go         - TransferPoints 메시지
│   ├── message_request_settlement.go      - RequestSettlement 메시지
│   └── params.go                          - 모듈 파라미터
├── module/
│   ├── autocli.go                         - CLI 자동 생성
│   ├── depinject.go                       - 의존성 주입
│   ├── module.go                          - 모듈 정의
│   └── simulation.go                      - 시뮬레이션
└── simulation/
    ├── issue_points.go
    ├── spend_points.go
    ├── transfer_points.go
    └── request_settlement.go

proto/scontract/points/
├── module/v1/
│   └── module.proto                       - 모듈 설정
└── v1/
    ├── genesis.proto                      - Genesis 상태
    ├── params.proto                       - 파라미터
    ├── tx.proto                           - 트랜잭션 메시지
    ├── query.proto                        - 쿼리 정의
    ├── point_balance.proto                - 잔액 타입
    ├── transaction.proto                  - 거래 타입
    └── settlement.proto                   - 정산 타입
```

---

## 메시지 타입

### 1. IssuePoints

**목적:** 관리자가 사용자에게 포인트를 발행합니다.

**Protobuf 정의:**
```protobuf
message MsgIssuePoints {
  string creator = 1;      // 발행자 (서명자)
  string recipient = 2;    // 받는 사람 주소
  uint64 amount = 3;       // 발행할 포인트 양
  string reason = 4;       // 발행 사유
}
```

**CLI 사용법:**
```bash
scontractd tx points issue-points [recipient] [amount] [reason] \
  --from alice \
  --chain-id scontract \
  --yes
```

**예시:**
```bash
scontractd tx points issue-points cosmos1abc...xyz 1000 "welcome bonus" \
  --from admin \
  --chain-id scontract \
  --yes
```

**구현 위치:** `x/points/keeper/msg_server_issue_points.go`

### 2. SpendPoints

**목적:** 사용자가 포인트를 사용합니다.

**Protobuf 정의:**
```protobuf
message MsgSpendPoints {
  string creator = 1;       // 사용자 (서명자)
  uint64 amount = 2;        // 사용할 포인트 양
  string description = 3;   // 사용 내역
}
```

**CLI 사용법:**
```bash
scontractd tx points spend-points [amount] [description] \
  --from alice \
  --chain-id scontract \
  --yes
```

**예시:**
```bash
scontractd tx points spend-points 100 "coffee purchase" \
  --from alice \
  --chain-id scontract \
  --yes
```

**구현 위치:** `x/points/keeper/msg_server_spend_points.go`

### 3. TransferPoints

**목적:** 사용자간 포인트를 전송합니다.

**Protobuf 정의:**
```protobuf
message MsgTransferPoints {
  string creator = 1;     // 보내는 사람 (서명자)
  string recipient = 2;   // 받는 사람 주소
  uint64 amount = 3;      // 전송할 포인트 양
}
```

**CLI 사용법:**
```bash
scontractd tx points transfer-points [recipient] [amount] \
  --from alice \
  --chain-id scontract \
  --yes
```

**예시:**
```bash
scontractd tx points transfer-points cosmos1def...xyz 50 \
  --from alice \
  --chain-id scontract \
  --yes
```

**구현 위치:** `x/points/keeper/msg_server_transfer_points.go`

### 4. RequestSettlement

**목적:** 가맹점이 포인트를 현금으로 정산 요청합니다.

**Protobuf 정의:**
```protobuf
message MsgRequestSettlement {
  string creator = 1;   // 요청자 (서명자)
  uint64 amount = 2;    // 정산 요청 금액
}
```

**CLI 사용법:**
```bash
scontractd tx points request-settlement [amount] \
  --from merchant \
  --chain-id scontract \
  --yes
```

**예시:**
```bash
scontractd tx points request-settlement 5000 \
  --from merchant \
  --chain-id scontract \
  --yes
```

**구현 위치:** `x/points/keeper/msg_server_request_settlement.go`

---

## 쿼리

### 1. PointBalance 조회

**목적:** 특정 계정의 포인트 잔액을 조회합니다.

**CLI 사용법:**
```bash
scontractd query points show-point-balance [address] \
  --chain-id scontract
```

**예시:**
```bash
scontractd query points show-point-balance cosmos1abc...xyz
```

**응답 예시:**
```json
{
  "point_balance": {
    "address": "cosmos1abc...xyz",
    "balance": "1000"
  }
}
```

### 2. Transaction 조회

**목적:** 거래 내역을 조회합니다.

#### 전체 목록 조회
```bash
scontractd query points list-transaction \
  --chain-id scontract
```

#### 특정 거래 조회
```bash
scontractd query points show-transaction [id] \
  --chain-id scontract
```

**응답 예시:**
```json
{
  "transaction": {
    "id": "1",
    "sender": "cosmos1abc...xyz",
    "recipient": "cosmos1def...xyz",
    "amount": "100",
    "txType": "transfer",
    "timestamp": "1703001234"
  }
}
```

### 3. Settlement 조회

**목적:** 정산 내역을 조회합니다.

#### 전체 목록 조회
```bash
scontractd query points list-settlement \
  --chain-id scontract
```

#### 특정 정산 조회
```bash
scontractd query points show-settlement [id] \
  --chain-id scontract
```

**응답 예시:**
```json
{
  "settlement": {
    "id": "1",
    "requester": "cosmos1merchant...xyz",
    "amount": "5000",
    "status": "pending",
    "timestamp": "1703001234"
  }
}
```

---

## 사용 방법

### 전체 워크플로우

#### 1. 블록체인 시작

```bash
ignite chain serve
```

#### 2. 포인트 발행 (관리자)

```bash
# Alice에게 1000 포인트 발행
scontractd tx points issue-points \
  $(scontractd keys show alice -a) \
  1000 \
  "welcome bonus" \
  --from admin \
  --chain-id scontract \
  --yes
```

#### 3. 잔액 확인

```bash
# Alice의 잔액 조회
scontractd query points show-point-balance \
  $(scontractd keys show alice -a)
```

#### 4. 포인트 사용

```bash
# Alice가 100 포인트 사용
scontractd tx points spend-points 100 "coffee" \
  --from alice \
  --chain-id scontract \
  --yes
```

#### 5. 포인트 전송

```bash
# Alice가 Bob에게 50 포인트 전송
scontractd tx points transfer-points \
  $(scontractd keys show bob -a) \
  50 \
  --from alice \
  --chain-id scontract \
  --yes
```

#### 6. 정산 요청

```bash
# 가맹점이 5000 포인트 정산 요청
scontractd tx points request-settlement 5000 \
  --from merchant \
  --chain-id scontract \
  --yes
```

#### 7. 거래 내역 조회

```bash
# 모든 거래 내역 조회
scontractd query points list-transaction

# 정산 내역 조회
scontractd query points list-settlement
```

### 스크립트 예제

#### 전체 테스트 스크립트

```bash
#!/bin/bash

CHAIN_ID="scontract"
ADMIN="admin"
ALICE="alice"
BOB="bob"

echo "=== 1. 포인트 발행 ==="
scontractd tx points issue-points \
  $(scontractd keys show $ALICE -a) \
  1000 \
  "welcome bonus" \
  --from $ADMIN \
  --chain-id $CHAIN_ID \
  --yes

sleep 3

echo "=== 2. 잔액 확인 ==="
scontractd query points show-point-balance \
  $(scontractd keys show $ALICE -a)

sleep 1

echo "=== 3. 포인트 사용 ==="
scontractd tx points spend-points 100 "coffee" \
  --from $ALICE \
  --chain-id $CHAIN_ID \
  --yes

sleep 3

echo "=== 4. 포인트 전송 ==="
scontractd tx points transfer-points \
  $(scontractd keys show $BOB -a) \
  50 \
  --from $ALICE \
  --chain-id $CHAIN_ID \
  --yes

sleep 3

echo "=== 5. 최종 잔액 확인 ==="
echo "Alice:"
scontractd query points show-point-balance \
  $(scontractd keys show $ALICE -a)

echo "Bob:"
scontractd query points show-point-balance \
  $(scontractd keys show $BOB -a)

echo "=== 6. 거래 내역 ==="
scontractd query points list-transaction
```

---

## 테스트

### 1. 빌드 및 실행

```bash
# 프로토콜 버퍼 컴파일 및 빌드
ignite chain build

# 블록체인 시작
ignite chain serve
```

### 2. 단위 테스트

Ignite가 자동으로 생성한 테스트 파일들:
- `x/points/keeper/msg_server_issue_points_test.go`
- `x/points/keeper/msg_server_spend_points_test.go`
- `x/points/keeper/msg_server_transfer_points_test.go`
- `x/points/keeper/msg_server_request_settlement_test.go`
- `x/points/keeper/query_point_balance_test.go`
- `x/points/keeper/query_transaction_test.go`
- `x/points/keeper/query_settlement_test.go`

**테스트 실행:**
```bash
# 전체 테스트
go test ./x/points/...

# 특정 패키지 테스트
go test ./x/points/keeper

# verbose 모드
go test -v ./x/points/keeper
```

### 3. 통합 테스트

```bash
# 블록체인 초기화 후 재시작
ignite chain serve --reset-once

# 테스트 계정 생성 (자동으로 생성됨)
# - alice
# - bob

# 수동 테스트 실행
./test-points-module.sh
```

---

## 다음 단계

### 구현 필요 사항

현재 Ignite CLI로 생성된 파일들은 기본 스켈레톤만 제공합니다.
다음 로직을 직접 구현해야 합니다:

#### 1. `msg_server_issue_points.go`
```go
// TODO: 구현 필요
// 1. recipient의 잔액 조회
// 2. 잔액에 amount 추가
// 3. 거래 기록 저장
// 4. 이벤트 발생
```

#### 2. `msg_server_spend_points.go`
```go
// TODO: 구현 필요
// 1. creator의 잔액 조회
// 2. 잔액 충분한지 확인
// 3. 잔액에서 amount 차감
// 4. 거래 기록 저장
// 5. 이벤트 발생
```

#### 3. `msg_server_transfer_points.go`
```go
// TODO: 구현 필요
// 1. creator의 잔액 조회 및 확인
// 2. creator 잔액 차감
// 3. recipient 잔액 증가
// 4. 거래 기록 저장
// 5. 이벤트 발생
```

#### 4. `msg_server_request_settlement.go`
```go
// TODO: 구현 필요
// 1. creator의 잔액 조회 및 확인
// 2. 정산 요청 기록 저장 (status: pending)
// 3. 이벤트 발생
```

### 추가 기능

1. **정산 승인 메시지** (ApproveSettlement)
   ```bash
   ignite scaffold message approve-settlement \
     settlement-id:uint64 \
     approved:bool \
     --module points \
     --yes
   ```

2. **포인트 소각 메시지** (BurnPoints)
   ```bash
   ignite scaffold message burn-points \
     address:string \
     amount:uint64 \
     reason:string \
     --module points \
     --yes
   ```

3. **포인트 만료 처리**
   - EndBlocker에서 만료된 포인트 자동 소각

4. **권한 관리**
   - 관리자 권한 체크
   - 발행 한도 설정

---

## 참고 자료

- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Ignite CLI Documentation](https://docs.ignite.com/)
- [Protocol Buffers](https://protobuf.dev/)

## 트러블슈팅

### 빌드 오류

```bash
# 의존성 업데이트
go mod tidy

# 프로토 재생성
ignite generate proto-go --yes

# 클린 빌드
ignite chain build --clear-cache
```

### 모듈 등록 확인

`app/app.go`에서 Points 모듈이 올바르게 등록되었는지 확인:

```go
// Points 모듈이 포함되어 있어야 함
import (
    pointsmodule "scontract/x/points/module"
)
```

---

## 변경 이력

| 날짜 | 버전 | 설명 |
|------|------|------|
| 2025-12-19 | 1.0 | 초기 모듈 생성 (스켈레톤) |
