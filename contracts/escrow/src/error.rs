use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Already funded")]
    AlreadyFunded {},

    #[error("Not funded")]
    NotFunded {},

    #[error("Insufficient funds sent")]
    InsufficientFunds {},

    #[error("Escrow expired")]
    Expired {},

    #[error("Escrow not expired")]
    NotExpired {},
}
