package usecase

import (
	"context"
	"errors"

	"nexus-video-games/internal/entity"
	"nexus-video-games/internal/gateway/game"
	"nexus-video-games/internal/model"
	"nexus-video-games/internal/repository/postgresql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameUsecase struct {
	db          *gorm.DB
	gameRepo    *postgresql.GameRepository
	balanceRepo *postgresql.BalanceRepository
	itadGW      *game.ITADGateway
}

func NewGameUsecase(
	db *gorm.DB,
	gameRepo *postgresql.GameRepository,
	balanceRepo *postgresql.BalanceRepository,
	itadGW *game.ITADGateway,
) *GameUsecase {
	return &GameUsecase{
		db:          db,
		gameRepo:    gameRepo,
		balanceRepo: balanceRepo,
		itadGW:      itadGW,
	}
}

func (u *GameUsecase) SyncCatalog(ctx context.Context) error {
	catalogItems, err := u.itadGW.FetchDealsCatalog(ctx)
	if err != nil {
		return err
	}
	return u.gameRepo.UpsertCatalog(ctx, catalogItems)
}

func (u *GameUsecase) ExploreGames(ctx context.Context) ([]model.GameCatalogItemResponse, error) {
	games, err := u.gameRepo.FindAvailableCatalog(ctx)
	if err != nil {
		return nil, err
	}

	var res []model.GameCatalogItemResponse
	for _, g := range games {
		res = append(res, model.GameCatalogItemResponse{
			ExternalGameID: g.ExternalGameID,
			Title:          g.Title,
			PriceInBalance: g.PriceInBalance,
			IsAvailable:    g.IsAvailable,
		})
	}
	return res, nil
}

func (u *GameUsecase) BuyGame(ctx context.Context, userID uuid.UUID, req *model.BuyGameRequest) error {
	game, err := u.gameRepo.FindCatalogByID(ctx, req.ExternalGameID)
	if err != nil {
		return errors.New("game not found in catalog")
	}

	if !game.IsAvailable {
		return errors.New("game is currently unavailable for purchase")
	}

	owned, _ := u.gameRepo.FindOwnedGame(ctx, userID, req.ExternalGameID)
	if owned != nil {
		return errors.New("game is already in your library")
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.balanceRepo.DeductBalance(ctx, tx, userID, game.PriceInBalance); err != nil {
			return errors.New("insufficient wallet balance")
		}

		newOwned := &entity.GamesOwned{
			UserID:          userID,
			ExternalGameID:  game.ExternalGameID,
			PlaytimeMinutes: 0,
		}
		return u.gameRepo.CreateOwned(ctx, tx, newOwned)
	})
}

func (u *GameUsecase) GetUserLibrary(ctx context.Context, userID uuid.UUID) ([]model.UserGameLibraryResponse, error) {
	owned, err := u.gameRepo.FindUserLibrary(ctx, userID)
	if err != nil {
		return nil, err
	}

	var library []model.UserGameLibraryResponse
	for _, item := range owned {
		title := item.ExternalGameID
		if item.Game.Title != "" {
			title = item.Game.Title
		}

		library = append(library, model.UserGameLibraryResponse{
			ID:              item.ID,
			ExternalGameID:  item.ExternalGameID,
			Title:           title,
			AcquiredAt:      item.AcquiredAt,
			PlaytimeMinutes: item.PlaytimeMinutes,
		})
	}
	return library, nil
}
