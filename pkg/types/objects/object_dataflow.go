package objects

import (
	"fmt"
	"strings"

	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type ObjectDataflow struct {
	Datastore      string
	Service        string
	Write          bool
	Direct         bool
	Variable       Object
	IndirectSource Object
	Field          *datastores.Field
	RequestIdx     int
	perm           bool
}

func (df *ObjectDataflow) GetOpString() string {
	if df.Write {
		return "write"
	}
	return "read"
}

func (df *ObjectDataflow) InRequestIndex(requestIdx int) bool {
	return df.RequestIdx == requestIdx
}

func (df *ObjectDataflow) HasVariable(variable Object) bool {
	return df.Variable == variable
}

func (df *ObjectDataflow) HasAnyVariable(variables []Object) bool {
	for _, v := range variables {
		if df.Variable == v {
			return true
		}
		if df.IndirectSource == v {
			return true
		}
	}
	return false
}

func (df *ObjectDataflow) IsWriteOp() bool {
	return df.Write
}

func (df *ObjectDataflow) IsOpInDatastore(datastore string) bool {
	return df.Datastore == datastore
}

func (df *ObjectDataflow) String() string {
	prefix := fmt.Sprintf("%s %s",
		func() string {
			if df.Direct {
				return "direct"
			} else {
				return "indirect"
			}
		}(),
		func() string {
			if df.Write {
				return "write"
			} else {
				return "read"
			}
		}(),
	)
	if df.Direct {
		return fmt.Sprintf("%s <%s> @ (%s, %s)", prefix, df.Variable.String(), df.Service, df.Datastore)
	}
	return fmt.Sprintf("%s <%s> by <%s> @ (%s, %s)", prefix, df.Variable.String(), df.IndirectSource.String(), df.Service, df.Datastore)
}

func (df *ObjectDataflow) Copy(force bool) *ObjectDataflow {
	return df
}

func (df *ObjectDataflow) DeepCopy() *ObjectDataflow {
	return df
}

func (df *ObjectDataflow) ShortString() string {
	if df.Direct {
		return fmt.Sprintf("%s <%s (%s)>", df.GetOpString(), df.Variable.String(), utils.GetType(df.Variable))
	}
	return fmt.Sprintf("%s <%s (%s)> from <%s (%s)>", df.GetOpString(), df.Variable.String(), utils.GetType(df.Variable), df.IndirectSource.String(), utils.GetType(df.IndirectSource))
}

func (df *ObjectDataflow) GetVariable() Object {
	if df.Direct {
		return df.Variable
	}
	return df.IndirectSource
}

func (df *ObjectDataflow) GetDatastore() string {
	return df.Datastore
}

func (df *ObjectDataflow) EnablePermanent() {
	df.perm = true
}

func (df *ObjectDataflow) IsPermanent() bool {
	return df.perm
}

func (v *ObjectInfo) SetDirectDataflow(datastore string, service string, variable Object, field *datastores.Field, write bool, requestIdx int) *ObjectDataflow {
	df := &ObjectDataflow{
		Direct:     true,
		Variable:   variable,
		Datastore:  datastore,
		Service:    service,
		Field:      field,
		Write:      write,
		RequestIdx: requestIdx,
	}
	v.Dataflows = append(v.Dataflows, df)
	logger.Logger.Warnf("[DIRECT DATAFLOW (REQ %d) - %s @ %s] (%d) [%s] %s", requestIdx, strings.ToUpper(df.GetOpString()), field.GetFullName(), df.Variable.GetId(), utils.GetType(df.Variable), df.Variable.String())
	return df
}

func (v *ObjectInfo) SetIndirectDataflow(datastore string, service string, current Object, source Object, field *datastores.Field, write bool, requestIdx int) *ObjectDataflow {
	df := &ObjectDataflow{
		Variable:       current,
		IndirectSource: source,
		Datastore:      datastore,
		Service:        service,
		Field:          field,
		Write:          write,
		RequestIdx:     requestIdx,
	}
	current.GetVariableInfo().IndirectDataflows = append(current.GetVariableInfo().IndirectDataflows, df)
	v.Dataflows = append(v.Dataflows, df)
	logger.Logger.Warnf("\t\t[INDIRECT DATAFLOW (REQ %d) - %s @ %s] (%d) [%s] %s", requestIdx, strings.ToUpper(df.GetOpString()), field.GetFullName(), df.Variable.GetId(), utils.GetType(df.Variable), df.Variable.String())
	return df
}

func (v *ObjectInfo) GetForeignDataflows(currentDS *datastores.Datastore) []*ObjectDataflow {
	var dataflows []*ObjectDataflow
	for _, df := range v.GetDataflows() {
		// only append if it a different datastore
		if df.Datastore != currentDS.Name {
			dataflows = append(dataflows, df)
		}
	}
	return dataflows
}
