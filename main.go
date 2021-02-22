package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"

	. "loan-books/domain"
	"loan-books/util"
)

var (
	bankFacilityMap     map[int][]*Facility
	facilityCovenantMap map[int][]*Covenant

	allFacilities []*Facility
	allLoans      []*Loan
)

// init initializes by loading data from large dataset and setting up
// lobal variables with relevant data
func init() {
	var (
		banks      []*Bank
		covenants  []*Covenant
		facilities []*Facility
		loans      []*Loan

		err error
	)

	if banks, err = util.GetBanks("./data/large/banks.csv"); err != nil {
		panic(err)
	}
	// Prompt doesn't make use of banks (leaving this for now)
	_ = banks

	if covenants, err = util.GetCovenants("./data/large/covenants.csv"); err != nil {
		panic(err)
	}

	if facilities, err = util.GetFacilities("./data/large/facilities.csv"); err != nil {
		panic(err)
	}

	if loans, err = util.GetLoans("./data/large/loans.csv"); err != nil {
		panic(err)
	}
	allLoans = loans

	if allFacilities == nil {
		allFacilities = make([]*Facility, 0, len(facilities))
	}
	if bankFacilityMap == nil {
		bankFacilityMap = make(map[int][]*Facility)
	}
	for _, f := range facilities {
		if _, ok := bankFacilityMap[f.BankId]; !ok {
			bankFacilityMap[f.BankId] = make([]*Facility, 0, 0)
		}
		bankFacilityMap[f.BankId] = append(bankFacilityMap[f.BankId], f)
		allFacilities = append(allFacilities, f)
	}
	sort.SliceStable(allFacilities, func(i, j int) bool {
		return allFacilities[i].InterestRate < allFacilities[j].InterestRate
	})

	if facilityCovenantMap == nil {
		facilityCovenantMap = make(map[int][]*Covenant)
	}
	for _, c := range covenants {
		if c.FacilityId != 0 {
			if _, ok := facilityCovenantMap[c.FacilityId]; !ok {
				facilityCovenantMap[c.FacilityId] = make([]*Covenant, 0, 1)
			}
			facilityCovenantMap[c.FacilityId] = append(facilityCovenantMap[c.FacilityId], c)
		} else {
			for _, f := range bankFacilityMap[c.BankId] {
				if _, ok := facilityCovenantMap[f.Id]; !ok {
					facilityCovenantMap[f.Id] = make([]*Covenant, 0, 1)
				}
				facilityCovenantMap[f.Id] = append(facilityCovenantMap[f.Id], c)
			}
		}
	}
}

// getYield returns the yield for a specific loan and facility.
// yield is returned as an int, rounded to the nearest cent.
func getYield(loan *Loan, facility *Facility) int {
	y := ((1 - loan.DefaultLikelihood) * loan.InterestRate * float64(loan.Amount)) -
		(loan.DefaultLikelihood * float64(loan.Amount)) -
		(facility.InterestRate * float64(loan.Amount))

	y = math.Round(y)
	return int(y)
}

// fundLoan finds the facility to find a loan and calculates expected yield.
// fundLoan returns a pointers to a Assignment and Yield objects.
func fundLoan(loan *Loan) (*Assignment, *Yield) {
	for i, facility := range allFacilities {
		if loan.InterestRate < facility.InterestRate {
			continue
		}
		if loan.Amount > facility.Amount {
			continue
		}
		covenants := facilityCovenantMap[facility.Id]
		for _, covenant := range covenants {
			if loan.State == covenant.BannedState {
				break
			}
			if loan.DefaultLikelihood > covenant.MaxDefaultLikelihood {
				break
			}

			expectedYield := getYield(loan, facility)

			// Remove facility from local dataset when it no longer has loan
			// amount capacity left.
			facility.Amount = facility.Amount - loan.Amount
			if facility.Amount == 0 {
				allFacilities = append(allFacilities[:i-1], allFacilities[i+1:]...)
				delete(facilityCovenantMap, facility.Id)
			}

			return &Assignment{loan.Id, facility.Id}, &Yield{facility.Id, expectedYield}
		}
	}
	return nil, nil
}

func main() {
	// initialize CSV writers for assignments.csv and yields.csv files
	assignmentFile, err := os.Create("./data/results/assignments.csv")
	if err != nil {
		panic(err)
	}
	assignmentWriter := csv.NewWriter(assignmentFile)
	defer func() {
		assignmentWriter.Flush()
		assignmentFile.Close()
	}()
	if err = assignmentWriter.Write([]string{"loan_id", "facility_id"}); err != nil {
		panic(err)
	}

	yieldFile, err := os.Create("./data/results/yields.csv")
	if err != nil {
		panic(err)
	}
	yieldWriter := csv.NewWriter(yieldFile)
	defer func() {
		yieldWriter.Flush()
		yieldFile.Close()
	}()
	if err = yieldWriter.Write([]string{"facility_id", "expected_yield"}); err != nil {
		panic(err)
	}

	// for each loan, fund and calculate the yield and write the result into
	// assignments.csv and yields.csv files
	facilityYieldMap := make(map[int]int)
	for _, loan := range allLoans {
		assignment, yield := fundLoan(loan)

		// accumulate yield per facility
		if _, found := facilityYieldMap[yield.FacilityId]; !found {
			facilityYieldMap[yield.FacilityId] = 0
		}
		facilityYieldMap[yield.FacilityId] = facilityYieldMap[yield.FacilityId] + yield.ExpectedYield

		record := []string{fmt.Sprint(assignment.LoanId), fmt.Sprint(assignment.FacilityId)}
		if err = assignmentWriter.Write(record); err != nil {
			panic(err)
		}
	}
	for facilityId, expectedYield := range facilityYieldMap {
		record := []string{fmt.Sprint(facilityId), fmt.Sprint(expectedYield)}
		if err = yieldWriter.Write(record); err != nil {
			panic(err)
		}
	}
}
