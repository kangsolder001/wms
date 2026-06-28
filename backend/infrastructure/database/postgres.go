package database

import (
	"database/sql"
	"fmt"
	"time"

	"wms/config"

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
			status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending','partial','received','cancelled')),
			expected_date DATE,
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

	return nil
}
