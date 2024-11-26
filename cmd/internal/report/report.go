package report

import (
  "fmt"
  "os"

  "github.com/jedib0t/go-pretty/v6/table"
  "github.com/jedib0t/go-pretty/v6/text"

  "elf_cmp/cmd/internal/file"
)

type SectionPair struct {
  A *file.Section
  B *file.Section
}

type SecCompare struct {
  Asections []*file.Section // Sections only in A
  Bsections []*file.Section // Sections only in B

  // Sections, that are present in both files
  ComonSections []*SectionPair
}

// Results of compared binaries
type Compare struct {
  Secs map[string]*SecCompare // Key - group name
}

func CountPercent(a, b uint64) float64 {
  if a > b {
    return -(1 - float64(b) / float64(a)) * 100
  }
  return (float64(b) / float64(a)) * 100
}

type Report struct {
  Stat      table.Writer

  F1, F2    *file.FileInfo
}

func New(A, B *file.FileInfo) *Report{
  r := Report{}
  r.F1 = A
  r.F2 = B

  r.Stat     = table.NewWriter()
  r.Stat.SetColumnConfigs([]table.ColumnConfig{
    {Number:    1,
     AutoMerge: true,
     Align:     text.AlignLeft,
     VAlign:    text.VAlignMiddle}})
  r.Stat.SetOutputMirror(os.Stdout)
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  r.Stat.AppendRow(table.Row{"A", "A", A.Name, A.Name, A.Name}, rowConfigAutoMerge)
  r.Stat.AppendRow(table.Row{"B", "B", B.Name, B.Name, B.Name}, rowConfigAutoMerge)
  r.Stat.AppendSeparator()
  r.Stat.AppendRow(table.Row{"", "", "A", "B", "Diff"}, rowConfigAutoMerge)
  r.Stat.AppendSeparator()
  r.Stat.AppendRow(table.Row{"General info", "General info",
    "General info", "General info", "General info"}, rowConfigAutoMerge)
  r.Stat.AppendSeparator()

  return &r
}

func (r *Report) AddTextRow(name, A, B string) {
  d := ""
  if A != B {
    d = "!"
  }
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  r.Stat.AppendRow([]interface{}{name, name, A, B, d}, rowConfigAutoMerge)
}

func (r *Report) AddIntRow(name string, A, B uint64) {
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  r.Stat.AppendRow([]interface{}{
    name,
    name,
    fmt.Sprintf("%d", A),
    fmt.Sprintf("%d", B),
    fmt.Sprintf("%+.2f%%", CountPercent(A, B))},
    rowConfigAutoMerge)
}

func (r *Report) AddIntRowGroup(group, name string, A, B uint64) {
  r.Stat.AppendRow([]interface{}{
    group,
    name,
    fmt.Sprintf("%d", A),
    fmt.Sprintf("%d", B),
    fmt.Sprintf("%+.2f%%", CountPercent(A, B))})
}

func (r *Report) AddIntRow1(name string, A uint64) {
  r.Stat.AppendRow([]interface{}{
    name,
    name,
    fmt.Sprintf("%d", A),
    "",
    ""})
}

func (r *Report) AddIntRow1Group(group, name string, A uint64) {
  r.Stat.AppendRow([]interface{}{
    group,
    name,
    fmt.Sprintf("%d", A),
    "",
    ""})
}

func (r *Report) AddIntRow2(name string, B uint64) {
  r.Stat.AppendRow([]interface{}{
    name,
    name,
    "",
    fmt.Sprintf("%d", B),
    ""})
}

func (r *Report) AddIntRow2Group(group, name string, B uint64) {
  r.Stat.AppendRow([]interface{}{
    group,
    name,
    "",
    fmt.Sprintf("%d", B),
    ""})
}

func (r *Report) AddSubtitle(name string) {
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  r.Stat.AppendSeparator()
  r.Stat.AppendRow(table.Row{name, name, name, name, name}, rowConfigAutoMerge)
  r.Stat.AppendSeparator()
}

func (r *Report) AddSeparator() {
  r.Stat.AppendSeparator()
}


func (r *Report) AddStatRow(name, A, B, d string) {
  r.Stat.AppendRow([]interface{}{name, A, B, d})
}

func (r *Report) Print() {
  r.Stat.Render()
}

func generateGeneralInfoHtml(A, B *file.FileInfo) string {
  return fmt.Sprintf(`
  <table>
    <tr>
      <th>Type</th>
      <td>%s</td>
      <td>%s</td>
    </tr>
    <tr>
      <th>Debug info</th>
      <td>%s</td>
      <td>%s</td>
    </tr>
    <tr>
      <th>Sections</th>
      <td>%d</td>
      <td>%d</td>
    </tr>
    <tr>
      <th>Symbols</th>
      <td>%d</td>
      <td>%d</td>
    </tr>
    <tr>
    <th>Size</th>
      <td>%d</td>
      <td>%d</td>
    </tr>
  </table>
`, A.ElfType(),    B.ElfType(),
   A.Dbg,          B.Dbg,
   A.SectionNum(), B.SectionNum(),
   A.SymbolNum(),  B.SymbolNum(),
   A.Size,         B.Size)
}

func generateSectionsTableHtml(cmp *Compare, A, B *file.FileInfo) string {
  secTbl := "" // Table of sections
  groups := [...]string{
    "Instr",
    "User data",
    "Go data",
    "Compiler data",
    "Debug info",
    "Other",
  }

  for _, gName := range groups {
    secs := cmp.Secs[gName]
    secRow := ""
    aSize := uint64(0)
    bSize := uint64(0)
    aSyms := 0
    bSyms := 0
    for _, aSec := range secs.Asections {
      secRow += fmt.Sprintf("    <tr><td><a href=\"sections/%s.html\">%s</a></td><td>%d</td><td></td><td></td>  <td>%d</td><td></td><td></td> </tr>\n", aSec.Info.Name, aSec.Info.Name, aSec.Info.Size, len(aSec.Symbols))
      aSize += aSec.Info.Size
      aSyms += len(aSec.Symbols)
    }
    for _, sec := range secs.ComonSections {
      aSymNum := len(sec.A.Symbols)
      bSymNum := len(sec.B.Symbols)
      secRow += fmt.Sprintf("    <tr><td><a href=\"sections/%s.html\">%s</a></td><td>%d</td><td>%d</td><td>%+.4f</td>  <td>%d</td><td>%d</td><td>%+.4f</td> </tr>\n", sec.A.Info.Name, sec.A.Info.Name, sec.A.Info.Size, sec.B.Info.Size, CountPercent(sec.A.Info.Size, sec.B.Info.Size), aSymNum, bSymNum, CountPercent(uint64(aSymNum), uint64(bSymNum)))
      aSize += sec.A.Info.Size
      bSize += sec.B.Info.Size
      aSyms += aSymNum
      bSyms += bSymNum
    }
    for _, bSec := range secs.Bsections {
      secRow += fmt.Sprintf("    <tr><td><a href=\"sections/%s.html\">%s</a></td><td></td><td>%d</td><td></td>  <td></td><td>%d</td><td></td> </tr>\n", bSec.Info.Name, bSec.Info.Name, bSec.Info.Size, len(bSec.Symbols))
      bSize += bSec.Info.Size
      bSyms += len(bSec.Symbols)
    }
    secRows := len(secs.Asections) + len(secs.ComonSections) + len(secs.Bsections) + 2
    secRow = fmt.Sprintf(`%s
`, secRow)
    secTbl += fmt.Sprintf(`    <tr><th rowspan=%d>%s</th></tr>
%s`, secRows, gName, secRow)
    secTbl += fmt.Sprintf("    <tr><td>Total</td><td>%d</td><td>%d</td><td>%+.4f</td>  <td>%d</td><td>%d</td><td>%+.4f</td> </tr>\n", aSize, bSize, CountPercent(aSize, bSize), aSyms, bSyms, CountPercent(uint64(aSyms), uint64(bSyms)))
  }

  secTbl = fmt.Sprintf(
`
  <table>
    <tr><th></th><th>Section name</th><th>Size A</th><th>Size B</th><th>Diff</th>  <th>Symbols A</th><th>Symbols B</th><th>Diff</th> </tr>
%s  </table>
`, secTbl)
  return secTbl
}

/*func generateCompareSectionsTableHtml(cmp *Compare, A, B *file.FileInfo) string {
  tbl := fmt.Sprintf("    <tr><th>Address A</th><th>Address B</th><th>Name</th><th>Size</th></tr>\n")
  for symName, sym := range s.Symbols {
    tbl += fmt.Sprintf("    <tr><td>%x</td><td>%s</td><td>%x</td></tr>", sym.S.Value, symName, sym.S.Size)
  }
  tbl = fmt.Sprintf(
`
  <table>
%s
  </table>
`, tbl)
  return tbl
}*/

func generateSingleSectionsTableHtml(s *file.Section) string {
  tbl := fmt.Sprintf("    <tr><th>Address</th><th>Name</th><th>Size</th></tr>\n")
  for symName, sym := range s.Symbols {
    tbl += fmt.Sprintf("    <tr><td>%x</td><td>%s</td><td>%x</td></tr>", sym.S.Value, symName, sym.S.Size)
  }
  tbl = fmt.Sprintf(
`
  <table>
%s
  </table>
`, tbl)
  return tbl
}

func (r *Report) PrintHtml(cmp *Compare) {
  os.Mkdir("report", 0750)
  ind, err := os.Create("report/index.html")
  if err != nil {
    panic(err)
  }
  defer ind.Close()

  genTbl := generateGeneralInfoHtml(r.F1, r.F2)
  secTbl := generateSectionsTableHtml(cmp, r.F1, r.F2)
  str := index(
    r.F1.Name,
    r.F2.Name,
    genTbl,
    secTbl)
  ind.Write([]byte(str))

  os.Mkdir("report/sections", 0750)
  for _, secs := range cmp.Secs {
    for _, aSec := range secs.Asections {
      secName := aSec.Info.Name
      filePath := fmt.Sprintf("report/sections/%s.html", secName)
      ind, err := os.Create(filePath)
      if err != nil {
        panic(err)
      }
      defer ind.Close()

      str := singleSection(secName, r.F1.Name, generateSingleSectionsTableHtml(aSec))
      ind.Write([]byte(str))
    }
    for _, sec := range secs.ComonSections {
      secName := sec.A.Info.Name
      filePath := fmt.Sprintf("report/sections/%s.html", secName)
      ind, err := os.Create(filePath)
      if err != nil {
        panic(err)
      }
      defer ind.Close()

      str := compareSections(secName, r.F1.Name, r.F2.Name)
      ind.Write([]byte(str))
    }
    for _, bSec := range secs.Bsections {
      secName := bSec.Info.Name
      filePath := fmt.Sprintf("report/sections/%s.html", secName)
      ind, err := os.Create(filePath)
      if err != nil {
        panic(err)
      }
      defer ind.Close()

      str := singleSection(secName, r.F2.Name, generateSingleSectionsTableHtml(bSec))
      ind.Write([]byte(str))
    }
  }
}

