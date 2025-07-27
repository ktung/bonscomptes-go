package service_test

import (
	"bonscomptes/domain"
	"bonscomptes/service"
	"bonscomptes/util"
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

	balances, err := service.CalculateBalances(expenses)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

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
	totalBalance := 0.0
	for _, balance := range balances {
		totalBalance += balance
	}
	if !util.IsZero(totalBalance) {
		t.Errorf("Total balance should be zero, got: %f", totalBalance)
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

	balances, err := service.CalculateBalances(expenses)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

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
	totalBalance := 0.0
	for _, balance := range balances {
		totalBalance += balance
	}
	if !util.IsZero(totalBalance) {
		t.Errorf("Total balance should be zero, got: %f", totalBalance)
	}
}

func TestCalculateBalances_IncorrectRatio(t *testing.T) {
	expenses := []domain.Expense{
		{User: "user1", Amount: 20, SplitRatios: []domain.SplitRatio{
			{User: "user1", Ratio: 1.0 / 3.0},
			{User: "user2", Ratio: 4.0 / 3.0},
			{User: "user3", Ratio: 1.0 / 3.0},
		}},
		{User: "user2", Amount: 10, SplitRatios: []domain.SplitRatio{
			{User: "user1", Ratio: 1.0 / 3.0},
			{User: "user2", Ratio: 1.0 / 3.0},
			{User: "user3", Ratio: 1.0 / 3.0},
		}},
	}

	balances, err := service.CalculateBalances(expenses)
	if err == nil {
		t.Fatal("Expected error for invalid split ratio, got none")
	}
	if balances != nil {
		t.Fatal("Expected nil balances for invalid split ratio, got non-nil")
	}
}

func TestCalculateSuggestedReimbursements(t *testing.T) {
	balances := map[string]float64{
		"user1": 111.173333,
		"user2": 28.163333,
		"user3": -139.336667,
	}

	suggestedReimbursements, err := service.CalculateSuggestedReimbursements(balances)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

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

func TestCalculateSuggestedReimbursements_IncorrectBalances(t *testing.T) {
	balances := map[string]float64{
		"user1": 100,
		"user2": 0,
		"user3": -20,
	}

	suggestedReimbursements, err := service.CalculateSuggestedReimbursements(balances)
	if err == nil {
		t.Fatal("Expected error for non-zero total balance, got none")
	}
	if suggestedReimbursements != nil {
		t.Fatal("Expected nil suggested reimbursements for non-zero total balance, got non-nil")
	}
}
