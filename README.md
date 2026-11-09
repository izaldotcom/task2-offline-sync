# Task 2: Offline-First Sync API

This is the implementation of the Server-Side Sync Endpoint for the offline-first clinic application.

### Tech Stack

- **Language:** Golang (Go 1.21)
- **Framework:** Echo (v4)
- **Database:** SQLite3
- **Architecture:** Controller-Service Pattern

---

## How to Run

The easiest way to run this API cleanly on port 8000 is by using Docker.

From the root of this `task2-offline-sync` directory, run the following command to build and start the server in the background:

```bash
docker compose up -d --build
```

The server will start instantly and listen on `http://localhost:8000`. The SQLite database will be automatically created inside the `data/` directory.

_(To stop the server later, simply run `docker compose down`)_

---

## How to Test

Once the server is running on port 8000, you can run the provided grading tests. The tests will first hit the `/reset` endpoint to seed initial clinics and stock, and then send the 7 sample payloads in sequence.

If you are on **Windows**, use the Python script:

```bash
python send_samples.py
```

If you are on **macOS/Linux** (or using Git Bash), you can use the shell script:

```bash
chmod +x send-samples.sh
./send-samples.sh
```

### Expected Output

- **Normal & Partial Batches:** Returns HTTP 200 with per-record statuses (`accepted`, `rejected`, `duplicate`).
- **Duplicate Batch:** Successfully identifies idempotency via `client_uuid` and marks duplicates.
- **Duplicate Patients:** Fuzzy matching flags patients without NIKs as `flagged_possible_duplicate_patient`.
- **Stock Conflicts:** Dispensing is always accepted, bringing the `MED-003` stock to negative, which is then flagged in the review queue.
