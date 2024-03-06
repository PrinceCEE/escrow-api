package testutils

// type testServer struct {
// 	Server *httptest.Server
// 	Config *config.Config
// }

// func NewTestServer() *testServer {
// 	c := config.NewConfig()
// 	c.Env.DSN = "postgres://postgres:password@localhost/bupher-test?sslmode=disable"
// 	r := routes.GetRouter(c)

// 	s := httptest.NewServer(r)

// 	if err := createTypes(c.DB); err != nil {
// 		panic(err)
// 	}

// 	return &testServer{s, c}
// }

// func (ts *testServer) DropTables() {
// 	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`

// 	ctx := context.Background()

// 	rows, err := ts.Config.DB.Query(ctx, query)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for rows.Next() {
// 		var tableName string
// 		if err := rows.Scan(&tableName); err != nil {
// 			panic(err)
// 		}

// 		_, err := ts.Config.DB.Exec(ctx, `DROP TABLE %s CASCADE`, tableName)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }
