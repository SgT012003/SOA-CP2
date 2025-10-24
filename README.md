# Project Hotel API (SOA)

This is a basic api with a swagger built-in.

## Group
|Name|RM|
|:-:|:-:|
|Diogo Julio|553837|
|Jonata Rafael|552939|
|Matheus Zottis|94119|
|Victor Didoff|552965|
|Vinicius da Silva|553240|

## Requirements

|Software|Version|Build Type|
|:-:|:-:|:-:|
|Docker|20+|Docker|
|Docker Compose|1.29+|Docker|
|Golang|1.24+|Local|
|Postgres|17+|Local|
|Swag CLI|1.16+|Local - Swagger Rebuild|

> Observation : the seeding data to db its already on the application building pipe-line.

## First Boot Into Project

1. Clone the Repo and acess the project dir

```bash
    git clone https://github.com/SgT012003/SOA-CP2.git && cd SOA-CP2
```

2. After Downloading Docker and Docker Composer run the following command on your bash or cmd terminal

```bash
    docker-compose up --build
```

>this command will build the project into the container and start it.

3. Access the [Swagger UI](#accessing-the-api)


## Docker

This is a basic tutorial to how to setup docker and basic commands

### Start
- flag: --build (mandatory: on first start)
```bash 
    docker-compose up --build
```

### Stop
- flag: -v (optional: dont keep data persistance [docker images])
```bash 
    docker-compose down
```

## Configs

1. **Defalt Config**
```dockerfile
    # Comando para rodar a aplicação
    # CMD ["./main"]

    # Comando para rodar a aplicação com dados de teste
    CMD ["sh", "-c", "./setup && ./main"]
```
---

2. **No DB setup**

> [Dockerfile](./Dockerfile#L32) -> IF you don't want to start DB with test data.

```dockerfile
    # Comando para rodar a aplicação
    CMD ["./main"]
    # CMD ["sh", "-c", "./setup && ./main"]
```

- Don't Create Basis Tables (products and clients)
- Don't Seed Tables with initial Data for Testing Porpouse

## Envs

**Default Docker-Postgres DB** -> [docker-compose.yml](./docker-compose.yml#L7)
```yml
    environment:
        POSTGRES_USER: hotel_dba
        POSTGRES_PASSWORD: 12345678
        POSTGRES_DB: hotel
```

**Default Golang DTO** -> [docker-compose.yml](./docker-compose.yml#L23)
```yml
    environment:
        DB_HOST: db
        DB_PORT: 5432
        DB_USER: hotel_dba
        DB_PASSWORD: 12345678
        DB_NAME: hotel
```

**or**

Set local `.env` to this values with your remote or local postgres db:
- DB_PORT
- DB_PORT
- DB_USER
- DB_PASSWORD
- DB_NAME


## Accessing the API


1. **Redirect EZ to Swagger UI** : this root access redirect you to the swagger ui page.
```uri
http://localhost:8080
``` 

2. **Swagger Full Path**
```uri
http://localhost:8080/swagger/index.html
``` 