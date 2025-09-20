Evento is a proof of concept for a reservation system that handles concurrent reservation requests consistently. It consists of HTTP endpoints and is highly concurrent, connecting to a Postgres database.

### Objective
The objective of this POC is to validate that a concurrently consistent system is possible using Go and Postgres. As a side product, the repo demonstrates the means required to achieve such consistency and allows benchmarking of the system.

### Running the POC

Ensure you've cloned this repo and your current working directory is the root of it. To run Evento you need to have Go installed on your machine.
Then, at the root folder run Evento using Go with:
```
> go run ./cmd/ -clients 200 -mode naive -rooms 200
```

Where `200` is the number of concurrent clients to simulate, `naive` is the strategy to use, and `200` is the number of rooms per hotel. The strategies available are:
- naive: No concurrency control at all
- safe: Pessimistic locking using `SELECT ... FOR UPDATE`
- atomic: Transactional approach without locking
- optimistic: Optimistic locking using `updated_at` timestamp

Database connection parameters can be set using `DATABASE_URL` environment variable. By default it will connect to `postgres://postgres@localhost:5432/evento`.

### Statements
- Evento has a set of rooms available
- At any given time there might be more than one instance of Evento running
- There is only ONE instance of the database.
- Rooms are reserved concurrently
- Evento should NOT allow to reserve more than the rooms available

### Database
At the database level we have a few tables that store the reservation data.

- events (id, name)
- hotels (id, name)
- event_hotel_rooms (id, event_id, hotel_id, assigned, reserved, locked)
- reservations (id, name, event_hotel_rooms_id, number)

### Possible improvements / TODOs

- Locking rooms (part of the reservation)
- Better TUI, including progress.
