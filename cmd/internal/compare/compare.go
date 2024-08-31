package compare

import (
  "fmt"

  "elf_cmp/cmd/internal/file"
  "elf_cmp/cmd/internal/report"
)

var A, B *file.FileInfo

var Report *report.Report

// analyzeSectionGroup takes a particular group of sections, counts their total
// size and adds the rows with sections size and total size into table
func analyzeSectionGroup(cmp *report.Compare, aS, bS map[string]*file.Section, gName string) {
  secSize1 := uint64(0)
  secSize2 := uint64(0)
  cmp.Secs[gName] = &report.SecCompare{}
  for _, s1 := range aS {
    sName := s1.Info.SectionHeader.Name
    size1 := s1.Info.SectionHeader.Size
    secSize1 += size1
    s2, ok := bS[sName]
    if ok {
      size2 := s2.Info.SectionHeader.Size
      Report.AddIntRowGroup(gName, sName, size1, size2)
      secSize2 += size2
      cmp.Secs[gName].ComonSections = append(cmp.Secs[gName].ComonSections,
        &report.SectionPair{s1, s2})
    } else {
      Report.AddIntRow1Group(gName, sName, size1)
      cmp.Secs[gName].Asections = append(cmp.Secs[gName].Asections, s1)
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
    cmp.Secs[gName].Bsections = append(cmp.Secs[gName].Bsections, s2)
  }
  Report.AddIntRowGroup(gName, "Total", secSize1, secSize2)
  Report.AddSeparator()
}

// analyzeSymbolGroup 
func analyzeSymbolGroup(cmp *report.Compare, aS, bS map[string]*file.Section, gName string) {
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

func fillTable(cmp *report.Compare) {
  Report.AddTextRow("Type", A.ElfType(), B.ElfType())
  Report.AddTextRow("Debug info", A.Dbg, B.Dbg)
  Report.AddTextRow ("Sections",
    fmt.Sprintf("%d", A.SectionNum()),
    fmt.Sprintf("%d", B.SectionNum()))
  Report.AddTextRow ("Symbols",
    fmt.Sprintf("%d", A.SymbolNum()),
    fmt.Sprintf("%d", B.SymbolNum()))
  Report.AddIntRow ("Size", A.Size, B.Size)
  Report.AddSubtitle("Sections size (bytes)")

  cmp.Secs = make(map[string]*report.SecCompare)
  analyzeSectionGroup(cmp, A.InstrSec, B.InstrSec, "Instr")
  analyzeSectionGroup(cmp, A.UDataSec, B.UDataSec, "User data")
  analyzeSectionGroup(cmp, A.GoSec, B.GoSec, "Go data")
  analyzeSectionGroup(cmp, A.CompilerSec, B.CompilerSec, "Compiler data")
  analyzeSectionGroup(cmp, A.DebugSec, B.DebugSec, "Debug info")
  analyzeSectionGroup(cmp, A.OtherSec, B.OtherSec, "Other")

  Report.AddSubtitle("Sections symbols number")

  analyzeSymbolGroup(cmp, A.InstrSec, B.InstrSec, "Instr")
  analyzeSymbolGroup(cmp, A.UDataSec, B.UDataSec, "User data")
  analyzeSymbolGroup(cmp, A.GoSec, B.GoSec, "Go data")
  analyzeSymbolGroup(cmp, A.CompilerSec, B.CompilerSec, "Compiler data")
  analyzeSymbolGroup(cmp, A.DebugSec, B.DebugSec, "Debug info")
  analyzeSymbolGroup(cmp, A.OtherSec, B.OtherSec, "Other")
}

func Compare(fname1, fname2 string, html bool) error {
  A, _ = file.CreateFileInfo(fname1)
  B, _ = file.CreateFileInfo(fname2)
  Report = report.New(A, B)
  var cmp report.Compare

  fillTable(&cmp)
  if html {
    Report.PrintHtml(&cmp)
  } else {
    Report.Print()
  }
  return nil
}
