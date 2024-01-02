Este projeto faz parte dos módulos `Arquitetura baseada em microsserviços` e `EDA -Event Driven Architeture` do
curso FullCicle

A idéia do projeto é a construção de uma `Wallet`, que mantém contas de clientes, e a principal função e efetuar
transações de transferência entre contas. Estas transações serão disponibilizadas no Kafka para serem consumidas.
O projeto é composto por 2 microsserviços `walletcore_app` e `balances_app`

## Como Rodar a Aplicação

```
git clone https://github.com/marcelors1977/curso_ArqBaseadaMS.git
docker compose up -d
```

## Microsserviço walletcore_app
* Microsserviço responsável por receber requisições HTTP de criação de cliente, criação de conta e realização de 
transferência entre contas. 
* Cada transação é persistida no banco de dados cujas tabelas são `clients` para clientes, `accounts` para contas e 
`transactions` para as transações.
* Os end points para cada requisição mencionada são:

  *  criar cliente
      * request
          ```
          POST http://localhost:8080/client HTTP/1.1
          Content-Type: application/json
          
          {
              "Name": "Client1",
              "Email": "client1@email"
          }
          ```
      * response
          ```
          {
            "ID": "8407aec2-eb06-4953-9e80-617a47973d66",
            "Name": "Client1",
            "Email": "client1@email",
            "CreatedAt": "2023-12-27T19:29:47.298204927Z",
            "UpdatedAt": "0001-01-01T00:00:00Z"
          }
          ```
  * criar conta
    * request
        ```
        POST http://localhost:8080/account HTTP/1.1
        Content-Type: application/json
        
        {
            "client_id": "8407aec2-eb06-4953-9e80-617a47973d66"
        }
        ```
    * response
      ```
           {
               "ID": "0d84d645-07c3-46d7-bc46-dfc645f3839d"
           }
       ```
  * atualizar saldo da conta
    * request
      ```
         POST http://localhost:8080/account/update/0d84d645-07c3-46d7-bc46-dfc645f3839d
          Content-Type: application/json

          {
             "balance": 100
          }
      ```
  * efetuar transferência
    ```
    POST http://localhost:8080/transaction HTTP/1.1
    Content-Type: application/json
    
    {
        "account_id_from": "0d84d645-07c3-46d7-bc46-dfc645f3839d",
        "account_id_to": "336a0bf2d1804fd19fd7201a786e2a12",
        "amount": 1
    }
    ```

## Microsserviço balances_app
* Microsserviço consome os eventos enviados ao kafka por `walletcore_app` e realiza a atualização do balance do cliente
em sua base local, além de armazenar o histórico dos eventos recebidos. Também responsável por disponibilizar um
end point de consulta do saldo atualizado de uma conta. 
* Cada evento lido é persistido na base de dados `balances`, cujas tabelas são `account_balances` para armazenar o 
saldo atualizado da conta, `account_balances_history` para guardar os eventos de atualização de saldo recebidos e 
`transactions` para armazenar o histórico de transações realizadas no microsserviço walletcore_app.

* Os end points para cada requisição mencionada são:
  * cria/atualiza saldo de conta
    ```
    POST http://localhost:3003/balance HTTP/1.1
    Content-Type: application/json

    {
        "AccountID":"99999999999999999", 
        "AccountBalance":900, 
        "DateTransaction": "2023-12-02T00:00:00Z"
    }
    ```
  * busca saldo atualizado de conta
    ```
    GET http://localhost:3003/balance/search/99999999999999999 HTTP/1.1
    ```
   * insere histórico transação realizada em walletcore_app
        ```
        POST http://localhost:3003/balance HTTP/1.1
        Content-Type: application/json
    
        {
            "AccountID":"99999999999999999", 
            "AccountBalance":900, 
            "DateTransaction": "2023-12-02T00:00:00Z"
        }
        ```    
    
## Containeres
Ao executar o docker compose up 7 containeres são inicializados

----------------------------------------------------------------------------------------------------
| container              | descrição                                                               |
| ---------------------- | ----------------------------------------------------------------------- |
| walletcore_app         | microsserviço de controle de contas financeiras de clientes             |
| balances_app           | microsserviço de recebimento de eventos e visualização saldo            |
| mysql                  | banco de dados que contém os databases wallet_core e balances           |
| zookeeper              | usado pelo kafka para sincronizar configurações                         |
| kafka                  | plataforma de streaming de eventos distribuída                          |
| kafka_create_topics    | container que inicia apenas para criar os tópicos no kafka              |
| control-center         | gerenciar e monitorar o kafka através de interface web                  |
----------------------------------------------------------------------------------------------------

## Migrações e dados fake
Quando o microsserviço `walletcore_app` é inicializado, dados fakes de clientes, contas e transações são gerados e
populados nas tabelas do database `wallet_core`. Estes dados gerados e persistidos são eviados via eventos kafka
para o microsserviço `balances_app` que persiste os dados conforme o modelo de dados do database `balances`.

## Banco de dados

### wallet_core

São 3 tabelas existentes neste database usadas pelo microsserviços walletcore_app.

###### accounts

persiste dados de contas de clientes

| coluna                 | descrição                                                               |
| ---------------------- | ----------------------------------------------------------------------- |
| id                     | id único para registro de contas                                        |
| client_id              | chave estrangeira para tabela clients                                   |
| balance                | saldo disponível para a conta                                           |
| created_at             | data de criação do registro                                             |
| updated_at             | data de atualização do registro                                         |
----------------------------------------------------------------------------------------------------
    
###### clients

persiste dados referentes a clientes

| coluna                 | descrição                                                               |
| ---------------------- | ----------------------------------------------------------------------- |
| id                     | id único para registro de cliente                                       |
| name                   | nome do cliente                                                         |
| email                  | email do cliente                                                        |
| created_at             | data de criação do registro                                             |
| updated_at             | data de atualização do registro                                         |
----------------------------------------------------------------------------------------------------

###### transactions

persiste dados referente a transações de transferências entre contas

| coluna                 | descrição                                                               |
| ---------------------- | ----------------------------------------------------------------------- |
| id                     | id único pare registro de transações                                    |
| account_id_from        | id da conta de origem do valor transferido                              |
| account_id_to          | id da conta de destino do valor transferido                             |
| amount                 | montante a ser transferido                                              |
| created_at             | data de criação do registro                                             |
| updated_at             | data de atualização do registro                                         |
----------------------------------------------------------------------------------------------------

### balances

São 3 tabelas existentes neste database usadas pelo microsserviços balances_app.

###### account_balances

através dos eventos recebidos pelo microsserviço balances_app, mantém o saldo atualizado da contas
       
| coluna                 | descrição                                                               |
| ---------------------- | ----------------------------------------------------------------------- |
| id                     | id único para registro de contas                                        |
| client_id              | chave estrangeira para tabela clients                                   |
| balance                | saldo disponível para a conta                                           |
| created_at             | data de criação do registro                                             |
| updated_at             | data de atualização do registro                                         |
----------------------------------------------------------------------------------------------------
    
###### account_balances_history

persiste dados de saldo recebidos pelos eventos

| coluna                 | descrição                                                               |
| ---------------------- | ----------------------------------------------------------------------- |
| id                     | id único para registro de cliente                                       |
| name                   | nome do cliente                                                         |
| email                  | email do cliente                                                        |
| created_at             | data de criação do registro                                             |
| updated_at             | data de atualização do registro                                         |
----------------------------------------------------------------------------------------------------

###### transactions

persiste dados de transações recebidos pelos eventos

| coluna                 | descrição                                                               |
| ---------------------- | ----------------------------------------------------------------------- |
| id                     | id único pare registro de transações                                    |
| account_id_from        | id da conta de origem do valor transferido                              |
| account_id_to          | id da conta de destino do valor transferido                             |
| amount                 | montante a ser transferido                                              |
| created_at             | data de criação do registro                                             |
| updated_at             | data de atualização do registro                                         |
----------------------------------------------------------------------------------------------------

