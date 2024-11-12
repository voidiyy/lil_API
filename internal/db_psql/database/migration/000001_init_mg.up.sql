CREATE TABLE worker (
                       id uuid PRIMARY KEY UNIQUE NOT NULL,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       password varchar(150) NOT NULL,
                       email VARCHAR(50) UNIQUE NOT NULL,
                       role varchar(10) CHECK (role IN ('Junior', 'Middle', 'Senior')),
                        isActive bool default true,
                       created_at TIMESTAMPTZ DEFAULT NOW() --not serialize
);

CREATE TABLE projects (
                          id uuid PRIMARY KEY UNIQUE NOT NULL,
                          name VARCHAR(55) NOT NULL,
                          description TEXT,
                          start_date DATE,
                          end_date DATE,
                          team uuid[],
                          metadata JSONB, --serve users who take a part
                          created_by uuid REFERENCES worker(id),
                          created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tasks (
                       id uuid PRIMARY KEY UNIQUE NOT NULL,
                       project_id uuid REFERENCES projects(id) ON DELETE CASCADE,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       deadline DATE,
                       tags TEXT[],
                       status VARCHAR(50) CHECK (status IN ('To Do', 'In Progress', 'Done')), -- To Do, In Progress, Done
                       assigned_to uuid REFERENCES worker(id),
                       created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE comments (
                          id uuid PRIMARY KEY UNIQUE NOT NULL,
                          task_id uuid REFERENCES tasks(id) ON DELETE CASCADE,
                          author_id uuid REFERENCES worker(id),
                          content TEXT NOT NULL,
                          created_at TIMESTAMPTZ DEFAULT NOW()
);
