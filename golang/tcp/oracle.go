package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

func main() {
	// var dbUser, passwd string

	// dbUser = "admin"
	// passwd = "Liujd.920928"
	// dataSourceUrl = fmt.Sprintf("oracle://%s:%s@%s/?sid=%s&", dbUser, passwd, "adb.ap-seoul-1.oraclecloud.com:1522", "bpxauahrflyyzbz_ku8dxupc5djaayss_high.adb.oraclecloud.com")

	// urlOptions := map[string]string{
	// 	"TRACE FILE":    "trace.log",
	// 	"SERVER":        "server2, server3",
	// 	"PREFETCH_ROWS": "500",
	// 	//"SSL": "enable",
	// 	//"SSL Verify": "false",
	// }
	walletPath := "/root/codehub/github/redis-operator/tmp/tcp/wallet"
	dbUrl := fmt.Sprintf("oracle://admin:Liujd.920928@adb.ap-seoul-1.oraclecloud.com:1522/bpxauahrflyyzbz_ku8dxupc5djaayss_high.adb.oraclecloud.com?TRACE FILE=trace.log&wallet=", walletPath)
	// dsn := "(description= (retry_count=20)(retry_delay=3)(address=(protocol=tcps)(port=1522)(host=adb.ap-seoul-1.oraclecloud.com))(connect_data=(service_name=bpxauahrflyyzbz_ku8dxupc5djaayss_high.adb.oraclecloud.com))(security=(ssl_server_dn_match=yes)))"
	// dbUrl := ora.BuildJDBC(dbUser, passwd, dsn, urlOptions)
	fmt.Println(dbUrl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	db, err := sql.Open("oracle", dbUrl)
	if err != nil {
		panic(err)
	}
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

}

// (description= (retry_count=20)(retry_delay=3)(address=(protocol=tcps)(port=1522)(host=adb.ap-seoul-1.oraclecloud.com))(connect_data=(service_name=bpxauahrflyyzbz_ku8dxupc5djaayss_high.adb.oraclecloud.com))(security=(ssl_server_dn_match=yes)))
