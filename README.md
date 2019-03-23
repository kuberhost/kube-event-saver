Save kubernetes events to PostgreSQL

create table:
```sql
CREATE TABLE public.k8s_events (
    id uuid NOT NULL,
    namespace character varying NOT NULL,
    created_at timestamp with time zone NOT NULL,
    last_update_at timestamp with time zone NOT NULL,
    reason character varying NOT NULL,
    message character varying NOT NULL,
    body jsonb NOT NULL
);

CREATE UNIQUE INDEX events_id_idx ON public.k8s_events USING btree (id);
```

run kube-event-saver:
```sh
export PGUSER=user
export PGPASSWORD=xxxx
export PGDATABASE=cluster_data
export PGHOST=127.0.0.1
export KUBE_CONFIG=/some/path

./kube-event-saver
```

