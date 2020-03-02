package infinote

import (
	"infinote/canlog"
	"infinote/db"
	"context"
	"fmt"

	"github.com/gofrs/uuid"
)

// Companys will get all Companys from the DB
func Companys(ctx context.Context, os CompanyStorer) ([]*db.Company, error) {
	result, err := os.All()
	if err != nil {
		canlog.AppendErr(ctx, "147c143f-d0e8-44cd-9b6b-9e90f6c3396f")
		return nil, fmt.Errorf("list Company: %w", err)
	}
	return result, nil
}

// Company will get a single Company given an ID
func Company(ctx context.Context, us UserStorer, bs BlacklistProvider, orgID uuid.UUID) (*db.Company, error) {
	result, err := CompanyLoaderFromContext(ctx, orgID)
	if err != nil {
		canlog.AppendErr(ctx, "d63a97d4-9431-48d5-bd7f-35a6e3c4eea3")
		return nil, fmt.Errorf("get Company: %w", err)
	}
	return result, nil
}

// CompanyCreate will create a new Company
func CompanyCreate(ctx context.Context, os CompanyStorer, name string) (*db.Company, error) {

	result, err := os.Insert(&db.Company{Name: name})
	if err != nil {
		canlog.AppendErr(ctx, "fdb6ebed-357b-49d5-8330-95efc29aa51e")
		return nil, fmt.Errorf("insert Company: %w", err)
	}
	return result, nil
}

// CompanyUpdate will update an existing Company
func CompanyUpdate(ctx context.Context, os CompanyStorer, id uuid.UUID, name string) (*db.Company, error) {
	result, err := os.Update(&db.Company{Name: name})
	if err != nil {
		canlog.AppendErr(ctx, "d61a761c-d33a-4830-b636-09fc9cbb4b11")
		return nil, fmt.Errorf("update Company: %w", err)
	}
	return result, nil
}
