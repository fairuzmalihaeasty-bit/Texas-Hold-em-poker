# Texas Hold'em (Monorepo)

Repository layout (recommended):

- backend/    -> Go REST API service (server, internal packages, Dockerfile)
- frontend/   -> Flutter client (flutter_app)
- README.md   -> Project overview and usage
- LICENSE     -> Project license (MIT by default)

This repository contains a Go backend (Texas Hold'em evaluator, comparator, simulator) and a Flutter frontend client. The current workspace contains the code under various folders; you can move backend source into `backend/` and the Flutter app into `frontend/` for a clearer monorepo layout.
