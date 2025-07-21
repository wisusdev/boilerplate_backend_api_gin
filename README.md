# SEMITA

## Guia de instalación

```bash
git clone https://github.com/wisusdev/boilerplate_backend_api_gin.git
cd boilerplate_backend_api_gin
cp .env.example .env
go run . key:generate

```

## Migraciones

```bash
go run . migrate

go run . migrate:refresh
go run . db:seeder
```

## OAuth2

```bash
go run . oauth:keys
go run . oauth:client
```

## Ejecutar el servidor con [Air](https://github.com/air-verse/air)

```bash
air
```
