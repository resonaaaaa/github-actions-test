package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"example.com/go_backend/internal/repository"
	"example.com/go_backend/internal/service"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("MYSQL_GO_DSN")
	if dsn == "" {
		t.Fatal("MYSQL_GO_DSN is required")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	var databaseName string
	err = db.QueryRowContext(
		ctx,
		"SELECT DATABASE()",
	).Scan(&databaseName)

	if err != nil {
		db.Close()
		t.Fatalf("failed to get database name: %v", err)
	}

	if databaseName != "go_backend_test" {
		db.Close()
		t.Fatalf(
			"unsafe database: expected go_backend_test, got %s",
			databaseName,
		)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func createTestUser(t *testing.T, db *sql.DB) int64 {
	t.Helper()

	var userID int64
	result, err := db.ExecContext(
		context.Background(),
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Test User",
		"test12@example.com",
	)

	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}

	err = db.QueryRowContext(
		context.Background(),
		"SELECT id FROM users WHERE name = ? AND email = ?",
		"Test User",
		"test12@example.com",
	).Scan(&userID)

	if err != nil {
		t.Fatalf("failed to get test user ID: %v", err)
	}

	return userID
}

func createTestProduct(t *testing.T, db *sql.DB) int64 {
	t.Helper()

	var productID int64
	result, err := db.ExecContext(
		context.Background(),
		"INSERT INTO products (name, price_cents, stock) VALUES (?, ?, ?)",
		"Test Product",
		1000,
		10,
	)
	if err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}

	err = db.QueryRowContext(
		context.Background(),
		"SELECT id FROM products WHERE name = ? AND price_cents = ?",
		"Test Product",
		1000,
	).Scan(&productID)

	if err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	return productID
}

func TestCreateOrderIntegration(t *testing.T) {
	t.Helper()
	db := openTestDB(t)

	// 清理测试数据
	resetTestData(t, db)

	// 准备测试商品与用户
	userID := createTestUser(t, db)
	productID := createTestProduct(t, db)

	// 创建OrderService实例
	uow := repository.NewUnitOfWork(db)
	orderReader := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(uow, orderReader)

	order, err := orderService.CreateOrder(context.Background(), userID, productID, 2, "test-request-id")
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	if order.UserID != userID || order.ProductID != productID || order.Quantity != 2 {
		t.Fatalf("unexpected order data: %+v", order)
	}

	err = db.QueryRowContext(context.Background(), "SELECT id, user_id, product_id, quantity FROM orders WHERE id = ?", order.ID).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity)
	if err != nil {
		t.Fatalf("failed to query order from database: %v", err)
	}
	if order.UserID != userID || order.ProductID != productID || order.Quantity != 2 {
		t.Fatalf("unexpected order data in database: %+v", order)
	}

}

func TestCreateOrderIdempotency(t *testing.T) {
	t.Helper()
	db := openTestDB(t)

	resetTestData(t, db)

	userID := createTestUser(t, db)
	productID := createTestProduct(t, db)

	uow := repository.NewUnitOfWork(db)
	orderReader := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(uow, orderReader)

	// 第一次创建订单
	order1, err := orderService.CreateOrder(context.Background(), userID, productID, 2, "test-request-id")
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	// 第二次使用相同的RequestID创建订单，应该返回相同的订单
	order2, err := orderService.CreateOrder(context.Background(), userID, productID, 2, "test-request-id")
	if err != nil {
		t.Fatalf("failed to create order with same request ID: %v", err)
	}
	if order1.ID != order2.ID {
		t.Fatalf("expected same order ID for idempotent request, got %d and %d", order1.ID, order2.ID)
	}
	if order1.RequestID != order2.RequestID {
		t.Fatalf("expected same request ID for idempotent request, got %s and %s", order1.RequestID, order2.RequestID)
	}
	if order1.UserID != order2.UserID || order1.ProductID != order2.ProductID || order1.Quantity != order2.Quantity {
		t.Fatalf("expected same order data for idempotent request, got %+v and %+v", order1, order2)
	}
	if order1.UnitPrice != order2.UnitPrice {
		t.Fatalf("expected same unit price for idempotent request, got %d and %d", order1.UnitPrice, order2.UnitPrice)
	}

}

func TestCreateOrderInsufficientStock(t *testing.T) {
	t.Helper()
	db := openTestDB(t)
	resetTestData(t, db)

	userID := createTestUser(t, db)
	productID := createTestProduct(t, db)

	uow := repository.NewUnitOfWork(db)
	orderReader := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(uow, orderReader)

	// 尝试创建一个数量超过库存的订单
	_, err := orderService.CreateOrder(context.Background(), userID, productID, 100, "test-request-id")
	if err == nil {
		t.Fatalf("expected error when creating order with insufficient stock, but got none")
	}
	if !errors.Is(err, service.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
	row := db.QueryRowContext(context.Background(), "SELECT stock FROM products WHERE id = ?", productID)
	var stock int
	if err := row.Scan(&stock); err != nil {
		t.Fatalf("failed to query product stock: %v", err)
	}
	if stock != 10 {
		t.Fatalf("expected stock to remain 10 after failed order creation, got %d", stock)
	}

}

func TestCreateOrderConcurrent(t *testing.T) {
	t.Helper()
	db := openTestDB(t)
	resetTestData(t, db)

	userID := createTestUser(t, db)
	productID := createTestProduct(t, db)

	uow := repository.NewUnitOfWork(db)
	orderReader := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(uow, orderReader)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := orderService.CreateOrder(context.Background(), userID, productID, 1, fmt.Sprintf("request-%d", i))
			if err != nil {
				t.Errorf("failed to create order: %v", err)
			}
		}(i)
	}
	wg.Wait()

	row := db.QueryRowContext(context.Background(), "SELECT stock FROM products WHERE id = ?", productID)
	var stock int
	if err := row.Scan(&stock); err != nil {
		t.Fatalf("failed to query product stock: %v", err)
	}
	if stock != 0 {
		t.Fatalf("expected stock to be 0 after concurrent order creation, got %d", stock)
	}

}

func TestCreateOrderConcurrentInsufficientStock(t *testing.T) {
	t.Helper()
	db := openTestDB(t)
	resetTestData(t, db)

	userID := createTestUser(t, db)
	productID := createTestProduct(t, db)

	uow := repository.NewUnitOfWork(db)
	orderReader := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(uow, orderReader)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := orderService.CreateOrder(context.Background(), userID, productID, 2, fmt.Sprintf("request-%d", i))
			if err != nil && !errors.Is(err, service.ErrInsufficientStock) {
				if !errors.Is(err, service.ErrMaxRetries) {
					t.Errorf("failed to create order: %v", err)
				}
			}
		}(i)
	}
	wg.Wait()

	rows, err := db.QueryContext(context.Background(), "SELECT COUNT(*) FROM orders WHERE user_id = ?", userID)
	if err != nil {
		t.Fatalf("failed to query order count: %v", err)
	}
	defer rows.Close()

	var count int

	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			t.Fatalf("failed to scan order count: %v", err)
		}
	}

	if count != 5 {
		t.Errorf("expected 5 orders, got %d", count)
	}

}

func MigrateUp() error {
	dbURL := os.Getenv("MIGRATE_TEST_DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("MIGRATE_TEST_DATABASE_URL is required")
	}

	goDSN := os.Getenv("MYSQL_GO_DSN")
	if goDSN == "" {
		return fmt.Errorf("MYSQL_GO_DSN is required")
	}
	goConfig, err := mysql.ParseDSN(goDSN)
	if err != nil {
		return fmt.Errorf("invalid MYSQL_GO_DSN: %w", err)
	}

	// 确保两个连接都指向 go_backend_test 数据库
	migrateConfig, err := mysql.ParseDSN(strings.TrimPrefix(dbURL, "mysql://"))
	if err != nil {
		return fmt.Errorf("invalid MIGRATE_TEST_DATABASE_URL: %w", err)
	}
	if goConfig.DBName != "go_backend_test" || migrateConfig.DBName != goConfig.DBName {
		return fmt.Errorf("unsafe test database: both connections must target go_backend_test")
	}

	mig, err := migrate.New("file://../../migrations", dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer mig.Close()

	err = mig.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

func resetTestData(t *testing.T, db *sql.DB) {
	if _, err := db.ExecContext(context.Background(), "DELETE FROM orders"); err != nil {
		t.Fatalf("failed to delete orders: %v", err)
	}

	if _, err := db.ExecContext(context.Background(), "DELETE FROM users"); err != nil {
		t.Fatalf("failed to delete users: %v", err)
	}

	if _, err := db.ExecContext(context.Background(), "DELETE FROM products"); err != nil {
		t.Fatalf("failed to delete products: %v", err)
	}

}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("failed to load .env: %v\n", err)
		os.Exit(1)
	}
	if os.Getenv("MYSQL_GO_DSN") == "" {
		if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
			if err := os.Setenv("MYSQL_GO_DSN", dsn); err != nil {
				fmt.Printf("failed to set MYSQL_GO_DSN: %v\n", err)
				os.Exit(1)
			}
		}
	}

	if err := MigrateUp(); err != nil {
		fmt.Printf("failed to apply migrations: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}
