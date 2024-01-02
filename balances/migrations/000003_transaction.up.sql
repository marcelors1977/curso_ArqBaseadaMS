CREATE TABLE IF NOT EXISTS transactions (
  account_id_from varchar(255),
  account_id_to varchar(255),
  amount float,
  date_transaction datetime
);