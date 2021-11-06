#!/bin/bash
tern migrate -m /app/migrations/ --config /app/migrations/tern.conf
tern code install /app/migrations --config /app/migrations/tern.conf
/app/bin/bc_info --port 8080 --host 0.0.0.0