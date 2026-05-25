package services

import (
	"fmt"
	"testing"

	"noams/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupDeviceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Device{},
		&models.DeviceGroup{},
		&models.Credential{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestDeviceCreate(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	err := svc.Create(&models.Device{
		Name:         "test-switch",
		DeviceType:   "hp_comware",
		Vendor:       "h3c",
		Role:         "access",
		ManagementIP: "10.0.0.1",
		SSHPort:      22,
		Building:     "test-building",
	})
	if err != nil {
		t.Fatalf("DeviceService.Create failed: %v", err)
	}
}

func TestDeviceList(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	svc.Create(&models.Device{Name: "dev1", DeviceType: "hp_comware", Vendor: "h3c", ManagementIP: "10.0.0.1"})
	svc.Create(&models.Device{Name: "dev2", DeviceType: "huawei", Vendor: "huawei", ManagementIP: "10.0.0.2", Building: "A"})
	svc.Create(&models.Device{Name: "dev3", DeviceType: "cisco_ios", Vendor: "cisco", ManagementIP: "10.0.0.3", Role: "core"})

	devices, total, err := svc.List(1, 10, map[string]interface{}{})
	if err != nil {
		t.Fatalf("DeviceService.List failed: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected 3 devices, got %d", total)
	}
	if len(devices) != 3 {
		t.Fatalf("expected 3 devices in list, got %d", len(devices))
	}
}

func TestDeviceList_FilterByVendor(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	svc.Create(&models.Device{Name: "d1", DeviceType: "hp_comware", Vendor: "h3c", ManagementIP: "10.0.0.1"})
	svc.Create(&models.Device{Name: "d2", DeviceType: "huawei", Vendor: "huawei", ManagementIP: "10.0.0.2"})
	svc.Create(&models.Device{Name: "d3", DeviceType: "hp_comware", Vendor: "h3c", ManagementIP: "10.0.0.3"})

	_, total, err := svc.List(1, 10, map[string]interface{}{"vendor": "h3c"})
	if err != nil {
		t.Fatalf("List with filter failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 h3c devices, got %d", total)
	}
}

func TestDeviceGetByID(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	device := &models.Device{Name: "findme", DeviceType: "hp_comware", ManagementIP: "10.0.0.99"}
	svc.Create(device)

	found, err := svc.GetByID(device.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found.Name != "findme" {
		t.Fatalf("expected name 'findme', got '%s'", found.Name)
	}
	if found.ManagementIP != "10.0.0.99" {
		t.Fatalf("expected IP '10.0.0.99', got '%s'", found.ManagementIP)
	}
}

func TestDeviceUpdate(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	device := &models.Device{Name: "oldname", DeviceType: "hp_comware", ManagementIP: "10.0.0.1", Building: "A"}
	svc.Create(device)

	err := svc.Update(device.ID, map[string]interface{}{"name": "newname", "building": "B"})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := svc.GetByID(device.ID)
	if found.Name != "newname" {
		t.Fatalf("expected updated name 'newname', got '%s'", found.Name)
	}
	if found.Building != "B" {
		t.Fatalf("expected updated building 'B', got '%s'", found.Building)
	}
}

func TestDeviceDelete(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	device := &models.Device{Name: "todelete", DeviceType: "hp_comware", ManagementIP: "10.0.0.1"}
	svc.Create(device)

	err := svc.Delete(device.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = svc.GetByID(device.ID)
	if err == nil {
		t.Fatal("GetByID should fail after delete")
	}
}

func TestDeviceCountByStatus(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	svc.Create(&models.Device{Name: "online1", DeviceType: "hp_comware", ManagementIP: "10.0.0.1", Status: 1})
	svc.Create(&models.Device{Name: "online2", DeviceType: "hp_comware", ManagementIP: "10.0.0.2", Status: 1})
	svc.Create(&models.Device{Name: "offline1", DeviceType: "hp_comware", ManagementIP: "10.0.0.3", Status: 0})

	online, offline, err := svc.CountByStatus()
	if err != nil {
		t.Fatalf("CountByStatus failed: %v", err)
	}
	if online != 2 {
		t.Fatalf("expected 2 online, got %d", online)
	}
	if offline != 1 {
		t.Fatalf("expected 1 offline, got %d", offline)
	}
}

func TestDevicePagination(t *testing.T) {
	db := setupDeviceTestDB(t)
	svc := NewDeviceService(db)

	for i := 1; i <= 25; i++ {
		svc.Create(&models.Device{
			Name:         fmt.Sprintf("dev-%d", i),
			DeviceType:   "hp_comware",
			ManagementIP: fmt.Sprintf("10.0.0.%d", i),
		})
	}

	// Page 1 (20 items)
	devices, total, err := svc.List(1, 20, map[string]interface{}{})
	if err != nil {
		t.Fatalf("List page 1 failed: %v", err)
	}
	if total != 25 {
		t.Fatalf("expected total 25, got %d", total)
	}
	if len(devices) != 20 {
		t.Fatalf("expected 20 devices on page 1, got %d", len(devices))
	}

	// Page 2 (5 items)
	devices, total, err = svc.List(2, 20, map[string]interface{}{})
	if err != nil {
		t.Fatalf("List page 2 failed: %v", err)
	}
	if len(devices) != 5 {
		t.Fatalf("expected 5 devices on page 2, got %d", len(devices))
	}
}
