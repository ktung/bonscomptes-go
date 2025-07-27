package service_test

import (
	"bonscomptes/domain"
	"bonscomptes/service"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func round2(f float64) float64 {
	return math.Round(f*100) / 100
}

func TestCalculateBalances_SimpleExpenses(t *testing.T) {
	expenses := []domain.Expense{
		{User: "user1", Amount: 20, SplitRatios: []domain.SplitRatio{
			{User: "user1", Ratio: 1.0 / 3.0},
			{User: "user2", Ratio: 1.0 / 3.0},
			{User: "user3", Ratio: 1.0 / 3.0},
		}},
		{User: "user2", Amount: 10, SplitRatios: []domain.SplitRatio{
			{User: "user1", Ratio: 1.0 / 3.0},
			{User: "user2", Ratio: 1.0 / 3.0},
			{User: "user3", Ratio: 1.0 / 3.0},
		}},
	}

	balances := service.CalculateBalances(expenses)

	expectedBalances := map[string]float64{
		"user1": 10.0,
		"user2": 0.0,
		"user3": -10.0,
	}
	for user, expectedBalance := range expectedBalances {
		if balance, exists := balances[user]; !exists || round2(balance) != round2(expectedBalance) {
			t.Errorf("Expected balance for %s: %f, got: %f", user, expectedBalance, balance)
		}
	}
}

func TestCalculateBalances_ComplexExpenses(t *testing.T) {
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

	balances := service.CalculateBalances(expenses)

	expectedBalances := map[string]float64{
		"user1": 111.173333,
		"user2": 28.163333,
		"user3": -139.336667,
	}
	for user, expectedBalance := range expectedBalances {
		if balance, exists := balances[user]; !exists || round2(balance) != round2(expectedBalance) {
			t.Errorf("Expected balance for %s: %f, got: %f", user, expectedBalance, balance)
		}
	}
}

func TestCalculateSuggestedReimbursements(t *testing.T) {
	balances := map[string]float64{
		"user1": 111.173333,
		"user2": 28.163333,
		"user3": -139.336667,
	}

	suggestedReimbursements := service.CalculateSuggestedReimbursements(balances)

	expectedResult := []domain.SuggestedReimbursement{
		{
			From:   "user3",
			To:     "user1",
			Amount: 111.173333,
		},
		{
			From:   "user3",
			To:     "user2",
			Amount: 28.163333,
		},
	}
	assert.ElementsMatch(t, expectedResult, suggestedReimbursements)
}
