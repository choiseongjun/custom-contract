use cosmwasm_std::{
    entry_point, to_json_binary, BankMsg, Binary, Deps, DepsMut, Env, MessageInfo, Response,
    StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ConfigResponse, ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{EscrowStatus, State, STATE};

const CONTRACT_NAME: &str = "crates.io:escrow";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    let state = State {
        buyer: info.sender,
        seller: deps.api.addr_validate(&msg.seller)?,
        amount: msg.amount,
        expiration: env.block.time.seconds() + msg.lock_time,
        status: EscrowStatus::Idle,
    };
    STATE.save(deps.storage, &state)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("buyer", state.buyer)
        .add_attribute("seller", state.seller)
        .add_attribute("expiration", state.expiration.to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::Deposit {} => execute_deposit(deps, info),
        ExecuteMsg::Release {} => execute_release(deps, env, info),
        ExecuteMsg::Refund {} => execute_refund(deps, env, info),
    }
}

/// 자금 예치 (Deposit)
/// 구매자가 실행하며, Instantiate 때 약속한 금액을 보내야 합니다.
fn execute_deposit(deps: DepsMut, info: MessageInfo) -> Result<Response, ContractError> {
    let mut state = STATE.load(deps.storage)?;

    // 1. 이미 예치되었는지 확인
    if state.status != EscrowStatus::Idle {
        return Err(ContractError::AlreadyFunded {});
    }

    // 2. 구매자가 맞는지 확인 (선택사항이지만 안전을 위해)
    if info.sender != state.buyer {
        return Err(ContractError::Unauthorized {});
    }

    // 3. 보낸 금액이 약속한 금액과 일치하는지 확인
    // coins는 리스트이므로 순회하며 확인하거나 find를 씁니다.
    // 여기서는 간단히 첫 번째 코인만 확인하거나, 전체 funds 중에서 찾습니다.
    let payment = info
        .funds
        .iter()
        .find(|c| c.denom == state.amount.denom)
        .ok_or(ContractError::InsufficientFunds {})?;

    if payment.amount < state.amount.amount {
        return Err(ContractError::InsufficientFunds {});
    }

    // 4. 상태 업데이트 (Idle -> Funded)
    state.status = EscrowStatus::Funded;
    STATE.save(deps.storage, &state)?;

    Ok(Response::new()
        .add_attribute("action", "deposit")
        .add_attribute("amount", payment.to_string()))
}

/// 자금 지급 (Release)
/// 구매자가 물건을 받고 확인(Release)하면 판매자에게 돈을 줍니다.
fn execute_release(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
) -> Result<Response, ContractError> {
    let mut state = STATE.load(deps.storage)?;

    // 1. 구매자만 실행 가능
    if info.sender != state.buyer {
        return Err(ContractError::Unauthorized {});
    }

    // 2. 예치된 상태인지 확인
    if state.status != EscrowStatus::Funded {
        return Err(ContractError::NotFunded {});
    }

    // 3. 만료 시간이 지났는지 확인 (만료되면 Refund만 가능)
    if env.block.time.seconds() >= state.expiration {
        return Err(ContractError::Expired {});
    }

    // 4. 상태 업데이트 (Funded -> Released)
    state.status = EscrowStatus::Released;
    STATE.save(deps.storage, &state)?;

    // 5. 판매자에게 송금
    let msg = BankMsg::Send {
        to_address: state.seller.to_string(),
        amount: vec![state.amount],
    };

    Ok(Response::new()
        .add_message(msg)
        .add_attribute("action", "release")
        .add_attribute("to", state.seller))
}

/// 자금 환불 (Refund)
/// 만료 시간이 지나면 구매자가 돈을 다시 찾아갈 수 있습니다.
fn execute_refund(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
) -> Result<Response, ContractError> {
    let mut state = STATE.load(deps.storage)?;

    // 1. 구매자만 실행 가능 (혹은 판매자도 가능하게 할 수 있음)
    if info.sender != state.buyer {
        return Err(ContractError::Unauthorized {});
    }

    // 2. 예치된 상태인지 확인
    if state.status != EscrowStatus::Funded {
        return Err(ContractError::NotFunded {});
    }

    // 3. 만료 시간이 지났는지 확인
    if env.block.time.seconds() < state.expiration {
        return Err(ContractError::NotExpired {});
    }

    // 4. 상태 업데이트 (Funded -> Refunded)
    state.status = EscrowStatus::Refunded;
    STATE.save(deps.storage, &state)?;

    // 5. 구매자에게 환불
    let msg = BankMsg::Send {
        to_address: state.buyer.to_string(),
        amount: vec![state.amount],
    };

    Ok(Response::new()
        .add_message(msg)
        .add_attribute("action", "refund")
        .add_attribute("to", state.buyer))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetConfig {} => to_json_binary(&query_config(deps)?),
    }
}

fn query_config(deps: Deps) -> StdResult<ConfigResponse> {
    let state = STATE.load(deps.storage)?;
    Ok(ConfigResponse {
        buyer: state.buyer.to_string(),
        seller: state.seller.to_string(),
        amount: state.amount,
        expiration: state.expiration,
        status: format!("{:?}", state.status),
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_json, Timestamp};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg {
            seller: "seller".to_string(),
            amount: coins(100, "token")[0].clone(),
            lock_time: 100,
        };
        let info = mock_info("buyer", &[]);
        instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetConfig {}).unwrap();
        let value: ConfigResponse = from_json(&res).unwrap();
        assert_eq!("buyer", value.buyer);
        assert_eq!("seller", value.seller);
        assert_eq!("Idle", value.status);
    }

    #[test]
    fn deposit_and_release() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg {
            seller: "seller".to_string(),
            amount: coins(100, "token")[0].clone(),
            lock_time: 100,
        };
        let info = mock_info("buyer", &[]);
        instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Deposit
        let info = mock_info("buyer", &coins(100, "token"));
        execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Deposit {}).unwrap();

        // Check Status
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetConfig {}).unwrap();
        let value: ConfigResponse = from_json(&res).unwrap();
        assert_eq!("Funded", value.status);

        // Release
        let info = mock_info("buyer", &[]);
        let res = execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Release {}).unwrap();
        assert_eq!(1, res.messages.len()); // Check BankMsg

        // Check Status
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetConfig {}).unwrap();
        let value: ConfigResponse = from_json(&res).unwrap();
        assert_eq!("Released", value.status);
    }

    #[test]
    fn refund_after_expiry() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg {
            seller: "seller".to_string(),
            amount: coins(100, "token")[0].clone(),
            lock_time: 100,
        };
        let info = mock_info("buyer", &[]);
        instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Deposit
        let info = mock_info("buyer", &coins(100, "token"));
        execute(deps.as_mut(), mock_env(), info, ExecuteMsg::Deposit {}).unwrap();

        // Try Refund before expiry -> Fail
        let info = mock_info("buyer", &[]);
        let err = execute(deps.as_mut(), mock_env(), info.clone(), ExecuteMsg::Refund {}).unwrap_err();
        match err {
            ContractError::NotExpired {} => {}
            _ => panic!("Expected NotExpired error"),
        }

        // Advance time
        let mut env = mock_env();
        env.block.time = Timestamp::from_seconds(env.block.time.seconds() + 200);

        // Refund after expiry -> Success
        let res = execute(deps.as_mut(), env, info, ExecuteMsg::Refund {}).unwrap();
        assert_eq!(1, res.messages.len());
    }
}
