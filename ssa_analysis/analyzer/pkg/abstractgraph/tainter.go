package abstractgraph

import (
	"github.com/sirupsen/logrus"

	"analyzer/pkg/app/backends"
	"analyzer/pkg/config"
	"analyzer/pkg/utils"
)

type MergeMode int

// i.e., triggers for the merge
const (
	MERGE_MODE_PARSE MergeMode = iota
	MERGE_MODE_TAINT
	MERGE_MODE_TRACE
	MERGE_MODE_DEBUG
)

func mergeModeToString(mode MergeMode) string {
	switch mode {
	case MERGE_MODE_PARSE:
		return "PARSE"
	case MERGE_MODE_TAINT:
		return "TAINT"
	case MERGE_MODE_TRACE:
		return "TRACE"
	}
	return ""
}

func mergeExistingTaintsWithNewTaints(obj *AbstractObject, objpath string, subpath string, newTaint *AbstractTaint, taintMapping *TaintMapping, mode MergeMode, t string) {
	for _, existingTaint := range obj.taints[objpath] {
		// filter by writes to reduce number of foreign keys for now
		// TODO: remove IsPrimary() condition?
		if existingTaint.IsPrimary() {
			lowerTaint := existingTaint.Copy()
			lowerTaint.AddSuffixToDatabasePath(subpath)

			if mode != MERGE_MODE_DEBUG {
				taintMapping.AddIfNotExists(*lowerTaint, *newTaint, true, false)
			}
			// EVAL: logrus.Tracef("\t\t[TAINTMAPPING] [MERGE] [OBJ={%s}] [1] upperpath={%s} // subpath={%s} // existingTaint={%s} // mode={%s}\n", obj.String(), objpath, subpath, existingTaint.LongString(), mergeModeToString(mode))
		} else {
			// sometimes it is not possible that taints are primary
			// for example, when there is a service that acts as a gateway for two service
			// e.g., dsb mediamicroservices:
			// [write, traced] [t34] @ movie_info_db.movie_info.Casts[*].CastInfoID
			// [write, traced] [t55] @ cast_info_db.cast.CastInfoID
			// [rpc] [t34] @ MovieInfoService.WriteMovieInfo.t9[*].CastInfoID
			// [rpc] [t55] @ CastInfoService.WriteCastInfo.t48
			if mode == MERGE_MODE_TRACE {
				if existingTaint.GetT() == t {
					// if T values are equal, then we skip since they
					// come from the same source and will eventually be matched there
					// EVAL: logrus.Tracef("\t\t[TAINTMAPPING] [MERGE] [TRACE] skipping for existingTaint={%s} and newTaint={%s} since T values (%s) are equal\n", existingTaint.LongString(), newTaint.LongString(), t)
					continue
				}
				lowerTaint := existingTaint.Copy()
				lowerTaint.AddSuffixToDatabasePath(subpath)
				taintMapping.AddIfNotExists(*lowerTaint, *newTaint, true, false)
				// EVAL: logrus.Tracef("\t\t[TAINTMAPPING] [MERGE] [TRACE] [OBJ={%s}] [1] upperpath={%s} // subpath={%s} // existingTaint={%s} // mode={%s}\n", obj.String(), objpath, subpath, existingTaint.LongString(), mergeModeToString(mode))
			}
		}
	}
}

func MergeTaints(obj *AbstractObject, otherTaintsMap map[string][]*AbstractTaint, otherTaintsMapKeys []string, mode MergeMode, t string, readOnly bool) *TaintMapping {
	// EVAL: logrus.Tracef("[TAINTMAPPING] merging taints mode={%s}: %v\n", mergeModeToString(mode), otherTaintsMap)
	var taintMapping *TaintMapping

	taintMapping = &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
	// when it's not nil its because we want to maintain the order
	if otherTaintsMapKeys == nil {
		for key := range otherTaintsMap {
			otherTaintsMapKeys = append(otherTaintsMapKeys, key)
		}
	}

	for _, objpath := range otherTaintsMapKeys {
		// EVAL: logrus.Tracef("[TAINTMAPPING] checking existing taints for objpath (%s)\n", objpath)
		existingTaints := obj.taints[objpath]

		var taintsToAdd []*AbstractTaint

		taintExists := func(otherTaint *AbstractTaint) (string, bool) {
			for _, existingTaint := range existingTaints {
				if existingTaint.Similar(otherTaint) {
					// EVAL: logrus.Tracef("[TAINTMAPPING] [EXISTS] returning true...\n")
					return objpath, true
				}
			}
			// EVAL: logrus.Tracef("[TAINTMAPPING] [EXISTS] returning false...\n")
			return objpath, false
		}

		// EVAL: logrus.Tracef("\t[TAINTMAPPING] existing taints on objpath=%s: %v\n", objpath, obj.taints[objpath])
		for _, otherTaint := range otherTaintsMap[objpath] {
			if config.Global.DualPassSchemaBuilding && readOnly && !otherTaint.IsRead() {
				logrus.WithField("dual_pass", config.Global.DualPassSchemaBuilding).WithField("read", otherTaint.IsRead()).
					Tracef("[TAINTMAPPING] skipping read taint...")
				continue
			}

			if objpath, exists := taintExists(otherTaint); !exists {
				if mode == MERGE_MODE_PARSE {
					// parameter "t" is empty for this mode
					t = otherTaint.GetT()
				}
				// need to create new AbstractTaint to avoid just storing the pointer and modifying its fields
				newTaint := NewAbstractTaint(t, otherTaint.fieldpath, otherTaint.dbcallID, otherTaint.dbOpType, mode == MERGE_MODE_PARSE, mode == MERGE_MODE_TRACE, otherTaint.IsReadKey(), otherTaint.IsReadValue())

				if mode != MERGE_MODE_DEBUG {
					taintsToAdd = append(taintsToAdd, newTaint)
				}

				// EVAL: logrus.Tracef("\t[TAINTMAPPING] [OBJ={%s}] added new taint (%s, traced=%t) on obj path (%s): %v\n", obj.String(), common.OperationTypeToString(newTaint.dbOpType), newTaint.traced, objpath, newTaint)

				// it is not necessary to be ran for MERGE_MODE_PARSE
				if mode == MERGE_MODE_PARSE {
					continue
				}

				// EVAL: logrus.Tracef("\t[TAINTMAPPING] [OBJ={%s}] attempting to add mapping for objpath={%s} // taint={%s} // mode={%s}\n", obj.String(), objpath, newTaint.LongString(), mergeModeToString(mode))

				mergeExistingTaintsWithNewTaints(obj, objpath, "", newTaint, taintMapping, mode, t)

				if mode != MERGE_MODE_TRACE {
					// The logic below for upper paths (and lower paths) cannot be ran for MERGE_MODE_TRACE because they
					// are not exact matches such as, for example, arg-params, which is necessary for computing upper taints
					// and we already extracted the selected taints (which also includes lower paths) prior to calling MergeTaints

					// 1. explore all upper paths
					var subpath string
					var ok bool
					for {
						objpath, subpath, ok = utils.ExtractUpperPath(objpath)
						if !ok {
							break
						}
						mergeExistingTaintsWithNewTaints(obj, objpath, subpath, newTaint, taintMapping, mode, t)
					}

					// 2. explore all lower paths
					/* fromObjpath := objpath
					fromTaint := otherTaint
					for _, toLocation := range obj.GetAllAbstractLocationsWithTaints() {
						if ok, diff := utils.IsUpperPath(fromObjpath, toLocation); ok {
							newDbpath := fromTaint.GetDatabasePath() + diff
							newTaint = newTaint.Copy()
							newTaint.SetDatabasepath(newDbpath)
							if mode != MERGE_MODE_DEBUG {
								obj.AddTaintIfNotExists(toLocation, newTaint)
							}
						}
					} */
				}
			}

			// 1. The logic below does not need to be ran for MERGE_MODE_TRACE because we already extracted the
			// selected taints (which include lower paths) prior to calling MergeTaints and added them above
			// 2. It is also not necessary to be ran for MERGE_MODE_PARSE
			if mode == MERGE_MODE_PARSE || mode == MERGE_MODE_TRACE {
				continue
			}

			// 1. The goal here is not to propagate new traces, but to make sure
			// the new taints are present in all abstract locations within the object,
			// which may only be annotated by traces and not taints
			// 2. This logic is needed, for example, in TrainTicket
			// 3. No need to add to taint mapping
			fromObjpath := objpath
			fromTaint := otherTaint
			for _, toLocation := range obj.GetAllAbstractLocationsWithTraces() {
				// e.g.,
				// from path: 	_obj 	=> taint: 			my_db.MyObject
				// to path: 	_obj.ID => taint to add: 	my_db.MyObject.ID (diff = .ID)

				// it is ok if fromObjpath is always an upper path of toLocation
				// in other words, toLocation is lowerpath of fromObjpath
				if ok, diff := utils.IsUpperOrEqualPath(fromObjpath, toLocation); ok {
					newDbpath := fromTaint.GetDatabasePath() + diff
					newTaint := NewAbstractTaint(t, newDbpath, fromTaint.dbcallID, fromTaint.dbOpType, mode == MERGE_MODE_PARSE, mode == MERGE_MODE_TRACE, fromTaint.IsReadKey(), fromTaint.IsReadValue())
					if mode != MERGE_MODE_DEBUG {
						obj.AddTaintIfNotExists(toLocation, newTaint)
					}
				}
			}
		}

		for _, newTaint := range taintsToAdd {
			obj.AddTaintIfNotExists(objpath, newTaint)
		}
	}
	return taintMapping
}

func MergeTraces(obj *AbstractObject, otherTracesMap map[string][]*AbstractTrace) {
	for otherKey, otherTracesLst := range otherTracesMap {
		for _, otherTrace := range otherTracesLst {
			var exists bool
			for _, existingTrace := range obj.traces[otherKey] {
				if existingTrace.svcallID == otherTrace.svcallID && existingTrace.svpath == otherTrace.svpath {
					exists = true
					break
				}
			}
			if !exists {
				obj.traces[otherKey] = append(obj.traces[otherKey], otherTrace)
			}
		}
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
func updateTransitiveReferencesTriggeredByCurrent(graph *AbstractCallGraph, schema *backends.Schema, current *backends.Constraint) {
	if !config.Global.EnableTransitiveReferences {
		return
	}

	if !config.Global.UpdateTransitiveReferencesTriggeredByCurrent {
		return
	}

	// EVAL: logrus.Tracef("[TRANSITIVE REFS] current: %s\n", current.String())

	for otherSchema := range schema.GetAllSchemasRefBy() {
		var toDelete []*backends.Constraint
		var toAdd []*backends.Constraint
		for _, old := range otherSchema.GetAllForeignKeyConstraints() {
			if old.IsMandatory() {
				if !config.Global.UpdateTransitiveReferencesTriggeredByCurrentOnMandatory {
					continue
				}
			}
			if old.GetFieldAt(1) == current.GetFieldAt(0) {
				if current.GetFieldAt(0).GetDatabase() == old.GetFieldAt(1).GetDatabase() {
					continue
				}

				new := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, old.GetFieldAt(0), current.GetFieldAt(1))
				new.SetTransitive()
				new.CopyMandatory(current)

				if config.Global.DeleteOldOnTransitiveReferences {
					toDelete = append(toDelete, old)
					// EVAL: logrus.Tracef("\t[TRANSITIVE REFS] to delete: %s\n", old.String())
				}
				toAdd = append(toAdd, new)
				// EVAL: logrus.Tracef("\t[TRANSITIVE REFS] to add: %s\n", new.String())
			}
		}
		if config.Global.DeleteOldOnTransitiveReferences {
			for _, constraint := range toDelete {
				constraint.GetFieldAt(0).RemoveConstraint(constraint)
				constraint.GetFieldAt(0).GetSchema().RemoveConstraint(constraint)
			}
		}
	
		for _, constraint := range toAdd {
			constraint.GetFieldAt(0).AddConstraint(constraint)
			schema.AddConstraint(constraint)
		}
	}

}

// RULE:
// (a) X references Y (OLD)
// (b) Y references Z (CURRENT)
//
// if (a) and (b), then (c)
// (c) X references Z (NEW)
func createTransitiveReferenceIfExists(field1 *backends.Field, field2 *backends.Field, reqIdx int, writewrite bool) bool {
	if !config.Global.EnableTransitiveReferences {
		return false
	}

	var seen = make(map[[2]*backends.Field]bool)

	for _, current := range field2.GetConstraintsForeignKeys() {
		field3 := current.GetFieldAt(1)
		if field2 != field3 {
			if c := field3.GetSchema().GetForeignKeyForPair(field3, field1); c != nil {
				// skip since it already exists the other way around
				continue
			} else if field1.GetDatabase() == field3.GetDatabase() {
				continue
			} else {
				pair := [2]*backends.Field{field1, field3}
				if _, exists := seen[pair]; !exists {
					new := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field1, field3)
					if writewrite && current.IsMandatory() {
						new.EnableMandatory(reqIdx)
					} else {
						new.DisableMandatory(reqIdx)
					}
					new.SetTransitive()
					if ok := field1.GetSchema().AddConstraint(new); ok {
						field1.AddConstraint(new)
						// EVAL: logrus.Tracef("[TRANSITIVE] added new transitive constraint: %s\n", new.String())
					}
					seen[pair] = true
				}
			}
		}
	}
	return len(seen) > 0
}

func PropagateNewTaintsToDatabaseSchemas(graph *AbstractCallGraph, reqIdx int, taintMapping *TaintMapping, readOnly bool) bool {
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

			var field1, field2 *backends.Field
			var db1, db2 *backends.Database
			var taint1, taint2 AbstractTaint

			field1 = otherField
			taint1 = otherTaint
			db1 = otherDb

			field2 = currField
			taint2 = currTaint
			db2 = currDb

			if utils.GreaterT(otherTaint.GetT(), currTaint.GetT()) {
				field1 = currField
				taint1 = currTaint
				db1 = currDb

				field2 = otherField
				taint2 = otherTaint
				db2 = otherDb
			}

			if utils.EqualT(taint1.GetT(), taint2.GetT()) {
				logrus.WithField("taint1", taint1.LongString()).WithField("taint2", taint2.LongString()).
					Warnf("[TAINTER] found taints with equal T numbers (%s) vs (%s)\n", taint1.GetT(), taint2.GetT())
			}

			if db1 == db2 {
				continue
			}

			if !config.Global.DualPassSchemaBuilding || (config.Global.DualPassSchemaBuilding && readOnly) {
				if taint1.IsRead() && taint2.IsRead() {
					logrus.WithField("taint1", taint1.String()).WithField("taint2", taint2.String()).
						Tracef("[TAINTER] found read-read taint pair")
					if propagateTaintsReadReadPair(graph, reqIdx, taint2, taint1, db2, db1, field2, field1) {
						modified = true
					}
				}
			}
			if !config.Global.DualPassSchemaBuilding || (config.Global.DualPassSchemaBuilding && !readOnly) {
				if taint1.IsWriteOrUpdate() && taint2.IsWriteOrUpdate() {
					if propagateTaintsWriteWritePair(graph, reqIdx, taint2, taint1, db2, db1, field2, field1) {
						modified = true
					}
				} else if taint1.IsRead() && taint2.IsWriteOrUpdate() {
					if propagateTaintsReadWritePair(graph, reqIdx, taint2, taint1, db2, db1, field2, field1) {
						modified = true
					}
				} else if taint1.IsWriteOrUpdate() && taint2.IsRead() {
					if propagateTaintsWriteReadPair(graph, reqIdx, taint2, taint1, db2, db1, field2, field1) {
						modified = true
					}
				} else if taint1.IsDelete() && (taint2.IsRead() || taint2.IsWrite() || taint2.IsDelete()) {
					// nothing to do
				} else if (taint1.IsRead() || taint1.IsWrite() || taint1.IsDelete()) && taint2.IsDelete() {
					// nothing to do
				} else if taint2.IsUpdate() || taint1.IsUpdate() {
					// nothing to do
				}
			}
		}
	}
	return modified
}

func propagateTaintsWriteWritePair(graph *AbstractCallGraph, reqIdx int, taint2_write AbstractTaint, taint1_write AbstractTaint, db2_write *backends.Database, db1_write *backends.Database, field2_write *backends.Field, field1_write *backends.Field) bool {
	var modified bool
	// EVAL: logrus.Tracef("[TAINTER] [WRITE-WRITE] pair (%s: %s) -> (%s: %s)\n", taint2_write.GetT(), taint2_write.String(), taint1_write.GetT(), taint1_write.String())
	if constraint := field2_write.GetConstraintForeignKeyToField(field1_write); constraint != nil {
		if taint1_write.IsWrite() && taint2_write.IsWrite() {
			if ok := constraint.EnableMandatory(reqIdx); ok {
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [WRITE-WRITE] (A) enabled mandatory: %s\n", constraint)
			}
		}
	} else if constraint := field1_write.GetConstraintForeignKeyToField(field2_write); constraint != nil {
		if taint1_write.IsWrite() && taint2_write.IsWrite() {
			if ok := constraint.EnableMandatory(reqIdx); ok {
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [WRITE-WRITE] (B) enabled mandatory: %s\n", constraint)
			}
		}
	} else if !field2_write.HasConstraintForeignKeyToField(field1_write) && !field1_write.HasConstraintForeignKeyToField(field2_write) {
		// 2nd condition is for sanity check
		// may happen when iterating queue.Push() --> queue.Pop()

		if ok := createTransitiveReferenceIfExists(field2_write, field1_write, reqIdx, true); ok {
			modified = true
		} else {
			constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field2_write, field1_write)
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			constraint.EnableMandatory(reqIdx)
			field2_write.AddConstraint(constraint)
			schema := db2_write.GetLastSchema()
			schema.AddConstraint(constraint)
			updateTransitiveReferencesTriggeredByCurrent(graph, schema, constraint)
			modified = true
			// EVAL: logrus.Tracef("\t\t[ITERATOR] [WRITE-WRITE] added new constraint: %s\n", constraint)
		}
	}
	return modified
}

func propagateTaintsReadWritePair(graph *AbstractCallGraph, reqIdx int, taint2_write AbstractTaint, taint1_read AbstractTaint, db2_write *backends.Database, db1_read *backends.Database, field2_write *backends.Field, field1_read *backends.Field) bool {
	var modified bool
	if constraint := field2_write.GetConstraintForeignKeyToField(field1_read); constraint != nil {
		if taint2_write.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [READ-WRITE] (A) disabled mandatory: %s\n", constraint)
			}
		}
	} else if constraint := field1_read.GetConstraintForeignKeyToField(field2_write); constraint != nil {
		if taint2_write.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [READ-WRITE] (B) disabled mandatory: %s\n", constraint)
			}
		}
	} else if !field2_write.HasConstraintForeignKeyToField(field1_read) && !field1_read.HasConstraintForeignKeyToField(field2_write) {
		// 2nd condition is for sanity check
		// may happen when iterating queue.Push() --> queue.Pop()

		if ok := createTransitiveReferenceIfExists(field2_write, field1_read, reqIdx, false); ok {
			modified = true
		} else {
			constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field2_write, field1_read)
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			constraint.DisableMandatory(reqIdx)
			field2_write.AddConstraint(constraint)
			schema := db2_write.GetLastSchema()
			schema.AddConstraint(constraint)
			updateTransitiveReferencesTriggeredByCurrent(graph, schema, constraint)
			modified = true
			// EVAL: logrus.Tracef("\t\t[ITERATOR] [WRITE-READ] added new constraint: %s\n", constraint)
		}

	}
	return modified
}

func propagateTaintsWriteReadPair(graph *AbstractCallGraph, reqIdx int, taint2_read AbstractTaint, taint1_write AbstractTaint, db2_read *backends.Database, db1_write *backends.Database, field2_read *backends.Field, field1_write *backends.Field) bool {
	/* if field1_write.GetPath() == "order_db.order.FromStation" && field2_read.GetPath() == "station_db.station.Name" {
		// EVAL: logrus.Tracef("CURRENT TAINT: %s\n", taint2_read.LongString())
		// EVAL: logrus.Tracef("OTHER TAINT: %s\n", taint1_write.LongString())
		logrus.Fatalf("NOTE: THIS IS WHY WE NEED A SECOND SCHEMA BUILDER ITERATION!")
	} */

	var modified bool
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
	if constraint := field2_read.GetConstraintForeignKeyToField(field1_write); constraint != nil {
		if taint1_write.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [WRITE-READ] [0A] disabled mandatory: %s\n", constraint)
			}
		}
	} else if constraint := field1_write.GetConstraintForeignKeyToField(field2_read); constraint != nil {
		if taint1_write.IsWrite() {
			if ok := constraint.DisableMandatory(reqIdx); ok {
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [WRITE-READ] [0B] disabled mandatory: %s\n", constraint)
			}
		}
	} else if !field2_read.HasConstraintForeignKeyToField(field1_write) && !field1_write.HasConstraintForeignKeyToField(field2_read) {
		// WRITE .. READ
		// field_write --FK--> field_read
		if taint2_read.IsReadValue() {
			return false
		}

		if ok := createTransitiveReferenceIfExists(field1_write, field2_read, reqIdx, false); ok {
			modified = true
		} else {
			constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field1_write, field2_read)
			// must (un)set mandatory before calling GetSchema().AddConstraint()
			constraint.DisableMandatory(reqIdx)
			field1_write.AddConstraint(constraint)
			schema := db1_write.GetLastSchema()
			schema.AddConstraint(constraint)
			updateTransitiveReferencesTriggeredByCurrent(graph, schema, constraint)
			modified = true
			// EVAL: logrus.Tracef("\t\t[ITERATOR] [READ-WRITE] [2] added new constraint: %s\n", constraint)
		}
	}
	return modified
}

func propagateTaintsReadReadPair(graph *AbstractCallGraph, reqIdx int, taint2 AbstractTaint, taint1 AbstractTaint, db2 *backends.Database, db1 *backends.Database, field2 *backends.Field, field1 *backends.Field) bool {
	if !config.Global.CreateReferencesFromReadReadPair {
		return false
	}

	var modified bool
	if !field2.HasConstraintForeignKeyToField(field1) && !field1.HasConstraintForeignKeyToField(field2) {
		if taint1.IsReadKey() && taint2.IsReadKey() {
			// foreign key: field1 <--- field2

			if field2.HasConstraintForeignKey() {
				// original reference origin could actually be another field
				return false
			}

			if ok := createTransitiveReferenceIfExists(field2, field1, reqIdx, false); ok {
				modified = true
			} else {
				constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field2, field1)
				constraint.DisableMandatory(reqIdx)
				field2.AddConstraint(constraint)
				db2.GetLastSchema().AddConstraint(constraint)
				//updateTransitiveReferencesTriggeredByCurrent(graph, constraint)
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [READ-READ] [KEY-KEY] added new constraint: %s\n", constraint)
			}
		} else if taint1.IsReadValue() && taint2.IsReadKey() {
			// foreign key: field1 ---> field2
			if !config.Global.CreateReferencesFromReadReadPairAndValKey {
				return false
			}
			if ok := createTransitiveReferenceIfExists(field1, field2, reqIdx, false); ok {
				modified = true
			} else {
				constraint := backends.NewConstraint(backends.CONSTRAINT_FOREIGN_KEY, field1, field2)
				constraint.DisableMandatory(reqIdx)
				field1.AddConstraint(constraint)
				db1.GetLastSchema().AddConstraint(constraint)
				//updateTransitiveReferencesTriggeredByCurrent(graph, constraint)
				modified = true
				// EVAL: logrus.Tracef("\t\t[ITERATOR] [READ-READ] [VAL-KEY] added new constraint: %s\n", constraint)
			}
		} else {
			// sanity check
			// NOTE: it's happening e.g., caches in dsb_sn
			// logrus.Fatalf("\t\t[ITERATOR] [READ-READ] [VAL-VAL] unexpected val-val pair: (%s, %s)\n", taint1.String(), taint2.String())
		}
	}
	return modified
}

func PropagateNewTaintsToDatabaseCallObjects(graph *AbstractCallGraph, node *AbstractNode, taintMapping *TaintMapping, readOnly bool) {
	for _, edge := range graph.GetEdgesFromNode(node) {
		if edge.GetEdgeType() == EDGE_DATABASE_CALL {
			for _, obj := range edge.GetArguments() {
				for _, currTaint := range taintMapping.GetMappingKeys() {
					if config.Global.DualPassSchemaBuilding && readOnly && !currTaint.IsRead() {
						logrus.WithField("dual_pass", config.Global.DualPassSchemaBuilding).WithField("read", currTaint.IsRead()).
							Tracef("[PROPAGATE DBS] skipping read taint...")
						continue
					}
					otherTaintsLst := taintMapping.GetMappingForKey(currTaint)
					objpath, found := obj.FindObjectPathWithEqualOrUpperTaint(currTaint)
					for _, otherTaint := range otherTaintsLst {
						if found {
							obj.AddTaintIfSimilarNotExists(objpath, otherTaint)
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
func PropagateTaintsToServiceCallObjects(graph *AbstractCallGraph, currNode *AbstractNode, taintMapping *TaintMapping, currEdge *AbstractEdge, propagateFromNode bool, readOnly bool) {
	if propagateFromNode {
		for _, otherEdge := range graph.GetEdgesFromNode(currNode) {
			// propagate from params in current node to call arguments in other edge
			for _, param := range currNode.GetParams() {
				// EVAL: logrus.Tracef("[TRACE] [FROM_NODE] [PARAM] [NODE=%s] param={%s} // otherEdge={%s}\n", currNode.String(), param.String(), otherEdge.String())
				taintTracedObjectsOnEdge(param, currNode, otherEdge, taintMapping, true, readOnly)
			}
		}
	} else {
		// propagate from call arguments (1) or returns (2) in current edge to objects acting as:
		// (a) parameters or returns in the current node
		// (b) arguments in other edges

		// 1. propagate from call arguments
		for _, arg := range currEdge.GetArguments() {
			// 1a. to objects acting as parameters or returns in the current node
			taintTracedObjectsOnNode(arg, currNode, nil, taintMapping, true, readOnly)
		}
		// 2. propagate from call returns
		for _, ret := range currEdge.GetReturns() {
			// 2a. to objects acting as parameters or returns in the current node
			taintTracedObjectsOnNode(ret, currNode, nil, taintMapping, true, readOnly)
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
				// EVAL: logrus.Tracef("[TRACE] [FROM_EDGE] [ARG] [NODE=%s] arg={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", currNode.String(), arg.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				// 1b. to objects acting as arguments in other edges
				taintTracedObjectsOnEdge(arg, currNode, otherEdge, taintMapping, doTaintAfter, readOnly)
			}
			// 2. propagate from call returns
			for _, ret := range currEdge.GetReturns() {
				// EVAL: logrus.Tracef("[TRACE] [FROM_EDGE] [RET] [NODE=%s] ret={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", currNode.String(), ret.String(), currEdge.String(), otherEdge.String(), taintMapping.String())
				// 2b. to objects acting as arguments in other edges
				taintTracedObjectsOnEdge(ret, currNode, otherEdge, taintMapping, doTaintAfter, readOnly)
			}
		}
	}
}

// taintTracedObjectsOnEdge checks traces on objects used for other calls (aka edge)
func taintTracedObjectsOnEdge(currObj *AbstractObject, currNode *AbstractNode, otherEdge *AbstractEdge, taintMapping *TaintMapping, doTaintAfter bool, readOnly bool) {
	for currObjpath, tracesLst := range currObj.GetTraces() {
		// e.g.,
		// MediaMicroservices in APIService.ReadPage(...)
		//
		// movieId := movieIdService.ReadMovieId(title)
		// movieInfo := movieInfoService.ReadMovieInfo(movieId.ID)
		//
		// t4 = ReadMovieId(..) => currObjpath 	 (@ t4.MovieID) = _obj.MovieID
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
			// CURRENT OBJECT is t13 w/ currObjpath = _obj
			// TRACED OBJECT is t18 w/ tracedObjpath = _obj.ID
			//
			// we want to get the taints of t13 at _obj and propagate them to t18 on _obj.ID
			// REMINDER: we just simply associate the SAME dbfield to t18 on _obj.ID

			// we get exactly the matching object by looking for the trace argument name
			if tracedObj := otherEdge.GetArgumentByNameIfExists(trace.GetArgumentName()); tracedObj != nil {
				// EVAL: logrus.Tracef("[TRACE] [OBJ=%s // OBJPATH=%s] trace={%s}\n", currObj.String(), currObjpath, trace.LongString())
				tracedObjPath := trace.GetArgumentPath()
				taintTracedObjectsHelper(currObj, tracedObj, currObjpath, tracedObjPath, trace, taintMapping, true, doTaintAfter, readOnly)
			}
		}
	}
}

// taintTracedObjectsOnNode checks traces on objects used as parameters or returns in the current function (aka node)
func taintTracedObjectsOnNode(obj *AbstractObject, currNode *AbstractNode, otherEdge *AbstractEdge, taintMapping *TaintMapping, doTaintAfter bool, readOnly bool) {
	for objpath, tracesLst := range obj.GetTraces() {
		for _, trace := range tracesLst {
			var tracedObjPaths []string
			var tracedObjs []*AbstractObject

			// EVAL: logrus.Tracef("[TRACE] [OBJ=%s // OBJPATH=%s] trace={%s}\n", obj.String(), objpath, trace.LongString())

			for _, param := range currNode.GetParams() {
				for paramObjpath, paramTraceLst := range param.GetTraces() {
					for _, paramTrace := range paramTraceLst {
						if paramTrace.GetServiceCallID() == trace.GetServiceCallID() {
							if paramTrace.GetServicePath() == trace.GetServicePath() {
								tracedObjs = append(tracedObjs, param)
								tracedObjPaths = append(tracedObjPaths, paramObjpath)
								// EVAL: logrus.Tracef("[TRACE] [ON_NODE] [PARAM] param: %s\n", param.String())
								// EVAL: logrus.Tracef("[TRACE] [ON_NODE] [PARAM] param trace call ID: %s\n", paramTrace.GetServiceCallID())
								// EVAL: logrus.Tracef("[TRACE] [ON_NODE] [PARAM] param trace path: %s\n", paramTrace.GetServicePath())
								// EVAL: logrus.Tracef("[TRACE] [ON_NODE] [PARAM] trace path: %s\n", trace.GetServicePath())
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
								// EVAL: logrus.Tracef("[TRACE] [ON_NODE] [RET] ret trace call ID: %s\n", retTrace.GetServiceCallID())
								// EVAL: logrus.Tracef("[TRACE] [ON_NODE] [RET] ret trace path: %s\n", retTrace.GetServicePath())
								// EVAL: logrus.Tracef("[TRACE] [ON_NODE] [RET] trace path: %s\n", trace.GetServicePath())
							}
						}
					}
				}
			}

			for i, tracedObj := range tracedObjs {
				// REMINDER: traceObjPath is simply the objpath of the traced object
				taintTracedObjectsHelper(obj, tracedObj, objpath, tracedObjPaths[i], trace, taintMapping, false, doTaintAfter, readOnly)

			}
		}
	}
}

func taintTracedObjectsHelper(currObj *AbstractObject, tracedObj *AbstractObject, currObjPath string, tracedObjPath string, trace *AbstractTrace, taintMapping *TaintMapping, onEdge bool, after bool, readOnly bool) {
	// EVAL: logrus.Tracef("[TRACE] [ON_EDGE=%t] [OBJ=%s // OBJPATH=%s] corresponding trace obj (path=%s): %s\n", onEdge, currObj.String(), currObjPath, tracedObjPath, tracedObj.String())
	var selectedTaints = make(map[string][]*AbstractTaint)
	var selectedTaintsKeys []string

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
	// after going up, we get a new potential subtaint (movieid_db.movieid)
	// (that we don't save for the current obj but only for the traced obj)
	// ------------------------------------------
	// 			_obj.ID
	// [read, secondary] @ movieid_db.movieid.ID
	// ------------------------------------------

	// Example 1
	// currObjpath = _obj
	// tracedObjpath = _obj.ID

	// Example 2
	// currObjpath = _obj.ID
	// tracedObjpath = _obj.Users.ID
	//
	// if we want to propagate taints from current obj to traced obj, then we can only, at most,
	// propagate the taints from the lower paths from the current object but NEVER the upper paths
	// because the two objects (current and traced) do not exactly match like, for example, args-parms
	currTaints, pathsDiffs := currObj.GetTaintsForCurrentAndLowerPaths(currObjPath)
	for path, taintLst := range currTaints {
		// pathDiff can be empty when paths match when checking lower paths
		selectedPath := tracedObjPath + pathsDiffs[path]

		for _, taint := range taintLst {
			if config.Global.DualPassSchemaBuilding && readOnly && !taint.IsRead() {
				logrus.WithField("dual_pass", config.Global.DualPassSchemaBuilding).WithField("read", taint.IsRead()).
					Tracef("[TRACE] skipping read taint...")
				continue
			}
			// EVAL: logrus.Tracef("[TRACE] [ON_EDGE=%t] currObjpath=%s // tracedObjpath=%s // path=%s // selectedPath=%s // taint={%s}\n", onEdge, currObjPath, tracedObjPath, path, selectedPath, taint.LongString())
			selectedTaint := taint.Copy()
			selectedTaints[selectedPath] = append(selectedTaints[selectedPath], selectedTaint)
		}
		if len(taintLst) > 0 {
			selectedTaintsKeys = append(selectedTaintsKeys, selectedPath)
		}
	}

	taintMappingTmp := MergeTaints(tracedObj, selectedTaints, selectedTaintsKeys, MERGE_MODE_TRACE, trace.GetT(), readOnly)
	// EVAL: logrus.Tracef("[TRACE] [ON_EDGE=%t] taint mapping tmp = %s\n", onEdge, taintMappingTmp.String())

	if taintMapping != nil {
		// EVAL: logrus.Tracef("[TRACE] [ON_EDGE=%t] merging taint mapping tmp into main taint mapping\n", onEdge)
		taintMapping.Join(taintMappingTmp, after)
	}

}
