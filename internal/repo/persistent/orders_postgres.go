package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/andreyxaxa/order_svc/pkg/postgres"
)

const (
	ordersTable     = "orders"
	deliveriesTable = "deliveries"
	paymentsTable   = "payments"
	itemsTable      = "items"
)

type OrdersRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *OrdersRepo {
	return &OrdersRepo{pg}
}

func (r *OrdersRepo) Store(ctx context.Context, o entity.Order) error {
	now := time.Now().UTC()

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("OrdersRepo - Store - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	// orders
	sql, args, err := r.Builder.
		Insert(ordersTable).
		Columns(
			"order_uid, track_number, entry, locale, internal_signature, customer_id",
			"delivery_service, shardkey, sm_id, date_created, oof_shard",
		).
		Values(
			o.OrderUID,
			o.TrackNumber,
			o.Entry,
			o.Locale,
			o.InternalSignature,
			o.CustomerID,
			o.DeliveryService,
			o.ShardKey,
			o.SmID,
			now,
			o.OofShard,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrdersRepo - Store - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OrdersRepo - Store - tx.Exec: %w", err)
	}

	// deliveries
	sql, args, err = r.Builder.
		Insert(deliveriesTable).
		Columns("order_uid, name, phone, zip, city, address, region, email").
		Values(
			o.OrderUID,
			o.Delivery.Name,
			o.Delivery.Phone,
			o.Delivery.Zip,
			o.Delivery.City,
			o.Delivery.Address,
			o.Delivery.Region,
			o.Delivery.Email,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrdersRepo - Store - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OrdersRepo - Store - tx.Exec: %w", err)
	}

	// payments
	sql, args, err = r.Builder.
		Insert(paymentsTable).
		Columns(
			"transaction, order_uid, request_id, currency, provider, amount",
			"payment_dt, bank, delivery_cost, goods_total, custom_fee").
		Values(
			o.Payment.Transaction,
			o.OrderUID,
			o.Payment.RequestID,
			o.Payment.Currency,
			o.Payment.Provider,
			o.Payment.Amount,
			o.Payment.PaymentDT,
			o.Payment.Bank,
			o.Payment.DeliveryCost,
			o.Payment.GoodsTotal,
			o.Payment.CustomFee,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrdersRepo - Store - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OrdersRepo - Store - tx.Exec: %w", err)
	}

	// items
	for _, item := range o.Items {
		sql, args, err = r.Builder.
			Insert(itemsTable).
			Columns("order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status").
			Values(
				o.OrderUID,
				item.ChrtID,
				item.TrackNumber,
				item.Price,
				item.RID,
				item.Name,
				item.Sale,
				item.Size,
				item.TotalPrice,
				item.NmID,
				item.Brand,
				item.Status,
			).
			ToSql()
		if err != nil {
			return fmt.Errorf("OrdersRepo - Store - r.Builder: %w", err)
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("OrdersRepo - Store - tx.Exec: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("OrdersRepo - Store - tx.Commit: %w", err)
	}

	return nil
}

func (r *OrdersRepo) GetOrder(ctx context.Context, orderUID string) (entity.Order, error) {
	var o entity.Order

	// orders
	sql, args, err := r.Builder.
		Select(
			"order_uid, track_number, entry, locale, internal_signature, customer_id",
			"delivery_service, shardkey, sm_id, date_created, oof_shard").
		From(ordersTable).
		Where(squirrel.Eq{"order_uid": orderUID}).
		ToSql()
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&o.OrderUID,
		&o.TrackNumber,
		&o.Entry,
		&o.Locale,
		&o.InternalSignature,
		&o.CustomerID,
		&o.DeliveryService,
		&o.ShardKey,
		&o.SmID,
		&o.DateCreated,
		&o.OofShard,
	)
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - row.Scan: %w", err)
	}

	// deliveries
	sql, args, err = r.Builder.
		Select("name, phone, zip, city, address, region, email").
		From(deliveriesTable).
		Where(squirrel.Eq{"order_uid": orderUID}).
		ToSql()
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - r.Builder: %w", err)
	}

	row = r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&o.Delivery.Name,
		&o.Delivery.Phone,
		&o.Delivery.Zip,
		&o.Delivery.City,
		&o.Delivery.Address,
		&o.Delivery.Region,
		&o.Delivery.Email,
	)
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - row.Scan: %w", err)
	}

	// payments
	sql, args, err = r.Builder.
		Select(
			"transaction, request_id, currency, provider, amount",
			"payment_dt, bank, delivery_cost, goods_total, custom_fee").
		From(paymentsTable).
		Where(squirrel.Eq{"order_uid": orderUID}).
		ToSql()
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - r.Builder: %w", err)
	}

	row = r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(
		&o.Payment.Transaction,
		&o.Payment.RequestID,
		&o.Payment.Currency,
		&o.Payment.Provider,
		&o.Payment.Amount,
		&o.Payment.PaymentDT,
		&o.Payment.Bank,
		&o.Payment.DeliveryCost,
		&o.Payment.GoodsTotal,
		&o.Payment.CustomFee,
	)
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - row.Scan: %w", err)
	}

	// items
	sql, args, err = r.Builder.
		Select("chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status").
		From(itemsTable).
		Where(squirrel.Eq{"order_uid": orderUID}).
		ToSql()
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return o, fmt.Errorf("OrdersRepo - GetOrder - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		); err != nil {
			return o, fmt.Errorf("OrdersRepo - GetOrder - rows.Scan: %w", err)
		}
		o.Items = append(o.Items, item)
	}

	return o, nil
}
