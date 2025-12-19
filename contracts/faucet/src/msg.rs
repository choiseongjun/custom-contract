use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::Coin;

/// 컨트랙트 생성 시 초기화 메시지
/// Faucet의 설정값을 정의합니다.
#[cw_serde]
pub struct InstantiateMsg {
    /// 재청구 대기 시간 (시간 단위)
    pub cooldown_hours: u64,
    /// 한 번 청구 시 지급할 토큰의 양과 종류
    pub amount: Coin,
}

/// 트랜잭션 실행 메시지 (Write 작업)
#[cw_serde]
pub enum ExecuteMsg {
    /// 토큰 청구 기능을 실행합니다.
    /// 쿨타임이 지나지 않았으면 실패합니다.
    Claim {},
}

/// 상태 조회 메시지 (Read 작업)
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// 현재 Faucet의 설정(쿨타임, 지급액)을 조회합니다.
    #[returns(ConfigResponse)]
    GetConfig {},
    /// 특정 주소의 마지막 청구 시간을 조회합니다.
    #[returns(LastClaimResponse)]
    GetLastClaim { address: String },
}

/// 설정 조회 응답 구조체
#[cw_serde]
pub struct ConfigResponse {
    /// 쿨타임 (초 단위)
    pub cooldown_seconds: u64,
    /// 지급할 토큰 정보
    pub amount: Coin,
}

/// 마지막 청구 시간 응답 구조체
#[cw_serde]
pub struct LastClaimResponse {
    /// 마지막 청구 시각 (Unix Timestamp, 초 단위)
    pub last_claim: u64,
}
