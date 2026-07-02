package database

import (
	"database/sql"
	"fmt"
	"time"

	"wms/config"
	"wms/pkg/logger"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	lifetime, err := time.ParseDuration(cfg.ConnMaxLifetime)
	if err != nil {
		lifetime = 5 * time.Minute
	}
	db.SetConnMaxLifetime(lifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(100) NOT NULL,
			role VARCHAR(20) NOT NULL CHECK (role IN ('admin','manager','operator','viewer')),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			sku VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			category VARCHAR(100),
			unit_of_measure VARCHAR(20) NOT NULL,
			weight DECIMAL(10,2),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS locations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			code VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(100) NOT NULL,
			zone VARCHAR(50),
			aisle VARCHAR(20),
			rack VARCHAR(20),
			level VARCHAR(20),
			bin VARCHAR(20),
			type VARCHAR(20) CHECK (type IN ('storage','receiving','shipping','staging')),
			capacity DECIMAL(10,2),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS stock (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			item_id UUID REFERENCES items(id),
			location_id UUID REFERENCES locations(id),
			quantity DECIMAL(12,2) NOT NULL DEFAULT 0,
			reserved_quantity DECIMAL(12,2) NOT NULL DEFAULT 0,
			batch_number VARCHAR(50),
			expiry_date DATE,
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(item_id, location_id, batch_number)
		)`,
		`CREATE TABLE IF NOT EXISTS stock_movements (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			item_id UUID REFERENCES items(id),
			from_location_id UUID REFERENCES locations(id),
			to_location_id UUID REFERENCES locations(id),
			quantity DECIMAL(12,2) NOT NULL,
			movement_type VARCHAR(30) NOT NULL,
			reference_type VARCHAR(30),
			reference_id UUID,
			notes TEXT,
			created_by UUID REFERENCES users(id),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS purchase_orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			po_number VARCHAR(50) UNIQUE NOT NULL,
			supplier_name VARCHAR(255),
			status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending','approved','received','cancelled')),
			expected_date DATE,
			storage_location_id UUID REFERENCES locations(id),
			notes TEXT,
			created_by UUID REFERENCES users(id),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS purchase_order_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			po_id UUID REFERENCES purchase_orders(id) ON DELETE CASCADE,
			item_id UUID REFERENCES items(id),
			expected_quantity DECIMAL(12,2) NOT NULL,
			received_quantity DECIMAL(12,2) DEFAULT 0,
			unit_price DECIMAL(12,2)
		)`,
		`CREATE TABLE IF NOT EXISTS goods_receipts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			grn_number VARCHAR(50) UNIQUE NOT NULL,
			po_id UUID REFERENCES purchase_orders(id),
			received_by UUID REFERENCES users(id),
			received_at TIMESTAMPTZ DEFAULT NOW(),
			notes TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS goods_receipt_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			grn_id UUID REFERENCES goods_receipts(id) ON DELETE CASCADE,
			item_id UUID REFERENCES items(id),
			quantity DECIMAL(12,2) NOT NULL,
			batch_number VARCHAR(50),
			location_id UUID REFERENCES locations(id)
		)`,
		`CREATE TABLE IF NOT EXISTS sales_orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			so_number VARCHAR(50) UNIQUE NOT NULL,
			customer_name VARCHAR(255),
			status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending','picking','picked','packed','shipped','cancelled')),
			notes TEXT,
			created_by UUID REFERENCES users(id),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sales_order_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			so_id UUID REFERENCES sales_orders(id) ON DELETE CASCADE,
			item_id UUID REFERENCES items(id),
			ordered_quantity DECIMAL(12,2) NOT NULL,
			picked_quantity DECIMAL(12,2) DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS pick_lists (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			so_id UUID REFERENCES sales_orders(id),
			status VARCHAR(20) DEFAULT 'pending',
			picked_by UUID REFERENCES users(id),
			picked_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS shipments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			shipment_number VARCHAR(50) UNIQUE NOT NULL,
			so_id UUID REFERENCES sales_orders(id),
			carrier VARCHAR(100),
			tracking_number VARCHAR(100),
			status VARCHAR(20) DEFAULT 'pending',
			shipped_at TIMESTAMPTZ,
			created_by UUID REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS stock_transfers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			transfer_number VARCHAR(50) UNIQUE NOT NULL,
			from_location_id UUID REFERENCES locations(id),
			to_location_id UUID REFERENCES locations(id),
			item_id UUID REFERENCES items(id),
			quantity DECIMAL(12,2) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending','in_transit','completed','cancelled')),
			notes TEXT,
			created_by UUID REFERENCES users(id),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	// ALTER TABLE for existing databases
	alterations := []string{
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'purchase_orders' AND column_name = 'storage_location_id') THEN
				ALTER TABLE purchase_orders ADD COLUMN storage_location_id UUID REFERENCES locations(id);
			END IF;
		END $$`,
	}

	for _, a := range alterations {
		db.Exec(a)
	}

	return nil
}

func Seed(db *sql.DB, log logger.Logger) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Error("failed to check users count", "error", err)
		return
	}
	if count > 0 {
		log.Info("database already seeded, skipping")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error("failed to begin seed transaction", "error", err)
		return
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash seed password", "error", err)
		return
	}

	hashedManager, _ := bcrypt.GenerateFromPassword([]byte("manager123"), bcrypt.DefaultCost)
	hashedOperator, _ := bcrypt.GenerateFromPassword([]byte("operator123"), bcrypt.DefaultCost)

	adminID := "a0000000-0000-0000-0000-000000000001"
	managerID := "a0000000-0000-0000-0000-000000000002"
	operatorID := "a0000000-0000-0000-0000-000000000003"

	_, err = tx.Exec(`INSERT INTO users (id, username, email, password_hash, full_name, role) VALUES
		($1, 'admin', 'admin@wms.local', $2, 'Administrator', 'admin'),
		($3, 'manager', 'manager@wms.local', $4, 'Warehouse Manager', 'manager'),
		($5, 'operator', 'operator@wms.local', $6, 'Warehouse Operator', 'operator')`,
		adminID, string(hashedPassword), managerID, string(hashedManager), operatorID, string(hashedOperator))
	if err != nil {
		log.Error("failed to seed users", "error", err)
		return
	}

	itemIDs := make([]string, 10)
	for i := 0; i < 10; i++ {
		itemIDs[i] = fmt.Sprintf("b0000000-0000-0000-0000-%012d", i+1)
	}

	_, err = tx.Exec(`INSERT INTO items (id, sku, name, description, category, unit_of_measure, weight) VALUES
		($1,  'ELC-001', 'Laptop ASUS 14"',       'Laptop 14 inch AMD Ryzen 5',     'Electronics', 'unit', 1.50),
		($2,  'ELC-002', 'Mouse Logitech M331',   'Wireless mouse',                 'Electronics', 'unit', 0.10),
		($3,  'ELC-003', 'Keyboard Mechanical',    'RGB mechanical keyboard',         'Electronics', 'unit', 0.80),
		($4,  'ELC-004', 'Monitor LG 24"',         '24 inch Full HD IPS monitor',     'Electronics', 'unit', 3.50),
		($5,  'FUR-001', 'Meja Kerja Lipat',       'Meja lipat 120x60cm',            'Furniture',   'unit', 8.00),
		($6,  'FUR-002', 'Kursi Ergonomis',        'Kursi kantor adjustable',         'Furniture',   'unit', 5.50),
		($7,  'PKG-001', 'Kardus Box Kecil',       'Box kardus 30x20x15cm',           'Packaging',   'pcs',  0.05),
		($8,  'PKG-002', 'Kardus Box Besar',       'Box kardus 60x40x40cm',           'Packaging',   'pcs',  0.15),
		($9,  'STA-001', 'Tinta Printer Canon',    'Tinta hitam Canon GI-290',        'Stationery',  'unit', 0.20),
		($10, 'STA-002', 'Kertas A4 500 lembar',   'HVS A4 70gsm',                   'Stationery',  'pack', 1.00)`,
		itemIDs[0], itemIDs[1], itemIDs[2], itemIDs[3], itemIDs[4],
		itemIDs[5], itemIDs[6], itemIDs[7], itemIDs[8], itemIDs[9])
	if err != nil {
		log.Error("failed to seed items", "error", err)
		return
	}

	locIDs := make([]string, 10)
	for i := 0; i < 10; i++ {
		locIDs[i] = fmt.Sprintf("c0000000-0000-0000-0000-%012d", i+1)
	}

	_, err = tx.Exec(`INSERT INTO locations (id, code, name, zone, aisle, rack, level, bin, type, capacity) VALUES
		($1,  'RCV-01', 'Receiving Dock 1',     'Receiving', 'A', '1', '1', '1', 'receiving', 500),
		($2,  'RCV-02', 'Receiving Dock 2',     'Receiving', 'A', '1', '1', '2', 'receiving', 500),
		($3,  'STR-A1', 'Zone A - Rack 1',      'Zone A',    'A', '1', '1', '1', 'storage',  200),
		($4,  'STR-A2', 'Zone A - Rack 2',      'Zone A',    'A', '1', '2', '1', 'storage',  200),
		($5,  'STR-A3', 'Zone A - Rack 3',      'Zone A',    'A', '1', '3', '1', 'storage',  200),
		($6,  'STR-B1', 'Zone B - Rack 1',      'Zone B',    'B', '1', '1', '1', 'storage',  300),
		($7,  'STR-B2', 'Zone B - Rack 2',      'Zone B',    'B', '1', '2', '1', 'storage',  300),
		($8,  'STG-01', 'Staging Area',         'Staging',   'S', '1', '1', '1', 'staging',  100),
		($9,  'SHP-01', 'Shipping Dock 1',      'Shipping',  'C', '1', '1', '1', 'shipping', 500),
		($10, 'SHP-02', 'Shipping Dock 2',      'Shipping',  'C', '1', '1', '2', 'shipping', 500)`,
		locIDs[0], locIDs[1], locIDs[2], locIDs[3], locIDs[4],
		locIDs[5], locIDs[6], locIDs[7], locIDs[8], locIDs[9])
	if err != nil {
		log.Error("failed to seed locations", "error", err)
		return
	}

	stockEntries := []struct {
		itemIdx, locIdx int
		qty, reserved   float64
		batch           string
	}{
		{0, 2, 50, 5, "BATCH-A001"},
		{0, 5, 30, 0, "BATCH-A002"},
		{1, 2, 200, 10, "BATCH-B001"},
		{1, 6, 150, 0, "BATCH-B002"},
		{2, 3, 80, 0, "BATCH-C001"},
		{3, 6, 40, 3, "BATCH-D001"},
		{4, 4, 25, 0, "BATCH-E001"},
		{5, 5, 30, 2, "BATCH-F001"},
		{6, 3, 1000, 0, "BATCH-G001"},
		{6, 7, 500, 0, "BATCH-G002"},
		{7, 3, 300, 0, "BATCH-H001"},
		{8, 4, 100, 0, "BATCH-I001"},
		{9, 7, 200, 15, "BATCH-J001"},
	}

	for _, s := range stockEntries {
		_, err = tx.Exec(`INSERT INTO stock (item_id, location_id, quantity, reserved_quantity, batch_number) VALUES ($1, $2, $3, $4, $5)`,
			itemIDs[s.itemIdx], locIDs[s.locIdx], s.qty, s.reserved, s.batch)
		if err != nil {
			log.Error("failed to seed stock", "error", err, "batch", s.batch)
			return
		}
	}

	po1ID := "e0000000-0000-0000-0000-000000000001"
	po2ID := "e0000000-0000-0000-0000-000000000002"
	po3ID := "e0000000-0000-0000-0000-000000000003"

	_, err = tx.Exec(`INSERT INTO purchase_orders (id, po_number, supplier_name, status, expected_date, notes, created_by) VALUES
		($1, 'PO-2025-001', 'PT Teknologi Maju',   'received',  '2025-06-15', 'Pengadaan laptop & aksesoris', $4),
		($2, 'PO-2025-002', 'PT Furniture Jaya',    'pending',   '2025-07-10', 'Pengadaan meja & kursi',       $4),
		($3, 'PO-2025-003', 'PT Bahan Kemas',       'partial',   '2025-07-15', 'Pengadaan kardus & packing',   $5)`,
		po1ID, po2ID, po3ID, managerID, operatorID)
	if err != nil {
		log.Error("failed to seed purchase orders", "error", err)
		return
	}

	_, err = tx.Exec(`INSERT INTO purchase_order_items (po_id, item_id, expected_quantity, received_quantity, unit_price) VALUES
		($1, $4,  100, 80, 8500000),
		($1, $5,  300, 300, 85000),
		($1, $6,  200, 200, 150000),
		($1, $7,  50, 45, 2200000),
		($2, $8,  50, 0, 450000),
		($2, $9,  50, 0, 1200000),
		($3, $10, 2000, 1200, 3500),
		($3, $11, 1000, 600, 8000)`,
		po1ID, po2ID, po3ID,
		itemIDs[0], itemIDs[1], itemIDs[2], itemIDs[3],
		itemIDs[4], itemIDs[5],
		itemIDs[6], itemIDs[7])
	if err != nil {
		log.Error("failed to seed purchase order items", "error", err)
		return
	}

	so1ID := "f0000000-0000-0000-0000-000000000001"
	so2ID := "f0000000-0000-0000-0000-000000000002"

	_, err = tx.Exec(`INSERT INTO sales_orders (id, so_number, customer_name, status, notes, created_by) VALUES
		($1, 'SO-2025-001', 'PT Mandiri Sejahtera', 'pending',  'Order laptop & monitor',  $3),
		($2, 'SO-2025-002', 'CV Berkah Abadi',       'picking',  'Order kardus & tinta',    $3)`,
		so1ID, so2ID, managerID)
	if err != nil {
		log.Error("failed to seed sales orders", "error", err)
		return
	}

	_, err = tx.Exec(`INSERT INTO sales_order_items (so_id, item_id, ordered_quantity, picked_quantity) VALUES
		($1, $3,  10, 0),
		($1, $4,  20, 0),
		($2, $5,  500, 200),
		($2, $6,  30, 10)`,
		so1ID, so2ID,
		itemIDs[0], itemIDs[3],
		itemIDs[6], itemIDs[8])
	if err != nil {
		log.Error("failed to seed sales order items", "error", err)
		return
	}

	_, err = tx.Exec(`INSERT INTO stock_transfers (transfer_number, from_location_id, to_location_id, item_id, quantity, status, notes, created_by) VALUES
		('TRF-2025-001', $1, $2, $5, 20, 'completed', 'Restock dari Zone B ke Zone A', $6),
		('TRF-2025-002', $1, $3, $4, 10, 'pending',   'Transfer ke staging area',       $6)`,
		locIDs[5], locIDs[2], locIDs[4],
		itemIDs[1], itemIDs[4], managerID)
	if err != nil {
		log.Error("failed to seed stock transfers", "error", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit seed transaction", "error", err)
		return
	}

	log.Info("seeded all data successfully",
		"users", 3,
		"items", 10,
		"locations", 10,
		"stock_entries", len(stockEntries),
		"purchase_orders", 3,
		"sales_orders", 2,
		"stock_transfers", 2,
	)
}
