package migration

import (
	"database/sql"
	"fmt"
)

// helper: check table existence
func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		);
	`
	var exists bool
	err := db.QueryRow(query, tableName).Scan(&exists)
	return exists, err
}

// main function
func EnsureTablesExistWithLogs(db *sql.DB) ([]string, error) {
	logs := []string{}
	queries := map[string]string{

		// ---------- EXTENSION ----------
		"uuid_extension": `
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		`,

		// ---------- ADMINS ----------
		"admins": `
			CREATE TABLE admins (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash TEXT NOT NULL,
				full_name VARCHAR(255),
				is_active BOOLEAN DEFAULT true,
				created_at TIMESTAMP DEFAULT now(),
				updated_at TIMESTAMP DEFAULT now(),
				deleted_at TIMESTAMP
			);
		`,

		"admin_roles": `
			CREATE TABLE admin_roles (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				name VARCHAR(100) NOT NULL,
				deleted_at TIMESTAMP
			);
		`,

		"admin_user_roles": `
			CREATE TABLE admin_user_roles (
				admins_uuid UUID REFERENCES admins(uuid),
				admin_roles_uuid UUID REFERENCES admin_roles(uuid),
				PRIMARY KEY (admins_uuid, admin_roles_uuid)
			);
		`,

		// ---------- USERS ----------
		"users": `
			CREATE TABLE users (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				personal_no VARCHAR(50),
				full_name VARCHAR(255),
				email VARCHAR(255),
				dob DATE,
				gender VARCHAR(20),
				profile_photos TEXT,
				institution_name VARCHAR(255),
				is_subscribed BOOLEAN DEFAULT false,
				created_at TIMESTAMP DEFAULT now(),
				updated_at TIMESTAMP DEFAULT now(),
				deleted_at TIMESTAMP
			);
		`,

		"user_roles": `
			CREATE TABLE user_roles (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				code VARCHAR(50) UNIQUE NOT NULL,
				description TEXT
			);
		`,

		"user_role_assignments": `
			CREATE TABLE user_role_assignments (
				users_uuid UUID REFERENCES users(uuid),
				user_roles_uuid UUID REFERENCES user_roles(uuid),
				assigned_at TIMESTAMP DEFAULT now(),
				PRIMARY KEY (users_uuid, user_roles_uuid)
			);
		`,

		// ---------- MASTER TABLES ----------
		"designations": `
			CREATE TABLE designations (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				title VARCHAR(255)
			);
		`,

		"departments": `
			CREATE TABLE departments (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				name VARCHAR(255)
			);
		`,

		"means_of_transport": `
			CREATE TABLE means_of_transport (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				name VARCHAR(100)
			);
		`,

		"office_locations": `
			CREATE TABLE office_locations (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				name VARCHAR(255)
			);
		`,

		"residency_areas": `
			CREATE TABLE residency_areas (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				name VARCHAR(255)
			);
		`,

		// ---------- USER PROFILE ----------
		"user_profiles": `
			CREATE TABLE user_profiles (
				uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				users_uuid UUID REFERENCES users(uuid),
				designations_uuid UUID REFERENCES designations(uuid),
				departments_uuid UUID REFERENCES departments(uuid),
				means_of_transport_uuid UUID REFERENCES means_of_transport(uuid),
				office_locations_uuid UUID REFERENCES office_locations(uuid),
				residency_areas_uuid UUID REFERENCES residency_areas(uuid),
				home_to_office_distance DECIMAL,
				wfh_days INT,
				absent_days INT,
				signature TEXT,
				deleted_at TIMESTAMP
			);
		`,
	}

	// extension
	if _, err := db.Exec(queries["uuid_extension"]); err != nil {
		return logs, err
	}
	logs = append(logs, "uuid-ossp extension checked")

	for table, stmt := range queries {
		if table == "uuid_extension" {
			continue
		}

		exists, err := tableExists(db, table)
		if err != nil {
			return logs, err
		}

		if exists {
			logs = append(logs, fmt.Sprintf("table %s already exists", table))
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			return logs, err
		}
		logs = append(logs, fmt.Sprintf("table %s created", table))
	}

	return logs, nil
}
