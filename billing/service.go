package billing

import (
	"fmt"

	"encore.app/billing/activities"
	"encore.app/billing/repositories"
	"encore.app/billing/workflows"
	"encore.dev/storage/sqldb"
	"github.com/go-playground/validator/v10"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db = sqldb.NewDatabase("fee_db", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

//encore:service
type Service struct {
	client       client.Client
	worker       worker.Worker
	BillRepo     repositories.IBillRepo
	BillItemRepo repositories.IBillItemRepo
	CurrencyRepo repositories.ICurrencyRepo
	validator    *validator.Validate
}

func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}

	// repo
	BillRepo := repositories.NewBillRepo(db)
	BillItemRepo := repositories.NewBillItemRepo(db)
	CurrencyRepo := repositories.NewCurrencyRepo(db)

	// activities
	billActivity := activities.NewBillActivity(BillRepo, BillItemRepo)

	// worker
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %v", err)
	}

	w := worker.New(c, "billing-task-queue", worker.Options{})

	// register workflow
	w.RegisterWorkflow(workflows.BillingWorkflow)
	w.RegisterActivity(billActivity)

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %v", err)
	}

	return &Service{
		client:       c,
		worker:       w,
		BillRepo:     BillRepo,
		BillItemRepo: BillItemRepo,
		CurrencyRepo: CurrencyRepo,
		validator:    validator.New(),
	}, nil
}
