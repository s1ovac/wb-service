
--INSERT INTO item (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
--VALUES ('0021e010-97ab-46db-8600-f7604ab52f92', 9934932, 'WBILMTESTTRACK2', 453, 'ab4219087a764ae0btest3', 'MascarasVascaras', 15, '2', 560, 238214, 'Nike', 205);

SELECT 
		o."order_uid",
		o."track_number",
		o."entry",
		d."name",
		d."phone",
		d."zip",
		d."city",
		d."address",
		d."region",
		d."email",
		p."transaction",
		p."request_id",
		p."currency",
		p."provider",
		p."amount",
		p."payment_dt",
		p."bank",
		p."delivery_cost",
		p."goods_total",
		p."custom_fee",
		o."locale",
		o."internal_signature",
		o."customer_id",
		o."delivery_service", 
		o."shardkey", 
		o."sm_id", 
		o."date_created", 
		o."oof_shard"
	FROM 
		"order" AS o 
		JOIN "delivery" AS d ON o."order_uid" = d."order_id"
		JOIN "payment" AS p ON d."order_id" = p."order_id"
		--JOIN "item" AS i ON p."order_id" = i."order_id"
	WHERE o.order_uid = '0021e010-97ab-46db-8600-f7604ab52f92';