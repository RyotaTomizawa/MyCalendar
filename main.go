package main

import (
    "net/http"
    "html/template"
    "time"
    "bytes"
    "strconv"
    "strings"
    "fmt"
)
type carendarInfo struct {
    NameTest string
    RemainingDays string
    SelectedYear int
    SelectedMonth int
    CalendarRows [5]template.HTML
}
var ci = carendarInfo{
}
var selectedTime = time.Now()

var nameTest string
var remainingDays int

func main() {
    http.HandleFunc("/", mainHandler)
    http.HandleFunc("/changeMonth", changeMonth)
    http.HandleFunc("/settingTest", settingTest)
    http.ListenAndServe(":8000", nil)
}

//初期表示ハンドル
func mainHandler(w http.ResponseWriter, r *http.Request) {
    if ci.SelectedMonth == 0{
        var timeNow = time.Now()
        updateCalendarInfo(timeNow)
    }
    //index.htmlを表示する
    tpl, err := template.ParseFiles("doc/index.html")
    if err != nil {
        panic(err.Error())
    }
    if err := tpl.Execute(w,ci); err != nil {
        panic(err.Error())
    }
}

//月変更ハンドル
func changeMonth(w http.ResponseWriter, r *http.Request) {
    var t2 time.Time
    if r.FormValue("whichMonth") == "nextMonth"{
        t2 = time.Date(selectedTime.Year(), selectedTime.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
        selectedTime = t2
    }else{
        t2 = time.Date(selectedTime.Year(), selectedTime.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, -1, 0)
        selectedTime = t2
    }
    updateCalendarInfo(t2)
    //index.htmlを表示する
    tpl, err := template.ParseFiles("doc/index.html")
    if err != nil {
        panic(err.Error())
    }
    if err := tpl.Execute(w,ci); err != nil {
        panic(err.Error())
    }
}

//試験設定ハンドル
func settingTest(w http.ResponseWriter, r *http.Request) {   
    var layout = "2006/01/02"
    var str = r.FormValue("dateTest")
    var value = strings.Replace(str, "-", "/", -1)
    var dateTest, _ = time.Parse(layout, value)
    var sa = dateTest.Sub(time.Now())
    ci.RemainingDays = fmtDuration(sa)
    ci.NameTest = r.FormValue("nameTest")
    //index.htmlを表示する
    tpl, err := template.ParseFiles("doc/index.html")
    if err != nil {
        panic(err.Error())
    }
    if err := tpl.Execute(w,ci); err != nil {
        panic(err.Error())
    }
}

//試験までの時間を算出する
func fmtDuration(d time.Duration) string {
    d = d.Round(time.Minute)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    return fmt.Sprintf("%02d時間%02d分", h, m)
}

//カレンダー情報を更新する
func updateCalendarInfo(t time.Time){
    //選択された年を取得
    var selectedYear  = int(t.Year()) 
    //選択された月を取得
    var selectedMonth  = int(t.Month()) 
    //カレンダー行情報リスト
    var calendarRowList[5]template.HTML
    //月の始まり曜日を取得
    var firstDay = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
    var intFirstDay = int(firstDay.Weekday())
    //月の日数を取得
    var daysInMonth = daysInMonth(t)
    //カレンダー行情報リストを作成
    var date = 1
    for i := 0; i < 5; i++ {
        var bufferCalendarRow = bytes.NewBufferString("")
        for j := 0; j < 7; j++ {
          if  i == 0 && j < intFirstDay  {
              bufferCalendarRow.WriteString("<td></td>")
          } else if date > daysInMonth {
              break;
          } else {
              bufferCalendarRow.WriteString("<td>")
              var strDate = strconv.Itoa(date)
              bufferCalendarRow.WriteString(strDate)
              bufferCalendarRow.WriteString("</td>")
              date++;
          }
        }
        var strBufferCalendarRow = bufferCalendarRow.String()
        var calendarRow = template.HTML(strBufferCalendarRow)
        calendarRowList[i] = calendarRow
    }
    //カレンダー情報
    ci.SelectedYear = selectedYear
    ci.SelectedMonth = selectedMonth
    ci.CalendarRows = calendarRowList
}

//月の日数を算出する
func daysInMonth(t1 time.Time) int{
  var t2 = time.Date(t1.Year(), t1.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
  var d = t2.AddDate(0, 0, -1)
  return d.Day()
}

