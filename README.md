# go-qb
Query builder and helpers for github.com/jmoiron/sqlx package.

## Install
````go get github.com/nirandas/go-qb````

````
type User struct {
ID int64 `db:"id"`
Name string `db:"name"`
Age int `db:"age"`
}

func (u *User) PK() *int64{
return &u.ID
}
	func (u *User) TableName() string{
	return "users"
	}
	func (u *User) Fields() []string{
	return []string{"id", "name", "age"}
	}
	
	
	func main(){
	//....
	user:=new(User)
	user.Name="Nirandas"
	err:=qb.InsertRow(db,user)
	fmt.Println(user.ID)//auto generated id field from database.
	//update user
	user.Name="Nirandas T"
	err=qb.UpdateRow(db,user)
	//delete user
	err=qb.DeleteRow(db,user)
	//...
err=qb.FindRowByPK(db,&user,1)
//fetch all users
users:=[]*User{}
q:=qb.NewQB()
err=qb.ListRows(db,&users,q)
//fetch all users with condition
users:=[]*User{}
q:=qb.NewQB()
q.Where("age >= ? and age < ?", 18, 60)
err=qb.ListRows(db,&users,q)
	}
````

@@ DBTable interface

A type implementing the DBTable interface can be easily processed using InsertRow, UpdateRow, DeleteRow and ListRows functions. The DBTable interface has 3 methods
* PK()*int64 Should returns a pointer to the ID field.
* TableName()string Should return the database table name.
* Fields()[]string Should return a slice of field names to be processed.

## Query Builder
The QB type is the core of the package. It can be created as 
````qb:=qb.NewQB("table name")````
Once created, modify the query using methods like Where, Order, Fields, Limit etc. Then the qb variable can be provided to the ListRows function for fetching rows.

## jmoiron/sqlx
This package is developed and tested against github.com/jmoiron/sqlx package and using MySql database. The sqlx database can be provided where DBGetter, DBSelector, DBExecer interfaces are  expected.
