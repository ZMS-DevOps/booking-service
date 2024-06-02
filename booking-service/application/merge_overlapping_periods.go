package application

import (
	"github.com/ZMS-DevOps/booking-service/domain"
	"time"
)

func mergeOverlappingPeriods(periods []domain.UnavailabilityPeriod) []domain.UnavailabilityPeriod {
	if len(periods) <= 1 {
		return periods
	}

	sortPeriodsByStartTime(periods)

	var mergedPeriods []domain.UnavailabilityPeriod
	currentPeriod := periods[0]

	for i := 1; i < len(periods); i++ {
		if periods[i].Start.Before(currentPeriod.End) || periods[i].Start.Equal(currentPeriod.End) {
			currentPeriod.End = maxTime(currentPeriod.End, periods[i].End)
		} else {
			mergedPeriods = append(mergedPeriods, currentPeriod)
			currentPeriod = periods[i]
		}
	}

	mergedPeriods = append(mergedPeriods, currentPeriod)

	return mergedPeriods
}

func sortPeriodsByStartTime(periods []domain.UnavailabilityPeriod) {
	for i := 0; i < len(periods)-1; i++ {
		for j := 0; j < len(periods)-1-i; j++ {
			if periods[j].Start.After(periods[j+1].Start) {
				periods[j], periods[j+1] = periods[j+1], periods[j]
			}
		}
	}
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func insertPeriod(newPeriod *domain.UnavailabilityPeriod, periods []domain.UnavailabilityPeriod) []domain.UnavailabilityPeriod {
	periods = append(periods, *newPeriod)
	return mergeOverlappingPeriods(periods)
}

func removePeriod(toRemove domain.UnavailabilityPeriod, periods []domain.UnavailabilityPeriod) []domain.UnavailabilityPeriod {
	var result []domain.UnavailabilityPeriod

	for _, period := range periods {
		if toRemove.Start.After(period.End) || toRemove.End.Before(period.Start) {
			result = append(result, period)
		} else {
			if toRemove.Start.After(period.Start) && toRemove.End.Before(period.End) {
				result = append(result, domain.UnavailabilityPeriod{Start: period.Start, End: toRemove.Start})
				result = append(result, domain.UnavailabilityPeriod{Start: toRemove.End, End: period.End})
			} else if toRemove.Start.After(period.Start) && toRemove.Start.Before(period.End) {
				result = append(result, domain.UnavailabilityPeriod{Start: period.Start, End: toRemove.Start})
			} else if toRemove.End.After(period.Start) && toRemove.End.Before(period.End) {
				result = append(result, domain.UnavailabilityPeriod{Start: toRemove.End, End: period.End})
			}
		}
	}
	return result
}
