package handler

import (
	"database/sql"
	"fmt"
	"log"
	types "opencsd-api-server/src/type"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	insertInstanceUserQuery = "INSERT INTO workbench_user (user_name) VALUES (?)"
)

func getInstanceInfoDbDsn(instanceName string, operationNode string) string {
	operationNodeIp := types.ManagementMaster_.OperationLayer[operationNode].NodeIP
	dbName := strings.ReplaceAll(instanceName, "-", "_")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", types.INSTANCE_METRIC_MYSQL_USER, types.INSTANCE_METRIC_MYSQL_ROOT_PASSWORD, operationNodeIp, types.INSTANCE_METRIC_MYSQL_PORT, dbName)
	return dsn
}

func AddUser(instanceName string, userName string, operationNode string) error {
	dsn := getInstanceInfoDbDsn(instanceName, operationNode)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		time.Sleep(1)
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	_, err = db.Exec(insertInstanceUserQuery, userName)
	if err != nil {
		return fmt.Errorf("failed to insert record: %v", err)
	}

	log.Printf("Record with user_name='%s' inserted successfully into table\n", userName)
	return nil
}
