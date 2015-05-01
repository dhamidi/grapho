package grapho

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type EventStore struct {
	storage EventStorage
}

type StoredEvent struct {
	Type       string
	RecordedAt time.Time
	Data       json.RawMessage

	event Event
}

func (self *StoredEvent) Event() Event {
	if self.event != nil {
		return self.event
	}

	self.event = NewEventForType(self.Type)
	if err := json.Unmarshal([]byte(self.Data), self.event); err != nil {
		panic(err)
	}
	return self.event
}

func NewEventStore(storage EventStorage) *EventStore {
	return &EventStore{storage}
}

func (self *EventStore) Store(events Events) error {
	tx := self.storage.Begin()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("EventStore.Store: %s", err)
			tx.Rollback()
		}
	}()

	for _, event := range events {
		tx.Add(event)
	}

	return tx.Commit()
}

func (self *EventStore) ReplayFunc(fn func(Event) error) error {
	view, err := self.storage.View()
	if err != nil {
		return err
	}
	defer view.Close()

	return view.ForEach(fn)
}

type EventStorage interface {
	Begin() EventTransaction
	View() (EventView, error)
}

type EventTransaction interface {
	Add(event Event)
	Commit() error
	Rollback()
}

type EventView interface {
	ForEach(fn func(Event) error) error
	Close() error
}

type EventsOnDisk struct {
	logFile *os.File
}

func NewEventsOnDisk(filename string) (EventStorage, error) {
	os.MkdirAll(filepath.Dir(filename), 0755)
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &EventsOnDisk{logFile: logFile}, nil
}

func (self *EventsOnDisk) Begin() EventTransaction {
	commitTo := self.logFile
	return NewDiskTransaction(commitTo)
}

func (self *EventsOnDisk) View() (EventView, error) {
	offset, err := self.logFile.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	return NewDiskView(offset, self.logFile.Name())
}

type DiskView struct {
	in        *os.File
	maxOffset int64
}

func NewDiskView(maxOffset int64, filename string) (*DiskView, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &DiskView{file, maxOffset}, nil
}

func (self *DiskView) currentOffset() (int64, error) {
	return self.in.Seek(0, 1)
}

func (self *DiskView) rewind() (int64, error) {
	return self.in.Seek(0, 0)
}

func (self *DiskView) Close() error {
	return self.in.Close()
}

func (self *DiskView) ForEach(fn func(Event) error) error {
	_, err := self.rewind()
	if err != nil {
		return err
	}

	dec := json.NewDecoder(io.LimitReader(self.in, self.maxOffset))

	for {
		envelope := StoredEvent{}
		err := dec.Decode(&envelope)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := fn(envelope.Event()); err != nil {
			return err
		}
	}

	return nil
}

type DiskTransaction struct {
	out      *os.File
	commitTo *os.File
	encoder  *json.Encoder
}

func NewDiskTransaction(commitTo *os.File) *DiskTransaction {
	out, err := ioutil.TempFile("", "tx-")
	if err != nil {
		panic(err)
	}

	return &DiskTransaction{out, commitTo, json.NewEncoder(out)}
}

func (self *DiskTransaction) remove() {
	os.Remove(self.out.Name())
}

func (self *DiskTransaction) write(event *StoredEvent) {
	if err := self.encoder.Encode(event); err != nil {
		panic(err)
	}
	self.out.Sync()
}

func (self *DiskTransaction) Add(event Event) {
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		panic(err)
	}

	envelope := &StoredEvent{
		Type:       event.EventType(),
		RecordedAt: time.Now(),
		Data:       json.RawMessage(data),
	}

	self.write(envelope)
}

func (self *DiskTransaction) Commit() error {
	defer self.remove()

	if _, err := self.out.Seek(0, 0); err != nil {
		return err
	}

	if _, err := io.Copy(self.commitTo, self.out); err != nil {
		return err
	}

	return self.commitTo.Sync()
}

func (self *DiskTransaction) Rollback() {
	self.remove()
}
