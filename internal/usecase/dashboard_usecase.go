package usecase

import (
	"context"

	"nexus-video-games/internal/model"
	"nexus-video-games/internal/repository/postgresql"

	"github.com/google/uuid"
)

type DashboardUsecase struct {
	balanceRepo *postgresql.BalanceRepository
	topupRepo   *postgresql.TopupRepository
	gameRepo    *postgresql.GameRepository
}

func NewDashboardUsecase(
	balanceRepo *postgresql.BalanceRepository,
	topupRepo *postgresql.TopupRepository,
	gameRepo *postgresql.GameRepository,
) *DashboardUsecase {
	return &DashboardUsecase{
		balanceRepo: balanceRepo,
		topupRepo:   topupRepo,
		gameRepo:    gameRepo,
	}
}

func (u *DashboardUsecase) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*model.DashboardStatsResponse, error) {
	balance, err := u.balanceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalGames, err := u.gameRepo.CountOwnedGames(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalCashSpent, err := u.topupRepo.GetTotalCashSpentByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	monthlyCashSpent, err := u.topupRepo.GetMonthlyCashSpentByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	library, err := u.gameRepo.FindUserLibrary(ctx, userID)
	if err != nil {
		return nil, err
	}

	var monthlyGames []model.MonthlyGamePlayResponse
	for _, item := range library {
		if item.PlaytimeMinutes > 0 {
			title := item.ExternalGameID
			if item.Game.Title != "" {
				title = item.Game.Title
			}
			monthlyGames = append(monthlyGames, model.MonthlyGamePlayResponse{
				ExternalGameID:  item.ExternalGameID,
				Title:           title,
				PlaytimeMinutes: item.PlaytimeMinutes,
			})
		}
	}

	return &model.DashboardStatsResponse{
		CurrentBalance:     balance.BalanceAmount,
		Currency:           balance.Currency,
		TotalGamesOwned:    totalGames,
		TotalCashSpent:     totalCashSpent,
		MonthlyCashSpent:   monthlyCashSpent,
		MonthlyGamesPlayed: monthlyGames,
	}, nil
}
