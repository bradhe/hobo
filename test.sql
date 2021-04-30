CREATE TEMPORARY TABLE orders (
  user_id INT,
  order_id INT,
  order_date DATE
);

INSERT INTO orders VALUES
  (1, 1, '2020-01-01'),
  (1, 2, '2020-01-02'),
  (2, 3, '2020-01-03'),
  (2, 4, '2020-01-04'),
  (2, 5, '2020-01-05'),
  (2, 6, '2020-01-09'),
  (3, 7, '2020-01-10'),
  (3, 9, '2020-02-10'),
  (3, 10, '2020-03-10'),
  (3, 11, '2020-04-10'),
  (3, 12, '2020-05-10'),
  (3, 13, '2020-06-10'),
  (3, 14, '2020-07-10'),
  (3, 15, '2020-08-10');

SELECT user_id, MAX(order_seq_id) AS sequential_orders FROM (
  SELECT
    user_id,
    order_id,
    order_date,
    COALESCE(order_date - LAG(order_date) OVER (PARTITION BY user_id ORDER BY order_date), 0) AS days_since_last_order,
    RANK() OVER (PARTITION BY user_id ORDER BY order_date) AS order_seq_id
  FROM orders
) dt WHERE days_since_last_order < 2
GROUP BY 1
ORDER BY 2 DESC;
