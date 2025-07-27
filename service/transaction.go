package service

import "bonscomptes/domain"

func CalculateBalances(expenses []domain.Expense) map[string]float64 {
	balances := make(map[string]float64)
	for _, expense := range expenses {
		for _, splitRatio := range expense.SplitRatios {
			split := expense.Amount * splitRatio.Ratio

			if expense.User == splitRatio.User {
				balances[splitRatio.User] -= (split - expense.Amount)
			} else {
				balances[splitRatio.User] -= split
			}
		}
	}
	return balances
}

// max negative balance (debtor) should reimburse max positive balance (creditor)
func CalculateSuggestedReimbursements(balances map[string]float64) []domain.SuggestedReimbursement {
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

	return suggestedReimbursements
}
