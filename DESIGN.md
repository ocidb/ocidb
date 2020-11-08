# OCI DB Design

Each database will be a single image.
Each table will be a new image tag.

Each write (insert/update) will be local and automatically create a transaction. 
Caller must "commit" to push to the registry backing the database.

Each instance of OCI DB will automatically "subscribe" to the registry for notification of updated tables (tags being overwritten), new tables (new tags), or deleted tables (deleted tags).
Schema will be managed by SchemaHero initially, built into the library.

