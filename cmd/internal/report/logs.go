package report

import (
  "fmt"
  "os"

  "elf_cmp/cmd/internal/compare"
)

func GcTraceReport(cycles []compare.GcCycle) error {

  heapSize := make([]int64, 0)
  heapLive := make([]int64, 0)
  gcTimes  := make([]float64, 0)
  for _, c := range cycles {
    heapSize = append(heapSize, c.HeapStart)
    heapSize = append(heapSize, c.HeapEnd)
    heapLive = append(heapLive, c.HeapLive)
    gcTimes = append(gcTimes, c.Time)
  }

  line1 := ""
  line2 := ""
  maxMemory := int64(0)
  for i := 0; i < len(heapSize) ; i+=2 {
    line1 += fmt.Sprintf("{val:%d,cycle:%d},\n", heapSize[i], i/2)
    line1 += fmt.Sprintf("{val:%d,cycle:%d},\n", heapSize[i+1], i/2)
    line2 += fmt.Sprintf("{val:%d,cycle:%d},\n", heapLive[i/2], i/2)
    if heapSize[i] > maxMemory {
      maxMemory = heapSize[i]
    }
    if heapSize[i+1] > maxMemory {
      maxMemory = heapSize[i+1]
    }
  }
  page := fmt.Sprintf(`
<!DOCTYPE html>
<div id="container"></div>
<div id="ca"></div>

<script src="d3.js"></script>
<script type="module">
const data1 = Object.values([
    [%s],
    [%s]
]);

var line = d3.line()
  .x((d, i) => x(d.cycle))
  .y((d) => y(d.val))

// set the dimensions and margins of the graph
var margin = {
    top: 50,
    right: 100,
    bottom: 130,
    left: 120
  },
  width = 900 - margin.left - margin.right,
  height = 400 - margin.top - margin.bottom;

// append the svg object to the body of the page
var svg = d3.select("#ca")
  .append("svg")
  .attr("width", width + margin.left + margin.right)
  .attr("height", height + margin.top + margin.bottom)
  .append("g")
  .attr("transform", ` + "`translate(${margin.left}, ${margin.top})`" + `);

// Add X axis
var x = d3.scaleLinear()
  .domain([0, d3.max(data1, (d) => %d)])
  .range([0, width]);

svg.append("g")
  .attr("transform", "translate(0," + height + ")")
  .call(d3.axisBottom(x).ticks(5));

// Add Y axis
var y = d3.scaleLinear()
  .domain([0, d3.max(data1, (d) => %d)])
  .range([height, 0]);

svg.append("g")
  .call(d3.axisLeft(y));

// Draw the line
svg.selectAll(".line")
  .data(data1)
  .enter()
  .append("path")
  .attr("fill", "none")
  .attr("stroke", "black")
  .attr("stroke-width", 1.5)
  .attr("d", (d) => line(d));

</script>
`, line1, line2, len(heapLive) + 10, maxMemory)

  os.Mkdir("report", 0750)
  ind, err := os.Create("report/index.html")
  if err != nil {
    panic(err)
  }
  defer ind.Close()

  ind.Write([]byte(page))

  return nil
}
