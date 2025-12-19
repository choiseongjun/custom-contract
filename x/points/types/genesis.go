package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		PointBalanceMap: []PointBalance{}, TransactionList: []Transaction{}, SettlementList: []Settlement{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	pointBalanceIndexMap := make(map[string]struct{})

	for _, elem := range gs.PointBalanceMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := pointBalanceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for pointBalance")
		}
		pointBalanceIndexMap[index] = struct{}{}
	}
	transactionIdMap := make(map[uint64]bool)
	transactionCount := gs.GetTransactionCount()
	for _, elem := range gs.TransactionList {
		if _, ok := transactionIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for transaction")
		}
		if elem.Id >= transactionCount {
			return fmt.Errorf("transaction id should be lower or equal than the last id")
		}
		transactionIdMap[elem.Id] = true
	}
	settlementIdMap := make(map[uint64]bool)
	settlementCount := gs.GetSettlementCount()
	for _, elem := range gs.SettlementList {
		if _, ok := settlementIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for settlement")
		}
		if elem.Id >= settlementCount {
			return fmt.Errorf("settlement id should be lower or equal than the last id")
		}
		settlementIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}
