package types

import (
	"fmt"

	"golang.org/x/tools/go/cfg"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type CFG struct {
	Cfg          *cfg.CFG
	ParsedBlocks []*Block
	Package      string
	FileHash     string
	FullMethod   string
	HasReceiver  bool
	ReceiverType gotypes.Type
}

func InitParsedCFG(cfg *cfg.CFG, fullMethod string) *CFG {
	parsedCfg := &CFG{
		Cfg:        cfg,
		FullMethod: fullMethod,
	}
	for i, block := range cfg.Blocks {
		parsedCfg.ParsedBlocks = append(parsedCfg.ParsedBlocks, &Block{Block: block, Index: i})
	}
	for _, parsedBlock := range parsedCfg.ParsedBlocks {
		for _, succ := range parsedBlock.Block.Succs {
			parsedBlock.Successors = append(parsedBlock.Successors, parsedCfg.ParsedBlocks[succ.Index])
		}
	}
	return parsedCfg
}

func (cfg *CFG) DeepCopy() *CFG {
	var parsedBlockCopies []*Block
	for _, b := range cfg.ParsedBlocks {
		parsedBlockCopies = append(parsedBlockCopies, b.DeepCopy())
	}
	return &CFG{
		Cfg:          cfg.Cfg,
		ParsedBlocks: parsedBlockCopies,
		Package:      cfg.Package,
		FileHash:     cfg.FileHash,
		FullMethod:   cfg.FullMethod,
		HasReceiver:  cfg.HasReceiver,
		ReceiverType: cfg.ReceiverType,
	}
}

func (cfg *CFG) String() string {
	str := fmt.Sprintf("%s: ", cfg.FullMethod)
	for i, block := range cfg.ParsedBlocks {
		str += block.String()
		if i < len(cfg.ParsedBlocks)-1 {
			str += ", "
		}
	}
	return str + "\n"
}

func (cfg *CFG) LongString() string {
	str := "\n"
	for _, block := range cfg.ParsedBlocks {
		str += block.LongString()
	}
	return str + "\n"
}

func (cfg *CFG) GetEntryParsedBlock() *Block {
	if len(cfg.ParsedBlocks) > 0 {
		return cfg.ParsedBlocks[0]
	}
	logger.Logger.Fatalf("[CFG] unexpected empty CFG %s", cfg.String())
	return nil
}

//FIXME: THIS SHOULD NOT BE USED IN THE FUTURE
// BUT CURRENTLY WE ARE NOT CONCERNED ABOUT CONDITIONAL BRANCHES RESULTING IN OTHER BLOCKS
func (cfg *CFG) GetLastLiveBlock() *Block {
	// interate in reverse
	for i := len(cfg.ParsedBlocks) - 1; i >= 0; i-- {
		block := cfg.ParsedBlocks[i]
		if block.Block.Live {
			return block
		}
	}
	logger.Logger.Fatalf("[CFG] no live blocks for CFG: %s", cfg.String())
	return nil
}

func (cfg *CFG) GetParsedBlocks() []*Block {
	return cfg.ParsedBlocks
}

func (cfg *CFG) GetParsedBlockAtIndex(index int32) *Block {
	if len(cfg.ParsedBlocks) > int(index) {
		return cfg.ParsedBlocks[index]
	}
	return nil
}
