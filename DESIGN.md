# OCI DB Design

Each database will be a single image.
Each table will be a new image tag.

Internally, the database is a SQLite database.
We chose this so that we can have an easy-to-use SQL langage, and not a key value store.

Each write (insert/update) will be local and automatically create a transaction. 
Caller must "commit" to push to the registry backing the database.

Each instance of OCI DB will automatically "subscribe" to the registry for notification of updated tables (tags being overwritten), new tables (new tags), or deleted tables (deleted tags).
Schema will be managed by SchemaHero initially, built into the library.

The caller should not be responsible for the quorum or anything about scaling.
This should be handled by the library automatically.
The idea is to create a new tag that stores all known connections, and they can create locks and manage state using this "table".

Open questions:
1. Should we leave a database as a single file, or split each table into a separate sqlite database and "attach" them at query time?
2. How will we "subscribe" to know when the image is updated? Maybe we can HEAD or a quick GET to know on a _very_ frequent polling? Is there anything in the OCI spec for this?
