package main

import (
	"bonscomptes/domain"
	"bonscomptes/service"
	"fmt"
)

func main() {

	expenses := []domain.Expense{
		{User: "user1", Amount: 250.51, Description: "User 1 pays for all users", SplitRatios: []domain.SplitRatio{
			{User: "user1", Ratio: 1.0 / 3.0},
			{User: "user2", Ratio: 1.0 / 3.0},
			{User: "user3", Ratio: 1.0 / 3.0},
		}},
		{User: "user2", Amount: 152.5, Description: "User 2 pays for all users", SplitRatios: []domain.SplitRatio{
			{User: "user1", Ratio: 1.0 / 3.0},
			{User: "user2", Ratio: 1.0 / 3.0},
			{User: "user3", Ratio: 1.0 / 3.0},
		}},
		{User: "user2", Amount: 10, Description: "User 2 pays for User 3", SplitRatios: []domain.SplitRatio{
			{User: "user2", Ratio: 0},
			{User: "user3", Ratio: 1.0},
		}},
		{User: "user3", Amount: 20, Description: "User 3 split expense with User 1", SplitRatios: []domain.SplitRatio{
			{User: "user1", Ratio: 0.25},
			{User: "user3", Ratio: 0.75},
		}},
	}
	fmt.Printf("Expenses: %v\n", expenses)

	balances := service.CalculateBalances(expenses)
	fmt.Printf("Balances: %v\n", balances)

	balancesCopy := make(map[string]float64, len(balances))
	for k, v := range balances {
		balancesCopy[k] = v
	}

	suggestedReimbursements := service.CalculateSuggestedReimbursements(balancesCopy)
	fmt.Printf("Suggested reimbursements: %v\n", suggestedReimbursements)
}
