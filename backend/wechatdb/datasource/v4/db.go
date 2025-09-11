package v4

import (
	"database/sql"
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"wechat-secretary/backend/pkg/filecopy"
	"wechat-secretary/backend/pkg/filemonitor"
)

type MessageDBInfo struct {
	FilePath  string
	StartTime time.Time
	EndTime   time.Time
}

type Db interface {
	InitDb(path string, group *Group, ch chan struct{}) (*filemonitor.FileGroup, error)
	GetDb(start, end time.Time) []*sql.DB
	Close()
}

type ContactDb struct {
	id       string
	fg       *filemonitor.FileGroup
	mutex    sync.RWMutex
	paths    []string
	dbs      []*sql.DB
	dbOpen   bool
	changeCh chan struct{}
}

func (c *ContactDb) Close() {
	_ = c.dbs[0].Close()
}

func (c *ContactDb) InitDb(path string, group *Group, ch chan struct{}) (*filemonitor.FileGroup, error) {
	fg, err := filemonitor.NewFileGroup(group.Name, path, group.Pattern, group.BlackList)
	if err != nil {
		return nil, err
	}
	fg.AddCallback(c.callback)
	c.changeCh = ch
	c.id = filepath.Base(path)
	c.fg = fg
	list, err := fg.List()
	if err != nil {
		log.Errorf("Wechat db file not exist %s:%s", path, group.Pattern)
	}
	if len(list) != 0 {
		c.paths = list
	}
	return fg, nil
}

func (c *ContactDb) GetDb(start, end time.Time) []*sql.DB {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.dbOpen {
		return c.dbs
	} else {
		paths := c.getPaths()

		dbs := make([]*sql.DB, 0)
		for _, path := range paths {
			if db, err := c.openDb(path); err != nil {
				log.Errorf("Failed to open contact db,err:%s", err)
			} else {
				dbs = append(dbs, db)
			}
		}
		c.dbs = dbs
		return c.dbs
	}
}

func (c *ContactDb) openDb(path string) (*sql.DB, error) {
	var err error
	tempPath := path
	if runtime.GOOS == "windows" {
		tempPath, err = filecopy.GetTempCopy(c.id, path)
		if err != nil {
			log.Errorf("Failed to get contact db copy temp,err:%s", err)
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", tempPath)
	if err != nil {
		log.Errorf("Failed to connect contact db,err:%s", err)
		return nil, err
	}
	return db, nil
}

func (c *ContactDb) getPaths() []string {
	if len(c.paths) != 0 {
		return c.paths
	}
	list, err := c.fg.List()
	if err != nil {
		log.Errorf("Failed to get contact db path,err:%s", err)
	}
	return list
}

func (c *ContactDb) callback(event fsnotify.Event) error {
	if !event.Op.Has(fsnotify.Create) {
		return nil
	}
	c.mutex.Lock()
	if c.dbOpen {
		c.dbOpen = false
	}
	_ = c.dbs[0].Close()

	c.changeCh <- struct{}{}
	c.mutex.Unlock()
	return nil
}

type MessageDb struct {
}

func (m MessageDb) Close() {
	// TODO implement me
	panic("implement me")
}

func (m MessageDb) InitDb(path string, group *Group, ch chan struct{}) (*filemonitor.FileGroup, error) {
	// TODO implement me
	panic("implement me")
}

func (m MessageDb) GetDb(start, end time.Time) []*sql.DB {
	// TODO implement me
	panic("implement me")
}

// func (ds *DataSource) initMessageDbs() error {
// 	dbPaths, err := ds.db.GetDBPath(Message)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "Wechat db file not exist") {
// 			ds.messageInfos = make([]MessageDBInfo, 0)
// 			return nil
// 		}
// 		return err
// 	}
//
// 	// 处理每个数据库文件
// 	infos := make([]MessageDBInfo, 0)
// 	for _, filePath := range dbPaths {
// 		db, err := ds.db.OpenDB(filePath)
// 		if err != nil {
// 			log.Errorf("Failed to open message DB path:%s", filePath)
// 			continue
// 		}
//
// 		// 获取 Timestamp 表中的开始时间
// 		var startTime time.Time
// 		var timestamp int64
//
// 		row := db.QueryRow("SELECT timestamp FROM Timestamp LIMIT 1")
// 		if err := row.Scan(&timestamp); err != nil {
// 			log.Errorf("Failed to get message DB timestamp:%s", filePath)
// 			continue
// 		}
// 		startTime = time.Unix(timestamp, 0)
//
// 		// 保存数据库信息
// 		infos = append(infos, MessageDBInfo{
// 			FilePath:  filePath,
// 			StartTime: startTime,
// 		})
// 	}
//
// 	// 按照 StartTime 排序数据库文件
// 	sort.Slice(infos, func(i, j int) bool {
// 		return infos[i].StartTime.Before(infos[j].StartTime)
// 	})
//
// 	// 设置结束时间
// 	for i := range infos {
// 		if i == len(infos)-1 {
// 			infos[i].EndTime = time.Now()
// 		} else {
// 			infos[i].EndTime = infos[i+1].StartTime
// 		}
// 	}
// 	if len(ds.messageInfos) > 0 && len(infos) < len(ds.messageInfos) {
// 		log.Warnf("message db count decreased from %d to %d, skip init", len(ds.messageInfos), len(infos))
// 		return nil
// 	}
// 	ds.messageInfos = infos
// 	return nil
// }

type SessionDb struct {
	id       string
	fg       *filemonitor.FileGroup
	mutex    sync.RWMutex
	paths    []string
	dbs      []*sql.DB
	dbOpen   bool
	changeCh chan struct{}
}

func (s *SessionDb) Close() {
	_ = s.dbs[0].Close()
}

func (s *SessionDb) InitDb(path string, group *Group, ch chan struct{}) (*filemonitor.FileGroup, error) {
	fg, err := filemonitor.NewFileGroup(group.Name, path, group.Pattern, group.BlackList)
	if err != nil {
		return nil, err
	}
	fg.AddCallback(s.callback)
	s.changeCh = ch
	s.id = filepath.Base(path)
	s.fg = fg
	list, err := fg.List()
	if err != nil {
		log.Errorf("Wechat db file not exist %s:%s", path, group.Pattern)
	}
	if len(list) != 0 {
		s.paths = list
	}
	return fg, nil
}

func (s *SessionDb) GetDb(start, end time.Time) []*sql.DB {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.dbOpen {
		return s.dbs
	} else {
		paths := s.getPaths()

		dbs := make([]*sql.DB, 0)
		for _, path := range paths {
			if db, err := s.openDb(path); err != nil {
				log.Errorf("Failed to open contact db,err:%s", err)
			} else {
				dbs = append(dbs, db)
			}
		}
		s.dbs = dbs
		return s.dbs
	}
}

func (s *SessionDb) openDb(path string) (*sql.DB, error) {
	var err error
	tempPath := path
	if runtime.GOOS == "windows" {
		tempPath, err = filecopy.GetTempCopy(s.id, path)
		if err != nil {
			log.Errorf("Failed to get contact db copy temp,err:%s", err)
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", tempPath)
	if err != nil {
		log.Errorf("Failed to connect contact db,err:%s", err)
		return nil, err
	}
	return db, nil
}

func (s *SessionDb) getPaths() []string {
	if len(s.paths) != 0 {
		return s.paths
	}
	list, err := s.fg.List()
	if err != nil {
		log.Errorf("Failed to get contact db path,err:%s", err)
	}
	return list
}

func (s *SessionDb) callback(event fsnotify.Event) error {
	if !event.Op.Has(fsnotify.Create) {
		return nil
	}
	s.mutex.Lock()
	if s.dbOpen {
		s.dbOpen = false
	}
	_ = s.dbs[0].Close()

	s.changeCh <- struct{}{}
	s.mutex.Unlock()
	return nil
}

type MediaDb struct {
	id       string
	fg       *filemonitor.FileGroup
	mutex    sync.RWMutex
	paths    []string
	dbs      []*sql.DB
	dbOpen   bool
	changeCh chan struct{}
}

func (s *MediaDb) Close() {
	_ = s.dbs[0].Close()
}

func (s *MediaDb) InitDb(path string, group *Group, ch chan struct{}) (*filemonitor.FileGroup, error) {
	fg, err := filemonitor.NewFileGroup(group.Name, path, group.Pattern, group.BlackList)
	if err != nil {
		return nil, err
	}
	fg.AddCallback(s.callback)
	s.changeCh = ch
	s.id = filepath.Base(path)
	s.fg = fg
	list, err := fg.List()
	if err != nil {
		log.Errorf("Wechat db file not exist %s:%s", path, group.Pattern)
	}
	if len(list) != 0 {
		s.paths = list
	}
	return fg, nil
}

func (s *MediaDb) GetDb(start, end time.Time) []*sql.DB {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.dbOpen {
		return s.dbs
	} else {
		paths := s.getPaths()

		dbs := make([]*sql.DB, 0)
		for _, path := range paths {
			if db, err := s.openDb(path); err != nil {
				log.Errorf("Failed to open contact db,err:%s", err)
			} else {
				dbs = append(dbs, db)
			}
		}
		s.dbs = dbs
		return s.dbs
	}
}

func (s *MediaDb) openDb(path string) (*sql.DB, error) {
	var err error
	tempPath := path
	if runtime.GOOS == "windows" {
		tempPath, err = filecopy.GetTempCopy(s.id, path)
		if err != nil {
			log.Errorf("Failed to get contact db copy temp,err:%s", err)
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", tempPath)
	if err != nil {
		log.Errorf("Failed to connect contact db,err:%s", err)
		return nil, err
	}
	return db, nil
}

func (s *MediaDb) getPaths() []string {
	if len(s.paths) != 0 {
		return s.paths
	}
	list, err := s.fg.List()
	if err != nil {
		log.Errorf("Failed to get contact db path,err:%s", err)
	}
	return list
}

func (s *MediaDb) callback(event fsnotify.Event) error {
	if !event.Op.Has(fsnotify.Create) {
		return nil
	}
	s.mutex.Lock()
	if s.dbOpen {
		s.dbOpen = false
	}
	_ = s.dbs[0].Close()

	s.changeCh <- struct{}{}
	s.mutex.Unlock()
	return nil
}

type VoiceDb struct {
	id       string
	fg       *filemonitor.FileGroup
	mutex    sync.RWMutex
	paths    []string
	dbs      []*sql.DB
	dbOpen   bool
	changeCh chan struct{}
}

func (s *VoiceDb) Close() {
	for _, db := range s.dbs {
		_ = db.Close()
	}
}

func (s *VoiceDb) InitDb(path string, group *Group, ch chan struct{}) (*filemonitor.FileGroup, error) {
	fg, err := filemonitor.NewFileGroup(group.Name, path, group.Pattern, group.BlackList)
	if err != nil {
		return nil, err
	}
	fg.AddCallback(s.callback)
	s.changeCh = ch
	s.id = filepath.Base(path)
	s.fg = fg
	list, err := fg.List()
	if err != nil {
		log.Errorf("Wechat db file not exist %s:%s", path, group.Pattern)
	}
	if len(list) != 0 {
		s.paths = list
	}
	return fg, nil
}

func (s *VoiceDb) GetDb(start, end time.Time) []*sql.DB {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.dbOpen {
		return s.dbs
	} else {
		paths := s.getPaths()

		dbs := make([]*sql.DB, 0)
		for _, path := range paths {
			if db, err := s.openDb(path); err != nil {
				log.Errorf("Failed to open contact db,err:%s", err)
			} else {
				dbs = append(dbs, db)
			}
		}
		s.dbs = dbs
		return s.dbs
	}
}

func (s *VoiceDb) openDb(path string) (*sql.DB, error) {
	var err error
	tempPath := path
	if runtime.GOOS == "windows" {
		tempPath, err = filecopy.GetTempCopy(s.id, path)
		if err != nil {
			log.Errorf("Failed to get contact db copy temp,err:%s", err)
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", tempPath)
	if err != nil {
		log.Errorf("Failed to connect contact db,err:%s", err)
		return nil, err
	}
	return db, nil
}

func (s *VoiceDb) getPaths() []string {
	if len(s.paths) != 0 {
		return s.paths
	}
	list, err := s.fg.List()
	if err != nil {
		log.Errorf("Failed to get contact db path,err:%s", err)
	}
	return list
}

func (s *VoiceDb) callback(event fsnotify.Event) error {
	if !event.Op.Has(fsnotify.Create) {
		return nil
	}
	s.mutex.Lock()
	if s.dbOpen {
		s.dbOpen = false
	}
	_ = s.dbs[0].Close()

	s.changeCh <- struct{}{}
	s.mutex.Unlock()
	return nil
}
