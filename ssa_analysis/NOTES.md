# MS CONSISTENCY ANALYZER (SSA PARSER)

## TODOS (MINOR)
1. (nosql): parse bson projections (see dsb_sn2 on UserMentionService.ComposeUserMentions)
2. (cache): fix fields (only needs KEY and VALUE... nothing more ??)
3. (sql):   parse mysql queries that extract multiple filters (instead of using *)
4. (nosql): parse NoSQL constraint file
5. (nosql): parse advanced filters for mongodb UPDATES
6. (nosql): get tags from struct fields and apply to dbfield taints (e.g., postid insteand of PostID)
7. (detection): improve precision about affected fields/constraints

## TODOS (MAJOR)
1. implement inter-procedural analysis
2. improve ssagraph taint algorithm
3. reduce number of foreign keys (and thus absence of cascade delestes) by filtering by unique dbfields
4. implement transitivity for inference of foreign keys with writes in diff requests (if A references B and B references C, then A references C)

## ASSUMPTIONS
- inline functions or go routines are ignored by parser
- for mongodb, the name in blueprint wiring must match the database name used in `GetCollection()`
- the service structure name that implements exposed methods must be `<service_interface_name>Impl`
- the service constructor must return the service interface and not the service struct
