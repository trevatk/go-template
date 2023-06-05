# Go Template

Project to be used as a template for future projects. Template is a hybrid flat architecture with DDD driven packaging. 
This allows for a loose architecture but strict package locations. This architecture is also great as you get all the benefits 
of DDD without all the boilerplate required. To accomplish this the application layer is coded while the business layer is 
generated using [sqlc](https://sqlc.dev/) and migrated using [golang-migrate](https://github.com/golang-migrate/migrate).
