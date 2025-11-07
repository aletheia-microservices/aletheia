package abstractgraph

import (
	"fmt"
	"log"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/config"
	"analyzer/pkg/utils"
)

func isTransitiveReference(constraint *backends.Constraint) (bool, []*backends.Constraint) {
	if !config.Global.EnableTransitiveReferences {
		return false, nil
	}

	var transitiveConstraints []*backends.Constraint
	field1 := constraint.GetFieldAt(0)
	field2 := constraint.GetFieldAt(1)

	for _, constraint2 := range field2.GetConstraintForeignKey() {
		field3 := constraint2.GetFieldAt(1)
		if field2 != field3 {
			if c := field3.GetConstraintForeignKeyToField(field1); c != nil {
				// skip since it already exists the other way around
				continue
			} else {
				transitiveConstraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field1, field3)
				transitiveConstraint.SetTransitive()
				transitiveConstraints = append(transitiveConstraints, transitiveConstraint)
			}
		}
	}
	return len(transitiveConstraints) > 0, transitiveConstraints
}

func PropagateNewTaintsToDatabaseSchemas(graph *AbstractCallGraph, reqIdx int, taintMapping *TaintMapping) bool {
	var modified bool

	mappingKeys := taintMapping.GetMappingKeys()
	/* sort.Slice(mappingKeys, func(i, j int) bool {
		return mappingKeys[i].LongString() < mappingKeys[j].LongString()
	}) */

	for _, currTaint := range mappingKeys {
		otherTaintsLst := taintMapping.GetMappingForKey(currTaint)
		/* sort.Slice(otherTaintsLst, func(i, j int) bool {
			return otherTaintsLst[i].LongString() < otherTaintsLst[j].LongString()
		}) */

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
				// nothing to do
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

		if isTransitive, transitiveConstraints := isTransitiveReference(constraint); isTransitive {
			for _, constraint := range transitiveConstraints {
				// must (un)set mandatory before calling GetSchema().AddConstraint()
				if otherTaint.IsWrite() && currTaint.IsWrite() {
					constraint.EnableMandatory(reqIdx)
				}
				field1 := constraint.GetFieldAt(0)
				field1.GetSchema().AddConstraint(constraint)
				modified = true
				fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] [TRANSITIVE] added new constraint: %s\n", constraint)
			}
		} else {
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			if otherTaint.IsWrite() && currTaint.IsWrite() {
				constraint.EnableMandatory(reqIdx)
			}
			currField.AddConstraint(constraint)
			currDb.GetLastSchema().AddConstraint(constraint)
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

		if isTransitive, transitiveConstraints := isTransitiveReference(constraint); isTransitive {
			for _, constraint := range transitiveConstraints {
				// must (un)set mandatory before calling GetSchema().AddConstraint()
				if currTaint.IsWrite() {
					constraint.DisableMandatory(reqIdx)
				}
				field1 := constraint.GetFieldAt(0)
				field1.GetSchema().AddConstraint(constraint)
				modified = true
				fmt.Printf("\t\t[ITERATOR] [WRITE-READ] [TRANSITIVE] added new constraint: %s\n", constraint)
			}
		} else {
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			if currTaint.IsWrite() {
				constraint.DisableMandatory(reqIdx)
			}
			currField.AddConstraint(constraint)
			currDb.GetLastSchema().AddConstraint(constraint)
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

		if isTransitive, transitiveConstraints := isTransitiveReference(constraint); isTransitive {
			for _, constraint := range transitiveConstraints {
				// must (un)set mandatory before calling GetSchema().AddConstraint()
				if currTaint.IsWrite() {
					constraint.DisableMandatory(reqIdx)
				}
				field1 := constraint.GetFieldAt(0)
				field1.GetSchema().AddConstraint(constraint)
				modified = true
				fmt.Printf("\t\t[ITERATOR] [READ-WRITE] [TRANSITIVE] added new constraint: %s\n", constraint)
			}
		} else {
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			if otherTaint.IsWrite() {
				constraint.DisableMandatory(reqIdx)
			}
			currReadField.AddConstraint(constraint)
			currDb.GetLastSchema().AddConstraint(constraint)
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
		}
	} else if true && !currReadField.HasConstraintForeignKeyToField(otherWriteField) && !otherWriteField.HasConstraintForeignKeyToField(currReadField) {
		// WRITE .. READ
		// field_write --FK--> field_read

		// 2nd condition is for sanity check
		constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherWriteField, currReadField)

		// FOREIGN_KEY user_db.user.UserID REFERENCES order_db.order.AccountID [MANDATORY]
		// FOREIGN_KEY post_db.post.PostID REFERENCES hometimeline_cache.*.Value[*].PostID [MANDATORY]
		// FOREIGN_KEY post_db.post.PostID REFERENCES usertimeline_cache.*.Value[*].PostID [MANDATORY]
		/* if currReadField.GetPath() == "user_db.user.UserID" && otherWriteField.GetPath() == "order_db.order.AccountID" {
			log.Fatalf("[1] HERE!")
		}
		if currReadField.GetPath() == "order_db.order.AccountID" && otherWriteField.GetPath() == "user_db.user.UserID" {
			log.Fatalf("[2] HERE!")
		} */
		if currReadField.GetPath() == "hometimeline_cache.*.Value[*].PostID" && otherWriteField.GetPath() == "post_db.post.PostID" {
			fmt.Printf("otherTaint: %s\n", otherTaint.LongString())
			fmt.Printf("currTaint: %s\n", currTaint.LongString())
			//log.Fatalf("[2] HERE HOMETIMELINE!")
		}
		if currReadField.GetPath() == "usertimeline_cache.*.Value[*].PostID" && otherWriteField.GetPath() == "post_db.post.PostID" {
			fmt.Printf("otherTaint: %s\n", otherTaint.LongString())
			fmt.Printf("currTaint: %s\n", currTaint.LongString())
			//log.Fatalf("[2] HERE USERTIMELINE!")
		}

		if isTransitive, transitiveConstraints := isTransitiveReference(constraint); isTransitive {
			for _, constraint := range transitiveConstraints {
				// must (un)set mandatory before calling GetSchema().AddConstraint()
				if currTaint.IsWrite() {
					constraint.DisableMandatory(reqIdx)
				}
				field1 := constraint.GetFieldAt(0)
				field1.GetSchema().AddConstraint(constraint)
				modified = true
				fmt.Printf("\t\t[ITERATOR] [READ-WRITE] [TRANSITIVE] added new constraint: %s\n", constraint)
			}
		} else {
			otherWriteField.AddConstraint(constraint)
			otherDb.GetLastSchema().AddConstraint(constraint)
			if otherTaint.IsWrite() {
				constraint.DisableMandatory(reqIdx)
			}
			modified = true
			fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
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

// propagate taints to traced objects within current service
// (i.e. objects passed as arguments in 'otherEdge' to other service calls)
func PropagateNewTaintsToTracedObjects(graph *AbstractCallGraph, node *AbstractNode, taintMapping *TaintMapping, currEdge *AbstractEdge, propagateFromParams bool) {
	if propagateFromParams {
		for _, otherEdge := range graph.GetEdgesFromNode(node) {
			for _, param := range node.GetParams() {
				fmt.Printf("[TRACE] [PARAM] [NODE=%s] param={%s} // otherEdge={%s}\n", node.String(), param.String(), otherEdge.String())
				taintTracedObjects(param, otherEdge, taintMapping, true)
			}
		}
	} else {
		var after bool
		for _, otherEdge := range graph.GetEdgesFromNode(node) {
			if otherEdge == currEdge {
				after = true
				continue
			}
			for _, arg := range currEdge.GetArguments() {
				fmt.Printf("[TRACE] [ARG] [NODE=%s] arg={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", node.String(), arg.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				taintTracedObjects(arg, otherEdge, taintMapping, after)
			}
			for _, ret := range currEdge.GetReturns() {
				fmt.Printf("[TRACE] [RET] [NODE=%s] ret={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", node.String(), ret.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				taintTracedObjects(ret, otherEdge, taintMapping, after)
			}
		}
	}
}

func taintTracedObjects(obj *AbstractObject, otherEdge *AbstractEdge, taintMapping *TaintMapping, after bool) {
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
			fmt.Printf("[TRACE] [OBJ=%s // OBJPATH=%s] trace={%s}\n", obj.String(), objpath, trace.LongString())
			tracedObj := otherEdge.GetArgumentByNameIfExists(trace.GetArgumentName())
			if tracedObj != nil {
				tracedObjPath := trace.GetArgumentPath()
				fmt.Printf("[TRACE] [OBJ=%s // OBJPATH=%s] corresponding trace obj (path=%s): %s\n", obj.String(), objpath, tracedObjPath, tracedObj.String())
				var selectedTaints = make(map[string][]*AbstractTaint)

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
					}
					objpath, subpath, ok = utils.ExtractUpperPath(objpath)
				}

				taintMappingTmp := MergeTaints(tracedObj, selectedTaints, false, true)
				fmt.Printf("[TRACE] taint mapping tmp = %s\n", taintMappingTmp.String())

				if taintMapping != nil {
					taintMapping.Merge(taintMappingTmp, after)
				}
			}

		}
	}
}
