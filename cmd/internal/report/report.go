package report

import (
  "fmt"
  "os"
  "strings"

  "github.com/jedib0t/go-pretty/v6/table"
  "github.com/jedib0t/go-pretty/v6/text"

  "elf_cmp/cmd/internal/compare"
  "elf_cmp/cmd/internal/file"
)

type secGroup int

const (
  DEBUG_GROUP secGroup = iota // Sections with debug information       
  INSTR_GROUP                 // Sections with executable instructions
  DATA_GROUP                  // Sections with user data
  GO_GROUP                    // Sections related to Go lang
  COMPILER_GROUP              // Sections with compiler data
  OTHER_GROUP                 // All other sections
)

func sectionGroup(name string) secGroup {
  if strings.Contains(name, ".debug") {
    return DEBUG_GROUP
  } else if name == ".data" ||
    name == ".bss" ||
    name == ".rodata" {
    return DATA_GROUP
  } else if name == ".typelink" ||
    name == ".gosymtab" ||
    name == ".noptrdata" ||
    name == ".gopclntab" ||
    name == ".noptrbss" ||
    name == ".itablink" {
    return GO_GROUP
  } else if name == ".note.go.buildid" ||
    name == ".go.buildinfo" ||
    name == ".note.gnu.property" ||
    name == ".note.ABI-tag" ||
    name == ".gnu.version" ||
    name == ".gnu.version_r" ||
    name == ".gnu.hash" ||
    name == ".gcc_except_table" {
    return COMPILER_GROUP
  }
  return OTHER_GROUP
}

func CountPercent(a, b uint64) float64 {
  if a == 0 || b == 0 {
    return 0
  }
  return float64(b) / float64(a)
}

func CountRatio(a, b uint64) string {
  if a == b {
    return "       ~"
  }

  if a == 0 || b == 0 {
    return "0"
  }

  r := float64(b) / float64(a)

  if r > 1 {
    return fmt.Sprintf("%.4f ^", r)
  } else if r < 1 {
    return fmt.Sprintf("%.4f V", r)
  }
  return "???"
}

func AddTextRow(name, A, B string, w table.Writer) {
  d := ""
  if A != B {
    d = "!"
  }
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  w.AppendRow([]interface{}{name, name, A, B, d}, rowConfigAutoMerge)
}

func AddIntRow(name string, A, B uint64, w table.Writer) {
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  w.AppendRow([]interface{}{
    name,
    name,
    fmt.Sprintf("%d", A),
    fmt.Sprintf("%d", B),
    fmt.Sprintf("%s", CountRatio(A, B))},
    rowConfigAutoMerge)
}

func AddIntRowGroup(group, name string, A, B uint64, w table.Writer) {
  w.AppendRow([]interface{}{
    group,
    name,
    fmt.Sprintf("%d", A),
    fmt.Sprintf("%d", B),
    fmt.Sprintf("%s", CountRatio(A, B))})
}

func AddIntRow1(name string, A uint64, w table.Writer) {
  w.AppendRow([]interface{}{
    name,
    name,
    fmt.Sprintf("%d", A),
    "",
    ""})
}

func AddIntRow1Group(group, name string, A uint64, w table.Writer) {
  w.AppendRow([]interface{}{
    group,
    name,
    fmt.Sprintf("%d", A),
    "",
    ""})
}

func AddIntRow2(name string, B uint64, w table.Writer) {
  w.AppendRow([]interface{}{
    name,
    name,
    "",
    fmt.Sprintf("%d", B),
    ""})
}

func AddIntRow2Group(group, name string, B uint64, w table.Writer) {
  w.AppendRow([]interface{}{
    group,
    name,
    "",
    fmt.Sprintf("%d", B),
    ""})
}

func AddSubtitle(name string, w table.Writer) {
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  w.AppendSeparator()
  w.AppendRow(table.Row{name, name, name, name, name}, rowConfigAutoMerge)
  w.AppendSeparator()
}

func AddSeparator(w table.Writer) {
  w.AppendSeparator()
}


func AddStatRow(name, A, B, d string, w table.Writer) {
  w.AppendRow([]interface{}{name, A, B, d})
}

func printHeader(cmp *compare.Compare, w table.Writer) {
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  w.AppendRow(table.Row{"A", "A", cmp.A.Name, cmp.A.Name, cmp.A.Name}, rowConfigAutoMerge)
  w.AppendRow(table.Row{"B", "B", cmp.B.Name, cmp.B.Name, cmp.B.Name}, rowConfigAutoMerge)
  w.AppendSeparator()
}

func printGeneralInfo(cmp *compare.Compare, w table.Writer) {
  rowConfigAutoMerge := table.RowConfig{
    AutoMerge:      true,
    AutoMergeAlign: text.AlignLeft}
  w.AppendRow(table.Row{"", "", "A", "B", "Diff (B/A)"}, rowConfigAutoMerge)
  w.AppendSeparator()
  w.AppendRow(table.Row{"General info", "General info",
    "General info", "General info", "General info"}, rowConfigAutoMerge)
  w.AppendSeparator()
  AddTextRow("Type", cmp.A.ElfType(), cmp.B.ElfType(), w)
  AddTextRow("Debug info", cmp.A.Dbg, cmp.B.Dbg, w)
  AddTextRow ("Sections",
    fmt.Sprintf("%d", cmp.A.SectionNum()),
    fmt.Sprintf("%d", cmp.B.SectionNum()), w)
  AddTextRow ("Symbols",
    fmt.Sprintf("%d", cmp.A.SymbolNum()),
    fmt.Sprintf("%d", cmp.B.SymbolNum()), w)
  AddIntRow ("Size", cmp.A.Size, cmp.B.Size, w)
  AddSubtitle("Sections size (bytes)", w)
}

func printSectionGroup(cmp *compare.Compare, gName string, w table.Writer) {
}

func Print(cmp *compare.Compare) {
  w := table.NewWriter()
  w.SetColumnConfigs([]table.ColumnConfig{
    {Number:    1,
     AutoMerge: true,
     Align:     text.AlignLeft,
     VAlign:    text.VAlignMiddle}})

  printHeader(cmp, w)
  printGeneralInfo(cmp, w)

//  printSectionGroup(cmp, A.InstrSec, B.InstrSec, "Instr", w)
//  printSectionGroup(cmp, A.UDataSec, B.UDataSec, "User data", w)
//  printSectionGroup(cmp, A.GoSec, B.GoSec, "Go data", w)
//  printSectionGroup(cmp, A.CompilerSec, B.CompilerSec, "Compiler data", w)
//  printSectionGroup(cmp, A.DebugSec, B.DebugSec, "Debug info", w)
//  printSectionGroup(cmp, A.OtherSec, B.OtherSec, "Other", w)

  AddSubtitle("Sections symbols number", w)

//  printSymbolGroup(cmp, A.InstrSec, B.InstrSec, "Instr", w)
//  printSymbolGroup(cmp, A.UDataSec, B.UDataSec, "User data", w)
//  printSymbolGroup(cmp, A.GoSec, B.GoSec, "Go data", w)
//  printSymbolGroup(cmp, A.CompilerSec, B.CompilerSec, "Compiler data", w)
//  printSymbolGroup(cmp, A.DebugSec, B.DebugSec, "Debug info", w)
//  printSymbolGroup(cmp, A.OtherSec, B.OtherSec, "Other", w)

  w.SetOutputMirror(os.Stdout)
  w.Render()
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

func generateSectionsTableHtml(cmp *compare.Compare, A, B *file.FileInfo) string {
/*  secTbl := "" // Table of sections
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
      secRow += fmt.Sprintf("    <tr><td>%s</td><td>%d</td><td></td><td></td>  <td>%d</td><td></td><td></td> </tr>\n", aSec.Name(), aSec.Info.Size, len(aSec.Symbols))
      aSize += aSec.Info.Size
      aSyms += len(aSec.Symbols)
    }
    for _, sec := range secs.ComonSections {
      aSymNum := len(sec.A.Symbols)
      bSymNum := len(sec.B.Symbols)
      secRow += fmt.Sprintf("    <tr><td>%s</td><td>%d</td><td>%d</td><td>%.4f</td>  <td>%d</td><td>%d</td><td>%.4f</td> </tr>\n", sec.A.Name(), sec.A.Info.Size, sec.B.Info.Size, CountPercent(sec.A.Info.Size, sec.B.Info.Size), aSymNum, bSymNum, CountPercent(uint64(aSymNum), uint64(bSymNum)))
      aSize += sec.A.Info.Size
      bSize += sec.B.Info.Size
      aSyms += aSymNum
      bSyms += bSymNum
    }
    for _, bSec := range secs.Bsections {
      secRow += fmt.Sprintf("    <tr><td>%s</td><td></td><td>%d</td><td></td>  <td></td><td>%d</td><td></td> </tr>\n", bSec.Name(), bSec.Info.Size, len(bSec.Symbols))
      bSize += bSec.Info.Size
      bSyms += len(bSec.Symbols)
    }
    secRows := len(secs.Asections) + len(secs.ComonSections) + len(secs.Bsections) + 2
    secRow = fmt.Sprintf(`%s
`, secRow)
    secTbl += fmt.Sprintf(`    <tr><th rowspan=%d>%s</th></tr>
%s`, secRows, gName, secRow)
    secTbl += fmt.Sprintf("    <tr><td>Total</td><td>%d</td><td>%d</td><td>%.4f</td>  <td>%d</td><td>%d</td><td>%.4f</td> </tr>\n", aSize, bSize, CountPercent(aSize, bSize), aSyms, bSyms, CountPercent(uint64(aSyms), uint64(bSyms)))
  }

  secTbl = fmt.Sprintf(
`
  <table>
    <tr><th></th><th>Section name</th><th>Size A</th><th>Size B</th><th>Diff</th>  <th>Symbols A</th><th>Symbols B</th><th>Diff</th> </tr>
%s  </table>
`, secTbl)
  return secTbl
*/
  return ""
}

func PrintHtml(cmp *compare.Compare) {
  return
/*  genTbl := generateGeneralInfoHtml(r.F1, r.F2)
  secTbl := generateSectionsTableHtml(cmp, r.F1, r.F2)
  str := index(
    r.F1.Name,
    r.F2.Name,
    genTbl,
    secTbl)

  os.Mkdir("report", 0750)
  ind, err := os.Create("report/index.html")
  if err != nil {
    panic(err)
  }
  defer ind.Close()

  ind.Write([]byte(str))*/
}

