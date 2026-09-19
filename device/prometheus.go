package device

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Registry *prometheus.Registry
}

// Messung
type MessungGauges struct {
	F1 prometheus.Gauge
	F2 prometheus.Gauge
	Y1 prometheus.Gauge
	Y2 prometheus.Gauge
}

func (m *Messung) setGauges() {
	m.gauges.F1.Set(float64(m.F1))
	m.gauges.F2.Set(float64(m.F2))
	m.gauges.Y1.Set(float64(m.Y1 / 1000))
	m.gauges.Y2.Set(float64(m.Y2 / 1000))
}

const metricName = "ufg_messwerte"

func newMessungGauge(base string, serial string, ort string, messung string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricName,
		Help: "Messwerte eines FruitGuard Sensors",
		ConstLabels: prometheus.Labels{
			"base":    base,
			"serial":  serial,
			"ort":     ort,
			"messung": messung,
		},
	})
}

func (s *State) initMetrics() {
	fmt.Println("Initalisiere Metriken")
	s.Metrics = Metrics{}
	s.Metrics.Registry = prometheus.NewRegistry()

	// Initialize Messung Gauges
	for _, line := range s.Cfg.Lines {
		for _, serial := range line.Id {
			// Initialize the MessungGauges for this serial number and line
			id := fmt.Sprintf("%d", serial)
			g := MessungGauges{
				F1: newMessungGauge(s.Cfg.Base, id, line.Ort, "F1"),
				F2: newMessungGauge(s.Cfg.Base, id, line.Ort, "F2"),
				Y1: newMessungGauge(s.Cfg.Base, id, line.Ort, "Y1"),
				Y2: newMessungGauge(s.Cfg.Base, id, line.Ort, "Y2"),
			}

			s.Metrics.Registry.MustRegister(g.F1)
			s.Metrics.Registry.MustRegister(g.F2)
			s.Metrics.Registry.MustRegister(g.Y1)
			s.Metrics.Registry.MustRegister(g.Y2)
		}
	}
}

func (m *Metrics) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{Registry: m.Registry})
}
