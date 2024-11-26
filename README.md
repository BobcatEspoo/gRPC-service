# GRPCBank
Функционал данного микросервиса заключается в отправке запросов через gRPC для получения курса USDT с биржи Garantex.
Для запуска приложение использовать команды(убедитесь что у вас установлен docker):
```bash
git clone https://github.com/BobcatEspoo/gRPC-service.git
cd gRPC-service
docker compose up

Для изменения базы данных используется .env файл, пример:
```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=grpcdb
DB_HOST=db
DB_PORT=5432

DB_URL=postgres://postgres:yourpassword@db:5432/grpcdb?sslmode=disable

 
