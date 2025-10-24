package enforceenumtypes

type MetricType int // want MetricType:"^MetricTypeUnset,MetricTypeGauge,MetricTypeSum$"

const (
	MetricTypeUnset MetricType = iota
	MetricTypeGauge
	MetricTypeSum
)

type NotConcerned int // want NotConcerned:"^NotConcernedOne,NotConcernedTwo$"

const (
	NotConcernedOne NotConcerned = iota
	NotConcernedTwo
)

func CheckMetricSwitch(v MetricType) {
	switch v { // want "^missing cases in switch of type enforceenumtypes.MetricType: enforceenumtypes.MetricTypeGauge, enforceenumtypes.MetricTypeSum$"
	case MetricTypeUnset:
	}
}

func CheckOtherSwitch(v NotConcerned) {
	switch v {
	case NotConcernedOne:
	}
}

func CheckMetricMap() {
	_ = map[MetricType]int{ // want "^missing keys in map of key type enforceenumtypes.MetricType: enforceenumtypes.MetricTypeGauge, enforceenumtypes.MetricTypeSum$"
		MetricTypeUnset: 0,
	}
}

func CheckOtherMap() {
	_ = map[NotConcerned]int{
		NotConcernedOne: 0,
	}
}
