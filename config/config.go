package config

import (
    "encoding/json"
    "os"
)

type LoggerConfig struct {
    LogLevel		string			`json:"log_level"`
}

type Server struct {
	ServerIp		string			`json:"server_ip"`
	Ports			[]int			`json:"ports"`
	Services		[]string		`json:"services"`
}

type Monitoring struct {
	Servers			[]Server		`json:"servers"`
}

type SelfMonitoring struct {
	Enabled			bool			`json:"enabled"`
	Ports			[]int			`json:"ports"`
	Services		[]string		`json:"services"`
}

type Database struct {
	Host			string			`json:"host"`
	Port			int				`json:"port"`
	User			string			`json:"user"`
	Password		string			`json:"password"`
	Database		string			`json:"database"`
}

type Backup struct {
	SourceDB		Database		`json:"source_database"`
	DestinationDBs	[]Database		`json:"destination_databases"`
} 

type Email struct {
	Recipients		[]string		`json:"recipients"`
}

type Telegram struct {
	ChatID			string			`json:"chat_id"`
	Token			string			`json:"token"`
} 

type Notifier struct {
	Email			Email			`json:"email"`
	Telegram		Telegram		`json:"telegram"`
}

type Interval struct {
	Availability	int				`json:"availability"`
	Port			int				`json:"port"`
	Service			int				`json:"service"`
}

type Config struct {
	Database		Database 		`json:"database"`
    Monitoring		Monitoring		`json:"monitoring"`	
	SelfMonitoring	SelfMonitoring	`json:"self_monitoring"`
    Backup			Backup			`json:"backup"`
    Logger 			LoggerConfig 	`json:"logger"`
    Notifier		Notifier		`json:"notifier"`
	Interval		Interval		`json:"interval"`
}

func LoadConfig(filename string) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    cfg := Config{}
    err = decoder.Decode(&cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}



type Table struct {
    TableName	string
    OrderBy		string
	Skip		[]string
}

// var DatabaseTableMap = map[string][]TableOrder{
//     "norway": {
//         {TableName: "admins", OrderBy: "org_id", Skip: nil},
//         {TableName: "app_configs", OrderBy: "id", Skip: nil},
//         // Add more tables as needed
//     },
//     "sweden": {
//         {TableName: "admins", OrderBy: "org_id", Skip: nil},
//         {TableName: "app_configs", OrderBy: "id", Skip: nil},
//         // Add more tables as needed
//     },
//     "languages": {
//         {TableName: "backend_strings", OrderBy: "key", Skip: nil},
//         {TableName: "email_strings", OrderBy: "key", Skip: nil},
//         // Add more tables as needed
//     },
//     "safety": {
//         {TableName: "backup_servers", OrderBy: "ip", Skip: nil},
//         {TableName: "blacklisted_ips", OrderBy: "id", Skip: nil},
//         // Add more tables as needed
//     },
//     // Add more databases as needed
// }


var DatabaseTableMap = map[string][]Table{
    "vs_test_database": {
        {TableName: "test", OrderBy: "value"},
    },
}



