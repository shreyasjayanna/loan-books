package main

import (
	"fmt"

	. "loan-books/domain"
	"loan-books/util"
)

func getYield(loan *Loan, facility *Facility) float64 {
	return ((1 - loan.DefaultLikelihood) * loan.InterestRate * float64(loan.Amount)) -
		(loan.DefaultLikelihood * float64(loan.Amount)) -
		(facility.InterestRate * float64(loan.Amount))
}

func main() {
	var (
		banks      []*Bank
		covenants  []*Covenant
		facilities []*Facility
		loans      []*Loan

		err error
	)

	if banks, err = util.GetBanks("./data/small/banks.csv"); err != nil {
		panic(err)
	}

	if covenants, err = util.GetCovenants("./data/small/covenants.csv"); err != nil {
		panic(err)
	}
	if facilities, err = util.GetFacilities("./data/small/facilities.csv"); err != nil {
		panic(err)
	}
	if loans, err = util.GetLoans("./data/small/loans.csv"); err != nil {
		panic(err)
	}

	for _, b := range banks {
		fmt.Println(b)
	}

	for _, c := range covenants {
		fmt.Println(c)
	}

	for _, f := range facilities {
		fmt.Println(f)
	}

	for _, l := range loans {
		fmt.Println(l)
	}

	assignments := make([]*Assignment, 0, len(loans))
	yields := make([]*Yield, 0, len(facilities))

	_ = assignments
	_ = yields

}
