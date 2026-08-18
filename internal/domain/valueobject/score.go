package valueobject

type Score struct {
	Value     float64 `json:"value"`
	MaxValue  float64 `json:"max_value"`
	Label     string  `json:"label"`
	Category  string  `json:"category"`
}

func NewScore(value, maxValue float64, label, category string) *Score {
	return &Score{
		Value:    value,
		MaxValue: maxValue,
		Label:    label,
		Category: category,
	}
}

func (s *Score) Percentage() float64 {
	if s.MaxValue == 0 {
		return 0
	}
	return (s.Value / s.MaxValue) * 100
}
