package models

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type log_t struct {
	perfs    *sync.Map
	mdy_flag int
	statf    *os.File
	warnf    *os.File
	logdir   string
	svrname  string
	logday   string
	warnstr  string
	level    int
	uids     []uint64
	ufmap    *sync.Map
	errcb    func(s string)
}

type perfstat_t struct {
	val       int64
	total_val int64
}

var (
	g_log        = &log_t{}
	g_stdout     = os.Stdout
	LogPrefixFn  = LogPrefixDay
	LogKeepFiles = 3
)

func init() {
	g_log.logdir = "./"
	g_log.perfs = new(sync.Map)
	g_log.ufmap = new(sync.Map)
}

func KeepLogFiles(fdir string, pattern string, nfile int) {
	logfs := []string{}
	filepath.Walk(fdir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			if f.Name() == filepath.Base(fdir) {
				return nil
			}
			return filepath.SkipDir
		}
		if ok, _ := filepath.Match(pattern, f.Name()); ok {
			logfs = append(logfs, f.Name())
		}
		return nil
	})

	sort.Slice(logfs, func(i, j int) bool {
		return strings.Compare(logfs[i], logfs[j]) < 0
	})

	for i := 0; i < len(logfs)-nfile; i++ {
		fn := filepath.Join(fdir, logfs[i])
		os.Remove(fn)
	}
}

func LogPrefixDay(t time.Time) (prefix string, patternstr string) {
	return t.Format("20060102"), "????????"
}

func LogPrefixHour(t time.Time) (prefix string, patternstr string) {
	return t.Format("2006010215"), "??????????"
}

func InitLog(logdir string, svrname string, warnstr string, fn func(s string), intervals ...int64) {
	g_log.logdir = logdir
	g_log.svrname = svrname
	g_log.warnstr = warnstr
	g_log.errcb = fn

	chgf := func(t time.Time) {
		tfmt, tpet := LogPrefixFn(t)
		if tfmt != g_log.logday {
			fname := filepath.Join(g_log.logdir, fmt.Sprintf("%s.%s.log", g_log.svrname, tfmt))
			f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return
			}
			g_log.logday = tfmt
			tmpf := g_stdout
			g_stdout = f
			if tmpf != os.Stdout {
				tmpf.Close() //最多就是出错
			}
			KeepLogFiles(g_log.logdir, fmt.Sprintf("%s.%s.log", g_log.svrname, tpet), LogKeepFiles)
		}
	}
	chgf(time.Now()) //即时创建
	fname := filepath.Join(g_log.logdir, "warn.log")
	g_log.warnf, _ = os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	interval := int64(60)
	if len(intervals) > 0 {
		interval = intervals[0]
	}
	if interval > 0 {
		//启动stat go
		go func() {
			for {
				tnow := time.Now()
				if tnow.Unix()%interval == 0 {
					logStat()
				}
				//切换g_stdout
				chgf(tnow)

				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func LogWrite(w io.Writer, level string, format string, v ...interface{}) {
	var nowTime string
	if level == "1" { //仅debug日志不加时间
		nowTime = time.Now().Format("15:04:05.000") //"2006-01-02 15:04:05.000")
	} else {
		nowTime = time.Now().Format("2006-01-02 15:04:05.000")
	}
	_, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)

	msg := format
	if len(v) > 0 {
		msg = fmt.Sprintf(format, v...)
	}
	fmt.Fprintf(w, "[%s]%s<%s>: %s:%d %s\r\n", nowTime, g_log.svrname, level, fileName, line, strings.TrimRight(msg, "\r\n"))
}

func LogD(format string, v ...interface{}) {
	if g_log.level > 1 {
		return
	}
	LogWrite(g_stdout, "1", format, v...)
}

func LogW(format string, v ...interface{}) {
	LogWrite(g_stdout, "4", format, v...)

	istr := fmt.Sprintf("\t%s\r\n", g_log.warnstr)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, fileName := path.Split(file)
		fn := runtime.FuncForPC(pc).Name()
		istr = istr + fmt.Sprintf("\t[-%d] %s(%s):%d\r\n", i, fileName, fn, line)
	}

	format = format + "\r\n%s"
	v = append(v, istr)

	LogWrite(g_log.warnf, "4", format, v...)

	if g_log.errcb != nil {
		sw := bytes.NewBufferString("")
		LogWrite(sw, "4", format, v...)
		go g_log.errcb(sw.String())
	}
}

func LogUid(uid uint64, format string, v ...interface{}) {
	if g_log.level <= 1 {
		LogWrite(g_stdout, fmt.Sprintf("u%d", uid), format, v...)
	}
	//save to u_xxx.log
	if InArray(uid, g_log.uids) {
		var f *os.File
		pf, loaded := g_log.ufmap.LoadOrStore(uid, &f)
		if loaded {
			ff, ok := pf.(**os.File)
			if !ok || *ff == nil {
				return
			}
			f = *ff
		} else {
			fname := filepath.Join(g_log.logdir, fmt.Sprintf("u_%d.log", uid))
			var err error
			f, err = os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				fmt.Printf("err=%v", err)
				g_log.ufmap.Delete(uid)
				return
			}
			//use f
		}

		LogWrite(f, fmt.Sprintf("u%d", uid), format, v...)
	}
}

func LogBool(cond bool, args ...interface{}) {
	format := "SUCC: "
	if !cond {
		format = "FATAL: "
	}
	if len(args) > 0 {
		format += args[0].(string)
		args = args[1:]
	}
	LogWrite(g_stdout, "1", format, args...)

	if !cond && g_log.errcb != nil {
		sw := bytes.NewBufferString("")
		LogWrite(sw, "4", format, args...)
		go g_log.errcb(sw.String())
	}
}

func LogError(err error, args ...interface{}) {
	format := "SUCC: "
	if err != nil {
		format = fmt.Sprintf("ERR %v", err)
	}
	if len(args) > 0 {
		format += args[0].(string)
		args = args[1:]
	}
	LogWrite(g_stdout, "1", format, args...)

	if err != nil && g_log.errcb != nil {
		sw := bytes.NewBufferString("")
		LogWrite(sw, "4", format, args...)
		go g_log.errcb(sw.String())
	}
}

func logStat() {
	if g_log.mdy_flag == 0 {
		return
	}

	var s string
	if g_log.mdy_flag == 2 {
		g_log.mdy_flag = 1

		s = "\r\n"
		g_log.perfs.Range(func(key, val interface{}) bool {
			if v, ok := val.(*perfstat_t); ok && (v.val != 0 || v.total_val != 0) {
				s = s + fmt.Sprintf("\t%v\t: %d / %d\r\n", key, v.val, v.total_val)
				v.val = 0
			}
			return true
		})
	} else {
		s = "="
	}
	if len(s) == 0 {
		return
	}
	//save to stat.log
	if g_log.statf == nil {
		fname := filepath.Join(g_log.logdir, "stat.log")
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return
		}
		g_log.statf = f
	}
	LogWrite(g_log.statf, "s", s)
}
