# Description
A ScyllaX PJP Golang

# Features
This ScyllaX PJP Golang Project with Go, Gin-Gonic, Gorm, JWT, Postgresql, Docker
 
# Tech Used
 ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
      
# Getting Start:
Before you running the program, make sure you've run this command:
```bash
 make install
```
```bash
 docker compose --file ./docker-compose.yml --env-file development.env up -d
```

### Run the program
```bash
 make dev
```
### Run program live reload
```bash
 make dev-reload
```

### Re-Init Docs Swagger
```bash
 make doc
```

### Create Table Sql
```bash
  make migration name={name table and alter table}
```

### Migrate Up
```bash
  make migrateUp
```

### Migrate Down
```bash
  make migrateDown
```

### Migrate Fix
```bash
  make migrateForce version={version}
```

### Migrate Drop
```bash
  make migrateDrop
```

### Check Docs Swagger
```bash
 http://localhost:8888/docs/index.html#/
```

The program will run on http://localhost:8888 
<!-- </> with 💛 by readMD (https://readmd.itsvg.in) -->
