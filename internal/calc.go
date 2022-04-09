package data

import "math"

const avgDaysPerMonth float64 = 30.43

func (s *Scenario) moneyLostRenting(days float64) float64 {
	totalBuyCost := s.downPayment() + s.House.ClosingCostBuy
	stockMarketGain :=
		math.Pow(s.Assumptions.MarketGrowthYearly, days/365.0)*totalBuyCost - totalBuyCost

	return stockMarketGain - s.Rent/avgDaysPerMonth*days
}

func (s *Scenario) moneyLostBuying(days float64) float64 {
	equityGain := s.equityGain(days)
	maintenanceCost := s.House.MaintenanceMonthly / avgDaysPerMonth * days
	interestPaid := s.totalInterestPaid(days)
	return equityGain - maintenanceCost - interestPaid - s.House.ClosingCostBuy - s.House.ClosingCostSell
}

func (s *Scenario) equityGain(days float64) {
	downPaymentGain :=
		s.downPayment()*math.Pow(s.Assumptions.RealEstateGrowthYearly, days/365.0) - s.downPayment()
	equity := s.downPayment()
	monthlyPayment := PMT(s.loanPrincipal(), s.House.InterestRate/12.0, 12.0*s.House.LoanTermYears)
	principal = s.loanPrincipal()
	for y := 0; y < days/365.0; y++ {
		equity *= s.Assumptions.RealEstateGrowthYearly
		for m := 0; m < 12 && principal > 0; m++ {
			interestPayment := principal * s.House.InterestRate / 12.0
			equity += (monthlyPayment - interestPayment) * s.Assumptions.RealEstateGrowthYearly * (12.0 - m)
			principal -= (monthlyPayment - interestPayment)
		}
	}
	return equity - s.downPayment() - (s.loanPrincipal() - principal)
}

func (s *Scenario) totalInterestPaid(days float64) float64 {
	monthlyPayment := PMT(s.loanPrincipal(), s.House.InterestRate/12.0, 12.0*s.House.LoanTermYears)
	out := 0.0
	principal = s.loanPrincipal()
	for m := 0; m < days/avgDaysPerMonth && principal > 0; m++ {
		interestPayment := principal * s.House.InterestRate / 12.0
		out += interestPayment
		principal -= (monthlyPayment - interestPayment)
	}
	return out
}

// Gets monthly/annual payment for an amortized loan
func PMT(principal, interest, term float64) float64 {
	return (principal * interest) / (1 - math.Pow((1+interest), -term))
}
