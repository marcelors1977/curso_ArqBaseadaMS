CREATE TABLE IF NOT EXISTS transactions (
    id varchar(255), 
    account_id_from varchar(255), 
    account_id_to varchar(255), 
    amount float, 
    created_at date, 
    updated_at date
);
