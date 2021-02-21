package util

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	. "loan-books/domain"
)

func getFileContents(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	r := csv.NewReader(file)
	return r.ReadAll()
}

func GetBanks(filepath string) ([]*Bank, error) {
	records, err := getFileContents(filepath)
	if err != nil {
		return nil, err
	}

	banks := make([]*Bank, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			// skip csv header
			continue
		}

		var (
			id int

			err error
		)

		if id, err = strconv.Atoi(record[0]); err != nil {
			return nil, err
		}

		banks = append(banks, &Bank{Id: id, Name: record[1]})
	}

	return banks, nil
}

func GetCovenants(filepath string) ([]*Covenant, error) {
	records, err := getFileContents(filepath)
	if err != nil {
		return nil, err
	}

	covenants := make([]*Covenant, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			// skip csv header
			continue
		}

		var (
			facilityId           int
			bankId               int
			maxDefaultLikelihood float64

			err error
		)

		if record[0] != "" {
			if facilityId, err = strconv.Atoi(record[0]); err != nil {
				return nil, err
			}
		}

		if record[1] != "" {
			if maxDefaultLikelihood, err = strconv.ParseFloat(record[1], 64); err != nil {
				return nil, err
			}
		}

		if bankId, err = strconv.Atoi(record[2]); err != nil {
			return nil, err
		}

		covenants = append(covenants, &Covenant{
			FacilityId:           facilityId,
			MaxDefaultLikelihood: maxDefaultLikelihood,
			BankId:               bankId,
			BannedState:          record[3],
		})
	}

	return covenants, nil

}

func GetFacilities(filepath string) ([]*Facility, error) {
	records, err := getFileContents(filepath)
	if err != nil {
		return nil, err
	}

	facilities := make([]*Facility, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			// skip csv header
			continue
		}

		var (
			id           int
			bankId       int
			amount       int
			interestRate float64

			err error
		)

		amountComponents := strings.Split(record[0], ".")
		if amount, err = strconv.Atoi(amountComponents[0]); err != nil {
			return nil, err
		}

		if interestRate, err = strconv.ParseFloat(record[1], 64); err != nil {
			return nil, err
		}

		if id, err = strconv.Atoi(record[2]); err != nil {
			return nil, err
		}

		if bankId, err = strconv.Atoi(record[3]); err != nil {
			return nil, err
		}

		facilities = append(facilities, &Facility{
			Id:           id,
			BankId:       bankId,
			Amount:       amount,
			InterestRate: interestRate,
		})
	}

	return facilities, nil
}

func GetLoans(filepath string) ([]*Loan, error) {
	records, err := getFileContents(filepath)
	if err != nil {
		return nil, err
	}

	loans := make([]*Loan, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			// skip csv header
			continue
		}

		var (
			id                int
			interestRate      float64
			amount            int
			defaultLikelihood float64

			err error
		)

		if interestRate, err = strconv.ParseFloat(record[0], 64); err != nil {
			return nil, err
		}

		amountComponents := strings.Split(record[1], ".")
		if amount, err = strconv.Atoi(amountComponents[0]); err != nil {
			return nil, err
		}

		if id, err = strconv.Atoi(record[2]); err != nil {
			return nil, err
		}

		if defaultLikelihood, err = strconv.ParseFloat(record[3], 64); err != nil {
			return nil, err
		}

		loans = append(loans, &Loan{
			Id:                id,
			InterestRate:      interestRate,
			Amount:            amount,
			DefaultLikelihood: defaultLikelihood,
			State:             record[4],
		})
	}

	return loans, nil
}
