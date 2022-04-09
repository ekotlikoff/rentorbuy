package data

import (
	"encoding/json"
	"log"
)

type InputData struct {
	Scenarios []Scenario
}

type Assumptions struct {
	MarketGrowthYearly       float64
	RealEstateGrowthYearly   float64
	IncomeTaxStateAndFederal float64
}

type Scenario struct {
	Assumptions Assumptions
	House       House
	Rent        float64
}

type House struct {
	Address               string
	Cost                  float64
	DownPaymentProportion float64
	InterestRate          float64
	LoanTermYears         float64
	MaintenanceMonthly    float64
	ClosingCostBuy        float64
	ClosingCostSell       float64
}

func (s *Scenario) downPayment() float64 {
	return s.House.Cost * s.House.DownPaymentProportion
}

func (s *Scenario) loanPrincipal() float64 {
	return s.House.Cost - s.downPayment()
}

const defaultScenario = Scenario{
	Assumptions: Assumptions{
		MarketGrowthYearly:       1.06,
		RealEstateGrowthYearly:   1.028,
		IncomeTaxStateAndFederal: 0.3465,
	},
	House: House{
		DownPaymentProportion: 0.20,
		InterestRate:          0.0325,
		LoanTermYears:         30.0,
		ClosingCostBuy:        8500.0,
	},
	Rent: 3550.0,
}

func LoadScenario(d []byte) *Scenario {
	var s Scenario = defaultScenario
	err := json.Unmarshal(d, &s)
	if err != nil {
		log.Fatalf(err)
	}
	return i
}
