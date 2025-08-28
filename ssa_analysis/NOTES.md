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

### 1.3. CRITICAL FIXES
1. (detection): add lineage info to remove false-positive on ref. integrity uncoord. repl. when reading before/after inserting the record (check trainticket)
2. (detection): fix false-positive on ref. integrity uncoord. repl. (check trainticket) by restricting to MANDATORY constraints on CancelOrder w/ READ_1=UserService.FindByUserID and READ_2=OrderService.Find
3. (detection): fix false-positive on ref. integrity uncoord. repl. (check dsb_socialnetwork) when there's a read to cache before read/write to database
4. (detection): improve precision about affected fields/constraints in results

## 2. CURRENT ASSUMPTIONS
- inline functions or go routines are ignored by parser
- for mongodb, the name in blueprint wiring must match the database name used in `GetCollection()`
- the service structure name that implements exposed methods must be `<service_interface_name>Impl`
- the service constructor must return the service interface and not the service struct
