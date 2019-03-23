package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx"
	_ "github.com/kr/pretty"
	"k8s.io/klog"
)

var pool *pgx.ConnPool

func SaveEvent(incEvent IncomingEvent) {
	checkPgPool()

	event_title := fmt.Sprintf("ns: %s\tmsg: %s\tid: %s",
		incEvent.Event.ObjectMeta.Namespace,
		incEvent.Event.Message,
		incEvent.Event.ObjectMeta.UID)

	klog.Infof("EVENT %s", event_title)

	_, err := pool.Exec("insertEvent",
		strings.ToUpper(fmt.Sprintf("%s", incEvent.Event.ObjectMeta.UID)),
		incEvent.Event.ObjectMeta.Namespace,
		incEvent.Event.ObjectMeta.CreationTimestamp.Time,
		incEvent.Event.LastTimestamp.Time,
		incEvent.Event.Reason,
		incEvent.Event.Message,
		incEvent.Event)

	//fmt.Printf("res %# v\n", pretty.Formatter(res))
	if err != nil {
		klog.Error("Can not insert event: ", err)
	}
}

func afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("insertEvent", `
        INSERT INTO k8s_events (id, namespace, created_at, last_update_at, reason, message, body)
        VALUES ($1::uuid, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            namespace = EXCLUDED.namespace,
            created_at = EXCLUDED.created_at,
            last_update_at = EXCLUDED.last_update_at,
            reason = EXCLUDED.reason,
            message = EXCLUDED.message,
            body = EXCLUDED.body
        RETURNING id`)

	if err != nil {
		klog.Error("Unable to parse insertEvent query:", err)
	}
	return err
}

func checkPgPool() {
	if pool != nil {
		return
	}

	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		klog.Error("Unable to parse environment:", err)
		os.Exit(1)
	}

	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 5,
		AfterConnect:   afterConnect,
	})
	if err != nil {
		klog.Error("Unable to connect to database:", err)
		os.Exit(1)
	}

}
