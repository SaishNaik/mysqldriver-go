package mysqldriver

import (
	"strconv"
	"testing"
	"time"

	"github.com/pubnative/mysqlproto-go"
	"github.com/stretchr/testify/assert"
)

func TestQueryError(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Query("SELECT * FROM unknown_table")
		assert.NotNil(t, err)
		assert.True(t, conn.valid)
		pkt, ok := err.(mysqlproto.ERRPacket)
		assert.True(t, ok)
		assert.True(t, ok)
		assert.Equal(t, pkt.Header, mysqlproto.ERR_PACKET)
		assert.Equal(t, pkt.ErrorCode, mysqlproto.ER_NO_SUCH_TABLE)
		assert.Equal(t, pkt.SQLStateMarker, "#")
		assert.Equal(t, pkt.SQLState, "42S02")
		assert.Equal(t, pkt.ErrorMessage, "Table 'test.unknown_table' doesn't exist")
	})
}

func TestQuerySelectValues(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`
			INSERT INTO people(firstname,lastname,cars,houses,cats,dogs,age,married,grade,score)
			VALUES("bob","ben",2,8,16,32,64,1,4.5,3.7)
		`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)

		rows, err := conn.Query("SELECT * FROM people")
		assert.Nil(t, err)
		assert.True(t, rows.Next())
		assert.Equal(t, rows.Int(), 1)
		assert.Equal(t, rows.String(), "bob")
		assert.Equal(t, rows.Bytes(), []byte("ben"))
		assert.Equal(t, rows.Int8(), int8(2))
		assert.Equal(t, rows.Int16(), int16(8))
		assert.Equal(t, rows.Int32(), int32(16))
		assert.Equal(t, rows.Int64(), int64(32))
		assert.Equal(t, rows.Int(), 64)
		assert.Equal(t, rows.Bool(), true)
		assert.Equal(t, rows.Float32(), float32(4.5))
		assert.Equal(t, rows.Float64(), float64(3.7))
		assert.NoError(t, rows.LastError())

		// read non-exist columns
		assert.Equal(t, rows.Int(), 0)
		assert.Equal(t, rows.Int8(), int8(0))
		assert.Equal(t, rows.Int16(), int16(0))
		assert.Equal(t, rows.Int32(), int32(0))
		assert.Equal(t, rows.Int64(), int64(0))
		assert.Equal(t, rows.String(), "")
		assert.Equal(t, rows.Bool(), false)
		assert.Equal(t, rows.Float32(), float32(0.0))
		assert.Equal(t, rows.Float64(), float64(0.0))
		assert.NoError(t, rows.LastError())

		assert.False(t, rows.Next())
	})
}

func TestQuerySelectValuesWithNULL(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`
			INSERT INTO people(firstname,lastname,cars,houses,cats,dogs,age,married,grade,score)
			VALUES("bob","ben",2,8,16,32,64,1,4.5,3.7)
		`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)

		rows, err := conn.Query("SELECT * FROM people")
		assert.Nil(t, err)
		assert.True(t, rows.Next())
		assert.Equal(t, rows.Int(), 1)
		firstname, null := rows.NullString()
		assert.Equal(t, firstname, "bob")
		assert.False(t, null)
		lastname, null := rows.NullBytes()
		assert.Equal(t, lastname, []byte("ben"))
		assert.False(t, null)
		cars, null := rows.NullInt8()
		assert.Equal(t, cars, int8(2))
		assert.False(t, null)
		houses, null := rows.NullInt16()
		assert.Equal(t, houses, int16(8))
		assert.False(t, null)
		cats, null := rows.NullInt32()
		assert.Equal(t, cats, int32(16))
		assert.False(t, null)
		dogs, null := rows.NullInt64()
		assert.Equal(t, dogs, int64(32))
		assert.False(t, null)
		age, null := rows.NullInt()
		assert.Equal(t, age, 64)
		assert.False(t, null)
		married, null := rows.NullBool()
		assert.Equal(t, married, true)
		assert.False(t, null)
		grade, null := rows.NullFloat32()
		assert.Equal(t, grade, float32(4.5))
		assert.False(t, null)
		score, null := rows.NullFloat64()
		assert.Equal(t, score, float64(3.7))
		assert.False(t, null)
		assert.NoError(t, rows.LastError())

		// read non-exist columns
		assert.Equal(t, rows.Int(), 0)
		assert.Equal(t, rows.Int8(), int8(0))
		assert.Equal(t, rows.Int16(), int16(0))
		assert.Equal(t, rows.Int32(), int32(0))
		assert.Equal(t, rows.Int64(), int64(0))
		assert.Equal(t, rows.String(), "")
		assert.Equal(t, rows.Bool(), false)
		assert.Equal(t, rows.Float32(), float32(0.0))
		assert.Equal(t, rows.Float64(), float64(0.0))
		assert.NoError(t, rows.LastError())

		assert.False(t, rows.Next())
	})
}

func TestQuerySelectNULLValues(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`
			INSERT INTO people(firstname,lastname,cars,houses,cats,dogs,age,married,grade,score)
			VALUES(NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL)
		`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)

		rows, err := conn.Query("SELECT * FROM people")
		assert.Nil(t, err)
		assert.True(t, rows.Next())
		assert.Equal(t, rows.Int(), 1)
		firstname, null := rows.NullString()
		assert.Equal(t, firstname, "")
		assert.True(t, null)
		lastname, null := rows.NullBytes()
		assert.Equal(t, lastname, []byte{})
		assert.True(t, null)
		cars, null := rows.NullInt8()
		assert.Equal(t, cars, int8(0))
		assert.True(t, null)
		houses, null := rows.NullInt16()
		assert.Equal(t, houses, int16(0))
		assert.True(t, null)
		cats, null := rows.NullInt32()
		assert.Equal(t, cats, int32(0))
		assert.True(t, null)
		dogs, null := rows.NullInt64()
		assert.Equal(t, dogs, int64(0))
		assert.True(t, null)
		age, null := rows.NullInt()
		assert.Equal(t, age, 0)
		assert.True(t, null)
		married, null := rows.NullBool()
		assert.Equal(t, married, false)
		assert.True(t, null)
		grade, null := rows.NullFloat32()
		assert.Equal(t, grade, float32(0))
		assert.True(t, null)
		score, null := rows.NullFloat64()
		assert.Equal(t, score, float64(0))
		assert.True(t, null)
		assert.False(t, rows.Next())
	})
}

func TestQueryRowReader(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`
			INSERT INTO people(firstname,lastname,cars,houses,cats,dogs,age,married,grade,score,note)
			VALUES
			("bob","ben",2,8,16,32,64,1,4.5,3.7,"good"),
			("one","two",1,2,33,44,55,0,7.7,8.8,"best")
		`)
		assert.NoError(t, err)
		assert.True(t, conn.valid)

		_, err = conn.Exec(`
			INSERT INTO categories(name)
			VALUES ("books"),(NULL),("cars")
		`)
		assert.NoError(t, err)
		assert.True(t, conn.valid)

		rows, err := conn.Query(`
			SELECT id, firstname as name, lastname as lastName,
			p.cars, p.houses as houses, cats, dogs,
			age, married, grade, score, note
			FROM people as p
		`)
		assert.NoError(t, err)

		// switch cursor to first row
		assert.True(t, rows.Next())

		// check sequential calls
		row := rows.Row()
		for i := 0; i < 2; i++ {
			assert.Equal(t, row.Int("id"), 1)
			assert.Equal(t, row.String("name"), "bob")
			assert.Equal(t, row.String("lastName"), "ben")
			assert.Equal(t, row.Int("cars"), 2)
			assert.Equal(t, row.Int8("houses"), int8(8))
			assert.Equal(t, row.Int16("cats"), int16(16))
			assert.Equal(t, row.Int32("dogs"), int32(32))
			assert.Equal(t, row.Int64("age"), int64(64))
			assert.Equal(t, row.Bool("married"), true)
			assert.Equal(t, row.Float32("grade"), float32(4.5))
			assert.Equal(t, row.Float64("score"), float64(3.7))
			assert.Equal(t, row.String("note"), "good")

			row = rows.Row()
		}
		assert.Equal(t, row.String("id"), "1")
		assert.Equal(t, row.String("cars"), "2")
		assert.NoError(t, rows.LastError())

		// switch cursor to the second row
		assert.True(t, rows.Next())

		row = rows.Row()
		assert.Equal(t, row.Int("id"), 2)
		assert.Equal(t, row.String("name"), "one")
		assert.Equal(t, row.String("lastName"), "two")
		assert.Equal(t, row.Int("cars"), 1)
		assert.Equal(t, row.Int8("houses"), int8(2))
		assert.Equal(t, row.Int16("cats"), int16(33))
		assert.Equal(t, row.Int32("dogs"), int32(44))
		assert.Equal(t, row.Int64("age"), int64(55))
		assert.Equal(t, row.Bool("married"), false)
		assert.Equal(t, row.Float32("grade"), float32(7.7))
		assert.Equal(t, row.Float64("score"), float64(8.8))
		assert.Equal(t, row.String("note"), "best")
		assert.NoError(t, rows.LastError())

		// close reading
		assert.False(t, rows.Next())

		rows, err = conn.Query(`SELECT id, cat.name FROM categories AS cat`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		row = rows.Row()
		id, null := row.NullInt("id")
		assert.False(t, null)
		assert.Equal(t, id, 1)
		name, null := row.NullString("name")
		assert.False(t, null)
		assert.Equal(t, name, "books")
		assert.NoError(t, rows.LastError())

		assert.True(t, rows.Next())
		row = rows.Row()
		id2, null := row.NullInt8("id")
		assert.False(t, null)
		assert.Equal(t, id2, int8(2))
		name, null = row.NullString("name")
		assert.True(t, null)
		assert.Equal(t, name, "")
		assert.NoError(t, rows.LastError())

		assert.True(t, rows.Next())
		row = rows.Row()
		id3, null := row.NullInt16("id")
		assert.False(t, null)
		assert.Equal(t, id3, int16(3))
		name, null = row.NullString("name")
		assert.False(t, null)
		assert.Equal(t, name, "cars")
		func() {
			defer func() {
				err := recover()
				assert.Equal(t, err, `mysqldriver: column "id2" doesn't exist. Available columns are: "id", "name"`)
			}()
			row.Int("id2")
		}()
		assert.NoError(t, rows.LastError())

		assert.Equal(t, row.Int("name"), 0)
		assert.EqualError(t, rows.LastError(), `strconv.Atoi: parsing "cars": invalid syntax`)

		assert.False(t, rows.Next())

		rows, err = conn.Query(`SELECT MAX(id), min(id) FROM categories AS cat`)
		assert.NoError(t, err)
		assert.True(t, rows.Next())
		row = rows.Row()
		assert.Equal(t, row.Int("MAX(id)"), 3)
		assert.Equal(t, row.Int("min(id)"), 1)
		assert.False(t, rows.Next())
	})
}

func TestQueryMarkConnInvalidWhenStreamIsBroken(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 10, time.Duration(0))
	conn, err := db.GetConn()
	assert.Nil(t, err)

	assert.Nil(t, conn.Close())
	_, err = conn.Query(`SELECT * FROM people`)
	assert.NotNil(t, err)
	assert.False(t, conn.valid)
}

func TestExecInsertSuccess(t *testing.T) {
	setup(t, func(conn *Conn) {
		pkt, err := conn.Exec(`INSERT INTO people(firstname) VALUES("bob")`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)
		assert.Equal(t, pkt.Header, mysqlproto.OK_PACKET)
		assert.Equal(t, pkt.AffectedRows, uint64(1))
		assert.Equal(t, pkt.LastInsertID, uint64(1))
		assert.Equal(t, pkt.Warnings, uint16(0))
		assert.Equal(t, pkt.Info, "")

		rows, err := conn.Query("SELECT firstname FROM people WHERE id = " + strconv.Itoa(int(pkt.LastInsertID)))
		assert.Nil(t, err)
		assert.True(t, conn.valid)
		assert.True(t, rows.Next())
		assert.Equal(t, rows.String(), "bob")
		assert.False(t, rows.Next())
	})
}

func TestExecInsertError(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`INSERT INTO people(firstname)`)
		assert.NotNil(t, err)
		assert.True(t, conn.valid)
		pkt, ok := err.(mysqlproto.ERRPacket)
		assert.True(t, ok)
		assert.Equal(t, pkt.Header, mysqlproto.ERR_PACKET)
		assert.Equal(t, pkt.ErrorCode, mysqlproto.ER_PARSE_ERROR)
		assert.Equal(t, pkt.SQLStateMarker, "#")
		assert.Equal(t, pkt.SQLState, "42000")
		assert.Equal(t, pkt.ErrorMessage, "You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near '' at line 1")
	})
}

func TestExecDeleteSuccess(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`INSERT INTO people(firstname) VALUES("bob")`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)

		pkt, err := conn.Exec(`DELETE FROM people WHERE firstname = "bob"`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)
		assert.Equal(t, pkt.Header, mysqlproto.OK_PACKET)
		assert.Equal(t, pkt.AffectedRows, uint64(1))
		assert.Equal(t, pkt.LastInsertID, uint64(0))
		assert.Equal(t, pkt.Warnings, uint16(0))
		assert.Equal(t, pkt.Info, "")
	})
}

func TestExecDeleteNotFound(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`INSERT INTO people(firstname) VALUES("bob")`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)

		pkt, err := conn.Exec(`DELETE FROM people WHERE firstname = "ben"`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)
		assert.Equal(t, pkt.Header, mysqlproto.OK_PACKET)
		assert.Equal(t, pkt.AffectedRows, uint64(0))
		assert.Equal(t, pkt.LastInsertID, uint64(0))
		assert.Equal(t, pkt.Warnings, uint16(0))
		assert.Equal(t, pkt.Info, "")
	})
}

func TestExecUpdateSuccess(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`INSERT INTO people(firstname) VALUES("bob")`)
		assert.Nil(t, err)
		_, err = conn.Exec(`INSERT INTO people(firstname) VALUES("bin")`)
		assert.Nil(t, err)

		pkt, err := conn.Exec(`UPDATE people SET firstname = "ben" WHERE firstname = "bob"`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)
		assert.Equal(t, pkt.Header, mysqlproto.OK_PACKET)
		assert.Equal(t, pkt.AffectedRows, uint64(1))
		assert.Equal(t, pkt.LastInsertID, uint64(0))
		assert.Equal(t, pkt.Warnings, uint16(0))
		assert.Equal(t, pkt.Info, "Rows matched: 1  Changed: 1  Warnings: 0")

		rows, err := conn.Query("SELECT firstname FROM people ORDER BY id")
		assert.Nil(t, err)
		assert.True(t, conn.valid)
		assert.True(t, rows.Next())
		assert.Equal(t, rows.String(), "ben")
		assert.True(t, rows.Next())
		assert.Equal(t, rows.String(), "bin")
		assert.False(t, rows.Next())
	})
}

func TestExecUpdateNotFound(t *testing.T) {
	setup(t, func(conn *Conn) {
		_, err := conn.Exec(`INSERT INTO people(firstname) VALUES("bob")`)
		assert.Nil(t, err)

		pkt, err := conn.Exec(`UPDATE people SET firstname = "ben" WHERE firstname = "bin"`)
		assert.Nil(t, err)
		assert.True(t, conn.valid)
		assert.Equal(t, pkt.Header, mysqlproto.OK_PACKET)
		assert.Equal(t, pkt.AffectedRows, uint64(0))
		assert.Equal(t, pkt.LastInsertID, uint64(0))
		assert.Equal(t, pkt.Warnings, uint16(0))
		assert.Equal(t, pkt.Info, "Rows matched: 0  Changed: 0  Warnings: 0")
	})
}

func TestExecMarkConnInvalidWhenStreamIsBroken(t *testing.T) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 10, time.Duration(0))
	conn, err := db.GetConn()
	assert.Nil(t, err)

	assert.Nil(t, conn.Close())
	_, err = conn.Exec(`INSERT INTO people(firstname) VALUES("bob")`)
	assert.NotNil(t, err)
	assert.False(t, conn.valid)
}

func setup(t *testing.T, fn func(conn *Conn)) {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 10, time.Duration(0))
	conn, err := db.GetConn()
	assert.Nil(t, err)

	_, err = conn.Exec(`CREATE TABLE people (
		id int NOT NULL AUTO_INCREMENT,
		firstname varchar(255),
		lastname varchar(255),
		cars tinyint,
		houses tinyint,
		cats int,
		dogs int,
		age int,
		married tinyint,
		grade decimal(6,2),
		score decimal(6,2),
		note text,
		PRIMARY KEY (id)
	)`)
	assert.Nil(t, err)

	_, err = conn.Exec(`CREATE TABLE categories (
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255),
		PRIMARY KEY (id)
	)`)
	assert.Nil(t, err)

	fn(conn)

	defer func() {
		assert.Nil(t, db.PutConn(conn))
		_, err = conn.Exec(`DROP TABLE people, categories`)
		assert.Nil(t, err)
	}()
}

func ExampleConn_Query_default() {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 10, time.Duration(0))
	conn, err := db.GetConn()
	if err != nil {
		// handle error
	}
	rows, err := conn.Query("SELECT id,badge,age,honors,length,weight,height,male,name,info FROM dogs")
	if err != nil {
		// handle error
	}
	for rows.Next() {
		_ = rows.Int()     // id
		_ = rows.Int8()    // badge
		_ = rows.Int16()   // age
		_ = rows.Int32()   // honors
		_ = rows.Int64()   // length
		_ = rows.Float32() // weight
		_ = rows.Float64() // height
		_ = rows.Bool()    // male
		_ = rows.String()  // name
		_ = rows.Bytes()   // info
	}
	if err = rows.LastError(); err != nil {
		// handle error

		// when error occurred during reading from the stream
		// connection must be manually closed to prevent further reuse
		conn.Close()
	}
}

func ExampleConn_Query_null() {
	db := NewDB("root@tcp(127.0.0.1:3306)/test", 10, time.Duration(0))
	conn, err := db.GetConn()
	if err != nil {
		// handle error
	}
	rows, err := conn.Query("SELECT id,badge,age,honors,length,weight,height,male,name,info FROM dogs")
	if err != nil {
		// handle error
	}
	for rows.Next() {
		_, _ = rows.NullInt()     // id
		_, _ = rows.NullInt8()    // badge
		_, _ = rows.NullInt16()   // age
		_, _ = rows.NullInt32()   // honors
		_, _ = rows.NullInt64()   // length
		_, _ = rows.NullFloat32() // weight
		_, _ = rows.NullFloat64() // height
		_, _ = rows.NullBool()    // male
		_, _ = rows.NullString()  // name
		_, _ = rows.NullBytes()   // info
	}
	if err = rows.LastError(); err != nil {
		// handle error

		// when error occurred during reading from the stream
		// connection must be manually closed to prevent further reuse
		conn.Close()
	}
}
