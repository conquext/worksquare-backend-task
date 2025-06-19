package pagination

import (
	"housing-api/internal/models"
	"math"
)

// CalculateMetadata calculates pagination metadata
func CalculateMetadata(page, limit int, total int64) models.MetaInfo {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return models.MetaInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}