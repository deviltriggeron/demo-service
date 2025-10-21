package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	e "demo-service/internal/entity"
	p "demo-service/internal/interface"
)

type postgres struct {
	db *sql.DB
}

func InitDB(cfg *e.ConfigDB) (p.OrderRepository, error) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := RunMigration(db, "up"); err != nil {
		return nil, err
	}

	fmt.Println("Database initialized successfully!")
	return &postgres{db: db}, nil
}

func (p *postgres) CloseDB() error {
	return p.db.Close()
}

func (p *postgres) MethodSelect(id string) (e.Order, error) {
	var order e.Order

	tx, err := p.db.Begin()
	if err != nil {
		return order, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	queryOrder := `
		SELECT order_uid, track_number, entry, locale, internal_signature, 
		       customer_id, delivery_service, shardkey, sm_id, 
		       date_created, oof_shard
		FROM orders WHERE order_uid = $1`

	err = p.db.QueryRow(queryOrder, id).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSig,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order, fmt.Errorf("order %s not found", id)
		}
		return order, err
	}

	queryDelivery := `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries WHERE order_uid = $1`

	err = p.db.QueryRow(queryDelivery, id).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
	)
	if err != nil {
		return order, err
	}

	queryPayment := `
		SELECT transaction, request_id, currency, provider, amount,
		       payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uid = $1`
	err = p.db.QueryRow(queryPayment, id).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil {
		return order, err
	}
	order.Payment.OrderUID = id

	queryItems := `
		SELECT chrt_id, track_number, price, rid, name, sale, size, 
		       total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1`
	rows, err := p.db.Query(queryItems, id)
	if err != nil {
		return order, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Fatalf("error close rows: %v", err)
		}
	}()

	for rows.Next() {
		var item e.Item
		item.OrderUID = id
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}

	order.Delivery.OrderUID = id

	if err := tx.Commit(); err != nil {
		return order, fmt.Errorf("transaction commit failed: %w", err)
	}

	return order, nil
}

func (p *postgres) MethodSelectAll() ([]e.Order, error) {
	var orders []e.Order

	tx, err := p.db.Begin()
	if err != nil {
		return orders, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	rows, err := p.db.Query(`
        SELECT order_uid, track_number, entry, locale, internal_signature,
               customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders
    `)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Fatalf("error close rows: %v", err)
		}
	}()

	for rows.Next() {
		var o e.Order
		err := rows.Scan(
			&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSig,
			&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard,
		)
		if err != nil {
			return nil, err
		}

		err = p.db.QueryRow(`
            SELECT name, phone, zip, city, address, region, email
            FROM delivery WHERE order_uid = $1
        `, o.OrderUID).Scan(
			&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City,
			&o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
		)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		err = p.db.QueryRow(`
            SELECT transaction, request_id, currency, provider, amount, payment_dt,
                   bank, delivery_cost, goods_total, custom_fee
            FROM payment WHERE order_uid = $1
        `, o.OrderUID).Scan(
			&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency,
			&o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDT,
			&o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
		)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		itemRows, err := p.db.Query(`
            SELECT chrt_id, track_number, price, rid, name, sale,
                   size, total_price, nm_id, brand, status
            FROM items WHERE order_uid = $1
        `, o.OrderUID)
		if err != nil {
			return nil, err
		}

		defer func() {
			if err = itemRows.Close(); err != nil {
				log.Fatalf("error close rows: %v", err)
			}
		}()
		for itemRows.Next() {
			var it e.Item
			if err := itemRows.Scan(
				&it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid, &it.Name,
				&it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status,
			); err != nil {
				return nil, err
			}
			o.Items = append(o.Items, it)
		}

		orders = append(orders, o)
	}

	if err := tx.Commit(); err != nil {
		return orders, fmt.Errorf("transaction commit failed: %w", err)
	}

	return orders, nil
}

func (p *postgres) MethodDelete(orderUID string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = p.db.Exec(`DELETE FROM items WHERE order_uid = $1`, orderUID); err != nil {
		return err
	}

	if _, err = p.db.Exec(`DELETE FROM payment WHERE order_uid = $1`, orderUID); err != nil {
		return err
	}

	if _, err = p.db.Exec(`DELETE FROM delivery WHERE order_uid = $1`, orderUID); err != nil {
		return err
	}

	if _, err = p.db.Exec(`DELETE FROM orders WHERE order_uid = $1`, orderUID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (p *postgres) MethodInsert(orders e.Order) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = p.db.Exec(`
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature,
                            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
    `, orders.OrderUID, orders.TrackNumber, orders.Entry, orders.Locale, orders.InternalSig,
		orders.CustomerID, orders.DeliveryService, orders.Shardkey, orders.SmID, orders.DateCreated, orders.OofShard)
	if err != nil {
		return fmt.Errorf("insert orders: %w", err)
	}

	_, err = p.db.Exec(`
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, orders.Delivery.OrderUID, orders.Delivery.Name, orders.Delivery.Phone, orders.Delivery.Zip, orders.Delivery.City, orders.Delivery.Address, orders.Delivery.Region, orders.Delivery.Email)
	if err != nil {
		return fmt.Errorf("insert delivery: %w", err)
	}

	_, err = p.db.Exec(`
		INSERT INTO payment (order_uid, "transaction", request_id, currency, provider, amount, payment_dt,
							bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		`, orders.Payment.OrderUID, orders.Payment.Transaction, orders.Payment.RequestID, orders.Payment.Currency, orders.Payment.Provider, orders.Payment.Amount, orders.Payment.PaymentDT, orders.Payment.Bank, orders.Payment.DeliveryCost, orders.Payment.GoodsTotal, orders.Payment.CustomFee)

	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	for _, it := range orders.Items {
		_, err = p.db.Exec(`
		INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale,
							size, total_price, nm_id, brand, status) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		`, it.OrderUID, it.ChrtID, it.TrackNumber, it.Price, it.Rid, it.Name, it.Sale,
			it.Size, it.TotalPrice, it.NmID, it.Brand, it.Status,
		)
		if err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (p *postgres) MethodUpdate(order e.Order) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	queryOrder := `
		UPDATE orders SET
			track_number = $1,
			entry = $2,
			locale = $3,
			internal_signature = $4,
			customer_id = $5,
			delivery_service = $6,
			shardkey = $7,
			sm_id = $8,
			date_created = $9,
			oof_shard = $10
		WHERE order_uid = $11`
	_, err = tx.Exec(queryOrder,
		order.TrackNumber, order.Entry, order.Locale, order.InternalSig,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID,
		order.DateCreated, order.OofShard, order.OrderUID,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Fatalf("error rollback: %v", rollbackErr)
		}
		return fmt.Errorf("update orders: %w", err)
	}

	queryDelivery := `
		UPDATE delivery SET
			name = $1, phone = $2, zip = $3, city = $4,
			address = $5, region = $6, email = $7
		WHERE order_uid = $8`
	_, err = tx.Exec(queryDelivery,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, order.OrderUID,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Fatalf("error rollback: %v", rollbackErr)
		}
		return fmt.Errorf("update delivery: %w", err)
	}

	queryPayment := `
		UPDATE payment SET
			transaction = $1, request_id = $2, currency = $3, provider = $4,
			amount = $5, payment_dt = $6, bank = $7, delivery_cost = $8,
			goods_total = $9, custom_fee = $10
		WHERE order_uid = $11`
	_, err = tx.Exec(queryPayment,
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee, order.OrderUID,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Fatalf("error rollback: %v", rollbackErr)
		}
		return fmt.Errorf("update payment: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM items WHERE order_uid = $1`, order.OrderUID); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Fatalf("error rollback: %v", rollbackErr)
		}
		return fmt.Errorf("delete old items: %w", err)
	}

	insertItem := `
		INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale,
						   size, total_price, nm_id, brand, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	for _, it := range order.Items {
		_, err := tx.Exec(insertItem,
			order.OrderUID, it.ChrtID, it.TrackNumber, it.Price, it.Rid, it.Name, it.Sale,
			it.Size, it.TotalPrice, it.NmID, it.Brand, it.Status,
		)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Fatalf("error rollback: %v", rollbackErr)
			}
			return fmt.Errorf("insert item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
