package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PackagesService struct {
	DB *gorm.DB
}

func NewPackagesService(db *gorm.DB) *PackagesService {
	return &PackagesService{DB: db}
}

// CREATE
func (s *PackagesService) CreatePackages(pkg *model.Packages) error {
	if pkg.ID == "" {
		pkg.ID = uuid.New().String()
	}

	if pkg.Name == "" {
		return fmt.Errorf("name is required")
	}

	if pkg.Price <= 0 {
		return fmt.Errorf("price must be > 0")
	}

	if pkg.DurationDays <= 0 {
		return fmt.Errorf("duration_days must be > 0")
	}

	if pkg.Status == "" {
		pkg.Status = "active"
	}

	// Validasi user_id ada di tabel users (opsional, jika ingin foreign key constraint)
	var user model.User
	if err := s.DB.First(&user, "id = ?", pkg.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	// Validasi packages_id jika ada (opsional, self-referencing)
	if pkg.PackagesID != "" {
		var parentPkg model.Packages
		if err := s.DB.First(&parentPkg, "id = ?", pkg.PackagesID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("packages not found")
			}
			return err
		}
	}

	return s.DB.Create(pkg).Error
}

// READ ALL
func (s *PackagesService) GetAllPackages() ([]model.Packages, error) {
	var packages []model.Packages
	result := s.DB.Find(&packages)
	return packages, result.Error
}

// READ BY ID
func (s *PackagesService) GetPackagesByID(id string) (*model.Packages, error) {
	var pkg model.Packages
	result := s.DB.First(&pkg, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pkg, result.Error
}

// READ BY USER ID
func (s *PackagesService) GetPackagesByUserID(userID string) ([]model.Packages, error) {
	var packages []model.Packages
	result := s.DB.Where("user_id = ?", userID).Find(&packages)
	return packages, result.Error
}

// UPDATE
func (s *PackagesService) UpdatePackages(pkg *model.Packages) error {
	if pkg.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldPkg model.Packages
	if err := s.DB.First(&oldPkg, "id = ?", pkg.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("packages not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if pkg.Name != "" && pkg.Name != oldPkg.Name {
		updateData["name"] = pkg.Name
	}

	if pkg.Description != "" && pkg.Description != oldPkg.Description {
		updateData["description"] = pkg.Description
	}

	if pkg.Price > 0 && pkg.Price != oldPkg.Price {
		updateData["price"] = pkg.Price
	}

	if pkg.DurationDays > 0 && pkg.DurationDays != oldPkg.DurationDays {
		updateData["duration_days"] = pkg.DurationDays
	}

	if pkg.Status != "" && pkg.Status != oldPkg.Status {
		updateData["status"] = pkg.Status
	}

	if pkg.UserID != "" && pkg.UserID != oldPkg.UserID {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", pkg.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = pkg.UserID
	}

	if pkg.PackagesID != "" && pkg.PackagesID != oldPkg.PackagesID {
		// Validasi packages_id baru ada
		var parentPkg model.Packages
		if err := s.DB.First(&parentPkg, "id = ?", pkg.PackagesID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("packages not found")
			}
			return err
		}
		updateData["packages_id"] = pkg.PackagesID
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.Packages{}).Where("id = ?", pkg.ID).Updates(updateData).Error
}

// DELETE
func (s *PackagesService) DeletePackages(id string) error {
	return s.DB.Delete(&model.Packages{}, "id = ?", id).Error
}
