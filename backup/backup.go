package backup

import (
    // "os"
    "fmt"
    "sync"
    "context"
    "database/sql"
    _ "github.com/lib/pq"
	"github.com/habx/pg-commands"

	"root/logger"
	"root/config"
    "root/connection"
    "root/utils"
)


func DumpDatabase(database *config.Database, databaseName string, folderPath string) (string, string, error) {
    sourceDB := &pgcommands.Postgres{
        Host:     database.Host,
        Port:     database.Port,
        DB:       databaseName,
        Username: database.User,
        Password: database.Password,
    }

    dump, err := pgcommands.NewDump(sourceDB)
    if err != nil {
        return "", "",err
    }
    
    logger.Info("Going to dump")

    dump.SetPath(folderPath)
    dumpExec := dump.Exec(pgcommands.ExecOptions{StreamPrint: false})
    if dumpExec.Error != nil {
        logger.Error("Error here")
        return "", "",dumpExec.Error.Err
    }


    return folderPath, dumpExec.File, nil
}


func RestoreDatabase(database *config.Database, databaseName string, filePath string) (error) {
    db, err := connection.GetDatabase(database)
    if err != nil {
        return err
    }
    _, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", databaseName))
    if err != nil {
        return err
    }

    _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", databaseName))
    if err != nil {
        return err
    }
    db.Close()

    restoreDB := &pgcommands.Postgres{
        Host:     database.Host,
        Port:     database.Port,
        DB:       databaseName,
        Username: database.User,
        Password: database.Password,
    }

    restore, err := pgcommands.NewRestore(restoreDB)
    if err != nil {
        return err
    }
    restore.Options = append(restore.Options, "--if-exists")
    restore.Role = database.User
    
    logger.Info("Going to restore")

    restoreExec := restore.Exec(filePath, pgcommands.ExecOptions{
        StreamPrint: false,
    })
    
    if restoreExec.Error != nil {
        logger.Info(restoreExec.Output)
        return restoreExec.Error.Err
    }

    return nil
}

func RemoteTransferBackup(user string, password string, host string, port int, sourcePath string, destinationPath string) (error) {
    sshClient, err := utils.SSHConnect(user, password, host, port)
    if err != nil {
        logger.Error("Error in ssh connection", err)
        return err
    }

    err = utils.FileTransfer(sshClient, sourcePath, destinationPath)
    if err != nil {
        logger.Error("Error in remote transport", err)
        return err
    }

    return nil
}

func DumpAndRestore(dumpDatabase *config.Database, restoreDatabase *config.Database, databaseName string, folderPath string) (error){
    folderPath, fileName, err := DumpDatabase(dumpDatabase, databaseName, folderPath)
    if err != nil {
        return err 
    }

    sourcePath := folderPath + fileName
    destinationPath := "/var/lib/postgresql/vs_test_backups/" + fileName

    err = RemoteTransferBackup("root", restoreDatabase.Password, restoreDatabase.Host, 22, sourcePath, destinationPath)
    if err != nil {
        return err 
    }
    
    // err = RestoreDatabase(restoreDatabase, databaseName ,folderPath + fileName)
    // if err != nil {
    //     return err
    // }

    err = CompareDBs(dumpDatabase, restoreDatabase, databaseName)
    if err != nil {
        return err
    }
    logger.Info("Successfully compared")


    return nil
}

func CompareDBs(sourceDatabase *config.Database, destinationDatabase *config.Database, databaseName string) (error) {
    sourceDatabase.Database = databaseName
    sourceDB, err := connection.GetDatabase(sourceDatabase)
    if err != nil {
        return err
    }
    defer sourceDB.Close()

    destinationDatabase.Database = databaseName
    destinationDB, err := connection.GetDatabase(destinationDatabase)
    if err != nil {
        return err
    }
    defer destinationDB.Close()

    tables := utils.GetTables(databaseName)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var wg sync.WaitGroup
    errCh := make(chan error, len(tables)) 

    for _, table := range tables {
        wg.Add(1)
        go func(table config.Table) {
            defer wg.Done()
            err := CompareTables(ctx, sourceDB, destinationDB, table)
            if err != nil {
                errCh <- err
                cancel()
            }
        }(table)
    }

    go func() {
        wg.Wait()
        close(errCh) 
    }()

    for err := range errCh {
        return fmt.Errorf(err.Error())
        break
    }

    return nil
}

func CompareTables(ctx context.Context, sourceDB *sql.DB, destinationDB *sql.DB, table config.Table) (error){
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        tableName := table.TableName
        orderBy := table.OrderBy

        query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s", tableName, orderBy)

        sourceRows, err := sourceDB.Query(query)
        if err != nil {
            return fmt.Errorf("error querying %s from sourceDB: %v", tableName, err)
        }
        defer sourceRows.Close()

        destinationRows, err := destinationDB.Query(query)
        if err != nil {
            return fmt.Errorf("error querying %s from destinationDB: %v", tableName, err)
        }
        defer destinationRows.Close()

        var sourceRowValues []interface{}
        columns, err := sourceRows.Columns()
        if err != nil {
            return fmt.Errorf("error getting columns for %s from sourceDB: %v", tableName, err)
        }

        for sourceRows.Next() {
            sourceRow := make([]interface{}, len(columns))
            sourceRowPtrs := make([]interface{}, len(columns))
            for i := range sourceRow {
                sourceRowPtrs[i] = &sourceRow[i]
            }
            if err := sourceRows.Scan(sourceRowPtrs...); err != nil {
                return fmt.Errorf("error scanning rows for %s from sourceDB: %v", tableName, err)
            }
            sourceRowValues = append(sourceRowValues, sourceRow)
        }

        var destinationRowValues []interface{}
        for destinationRows.Next() {
            destinationRow := make([]interface{}, len(columns))
            destinationRowPtrs := make([]interface{}, len(columns))
            for i := range destinationRow {
                destinationRowPtrs[i] = &destinationRow[i]
            }
            if err := destinationRows.Scan(destinationRowPtrs...); err != nil {
                return fmt.Errorf("error scanning rows for %s from destinationDB: %v", tableName, err)
            }
            destinationRowValues = append(destinationRowValues, destinationRow)
        }

        if len(sourceRowValues) != len(destinationRowValues) {
            return fmt.Errorf("%s has unequal rows", tableName)
        }

        for i := 0; i < len(sourceRowValues); i++ {
            sourceRow := sourceRowValues[i].([]interface{})
            destinationRow := destinationRowValues[i].([]interface{})

            for j := 0; j < len(sourceRow); j++ {
                column := columns[j]
                if utils.Contains(table.Skip, column) {
                    continue
                }

                if sourceRow[j] != destinationRow[j] {
                    return fmt.Errorf("%s column in %s is unequal for row %d -- %v | %v", column, tableName, i, sourceRow[j], destinationRow[j])
                }
            }
        }
    }
    return nil
}


