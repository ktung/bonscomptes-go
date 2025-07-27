package domain

type Expense struct {
	User        string
	Amount      float64
	Description string
	SplitRatios []SplitRatio
}

type SplitRatio struct {
	User  string
	Ratio float64
}
