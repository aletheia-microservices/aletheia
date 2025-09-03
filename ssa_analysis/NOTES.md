# MS CONSISTENCY ANALYZER (SSA PARSER)

## 1. TODOS

### 1.1. MINOR FEATURES
1. (nosql): parse bson projections (see dsb_sn2 on UserMentionService.ComposeUserMentions)
2. (cache): fix fields (only needs KEY and VALUE... nothing more ??)
3. (sql):   parse mysql queries that extract multiple filters (instead of using *)
4. (nosql): parse NoSQL constraint file
5. (nosql): parse advanced filters for mongodb UPDATES
6. (nosql): get tags from struct fields and apply to dbfield taints (e.g., postid insteand of PostID)

### 1.2. MAJOR FEATURES
1. (analysis) implement inter-procedural analysis
2. (ssagraph) improve ssagraph taint algorithm
3. (schema) reduce number of foreign keys (and thus absence of cascade delestes) by filtering by unique dbfields
4. (schema) implement transitivity for inference of foreign keys with writes in diff requests (if A references B and B references C, then A references C)

### 1.3. MINOR FIXES
1. (detection): fix existence of duplicated foreign key reads -- temporary fix now concerns adding an additional condition that verifies if op calls are equal in the KeyCoordinationDetector's `hasForeignRead` method 
2. (detection): improve precision about affected fields/constraints in results

### 1.4. CRITICAL FIXES
1. (iterator) figure out why trainticket needs 2 passes on first phase to build schema in order to get following foreign keys:
```sql
"FOREIGN_KEY order_db.order.FromStation REFERENCES station_db.station.Name",
"FOREIGN_KEY order_db.order.ToStation REFERENCES station_db.station.Name",
```
2. (tainter) fix propagation to distinguish between keys used for reads and values retrived from reads (check socialnetwork on `HomeTimelineService.WriteHomeTimeline` where tainter is assuming that the `postID` method parameter was used to read the cache in `h.homeTimelineCache.Get(ctx, id_str, &posts)` due to the append afterwards `posts = append(posts, PostInfo{PostID: postID, Timestamp: timestamp})` that includes the `postID`)

## 2. CURRENT ASSUMPTIONS
- inline functions or go routines are ignored by parser
- for mongodb, the name in blueprint wiring must match the database name used in `GetCollection()`
- the service structure name that implements exposed methods must be `<service_interface_name>Impl`
- the service constructor must return the service interface and not the service struct


	[FOREIGN KEY CONCURRENCY | CHECKER] delete = ContactsService.Delete() --> contacts_db.contacts.DeleteOne()
