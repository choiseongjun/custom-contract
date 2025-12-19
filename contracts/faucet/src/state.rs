use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{Addr, Coin};
use cw_storage_plus::{Item, Map};

/// Faucet 설정 정보를 저장하는 구조체
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct Config {
    /// 쿨타임 (초 단위)
    pub cooldown_seconds: u64,
    /// 지급할 토큰 정보
    pub amount: Coin,
}

/// 전역 설정을 저장소(Storage)에 저장하기 위한 키 ("config")
pub const CONFIG: Item<Config> = Item::new("config");

/// 사용자별 마지막 청구 시간을 저장하는 맵
/// Key: 사용자 주소 (&Addr)
/// Value: 마지막 청구 시각 (u64, Unix Timestamp)
pub const LAST_CLAIM: Map<&Addr, u64> = Map::new("last_claim");
