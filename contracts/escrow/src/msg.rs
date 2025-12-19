use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::Coin;

/// 에스크로 초기화 메시지
#[cw_serde]
pub struct InstantiateMsg {
    /// 판매자 주소
    pub seller: String,
    /// 거래 금액
    pub amount: Coin,
    /// 잠금 기간 (초 단위, 현재 시간 기준)
    pub lock_time: u64,
}

/// 트랜잭션 실행 메시지
#[cw_serde]
pub enum ExecuteMsg {
    /// 구매자가 자금을 예치 (한 번만 가능)
    Deposit {},
    /// 구매자가 자금을 판매자에게 지급 (조건: 예치됨 + 만료 전)
    Release {},
    /// 구매자가 자금을 환불 (조건: 예치됨 + 만료 후 + 릴리즈 안됨)
    Refund {},
}

/// 상태 조회 메시지
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// 현재 에스크로 설정 및 상태 조회
    #[returns(ConfigResponse)]
    GetConfig {},
}

/// 설정 조회 응답 구조체
#[cw_serde]
pub struct ConfigResponse {
    pub buyer: String,
    pub seller: String,
    pub amount: Coin,
    pub expiration: u64,
    pub status: String,
}
