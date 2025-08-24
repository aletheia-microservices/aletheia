# MS CONSISTENCY ANALYZER (SSA PARSER)

## TODOS (MINOR)
1. (nosql): parse bson projections (see dsb_sn2 on UserMentionService.ComposeUserMentions)
2. (cache): fix fields (only needs KEY and VALUE... nothing more ??)
3. (sql):   parse mysql queries that extract multiple filters (instead of using *)
4. (nosql): parse NoSQL constraint file
5. (nosql): parse advanced filters for mongodb UPDATES
6. (nosql): get tags from struct fields and apply to dbfield taints (e.g., postid insteand of PostID)

## TODOS (MAJOR)
1. implement inter-procedural analysis
2. improve ssagraph taint algorithm

## ASSUMPTIONS
- inline functions or go routines are ignored by parser
- for mongodb, the name in blueprint wiring must match the database name used in `GetCollection()`
- the service structure name that implements exposed methods must be `<service_interface_name>Impl`
- the service constructor must return the service interface and not the service struct


Block #26: socialnetwork2.UserTimelineService.ReadUserTimeline.if.done
			00: t122 = &t107.PostID [#0]
			01: t123 = *t122
			02: t124 = new [1]int64 (varargs)
			03: t125 = &t124[0:int]
			04: *t125 = t123
			05: t126 = slice t124[:]
			06: t127 = append(t111, t126...)
			07: jump 24
