# Architecture Overview

## Project Structure
```
froggit/
├── main.go           # Application entry point
├── internal/         # Internal packages
│   ├── git/         # Git operations wrapper
│   └── tui/         # Terminal UI components
│       ├── model/   # Application state
│       ├── update/  # State updates
│       └── view/    # UI rendering
└── docs/            # Documentation
```

## Core Components

### Git Package
The `git` package provides a clean interface to Git operations, abstracting the complexity of Git commands and their execution.

### TUI Package
Built using Bubble Tea framework, the TUI package follows the Model-View-Update (MVU) architecture:

- **Model**: Maintains application state
- **Update**: Handles state transitions
- **View**: Renders the UI based on current state

## Data Flow
1. User input triggers messages
2. Messages are processed by Update functions
3. State changes are reflected in the Model
4. View renders updated state