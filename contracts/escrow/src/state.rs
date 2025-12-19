use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{Addr, Coin};
use cw_storage_plus::Item;

/// 에스크로 상태 Enum
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub enum EscrowStatus {
    /// 초기 상태 (자금 미예치)
    Idle,
    /// 예치 완료 (진행 중)
    Funded,
    /// 자금 지급 완료 (종료)
    Released,
    /// 환불 완료 (종료)
    Refunded,
}

/// 에스크로 전체 상태 저장
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct State {
    /// 구매자 (컨트랙트 생성자)
    pub buyer: Addr,
    /// 판매자
    pub seller: Addr,
    /// 거래 금액
    pub amount: Coin,
    /// 만료 시간 (Unix Timestamp)
    pub expiration: u64,
    /// 현재 진행 상태
    pub status: EscrowStatus,
}

pub const STATE: Item<State> = Item::new("state");
