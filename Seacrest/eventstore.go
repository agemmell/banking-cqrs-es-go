package Seacrest

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"os"
	"time"
)

type EventEnvelope struct {
	EventID     string
	Order       uint
	AggregateID string
	EventType   string
	Payload     []byte
	RecordedAt  int64
}

type EventStore struct {
	orderedEvents []EventEnvelope                   // <global order> -> EventEnvelope
	eventsByID    map[string]map[uint]EventEnvelope // <aggregateID> -> <version> -> EventEnvelope
	globalOrder   uint
}

func NewEventStore() *EventStore {
	eventsByID := make(map[string]map[uint]EventEnvelope, 0)
	orderedEvents := make([]EventEnvelope, 0)
	es := EventStore{orderedEvents, eventsByID, 0}
	return &es
}

func (es *EventStore) GetAllEvents() []EventEnvelope {
	return es.orderedEvents
}

func (es *EventStore) PersistEvent(aggregateID string, eventType string, payload []byte) error {
	UUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	eventEnvelope := EventEnvelope{
		EventID:     UUID.String(),
		Order:       es.GlobalOrder() + 1,
		AggregateID: aggregateID,
		EventType:   eventType,
		Payload:     payload,
		RecordedAt:  time.Now().UnixNano(),
	}

	return es.PersistEventEnvelope(eventEnvelope)
}

func (es *EventStore) GetEventsByAggregateID(aggregateID string) map[uint]EventEnvelope {
	if aggregateEvents, ok := es.eventsByID[aggregateID]; ok {
		return aggregateEvents
	}
	return nil
}

func closeFileHandle(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func (es *EventStore) WriteEventsToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer closeFileHandle(f)

	for _, event := range es.orderedEvents {
		eventJson, err := json.Marshal(event)
		if err != nil {
			return err
		}
		_, err = f.WriteString(string(eventJson) + "\n")
		if err != nil {
			return err
		}
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (es *EventStore) LoadEventsFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFileHandle(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		eventEnvelope := EventEnvelope{}
		err := json.Unmarshal(scanner.Bytes(), &eventEnvelope)
		if err != nil {
			return err
		}
		err = es.PersistEventEnvelope(eventEnvelope)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (es *EventStore) GlobalOrder() uint {
	return es.globalOrder
}

func (es *EventStore) IncrementGlobalOrder() {
	es.globalOrder++
}

func (es *EventStore) PersistEventEnvelope(envelope EventEnvelope) error {
	version := uint(0)
	if _, ok := es.eventsByID[envelope.AggregateID]; ok {
		version = uint(len(es.eventsByID[envelope.AggregateID]))
		if _, ok := es.eventsByID[envelope.AggregateID][version]; ok {
			return errors.New(fmt.Sprintf("event version %d already exists for envelope %+v", version, envelope))
		}
		es.eventsByID[envelope.AggregateID][version] = envelope
	} else {
		es.eventsByID[envelope.AggregateID] = map[uint]EventEnvelope{
			0: envelope,
		}
	}
	es.orderedEvents = append(es.orderedEvents, envelope)
	es.IncrementGlobalOrder()

	return nil
}
