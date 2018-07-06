
import (
	"sync"
	"time"
)

// NewDedupe wraps an existing logging Interface and returns one which dedupe
// log lines.
func NewDedupe(level Level, interval time.Duration, inner Interface) Interface {
	return &dedupe{
		inner:    inner,
		level:    level,
		interval: interval,
		seen:     map[string]entryCount{},
	}
}

type dedupe struct {
	inner    Interface
	level    Level
	interval time.Duration
	seen     map[string]entryCount
	lock     sync.Mutex
}

/*func (d *dedupe) shouldLog(level Level, entry Entry) bool {
	if _, ok := entry.Data["deduplicated"]; ok {
		// ignore our own logs about deduped messages
		return true
	}
	if level.logrus < d.level.logrus {
		// ignore logs more severe than our level
		return true
	}
	key := fieldsToString(entry)
	d.lock.Lock()
	defer d.lock.Unlock()
	if ec, ok := d.seen[key]; ok {
		// already seen, increment count and do not log
		ec.count++
		d.seen[key] = ec
		return false
	}
	// New message, log it but add it to seen.
	// We need to copy because the pointer ceases to be valid after we return from Format
	d.seen[key] = entryCount{entry: *entry}
	go d.evictEntry(key) // queue to evict later
	return true
}

// Wait for interval seconds then evict the entry and send the log
func (d *dedupe) evictEntry(key string) {
	time.Sleep(d.interval)
	var ec entryCount
	func() {
		d.lock.Lock()
		defer d.lock.Unlock()
		ec = d.seen[key]
		delete(d.seen, key)
	}()
	if ec.count == 0 {
		return
	}
	entry := d.inner.WithFields(ec.entry.Data).WithField("deduplicated", ec.count)
	message := fmt.Sprintf("Repeated %d times: %s", ec.count, ec.entry.Message)
	// There's no way to choose the log level dynamically, so we have to do this hack
	map[log.Level]func(args ...interface{}){
		log.PanicLevel: entry.Panic,
		log.FatalLevel: entry.Fatal,
		log.ErrorLevel: entry.Error,
		log.WarnLevel:  entry.Warn,
		log.InfoLevel:  entry.Info,
		log.DebugLevel: entry.Debug,
	}[ec.entry.Level](message)
}*/

func (d *dedupe) Debugf(format string, args ...interface{}) {
	d.inner.Debugf(format, args...)
}
func (d *dedupe) Debugln(args ...interface{}) {
	d.inner.Debugln(args...)
}

func (d *dedupe) Infof(format string, args ...interface{}) {
	d.inner.Infof(format, args...)
}
func (d *dedupe) Infoln(args ...interface{}) {
	d.inner.Infoln(args...)
}

func (d *dedupe) Warnf(format string, args ...interface{}) {
	d.inner.Warnf(format, args...)
}
func (d *dedupe) Warnln(args ...interface{}) {
	d.inner.Warnln(args...)
}

func (d *dedupe) Errorf(format string, args ...interface{}) {
	d.inner.Errorf(format, args...)
}
func (d *dedupe) Errorln(args ...interface{}) {
	d.inner.Errorln(args...)
}

func (d *dedupe) WithField(key string, value interface{}) Interface {
	return d.inner.WithField(key, value)
}

func (d *dedupe) WithFields(fields Fields) Interface {
	return d.inner.WithFields(fields)
}
