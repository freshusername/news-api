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