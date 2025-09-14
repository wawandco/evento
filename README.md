Evento is a proof of concept for a reservation system that handles concurrent reservation requests consistently. For the purpose of the test it may not have a web interface but consist of a http endpoints, however, it will be highly concurrent and be connected to a Postgres database.

### Statements
- Evento has a set of rooms available
- At any given time there might be more than one instance of Evento running
- There is only ONE instance of the database.
- Rooms are reserved concurrently
- Evento should not allow to reserve more than the rooms available
- Rooms are picked by users and locked while picked but not reserved
- Locked rooms go back to inventory after some time

### Database
At the database level we have a few tables that store the reservation data.

- events (id, name)
- hotels (id, name)
- event_hotel_rooms (id, event_id, hotel_id, assigned, reserved, locked)
- reservations (id, name, event_hotel_rooms_id, number)

### Objective
The objective of this POC is to validate that such concurrently consistent system is possible combining Go and Postgres. As a side product the repo will show the means required to achieve such consistency and allow us to do some benchmarking of the system.
