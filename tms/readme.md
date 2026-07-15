# Description

ScyllaX-TMS Api Golang

# Features

This ScyllaX-TMS Api Golang Project with Go, Fiber, Gorm, JWT, Postgresql, Swagger

# Tech Used

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white) ![Gin Badge](https://img.shields.io/badge/Gin-008ECF?logo=gin&logoColor=fff&style=for-the-badge)

# Getting Start:

Before you running the program, make sure you've run this command:

```bash
 make install
```

### Run the program

```bash
 make dev
```

### Re-Init Docs Swagger

```bash
 make doc
```

### Create Table Sql

```bash
 make migration command={create_table_name or alter_table_name_add_email}
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
  make migrateForce command={version}
```

### Migrate Drop

```bash
  make migrateDrop
```

### Check Docs Swagger

```bash
 http://localhost:3000/docs/index.html#/
```

The program will run on http://localhost:3000

<!-- </> with 💛 by readMD (https://readmd.itsvg.in) -->
