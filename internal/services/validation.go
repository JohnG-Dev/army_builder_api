package services

import (
	"context"
	"fmt"

	"github.com/JohnG-Dev/army_builder_api/internal/models"
	"github.com/JohnG-Dev/army_builder_api/internal/state"
)

func ValidateArmy(s *state.State, ctx context.Context, req models.ArmyValidationRequest) (models.ValidationResponse, error) {
	resp := models.ValidationResponse{
		IsValid:     true,
		Errors:      []string{},
		TotalPoints: 0,
	}

	for _, u := range req.Units {
		currentID := u.UnitID
		currentQTY := u.Quantity

		unit, err := GetUnitByID(s, ctx, currentID)
		if err != nil {
			resp.Errors = append(resp.Errors, fmt.Sprintf("Unit ID %v, not found", u.UnitID))
			resp.IsValid = false
			continue
		}

		resp.TotalPoints += int(unit.Points) * currentQTY

		if unit.FactionID != req.FactionID {
			resp.Errors = append(resp.Errors, fmt.Sprintf("unit %s does not belong to the selected faction"), unit.Name)
		}

		if unit.IsUnique && currentQTY > 1 {
			resp.Errors = append(resp.Errors, fmt.Sprintf("Unit %s is unique and unable to have more than 1 in army"), unit.Name)
		}
	}

	if resp.TotalPoints > req.PointsLimit {
		msg := fmt.Sprintf("Total Points %d exceeds point limit %d", resp.TotalPoints, req.PointsLimit)
		resp.Errors = append(resp.Errors, msg)
	}

	if len(resp.Errors) > 0 {
		resp.IsValid = false
	}
	return resp, nil
}
