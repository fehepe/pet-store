//go:generate go run github.com/99designs/gqlgen generate

package graph

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fehepe/pet-store/backend/internal/auth"
	"github.com/fehepe/pet-store/backend/internal/graph/model"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/fehepe/pet-store/backend/internal/service"
	"github.com/google/uuid"
)

type Resolver struct {
	storeService *service.StoreService
	petService   *service.PetService
	orderService *service.OrderService
}

func NewResolver(storeService *service.StoreService, petService *service.PetService, orderService *service.OrderService) *Resolver {
	return &Resolver{
		storeService: storeService,
		petService:   petService,
		orderService: orderService,
	}
}

// Resolver implements the ResolverRoot interface
func (r *Resolver) Query() QueryResolver {
	return r
}

func (r *Resolver) Mutation() MutationResolver {
	return r
}

// Query resolvers

func (r *Resolver) MyStore(ctx context.Context) (*model.Store, error) {
	store, err := r.getStoreForMerchant(ctx)
	if err != nil {
		return nil, err
	}

	return &model.Store{
		ID:        store.ID,
		Name:      store.Name,
		CreatedAt: store.CreatedAt,
	}, nil
}

func (r *Resolver) ListPets(ctx context.Context, filter *model.PetFilterInput, pagination *model.PaginationInput) (*model.PetConnection, error) {
	store, err := r.getStoreForMerchant(ctx)
	if err != nil {
		return nil, err
	}

	// Build filter
	petFilter := models.PetFilter{
		StoreID: &store.ID,
		Limit:   50, // Default limit
		Offset:  0,
	}

	if filter != nil {
		if filter.Status != nil {
			status := models.PetStatus(*filter.Status)
			petFilter.Status = &status
		}
		if filter.StartDate != nil {
			petFilter.StartDate = filter.StartDate
		}
		if filter.EndDate != nil {
			petFilter.EndDate = filter.EndDate
		}
	}

	r.applyPagination(&petFilter, pagination)

	pets, totalCount, err := r.petService.ListPets(ctx, petFilter)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL types
	var edges []*model.Pet
	for _, pet := range pets {
		edges = append(edges, r.petToGraphQLModel(pet, true))
	}

	hasNextPage := len(pets) == petFilter.Limit
	endCursor := ""
	if hasNextPage {
		endCursor = strconv.Itoa(petFilter.Offset + len(pets))
	}

	return &model.PetConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: petFilter.Offset > 0,
			EndCursor:       &endCursor,
		},
		TotalCount: int32(totalCount),
	}, nil
}

func (r *Resolver) GetPet(ctx context.Context, id uuid.UUID) (*model.Pet, error) {
	store, err := r.getStoreForMerchant(ctx)
	if err != nil {
		return nil, err
	}

	pet, err := r.petService.GetPetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if pet.StoreID != store.ID {
		return nil, fmt.Errorf("pet not found")
	}

	// Decrypt email for merchants
	decryptedEmail := "[Hidden]"
	if email, err := r.petService.DecryptBreederEmail(pet.BreederEmailEncrypted); err == nil {
		decryptedEmail = email
	}

	return &model.Pet{
		ID:           pet.ID,
		Name:         pet.Name,
		Species:      model.PetSpecies(pet.Species),
		Age:          int32(pet.Age),
		PictureURL:   pet.PictureURL,
		Description:  pet.Description,
		BreederName:  pet.BreederName,
		BreederEmail: decryptedEmail,
		Status:       model.PetStatus(pet.Status),
		CreatedAt:    pet.CreatedAt,
	}, nil
}

func (r *Resolver) AvailablePets(ctx context.Context, storeID uuid.UUID, pagination *model.PaginationInput) (*model.PetConnection, error) {
	// This is now a public endpoint for demo purposes

	// Build filter for available pets only
	status := models.PetStatusAvailable
	petFilter := models.PetFilter{
		StoreID: &storeID,
		Status:  &status,
		Limit:   50, // Default limit
		Offset:  0,
	}

	r.applyPagination(&petFilter, pagination)

	pets, totalCount, err := r.petService.ListPets(ctx, petFilter)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL types (hide breeder email for customers)
	var edges []*model.Pet
	for _, pet := range pets {
		edges = append(edges, r.petToGraphQLModel(pet, false)) // Hide email for customers
	}

	hasNextPage := len(pets) == petFilter.Limit
	endCursor := ""
	if hasNextPage {
		endCursor = strconv.Itoa(petFilter.Offset + len(pets))
	}

	return &model.PetConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: petFilter.Offset > 0,
			EndCursor:       &endCursor,
		},
		TotalCount: int32(totalCount),
	}, nil
}


func (r *Resolver) ListStores(ctx context.Context) ([]*model.Store, error) {
	stores, err := r.storeService.ListAllStores(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Store
	for _, store := range stores {
		result = append(result, &model.Store{
			ID:        store.ID,
			Name:      store.Name,
			CreatedAt: store.CreatedAt,
		})
	}

	return result, nil
}

func (r *Resolver) SoldPets(ctx context.Context, startDate time.Time, endDate time.Time, pagination *model.PaginationInput) (*model.PetConnection, error) {
	store, err := r.getStoreForMerchant(ctx)
	if err != nil {
		return nil, err
	}

	// Build filter for sold pets within date range
	status := models.PetStatusSold
	petFilter := models.PetFilter{
		StoreID:   &store.ID,
		Status:    &status,
		StartDate: &startDate,
		EndDate:   &endDate,
		Limit:     50, // Default limit
		Offset:    0,
	}

	r.applyPagination(&petFilter, pagination)

	pets, totalCount, err := r.petService.ListPets(ctx, petFilter)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL types
	var edges []*model.Pet
	for _, pet := range pets {
		edges = append(edges, r.petToGraphQLModel(pet, true))
	}

	hasNextPage := len(pets) == petFilter.Limit
	endCursor := ""
	if hasNextPage {
		endCursor = strconv.Itoa(petFilter.Offset + len(pets))
	}

	return &model.PetConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: petFilter.Offset > 0,
			EndCursor:       &endCursor,
		},
		TotalCount: int32(totalCount),
	}, nil
}

func (r *Resolver) UnsoldPets(ctx context.Context, pagination *model.PaginationInput) (*model.PetConnection, error) {
	store, err := r.getStoreForMerchant(ctx)
	if err != nil {
		return nil, err
	}

	// Build filter for available (unsold) pets
	status := models.PetStatusAvailable
	petFilter := models.PetFilter{
		StoreID: &store.ID,
		Status:  &status,
		Limit:   50, // Default limit
		Offset:  0,
	}

	r.applyPagination(&petFilter, pagination)

	pets, totalCount, err := r.petService.ListPets(ctx, petFilter)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL types
	var edges []*model.Pet
	for _, pet := range pets {
		edges = append(edges, r.petToGraphQLModel(pet, true))
	}

	hasNextPage := len(pets) == petFilter.Limit
	endCursor := ""
	if hasNextPage {
		endCursor = strconv.Itoa(petFilter.Offset + len(pets))
	}

	return &model.PetConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: petFilter.Offset > 0,
			EndCursor:       &endCursor,
		},
		TotalCount: int32(totalCount),
	}, nil
}

// Mutation resolvers

func (r *Resolver) CreatePet(ctx context.Context, input model.CreatePetInput) (*model.Pet, error) {
	store, err := r.getStoreForMerchant(ctx)
	if err != nil {
		return nil, err
	}

	createInput := models.CreatePetInput{
		StoreID:      store.ID,
		Name:         input.Name,
		Species:      models.PetSpecies(input.Species),
		Age:          int(input.Age),
		PictureURL:   input.PictureURL,
		Description:  input.Description,
		BreederName:  input.BreederName,
		BreederEmail: input.BreederEmail,
	}

	pet, err := r.petService.CreatePet(ctx, createInput)
	if err != nil {
		return nil, err
	}

	return &model.Pet{
		ID:           pet.ID,
		Name:         pet.Name,
		Species:      model.PetSpecies(pet.Species),
		Age:          int32(pet.Age),
		PictureURL:   pet.PictureURL,
		Description:  pet.Description,
		BreederName:  pet.BreederName,
		BreederEmail: input.BreederEmail, // Return original email
		Status:       model.PetStatus(pet.Status),
		CreatedAt:    pet.CreatedAt,
	}, nil
}

func (r *Resolver) DeletePet(ctx context.Context, id uuid.UUID) (bool, error) {
	store, err := r.getStoreForMerchant(ctx)
	if err != nil {
		return false, err
	}

	// Verify ownership
	pet, err := r.petService.GetPetByID(ctx, id)
	if err != nil {
		return false, err
	}

	if pet.StoreID != store.ID {
		return false, fmt.Errorf("pet not found")
	}

	err = r.petService.DeletePetByID(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Resolver) PurchasePet(ctx context.Context, petID uuid.UUID) (*model.Order, error) {
	if err := auth.RequireCustomer(ctx); err != nil {
		return nil, err
	}

	username, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	// Get the pet to find its store
	pet, err := r.petService.GetPetByID(ctx, petID)
	if err != nil {
		return nil, err
	}

	order, err := r.orderService.CreateOrder(ctx, models.CreateOrderInput{
		CustomerID: username,
		StoreID:    pet.StoreID,
		PetIDs:     []uuid.UUID{petID},
	})
	if err != nil {
		// Check if it's a pet availability error
		if strings.Contains(err.Error(), "no longer available") {
			return nil, fmt.Errorf("sorry, the pet '%s' is no longer available for purchase. It may have been purchased by another customer", pet.Name)
		}
		return nil, fmt.Errorf("unable to complete the purchase: %v", err)
	}

	// Get pets for the order
	pets, err := r.orderService.GetOrderPets(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	var modelPets []*model.Pet
	for _, pet := range pets {
		modelPets = append(modelPets, r.petToGraphQLModel(pet, false))
	}

	return &model.Order{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Pets:       modelPets,
		TotalPets:  int32(order.TotalPets),
		CreatedAt:  order.CreatedAt,
	}, nil
}

func (r *Resolver) PurchasePets(ctx context.Context, petIDs []uuid.UUID) (*model.Order, error) {
	if err := auth.RequireCustomer(ctx); err != nil {
		return nil, err
	}

	username, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	if len(petIDs) == 0 {
		return nil, fmt.Errorf("no pets specified")
	}

	// Get the first pet to find the store (all pets should be from same store)
	firstPet, err := r.petService.GetPetByID(ctx, petIDs[0])
	if err != nil {
		return nil, err
	}

	order, err := r.orderService.CreateOrder(ctx, models.CreateOrderInput{
		CustomerID: username,
		StoreID:    firstPet.StoreID,
		PetIDs:     petIDs,
	})
	if err != nil {
		// Check if it's a pet availability error
		if strings.Contains(err.Error(), "no longer available") {
			return nil, fmt.Errorf("some pets in your cart are no longer available for purchase. They may have been purchased by other customers. %v", err)
		}
		return nil, fmt.Errorf("unable to complete the purchase: %v", err)
	}

	// Get pets for the order
	pets, err := r.orderService.GetOrderPets(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	var modelPets []*model.Pet
	for _, pet := range pets {
		modelPets = append(modelPets, r.petToGraphQLModel(pet, false))
	}

	return &model.Order{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Pets:       modelPets,
		TotalPets:  int32(order.TotalPets),
		CreatedAt:  order.CreatedAt,
	}, nil
}

func (r *Resolver) CreateStore(ctx context.Context, input model.CreateStoreInput) (*model.Store, error) {
	if err := auth.RequireMerchant(ctx); err != nil {
		return nil, err
	}

	username, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	createInput := models.CreateStoreInput{
		Name:    input.Name,
		OwnerID: username,
	}

	store, err := r.storeService.CreateStore(ctx, createInput)
	if err != nil {
		return nil, err
	}

	return &model.Store{
		ID:        store.ID,
		Name:      store.Name,
		CreatedAt: store.CreatedAt,
	}, nil
}

// Helper method to convert models.Pet to model.Pet with email handling
func (r *Resolver) petToGraphQLModel(pet *models.Pet, showEmail bool) *model.Pet {
	var breederEmail string
	if showEmail {
		// Decrypt email for merchants
		if decrypted, err := r.petService.DecryptBreederEmail(pet.BreederEmailEncrypted); err == nil {
			breederEmail = decrypted
		} else {
			breederEmail = "[Hidden]"
		}
	} else {
		breederEmail = "[Hidden]"
	}

	return &model.Pet{
		ID:           pet.ID,
		Name:         pet.Name,
		Species:      model.PetSpecies(pet.Species),
		Age:          int32(pet.Age),
		PictureURL:   pet.PictureURL,
		Description:  pet.Description,
		BreederName:  pet.BreederName,
		BreederEmail: breederEmail,
		Status:       model.PetStatus(pet.Status),
		CreatedAt:    pet.CreatedAt,
	}
}

// Helper method to get store for authenticated merchant
func (r *Resolver) getStoreForMerchant(ctx context.Context) (*models.Store, error) {
	if err := auth.RequireMerchant(ctx); err != nil {
		return nil, err
	}

	username, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	return r.storeService.GetStoreByOwnerID(ctx, username)
}

// Helper method to apply pagination to pet filter
func (r *Resolver) applyPagination(petFilter *models.PetFilter, pagination *model.PaginationInput) {
	if pagination != nil {
		if pagination.First != nil {
			petFilter.Limit = int(*pagination.First)
		}
		if pagination.After != nil {
			offset, _ := strconv.Atoi(*pagination.After)
			petFilter.Offset = offset
		}
	}
}
