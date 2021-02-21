package domain

type Bank struct {
	Id   int    `csv:"bank_id"`
	Name string `csv:"bank_name"`
}

type Facility struct {
	Id           int     `csv:"facility_id"`
	BankId       int     `csv:"bank_id"`
	InterestRate float64 `csv:"interest_rate"`
	Amount       int     `csv:"amount"`
}

type Covenant struct {
	BankId               int     `csv:"bank_id"`
	FacilityId           int     `csv:"facility_id"`
	MaxDefaultLikelihood float64 `csv:"max_default_likelihood"`
	BannedState          string  `csv:"banned_state"`
}

type Loan struct {
	Id                int     `csv:"id"`
	Amount            int     `csv:"amount"`
	InterestRate      float64 `csv:"interest_rate"`
	DefaultLikelihood float64 `csv:"default_likelihood"`
	State             string  `csv:"state"`
}

type Assignment struct {
	LoanId     int `csv:"loan_id"`
	FacilityId int `csv:"facility_id"`
}

type Yield struct {
	FacilityId    int `csv:"facility_id"`
	ExpectedYield int `csv:"expected_yield"`
}
