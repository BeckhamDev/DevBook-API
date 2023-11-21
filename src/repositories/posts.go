package repositories

import (
	"api/src/models"
	"database/sql"
)

type Posts struct {
	db *sql.DB
}

func NewPostRep(db *sql.DB) *Posts {
	return &Posts{db}
}

func (p Posts) CreatePost(post models.Post) (uint64, error){
	sql, err := p.db.Prepare("insert into posts (title, content, author_id) values(?,?,?)")
	if err != nil {
		return 0, err
	}
	defer sql.Close()

	result, err := sql.Exec(post.Title, post.Content, post.AuthorID)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}

func (p Posts) GetOnePost(postID uint64) (models.Post, error){
	sql, err := p.db.Query("select p.*, u.nick from posts p inner join users u on u.id = p.author_id where p.id = ?", postID)
	if err != nil {
		return models.Post{}, err
	}
	defer sql.Close()

	var post models.Post

	if sql.Next() {
		if err = sql.Scan(
			&post.ID,
			&post.Title, 
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return models.Post{}, err
		}
	}

	return post, nil
}

func (p Posts) SearchPosts(userID uint64) ([]models.Post, error){
	sql, err := p.db.Query("select distinct p.*, u.nick from posts p inner join users u on u.id = p.author_id inner join followers f on p.author_id = f.user_id where u.id = ? or f.follower_id = ? order by 1 desc", userID, userID)
	if err != nil {
		return nil , err
	}
	defer sql.Close()

	var posts []models.Post

	for sql.Next(){
		var post models.Post
		if err = sql.Scan(
			&post.ID,
			&post.Title, 
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (p Posts) Update(postID uint64, post models.Post) error{
	sql, err := p.db.Prepare("update posts set title=?, content=? where id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _, err = sql.Exec(post.Title, post.Content, postID); err != nil {
		return err
	}

	return nil
}

func (p Posts) DeletePost(postID uint64) error{
	sql, err := p.db.Prepare("delete from posts where id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _, err = sql.Exec(postID); err != nil {
		return err
	}

	return nil
}

func (p Posts) GetUserPosts(userID uint64) ([]models.Post, error){
	sql, err := p.db.Query("select p.*, u.nick from posts p join users u on u.id = p.author_id where p.author_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer sql.Close()

	var posts []models.Post

	for sql.Next(){
		var post models.Post

		if err = sql.Scan(
			&post.ID,
			&post.Title, 
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (p Posts) Like(postID uint64) error{
	sql, err := p.db.Prepare("update posts p set likes = likes + 1 where id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _, err := sql.Exec(postID); err != nil {
		return err
	} 

	return nil
}

func (p Posts) Unlike(postID uint64) error{
	sql, err := p.db.Prepare("update posts p set likes = CASE WHEN likes > 0 THEN likes - 1 ELSE 0 END where id = ?")
	if err != nil {
		return err
	}
	defer sql.Close()

	if _, err := sql.Exec(postID); err != nil {
		return err
	}

	return nil
}