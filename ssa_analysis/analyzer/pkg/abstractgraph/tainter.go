package abstractgraph

import (
	"fmt"
	"log"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/common"
	"analyzer/pkg/config"
	"analyzer/pkg/utils"
)

func MergeTaints(obj *AbstractObject, otherTaintsMap map[string][]*AbstractTaint, otherTaintsMapKeys []string, primary bool, traced bool, triggeredByTraces bool) (*TaintMapping) {
	fmt.Printf("[TAINTMAPPING] merging taints (primary=%t, traced=%t): %v\n", primary, traced, otherTaintsMap)
	var taintMapping *TaintMapping

	taintMapping = &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}

	// when it's not nil its because we want to maintain the order
	if otherTaintsMapKeys == nil {
		for key := range otherTaintsMap {
			otherTaintsMapKeys = append(otherTaintsMapKeys, key)
		}
	}

	for _, objpath := range otherTaintsMapKeys {
		otherTaints := otherTaintsMap[objpath]
		fmt.Printf("[TAINTMAPPING] checking existing taints for objpath (%s)\n", objpath)
		existingTaints := obj.taints[objpath]

		exists := func(otherTaint *AbstractTaint) (string, bool) {
			for _, existingTaint := range existingTaints {
				if existingTaint.Equals(otherTaint) {
					fmt.Printf("[TAINTMAPPING] [EXISTS] returning true...\n")
					return objpath, true
				}
				fmt.Printf("[TAINTMAPPING] checking if upper path (%s) vs (%s)\n", existingTaint.fieldpath, otherTaint.fieldpath)
				if ok, subpath := existingTaint.IsUpperPath(otherTaint); ok {
					fmt.Printf("[TAINTMAPPING] [EXISTS] returning false...\n")
					return objpath + subpath, false
				}
			}
			fmt.Printf("[TAINTMAPPING] [EXISTS] returning false...\n")
			return objpath, false
		}

		fmt.Printf("\t[TAINTMAPPING] existing taints on objpath=%s: %v\n", objpath, obj.taints[objpath])
		for _, otherTaint := range otherTaints {
			if objpath, equal := exists(otherTaint); !equal {
				// need to create new AbstractTaint to avoid just
				// storing the pointer and modifying its fields
				newTaint := NewAbstractTaint(
					otherTaint.fieldpath,
					otherTaint.dbcallID,
					otherTaint.dbOpType,
					primary, traced,
				)

				fmt.Printf("\t[TAINTMAPPING] [OBJ={%s}] adding new taint (%s, traced=%t) on obj path (%s): %v\n", obj.String(), common.OperationTypeToString(newTaint.dbOpType), newTaint.traced, objpath, newTaint)
				obj.taints[objpath] = append(obj.taints[objpath], newTaint)

				if !primary {
					// 1. explore all upper paths
					//
					// trace info for arguments and (especially) returns
					// can still lead to secondary taints that we still want to track
					fmt.Printf("\t[TAINTMAPPING] [OBJ={%s}] adding mapping for objpath={%s} // taint={%s} // traced={%t}\n", obj.String(), objpath, newTaint.LongString(), traced)
					var ok = true
					var subpath = ""
					for ok {
						for _, existingTaint := range obj.taints[objpath] {
							// filter by writes to reduce number of foreign keys for now
							if existingTaint.IsPrimary() && !traced {
								if subpath == "" {
									taintMapping.AddIfNotExists(*existingTaint, *newTaint, true)
									fmt.Printf("\t\t[TAINTMAPPING] [OBJ={%s}] [1] upperpath={%s} // subpath={%s} // existingTaint={%s} // traced={%t}\n", obj.String(), objpath, subpath, existingTaint.LongString(), traced)
								} else {
									lowerTaint := existingTaint.Copy()
									// TODO: verify if this is needed (can't recall now)
									lowerTaint.AddSuffixToDatabasePath(subpath)
									taintMapping.AddIfNotExists(*lowerTaint, *newTaint, true)
								}
							} else if traced && !triggeredByTraces { // 2nd cond works on digota?
								// NOTE:
								// if we allow when checkTraces = false (i.e., when curr method is called by taintTracedObjectsHelper)
								// then, it is creating a new FK such as e.g., Digota:
								// > orders_db.orders.Items[*] --> skus_db.skus.ID
								// note that the correct original FK is
								// > orders_db.orders.Items[*].ParentID --> skus_db.skus.ID
								fmt.Printf("\t\t[TAINTMAPPING] [OBJ={%s}] [3] upperpath={%s} // subpath={%s} // existingTaint={%s} // traced={%t}\n", obj.String(), objpath, subpath, existingTaint.LongString(), traced)
								if subpath == "" {
									taintMapping.AddIfNotExists(*newTaint, *existingTaint, true)
								} else {
									lowerTaint := *existingTaint
									lowerTaint.fieldpath = lowerTaint.fieldpath + subpath
									// [TO BE IMPROVED]
									// for some reason it works better when we change the
									// position between newTaint and lowerTaint in call args
									// e.g., SockShop3: order_db.orders.ID REFERENCES ship_db.shipments.ID
									//
									// i think this is because of the order
									// when tainting primary vs. traced
									taintMapping.AddIfNotExists(*newTaint, lowerTaint, true)
								}
							}
						}
						objpath, subpath, ok = utils.ExtractUpperPath(objpath)
					}
				}
			}
			// WARNING:
			// THIS CAN ONLY BE APPLIED WHEN DOING
			// ARG-PARAM or RET-RET MATCHING
			// AND NOT WHEN MERGE TAINTS IS CALLED BY taintTracedObjectsHelper()
			// this is because the latter func has the following condition such as
			// paramTrace.GetServicePath() == trace.GetServicePath()
			// so we expect dont want to mess up with the leftobjpath and rightobjpath
			// because they are not related in this case ????
			// (I ACTUALLY THINK THIS COND IS NOT NEEDED... needs some testing to figure it out)
			//
			// 2. explore all lower paths
			// taint: left <-- right
			if !triggeredByTraces { // cond works on trainticket?
				rightObjpath := objpath
				rightTaint := otherTaint
				for leftObjpath, _ := range obj.GetAllTraces() {
					fmt.Printf("\t rightObjpath: %s\n", rightObjpath)
					fmt.Printf("\t leftObjpath: %s\n", leftObjpath)
					if ok, diff := utils.IsUpperPath(rightObjpath, leftObjpath); ok {
						fmt.Printf("DIFF = %s\n", diff)

						dbpath := rightTaint.fieldpath + diff
						newTaint := NewAbstractTaint(dbpath, rightTaint.dbcallID, rightTaint.dbOpType, primary, traced)

						obj.AddTaintIfNotExists(leftObjpath, *newTaint)
						//taintMapping.AddIfNotExists(*rightTaint, *newTaint, true)
					}
				}
			}
		}
	}
	return taintMapping
}

func MergeTraces(obj *AbstractObject, otherTracesMap map[string][]*AbstractTrace) {
	for otherKey, otherTracesLst := range otherTracesMap {
		obj.traces[otherKey] = append(obj.traces[otherKey], otherTracesLst...)
	}
}

// updateTransitiveReferencesTriggeredByCurrent creates a new transitive reference according to the rule above,
// where (b) is the **current** constraint received as parameter, and (a) is an **old** constraint
// that we want to upgrade to (c)
//
// RULE:
// (a) X references Y (OLD)
// (b) Y references Z (CURRENT)
//
// if (a) and (b), then (c)
// (c) X references Z (NEW)
func updateTransitiveReferencesTriggeredByCurrent(graph *AbstractCallGraph, current *backends.Constraint) {
	if !config.Global.EnableTransitiveReferences {
		return
	}

	for _, db := range graph.app.GetAllDatabases() {
		for _, schema := range db.GetAllSchemas() {
			var toDelete []*backends.Constraint
			var toAdd []*backends.Constraint
			for _, old := range schema.GetAllForeignKeyConstraints() {
				if old.GetFieldAt(1) == current.GetFieldAt(0) {
					new := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, old.GetFieldAt(0), current.GetFieldAt(1))
					new.SetTransitive()
					new.CopyMandatory(current)

					toDelete = append(toDelete, old)
					toAdd = append(toAdd, new)
				}
			}
			for _, constraint := range toDelete {
				schema.RemoveConstraint(constraint)
			}

			for _, constraint := range toAdd {
				schema.AddConstraint(constraint)
			}
		}
	}
}

// RULE:
// (a) X references Y (OLD)
// (b) Y references Z (CURRENT)
//
// if (a) and (b), then (c)
// (c) X references Z (NEW)
func createTransitiveReferenceIfExists(old *backends.Constraint) bool {
	if !config.Global.EnableTransitiveReferences {
		return false
	}

	var transitiveConstraints []*backends.Constraint
	field1 := old.GetFieldAt(0)
	field2 := old.GetFieldAt(1)

	for _, current := range field2.GetConstraintForeignKey() {
		field3 := current.GetFieldAt(1)
		if field2 != field3 {
			if c := field3.GetConstraintForeignKeyToField(field1); c != nil {
				// skip since it already exists the other way around
				continue
			} else {
				new := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field1, field3)
				new.SetTransitive()
				new.CopyMandatory(current)
				field1.GetSchema().AddConstraint(new)

				transitiveConstraints = append(transitiveConstraints, new)
			}
		}
	}
	return len(transitiveConstraints) > 0
}

func PropagateNewTaintsToDatabaseSchemas(graph *AbstractCallGraph, reqIdx int, taintMapping *TaintMapping) bool {
	var modified bool
	mappingKeys := taintMapping.GetMappingKeys()

	for _, currTaint := range mappingKeys {
		otherTaintsLst := taintMapping.GetMappingForKey(currTaint)
		currDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(currTaint.GetDatabasePath()))
		currField := currDb.GetLastSchema().GetOrCreateField(currDb, currTaint.GetDatabasePath())

		for _, otherTaint := range otherTaintsLst {
			otherDb := graph.GetApp().GetDatabaseByName(utils.ExtractDatabaseNameFromFieldPath(otherTaint.GetDatabasePath()))

			if currDb == otherDb {
				// skip if its the same
				// may happen when iterating queue.Push() --> queue.Pop()
				continue
			}
			otherField := otherDb.GetLastSchema().GetOrCreateField(otherDb, otherTaint.GetDatabasePath())

			if otherField.GetDatabase() == currField.GetDatabase() {
				continue
			}

			if otherTaint.IsWriteOrUpdate() && currTaint.IsWriteOrUpdate() {
				if propagateTaintsWriteWritePair(graph, reqIdx, currTaint, otherTaint, currDb, otherDb, currField, otherField) {
					modified = true
				}
			} else if otherTaint.IsRead() && currTaint.IsWriteOrUpdate() {
				if propagateTaintsReadWritePair(graph, reqIdx, currTaint, otherTaint, currDb, otherDb, currField, otherField) {
					modified = true
				}
			} else if otherTaint.IsWriteOrUpdate() && currTaint.IsRead() {
				if propagateTaintsWriteReadPair(graph, reqIdx, currTaint, otherTaint, currDb, otherDb, currField, otherField) {
					modified = true
				}
			} else if otherTaint.IsRead() && currTaint.IsRead() {
				if propagateTaintsReadReadPair(graph, reqIdx, currTaint, otherTaint, currDb, otherDb, currField, otherField) {
					modified = true
				}
			} else if otherTaint.IsDelete() && (currTaint.IsRead() || currTaint.IsWrite() || currTaint.IsDelete()) {
				// nothing to do
			} else if (otherTaint.IsRead() || otherTaint.IsWrite() || otherTaint.IsDelete()) && currTaint.IsDelete() {
				// nothing to do
			} else if currTaint.IsUpdate() || otherTaint.IsUpdate() {
				// nothing to do
			} else {
				log.Fatalf("\t\t[ABSTRACTGRAPH] unexpected taint mapping:\nOTHER TAINT: %s\nCURR TAINT:%s", otherTaint.LongString(), currTaint.LongString())
			}
		}
	}
	return modified
}

func propagateTaintsWriteWritePair(graph *AbstractCallGraph, reqIdx int, currTaint AbstractTaint, otherTaint AbstractTaint, currDb *backends.Database, otherDb *backends.Database, currField *backends.Field, otherField *backends.Field) bool {
	var modified bool
	if constraint := currField.GetConstraintForeignKeyToField(otherField); constraint != nil {
		if otherTaint.IsWrite() && currTaint.IsWrite() {
			if ok := constraint.EnableMandatory(reqIdx); ok {
				modified = true
				fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] (A) enabled mandatory: %s\n", constraint)
			}
		}
	} else if constraint := otherField.GetConstraintForeignKeyToField(currField); constraint != nil {
		if otherTaint.IsWrite() && currTaint.IsWrite() {
			if ok := constraint.EnableMandatory(reqIdx); ok {
				modified = true
				fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] (B) enabled mandatory: %s\n", constraint)
			}
		}
	} else if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
		// 2nd condition is for sanity check
		// may happen when iterating queue.Push() --> queue.Pop()
		constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)

		if ok := createTransitiveReferenceIfExists(constraint); ok {
			modified = true
			fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] [TRANSITIVE] added new transitive constraints\n")
		} else {
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			if otherTaint.IsWrite() && currTaint.IsWrite() {
				constraint.EnableMandatory(reqIdx)
			}
			currField.AddConstraint(constraint)
			currDb.GetLastSchema().AddConstraint(constraint)
			updateTransitiveReferencesTriggeredByCurrent(graph, constraint)
			modified = true
			fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] added new constraint: %s\n", constraint)
		}
	}
	return modified
}

func propagateTaintsReadWritePair(graph *AbstractCallGraph, reqIdx int, currTaint AbstractTaint, otherTaint AbstractTaint, currDb *backends.Database, otherDb *backends.Database, currField *backends.Field, otherField *backends.Field) bool {
	var modified bool
	if constraint := currField.GetConstraintForeignKeyToField(otherField); constraint != nil {
		if currTaint.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				fmt.Printf("\t\t[ITERATOR] [READ-WRITE] (A) disabled mandatory: %s\n", constraint)
			}
		}
	} else if constraint := otherField.GetConstraintForeignKeyToField(currField); constraint != nil {
		if currTaint.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				fmt.Printf("\t\t[ITERATOR] [READ-WRITE] (B) disabled mandatory: %s\n", constraint)
			}
		}
	} else if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
		// 2nd condition is for sanity check
		// may happen when iterating queue.Push() --> queue.Pop()
		constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)

		if ok := createTransitiveReferenceIfExists(constraint); ok {
			modified = true
			fmt.Printf("\t\t[ITERATOR] [WRITE-READ] [TRANSITIVE] added new transitive constraints\n")
		} else {
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			if currTaint.IsWrite() {
				constraint.DisableMandatory(reqIdx)
			}
			currField.AddConstraint(constraint)
			currDb.GetLastSchema().AddConstraint(constraint)
			updateTransitiveReferencesTriggeredByCurrent(graph, constraint)
			modified = true
			fmt.Printf("\t\t[ITERATOR] [WRITE-READ] added new constraint: %s\n", constraint)
		}

	}
	return modified
}

func propagateTaintsWriteReadPair(graph *AbstractCallGraph, reqIdx int, currTaint AbstractTaint, otherTaint AbstractTaint, currDb *backends.Database, otherDb *backends.Database, currField *backends.Field, otherField *backends.Field) bool {
	if otherField.GetPath() == "order_db.order.FromStation" && currField.GetPath() == "station_db.station.Name" {
		fmt.Printf("CURRENT TAINT: %s\n", currTaint.LongString())
		fmt.Printf("OTHER TAINT: %s\n", otherTaint.LongString())
		log.Fatalf("NOTE: THIS IS WHY WE NEED A SECOND SCHEMA BUILDER ITERATION!")
	}

	var modified bool
	currReadField := currField
	otherWriteField := otherField
	// e.g.,
	// postnotification: TODO
	// 		=> FOREIGN_KEY notifications_queue.notification.PostID REFERENCES posts_db.post.PostID [MANDATORY]
	// 		=> the constraint already exists so condition ahead is skipped
	//
	// sockshop3: shippingservice.ship_db.write(shipping)* // shippingservice.ship_queue.push() --> queuemaster.ship_queue.pop()*
	// 		=> FOREIGN_KEY ship_db.shipments REFERENCES ship_queue.notification
	//
	// digota: orderservice.orders_db.write(order)* <-- orderservice.skuservice.get(ctx, item.parent) // skuservice.skus_db.read(parent)*
	// 		=> FOREIGN_KEY orders_db.orders.Items[*].Parent REFERENCES skus_db.skus.Id

	// NOTE: verify this
	// not sure if we shoud leave the following conditions ahead with "nothing to do"
	// to also capture foreign keys for other combinations of operations
	if constraint := currReadField.GetConstraintForeignKeyToField(otherWriteField); constraint != nil {
		if otherTaint.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				fmt.Printf("\t\t[ITERATOR] [WRITE-READ] (A) disabled mandatory: %s\n", constraint)
			}
		}
	} else if constraint := otherWriteField.GetConstraintForeignKeyToField(currReadField); constraint != nil {
		if otherTaint.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				fmt.Printf("\t\t[ITERATOR] [WRITE-READ] (B) disabled mandatory: %s\n", constraint)
			}
		}
	} else if false && !otherWriteField.HasConstraintForeignKeyToField(currReadField) && !currReadField.HasConstraintForeignKeyToField(otherWriteField) {
		// VERSION 2
		// WRITE .. READ
		// field_write <--FK-- field_read
		// 2nd condition is for sanity check
		constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currReadField, otherWriteField)

		// FOREIGN_KEY user_db.user.UserID REFERENCES order_db.order.AccountID [MANDATORY]
		if currReadField.GetPath() == "user_db.user.UserID" && otherWriteField.GetPath() == "order_db.order.AccountID" {
			fmt.Printf("otherTaint: %s\n", otherTaint.LongString())
			fmt.Printf("currTaint: %s\n", currTaint.LongString())
			//log.Fatalf("[1] HERE!")
		}
		if currReadField.GetPath() == "order_db.order.AccountID" && otherWriteField.GetPath() == "user_db.user.UserID" {
			fmt.Printf("otherTaint: %s\n", otherTaint.LongString())
			fmt.Printf("currTaint: %s\n", currTaint.LongString())
			//log.Fatalf("[2] HERE!")
		}

		if ok := createTransitiveReferenceIfExists(constraint); ok {
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-WRITE] [TRANSITIVE] added new transitive constraints\n")
		} else {
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			if otherTaint.IsWrite() {
				constraint.DisableMandatory(reqIdx)
			}
			currReadField.AddConstraint(constraint)
			currDb.GetLastSchema().AddConstraint(constraint)
			updateTransitiveReferencesTriggeredByCurrent(graph, constraint)
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
		}
	} else if true && !currReadField.HasConstraintForeignKeyToField(otherWriteField) && !otherWriteField.HasConstraintForeignKeyToField(currReadField) {
		// WRITE .. READ
		// field_write --FK--> field_read

		// 2nd condition is for sanity check
		constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherWriteField, currReadField)

		if currReadField.GetPath() == "hometimeline_cache.*.Value[*].PostID" && otherWriteField.GetPath() == "post_db.post.PostID" {
			fmt.Printf("otherTaint: %s\n", otherTaint.LongString())
			fmt.Printf("currTaint: %s\n", currTaint.LongString())
		}
		if currReadField.GetPath() == "usertimeline_cache.*.Value[*].PostID" && otherWriteField.GetPath() == "post_db.post.PostID" {
			fmt.Printf("otherTaint: %s\n", otherTaint.LongString())
			fmt.Printf("currTaint: %s\n", currTaint.LongString())
		}

		if ok := createTransitiveReferenceIfExists(constraint); ok {
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-WRITE] [TRANSITIVE] added new transitive constraints\n")
		} else {
			otherWriteField.AddConstraint(constraint)
			otherDb.GetLastSchema().AddConstraint(constraint)
			if otherTaint.IsWrite() {
				constraint.DisableMandatory(reqIdx)
			}
			updateTransitiveReferencesTriggeredByCurrent(graph, constraint)
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
		}
	}
	return modified
}

func propagateTaintsReadReadPair(graph *AbstractCallGraph, reqIdx int, currTaint AbstractTaint, otherTaint AbstractTaint, currDb *backends.Database, otherDb *backends.Database, currField *backends.Field, otherField *backends.Field) bool {
	if !config.Global.CreateReferencesFromReadReadPair {
		return false
	}

	var modified bool
	if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
		constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherField, currField)

		if ok := createTransitiveReferenceIfExists(constraint); ok {
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-READ] [TRANSITIVE] added new transitive constraints\n")
		} else {
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			constraint.DisableMandatory(reqIdx)
			otherField.AddConstraint(constraint)
			otherDb.GetLastSchema().AddConstraint(constraint)
			updateTransitiveReferencesTriggeredByCurrent(graph, constraint)
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-READ] added new constraint: %s\n", constraint)
		}

	}
	return modified
}

func PropagateNewTaintsToDatabaseCallObjects(graph *AbstractCallGraph, node *AbstractNode, taintMapping *TaintMapping) {
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetEdgeType() == EDGE_DATABASE_CALL {
			for _, obj := range edge.GetArguments() {
				for _, currTaint := range taintMapping.GetMappingKeys() {
					otherTaintsLst := taintMapping.GetMappingForKey(currTaint)
					fmt.Printf("[PROPAGATE DB OBJECTS] taint mapping: %s\n", taintMapping)
					objpath, found := obj.FindObjectPathWithEqualOrUpperTaint(currTaint)
					for _, otherTaint := range otherTaintsLst {
						if found {
							obj.AddTaintIfNotExists(objpath, otherTaint)
						}
					}
				}
			}
		}
	}
}

// PropagateTaintsToServiceCallObjects propagates taints to traced objects within current service
//
// Reminders:
// - if current edge != nil, then the current node is acting as a callee for the current edge
func PropagateTaintsToServiceCallObjects(graph *AbstractCallGraph, currNode *AbstractNode, taintMapping *TaintMapping, currEdge *AbstractEdge, propagateFromNode bool) {
	if propagateFromNode {
		for _, otherEdge := range graph.GetEdgesFromNode(currNode) {
			// propagate from params in current node to call arguments in other edge
			for _, param := range currNode.GetParams() {
				fmt.Printf("[TRACE] [FROM_NODE] [PARAM] [NODE=%s] param={%s} // otherEdge={%s}\n", currNode.String(), param.String(), otherEdge.String())
				taintTracedObjectsOnEdge(param, currNode, otherEdge, taintMapping, true)
			}
		}
	} else {
		// propagate from call arguments (1) or returns (2) in current edge to objects acting as:
		// (a) parameters or returns in the current node
		// (b) arguments in other edges

		// 1. propagate from call arguments
		for _, arg := range currEdge.GetArguments() {
			// 1a. to objects acting as parameters or returns in the current node
			taintTracedObjectsOnNode(arg, currNode, nil, taintMapping, true)
		}
		// 2. propagate from call returns
		for _, ret := range currEdge.GetReturns() {
			// 2a. to objects acting as parameters or returns in the current node
			taintTracedObjectsOnNode(ret, currNode, nil, taintMapping, true)
		}

		var doTaintAfter bool
		for _, otherEdge := range graph.GetEdgesFromNode(currNode) {
			if otherEdge == currEdge {
				// ignore current edge when propagating from call arguments or returns
				doTaintAfter = true
				continue
			}
			// 1. propagate from call arguments
			for _, arg := range currEdge.GetArguments() {
				fmt.Printf("[TRACE] [FROM_EDGE] [ARG] [NODE=%s] arg={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", currNode.String(), arg.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				// 1b. to objects acting as arguments in other edges
				taintTracedObjectsOnEdge(arg, currNode, otherEdge, taintMapping, doTaintAfter)
			}
			// 2. propagate from call returns
			for _, ret := range currEdge.GetReturns() {
				fmt.Printf("[TRACE] [FROM_EDGE] [RET] [NODE=%s] ret={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", currNode.String(), ret.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				// 2b. to objects acting as arguments in other edges
				taintTracedObjectsOnEdge(ret, currNode, otherEdge, taintMapping, doTaintAfter)
			}
		}
	}
}

// taintTracedObjectsOnEdge checks traces on objects used for other calls (aka edge)
func taintTracedObjectsOnEdge(obj *AbstractObject, currNode *AbstractNode, otherEdge *AbstractEdge, taintMapping *TaintMapping, doTaintAfter bool) {
	for objpath, tracesLst := range obj.GetTraces() {
		// e.g.,
		// MediaMicroservices in APIService.ReadPage(...)
		//
		// movieId := movieIdService.ReadMovieId(title)
		// movieInfo := movieInfoService.ReadMovieInfo(movieId.ID)
		//
		// t4 = ReadMovieId(..) => objpath 		 (@ t4.MovieID) = _obj.MovieID
		// ReadMovieInfo(t10) 	=> tracedObjPath (@ t10) 		= _obj
		//
		// REMINDER: traceObjPath is simply the objpath of the traced object
		for _, trace := range tracesLst {
			if trace.GetServiceCallID() != otherEdge.GetID() {
				continue
			}
			// e.g., SockShop3 @ Frontend.AddItem
			// AddItem(ctx, sessionID, Item{ID: itemID, Quantity: 1, UnitPrice: sock.Price})
			// <=> AddItem(ctx, sessionID, t18)
			// ------------------------------------
			// 		t12: local Item (complit)
			// ------------------------------------
			// 		t13: &t12.ID (#0)
			// ==== tainted ====
			// 		_obj
			// [rpc] @ CatalogueService.Get.itemID
			// [rpc] @ CatalogueService.AddItem.t18.ID
			// 	  	   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
			// ------------------------------------
			// 		t18: *t12
			// 		^^^^^^^^^
			// ==== tainted ====
			//		 _obj
			// [rpc] @ CartService.AddItem.t18
			// 		_obj.ID
			// 		^^^^^^^^
			// [rpc] @ CartService.Get.itemID
			// [rpc] @ CartService.AddItem.t18.ID
			// ------------------------------------
			//
			// CURRENT OBJECT is t13 w/ objpath = _obj
			// TRACED OBJECT is t18 w/ traceobjpath = _obj.ID
			//
			// we want to get the taints of t13 at _obj and propagate them to t18 on _obj.ID
			// REMINDER: we just simply associate the SAME dbfield to t18 on _obj.ID
			if tracedObj := otherEdge.GetArgumentByNameIfExists(trace.GetArgumentName()); tracedObj != nil {
				fmt.Printf("[TRACE] [OBJ=%s // OBJPATH=%s] trace={%s}\n", obj.String(), objpath, trace.LongString())
				// REMINDER: traceObjPath is simply the objpath of the traced object
				tracedObjPaths := trace.GetArgumentPath()
				taintTracedObjectsHelper(objpath, tracedObjPaths, obj, tracedObj, trace, taintMapping, true, doTaintAfter)
			}
		}
	}
}

// taintTracedObjectsOnNode checks traces on objects used as parameters or returns in the current function (aka node)
func taintTracedObjectsOnNode(obj *AbstractObject, currNode *AbstractNode, otherEdge *AbstractEdge, taintMapping *TaintMapping, doTaintAfter bool) {
	for objpath, tracesLst := range obj.GetTraces() {
		for _, trace := range tracesLst {
			var tracedObjPaths []string
			var tracedObjs []*AbstractObject
			
			fmt.Printf("[TRACE] [OBJ=%s // OBJPATH=%s] trace={%s}\n", obj.String(), objpath, trace.LongString())

			for _, param := range currNode.GetParams() {
				for paramObjpath, paramTraceLst := range param.GetTraces() {
					for _, paramTrace := range paramTraceLst {
						if paramTrace.GetServiceCallID() == trace.GetServiceCallID() {
							if paramTrace.GetServicePath() == trace.GetServicePath() {
								tracedObjs = append(tracedObjs, param)
								tracedObjPaths = append(tracedObjPaths, paramObjpath)
							} else {
								//TODO?
								fmt.Printf("param trace call ID: %s\n", paramTrace.GetServiceCallID())
								fmt.Printf("param trace path: %s\n", paramTrace.GetServicePath())
								fmt.Printf("trace path: %s\n", trace.GetServicePath())
							}
						}
					}
				}
			}
			for _, ret := range currNode.GetReturns() {
				for retObjpath, retTraceLst := range ret.GetTraces() {
					for _, retTrace := range retTraceLst {
						if retTrace.GetServiceCallID() == trace.GetServiceCallID() {
							if retTrace.GetServicePath() == trace.GetServicePath() {
								tracedObjs = append(tracedObjs, ret)
								tracedObjPaths = append(tracedObjPaths, retObjpath)
								fmt.Printf("ret trace call ID: %s\n", retTrace.GetServiceCallID())
								fmt.Printf("ret trace path: %s\n", retTrace.GetServicePath())
								fmt.Printf("trace path: %s\n", trace.GetServicePath())
							} else {
								//TODO?
								fmt.Printf("ret trace call ID: %s\n", retTrace.GetServiceCallID())
								fmt.Printf("ret trace path: %s\n", retTrace.GetServicePath())
								fmt.Printf("trace path: %s\n", trace.GetServicePath())
							}
						}
					}
				}
			}

			for i, tracedObj := range tracedObjs {
				// REMINDER: traceObjPath is simply the objpath of the traced object
				taintTracedObjectsHelper(objpath, tracedObjPaths[i], obj, tracedObj, trace, taintMapping, false, doTaintAfter)

			}
		}
	}
}

func taintTracedObjectsHelper(objpath string, tracedObjPath string, obj *AbstractObject, tracedObj *AbstractObject, trace *AbstractTrace, taintMapping *TaintMapping, onEdge bool, after bool) {
	fmt.Printf("[TRACE] [OBJ=%s // OBJPATH=%s] corresponding trace obj (path=%s): %s\n", obj.String(), objpath, tracedObjPath, tracedObj.String())
	var selectedTaints = make(map[string][]*AbstractTaint)
	var selectedTaintsKeys []string

	var ok = true
	var subpath = ""
	// if there is no taint for current objpath then it is possible that there are upper taints
	// so we go up, create a new subtaint and save to the selected taints of the traced object
	// e.g., MediaMicroservices: APIService.ReadPage():
	//
	// 				CURRENT OBJECT BELOW
	// ------------------------------------------
	// ==== arg 1 (movieID) tainted ====
	// 			_obj
	// [read, secondary] @ movieid_db.movieid
	// 			_obj.ID
	// [rpc] @ MovieIdService.ReadMovieId.movieID
	// ------------------------------------------
	// after going up, we get a new potential subtaint
	// (that we don't save for the current obj but only for the traced obj)
	// ------------------------------------------
	// 			_obj.ID
	// [read, secondary] @ movieid_db.movieid.ID
	// ------------------------------------------
	for ok {
		for _, taint := range obj.GetTaintsForObjectPath(objpath) {
			subTaint := taint.Copy()
			subTaint.AddSuffixToDatabasePath(subpath)
			selectedTaints[tracedObjPath] = append(selectedTaints[tracedObjPath], subTaint)
			selectedTaintsKeys = append(selectedTaintsKeys, tracedObjPath)
		}
		if onEdge {
			objpath, subpath, ok = utils.ExtractUpperPath(objpath)
		} else {
			// removing this is dangerous
			ok = false
		}
	}

	taintMappingTmp := MergeTaints(tracedObj, selectedTaints, selectedTaintsKeys, false, true, true)
	fmt.Printf("[TRACE] taint mapping tmp = %s\n", taintMappingTmp.String())

	if taintMapping != nil {
		taintMapping.Merge(taintMappingTmp, after)
	}

}
