package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	minder_model "github.com/AlvinTendio/minder/minder/model"
)

type minderRepositoryImpl struct {
	DB *sql.DB
}

func NewMinderRepositoryImpl(db *sql.DB) MinderRepository {
	return &minderRepositoryImpl{DB: db}
}

const (
	insertUsers = `INSERT INTO Users (username, email, phone_number, password, full_name, gender, date_of_birth, profile_picture)
					VALUES (?,?,?,?,?,?,?,?)`
	getUsersLoginData = `SELECT user_id, username, email, phone_number, full_name, gender, date_of_birth, profile_picture, is_upgraded FROM Users
					WHERE username = ? AND password = ?`
	upgradeAccount = `UPDATE Users set is_upgraded=true WHERE user_id =?`

	getUserUpgradeStatus = `SELECT is_upgraded FROM Users WHERE user_id=?`

	getUserViewCount = `SELECT COUNT(1) AS total FROM Swipes WHERE user_id=? AND DATE(created_at) = CURDATE()`

	getTargetUser = `SELECT u.user_id, u.username, u.email, u.phone_number, u.full_name, u.gender, u.date_of_birth, u.profile_picture
			FROM Users u
			WHERE u.user_id NOT IN (
				SELECT s.target_user_id
				FROM Swipes s
				WHERE s.user_id = ?
					AND DATE(s.created_at) = CURDATE()
			)
			AND u.gender != (
				SELECT gender
				FROM Users
				WHERE user_id = ?
			)
			LIMIT 1`

	insertSwipeLog = `INSERT INTO Swipes (user_id, target_user_id) 
						VALUES(?,?)`

	updateSwipeLog = `UPDATE Swipes SET swipe_action=? WHERE user_id=? AND target_user_id=?`
)

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Print(err)
	}
}

func (r *minderRepositoryImpl) Register(ctx context.Context, req *minder_model.RegisterReq) (data int64, err error) {
	stmt, err := r.DB.PrepareContext(ctx, insertUsers)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Register err", err)
		return
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, req.Username, req.Email, req.PhoneNumber, req.Password, req.FullName, req.Gender, req.DateOfBirth, req.ProfilePicture)
	if err != nil {
		log.Println(ctx, "[repository:minder] Insert Register err ", err)
		return
	}

	return result.RowsAffected()

}

func (r *minderRepositoryImpl) Login(ctx context.Context, req *minder_model.LoginReq) (data *minder_model.UserData, err error) {
	stmt, err := r.DB.PrepareContext(ctx, getUsersLoginData)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Get User Login Data Count err", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, req.Username, req.Password)
	if err != nil {
		log.Println(ctx, "[repository:minder] User Login Data err", err)
		return
	}

	defer func() {
		closeRows(rows)
		if err := rows.Err(); err != nil {
			log.Println("ERROR: ", err)
		}
	}()

	if !rows.Next() {
		return nil, fmt.Errorf("no user found with the given credentials")
	}

	var userData minder_model.UserData
	err = rows.Scan(
		&userData.UserId,
		&userData.Username,
		&userData.Email,
		&userData.PhoneNumber,
		&userData.FullName,
		&userData.Gender,
		&userData.DateOfBirth,
		&userData.ProfilePicture,
		&userData.IsUpgraded,
	)
	return &userData, err
}
func (r *minderRepositoryImpl) UpgradeAccount(ctx context.Context, id uint64) (data int64, err error) {
	stmt, err := r.DB.PrepareContext(ctx, upgradeAccount)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Upgrade Account err", err)
		return
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Println(ctx, "[repository:minder] Upgrade Account err", err)
		return
	}
	data, err = result.RowsAffected()
	if err != nil {
		log.Println(ctx, "[repository:minder] Get RowsEffected  err", err)
	}
	return
}

func (r *minderRepositoryImpl) GetUserUpgradeStatus(ctx context.Context, id uint64) (data bool, err error) {
	stmt, err := r.DB.PrepareContext(ctx, getUserUpgradeStatus)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Get User Upgrade Status err", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		log.Println(ctx, "[repository:minder] User Upgrade Status err", err)
		return
	}

	defer func() {
		closeRows(rows)
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
	}()

	if !rows.Next() {
		return false, fmt.Errorf("no data found")
	}

	err = rows.Scan(
		&data,
	)

	return
}

func (r *minderRepositoryImpl) GetUserViewCount(ctx context.Context, id uint64) (total *int64, err error) {
	stmt, err := r.DB.PrepareContext(ctx, getUserViewCount)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Get User View Count err", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		log.Println(ctx, "[repository:minder] User View Count err", err)
		return
	}

	defer func() {
		closeRows(rows)
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
	}()
	if !rows.Next() {
		return nil, fmt.Errorf("no data found")
	}
	err = rows.Scan(
		&total,
	)
	return
}
func (r *minderRepositoryImpl) GetTargetUser(ctx context.Context, id uint64) (data *minder_model.TargetUserData, err error) {
	stmt, err := r.DB.PrepareContext(ctx, getTargetUser)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Get Targer User Data err", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id, id)
	if err != nil {
		log.Println(ctx, "[repository:minder] Get Target User Data err", err)
		return
	}

	defer func() {
		closeRows(rows)
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
	}()
	if !rows.Next() {
		return nil, fmt.Errorf("no data found")
	}
	var userData minder_model.TargetUserData
	err = rows.Scan(
		&userData.UserId,
		&userData.Username,
		&userData.Email,
		&userData.PhoneNumber,
		&userData.FullName,
		&userData.Gender,
		&userData.DateOfBirth,
		&userData.ProfilePicture,
	)
	if err != nil {
		log.Println("[repository:minder] Error scanning row:", err)
		return nil, err
	}

	return &userData, err
}

func (r *minderRepositoryImpl) InsertSwipe(ctx context.Context, id uint64, targetId int64) (data int64, err error) {
	stmt, err := r.DB.PrepareContext(ctx, insertSwipeLog)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Swipe err", err)
		return
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, id, targetId)
	if err != nil {
		log.Println(ctx, "[repository:minder] Insert Swipe err ", err)
		return
	}

	return result.RowsAffected()
}

func (r *minderRepositoryImpl) UpdateSwipe(ctx context.Context, req *minder_model.SwipeReq) (data int64, err error) {
	stmt, err := r.DB.PrepareContext(ctx, updateSwipeLog)
	if err != nil {
		log.Println(ctx, "[repository:minder] Preparing Swipe err", err)
		return
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, req.Action, req.Id, req.TargetId)
	if err != nil {
		log.Println(ctx, "[repository:minder] Insert Swipe err ", err)
		return
	}

	return result.RowsAffected()
}
