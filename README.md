# E-card-catalog

Rest API with authentication

## Built with

-   [Go](https://golang.org/ 'The Go Programming language')
-   [PostgreSQL](https://www.postgresql.org/ 'PostgreSQL')

## API

```
sign in
POST /api/auth/signin

sign up
POST /api/auth/signup

authentication required

get all schools
GET /api/schools/getAll

get all books
GET /api/books/getAll

get books considering query params
GET /api/books/filter

get all cards
GET /api/cards/getAll

get one card (and books assigned to it)
GET /api/cards/getOne/{id}

get cards (and books assigned to them) considering query params
GET /api/cards/filter

create a card
POST /api/cards/create

edit a card
PUT /api/cards/edit/{id}

assign book to a card
GET /api/cards/attachBook/{book_id}

remove book from a card
GET /api/cards/detachBook/{book_id}

delete a card
DELETE /api/cards/delete/{id}
```
