# Blog API with social interactions (developing)

Basic personal blog, built it with Go used gintonic framework and postgrest. This is only part of my portfolio

## Structure

### Models
In this package are define all structures and services(functions) for CRUD operations. Also Manage the db connection.

### Controllers
Contains web controllers using gintonic framework that use services provided by models

### Config
Has the env variables how api port or db user.

