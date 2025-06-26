package app

import (
	"fmt"

	"github.com/fehepe/pet-store/backend/internal/cache"
	"github.com/fehepe/pet-store/backend/internal/config"
	"github.com/fehepe/pet-store/backend/internal/database"
	"github.com/fehepe/pet-store/backend/internal/graph"
	"github.com/fehepe/pet-store/backend/internal/repository"
	"github.com/fehepe/pet-store/backend/internal/service"
	"github.com/fehepe/pet-store/backend/pkg/encryption"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	Config       *config.Config
	DB           database.Repository
	Cache        cache.CacheInterface
	Encryptor    encryption.EncryptorInterface
	Repositories *Repositories
	Services     *Services
	Resolver     graph.ResolverRoot
}

// Repositories holds all repository instances
type Repositories struct {
	Pet   repository.PetRepositoryInterface
	Store repository.StoreRepositoryInterface
	Order repository.OrderRepositoryInterface
}

// Services holds all service instances
type Services struct {
	Pet   *service.PetService
	Store *service.StoreService
	Order *service.OrderService
}

// InitializeDependencies initializes all application dependencies
func InitializeDependencies(cfg *config.Config) (*Dependencies, error) {
	db, err := database.NewConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	redisCache, err := cache.NewRedisCache(cfg)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	encryptor, err := encryption.NewEncryptor(cfg.EncryptionKey)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize encryptor: %w", err)
	}

	repos := &Repositories{
		Pet:   repository.NewPetRepository(db),
		Store: repository.NewStoreRepository(db),
		Order: repository.NewOrderRepository(db),
	}

	services := &Services{
		Store: service.NewStoreService(repos.Store, redisCache),
		Pet:   service.NewPetService(repos.Pet, redisCache, encryptor),
	}
	services.Order = service.NewOrderService(repos.Order, repos.Pet, redisCache, services.Pet)

	resolver := graph.NewResolver(services.Store, services.Pet, services.Order)

	return &Dependencies{
		Config:       cfg,
		DB:           db,
		Cache:        redisCache,
		Encryptor:    encryptor,
		Repositories: repos,
		Services:     services,
		Resolver:     resolver,
	}, nil
}

// Close closes all closeable dependencies
func (d *Dependencies) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.Cache != nil {
		d.Cache.Close()
	}
}
