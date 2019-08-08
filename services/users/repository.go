package users

import "github.com/jinzhu/gorm"
import "fmt"

type Repository struct {
	db *gorm.DB
}

// InjectDep is a function for inject db to Repository object
func ProvideRepo(db *gorm.DB) *Repository {
	return &Repository{db}
}

// Find is a function to find list of object with parameter, offset and limit
// Using native query for get dynamic parameter
// visit http://gorm.io/docs/sql_builder.html for reference
func (repo *Repository) Find(param User, offset string, limit string) ([]User, bool) {
	var users []User

	var sql = "SELECT * FROM users WHERE 1=1"

	if param.Username != "" {
		sql += " AND username = '" + param.Username + "'"
	}

	if param.Name != "" {
		sql += " AND name LIKE '%" + param.Name + "%'"
	}

	if param.RoleId != 0 {
		sql += " AND role_id = " + fmt.Sprint(param.RoleId)
	}

	sql += " LIMIT " + offset + "," + limit

	result := repo.db.Raw(sql).Scan(&users).RecordNotFound()

	return users, result
}

// Count is a function to count length list of object with parameter
// Using native query for get dynamic parameter
// visit http://gorm.io/docs/sql_builder.html for reference
func (repo *Repository) Count(param User) uint {
	var result uint

	var sql = "SELECT * FROM user WHERE 1=1"

	if param.Name != "" {
		sql += " AND name LIKE '%" + param.Name + "%'"
	}

	if param.RoleId != 0 {
		sql += " AND role_id = " + fmt.Sprint(param.RoleId)
	}

	repo.db.Raw(sql).Scan(&result)

	return result
}

// FindByID is function to find specific object by id as a param
// visit http://gorm.io/docs/query.html for reference
func (repo *Repository) FindByID(id uint64) (User, bool) {
	var user User

	result := repo.db.Where("id = ?", id).First(&user).RecordNotFound()

	return user, result
}

// FindByUsername is function to find specific object by username as a param
// visit http://gorm.io/docs/query.html for reference
func (repo *Repository) FindByUsername(username string) (User, bool) {
	var user User

	result := repo.db.Where("username = ?", username).First(&user).RecordNotFound()

	return user, result
}

// Save is function to save data to table
// visit http://gorm.io/docs/create.html
func (repo *Repository) Save(entity User) (User, error) {
	err := repo.db.Create(&entity)

	return entity, err.Error
}

// Update is function to update data those changed fields
// visit http://gorm.io/docs/update.html
func (repo *Repository) Update(entity User) (User, int64) {
	result := repo.db.Model(&entity).Updates(User{Name: entity.Name, Password: entity.Password, RoleId: entity.RoleId})

	return entity, result.RowsAffected
}

// Delete is function to delete data (flagged)
// visit http://gorm.io/docs/delete.html
// using approach to not permanently delete data, just update on deleteAt column
// to delete permanently use db.Unscoped().Delete(&entity)
func (repo *Repository) Delete(entity User) (User, int64) {
	result := repo.db.Delete(&entity)

	return entity, result.RowsAffected
}
