package data

import (
	"fmt"
	"math"
)

const (
	closingCostSellProp float64 = 0.078
)

func (s *Scenario) monthlyMoneyMadeRenting(months float64) []float64 {
	monthlyCompoundingRate := math.Pow(s.Assumptions.MarketGrowthYearly, 1.0/12)
	totalInvested := (s.downPayment() + s.House.ClosingCostBuy) * s.Assumptions.DownPaymentInvestedPropIfRenting
	out := make([]float64, 0, int(months))
	value := totalInvested
	for m := 0.0; m < months; m++ {
		value *= monthlyCompoundingRate
		out = append(out, value-totalInvested-s.Rent*m)
		fmt.Printf("months: %d, equity gain: %.0f, rent paid: %.0f\n", int(m), value-totalInvested, s.Rent*m)
	}
	return out
}

func (s *Scenario) monthlyMoneyMadeBuying(months float64) []float64 {
	equityGain := s.equityGain(months)
	interestPaid := s.totalInterestPaid(months)
	sellClosingCosts := s.sellClosingCosts(months)
	out := make([]float64, 0, int(months))
	for m := 0; m < int(months); m++ {
		maintenanceCost := s.House.MaintenanceMonthly * float64(m)
		out = append(out,
			equityGain[m]-maintenanceCost-interestPaid[m]-s.House.ClosingCostBuy-sellClosingCosts[m])
		fmt.Printf("months: %d, equity gain: %.0f, interest paid: %.0f, maintenance cost: %.0f, buyingClosingCosts: %.0f, sellingClosingCosts: %.0f\n", m, equityGain[m], interestPaid[m], maintenanceCost, s.House.ClosingCostBuy, sellClosingCosts[m])
	}
	return out
}

func (s *Scenario) equityGain(months float64) []float64 {
	compoundingRealEstateRate := math.Pow(s.Assumptions.RealEstateGrowthYearly, 1.0/12)
	// Most mortgages compound monthly
	compoundingInterestRate := s.House.InterestRate / 12
	equity := s.downPayment()
	monthlyPayment := PMT(s.loanPrincipal(), s.House.InterestRate/12.0, 12.0*s.House.LoanTermYears)
	principal := s.loanPrincipal()
	out := make([]float64, 0, int(months))
	for m := 0.0; m < months; m++ {
		equity *= compoundingRealEstateRate
		interestPayment := principal * compoundingInterestRate
		if m < 12.0*s.House.LoanTermYears {
			equity += (monthlyPayment - interestPayment)
		}
		principal = math.Max(0, principal-(monthlyPayment-interestPayment))
		fmt.Printf("RE equity gain = equity %f - downpayment %f - principal paid %f\n", equity, s.downPayment(), s.loanPrincipal()-principal)
		out = append(out, equity-s.downPayment()-(s.loanPrincipal()-principal))
	}
	return out
}

func (s *Scenario) totalInterestPaid(months float64) []float64 {
	monthlyPayment := PMT(s.loanPrincipal(), s.House.InterestRate/12.0, 12.0*s.House.LoanTermYears)
	interest := 0.0
	principal := s.loanPrincipal()
	out := make([]float64, 0, int(months))
	for m := 0.0; m < months; m++ {
		interestPayment := principal * s.House.InterestRate / 12.0
		principal = math.Max(0, principal-(monthlyPayment-interestPayment))
		interest += interestPayment
		out = append(out, interest)
	}
	return out
}

func (s *Scenario) sellClosingCosts(months float64) []float64 {
	compoundingRealEstateRate := math.Pow(s.Assumptions.RealEstateGrowthYearly, 1.0/12)
	houseValue := s.House.Cost
	out := make([]float64, 0, int(months))
	for m := 0.0; m < months; m++ {
		houseValue *= compoundingRealEstateRate
		cost := closingCostSellProp * houseValue
		out = append(out, cost)
	}
	return out
}

// Gets monthly/annual payment for an amortized loan
func PMT(principal, interest, term float64) float64 {
	return (principal * interest) / (1 - math.Pow((1+interest), -term))
}
