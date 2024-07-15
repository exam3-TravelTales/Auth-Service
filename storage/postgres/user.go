package postgres

import (
	pb "auth/genproto/users"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*pb.UserInfo, error) {
	user := pb.UserInfo{Id: id}

	query := `
	SELECT
		username,
		email,
		password,
		full_name,
		bio,
		countries_visited
	FROM
		users
	WHERE
		id = $1 AND deleted_at = 0
	`
	row := r.DB.QueryRowContext(ctx, query, id)

	var bio sql.NullString
	err := row.Scan(&user.Username, &user.Email, &user.Password, &user.FullName, &bio, &user.CountriesVisited)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Assign the bio value based on whether it is valid or not
	if bio.Valid {
		user.Bio = bio.String
	} else {
		user.Bio = ""
	}

	return &user, nil
}

func (r *UserRepo) GetUserProfile(ctx context.Context, id *pb.UserId) (*pb.GetProfileResponse, error) {
	user := pb.GetProfileResponse{Id: id.Id}

	query := `
	SELECT
	username,
	email,
	full_name,
	bio,
	countries_visited,
	created_at,
	updated_at
	from
		users
	where
		deleted_at=0 and id = $1 
	`
	row := r.DB.QueryRowContext(ctx, query, id.Id)
	var bio sql.NullString
	err := row.Scan(&user.Username, &user.Email, &user.FullName, &bio, &user.CountriesVisited, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	if bio.Valid {
		user.Bio = bio.String
	} else {
		user.Bio = ""
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res := pb.RegisterResponse{Username: req.Username, Email: req.Email, FullName: req.FullName}

	query := `
	insert into users (
		username, email, password, full_name
	)
	values (
		$1, $2, $3, $4
	)
	returning id, created_at `
	err := r.DB.QueryRowContext(ctx, query, req.Username, req.Email, req.Password, req.FullName).Scan(&res.Id, &res.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*pb.UserInfo, error) {
	user := pb.UserInfo{Email: email}

	query := `
	SELECT id,
	username,
	password,
	full_name,
	bio,
	countries_visited
	from
		users
	where
		email = $1 and deleted_at=0
	`
	row := r.DB.QueryRowContext(ctx, query, email)
	var bio sql.NullString

	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.FullName, &bio, &user.CountriesVisited)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	if bio.Valid {
		user.Bio = bio.String
	} else {
		user.Bio = ""
	}
	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	query := `update users set `
	n := 1
	var arr []interface{}
	if len(req.Bio) > 0 {
		query += fmt.Sprintf("bio = $%d, ", n)
		arr = append(arr, req.Bio)
		n++
	}
	if len(req.FullName) > 0 {
		query += fmt.Sprintf("full_name = $%d, ", n)
		arr = append(arr, req.FullName)
		n++
	}

	if req.CountriesVisited > 0 {
		query += fmt.Sprintf("countries_visited = $%d, ", n)
		arr = append(arr, req.CountriesVisited)
		n++
	}
	query += fmt.Sprintf("updated_at=current_timestamp where id=$%d and deleted_at=0 ", n)
	arr = append(arr, req.Id)

	_, err := r.DB.Exec(query, arr...)
	if err != nil {
		return nil, err
	}

	var res pb.UpdateProfileResponse
	err = r.DB.QueryRow(`SELECT id,
	username,
	email,
	full_name,
	bio,
	countries_visited,
	updated_at
	from users where id = $1 and deleted_at=0`, req.Id).Scan(&res.Id, &res.Username, &res.Email, &res.FullName, &res.Bio, &res.CountriesVisited, &res.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &res, nil
}

func (r *UserRepo) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	query := "select id,username,full_name,countries_visited from users where deleted_at=0 "
	if req.Limit > 0 {
		query += fmt.Sprintf("LIMIT %d ", req.Limit)
	}
	if req.Offset > 0 {
		query += fmt.Sprintf("OFFSET %d ", req.Offset)
	}
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var n int64
	var res pb.GetUsersResponse
	for rows.Next() {
		no := pb.Users{}
		err = rows.Scan(&no.Id, &no.Username, &no.FullName, &no.CountriesVisited)
		n++
		if err != nil {
			return nil, err
		}
		res.Users = append(res.Users, &no)
	}
	res.Offset = req.Offset
	res.Limit = req.Limit
	res.Total = n
	return &res, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	query := `
	UPDATE
		users 
	SET
		deleted_at=date_part('epoch', current_timestamp)::INT
	WHERE
		deleted_at=0 and id=$1`

	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("failed to delete user")
		return err
	}

	rowAff, err := res.RowsAffected()
	if err != nil {
		log.Println("failed to get rows affected")
		return err
	}

	if rowAff < 1 {
		log.Println("user already deleted or not found")
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, req *pb.EmailRecoveryRequest) error {
	query := `update users set password=$1 where id=$2`
	_, err := r.DB.ExecContext(ctx, query, req.NewPassword, req.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetUserActivity(ctx context.Context, userID string) (*pb.ActivityResponse, error) {
	var activityResponse pb.ActivityResponse

	query := `
        SELECT
            u.id AS user_id,
            COUNT(DISTINCT s.id) AS stories_count,
            COALESCE(SUM(s.comments_count), 0) AS comments_count,
            COALESCE(SUM(s.likes_count), 0) AS likes_received,
            u.countries_visited,
            MAX(GREATEST(s.created_at, s.updated_at)) AS last_activity
        FROM
            users u
        LEFT JOIN
            stories s ON u.id = s.author_id
        WHERE
            u.id = $1
            AND u.deleted_at = 0
        GROUP BY
            u.id, u.countries_visited
    `

	var activity sql.NullString
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&activityResponse.UserId,
		&activityResponse.StoriesCount,
		&activityResponse.CommentsCount,
		&activityResponse.LikesReceived,
		&activityResponse.CountriesVisited,
		&activity,
	)

	if err != nil {
		return nil, err
	}
	if activity.Valid {
		activityResponse.LastActivity = activity.String
	} else {
		activityResponse.LastActivity = ""
	}
	return &activityResponse, nil
}

func (r *UserRepo) Follow(ctx context.Context, followerID string, followingID string) (*pb.FollowResponse, error) {
	// Prepare a variable to hold the followed_at timestamp
	var followedAt string

	// Execute the SQL query to insert into followers table and retrieve followed_at
	query := `
		INSERT INTO followers (follower_id, following_id)
		VALUES ($1, $2)
		RETURNING followed_at
	`
	err := r.DB.QueryRowContext(ctx, query, followerID, followingID).Scan(&followedAt)
	if err != nil {
		return nil, err
	}

	// Construct and return the response
	return &pb.FollowResponse{
		FollowerId:  followerID,
		FollowingId: followingID,
		FollowedAt:  followedAt,
	}, nil
}

func (r *UserRepo) GetFollowers(ctx context.Context, followerID string, limit, offset int64) (*pb.FollowersResponse, error) {
	// Prepare an empty response
	response := &pb.FollowersResponse{
		Followers: make([]*pb.Followers, 0),
		Total:     0,
		Offset:    offset,
		Limit:     limit,
	}

	// Base SQL query
	query := `
        SELECT
            u.id,
            u.username,
            u.full_name
        FROM
            users u
        JOIN
            followers f ON u.id = f.following_id
        WHERE
            f.follower_id = $1
        ORDER BY
            u.username
    `

	// Add LIMIT and OFFSET if provided
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	// Execute the SQL query
	rows, err := r.DB.QueryContext(ctx, query, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the response
	for rows.Next() {
		follower := &pb.Followers{}
		err := rows.Scan(&follower.Id, &follower.Username, &follower.FullName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		response.Followers = append(response.Followers, follower)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Query for total count of followers
	countQuery := `
        SELECT COUNT(*)
        FROM users u
        JOIN followers f ON u.id = f.following_id
        WHERE f.follower_id = $1
    `
	var total int64
	err = r.DB.QueryRowContext(ctx, countQuery, followerID).Scan(&total)
	if err != nil {
		return nil, err
	}
	response.Total = total

	return response, nil
}
