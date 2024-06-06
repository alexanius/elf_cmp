package file

import (
  "fmt"
  "os"
  "strings"

  "debug/elf"
)

type Section struct {
  Info    *elf.Section
  Symbols map[string]*elf.Symbol
}

type FileInfo struct {
  File        *elf.File // ELF file
  Type        elf.Type  // Elf type
  Dbg         string    // Has debug info
  Size        uint64    // Total size of file
  SectionNum  int       // Number of all sections
  SymbolNum   int       // Number of all symbols

  DebugSec    map[string]*Section   // Sections with debug information
  InstrSec    map[string]*Section   // Sections with executable instructions
  UDataSec    map[string]*Section   // Sections with user data
  GoSec       map[string]*Section   // Sections related to Go lang
  CompilerSec map[string]*Section   // Sections with compiler data
  OtherSec    map[string]*Section   // All other sections

  AllSections []*Section            // All the sections sorted by offset
  AllSymbols  []*elf.Symbol         // All other symbols sorted by value
}

func newFileInfo() *FileInfo {
  var res FileInfo
  res.Dbg = "no"
  res.DebugSec    = make(map[string]*Section)
  res.InstrSec    = make(map[string]*Section)
  res.UDataSec    = make(map[string]*Section)
  res.GoSec       = make(map[string]*Section)
  res.CompilerSec = make(map[string]*Section)
  res.OtherSec    = make(map[string]*Section)
  return &res
}

func (fi *FileInfo) readSection() {
  isDebug := func(name string) bool {
    return strings.Contains(name, ".debug")
  }

  isUserData := func(name string) bool {
    return name == ".data" || name == ".bss" || name == ".rodata"
  }

  isGoSpecific := func(name string) bool {
    return name == ".typelink" || name == ".gosymtab" || name == ".noptrdata" ||
      name == ".gopclntab" || name == ".noptrbss" || name == ".itablink"
  }

  isCompilerSpecific := func(name string) bool {
    return name == ".note.go.buildid" || name == ".go.buildinfo" ||
      name == ".note.gnu.property" || name == ".note.ABI-tag" ||
      name == ".gnu.version" || name == ".gnu.version_r" ||
      name == ".gnu.hash" || name == ".gcc_except_table"
  }

  fi.AllSections = make([]*Section, len(fi.File.Sections))

  for i, s := range fi.File.Sections {
    n := s.SectionHeader.Name
    s1 := &Section{Info : s}
    fi.AllSections[i] = s1

    if n == "" {
      continue
    }

    if isDebug(n) {
      fi.DebugSec[n] = s1
      fi.Dbg = "yes"
    } else if isUserData(n) {
      fi.UDataSec[n] = s1
    } else if isGoSpecific(n) {
      fi.GoSec[n] = s1
    } else if isCompilerSpecific(n) {
      fi.CompilerSec[n] = s1
    } else if s1.Info.SectionHeader.Flags & elf.SHF_EXECINSTR != 0 {
      fi.InstrSec[n] = s1
    } else {
      fi.OtherSec[n] = s1
    }
  }
}

func (fi *FileInfo) readSymbols() {
/*  curSym := 0
  symNum := len(fi.AllSymbols)
  for curSec, sec := range fi.AllSections {
    nextSec := sec.Info.Offset + sec.Info.Size
    for ; curSym < symNum ; curSym++ {
      sym := fi.AllSymbols[curSym]
      if sym.Value > nextSec {
        fi.AllSections[curSec+1].Symbols[sym.Name] = sym
        break
      }
      fi.AllSections[curSec].Symbols[sym.Name] = sym
    }
  }*/
}

func CreateFileInfo(name string) (*FileInfo, error) {
  f, err := os.Open(name)
  if err != nil {
    return nil, fmt.Errorf("Failed to open file: %w", err)
  }
  resElf, err := elf.NewFile(f)

  if err != nil {
    return nil, fmt.Errorf("Failed to read file: %w", err)
  }

  info := newFileInfo()
  info.File = resElf
  info.Type = resElf.Type
  stat, _ := f.Stat()
  info.Size = uint64(stat.Size())
  info.readSection()

  return info, nil
}

