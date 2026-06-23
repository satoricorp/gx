CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    created_at INTEGER NOT NULL,
    ended_at INTEGER,
    command TEXT NOT NULL,
    cwd TEXT NOT NULL,
    client_pid INTEGER,
    exit_code INTEGER,
    gx_version TEXT NOT NULL,
    source TEXT,
    process_name TEXT,
    parent_pid INTEGER,
    last_seen_at INTEGER,
    end_reason TEXT,
    repo_root TEXT
);

CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL REFERENCES sessions(id),
    created_at INTEGER NOT NULL,
    provider TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    model TEXT,
    request_body BLOB NOT NULL,
    request_headers TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_requests_session ON requests(session_id);
CREATE INDEX IF NOT EXISTS idx_requests_created ON requests(created_at);

CREATE TABLE IF NOT EXISTS responses (
    id TEXT PRIMARY KEY,
    request_id TEXT NOT NULL REFERENCES requests(id),
    created_at INTEGER NOT NULL,
    completed_at INTEGER NOT NULL,
    status_code INTEGER NOT NULL,
    response_body BLOB NOT NULL,
    response_headers TEXT NOT NULL,
    is_streaming INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    provider_request_id TEXT,
    finish_reason TEXT,
    input_tokens INTEGER,
    output_tokens INTEGER,
    cache_read_tokens INTEGER,
    cache_write_tokens INTEGER,
    error TEXT
);

CREATE INDEX IF NOT EXISTS idx_responses_request ON responses(request_id);

CREATE TABLE IF NOT EXISTS repos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    root_path TEXT NOT NULL UNIQUE,
    backend TEXT NOT NULL,
    default_remote TEXT,
    default_branch TEXT,
    authoring_base_ref TEXT,
    remote_url TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS initialized_repos (
    root_path TEXT PRIMARY KEY,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS changes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL REFERENCES repos(id),
    jj_change_id TEXT NOT NULL,
    current_commit_id TEXT NOT NULL,
    description TEXT NOT NULL,
    parent_change_id TEXT,
    status TEXT NOT NULL,
    first_seen_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    UNIQUE(repo_id, jj_change_id)
);

CREATE INDEX IF NOT EXISTS idx_changes_repo ON changes(repo_id);

CREATE TABLE IF NOT EXISTS stacks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL REFERENCES repos(id),
    name TEXT NOT NULL,
    bookmark_name TEXT NOT NULL,
    base_ref TEXT NOT NULL,
    base_commit_id TEXT NOT NULL,
    head_change_id TEXT,
    remote_name TEXT,
    remote_ref TEXT,
    github_pr_url TEXT,
    status TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    UNIQUE(repo_id, bookmark_name)
);

CREATE INDEX IF NOT EXISTS idx_stacks_repo ON stacks(repo_id);
CREATE INDEX IF NOT EXISTS idx_stacks_head_change ON stacks(repo_id, head_change_id);

CREATE TABLE IF NOT EXISTS stack_changes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stack_id INTEGER NOT NULL REFERENCES stacks(id),
    change_id INTEGER NOT NULL REFERENCES changes(id),
    position INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    UNIQUE(stack_id, change_id)
);

CREATE INDEX IF NOT EXISTS idx_stack_changes_stack ON stack_changes(stack_id, position);
CREATE INDEX IF NOT EXISTS idx_stack_changes_change ON stack_changes(change_id);

CREATE TABLE IF NOT EXISTS change_revisions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    change_id INTEGER NOT NULL REFERENCES changes(id),
    jj_commit_id TEXT NOT NULL UNIQUE,
    jj_operation_id TEXT NOT NULL,
    changed_files_json TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_change_revisions_change ON change_revisions(change_id);

CREATE TABLE IF NOT EXISTS change_bookmarks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    change_id INTEGER NOT NULL REFERENCES changes(id),
    bookmark_name TEXT NOT NULL,
    remote_name TEXT,
    remote_ref TEXT,
    last_pushed_commit_id TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    UNIQUE(change_id, bookmark_name)
);

CREATE INDEX IF NOT EXISTS idx_change_bookmarks_change ON change_bookmarks(change_id);

CREATE TABLE IF NOT EXISTS change_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    change_id INTEGER NOT NULL REFERENCES changes(id),
    session_id TEXT NOT NULL REFERENCES sessions(id),
    created_at INTEGER NOT NULL,
    UNIQUE(change_id, session_id)
);

CREATE INDEX IF NOT EXISTS idx_change_sessions_change ON change_sessions(change_id);
CREATE INDEX IF NOT EXISTS idx_change_sessions_session ON change_sessions(session_id);

CREATE TABLE IF NOT EXISTS change_session_provenance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    change_id INTEGER NOT NULL REFERENCES changes(id),
    session_id TEXT NOT NULL REFERENCES sessions(id),
    agent_tool TEXT NOT NULL,
    provider TEXT NOT NULL,
    model_id TEXT NOT NULL,
    source TEXT,
    process_name TEXT,
    created_at INTEGER NOT NULL,
    UNIQUE(change_id, session_id, agent_tool, provider, model_id)
);

CREATE INDEX IF NOT EXISTS idx_change_session_provenance_change ON change_session_provenance(change_id);
CREATE INDEX IF NOT EXISTS idx_change_session_provenance_session ON change_session_provenance(session_id);

CREATE TABLE IF NOT EXISTS demux_proposals (
    id TEXT PRIMARY KEY,
    repo_id INTEGER NOT NULL REFERENCES repos(id),
    base_change_id TEXT NOT NULL,
    status TEXT NOT NULL,
    payload_json TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    applied_at INTEGER
);

CREATE INDEX IF NOT EXISTS idx_demux_proposals_repo ON demux_proposals(repo_id, status, created_at);

CREATE TABLE IF NOT EXISTS change_demux_evidence (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    change_id INTEGER NOT NULL REFERENCES changes(id),
    demux_proposal_id TEXT NOT NULL REFERENCES demux_proposals(id),
    revision_proposal_id TEXT NOT NULL,
    intent TEXT NOT NULL,
    files_json TEXT NOT NULL,
    hunk_ids_json TEXT NOT NULL,
    use_hunks INTEGER NOT NULL,
    confidence REAL NOT NULL,
    provenance_status TEXT NOT NULL,
    evidence_json TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    UNIQUE(change_id)
);

CREATE INDEX IF NOT EXISTS idx_change_demux_evidence_proposal ON change_demux_evidence(demux_proposal_id);

CREATE TABLE IF NOT EXISTS modify_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL REFERENCES repos(id),
    target_change_id INTEGER NOT NULL REFERENCES changes(id),
    previous_current_change_id INTEGER,
    jj_operation_id TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_modify_events_repo ON modify_events(repo_id);

CREATE TABLE IF NOT EXISTS pushes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL REFERENCES repos(id),
    remote_name TEXT,
    branch_name TEXT,
    head_commit_id TEXT NOT NULL,
    current_change_id INTEGER REFERENCES changes(id),
    created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pushes_repo ON pushes(repo_id);

CREATE TABLE IF NOT EXISTS cursor_messages (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL REFERENCES sessions(id),
    created_at INTEGER NOT NULL,
    role TEXT NOT NULL,
    text TEXT NOT NULL,
    raw_json BLOB NOT NULL,
    input_tokens INTEGER,
    output_tokens INTEGER
);

CREATE INDEX IF NOT EXISTS idx_cursor_messages_session ON cursor_messages(session_id);

CREATE TABLE IF NOT EXISTS capture_extracts (
    id TEXT PRIMARY KEY,
    repo_root TEXT NOT NULL,
    ref_range TEXT NOT NULL,
    payload_json BLOB NOT NULL,
    created_at INTEGER NOT NULL,
    uploaded_at INTEGER,
    upload_error TEXT
);

CREATE INDEX IF NOT EXISTS idx_capture_extracts_pending ON capture_extracts(uploaded_at, created_at);

CREATE TABLE IF NOT EXISTS capture_sessions (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    tool TEXT NOT NULL,
    payload_json BLOB NOT NULL,
    created_at INTEGER NOT NULL,
    uploaded_at INTEGER,
    upload_error TEXT
);

CREATE INDEX IF NOT EXISTS idx_capture_sessions_pending ON capture_sessions(uploaded_at, created_at);
