package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	// 连接ClickHouse数据库
	dsn := "https://default:QhH_vObgVEGw6@ci5eaxwoe9.asia-southeast1.gcp.clickhouse.cloud:8443/default?secure=true"

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to ClickHouse!")

	// 读取SQL初始化脚本
	sqlContent, err := ioutil.ReadFile("SQL/init.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// 读取并清理SQL内容
	sqlStr := string(sqlContent)

	// 移除注释行
	lines := strings.Split(sqlStr, "\n")
	var cleanStatements []string
	var currentStatement strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过注释行和空行
		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		currentStatement.WriteString(line)
		currentStatement.WriteString("\n")

		// 如果遇到分号，执行当前语句
		if strings.Contains(line, ";") {
			statement := strings.TrimSpace(currentStatement.String())
			if statement != "" {
				cleanStatements = append(cleanStatements, statement)
			}
			currentStatement.Reset()
		}
	}

	// 执行每个SQL语句
	for i, statement := range cleanStatements {
		fmt.Printf("Executing statement %d: %s\n", i+1, statement[:min(100, len(statement))])

		if _, err := db.Exec(statement); err != nil {
			log.Printf("Failed to execute statement %d: %v", i+1, err)
		} else {
			fmt.Printf("Statement %d executed successfully\n", i+1)
		}
	}

	fmt.Println("Database initialization completed!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
