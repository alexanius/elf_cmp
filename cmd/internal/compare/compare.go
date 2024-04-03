package compare

import (
  "debug/elf"
  "fmt"
  "os"
  "strings"

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

type FileInfo struct {
  Type        elf.Type  // Elf type
  Dbg         string    // Has debug info
  Size        uint64    // Total size of file
  SectionNum  int       // Number of all sections

  DebugSec    map[string]*elf.Section   // Sections with debug information
  InstrSec    map[string]*elf.Section   // Sections with executable instructions
  UDataSec    map[string]*elf.Section   // Sections with user data
  GoSec       map[string]*elf.Section   // Sections related to Go lang
  CompilerSec map[string]*elf.Section   // Sections with compiler data
  OtherSec    map[string]*elf.Section   // All other sections
}

var A, B *FileInfo

func newFileInfo() *FileInfo{
  var res FileInfo
  res.Dbg = "no"
  res.DebugSec    = make(map[string]*elf.Section)
  res.InstrSec    = make(map[string]*elf.Section)
  res.UDataSec    = make(map[string]*elf.Section)
  res.GoSec       = make(map[string]*elf.Section)
  res.CompilerSec = make(map[string]*elf.Section)
  res.OtherSec    = make(map[string]*elf.Section)
  return &res
}

var Report *report.Report

func readElf(name string) (*os.File, *elf.File, error) {
    f, err := os.Open(name)
    if err != nil {
        return nil, nil, fmt.Errorf("Failed to open elf: %w", err)
    }
    resElf, err := elf.NewFile(f)
    return f, resElf, err
}

func compareStat(f1, f2 *os.File) {
  aStat, _ := f1.Stat()
  bStat, _ := f2.Stat()
  A.Size = uint64(aStat.Size())
  B.Size = uint64(bStat.Size())
}

func compareHeaders(f1, f2 *elf.File) {

  A.Type = f1.Type
  B.Type = f2.Type

/*  Report.AddHeaderRow(f1.Class.String(), "Class", f2.Class.String())
  Report.AddHeaderRow(f1.Data.String(), "Data", f2.Data.String())
  Report.AddHeaderRow(f1.Version.String(), "Version", f2.Version.String())
  Report.AddHeaderRow(f1.OSABI.String(), "OSABI", f2.OSABI.String())
  Report.AddHeaderRow(strconv.Itoa(int(f1.ABIVersion)), "ABIVersion", strconv.Itoa(int(f2.ABIVersion)))
  Report.AddHeaderRow(f1.ByteOrder.String(), "ByteOrder", f2.ByteOrder.String())
  Report.AddHeaderRow(f1.Machine.String(), "Machine", f2.Machine.String())
  Report.AddHeaderRow(strconv.FormatUint(f1.Entry, 16), "Entry", strconv.FormatUint(f2.Entry, 16))*/
}

func compareSections(f1, f2 *elf.File) {
  A.SectionNum = len(f1.Sections)
  B.SectionNum = len(f2.Sections)

  type SectionPair struct {
    s1 *elf.Section
    s2 *elf.Section
  }

  commonSections := make(map[string]SectionPair)
  f1OnlySections := make(map[string]*elf.Section)
  f2OnlySections := make(map[string]*elf.Section)

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

  for _, s1 := range f1.Sections {
    n1 := s1.SectionHeader.Name
    both := false

    if n1 == "" {
      continue
    }

    if isDebug(n1) {
      A.DebugSec[n1] = s1
      A.Dbg = "yes"
    } else if isUserData(n1) {
      A.UDataSec[n1] = s1
    } else if isGoSpecific(n1) {
      A.GoSec[n1] = s1
    } else if isCompilerSpecific(n1) {
      A.CompilerSec[n1] = s1
    } else if s1.SectionHeader.Flags & elf.SHF_EXECINSTR != 0 {
      A.InstrSec[n1] = s1
    } else {
      A.OtherSec[n1] = s1
    }

    for _, s2 := range f2.Sections {
      n2 := s2.SectionHeader.Name
      if n1 == n2 {
        commonSections[n1] = SectionPair{s1, s2}
        both = true
      }
    }

    if !both {
      f1OnlySections[n1] = s1
    }
  }

  for _, s2 := range f2.Sections {
    _, ok := commonSections[s2.SectionHeader.Name]
    n2 := s2.SectionHeader.Name

    if n2 == "" {
      continue
    }

    if isDebug(n2) {
      B.Dbg = "yes"
      B.DebugSec[n2] = s2
    } else if isUserData(n2) {
      B.UDataSec[n2] = s2
    } else if isGoSpecific(n2) {
      B.GoSec[n2] = s2
    } else if isCompilerSpecific(n2) {
      B.CompilerSec[n2] = s2
    } else if s2.SectionHeader.Flags & elf.SHF_EXECINSTR != 0 {
      B.InstrSec[n2] = s2
    } else {
      B.OtherSec[n2] = s2
    }

    if !ok {
      f2OnlySections[s2.SectionHeader.Name] = s2
    }
  }
}

// analyzeSectionGroup takes a particular group of sections, counts their total
// size and adds the rows with sections size and total size into table
func analyzeSectionGroup(aS, bS map[string]*elf.Section, gName string) {
  secSize1 := uint64(0)
  secSize2 := uint64(0)
  for _, s1 := range aS {
    sName := s1.SectionHeader.Name
    size1 := s1.SectionHeader.Size
    secSize1 += size1
    s2, ok := bS[sName]
    if ok {
      size2 := s2.SectionHeader.Size
      Report.AddIntRowGroup(gName, sName, size1, size2)
      secSize2 += size2
    } else {
      Report.AddIntRow1Group(gName, sName, size1)
    }
  }
  for _, s2 := range bS {
    sName := s2.SectionHeader.Name
    _, ok := aS[sName]
    if ok {
      continue
    }
    size2 := s2.SectionHeader.Size
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
    fmt.Sprintf("%d", A.SectionNum),
    fmt.Sprintf("%d", B.SectionNum))
  Report.AddIntRow ("Size", A.Size, B.Size)
  Report.AddSubtitle("Sections")

  analyzeSectionGroup(A.InstrSec, B.InstrSec, "Instr")
  analyzeSectionGroup(A.UDataSec, B.UDataSec, "User data")
  analyzeSectionGroup(A.GoSec, B.GoSec, "Go data")
  analyzeSectionGroup(A.CompilerSec, B.CompilerSec, "Compiler data")
  analyzeSectionGroup(A.DebugSec, B.DebugSec, "Debug info")
  analyzeSectionGroup(A.OtherSec, B.OtherSec, "Other")
}

func Compare(fname1, fname2 string) error {
  Report = report.New(fname1, fname2)
  A = newFileInfo()
  B = newFileInfo()

  f1, elf1, err := readElf(fname1)
  if err != nil {
    return fmt.Errorf("Failed to open elf: %w", err)
  }
  f2, elf2, err := readElf(fname2)
  if err != nil {
    return fmt.Errorf("Failed to open elf: %w", err)
  }

  compareHeaders(elf1, elf2)
  compareSections(elf1, elf2)
  compareStat(f1, f2)
  fillTable()
  Report.Print()
  return nil
}
