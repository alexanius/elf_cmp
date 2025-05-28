package compare

import (
  "bufio"
  "fmt"
  "os"
  "regexp"
  "strconv"
)

/*
export GODEBUG=gctrace=1

gc # @#s #%: #+...+# ms clock, #+...+# ms cpu, #->#-># MB, # MB goal, # P

gc #        the GC number, incremented at each GC
@#s         time in seconds since program start
#%          percentage of time spent in GC since program start
#+...+#     wall-clock/CPU times for the phases of the GC
#->#-># MB  heap size at GC start, at GC end, and live heap
# MB goal   goal heap size
# P         number of processors used
*/

type gc_cycle struct {
  Num       int64   // Number of GC cycle
  Time      float64 // Time from program start
  // percentage of time spent in GC since program start
  // wall-clock/CPU times for the phases of the GC
  HeapStart int64   // Heap size at GC start
  HeapEnd   int64   // Heap size at GC end
  HeapLive  int64   // Live heap
  // goal heap size
}

func AnalyzeLog(filePath string) error {
  file, err := os.Open(filePath)
  if err != nil {
    return err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }

  cycles := make([]gc_cycle, 0)
  gc_rx := regexp.MustCompile(`.*gc ([0-9]+) @([0-9]+\.[0-9]+)s.*, ([0-9]+)\-\>([0-9]+)\-\>([0-9]+)`)
  for _, l := range lines {
    if m := gc_rx.FindStringSubmatch(l) ; m != nil {
        m1, _ := strconv.ParseInt(m[1], 10, 64)
	m2, _ := strconv.ParseFloat(m[2], 64)
	m3, _ := strconv.ParseInt(m[3], 10, 64)
	m4, _ := strconv.ParseInt(m[4], 10, 64)
	m5, _ := strconv.ParseInt(m[5], 10, 64)

      cycle := gc_cycle {
        Num       : m1,
	Time      : m2,
	HeapStart : m3,
	HeapEnd   : m4,
	HeapLive  : m5,
      }
      cycles = append(cycles, cycle)
    }
  }

  for _, c := range cycles {
    fmt.Printf("%v\n", c)
  }

  return nil
}
