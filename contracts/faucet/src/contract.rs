use cosmwasm_std::{
    entry_point, to_json_binary, BankMsg, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ConfigResponse, ExecuteMsg, InstantiateMsg, LastClaimResponse, QueryMsg};
use crate::state::{Config, CONFIG, LAST_CLAIM};

// 마이그레이션을 위한 컨트랙트 이름과 버전 정보
const CONTRACT_NAME: &str = "crates.io:faucet";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

/// 컨트랙트 생성(Instantiate) 진입점
/// 초기 설정을 저장하고 버전을 명시합니다.
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    // 컨트랙트 버전 저장 (나중에 마이그레이션 할 때 필요함)
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    // 입력받은 시간(시간 단위)을 초 단위로 변환
    let config = Config {
        cooldown_seconds: msg.cooldown_hours * 3600,
        amount: msg.amount,
    };
    // 설정 저장
    CONFIG.save(deps.storage, &config)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("cooldown_seconds", config.cooldown_seconds.to_string())
        .add_attribute("amount", config.amount.to_string()))
}

/// 트랜잭션 실행(Execute) 진입점
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        // Claim 메시지가 오면 execute_claim 함수 실행
        ExecuteMsg::Claim {} => execute_claim(deps, env, info),
    }
}

/// 토큰 청구 로직
/// 1. 쿨타임 확인
/// 2. 마지막 청구 시간 업데이트
/// 3. 토큰 전송 (BankMsg)
pub fn execute_claim(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
) -> Result<Response, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    // 보낸 사람(sender)의 마지막 청구 시간을 가져옴 (없으면 None)
    let last_claim = LAST_CLAIM.may_load(deps.storage, &info.sender)?;
    // 현재 블록 시간 (초)
    let now = env.block.time.seconds();

    // 이전에 청구한 기록이 있다면 쿨타임 체크
    if let Some(last) = last_claim {
        if now < last + config.cooldown_seconds {
            // 아직 쿨타임이 안 끝났으면 에러 반환
            let remaining = (last + config.cooldown_seconds) - now;
            return Err(ContractError::ClaimCooldownNotExpired { remaining });
        }
    }

    // 현재 시간으로 마지막 청구 시간 업데이트
    LAST_CLAIM.save(deps.storage, &info.sender, &now)?;

    // Bank모듈을 사용하여 사용자에게 토큰 전송 메시지 생성
    let msg = BankMsg::Send {
        to_address: info.sender.to_string(),
        amount: vec![config.amount],
    };

    Ok(Response::new()
        .add_message(msg) // 생성한 토큰 전송 메시지를 응답에 포함
        .add_attribute("action", "claim")
        .add_attribute("sender", info.sender)
        .add_attribute("timestamp", now.to_string()))
}

/// 조회(Query) 진입점
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetConfig {} => to_json_binary(&query_config(deps)?),
        QueryMsg::GetLastClaim { address } => to_json_binary(&query_last_claim(deps, address)?),
    }
}

/// 설정 조회
fn query_config(deps: Deps) -> StdResult<ConfigResponse> {
    let config = CONFIG.load(deps.storage)?;
    Ok(ConfigResponse {
        cooldown_seconds: config.cooldown_seconds,
        amount: config.amount,
    })
}

/// 마지막 청구 시간 조회
fn query_last_claim(deps: Deps, address: String) -> StdResult<LastClaimResponse> {
    let addr = deps.api.addr_validate(&address)?;
    // 기록이 없으면 0 반환
    let last_claim = LAST_CLAIM.may_load(deps.storage, &addr)?.unwrap_or(0);
    Ok(LastClaimResponse { last_claim })
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_json, Addr};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg {
            cooldown_hours: 24,
            amount: coins(10, "token")[0].clone(),
        };
        let info = mock_info("creator", &coins(1000, "token"));

        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetConfig {}).unwrap();
        let value: ConfigResponse = from_json(&res).unwrap();
        assert_eq!(86400, value.cooldown_seconds);
    }

    #[test]
    fn claim() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg {
            cooldown_hours: 1, // 1 hour for testing
            amount: coins(10, "token")[0].clone(),
        };
        let info = mock_info("creator", &coins(1000, "token"));
        instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // user can claim
        let info = mock_info("user", &[]);
        let res = execute(deps.as_mut(), mock_env(), info.clone(), ExecuteMsg::Claim {}).unwrap();
        assert_eq!(1, res.messages.len()); // BankMsg

        // check last claim
        let res = query(
            deps.as_ref(),
            mock_env(),
            QueryMsg::GetLastClaim {
                address: "user".to_string(),
            },
        )
        .unwrap();
        let value: LastClaimResponse = from_json(&res).unwrap();
        assert_ne!(0, value.last_claim);

        // user cannot claim again immediately
        let err = execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Claim {}).unwrap_err();
        match err {
            ContractError::ClaimCooldownNotExpired { .. } => {}
            e => panic!("unexpected error: {:?}", e),
        }
    }
}
