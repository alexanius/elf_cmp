package report

import (
  "fmt"
  "os"

  "github.com/jedib0t/go-pretty/v6/table"
  "github.com/jedib0t/go-pretty/v6/text"
)

func CountPercent(a, b uint64) float64 {
  if a > b {
    return -(1 - float64(b) / float64(a)) * 100
  }
  return (float64(b) / float64(a)) * 100
}

type Report struct {
  Stat      table.Writer
}

func New(A, B string) *Report{
  r := Report{}

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
  r.Stat.AppendRow(table.Row{"A", "A", A, A, A}, rowConfigAutoMerge)
  r.Stat.AppendRow(table.Row{"B", "B", B, B, B}, rowConfigAutoMerge)
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
