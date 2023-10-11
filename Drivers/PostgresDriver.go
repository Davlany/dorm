package Drivers

import (
	"database/sql"
	"dorm/pkg"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
)

type PostgresDriver struct {
	conn *sqlx.DB
}

type PgTable struct {
	name string
	pd   *PostgresDriver
}

type MyTable struct {
	TableName string `db:"table_name"`
}

func (pt PgTable) InsertOne(entity interface{}) (int, error) {

	queryParam := pkg.ScanTagsFromKeyInStruct(entity, "db")
	query := fmt.Sprintf("INSERT INTO %s(", pt.name)
	k := 1
	var keys []string
	for key := range queryParam {
		if k == len(queryParam) {
			keys = append(keys, key)
			query += fmt.Sprintf("%s", key)
		} else {
			keys = append(keys, key)
			query += fmt.Sprintf("%s,", key)
		}
		k++
	}
	k = 1
	query += ") values("
	for _, value := range keys {
		if k == len(queryParam) {
			if reflect.TypeOf(queryParam[value]).Kind() == reflect.String {
				query += fmt.Sprintf("'%s'", queryParam[value])
			} else {
				query += fmt.Sprintf("%d", queryParam[value])
			}
		} else {
			if reflect.TypeOf(queryParam[value]).Kind() == reflect.String {
				query += fmt.Sprintf("'%s',", queryParam[value])
			} else {
				query += fmt.Sprintf("%d,", queryParam[value])
			}
		}
		k++
	}
	query += ") returning id;"
	var res int
	err := pt.pd.conn.QueryRowx(query).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil

}

func (pt PgTable) InsertMany(entities interface{}) error {
	num := reflect.ValueOf(entities).Len()
	query := fmt.Sprintf("INSERT INTO %s(", pt.name)
	var entitiesValue []map[string]interface{}
	for i := 0; i < num; i++ {
		val := pkg.ScanTagsFromKeyInStruct(reflect.ValueOf(entities).Index(i).Interface(), "db")
		entitiesValue = append(entitiesValue, val)
	}
	i := 1
	var keys []string
	for key, _ := range entitiesValue[0] {
		if i == len(entitiesValue[0]) {
			query += fmt.Sprintf("%s", key)
			keys = append(keys, key)
		} else {
			query += fmt.Sprintf("%s,", key)
			keys = append(keys, key)
		}
		i++
	}
	query += ") VALUES "
	n := 1
	for _, ent := range entitiesValue {
		query += "("
		for i := 0; i < len(keys); i++ {
			if reflect.TypeOf(ent[keys[i]]).Kind() == reflect.String {
				if i == len(keys)-1 {
					query += fmt.Sprintf("'%s'", ent[keys[i]])
				} else {
					query += fmt.Sprintf("'%s',", ent[keys[i]])
				}
			} else {
				if i == len(keys)-1 {
					query += fmt.Sprintf("%d", ent[keys[i]])
				} else {
					query += fmt.Sprintf("%d,", ent[keys[i]])
				}
			}

		}
		if n == len(entitiesValue) {
			query += ");"
		} else {
			query += "),"
		}
		n++
	}
	_, err := pt.pd.conn.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (pt PgTable) FindOne(id interface{}, dest interface{}) error {
	i := reflect.TypeOf(id)
	if i.Kind() == reflect.String {
		query := fmt.Sprintf("SELECT * from %s WHERE id = '%s'", pt.name, id)
		err := pt.pd.conn.Get(dest, query)
		return err
	} else {
		query := fmt.Sprintf("SELECT * from %s WHERE id = %d", pt.name, id)
		err := pt.pd.conn.Get(dest, query)
		return err
	}
}

func (pt PgTable) FindAll(dest interface{}) error {
	query := fmt.Sprintf("SELECT * from %s", pt.name)
	err := pt.pd.conn.Select(dest, query)
	if err != nil {
		return err
	}
	return nil
}

func (pt PgTable) UpdateOne(entity interface{}) error {
	tagsValue := pkg.ScanTagsFromKeyInStruct(entity, "db")
	id := tagsValue["id"]
	delete(tagsValue, "id")
	query := fmt.Sprintf("UPDATE %s SET ", pt.name)
	i := 1
	for key, value := range tagsValue {
		if i == len(tagsValue) {
			if reflect.TypeOf(value).Kind() == reflect.String {
				query += fmt.Sprintf("%s = '%s'", key, value)
			} else {
				query += fmt.Sprintf("%s = %d", key, value)
			}
		} else {
			if value == reflect.String {
				query += fmt.Sprintf("%s = '%s',", key, value)
			} else {
				query += fmt.Sprintf("%s = %d,", key, value)
			}
		}
		i++
	}
	if reflect.TypeOf(id).Kind() == reflect.String {
		query += fmt.Sprintf(" WHERE id = '%s'", id)
	} else {
		query += fmt.Sprintf(" WHERE id = %d", id)
	}

	_, err := pt.pd.conn.Query(query)

	return err
}

func (pt PgTable) UpdateMany(entities interface{}) error {
	num := reflect.TypeOf(entities).Len()
	sl := reflect.ValueOf(entities)
	for i := 0; i < num; i++ {
		err := pt.UpdateOne(sl.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (pd PostgresDriver) ConnTable(name string, strct interface{}) pkg.Table {
	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public'", name)
	err := pd.conn.Get(&MyTable{}, query)
	if err == sql.ErrNoRows {
		createTableQuery := fmt.Sprintf("CREATE TABLE %s (", name)
		tagsType := pkg.ScanTypeFromKeyInStruct(strct, "db")

		DbTypes := map[string]string{
			"int":     "INTEGER",
			"string":  "TEXT",
			"float32": "DOUBLE",
			"float64": "DOUBLE",
		}

		for name, dataType := range tagsType {
			createTableQuery += fmt.Sprintf("%s %s,", name, DbTypes[dataType])
		}
		createTableQuery = createTableQuery[:len(createTableQuery)-1] + ")"
		fmt.Println(createTableQuery)
		_, err := pd.conn.Query(createTableQuery)
		if err != nil {
			log.Fatalln(err, "sdsd")
		}
	}
	return PgTable{
		name: name,
		pd:   &pd,
	}
}

func NewPostgresDriver(user, password, dbName, port, sslMode string) (*PostgresDriver, error) {
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("user = %s password = %s dbname = %s sslmode = %s port = %s", user, password, dbName, sslMode, port))
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresDriver{
		conn: conn,
	}, nil
}
