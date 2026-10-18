# Reading Detection Warnings

Each file in `output/{app}/analysis/` starts with `[NUM_WARNINGS = N]`, followed by one block per warning. All files use the notation below.

## Notation

### Call Paths

A call path shows which service method runs a database operation. Calls in between are left out:

```txt
StorageService.ReadPost() ... posts_db.post.FindOne()
```

Here, `StorageService.ReadPost()` runs `FindOne()` on entity `post` in database `posts_db`. Some paths also start with the entry point of the request:

```txt
UploadService.UploadPost() ... StorageService.StorePost() ... posts_db.post.InsertOne()
```

### Fields

Fields are written as `database.entity.field`:

| Field | Meaning |
| ----- | ------- |
| `posts_db.post.Creator.Username` | The `Username` of the post's `Creator` |
| `posts_db.post.Mentions[*]` | Every element of the list `Mentions` |

### Foreign Keys

Foreign keys are inferred by Aletheia:

```txt
FOREIGN_KEY notif_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]
```

Here, each notification's `PostID` points to a post in `posts_db`. A foreign key can have two tags:

| Tag | Meaning |
| --- | ------- |
| `[MANDATORY]` | Every request that creates the reference also writes the record it points to. Here, `UploadService.UploadPost()` writes both the post and its notification. |
| `[T]` | Transitive: derived from two other foreign keys (if X references Y and Y references Z, then X references Z). |

## Example Warnings

Below is an example warning for each pattern, with an explanation of its parts. See [What Aletheia Detects](../patterns.md) for what each pattern means.

### RI-1: Referential integrity, absence of cascading deletes

```txt title="foreign-key-cascade.txt"
delete: ProductService.Delete() ... products_db.products.DeleteOne()
    missing cascade #1: database={skus_db}, entity={skus}, pending_fields={Parent}
```

- `delete:` is the operation that deletes a record.
- `missing cascade #N` names a database and entity whose records can still reference the deleted record, through the fields in `pending_fields`.

### RI-2: Referential integrity, concurrent operations

```txt title="foreign-key-concurrency.txt"
delete: ProductService.Delete() ... products_db.products.DeleteOne()
    write #1: SkuService.New() ... SkuService.New() ... skus_db.skus.InsertOne()
        - database={skus_db}, entity={skus}, written_fields={Parent}
```

- `delete:` is the operation that deletes a record.
- `write #N` is an operation in another request that can store a reference to that record at the same time, through the fields in `written_fields`.

### RI-3: Referential integrity, uncoordinated replication

```txt title="foreign-key-coordination.txt"
entry request: UploadService.UploadPost()
    FOREIGN KEY READS #1:
        READ (FOREIGN KEY): NotifyService.Run() ... notif_queue.notification.Pop()
            - field: notif_queue.notification.PostID
            - constraint: FOREIGN_KEY notif_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]
        READ (ORIGIN): StorageService.ReadPost() ... posts_db.post.FindOne()
            - field: posts_db.post.PostID
```

- `entry request` is the client request that leads to the two reads.
- `READ (FOREIGN KEY)` reads a record that holds a reference (`field`, following `constraint`).
- `READ (ORIGIN)` reads the record that the reference points to.

### EI-1: Entity integrity, uncoordinated replication

```txt title="primary-key-coordination.txt"
entry request: APIService.ReadPage()
    PRIMARY KEY READS #1:
        READ: MovieInfoService.ReadMovieInfo() ... movie_info_db.movie_info.FindOne()
            - field: movie_info_db.movie_info._id
            - constraint: PRIMARY KEY (movie_info_db.movie_info._id)
        READ: MovieIdService.ReadMovieId() ... movie_id_db.movie.FindOne()
            - field: movie_id_db.movie._id
            - constraint: PRIMARY KEY (movie_id_db.movie._id)
```

- `entry request` is the client request that leads to the two reads.
- Each `READ` reads one part of the entity from a different database, using the same primary key.

### Un-1: Uniqueness, conflicting writes

```txt title="uniqueness-concurrency.txt"
entry request: APIService.RegisterMovie()
write (origin): APIService.RegisterMovie() ... MovieIdService.RegisterMovieId() ... movie_id_db.movie.InsertOne()
        - field (constrained): movie_id_db.movie.Title (UNIQUE)
    - affected write #1: MovieInfoService.WriteMovieInfo() ... movie_info_db.movie_info.InsertOne()
```

- `write (origin)` writes the unique field shown in `field (constrained)`.
- `affected write #N` is a related write to another database in the same request. Its effect remains even if the unique write is discarded.

> [!TIP]
> To review a warning, follow its call path in the code. Fix real problems, and suppress false positives as described in [Suppressing Detection Warnings](../3-usage/suppressing-warnings.md).
