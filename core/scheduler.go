package core

import (
	"log"
	"scheduler-app/domain"
	"sync"
	"time"
)

// Event holds the task and the user ID to whom it belongs
type Event struct {
	UserID int
	Task   *domain.Task
}

// ClientConn wraps the user ID and the channel used to send data to the SSE stream
type ClientConn struct {
	Ch     chan *domain.Task
	UserID int
}

// Broker manages active client connections and broadcasts events to them
type Broker struct {
	mutex      sync.RWMutex
	clients    map[chan *domain.Task]int // map client channel to userID
	Register   chan ClientConn
	Unregister chan chan *domain.Task
	Broadcast  chan Event
}

func NewBroker() *Broker {
	return &Broker{
		clients:    make(map[chan *domain.Task]int),
		Register:   make(chan ClientConn),
		Unregister: make(chan chan *domain.Task),
		Broadcast:  make(chan Event, 100),
	}
}

func (b *Broker) Start() {
	log.Println("Event broker started")
	for {
		select {
		case conn := <-b.Register:
			b.mutex.Lock()
			b.clients[conn.Ch] = conn.UserID
			b.mutex.Unlock()
			log.Printf("SSE client registered for User ID %d. Active clients: %d", conn.UserID, len(b.clients))

		case ch := <-b.Unregister:
			b.mutex.Lock()
			if userID, ok := b.clients[ch]; ok {
				delete(b.clients, ch)
				close(ch)
				log.Printf("SSE client disconnected for User ID %d. Active clients: %d", userID, len(b.clients))
			}
			b.mutex.Unlock()

		case event := <-b.Broadcast:
			b.mutex.RLock()
			for ch, clientUserID := range b.clients {
				if clientUserID == event.UserID {
					// Non-blocking send in case a client buffer is full or blocked
					select {
					case ch <- event.Task:
					default:
						log.Printf("Warning: User %d client buffer full, dropping message", event.UserID)
					}
				}
			}
			b.mutex.RUnlock()
		}
	}
}

// TaskRepository interface defines the dependency requirements for the scheduler
type TaskRepository interface {
	GetPendingTasks() ([]*domain.Task, error)
	MarkTaskAsTriggered(taskID int) error
}

// Scheduler checks the database on a tick and triggers tasks when scheduled time passes
type Scheduler struct {
	repo   TaskRepository
	broker *Broker
	stop   chan struct{}
}

func NewScheduler(repo TaskRepository, broker *Broker) *Scheduler {
	return &Scheduler{
		repo:   repo,
		broker: broker,
		stop:   make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Second)
	log.Println("Background task scheduler started")
	go func() {
		for {
			select {
			case <-ticker.C:
				s.checkTasks()
			case <-s.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) checkTasks() {
	tasks, err := s.repo.GetPendingTasks()
	if err != nil {
		log.Printf("Scheduler error fetching pending tasks: %v", err)
		return
	}

	for _, task := range tasks {
		log.Printf("Triggering task: ID=%d, Title=%s, UserID=%d", task.ID, task.Title, task.UserID)

		// Mark as triggered in DB first to ensure single delivery
		if err := s.repo.MarkTaskAsTriggered(task.ID); err != nil {
			log.Printf("Scheduler error marking task %d as triggered: %v", task.ID, err)
			continue
		}

		// Broadcast event to active channels
		s.broker.Broadcast <- Event{
			UserID: task.UserID,
			Task:   task,
		}
	}
}
