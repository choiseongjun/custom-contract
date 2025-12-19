use cosmwasm_schema::{cw_serde, QueryResponses};

#[cw_serde]
pub struct InstantiateMsg {}

#[cw_serde]
pub enum ExecuteMsg {
    Create { key: String, value: String },
    Update { key: String, value: String },
    Delete { key: String },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(ReadResponse)]
    Read { key: String },
}

#[cw_serde]
pub struct ReadResponse {
    pub value: Option<String>,
}
