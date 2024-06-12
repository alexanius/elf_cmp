package compare

import (
  "debug/elf"
  "fmt"

  "elf_cmp/cmd/internal/file"
  "elf_cmp/cmd/internal/report"
)

var ElfType = map[elf.Type]string{
  elf.ET_NONE   : "No file type",
  elf.ET_REL    : "Relocatable file",
  elf.ET_EXEC   : "Executable file",
  elf.ET_DYN    : "Shared object file",
  elf.ET_CORE   : "Core file",
  elf.ET_LOOS   : "First operating system specific",
  elf.ET_HIOS   : "Last operating system-specific",
  elf.ET_LOPROC : "Processor-specific",
  elf.ET_HIPROC : "Processor-specific",
}

var A, B *file.FileInfo

var Report *report.Report

// analyzeSectionGroup takes a particular group of sections, counts their total
// size and adds the rows with sections size and total size into table
func analyzeSectionGroup(aS, bS map[string]*file.Section, gName string) {
  secSize1 := uint64(0)
  secSize2 := uint64(0)
  for _, s1 := range aS {
    sName := s1.Info.SectionHeader.Name
    size1 := s1.Info.SectionHeader.Size
    secSize1 += size1
    s2, ok := bS[sName]
    if ok {
      size2 := s2.Info.SectionHeader.Size
      Report.AddIntRowGroup(gName, sName, size1, size2)
      secSize2 += size2
    } else {
      Report.AddIntRow1Group(gName, sName, size1)
    }
  }
  for _, s2 := range bS {
    sName := s2.Info.SectionHeader.Name
    _, ok := aS[sName]
    if ok {
      continue
    }
    size2 := s2.Info.SectionHeader.Size
    secSize2 += size2
    Report.AddIntRow2Group(gName, sName, size2)
  }
  Report.AddIntRowGroup(gName, "Total", secSize1, secSize2)
  Report.AddSeparator()
}

// analyzeSymbolGroup 
func analyzeSymbolGroup(aS, bS map[string]*file.Section, gName string) {
  secSize1 := uint64(0)
  secSize2 := uint64(0)
  for _, s1 := range aS {
    sName := s1.Info.SectionHeader.Name
    size1 := uint64(len(s1.Symbols))
    secSize1 += size1
    s2, ok := bS[sName]
    if ok {
      size2 := uint64(len(s2.Symbols))
      Report.AddIntRowGroup(gName, sName, size1, size2)
      secSize2 += size2
    } else {
      Report.AddIntRow1Group(gName, sName, size1)
    }
  }
  for _, s2 := range bS {
    sName := s2.Info.SectionHeader.Name
    _, ok := aS[sName]
    if ok {
      continue
    }
    size2 := uint64(len(s2.Symbols))
    secSize2 += size2
    Report.AddIntRow2Group(gName, sName, size2)
  }
  Report.AddIntRowGroup(gName, "Total", secSize1, secSize2)
  Report.AddSeparator()
}

func fillTable() {
  Report.AddTextRow("Type", ElfType[A.Type], ElfType[B.Type])
  Report.AddTextRow("Debug info", A.Dbg, B.Dbg)
  Report.AddTextRow ("Sections",
    fmt.Sprintf("%d", A.SectionNum()),
    fmt.Sprintf("%d", B.SectionNum()))
  Report.AddTextRow ("Symbols",
    fmt.Sprintf("%d", A.SymbolNum()),
    fmt.Sprintf("%d", B.SymbolNum()))
  Report.AddIntRow ("Size", A.Size, B.Size)
  Report.AddSubtitle("Sections size (bytes)")

  analyzeSectionGroup(A.InstrSec, B.InstrSec, "Instr")
  analyzeSectionGroup(A.UDataSec, B.UDataSec, "User data")
  analyzeSectionGroup(A.GoSec, B.GoSec, "Go data")
  analyzeSectionGroup(A.CompilerSec, B.CompilerSec, "Compiler data")
  analyzeSectionGroup(A.DebugSec, B.DebugSec, "Debug info")
  analyzeSectionGroup(A.OtherSec, B.OtherSec, "Other")

  Report.AddSubtitle("Sections symbols number")

  analyzeSymbolGroup(A.InstrSec, B.InstrSec, "Instr")
  analyzeSymbolGroup(A.UDataSec, B.UDataSec, "User data")
  analyzeSymbolGroup(A.GoSec, B.GoSec, "Go data")
  analyzeSymbolGroup(A.CompilerSec, B.CompilerSec, "Compiler data")
  analyzeSymbolGroup(A.DebugSec, B.DebugSec, "Debug info")
  analyzeSymbolGroup(A.OtherSec, B.OtherSec, "Other")
}

func Compare(fname1, fname2 string) error {
  Report = report.New(fname1, fname2)
  A, _ = file.CreateFileInfo(fname1)
  B, _ = file.CreateFileInfo(fname2)

  fillTable()
  Report.Print()
  return nil
}
