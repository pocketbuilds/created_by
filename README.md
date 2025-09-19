# Created By

An [xpb](https://github.com/pocketbuilds/xpb) plugin for [Pocketbase](https://pocketbase.io/) that allows for easily configured created_by fields.

Configure single relation fields to auth collections to be automatically set to the auth record of the user making the post request.

## Installation

1. [Install XPB](https://docs.pocketbuilds.com/installing-xpb).
2. [Use the builder](https://docs.pocketbuilds.com/using-the-builder):

```sh
xpb build --with github.com/pocketbuilds/created_by@latest
```

## Setup
1. Create a single relation field to the auth collection.
2. Add the field to the [plugin config](#plugin-config).
3. Restart the pocketbase app.

## Plugin Config

```toml
# pocketbuilds.toml

[created_by]
# Array of single relation fields to auth collections to
#   automatically set on record create.
#   - format: "<collection_name>.<field_name>"
fields = [
  "my_collection.created_by",
  "posts.authored_by",
  "messages.user_id",
  # etc...
]
```
