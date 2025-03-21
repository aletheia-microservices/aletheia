package types

import (
	"fmt"
	"go/ast"

	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type Method interface {
	String() string
	LongString() string
	GetParams() []*MethodField
	GetReturns() []*MethodField
	IsQueueWrite() bool
	IsQueueRead() bool
	IsWrite() bool
	IsDelete() bool
	IsUpdate() bool
	IsRead() bool
	IsFetchNoSQLCollection() bool
	DeepCopy() Method
}

type ParsedMethod struct {
	Method
	Ast         *ast.FuncDecl                 `json:"-"`
	Name        string                        `json:"name"`
	Calls       []Call                        `json:"-"` // calls made within method
	Package     *Package                      `json:"-"`
	ParsedCfg   *CFG                          `json:"-"`
	DbInstances []datastores.DatabaseInstance `json:"-"`
	Exported    bool                          `json:"-"`

	// true if all blocks of the cfg are parsed
	// note that ParsedCfg can be != nil before being parsed
	Parsed          bool   `json:"-"`
	AttachedService string `json:"-"`
	Constructor     bool   `json:"-"`

	// used to fetch the params when generating the basic cfg
	// to store in the variables array of the function
	Params   []*MethodField
	Returns  []*MethodField
	Receiver *MethodField `json:"-"`
}

func (f *ParsedMethod) DeepCopy() Method {
	return &ParsedMethod{
		Ast:             f.Ast,
		Name:            f.Name,
		Calls:           f.Calls,
		Package:         f.Package,
		ParsedCfg:       f.ParsedCfg.DeepCopy(),
		DbInstances:     f.DbInstances,
		Exported:        f.Exported,
		Parsed:          f.Parsed,
		AttachedService: f.AttachedService,
		Constructor:     f.Constructor,
		Params:          f.Params,
		Returns:         f.Returns,
		Receiver:        f.Receiver,
	}
}
func (f *ParsedMethod) GetAst() *ast.FuncDecl {
	return f.Ast
}
func (f *ParsedMethod) GetName() string {
	return f.Name
}
func (f *ParsedMethod) SetParsed() {
	f.Parsed = true
}
func (f *ParsedMethod) IsParsed() bool {
	return f.Parsed
}
func (f *ParsedMethod) EnableConstructorFlag() {
	f.Constructor = true
}
func (f *ParsedMethod) IsConstructor() bool {
	return f.Constructor
}
func (f *ParsedMethod) AttachService(service string) {
	f.AttachedService = service
}
func (f *ParsedMethod) HasAttachedService() bool {
	return f.AttachedService != ""
}
func (f *ParsedMethod) GetReceiverIfExists() *MethodField {
	return f.Receiver
}
func (f *ParsedMethod) HasReceiver() bool {
	return f.Receiver != nil
}
func (f *ParsedMethod) GetReceiverType() gotypes.Type {
	return f.Receiver.GetType()
}
func (f *ParsedMethod) GetReceiver() *MethodField {
	return f.Receiver
}
func (f *ParsedMethod) GetBody() *ast.BlockStmt {
	return f.Ast.Body
}
func (f *ParsedMethod) SetParsedCFG(parsedCfg *CFG) {
	f.ParsedCfg = parsedCfg
}
func (f *ParsedMethod) GetParsedCfg() *CFG {
	if f.ParsedCfg == nil {
		logger.Logger.Fatalf("[METHOD] unexpected nil parsed cfg for method (%s)", f.String())
	}
	return f.ParsedCfg
}
func (*ParsedMethod) IsQueueWrite() bool {
	return false
}
func (*ParsedMethod) IsQueueRead() bool {
	return false
}

func (p *ParsedMethod) String() string {
	prefix := ""
	if p.AttachedService != "" {
		prefix += p.AttachedService
	}
	if p.Receiver != nil {
		// this just resets everything
		prefix = p.Receiver.GetTypeName()
	}
	repr := fmt.Sprintf("%s.%s(", prefix, p.Name)
	for i, arg := range p.Params {
		repr += arg.String()
		if i < len(p.Params)-1 {
			repr += ", "
		}
	}
	repr += ")"
	return repr
}

func (p *ParsedMethod) SimpleString() string {
	repr := fmt.Sprintf("%s(", p.Name)
	for i, arg := range p.Params {
		repr += arg.String()
		if i < len(p.Params)-1 {
			repr += ", "
		}
	}
	repr += ")"
	return repr
}

func (p *ParsedMethod) LongString() string {
	prefix := ""
	if p.Package != nil {
		prefix += p.Package.Name
	}
	if p.AttachedService != "" {
		prefix += p.AttachedService
	}
	if p.Receiver != nil {
		// this just resets everything
		prefix = p.Receiver.GetTypeLongName()
	}
	repr := fmt.Sprintf("%s.%s(", prefix, p.Name)
	for i, arg := range p.Params {
		repr += arg.String()
		if i < len(p.Params)-1 {
			repr += ", "
		}
	}
	repr += ")"
	if len(p.Returns) == 1 {
		repr += " " + p.Returns[0].String()
	}
	if len(p.Returns) > 1 {
		repr += " ("
		for i, ret := range p.Returns {
			repr += ret.String()
			if i < len(p.Returns)-1 {
				repr += ", "
			}
		}
		repr += ")"
	}
	return repr
}

// TODO: missing returns!!
func (p *ParsedMethod) Signature() string {
	repr := fmt.Sprintf("%s(", p.Name)
	for i, arg := range p.Params {
		repr += arg.String()
		if i < len(p.Params)-1 {
			repr += ", "
		}
	}
	repr += ")"
	if len(p.Returns) == 0 {
		return repr
	}
	if len(p.Returns) == 1 {
		return repr + " " + p.Returns[0].String()
	}

	repr += " ("
	for i, arg := range p.Returns {
		repr += arg.String()
		if i < len(p.Returns)-1 {
			repr += ", "
		}
	}
	return repr + ")"
}

func (p *ParsedMethod) GetParams() []*MethodField {
	return p.Params
}
func (p *ParsedMethod) GetReturns() []*MethodField {
	return p.Returns
}
