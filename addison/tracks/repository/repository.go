package repository

import(
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

type Track struct {
	Id      string
	Audio   string
}
type Repository struct {
	DB *sql.DB

}

var repo Repository

func Init() {
	//Database initialisation and check that it has happened without errors
	if db, err := sql.Open("sqlite3", "TracksDB.db"); err == nil {
		repo = Repository{DB: db}
	} else {
		log.Fatal("Database initialisation")
	}
}

func Create() int {
	//creatung table
    const sql = "CREATE TABLE IF NOT EXISTS Tracks(Id TEXT PRIMARY KEY, Audio TEXT)"
    if _, err := repo.DB.Exec(sql); err == nil {
        return 0
    } else {
        return -1
    }
}

func Clear() int {
	//clearing the table
    const sql = "DELETE FROM Tracks"
    if _, err := repo.DB.Exec(sql); err == nil {
        return 0
    } else {
        return -1
    }
}

func AddTrack(t Track) int64{
	//function that allows user to add track into database
	const sql = " INSERT INTO Tracks (Id , Audio) VALUES (? , ?)"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Id, t.Audio) ; err == nil {
			if n, err := res.RowsAffected() ; err == nil {
				return n
			}
		}
	}
	return -1
}

func GetAllTracks()([]Track, int){
	//function that returns all the tracks from the database
	const sql = "SELECT * FROM Tracks"
	if stmt, err := repo.DB.Prepare(sql); err == nil{
		defer stmt.Close()
		rows, err := stmt.Query()
		if err != nil {
			return []Track{}, -1
		}

		tracks := []Track{}
		for rows.Next() {
			track := Track{}
			err := rows.Scan(&track.Id, &track.Audio)
			if err != nil {
				return []Track{}, -1
			}
			tracks = append(tracks, track)
		}
		if err := rows.Err(); err != nil {
			return []Track{}, -1
		}
		return tracks, len(tracks)
		
	}else{
		return []Track{}, -1
	}
}
func GetTrack (id string) (Track, int64){
	//returns single track based on id
	const sql = " SELECT * FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql) ; err == nil {
		defer stmt.Close()
		var c Track
		row := stmt.QueryRow(id)
		if err := row.Scan(&c.Id, &c.Audio); err == nil {
			return c, 1
		} else {
			return Track{} , 0
		}
	}
	return Track{} , -1
}
func DeleteTrack(id string) int64{
	//function that deletes a track from the database based on id
	const sql = "DELETE FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}else{
				return -1
			}
		}else{
			return -1
		}
	}else{
		return -1
	}
}