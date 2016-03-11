package main

import (
    "os"
    "fmt"
    "log"
    "strings"
    "database/sql"
    "encoding/json"
    _ "github.com/mattn/go-sqlite3"
)

type T interface{};
type Dict map[string]T;

func parse(rows *sql.Rows) []Dict {
    columns, _ := rows.Columns()    
    scanArgs := make([]T, len(columns))
    values   := make([]T, len(columns))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    // Map
    rmap := make(map[int]Dict);
    index := 0
    for rows.Next() {
        rows.Scan(scanArgs...)
        // Map K-V
        record := make(Dict)
        for i, col := range values {
            key := columns[i]
            if col != nil {
                switch t := col.(type) {
                    default:
                        fmt.Printf("\tUnexpected type %T\n", t)
                    case bool:
                        record[key] = col.(bool)
                    case int:
                        record[key] = col.(int)
                    case int8:
                        record[key] = col.(int8)
                    case int16:
                        record[key] = col.(int16)
                    case int32:
                        record[key] = col.(int32)
                    case int64:
                        record[key] = col.(int64)
                    case uint:
                        record[key] = col.(uint)
                    case uint8:
                        record[key] = col.(uint8)
                    case uint16:
                        record[key] = col.(uint16)
                    case uint32:
                        record[key] = col.(uint32)
                    case uint64:
                        record[key] = col.(uint64)
                    case float32:
                        record[key] = col.(float32)
                    case float64:
                        record[key] = col.(float64)
                    case []uint8:
                        record[key] = string(col.([]uint8))
                    case string:
                        record[key] = col.(string)
                }
            } else {
                record[key] = nil
            }
        }
        rmap[index] = record
        index += 1
    }
    // Array
    result := make([]Dict, len(rmap))
    for i := 0; i < len(rmap); i++ {
        result[i] = rmap[i]
    }
    return result
}

func query(db *sql.DB, stmt string) T{
    stmt   =  strings.TrimLeft(stmt, " ")
    // INSERT DELETE UPDATE SELECT
    action := strings.ToUpper(stmt[0:6])
    // log.Printf("Action: %s", action)

    // Exec, Prepare, Query
    if action == "CREATE" {
        _, err := db.Exec(stmt)
        if err != nil {
            return false;
        }
        return true
    } else if action == "INSERT" {
        res, err := db.Exec(stmt)
        if err != nil {
            return false;
        }
        id, _  := res.LastInsertId()
        return id
    } else if action == "UPDATE"{
        res, err := db.Exec(stmt)
        if err != nil {
            return false;
        }
        id, _  := res.LastInsertId()
        return id
    } else if action == "DELETE"{
        res, err := db.Exec(stmt)
        if err != nil {
            return false;
        }
        id, _  := res.LastInsertId()
        return id
    } else if action == "SELECT"{
        rows, err := db.Query(stmt)
        defer rows.Close()
        if err != nil {
            log.Printf("ERROR: 该条SQL（%s）语句查询失败.", stmt)
            return false
        }
        // Parse sql.Rows
        result := parse(rows)
        return result
    } else {
        log.Printf("ERROR: 该条SQL（%s）语句无法执行.", stmt)
        return false
    }
}

func execute(dbpath string, stmt string) bool {
    db, err := sql.Open("sqlite3", dbpath)
    if err != nil {
        log.Fatal(err)
        return false
    }
    defer db.Close()

    result := query(db, stmt)

    // Output
    switch t := result.(type) {
        default:
            fmt.Println("Type: ", t, "\tValue: ", result)
        case []Dict:
            jstring, _ := json.Marshal(result)
            fmt.Println( "Result(Dict): ", result)
            fmt.Println(string(jstring))
        case int64:
            fmt.Println( "LastInsertId: ", result)
        case bool:
            fmt.Println( "fail", result)
    }
    return true
}

func init_schema(db *sql.DB) bool {
    stmt := "CREATE TABLE auth (uid INTEGER NOT NULL PRIMARY KEY, name TEXT, nickname TEXT) "
    result := query(db, stmt)
    if result == true {
        return true
    } else {
        return false
    }
}

func init_db(name string) bool {
    path := "./"+name
    os.Remove(path)
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        log.Fatal(err)
        return false
    }
    defer db.Close()
    status := init_schema(db)
    return status
}

func main() {
    database := "test.sqlite.db"
    args := os.Args

    if len(args) <= 1 {
        return
    }
    command := args[1]
    switch command {
        default:
            log.Fatal("ERROR")
        case "init":
            // if len(args) > 2 {
            //     database = args[2]
            // }
            init_db(database)
        case "query":
            // "CREATE TABLE auth (uid INTEGER NOT NULL PRIMARY KEY, name TEXT, nickname TEXT) "
            // "INSERT INTO auth (uid, name, nickname) VALUES(1, '名字', '昵称') "
            // "UPDATE auth SET name = '名字2' WHERE uid=1 "
            // "SELECT * FROM auth"
            if len(args) > 2 {
                stmt := args[2]
                execute("./"+database, stmt)
            }
        case "excute":
            if len(args) > 2 {
                stmt := args[2]
                execute("./"+database, stmt)
            }

    }
}
