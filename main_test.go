package main_test

import (
	"context"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedis(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ALLOW_EMPTY_PASSWORD": "true",
			"MYSQL_DATABASE":             "runner_admin_db_test",
		},
		Cmd: []string{
			"--general-log=ON",
			"--general-log-file=/var/log/mysql/query.log",
			"--slow-query-log=ON",
			"--slow-query-log-file=/var/log/mysql/slow.log",
			"--log-queries-not-using-indexes=ON",
			"--long-query-time=1",
			"--log-slow-extra=ON",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306  MySQL Community Server - GPL"),
			wait.ForListeningPort(nat.Port("3306/tcp")),
		),
	}
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := mysqlC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()
}
