use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Key already exists")]
    KeyAlreadyExists {},

    #[error("Key not found")]
    KeyNotFound {},
}
