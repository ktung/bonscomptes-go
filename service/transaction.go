package service

import (
	"bonscomptes/domain"
	"bonscomptes/util"
	"fmt"
	"math"
)

func CalculateBalances(expenses []domain.Expense) (map[string]float64, error) {
	balances := make(map[string]float64)
	for _, expense := range expenses {
		totalRatio := 0.0
		for _, splitRatio := range expense.SplitRatios {
			if splitRatio.Ratio < 0 || splitRatio.Ratio > 1 {
				return nil, fmt.Errorf("Invalid split ratio: %f for user: %s in expense: %s", splitRatio.Ratio, splitRatio.User, expense.Description)
			}
			totalRatio += splitRatio.Ratio
			split := expense.Amount * splitRatio.Ratio

			if expense.User == splitRatio.User {
				balances[splitRatio.User] -= (split - expense.Amount)
			} else {
				balances[splitRatio.User] -= split
			}
		}

		if totalRatio != 1.0 {
			return nil, fmt.Errorf("Total split ratio must equal 1.0, got: %f for expense: %s", totalRatio, expense.Description)
		}
	}
	return balances, nil
}

// max negative balance (debtor) should reimburse max positive balance (creditor)
func CalculateSuggestedReimbursements(balances map[string]float64) ([]domain.SuggestedReimbursement, error) {
	totalBalance := 0.0
	for _, balance := range balances {
		totalBalance += balance
	}
	if !util.IsZero(totalBalance) {
		return nil, fmt.Errorf("Total balance should be zero before calculating reimbursements, got %f", math.Abs(totalBalance))
	}

	suggestedReimbursements := make([]domain.SuggestedReimbursement, 0, len(balances)-1)
	for {
		maxCreditor := ""
		maxCreditorBalance := 0.0
		maxDebtor := ""
		maxDebtorBalance := 0.0
		for User, balance := range balances {
			if balance < maxDebtorBalance {
				maxDebtor = User
				maxDebtorBalance = balance
			} else if balance > maxCreditorBalance {
				maxCreditor = User
				maxCreditorBalance = balance
			}
		}

		if maxCreditorBalance <= 0 || maxDebtorBalance >= 0 {
			break
		}

		suggestedReimbursements = append(suggestedReimbursements, domain.SuggestedReimbursement{
			From:   maxDebtor,
			To:     maxCreditor,
			Amount: maxCreditorBalance,
		})

		balances[maxCreditor] -= maxCreditorBalance
		balances[maxDebtor] += maxCreditorBalance
	}

	return suggestedReimbursements, nil
}
