package services

// RatingCategory represents the user's broad sentiment for a rated item.
// It determines the score range; binary search within the category determines
// the specific score. Scores are fully isolated between categories — adding
// songs to one category never affects scores in another.
type RatingCategory int

const (
	CategoryVeryBad RatingCategory = iota // 0.0 – 2.0
	CategoryBad                           // 2.0 – 4.0
	CategoryOk                            // 4.0 – 6.0
	CategoryGood                          // 6.0 – 8.0
	CategoryAmazing                       // 8.0 – 10.0
)

const categoryRange = 2.0

// CategoryMin returns the minimum (worst) score within a category.
func CategoryMin(c RatingCategory) float64 {
	return float64(c) * categoryRange
}

// CategoryMax returns the maximum (best) score within a category.
func CategoryMax(c RatingCategory) float64 {
	return CategoryMin(c) + categoryRange
}

// ScoreInWindow distributes comparison items evenly within a score window.
//
// Manual ratings act as fixed boundaries; comparison items are placed at equal
// intervals within [lower, upper]. Scores approach but never reach upper, so
// a comparison item can never equal a manual score ranked above it, and 10.0
// is only reachable via manual entry.
//
// position: 0-indexed within the window (0 = best in window)
// nTotal:   total comparison items in this window, including this one
// lower:    score of the nearest manual below (or category minimum)
// upper:    score of the nearest manual above (or category maximum)
func ScoreInWindow(position, nTotal int, lower, upper float64) float64 {
	if nTotal <= 1 {
		return (upper + lower) / 2.0
	}
	step := (upper - lower) / float64(nTotal+1)
	return upper - step*float64(position+1)
}

// SortedItem represents a rated item within a category, sorted by score.
type SortedItem struct {
	ID     string
	Method string  // "comparison" or "manual"
	Score  float64 // current stored score
}

// RecomputeCategoryScores recalculates scores for all comparison-method items
// in a category after the sorted order changes (e.g. a new item is inserted).
//
// items must be sorted best→worst (score descending) and may include both
// comparison and manual ratings. Manual ratings divide the score range into
// windows; comparison items within each window are rescored using ScoreInWindow.
// Manual scores are never modified.
//
// Returns a map of item ID → updated score for comparison-method items only.
func RecomputeCategoryScores(category RatingCategory, items []SortedItem) map[string]float64 {
	updated := make(map[string]float64)

	// Process items top-to-bottom. Each manual rating closes the current window
	// and opens a new one. Comparison items accumulate into the current window.
	currentUpper := CategoryMax(category)
	var windowItems []SortedItem

	scoreWindow := func(upper, lower float64, compItems []SortedItem) {
		nTotal := len(compItems)
		for position, item := range compItems {
			updated[item.ID] = ScoreInWindow(position, nTotal, lower, upper)
		}
	}

	for _, item := range items {
		if item.Method == "manual" {
			scoreWindow(currentUpper, item.Score, windowItems)
			currentUpper = item.Score
			windowItems = nil
		} else {
			windowItems = append(windowItems, item)
		}
	}
	// Score the final (bottom) window
	scoreWindow(currentUpper, CategoryMin(category), windowItems)

	return updated
}
