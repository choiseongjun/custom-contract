use cosmwasm_std::{
    entry_point, to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, ReadResponse};
use crate::state::DATA;

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    Ok(Response::new().add_attribute("method", "instantiate"))
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::Create { key, value } => execute_create(deps, key, value),
        ExecuteMsg::Update { key, value } => execute_update(deps, key, value),
        ExecuteMsg::Delete { key } => execute_delete(deps, key),
    }
}

fn execute_create(deps: DepsMut, key: String, value: String) -> Result<Response, ContractError> {
    if DATA.has(deps.storage, key.clone()) {
        return Err(ContractError::KeyAlreadyExists {});
    }
    DATA.save(deps.storage, key.clone(), &value)?;
    Ok(Response::new()
        .add_attribute("action", "create")
        .add_attribute("key", key)
        .add_attribute("value", value))
}

fn execute_update(deps: DepsMut, key: String, value: String) -> Result<Response, ContractError> {
    if !DATA.has(deps.storage, key.clone()) {
        return Err(ContractError::KeyNotFound {});
    }
    DATA.save(deps.storage, key.clone(), &value)?;
    Ok(Response::new()
        .add_attribute("action", "update")
        .add_attribute("key", key)
        .add_attribute("value", value))
}

fn execute_delete(deps: DepsMut, key: String) -> Result<Response, ContractError> {
    if !DATA.has(deps.storage, key.clone()) {
        return Err(ContractError::KeyNotFound {});
    }
    DATA.remove(deps.storage, key.clone());
    Ok(Response::new()
        .add_attribute("action", "delete")
        .add_attribute("key", key))
}

#[entry_point]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::Read { key } => to_json_binary(&query_read(deps, key)?),
    }
}

fn query_read(deps: Deps, key: String) -> StdResult<ReadResponse> {
    let value = DATA.may_load(deps.storage, key)?;
    Ok(ReadResponse { value })
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{from_json};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg {};
        let info = mock_info("creator", &[]);
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());
    }

    #[test]
    fn create_read_update_delete() {
        let mut deps = mock_dependencies();
        let msg = InstantiateMsg {};
        let info = mock_info("creator", &[]);
        instantiate(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        // Create
        let msg = ExecuteMsg::Create {
            key: "name".to_string(),
            value: "Alice".to_string(),
        };
        let res = execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();
        assert_eq!(res.attributes[0].value, "create");

        // Read
        let msg = QueryMsg::Read {
            key: "name".to_string(),
        };
        let res = query(deps.as_ref(), mock_env(), msg).unwrap();
        let value: ReadResponse = from_json(&res).unwrap();
        assert_eq!(value.value, Some("Alice".to_string()));

        // Update
        let msg = ExecuteMsg::Update {
            key: "name".to_string(),
            value: "Bob".to_string(),
        };
        execute(deps.as_mut(), mock_env(), info.clone(), msg).unwrap();

        // Read updated
        let msg = QueryMsg::Read {
            key: "name".to_string(),
        };
        let res = query(deps.as_ref(), mock_env(), msg).unwrap();
        let value: ReadResponse = from_json(&res).unwrap();
        assert_eq!(value.value, Some("Bob".to_string()));

        // Delete
        let msg = ExecuteMsg::Delete {
            key: "name".to_string(),
        };
        execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Read after delete
        let msg = QueryMsg::Read {
            key: "name".to_string(),
        };
        let res = query(deps.as_ref(), mock_env(), msg).unwrap();
        let value: ReadResponse = from_json(&res).unwrap();
        assert_eq!(value.value, None);
    }
}
