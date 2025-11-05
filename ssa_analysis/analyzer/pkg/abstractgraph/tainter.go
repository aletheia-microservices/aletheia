package abstractgraph

import (
	"fmt"
	"log"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/utils"
)

func PropagateNewTaintsToDatabaseSchemas(graph *AbstractCallGraph, reqIdx int, taintMapping *TaintMapping) bool {
	var modified bool
	for currTaint, otherTaintsLst := range taintMapping.mapping {
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
					if otherTaint.IsWrite() && currTaint.IsWrite() {
						constraint.EnableMandatory(reqIdx)
					}
					currField.AddConstraint(constraint)
					currDb.GetLastSchema().AddConstraint(constraint)
					modified = true
					fmt.Printf("\t\t[ITERATOR] [WRITE-WRITE] added new constraint: %s\n", constraint)
				}
			} else if otherTaint.IsRead() && currTaint.IsWriteOrUpdate() {
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
					currField.AddConstraint(constraint)
					currDb.GetLastSchema().AddConstraint(constraint)
					if currTaint.IsWrite() {
						constraint.DisableMandatory(reqIdx)
					}
					modified = true
					fmt.Printf("\t\t[ITERATOR] [WRITE-READ] added new constraint: %s\n", constraint)
				}
			} else if otherTaint.IsWriteOrUpdate() && currTaint.IsRead() {
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
				if constraint := currField.GetConstraintForeignKeyToField(otherField); constraint != nil {
					if otherTaint.IsWrite() {
						if ok := constraint.DisableMandatory(reqIdx); ok {
							modified = true
							fmt.Printf("\t\t[ITERATOR] [WRITE-READ] (A) disabled mandatory: %s\n", constraint)
						}
					}
				} else if constraint := otherField.GetConstraintForeignKeyToField(currField); constraint != nil {
					if otherTaint.IsWrite() {
						if ok := constraint.DisableMandatory(reqIdx); ok {
							modified = true
							fmt.Printf("\t\t[ITERATOR] [WRITE-READ] (B) disabled mandatory: %s\n", constraint)
						}
					}
				} else if !currField.HasConstraintForeignKeyToField(otherField) && !otherField.HasConstraintForeignKeyToField(currField) {
					// 2nd condition is for sanity check
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, otherField, currField)
					otherField.AddConstraint(constraint)
					otherDb.GetLastSchema().AddConstraint(constraint)
					if otherTaint.IsWrite() {
						constraint.DisableMandatory(reqIdx)
					}
					modified = true
					fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
				}/*  else if !otherField.HasConstraintForeignKeyToField(currField) && !currField.HasConstraintForeignKeyToField(otherField) {
					// 2nd condition is for sanity check
					constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, currField, otherField)
					currField.AddConstraint(constraint)
					currDb.GetLastSchema().AddConstraint(constraint)
					if otherTaint.IsWrite() {
						constraint.DisableMandatory(reqIdx)
					}
					modified = true
					fmt.Printf("\t\t[ITERATOR] [READ-WRITE] added new constraint: %s\n", constraint)
				} */
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

func PropagateNewTaintsToDatabaseCallObjects(graph *AbstractCallGraph, node *AbstractNode, taintMapping *TaintMapping) {
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetEdgeType() == EDGE_DATABASE_CALL {
			for _, obj := range edge.GetArguments() {
				for currTaint, otherTaintsLst := range taintMapping.GetMapping() {
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
				taintTracedObjects(param, otherEdge, taintMapping)
			}
		}
	} else {
		for _, otherEdge := range graph.GetEdgesFromNode(node) {
			if otherEdge == currEdge {
				continue
			}
			for _, arg := range currEdge.GetArguments() {
				fmt.Printf("[TRACE] [ARG] [NODE=%s] arg={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", node.String(), arg.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				taintTracedObjects(arg, otherEdge, taintMapping)
			}
			for _, ret := range currEdge.GetReturns() {
				fmt.Printf("[TRACE] [RET] [NODE=%s] ret={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", node.String(), ret.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				taintTracedObjects(ret, otherEdge, taintMapping)
			}
		}
	}
}

func taintTracedObjects(obj *AbstractObject, otherEdge *AbstractEdge, taintMapping *TaintMapping) {
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
					taintMapping.Merge(taintMappingTmp)
				}
			}

		}
	}
}
