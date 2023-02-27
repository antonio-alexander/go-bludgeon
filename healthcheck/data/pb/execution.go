package pb

import "github.com/antonio-alexander/go-bludgeon/healthcheck/data"

func FromHealthCheck(h *data.HealthCheck) *HealthCheck {
	if h == nil {
		return nil
	}
	return &HealthCheck{
		Time: h.Time,
	}
}

func ToHealthCheck(h *HealthCheck) *data.HealthCheck {
	if h == nil {
		return nil
	}
	return &data.HealthCheck{
		Time: h.GetTime(),
	}
}
