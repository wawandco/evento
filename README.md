Evento is a proof of concept for a reservation system that handles concurrent reservation requests consistently. For the purpose of the test it may not have a web interface but consist of a http endpoints, however, it will be highly concurrent and be connected to a Postgres database.

### Objective
The objective of this POC is to validate that such concurrently consistent system is possible combining Go and Postgres. As a side product the repo will show the means required to achieve such consistency and allow us to do some benchmarking of the system.

### Running the POC

Ensure you've cloned this repo and your current working directory is the root of it. To run Evento you need to have Go installed on your machine.
Then, at the root folder run Evento using using Go with:
```
> go run ./cmd/ 200 naive
```

Where `200` is the number of concurrent clients to simulate and `naive` is the strategy to use. The strategies available are:
- naive: No concurrency control at all
- pessimistic: Pessimistic locking using `SELECT ... FOR UPDATE`

Database connection parameters can be set using `DATABASE_URL` environment variable. By default it will connect to `postgres://postgres@localhost:5432/evento`.

### Statements
- Evento has a set of rooms available
- At any given time there might be more than one instance of Evento running
- There is only ONE instance of the database.
- Rooms are reserved concurrently
- Evento should NOT allow to reserve more than the rooms available
- TODO: Rooms are picked by users and locked while picked but not reserved
- TODO: Locked rooms go back to inventory after some time

### Database
At the database level we have a few tables that store the reservation data.

- events (id, name)
- hotels (id, name)
- event_hotel_rooms (id, event_id, hotel_id, assigned, reserved, locked)
- reservations (id, name, event_hotel_rooms_id, number)

### Possible improvements / TODOs

- More strategies (optimistic locking, etc)
- Locking rooms (part of the reservation)
- Data customization (events, hotels, rooms): Currently depends on database/data.sql
- Multiple servers runnning (multiple containers)
