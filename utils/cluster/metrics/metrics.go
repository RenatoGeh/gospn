package metrics

type Metric func([]int, []int) float64
type MetricI func([]int, []int) int
type MetricF func([]float64, []float64) float64
