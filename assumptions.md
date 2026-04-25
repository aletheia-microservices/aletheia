# Aletheia: Automated Detection of Data Integrity Violations in Microservices

## Analysis Assumptions

This document describes the assumptions made by Aletheia’s current implementation that developers should take into account when analyzing their applications. These are temporary assumptions that can be addressed in the future by extending the implementation:

### 1. Database Naming

- The name of the database used in the wiring specification (e.g., `mongodb.Container`) must exactly match the name passed in service operations (e.g., `GetCollection`).

  ```Go
  // wiring
  posts_db := mongodb.Container(spec, "posts_db")
  // workflow
  s.postsDb.GetCollection(ctx, "posts_db", ...)
  ```

### 2. Service Structure

- The name of the service struct should be named `<service-interface>Impl`.

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

### 3. Schema Inference Constraints

- In NoSQL databases, bson tags/filters for writes/reads must match the corresponding Go struct field name. Note that the `_id` field is treated as a special case and is used to infer primary key constraints.
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
