package main

import (
  "bufio"
  "encoding/csv"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "path/filepath"
  "sort"
  "strconv"
  "strings"
  "sync"
  "time"
)

var waitGroup sync.WaitGroup

//var dataMap map[string][]float64

func main() {
  //n:=time.Now()
  //fmt.Print(n.Format("2006-01-02 15:04:05.000000"))
  //dataMap = make(map[string][]float64)

  for _, fd := range dirWalk(".") {
    waitGroup.Add(1)
    go worker(fd)
  }

  waitGroup.Wait()
}

func dirWalk(dir string) []string {
  files, err := ioutil.ReadDir(dir)
  if err != nil {
    panic(err)
  }

  var paths []string
  for _, file := range files {
    if file.IsDir() {
      paths = append(paths, filepath.Join(dir, file.Name()))
    }
  }

  return paths
}

func worker(fd string) {
  dataMap := make(map[string][]float64)
  defer waitGroup.Done()
  //ioutil.WriteFile(fd+"/"+srcName, bs, 0644)

  file, err := os.OpenFile(fd+"/simulation.log", os.O_RDONLY, 0644)
  defer file.Close()
  if err != nil {
    fmt.Printf("ファイル読み込みにエラーが発生した：　%s\n", err.Error())
    return
  }

  bf := bufio.NewScanner(file)
  for bf.Scan() {
    line := bf.Text()
    if !strings.Contains(line, "REQUEST") {
      continue
    }
    cell := strings.Split(line, "\t")
    if len(cell) != 9 {
      fmt.Printf("len(cell)　is ： %d\n", len(cell))
      continue
    }
    if dataMap[cell[4]] == nil {
      dataMap[cell[4]] = []float64{}
    }
    i2, _ := strconv.ParseInt(cell[6], 10, 64)
    //t2 := time.Unix(i2/1000, 0)
    t2 := time.Unix(0, i2*int64(1000000))
    i1, _ := strconv.ParseInt(cell[5], 10, 64)
    //t1 := time.Unix(i1/1000, 0)
    t1 := time.Unix(0, i1*int64(1000000))
    //t2, err := time.Parse("2006-01-02 15:04:05.000000", cell[6])
    //if err != nil {
    //	fmt.Printf("時間読み込みにエラーが発生した：　%s\n", err.Error())
    //	continue
    //}
    //t1, err := time.Parse("2006-01-02 15:04:05.000000", cell[5])
    //if err != nil {
    //	fmt.Printf("時間読み込みにエラーが発生した：　%s\n", err.Error())
    //	continue
    //}

    dataMap[cell[4]] = append(dataMap[cell[4]], float64(t2.Sub(t1).Milliseconds()))
  }

  sb := strings.Builder{}
  // BOM付きCSVを作成
  sb.Write([]byte{0xEF, 0xBB, 0xBF})
  sb.WriteString(",平均値(ms),中央値(ms)\n")
  for _, v := range dataMap {
    sort.Float64s(v)
  }
  var dataSortedSlice [][]float64
  dataSortedSlice = append(dataSortedSlice, dataMap["init"])
  dataSortedSlice = append(dataSortedSlice, dataMap["init Redirect 1"])
  dataSortedSlice = append(dataSortedSlice, dataMap["inputMailAddress"])
  dataSortedSlice = append(dataSortedSlice, dataMap["inputMailAddress Redirect 1"])
  dataSortedSlice = append(dataSortedSlice, dataMap["inputMailAddress Redirect 2"])
  dataSortedSlice = append(dataSortedSlice, dataMap["idp_login_submit"])
  dataSortedSlice = append(dataSortedSlice, dataMap["saml_redirect_back"])
  for ind, val := range dataSortedSlice {
    switch ind {
    case 0:
      sb.WriteString("init")
    case 1:
      sb.WriteString("init Redirect 1")
    case 2:
      sb.WriteString("inputMailAddress")
    case 3:
      sb.WriteString("inputMailAddress Redirect 1")
    case 4:
      sb.WriteString("inputMailAddress Redirect 2")
    case 5:
      sb.WriteString("idp_login_submit")
    case 6:
      sb.WriteString("saml_redirect_back")
    }
    sb.WriteString(",")
    sb.WriteString(strconv.FormatFloat(average(val), 'f', 0, 64))
    sb.WriteString(",")
    sb.WriteString(strconv.FormatFloat(medianOf(val), 'f', 0, 64))
    sb.WriteString("\n")
  }
  //file, _ = os.OpenFile(fd+"/result.csv", os.O_CREATE|os.O_WRONLY, 0644)
  //defer file.Close()
  //csvWr := newCsvWriter(file,true)
  //csvWr.Write([]string{sb.String()})
  //csvWr.Flush()
  //for k, v := range dataMap {
  //	sort.Float64s(v)
  //  fmt.Printf("key is %s,value is %s\n",k,v)
  //	sb.WriteString(k)
  //	sb.WriteString(",")
  //	sb.WriteString(strconv.FormatFloat(average(v), 'f', 0, 64))
  //	sb.WriteString(",")
  //	sb.WriteString(strconv.FormatFloat(medianOf(v), 'f', 0, 64))
  //	sb.WriteString("\n")
  //}

  ioutil.WriteFile(fd+"/result.csv", []byte(sb.String()), 0644)
}

func newCsvWriter(w io.Writer, bom bool) *csv.Writer {
  bw := bufio.NewWriter(w)
  if bom {
    bw.Write([]byte{0xEF, 0xBB, 0xBF})
  }
  return csv.NewWriter(bw)
}

func average(xs []float64) (avg float64) {
  sum := 0.0
  switch len(xs) {
  case 0:
    avg = 0
  default:
    for _, val := range xs {
      sum += val
    }
    avg = sum / float64(len(xs))
  }
  return
}

func medianOf2(nums []float64) (median float64) {
  sort.Float64s(nums)
  l := len(nums)
  if l == 0 {
    panic("The length is zero!!!")
  }

  if l%2 == 0 {
    median = nums[l/2] + nums[l/2-1]/2.0
  } else {
    median = nums[l/2]
  }
  return
}

// medianOf はfloat配列から中央値を算出する。
func medianOf(ns []float64) float64 {
  l := len(ns)
  if l <= 0 {
    return 0.0
  }
  if l%2 == 1 {
    return ns[l/2]
  }
  return ns[l/2-1]
}