# Python Architect — Recipes

## Database access: SQL files + `importlib.resources` (§10)

```python
from importlib.resources import files
import psycopg
from psycopg.rows import class_row

GET_USER_BY_ID = files(__package__).joinpath("queries/get_user_by_id.sql").read_text()

class UserRepo:
    def __init__(self, conn: psycopg.AsyncConnection) -> None:
        self._conn = conn

    async def get_by_id(self, user_id: int) -> User | None:
        async with self._conn.cursor(row_factory=class_row(User)) as cur:
            await cur.execute(GET_USER_BY_ID, (user_id,))
            return await cur.fetchone()
```

Layout:

```
src/myapp/userrepo/
├── __init__.py
├── repo.py
└── queries/
    ├── get_user_by_id.sql
    ├── insert_user.sql
    └── list_users.sql
```
