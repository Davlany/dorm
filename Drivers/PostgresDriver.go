package Drivers

import (
	"database/sql"
	"dorm"
	"dorm/pkg"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PostgresDriver struct {
	conn *sqlx.DB
}

type PgTable struct {
	name string `db:"table_name"`
	pd   *PostgresDriver
}

func (pt PgTable) InsertOne(entity interface{}) (int, error) {
	queryParam := pkg.ScanTagsFromKeyInStruct(entity, "db")
	query := fmt.Sprintf("INSERT INTO %s (", pt.name)
	fmt.Println(pkg.ScanTypeFromKeyInStruct(entity))
	var keys []string
	for name := range queryParam {
		if name == "id" && pkg.ScanTypeFromKeyInStruct(entity)["id"]["dataType"] == "serial" {
			continue
		}
		if name == "" {
			continue
		}
		keys = append(keys, name)
		query += fmt.Sprintf("%s,", name)
	}
	query = query[:len(query)-1] + ") values ("
	for _, key := range keys {
		if key == "id" && pkg.ScanTypeFromKeyInStruct(entity)["id"]["dataType"] == "serial" {
			continue
		}
		if reflect.TypeOf(queryParam[key]).Kind() == reflect.String {
			query += fmt.Sprintf("'%s',", queryParam[key])
		} else {
			query += fmt.Sprintf("%d,", queryParam[key])
		}

	}
	query = query[:len(query)-1] + ") returning id;"
	fmt.Println(query)
	var res int
	err := pt.pd.conn.QueryRowx(query).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil

}

func (pt PgTable) InsertMany(entities interface{}) error {
	num := reflect.ValueOf(entities).Len()
	for i := 0; i < num; i++ {
		entityValue := reflect.ValueOf(entities).Index(i).Interface()
		_, err := pt.InsertOne(entityValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pt PgTable) FindOne(id interface{}, dest interface{}) error {
	i := reflect.TypeOf(id)
	if i.Kind() == reflect.String {
		query := fmt.Sprintf("SELECT * from %s WHERE id = '%s'", pt.name, id)
		err := pt.pd.conn.Get(dest, query)
		fmt.Println(query)
		return err
	} else {
		query := fmt.Sprintf("SELECT * from %s WHERE id = %d", pt.name, id)
		err := pt.pd.conn.Get(dest, query)
		fmt.Println(query)
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

func (pt PgTable) DeleteOne(entity interface{}) error {
	entityId := pkg.ScanTagsFromKeyInStruct(entity, "db")["id"]
	var query string
	if pkg.ScanTypeFromKeyInStruct(entity)["id"]["dataType"] == "string" {
		query = fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", pt.name, entityId)
	} else {
		query = fmt.Sprintf("DELETE FROM %s WHERE id = %d", pt.name, entityId)
	}
	_, err := pt.pd.conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// CREATE TABLE {name}(
// >> ID
// >> Name
// >>
// >>

func (pd PostgresDriver) RegisterSchemas(schemas interface{}) error {
	num := reflect.ValueOf(schemas).Len()
	for i := 0; i < num; i++ {
		fmt.Println(i)
		schemaName := strings.ToLower(reflect.TypeOf(reflect.ValueOf(schemas).Index(i).Interface()).Name() + "s")
		query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public'", schemaName)
		err := pd.conn.Get(&PgTable{}, query)
		if errors.Is(err, sql.ErrNoRows) {
			createTableQuery := fmt.Sprintf("CREATE TABLE %s (", schemaName)
			tagsType := pkg.ScanTypeFromKeyInStruct(reflect.ValueOf(schemas).Index(i).Interface())

			DbTypes := map[string]string{
				"int":     "INTEGER",
				"serial":  "SERIAL PRIMARY KEY",
				"string":  "TEXT",
				"float32": "DOUBLE",
				"float64": "DOUBLE",
			}

			for fieldName, types := range tagsType {
				if fieldName == "rel" {
					continue
				}
				createTableQuery += fieldName + " " + DbTypes[types["dataType"]]
				skipTags := []string{"fk", "pk", "field"}
				for typeName, value := range types {

					if slices.Contains(skipTags, typeName) {
						continue
					}

					if typeName == "rel" && value != "" {
						if types["fk"] == "true" {
							if types["field"] != "" {
								fkQuery := fmt.Sprintf("FOREIGN KEY(%s) REFERENCES %s(%s)", fieldName, types["rel"], types["field"])
								createTableQuery += "," + fkQuery
							} else {
								log.Fatal("Empty field")
							}
						}
					} else {
						if value == "true" {
							createTableQuery += DbTypes[typeName]
						}
					}
				}
				createTableQuery += ","
			}

			createTableQuery = createTableQuery[:len(createTableQuery)-1] + ");"
			fmt.Println(createTableQuery)
			_, err = pd.conn.Query(createTableQuery)
			if err != nil {
				return err
			}
		} else if err != nil {
			log.Println(err)
		} else {
			continue
		}
	}
	return nil
}

func (pd PostgresDriver) ConnTable(name string, strct interface{}) dorm.Table {
	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public'", name)
	err := pd.conn.Get(&PgTable{}, query)
	if errors.Is(err, sql.ErrNoRows) {
		createTableQuery := fmt.Sprintf("CREATE TABLE %s (", name)
		tagsType := pkg.ScanTypeFromKeyInStruct(strct)

		DbTypes := map[string]string{
			"int":     "INTEGER",
			"serial":  "SERIAL PRIMARY KEY",
			"string":  "TEXT",
			"float32": "DOUBLE",
			"float64": "DOUBLE",
		}

		for fieldName, types := range tagsType {
			createTableQuery += fieldName + " " + DbTypes[types["dataType"]]
			skipTags := []string{"fk", "pk", "field"}
			for typeName, value := range types {

				if slices.Contains(skipTags, typeName) {
					continue
				}

				if typeName == "rel" && value != "" {
					if types["fk"] == "true" {
						if types["field"] != "" {
							fkQuery := fmt.Sprintf("FOREIGN KEY(%s) REFERENCES %s(%s)", fieldName, types["rel"], types["field"])
							createTableQuery += "," + fkQuery
						} else {
							log.Fatal("Empty field")
						}
					}
				} else {
					if value == "true" {
						createTableQuery += DbTypes[typeName]
					}
				}
			}
			createTableQuery += ","
		}

		createTableQuery = createTableQuery[:len(createTableQuery)-1] + ");"
		fmt.Println(createTableQuery)
		_, err = pd.conn.Query(createTableQuery)
		if err != nil {
			log.Fatalln("Executing query:", err)
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
