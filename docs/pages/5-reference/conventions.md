# Conventions

Aletheia currently relies on the following naming and code conventions when analyzing Blueprint applications. Check that your application follows them before running the analysis.

## Database Naming

- The name of the database used in the wiring specification (e.g., `mongodb.Container`) must exactly match the name passed in service operations (e.g., `GetCollection`). 

  ```Go
  // wiring
  posts_db := mongodb.Container(spec, "posts_db")
  // workflow
  s.postsDb.GetCollection(ctx, "posts_db", ...)
  ```

## Service Structure

- Name the service struct `<service-interface>Impl`.

  ```Go
  type StorageService interface {...}
  type StorageServiceImpl struct {...}
  ```

- Service implementation functions must return the interface type (e.g., `StorageService`) rather than the concrete service struct (e.g., `StorageServiceImpl`).
  ```Go
  func NewStorageServiceImpl(ctx context.Context, ...) (StorageService, error) {
    return &StorageServiceImpl{...}, nil
  }
  ```

## Schema Fields

- For NoSQL databases, BSON tags and filter keys must match the corresponding Go struct field names. The exception is `_id`, which Aletheia uses to infer primary keys.
  ```Go
  type Movie struct {
  	MovieID string `bson:"_id"`
  	Title   string `bson:"Title"`
  }
  ...
  movie := Movie{...}
  collection.InsertOne(ctx, movie)
  ...
  query := bson.D{{Key: "Title", Value: title}}
  collection.FindOne(ctx, query)
  ```
