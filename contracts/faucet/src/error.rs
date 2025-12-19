use cosmwasm_std::StdError;
use thiserror::Error;

/// 컨트랙트에서 사용할 커스텀 에러 정의
#[derive(Error, Debug)]
pub enum ContractError {
    /// CosmWasm 표준 에러 포함
    #[error("{0}")]
    Std(#[from] StdError),

    /// 권한 없음 에러
    #[error("Unauthorized")]
    Unauthorized {},

    /// 쿨타임이 아직 끝나지 않았을 때 발생하는 에러
    /// {remaining} 변수를 통해 남은 시간을 알려줍니다.
    #[error("Claim cooldown not expired. Remaining: {remaining} seconds")]
    ClaimCooldownNotExpired { remaining: u64 },
}
