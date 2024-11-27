package model

import "strconv"

type User struct {
	Index      string `json:"-"`
	EmployeeID string `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Title      string `json:"title"`
	Email      string `json:"email"`
	StartDate  string `json:"start_date"`
}

func UserFromList(rows []string) User {
	return User{
		Index:      rows[0],
		LastName:   rows[1],
		FirstName:  rows[2],
		Title:      rows[3],
		Email:      rows[4],
		EmployeeID: rows[5],
		StartDate:  rows[6],
	}
}

func ListFromUser(data User) []string {
	return []string{
		data.Index,
		data.LastName,
		data.FirstName,
		data.Title,
		data.Email,
		data.EmployeeID,
		data.StartDate,
	}
}

func UsersFromList(rows [][]string) []User {
	result := make([]User, 0)
	for _, row := range rows {
		result = append(result, User{
			LastName:   row[1],
			FirstName:  row[2],
			Title:      row[3],
			Email:      row[4],
			EmployeeID: row[5],
			StartDate:  row[6],
		})
	}
	return result
}

func ListFromUsers(data []User) [][]string {
	result := make([][]string, 0)
	for idx, user := range data {
		result = append(result, []string{
			strconv.Itoa(idx + 1),
			user.LastName,
			user.FirstName,
			user.Title,
			user.Email,
			user.EmployeeID,
			user.StartDate,
		})
	}
	return result
}
