# news-api
News CRUD microservice written in Golang

Requirements for the task:

- Write a Golang microservice for CRUD operations with news.

### Approximate structure of the API:
- POST /posts
- GET /posts
- PUT /posts/{id}
- GET /posts/{id}
- DELETE /posts/{id}
Post minimum contains: id, title, content, created_at,
updated_at.

Tests are required.

### It will be an advantage:
- Use PostgreSQL
- Makefile
- Data validation
- Migration
- OpenAPI specification
- REST microservice

Upload the code to Github.

There should be startup instructions. Ideally, this is a Makefile for assembly and a docker-compose file that should raise the DB instance

The following will be assessed:
- Organization of the API
- Tests (successful and failed cases)
---
### Tech stack
- Postgres for DB
- Docker compose for local development
- Air for hot reload
- 

### How to run
1. You should create .env.local in root folder and paste these secrets:
```
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=news
DB_PORT=5432
DB_HOST=db #changed during deployment, from docker-compose if local development

POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_DB=news
```
2. Run make docker-up to start api and db locally