package compare

import (
  "elf_cmp/cmd/internal/file"
)

var A, B *file.FileInfo

type PairSection struct {
  A *file.Section
  B *file.Section
}

type CompareSection struct {
  A []*file.Section  // Only in A
  B []*file.Section  // Only in B
  C []*PairSection   // In both
}

type Compare struct {
  A *file.FileInfo
  B *file.FileInfo
  Secs CompareSection
}

// analyzeSectionGroup takes a particular group of sections, counts their total
// size and adds the rows with sections size and total size into table
func analyzeSectionGroup(cmp *Compare, A, B *file.FileInfo) {
  for _, sA := range A.AllSections {
    var bSec *file.Section
    for _, sB := range B.AllSections {
      if sB.Info.Name == sA.Info.Name {
        bSec = sB
      }
    }
    if bSec != nil {
      cmp.Secs.C = append(cmp.Secs.C, &PairSection{sA, bSec})
    } else {
      cmp.Secs.A = append(cmp.Secs.A, sA)
    }
  }

  for _, sB := range B.AllSections {
    for _, sA := range A.AllSections {
      if sB.Info.Name == sA.Info.Name {
        // Already added as pair section
        break
      }
    }
    cmp.Secs.B = append(cmp.Secs.B, sB)
  }
}

func CompareFiles(fname1, fname2 string) *Compare {
  A, _ = file.CreateFileInfo(fname1)
  B, _ = file.CreateFileInfo(fname2)
  var cmp Compare

  cmp.A = A
  cmp.B = B
  analyzeSectionGroup(&cmp, A, B)

  return &cmp
}

