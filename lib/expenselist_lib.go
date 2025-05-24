package lib

import (
	"cmp"
	"slices"

	"github.com/fabianofski/equaly-backend/models"
)

func calculate_compensations(shares []*models.ExpenseListShare) []models.ExpenseListCompensation {
	compensations := []models.ExpenseListCompensation{}

	positiveDifferences := []models.ExpenseListShare{}
	negativeDifferences := []models.ExpenseListShare{}
	for _, share := range shares {
		if share.Difference >= 0 {
			positiveDifferences = append(positiveDifferences, *share)
		} else {
			negativeDifferences = append(negativeDifferences, *share)
		}
	}

	for _, negativeParticipant := range negativeDifferences {
		for _, positiveParticipant := range positiveDifferences {
			amount := min(-negativeParticipant.Difference, positiveParticipant.Difference)

			positiveParticipant.Difference -= amount
			negativeParticipant.Difference -= amount

			if amount <= 0 {
				continue
			}

			compensations = append(compensations, models.ExpenseListCompensation{
				From:   negativeParticipant.ParticipantId,
				To:     positiveParticipant.ParticipantId,
				Amount: amount,
			})
		}
	}

	return compensations
}

func calculate_shares(expenselist models.ExpenseList) []*models.ExpenseListShare {
	shares := []*models.ExpenseListShare{}

	for _, particpant := range expenselist.Participants {
		shares = append(shares, &models.ExpenseListShare{
			ParticipantId:    particpant.ID,
			NumberOfExpenses: 0,
			ExpenseAmount:    0,
			Share:            0,
			Difference:       0,
		})
	}

	for _, expense := range expenselist.Expenses {
		shareIndex := slices.IndexFunc(shares, func(s *models.ExpenseListShare) bool { return s.ParticipantId == expense.Buyer })
		if shareIndex == -1 {
			continue
		}
		share := shares[shareIndex]
		share.ExpenseAmount += expense.Amount
		share.NumberOfExpenses++
		share.Difference = share.ExpenseAmount - share.Share

		for _, participant := range expense.Participants {
			shareIndex := slices.IndexFunc(shares, func(s *models.ExpenseListShare) bool { return s.ParticipantId == participant })
			if shareIndex == -1 {
				continue
			}
			share := shares[shareIndex]
			share.Share += expense.Amount / float64(len(expense.Participants))
			share.Difference = share.ExpenseAmount - share.Share
		}
	}

	slices.SortFunc(shares, func(a, b *models.ExpenseListShare) int {
		return cmp.Compare(a.Difference, b.Difference)
	})
	return shares
}

func Calculate_shares_and_compensations(expenselist models.ExpenseList) models.ExpenseListWrapper {
	shares := calculate_shares(expenselist)

	listWrapper := models.ExpenseListWrapper{
		ExpenseList:   expenselist,
		Shares:        shares,
		Compensations: calculate_compensations(shares),
	}
	return listWrapper
}
