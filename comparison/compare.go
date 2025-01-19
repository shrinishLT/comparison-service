package comparison

import (
	"smartui-comparison-service/models"
)

func ProcessTask(task models.ComparisonTask) models.ComparisonResult {

	// process the comparison task over here

	result := models.ComparisonResult{
		RequestID: task.RequestID,
		Success:   true,
		Message:   "Comparison completed successfully",
	}

	return result
}
